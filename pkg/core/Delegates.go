package core

import (
	"reflect"
	"strings"
)

var (
	repository = "DefinitionsRepository"
	fieldLong  = "definitions"
	fieldShort = "def"
)

// IsDefinitions checks if the field is a definitions
func IsDefinitions(field reflect.StructField) bool {
	return strings.ToLower(field.Name) == fieldShort || strings.ToLower(field.Name) == fieldLong
}

// IsNotDefinitions checks if the field is not a definitions
func IsNotDefinitions(field reflect.StructField) bool {
	return strings.ToLower(field.Name) != fieldShort && strings.ToLower(field.Name) != fieldLong
}
