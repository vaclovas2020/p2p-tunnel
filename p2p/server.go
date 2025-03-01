package p2p

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
)

func StartServer(port int) {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Println("Error loading certificate:", err)
		os.Exit(1)
	}

	ln, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Printf("Listening on port %d...", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New secure connection established")

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	fmt.Println("Received message:", string(buffer))
}
