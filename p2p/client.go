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
		log.Fatal("Error connecting to p2p peer:", err)
	}

	defer conn.Close()

	for i := range 5 {
		sendMessageClient(conn, fmt.Sprintf("%d:%s", i, message))
		req, _ := receiveMessageClient(conn)
		message = req
	}
}

func receiveMessageClient(conn *tls.Conn) (string, error) {
	reqbuff := make([]byte, 1024)
	n, err := conn.Read(reqbuff)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Connection closed by the server (EOF detected)")
			return "", nil // Return empty string and nil error to indicate EOF
		}
		fmt.Println("Error reading data:", err)
		return "", err
	}

	reqStr := string(reqbuff[:n])
	fmt.Println("Received message:", reqStr)

	return reqStr, nil
}

func sendMessageClient(conn *tls.Conn, message string) {
	// Send message
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Fatal("Error sending message:", err)
	}

	fmt.Println("Message sent:", message)
}
