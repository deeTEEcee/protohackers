package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
