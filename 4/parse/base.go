package parse

import "strings"

func ParseMessage(msg string) (key string, value string, isInsert bool) {
	// Returns key, value if its insert
	// Returns key, value="" if its retrieve
	result := strings.SplitN(msg, "=", 2)
	if len(result) == 1 {
		return result[0], "", false
	} else {
		return result[0], result[1], true
	}
}
