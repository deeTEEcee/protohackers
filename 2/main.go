package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
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
	var smallest *int
	for lo <= hi {
		mid := (hi + lo) / 2
		val := f(arr[mid])
		if val > target {
			hi = mid - 1
		} else if val < target {
			smallest = &mid
			lo = mid + 1
		} else {
			return mid
		}
	}
	if smallest == nil {
		return 0
	} else {
		return *smallest
	}
}

func findMaxIndex[T any](arr []T, target int, f func(T) int) int {
	lo := 0
	hi := len(arr) - 1
	var largest *int
	for lo <= hi {
		mid := (hi + lo) / 2
		val := f(arr[mid])
		if val > target {
			largest = &mid
			hi = mid - 1
		} else if val < target {
			lo = mid + 1
		} else {
			return mid
		}
	}
	if largest == nil {
		return len(arr) - 1
	} else {
		return *largest
	}
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
	return arr[i]
}

/*
Just use a sorted array, append, and try sort.searchInts to get something initially working. We can
research performance after.
*/
func insert(store *[]InsertMessage, timestamp int32, price int32) {
	msg := InsertMessage{timestamp, price}
	*store = append(*store, msg)
}

func query(store []InsertMessage, minTs int32, maxTs int32) float32 {
	// Find the average price
	//getter := func(x InsertMessage) int { return int(x.Timestamp) }
	//minIndex := findMinIndex(store, int(minTs), getter)
	//maxIndex := findMaxIndex(store, int(maxTs), getter)

	return 0.0
}

func handleConnection(conn net.Conn) {
	defer func() {
		fmt.Println("Closing connection")
		conn.Close()
	}()

	buf := make([]byte, 9)
	store := make([]InsertMessage, 0)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		if n == 0 {
			return
		}
		fmt.Printf("Received: %s\n", buf[:n])

		switch buf[0] {
		case 'I':
			var msg InsertMessage
			err = binary.Read(bytes.NewBuffer(buf[1:]), binary.BigEndian, &msg)
			insert(&store, msg.Timestamp, msg.Price)
		case 'Q':
			var msg QueryMessage
			var avg float32
			err = binary.Read(bytes.NewBuffer(buf[1:]), binary.BigEndian, &msg)
			avg = query(store, msg.MinTimestamp, msg.MaxTimestamp)
			var output string
			if math.Trunc(float64(avg)) == float64(avg) {
				output = strconv.FormatInt(int64(avg), 10)
			} else {
				output = strconv.FormatFloat(float64(avg), 'f', -1, 32)
			}
			_, err = conn.Write([]byte(output))
			if err != nil {
				fmt.Println("Error writing:", err)
				return
			}
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
