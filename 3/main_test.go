package main

import (
	"github.com/stretchr/testify/assert"
	. "protohackers/chat"
	"sync"
	"testing"
	"time"
)

func TestServerBasics(t *testing.T) {
	server := Server{Mu: &sync.Mutex{}}
	client := Client{}
	server.AddClient(&client)
	assert.Equal(t, len(server.Clients), 1)
	// TODO: better way to handle these and force goroutine functions to run
	// sync?
	server.DeregisterUser(&client)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, len(server.Clients), 0)

}
