package tcp

import (
	"bufio"
	"log"
	"net"
)

func SetupUpstream(addr string) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Dial failed: %s\n", err.Error())
	}
	return conn
}

func ReadMessage(conn net.Conn) string {
	// We can use this for reading messages from any connection
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		// This can occur with EOF when the client disconnects intentionally or not.
		log.Printf("Error occurred while reading from client: %s\n", err)
		return ""
	}
	log.Printf("Read message '%s' from (%s)\n", message, conn.LocalAddr())
	return message
}

func WriteMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Printf("Error occurred while sending: %s\n", err)
		return
	}
	log.Printf("Writing '%s' to (%s)\n", message, conn.LocalAddr())
}
