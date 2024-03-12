package gobpmn_reflection

import (
	"reflect"
	"strings"
)

// Reflect ...
type Reflect struct {
	Anonym    map[int]string
	Config    map[int]string
	Rflct     map[int]string
	IF        interface{}
	Temporary reflect.Value
	Element   reflect.Value
}

// New initialize interface {} p
func New(p interface{}) *Reflect {
	return &Reflect{IF: p}
}

// Interface where Element is the interface {}
func (ref *Reflect) Interface() *Reflect {
	ref.Element = reflect.ValueOf(&ref.IF).Elem()
	return ref
}

// Allocate allocates a temporary variable with type of the struct.
// ref.Element.Elem() is the value contained in the interface
func (ref *Reflect) Allocate() *Reflect {
	ref.Temporary = reflect.New(ref.Element.Elem().Type()).Elem()
	ref.Temporary.Set(ref.Element.Elem())
	return ref
}

/*
Maps initializes maps to analyze then later

  - anonym: all anonymous fields
  - config: all boolean fields
  - rflct:  all reflection fields
*/
func (ref *Reflect) Maps() *Reflect {
	ref.Anonym = make(map[int]string)
	ref.Config = make(map[int]string)
	ref.Rflct = make(map[int]string)
	return ref
}

// Assign holds the methods to reflect fields of interface {}
func (ref *Reflect) Assign() {
	ref.anonymousFields()
	ref.reflectionType()
	ref.boolType()
}

// Set temporary variable values to interface {}
// This method is set inside at the end of a build method, where
// fields of a struct got reflected by names
func (ref *Reflect) Set() any {
	ref.Element.Set(ref.Temporary)
	return ref.Element.Interface()
}

/*
 * @private
 */

// anonymousFields takes all compounds which are anonymous
func (ref *Reflect) anonymousFields() {

	fields := reflect.VisibleFields(reflect.TypeOf(ref.IF))
	index := 0
	for _, field := range fields {
		if field.Anonymous {
			ref.Anonym[index] = field.Name
			index++
		}
	}

}

// reflectionType with three static filter options
func (ref Reflect) reflectionType() {

	fields := reflect.VisibleFields(reflect.TypeOf(ref.IF))
	index := 0
	for _, field := range fields {
		if !field.Anonymous && ref.isNotDefinitions(field) && field.Type.Name() == "Injection" {
			ref.Rflct[index] = field.Name
			index++
		}
	}

}

// boolType with one static filter option, which must be kind of reflect.Bool
// The bool type of a field in a struct describes mostly configuration settings
// The field must contain a sibling in title case, e.g. IsExecutable
func (ref Reflect) boolType() {

	fields := reflect.VisibleFields(reflect.TypeOf(ref.IF))
	index := 0
	for _, field := range fields {
		if field.Type.Kind() == reflect.Bool {
			ref.Config[index] = field.Name
			index++
		}
	}

}

// isDefinitions ...
func (ref *Reflect) isDefinitions(field reflect.StructField) bool {
	return strings.ToLower(field.Name) == "def" || strings.ToLower(field.Name) == "definitions"
}

// isNotDefinitions ...
func (ref *Reflect) isNotDefinitions(field reflect.StructField) bool {
	return strings.ToLower(field.Name) != "def" || strings.ToLower(field.Name) != "definitions"
}
