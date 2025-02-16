package main

import (
	"fmt"
	"log"
	"net"
	"os"
	. "protohackers/chat"
	"protohackers/validation"
	"strings"
)

// TODO: Add unit tests for tcp chat room

func handleClient(server *Server, client *Client) {
	defer func() {
		log.Println("Closing connection")
		client.Connection.Close()
		server.Exit <- client
	}()
	for {
		success := runOnce(server, client)
		if !success {
			return
		}
	}
}

func handleChatroom(s *Server) {
	// Handles messages sent from client to the server
	// TODO: Odd that we do publish here and 1:1 messages elsewhere but we can handle that later.
	// Just know that by doing group management here, we don't run into lock or synch issues
	// with data
	for {
		select {
		case client := <-s.Enter:
			log.Printf("Registering user %s", client.Name)
			s.AddClient(client)
			clientNames := make([]string, 0)
			for _, c := range s.Clients {
				if c != client {
					clientNames = append(clientNames, c.Name)
				}
			}
			s.Send(client, fmt.Sprintf("* The room contains: %s\n", strings.Join(clientNames, ", ")))
			s.Publish(fmt.Sprintf("* %s has entered the room\n", client.Name), client)
		case client := <-s.Exit:
			log.Printf("Deregistering user %s", client.Name)
			// Unregister the user and send messages
			s.RemoveClient(client)
			s.Publish(fmt.Sprintf("* %s has left the room\n", client.Name), client)
		case msg := <-s.Msg:
			log.Println("Running publish")
			client := msg.Sender
			s.Publish(fmt.Sprintf("[%s] %s\n", client.Name, msg.Message), client)
		}
	}
}

func runOnce(server *Server, client *Client) bool {
	if !client.HasName() {
		// TODO: This part is actually a bit. We tried limiting server <-> client behavior
		// to messages but initially, the client is still handling server work.
		server.Send(client, "Welcome to budgetchat! What shall I call you?\n")
		response := server.Wait(client)
		if validation.ValidateName(response) {
			client.Name = response
			server.Enter <- client
			return true
		} else {
			server.Send(client, fmt.Sprintf("Name '%s' is invalid. Disconnecting\n", response))
			return false
		}
	}
	message := server.Wait(client)
	if message == "" {
		return false
	}
	server.Msg <- Message{Sender: client, Message: message}
	return true
}

func startServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listening on :8080")

	server := Server{
		Msg:   make(chan Message),
		Enter: make(chan *Client),
		Exit:  make(chan *Client),
	}
	go handleChatroom(&server)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleClient(&server, &Client{Connection: conn})
	}
}

func main() {
	startServer(":8080")
}
