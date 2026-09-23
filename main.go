package main

import (
	"Build_my_own_redis/internal/command"
	"Build_my_own_redis/internal/protocol"
	"bufio"
	"fmt"
	"io"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		args, err := protocol.ReadCommand(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
			}
			return
		}
		if len(args) == 0 {
			continue
		}
		if _, err := conn.Write(command.Handle(args)); err != nil {
			fmt.Println("Error writing response:", err)
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
