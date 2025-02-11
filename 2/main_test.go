package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindMinMax(t *testing.T) {
	arr := []int{1, 5, 10, 18, 21}
	assert.Equal(t, -1, findMinInt(arr, 0))
	assert.Equal(t, 1, findMinInt(arr, 1))
	assert.Equal(t, 1, findMinInt(arr, 3))
	assert.Equal(t, 5, findMinInt(arr, 5))
	assert.Equal(t, 5, findMinInt(arr, 6))

	assert.Equal(t, 18, findMaxInt(arr, 11))
	assert.Equal(t, 18, findMaxInt(arr, 18))
	assert.Equal(t, 19, findMaxInt(arr, 21))
	assert.Equal(t, 21, findMaxInt(arr, 21))
	assert.Equal(t, -1, findMaxInt(arr, 24))
}
