package parse

import "strings"

func ParseMessage(msg string) (key string, value string, isInsert bool) {
	// Returns key, value if its insert
	// Returns key, value="" if its retrieve
	before, after, splitExists := strings.Cut(msg, "=")
	return before, after, splitExists
}
