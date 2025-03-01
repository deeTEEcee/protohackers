package main

import (
	"fmt"
	"log"
	"net"
	"os"
	. "protohackers/tcp"
	"strings"
)

func handleClient(client net.Conn, upstream net.Conn) {
	defer func() {
		log.Println("Closing connection client")
		client.Close()
		upstream.Close()
	}()
	// 1. Initialize the chatroom name first
	initServerMsg := ReadMessage(upstream)
	WriteMessage(client, initServerMsg)

	initClientResponse := ReadMessage(client)
	WriteMessage(upstream, initClientResponse)

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
	defer func() {
		log.Println("Closing connection upstream")
		client.Close()
		upstream.Close()
	}()

	for {
		message := ReadMessage(upstream)
		if message == "" {
			return
		}
		// Not sure why this is needed but I guess sometimes as a malicious proxy you only
		// proxy client OR server side?
		message = strings.TrimRight(message, "\n")
		message = Rewrite(message)
		message += "\n"
		err := WriteMessage(client, message)
		if err != nil {
			return
		}
	}
}

func runOnce(client net.Conn, upstream net.Conn) bool {
	message := ReadMessage(client)
	if message == "" {
		return false
	}
	// We rewrite messages that we send upstream to the client.
	// When `handleUpstream` happens, that message will be the rewritten message.
	message = strings.TrimRight(message, "\n")
	message = Rewrite(message)
	message += "\n"
	err := WriteMessage(upstream, message)
	if err != nil {
		return false
	}
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
