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
}

func TestProcess2(t *testing.T) {
	store := make([]InsertMessage, 0)
	insert(&store, 12346, 0)
	insert(&store, 12347, 0)
	insert(&store, 40960, 0)
	// Not sure what's expected but let's dot his.
	assert.Equal(t, 0.0, query(store, 12288, 0))
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
