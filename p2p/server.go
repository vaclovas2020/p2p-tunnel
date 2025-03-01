package p2p

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
)

func StartServer(port int) {
	// Load server certificate
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Println("Error loading server certificate:", err)
		os.Exit(1)
	}

	// Load CA certificate to verify client certificates
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		fmt.Println("Error reading CA certificate:", err)
		os.Exit(1)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure TLS
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool, // Used to verify client certificates
		ClientAuth:   tls.RequireAndVerifyClientCert,
		GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			// Return the server's certificate
			return &cert, nil
		},
	}

	ln, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), config)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Printf("Listening on port %d...\n", port)

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

	reqbuff := make([]byte, 1024)
	n, err := conn.Read(reqbuff)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	fmt.Println("Received message:", string(reqbuff[:n]))

	resbuff := make([]byte, 1024)
	n, err = conn.Write(resbuff)
	if err != nil {
		fmt.Println("Error writing data:", err)
		return
	}

	fmt.Println("Writed message:", string(resbuff[:n]))
}
