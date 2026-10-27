package main

import (
	"Build_my_own_redis/internal/command"
	"Build_my_own_redis/internal/protocol"
	"Build_my_own_redis/io_multiplexing"
	"bufio"
	"bytes"
	"fmt"
	"net"
	"syscall"
)

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	reader := bufio.NewReader(conn)
// 	for {
// 		args, err := protocol.ReadCommand(reader)
// 		if err != nil {
// 			if err == io.EOF {
// 				fmt.Println("Client disconnected")
// 			}
// 			return
// 		}
// 		if len(args) == 0 {
// 			continue
// 		}
// 		if _, err := conn.Write(command.Handle(args)); err != nil {
// 			fmt.Println("Error writing response:", err)
// 			return
// 		}
// 	}
// }

func handleConnectionFD(connFd int) {
	buf := make([]byte, 4096) // 4KB buffer
	// Read data from the connection
	n, err := syscall.Read(connFd, buf)
	// Handle the case where the client has closed the connection
	if err != nil {
		fmt.Println("read error:", err)
		syscall.Close(connFd)
		return
	}

	if n == 0 {
		fmt.Println("Client disconnected")
		syscall.Close(connFd)
		return
	}

	// Process the data (for demonstration, we just print it)
	reader := bufio.NewReader(bytes.NewReader(buf[:n]))
	for {
		// Read and handle commands from the buffer
		args, err := protocol.ReadCommand(reader)
		if err != nil {
			return
		}
		// If no arguments were read, continue to the next iteration
		if len(args) == 0 {
			continue
		}
		// Handle the command and send a response
		reply := command.Handle(args)
		if _, err := syscall.Write(connFd, reply); err != nil {
			fmt.Println("write error:", err)
			syscall.Close(connFd)
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080")

	// Use type assertion to get the underlying *net.TCPListener
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		fmt.Println("Failed to cast listener to TCPListener")
		return
	}

	// Get the file descriptor of the listener
	listenerFile, err := tcpListener.File()
	if err != nil {
		fmt.Println("Error getting listener file:", err)
		return
	}
	defer listenerFile.Close()

	severFd := int(listenerFile.Fd())

	//create an I/O multiplexer (epoll) to handle multiple connections
	multiplexer, err := io_multiplexing.CreateIOMultiplexer()
	if err != nil {
		fmt.Println("Error creating I/O multiplexer:", err)
		return
	}

	defer multiplexer.Close()

	if err := multiplexer.Add(io_multiplexing.Event{FD: severFd}); err != nil {
		fmt.Println("Error adding listener to multiplexer:", err)
		return
	}

	for {
		events, err := multiplexer.Wait()
		if err != nil {
			fmt.Println("Error waiting for events:", err)
			continue
		}

		for _, event := range events {
			// Check if the event is for the listener (new connection)
			if event.FD == severFd {
				// Accept the new connection
				connFD, _, err := syscall.Accept(severFd)
				if err != nil {
					fmt.Println("Error accepting connection:", err)
					continue
				}
				// Add the new connection to the multiplexer
				if err := multiplexer.Add(io_multiplexing.Event{FD: connFD}); err != nil {
					fmt.Println("Error adding connection to multiplexer:", err)
					syscall.Close(connFD)
				}
				continue
			}
			// Handle the event for an existing connection
			handleConnectionFD(event.FD)
		}
	}
	// threadPoolSize := 5                        // Set the desired thread pool size
	// sem := make(chan struct{}, threadPoolSize) // Create a semaphore channel to limit the number of concurrent goroutines

	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		fmt.Println("Error accepting connection:", err)
	// 		continue
	// 	}
	// 	// go handleConnection(conn) // 1 connect equals 1 goroutine, so we can handle multiple connections concurrently
	// 	sem <- struct{}{} // Acquire a slot in the semaphore
	// 	go func() {
	// 		handleConnection(conn)
	// 		<-sem // Release the slot in the semaphore
	// 	}()
	// }

}
