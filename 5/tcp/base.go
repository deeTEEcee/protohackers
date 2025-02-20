package tcp

import (
	"bufio"
	"log"
	"net"
	"regexp"
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
	log.Printf("Read message '%s' from (%s)", message, conn.LocalAddr())
	return message
}

func WriteMessage(conn net.Conn, message string) error {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Printf("Error occurred while sending: %s\n", err)
		return err
	}
	log.Printf("Writing '%s' to (%s)", message, conn.LocalAddr())
	return nil
}

var tonyAddress = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
var re = regexp.MustCompile(`(\s+)7[a-zA-Z0-9]{25,36}|7[a-zA-Z0-9]{25,36}(\s+)`)

func Rewrite(message string) string {
	// Rewrite boguscoin addresses as requested in https://protohackers.com/problem/5
	return re.ReplaceAllString(message, `${1}`+tonyAddress+`${2}`)
}
