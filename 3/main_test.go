package main

import (
	"github.com/stretchr/testify/assert"
	. "protohackers/chat"
	"sync"
	"testing"
)

func TestAddClient(t *testing.T) {
	server := Server{Mu: &sync.Mutex{}}
	client := Client{}
	server.AddClient(&client)
	assert.Equal(t, len(server.Clients), 1)
}
