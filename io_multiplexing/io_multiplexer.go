package io_multiplexing

const MaxEvents = 1024

type Event struct {
	FD     int
	Events uint32
}

type IOEventMultiplexer interface {
	// Add adds a file descriptor to the multiplexer with the specified events.
	Add(fd int, events uint32) error
	// Remove removes a file descriptor from the multiplexer.
	Remove(fd int) error
	// Wait waits for events on the file descriptors and returns the triggered events.
	Wait() ([]Event, error)
}
