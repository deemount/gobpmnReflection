package elements

import "strings"

var (
	structMessage = "message"
)

// IsMessage checks if the field is a message
func IsMessage(field string) bool {
	return strings.ToLower(field) == structMessage
}
