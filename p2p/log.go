package p2p

import (
	"crypto/tls"
	"log"
	"net"
)

func logServer(conn net.Conn, message string) {
	log.Printf("%s (remote: %s, local: %s)",
		message,
		conn.RemoteAddr().String(),
		conn.LocalAddr().String())
}

func logClient(conn *tls.Conn, message string) {
	log.Printf("%s (remote: %s, local: %s)",
		message,
		conn.RemoteAddr().String(),
		conn.LocalAddr().String())
}
