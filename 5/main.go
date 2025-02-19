package main

import (
	"fmt"
	"log"
	"net"
	"os"
	. "protohackers/tcp"
	"time"
)

func handleClient(client net.Conn, upstream net.Conn) {
	defer func() {
		log.Println("Closing connection")
		client.Close()
	}()
	// 1. Initialize the chatroom name first
	initServerMsg := ReadMessage(upstream)
	WriteMessage(client, initServerMsg)

	initClientResponse := ReadMessage(client)
	WriteMessage(upstream, initClientResponse)

	//serverMsg := ReadMessage(upstream)
	//WriteMessage(client, serverMsg)

	// 2. After step 1, we just have to send every upstream message back down to the client
	go handleUpstream(client, upstream)

	for {
		// 3. Client can handle sending messages upstream here
		success := runOnce(client, upstream)
		if !success {
			return
		}
	}
}

func handleUpstream(client net.Conn, upstream net.Conn) {
	// We don't know when to expect upstream messages. So this has to be handled
	// separately from client messages
	message := ReadMessage(upstream)
	// TODO: Filter message here and replace bitcoin
	WriteMessage(client, message)
	time.Sleep(100 * time.Millisecond)
}

func runOnce(client net.Conn, upstream net.Conn) bool {
	message := ReadMessage(client)
	// TODO: Filter message here and replace bitcoin
	WriteMessage(upstream, message)
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

	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		upstream := SetupUpstream("chat.protohackers.com:16963")
		go handleClient(client, upstream)
	}
}

func main() {
	startServer(":8080")
}
