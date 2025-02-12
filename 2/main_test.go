package main

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	store := make([]InsertMessage, 0)
	assert.Equal(t, 0.0, query(store, 0, 20))
	insert(&store, 0, 100)
	assert.Equal(t, 100.0, query(store, 0, 20))
	insert(&store, 5, 20)
	assert.Equal(t, 60.0, query(store, 0, 20))
	assert.Equal(t, 0.0, query(store, 10, 20))

	// Bad input values
	assert.Equal(t, 0.0, query(store, 12288, 0))
}

func TestProcessUnsorted(t *testing.T) {
	store := make([]InsertMessage, 0)
	insert(&store, 0, 100)
	insert(&store, 5, 20)
	insert(&store, -2, 1)
	insert(&store, -4, 5)
	insert(&store, 10, 50)
	assert.InDelta(t, 56.66, query(store, 0, 20), 0.1)
}

func TestFindMinMax(t *testing.T) {
	arr := []int{1, 5, 10, 18, 21}
	assert.Equal(t, 1, findMinInt(arr, 0))
	assert.Equal(t, 1, findMinInt(arr, 1))
	assert.Equal(t, 1, findMinInt(arr, 3))
	assert.Equal(t, 5, findMinInt(arr, 5))
	assert.Equal(t, 5, findMinInt(arr, 6))
	assert.Equal(t, 18, findMinInt(arr, 18))
	assert.Equal(t, 18, findMinInt(arr, 20))

	assert.Equal(t, 1, findMaxInt(arr, 0))
	assert.Equal(t, 1, findMaxInt(arr, 1))
	assert.Equal(t, 18, findMaxInt(arr, 11))
	assert.Equal(t, 18, findMaxInt(arr, 18))
	assert.Equal(t, 21, findMaxInt(arr, 19))
	assert.Equal(t, 21, findMaxInt(arr, 21))
	assert.Equal(t, 21, findMaxInt(arr, 24))
}

func insertMsg(ts int, price int) []byte {
	buf := [9]byte{}
	buf[0] = 'I'
	binary.BigEndian.PutUint32(buf[1:5], uint32(ts))
	binary.BigEndian.PutUint32(buf[5:], uint32(price))
	return buf[:]
}

func queryMsg(start int, end int) []byte {
	buf := [9]byte{}
	buf[0] = 'Q'
	binary.BigEndian.PutUint32(buf[1:5], uint32(start))
	binary.BigEndian.PutUint32(buf[5:], uint32(end))
	return buf[:]
}

//type MockConn struct {
//}
//
//func (m MockConn) Read(b []byte) (n int, err error) {
//
//	return 9, nil
//}
//
//func (m MockConn) Write(b []byte) (n int, err error) {
//	return 0, nil
//}

// Method 1: Mock connection interface and write in a unit test manner.
// Method 2: No mocks. Run the server and then send messages through the client.
// If we used `conn.Read(buf)` instead of `io.ReadFull(conn, buf)`, this
// test will timeout since it would hang.
func TestClient(t *testing.T) {
	timeout := time.After(5 * time.Second)
	done := make(chan bool)
	go func() {
		address := "localhost:8081"
		go startServer(address)
		conn, _ := net.Dial("tcp", address)
		buf := make([]byte, 4)
		randBytes := make([]byte, 9)
		var err error
		conn.Write(insertMsg(0, 0))
		conn.Write(insertMsg(5, 10))
		rand.Read(randBytes)
		randBytes[0] = 'X'    // Ensure this doesnt get processed
		conn.Write(randBytes) // Add random bytes.
		conn.Write(queryMsg(0, 10))
		_, err = io.ReadFull(conn, buf)
		assert.Equal(t, nil, err)
		assert.Equal(t, int(binary.BigEndian.Uint32(buf)), 5)

		// Do partial insert in multiple writes
		partialInsert := insertMsg(10, 20)
		conn.Write(partialInsert[:5])
		time.Sleep(1 * time.Second)
		conn.Write(partialInsert[5:])
		conn.Write(queryMsg(0, 10))
		_, err = io.ReadFull(conn, buf)
		assert.Equal(t, nil, err)
		assert.Equal(t, int(binary.BigEndian.Uint32(buf)), 10)
		done <- true
	}()
	select {
	case <-timeout:
		t.Fatal("Timeout (possibly due to partial insert byte read issues)")
	case <-done:
	}
}
