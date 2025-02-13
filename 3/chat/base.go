package chat

import (
	"bufio"
	"log"
	"net"
	"strings"
	"unicode"
)

type Client struct {
	Connection net.Conn
	// TODO: Learn mutxes and create an autoid incrementer (https://gobyexample.com/mutexes)
	// https://stackoverflow.com/questions/64631848/how-to-create-an-autoincrement-id-field
	//id
	Name string
}

func (c Client) HasName() bool {
	return c.Name != ""
}

func (c Client) Send(s Server, message string) {

}

type Server struct {
	clients []*Client
}

func (s Server) AddClient(c *Client) {
	s.clients = append(s.clients, c)
}

func (s Server) Publish(message string) {
	for _, client := range s.clients {
		s.Send(client, message)
	}
}

func (s Server) Send(client *Client, message string) {
	_, err := client.Connection.Write([]byte(message))
	if err != nil {
		log.Printf("Error occurred while sending: %s\n", err)
		return
	}
}

func (s Server) Wait(client *Client) string {
	reader := bufio.NewReader(client.Connection)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error occurred while reading from client: %s\n", err)
		return ""
	}
	message = strings.TrimRightFunc(message, unicode.IsSpace)
	return message
}
