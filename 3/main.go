package main

import (
	"fmt"
	"log"
	"net"
	"os"
	. "protohackers/chat"
	"protohackers/validation"
	"sync"
)

// TODO: Add unit tests for tcp chat room

func handleConnection(server *Server, client *Client) {
	defer func() {
		log.Println("Closing connection")
		client.Connection.Close()
		server.DeregisterUser(client)
	}()
	for {
		success := runOnce(server, client)
		if !success {
			return
		}
	}
}

func runOnce(server *Server, client *Client) bool {
	if !client.HasName() {
		server.Send(client, "Welcome to budgetchat! What shall I call you?\n")
		response := server.Wait(client)
		if validation.ValidateName(response) {
			server.RegisterUser(client, response)
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
	chatMessage := fmt.Sprintf("[%s] %s\n", client.Name, message)
	server.Publish(chatMessage, client)
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

	server := Server{Mu: &sync.Mutex{}}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(&server, &Client{Connection: conn})
	}
}

func main() {
	startServer(":8080")
}
