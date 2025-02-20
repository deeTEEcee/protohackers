package main

import (
	"github.com/stretchr/testify/assert"
	"protohackers/tcp"
	"strings"
	"testing"
)

func countSpaces(s string) int {
	return strings.Count(s, " ")
}

func TestRewriteFunc(t *testing.T) {
	//tonyAddress := "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
	//examples := []string{
	//	" 7F1u3wSD5RbOHQmupo9nx4TnhQ",
	//	"7iKDZEwPZSqIvDnHvVN2r0hUWXD5rHX ",
	//	" 7LOrwbDlS8NujgjddyogWgIM93MV5N2VR ",
	//	"7adNeSwJkMakpEcln9HEtthSRtxdmEHOT8T ",
	//}
	//for _, example := range examples {
	//	result := tcp.Rewrite(example)
	//	assert.Equal(t, countSpaces(example), countSpaces(result), example)
	//	assert.Equal(t, strings.TrimSpace(result), tonyAddress)
	//}
	//assert.Equal(t,
	//	"Hi alice, please send payment to 7YWHMfk9JZe0LM0g1ZauHuiSxhI",
	//	tcp.Rewrite("Hi alice, please send payment to 7iKDZEwPZSqIvDnHvVN2r0hUWXD5rHX"),
	//)
	tooLong := "This is too long: 7Yi1rFtaO2hucc97aanPVUBcvwc0Er11xkec"
	result := tcp.Rewrite(tooLong)
	assert.Equal(t,
		tooLong,
		result,
	)
}
