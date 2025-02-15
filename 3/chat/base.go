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

func (c *Client) HasName() bool {
	return c.Name != ""
}

func (c *Client) Send(s *Server, message string) {

}

type Server struct {
	Clients []*Client
	Mu      *sync.Mutex
}

func (s *Server) AddClient(c *Client) {
	s.Clients = append(s.Clients, c)
}

func (s *Server) Publish(message string, exclude *Client) {
	go func() {
		log.Println("Running publish")
		for _, client := range s.Clients {
			if exclude == nil || client != exclude {
				log.Println("Sending ", message, " to ", client.Name)
				s.Send(client, message)
			}
		}
	}()
}

func (s *Server) RegisterUser(client *Client, name string) {
	log.Printf("Registering user %s", name)
	client.Name = name
	s.Mu.Lock()
	s.AddClient(client)
	clientNames := make([]string, 0)
	for _, c := range s.Clients {
		clientNames = append(clientNames, c.Name)
	}
	s.Mu.Unlock()
	go func() {
		if len(clientNames) > 0 {
			s.Send(client, fmt.Sprintf("* The room contains: %s\n", strings.Join(clientNames, ", ")))
		}
		s.Publish(fmt.Sprintf("* %s has entered the room\n", client.Name), client)
	}()
}

func (s *Server) DeregisterUser(client *Client) {
	go func() {
		s.Mu.Lock()
		removeIndex := -1
		for i, c := range s.Clients {
			if c == client {
				removeIndex = i
			}
		}
		if removeIndex == -1 {
			panic(fmt.Sprintf("The client '%s' was not found", client.Name))
		}
		log.Printf("Dergistering user %s", client.Name)
		log.Println("Removing", client.Name)
		slices.Delete(s.Clients, removeIndex, removeIndex+1)
		s.Mu.Unlock()
	}()
	go func() {
		s.Publish(fmt.Sprintf("* %s has left the room", client.Name), client)
	}()
}

func (s *Server) Send(client *Client, message string) {
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
