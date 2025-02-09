package main

import (
	"bufio"
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
		if err := conn.Close(); err != nil {
			fmt.Println("Error closing")
		}
	}()
	var err error
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		in := scanner.Bytes()
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		fmt.Printf("Received: %s\n", in)
		responseBytes, ok := process(in)
		responseBytes = append(responseBytes, byte('\n'))
		fmt.Printf("About to write: %s\n", responseBytes)
		if _, err = conn.Write(responseBytes); err != nil {
			fmt.Printf("Error writing: %v", err)
		}
		if !ok {
			fmt.Println("Malformed response. Breaking out of infinite loop.")
			break
		}
	}
}

func isPrime(n int) bool {
	// This needs to be optimized. A slow approach may timeout.
	if n <= 3 {
		return n > 1
	} else if n%2 == 0 || n%3 == 0 {
		return false
	}

	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
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
		malformed := []byte("p")
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
