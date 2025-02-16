package main

import (
	"github.com/stretchr/testify/assert"
	. "protohackers/chat"
	"testing"
)

func TestServerBasics(t *testing.T) {
	server := Server{}
	client := Client{}
	server.AddClient(&client)
	assert.Equal(t, len(server.Clients), 1)
	server.RemoveClient(&client)
	assert.Equal(t, len(server.Clients), 0)
}

func TestServerEnterAndLeave(t *testing.T) {
	server := Server{}
	names := []string{"david", "ryo", "will"}
	for _, name := range names {
		client := Client{Name: name}
		server.AddClient(&client)
	}
	assert.Equal(t, len(names), len(server.Clients))
	for i := 0; i < len(names); i++ {
		server.RemoveClient(server.Clients[0])
	}

}
