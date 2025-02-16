package chat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"slices"
	"strings"
	"sync"
	"unicode"
)

type Client struct {
	Connection net.Conn
	Name       string
}

type Message struct {
	Sender  *Client
	Message string
}

func (c *Client) HasName() bool {
	return c.Name != ""
}

func (c *Client) Send(s *Server, message string) {

}

type Server struct {
	Clients []*Client
	Mu      *sync.Mutex

	// Channels
	Msg chan Message // A message we need to share to all clients other than the sender
	// Tells us when a client enters or exits
	Enter chan *Client
	Exit  chan *Client
}

func (s *Server) AddClient(c *Client) {
	s.Clients = append(s.Clients, c)
}

func (s *Server) Publish(message string, exclude *Client) {
	for _, client := range s.Clients {
		if exclude == nil || client != exclude {
			//log.Println("Sending ", message, " to ", client.Name)
			s.Send(client, message)
		}
	}
}

func (s *Server) RemoveClient(client *Client) {
	removeIndex := -1
	for i, c := range s.Clients {
		if c == client {
			removeIndex = i
		}
	}
	if removeIndex == -1 {
		panic(fmt.Sprintf("The client '%s' was not found", client.Name))
	}
	s.Clients = slices.Delete(s.Clients, removeIndex, removeIndex+1)
}

func (s *Server) Send(client *Client, message string) {
	log.Printf("%v\n", client)
	_, err := client.Connection.Write([]byte(message))
	if err != nil {
		log.Printf("Error occurred while sending: %s\n", err)
		return
	}
}

func (s *Server) Wait(client *Client) string {
	reader := bufio.NewReader(client.Connection)
	message, err := reader.ReadString('\n')
	if err != nil {
		// This can occur with EOF when the client disconnects intentionally or not.
		log.Printf("Error occurred while reading from client: %s\n", err)
		return ""
	}
	message = strings.TrimRightFunc(message, unicode.IsSpace)
	return message
}
