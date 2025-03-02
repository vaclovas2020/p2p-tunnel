package p2p

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"os"
)

func SendMessageToServer(host string, port int, message string) {
	// Load client's certificate and key
	clientCert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Fatal("Error loading client certificate:", err)
	}

	// Load CA certificate to verify the server's certificate
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Fatal("Error reading CA certificate:", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatal("Failed to append CA certificate")
	}

	// Secure TLS configuration
	config := &tls.Config{
		Certificates:       []tls.Certificate{clientCert}, // Use client's certificate
		RootCAs:            caCertPool,                    // Verify server certificate
		InsecureSkipVerify: false,                         // Ensure TLS verification is enabled
	}

	// Establish a secure TLS connection
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)

	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}

	defer func() {
		logClient(conn, "Secure Connection closed")

		conn.Close()
	}()

	logClient(conn, "New secure connection established")

	err = sendMessageClient(conn, message)

	if err != nil {
		logClient(conn, fmt.Sprintln("Write error:", err))

		return
	}

	_, err = receiveMessageClient(conn)

	if err == io.EOF {
		logClient(conn, "Connection closed by the server (EOF detected)")

		return
	}

	if err != nil {
		logClient(conn, fmt.Sprintln("Read error:", err))

		return
	}
}

func receiveMessageClient(conn *tls.Conn) (string, error) {
	resbuff := make([]byte, 1024)

	n, err := conn.Read(resbuff)

	if err != nil {
		return "", err
	}

	resStr := string(resbuff[:n])

	logClient(conn, fmt.Sprintln("Received message:", resStr))

	return resStr, nil
}

func sendMessageClient(conn *tls.Conn, message string) error {
	// Send message
	_, err := conn.Write([]byte(message))

	if err != nil {
		return err
	}

	logClient(conn, fmt.Sprintln("Message sent:", message))

	return nil
}
