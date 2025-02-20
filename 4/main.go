package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"protohackers/parse"
	"strings"
	"unicode"
)

type KeyStore struct {
	store map[string]string
}

func (ks *KeyStore) Get(k string) string {
	return ks.store[k]
}

func (ks *KeyStore) Put(k string, v string) {
	ks.store[k] = v
}

func handleClient(connection *net.UDPConn) {
	store := KeyStore{make(map[string]string)}
	defer func() {
		log.Println("Ending connection")
		err := connection.Close()
		if err != nil {
			log.Printf("Error occurred during close: %s\n", err)
		}
	}()
	buffer := make([]byte, 1024)
	store.Put("version", "Ken's Key-Value Store 1.0")
	for {
		runOnce(connection, &buffer, &store)
		//if !success {
		//	return
		//}
	}
}

func runOnce(conn *net.UDPConn, buffer *[]byte, store *KeyStore) bool {
	n, clientAddr, err := conn.ReadFromUDP(*buffer)
	if err != nil {
		log.Println("Error reading from UDP:", err)
		return false
	}

	message := string((*buffer)[:n])
	if strings.Contains(message, "\n") {
		message = strings.TrimRightFunc(message, unicode.IsSpace)
	}
	key, value, isInsert := parse.ParseMessage(message)
	//log.Printf("%s: %s\n", clientAddr, message)
	if isInsert {
		if key != "version" {
			store.Put(key, value)
		}
	} else {
		value = store.Get(key)
		_, err = conn.WriteToUDP([]byte(fmt.Sprintf("%s=%s", key, value)), clientAddr)
		if err != nil {
			log.Println("Error writing to UDP:", err)
			return false
		}
	}
	return true
}

func startServer(address string) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	fmt.Println("Server listening on " + address)
	conn, err := net.ListenUDP("udp", addr)

	for {
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		handleClient(conn)
	}
}

func main() {
	startServer("0.0.0.0:8080")
}
