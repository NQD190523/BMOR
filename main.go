package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		fmt.Println("Received:", string(buffer[:n]))
		conn.Write([]byte("Hello from server!\n"))
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

	threadPoolSize := 5                        // Set the desired thread pool size
	sem := make(chan struct{}, threadPoolSize) // Create a semaphore channel to limit the number of concurrent goroutines

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// go handleConnection(conn) // 1 connect equals 1 goroutine, so we can handle multiple connections concurrently
		sem <- struct{}{} // Acquire a slot in the semaphore
		go func() {
			handleConnection(conn)
			<-sem // Release the slot in the semaphore
		}()
	}

}
