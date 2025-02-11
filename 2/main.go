package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type InsertMessage struct {
	Timestamp int32
	Price     int32
}

type QueryMessage struct {
	MinTimestamp int32
	MaxTimestamp int32
}

type MessageType int8

const (
	Insert MessageType = 'I'
	Query  MessageType = 'Q'
)

func findMinIndex[T any](arr []T, target int, f func(T) int) int {
	// Returns target or value less than target
	lo := 0
	hi := len(arr) - 1
	smallest := -1
	for lo <= hi {
		mid := (hi + lo) / 2
		val := f(arr[mid])
		if val > target {
			hi = mid - 1
		} else if val < target {
			if lo > smallest {
				smallest = lo
			}
			lo = mid + 1
		} else {
			return mid
		}
	}
	return smallest
}

func findMaxIndex[T any](arr []T, target int, f func(T) int) int {
	lo := 0
	hi := len(arr) - 1
	largest := len(arr) + 1
	for lo <= hi {
		mid := (hi + lo) / 2
		val := f(arr[mid])
		if val > target {
			if hi < largest {
				largest = hi
			}
			hi = mid - 1
		} else if val < target {
			lo = mid + 1
		} else {
			return mid
		}
	}
	return largest
}

func findMinInt(arr []int, target int) int {
	i := findMinIndex(arr, target, func(x int) int { return x })
	if i == -1 {
		return -1
	} else {
		return arr[i]
	}
}

func findMaxInt(arr []int, target int) int {
	i := findMaxIndex(arr, target, func(x int) int { return x })
	if i > len(arr) {
		return -1
	} else {
		return arr[i]
	}
}

/*
Just use a sorted array, append, and try sort.searchInts to get something initially working. We can
research performance after.
*/
func processInsert(msg InsertMessage) {

}

func processQuery(msg QueryMessage) {

}

func handleConnection(conn net.Conn) {
	defer func() {
		fmt.Println("Closing connection")
		conn.Close()
	}()

	buf := make([]byte, 9)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		if n == 0 {
			return
		}
		switch buf[0] {
		case 'I':
			var msg InsertMessage
			err = binary.Read(bytes.NewBuffer(buf[1:]), binary.BigEndian, &msg)
			processInsert(msg)
		case 'Q':
			var msg QueryMessage
			err = binary.Read(bytes.NewBuffer(buf[1:]), binary.BigEndian, &msg)
			processQuery(msg)
		}

		if err != nil {
			fmt.Printf("Error unpacking %s\n", buf)
		}

		fmt.Printf("Received: %s\n", buf[:n])
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
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
