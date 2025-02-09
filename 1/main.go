package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func isValid(req Request) bool {
	return *req.Method == "isPrime" && req.Number != nil
}

type Request struct {
	Method *string `json:"method"`
	Number *int    `json:"number"`
}

type Response struct {
	Method  string `json:"method"`
	IsPrime bool   `json:"prime"`
}

func handleConnection(conn net.Conn) {
	defer func() {
		fmt.Println("Closing connection")
		conn.Close()
	}()

	buffer := make([]byte, 1024)
	for {
		var err error
		fmt.Println("Reading...")
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		if n == 0 {
			return
		} // {"number":-2,"method":"isPrime"}
		fmt.Printf("Received: %s\n", buffer[:n])
		responseBytes, _ := process(buffer[:n])
		fmt.Printf("About to write: %s\n", responseBytes)
		conn.Write(responseBytes)
		fmt.Println("Finished writing")
	}
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func process(requestBytes []byte) ([]byte, bool) {
	// Process a request and check if it's prime
	var request Request
	err := json.Unmarshal(requestBytes, &request)
	if err == nil && isValid(request) {
		response := Response{Method: "isPrime", IsPrime: isPrime(*request.Number)}
		var responseBytes []byte
		responseBytes, err = json.Marshal(response)
		return responseBytes, true
	} else {
		malformed := make([]byte, 4)
		return malformed, false
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn)
	}
}
