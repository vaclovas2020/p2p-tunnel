package p2p

import (
	"crypto/tls"
	"fmt"
	"log"
)

func SendMessageToPeer(host string, port int, message string) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))

	if err != nil {
		log.Fatal("Error sending message:", err)
	}

	fmt.Println("Message sent:", message)
}
