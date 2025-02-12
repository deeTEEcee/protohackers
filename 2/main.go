package main

import (
	"cmp"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
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
var getter = func(x InsertMessage) int { return int(x.Timestamp) }
var compare = func(a, b InsertMessage) int {
	return cmp.Compare(getter(a), getter(b))
}

func insert(store *[]InsertMessage, timestamp int32, price int32) {
	fmt.Printf("Inserting %d %d\n", timestamp, price)
	msg := InsertMessage{timestamp, price}
	*store = append(*store, msg)
}

func query(store []InsertMessage, minTs int32, maxTs int32) float64 {
	// Do lazy sorting to improve on insert performance.
	slices.SortFunc(store, compare)
	fmt.Printf("Querying %d %d\n", minTs, maxTs)
	// Find the average price
	if len(store) == 0 || minTs > maxTs {
		return 0.0
	}
	minIndex := findMinIndex(store, int(minTs), getter)
	maxIndex := findMaxIndex(store, int(maxTs), getter)
	var total int32 = 0
	var count int32 = 0
	for i := minIndex; i <= maxIndex; i++ {
		// Need this duplicate condition check because findIndex gets closest numbers,
		// allowing for flexible queries.
		if store[i].Timestamp >= minTs && store[i].Timestamp <= maxTs {
			total += store[i].Price
			count += 1
		}
	}
	if count == 0 {
		return 0.0
	}
	fmt.Printf("Size of store: %d, Calculation: %d/%d\n", len(store), total, count)
	return float64(total) / float64(count)
}

func handleConnection(conn net.Conn) {
	defer func() {
		fmt.Println("Closing connection")
		conn.Close()
	}()

	buf := make([]byte, 9)
	store := make([]InsertMessage, 0)
	for {
		// TODO: Create a local test for tcp server to deal with difference between io.ReadFull
		// and conn.Read(buf)
		_, err := io.ReadFull(conn, buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		fieldOne := int32(binary.BigEndian.Uint32(buf[1:5]))
		fieldTwo := int32(binary.BigEndian.Uint32(buf[5:]))

		switch buf[0] {
		case 'I':
			//var msg InsertMessage
			//err = binary.Read(bytes.NewBuffer(buf[1:]), binary.BigEndian, &msg)
			//if err != nil {
			//	fmt.Println("Error doing binary read:", err)
			//}
			insert(&store, fieldOne, fieldTwo)
		case 'Q':
			//var msg QueryMessage
			var avg float64
			//err = binary.Read(bytes.NewBuffer(buf[1:]), binary.BigEndian, &msg)
			//if err != nil {
			//	fmt.Println("Error doing binary read:", err)
			//}
			avg = query(store, fieldOne, fieldTwo)
			fmt.Printf("Sending %d\n", int32(avg))

			// Make sure we send 4 bytes.
			// Unsigned int32 is same binary rep as uint32?
			output := make([]byte, 4)
			binary.BigEndian.PutUint32(output, uint32(avg))
			_, err = conn.Write(output)
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
