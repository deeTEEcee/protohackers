package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"protohackers/parse"
	"strings"
	"time"
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
	start := time.Now()
	//t := time.Now()
	n, clientAddr, err := conn.ReadFrom(*buffer)
	if err != nil {
		log.Println("Error reading from UDP:", err)
		return false
	}
	//log.Printf("Time to read: %s", time.Since(t))

	message := string((*buffer)[:n])
	if strings.Contains(message, "\n") {
		message = strings.TrimRightFunc(message, unicode.IsSpace)
	}
	key, value, isInsert := parse.ParseMessage(message)
	//log.Printf("%s: %s\n", clientAddr, message)
	if isInsert {
		if key != "version" {
			//t = time.Now()
			store.Put(key, value)
			//log.Printf("Time to put: %s", time.Since(t))
		}
	} else {
		//t = time.Now()
		value = store.Get(key)
		//log.Printf("Time to get: %s", time.Since(t))
		//t = time.Now()
		_, err = conn.WriteTo([]byte(fmt.Sprintf("%s=%s", key, value)), clientAddr)
		//log.Printf("Time to write udp: %s", time.Since(t))
		if err != nil {
			log.Println("Error writing to UDP:", err)
			return false
		}
	}
	log.Printf("Total time: %s", time.Since(start))
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
	if err != nil {
		fmt.Println("Error accepting:", err)
		return
	}
	// This is timing out due to 3 seconds... suspicious?
	//err = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	//if err != nil {
	//	fmt.Println("Error configuring:", err)
	//	return
	//}

	for {
		handleClient(conn)
	}
}

func main() {
	startServer("0.0.0.0:8080")
}
