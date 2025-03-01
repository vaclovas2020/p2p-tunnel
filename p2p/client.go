package p2p

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
)

func SendMessageToPeer(host string, port int, message string) {
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

	// Send message
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatal("Error sending message:", err)
	}

	fmt.Println("Message sent:", message)

	// Received message
	resbuff := make([]byte, 1024)
	n, err := conn.Read(resbuff)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	fmt.Println("Received message:", string(resbuff[:n]))
}
