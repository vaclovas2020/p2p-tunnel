package p2p

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func StartServer(port int) {
	// Load server certificate
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println("Error loading server certificate:", err)
		os.Exit(1)
	}

	// Load CA certificate to verify client certificates
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Println("Error reading CA certificate:", err)
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
		log.Println("Error listening:", err)
		os.Exit(1)
	}

	defer ln.Close()

	log.Printf("Listening on port %d...\n", port)

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		logServer(conn, "Secure Connection closed")

		conn.Close()
	}()

	logServer(conn, "New secure connection established")

	for {
		req, err := receiveMessageServer(conn)

		if err == io.EOF {
			logServer(conn, "Connection closed by the client (EOF detected)")
			return
		}

		if err != nil {
			logServer(conn, fmt.Sprintln("Read error:", err))

			return
		}

		err = sendMessageServer(conn, fmt.Sprintf("echo %s", req))

		if err != nil {
			logServer(conn, fmt.Sprintln("Write error:", err))

			return
		}
	}
}

func receiveMessageServer(conn net.Conn) (string, error) {
	reqbuff := make([]byte, 1024)
	n, err := conn.Read(reqbuff)
	if err != nil {
		return "", err
	}

	reqStr := string(reqbuff[:n])

	logServer(conn, fmt.Sprintln("Received message:", reqStr))

	return reqStr, nil
}

func sendMessageServer(conn net.Conn, message string) error {
	// Send message
	_, err := conn.Write([]byte(message))

	if err != nil {
		return err
	}

	logServer(conn, fmt.Sprintln("Message sent:", message))

	return nil
}
