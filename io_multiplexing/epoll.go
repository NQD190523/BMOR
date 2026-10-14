//go:build linux

package io_multiplexing

import "syscall"

type Epoll struct {
	fileDescriptor int
	epollEvents    []syscall.EpollEvent
	genericEvents  []Event
}

func CreateIOMultiplexer() (*Epoll, error) {

	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &Epoll{
		fileDescriptor: fd,
		epollEvents:    make([]syscall.EpollEvent, 0),
		genericEvents:  make([]Event, 0),
	}, nil
}

func (e *Epoll) Add(event Event) error {
	var events uint32 = syscall.EPOLLIN
	epollEvent := syscall.EpollEvent{
		Events: events,
		Fd:     int32(event.FD),
	}
	// Add the file descriptor to the epoll instance
	return syscall.EpollCtl(e.fileDescriptor, syscall.EPOLL_CTL_ADD, event.FD, &epollEvent)
}

func (e *Epoll) Wait() ([]Event, error) {
	n, err := syscall.EpollWait(e.fileDescriptor, e.epollEvents, -1)
	if err != nil {
		return nil, err
	}
	e.genericEvents = e.genericEvents[:0] // Clear the slice
	for i := 0; i < n; i++ {
		e.genericEvents = append(e.genericEvents, Event{
			FD:     int(e.epollEvents[i].Fd),
			Events: e.epollEvents[i].Events,
		})
	}
	return e.genericEvents, nil
}

func (e *Epoll) Close() error {
	return syscall.Close(e.fileDescriptor)
}
