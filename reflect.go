package gobpmn_reflection

import (
	"log"
	"reflect"
	"strings"

	"github.com/deemount/gobpmnReflection/internals/utils"
)

// Maps ...
type Maps struct {
	Anonym  map[int]string
	Config  map[int]string
	Builder map[int]string
	Words   map[int][]string
}

// Reflect ...
type Reflect struct {
	Maps
	IF        interface{}
	Temporary reflect.Value
	Element   reflect.Value
}

// NewReflect initialize interface {} p
func NewReflect(p interface{}) *Reflect {
	return &Reflect{IF: p}
}

// Interface where Element is the interface {}
func (ref *Reflect) Interface() *Reflect {
	ref.Element = reflect.ValueOf(&ref.IF).Elem()
	return ref
}

// New allocates a temporary variable with type of the struct.
// ref.Element.Elem() is the value contained in the interface
func (ref *Reflect) New() *Reflect {
	ref.Temporary = reflect.New(ref.Element.Elem().Type()).Elem()
	ref.Temporary.Set(ref.Element.Elem())
	return ref
}

/*
Maps initializes maps to analyze then later

  - anonym: all anonymous fields
  - config: all boolean fields
  - builder: all builder fields
  - words: all collected words splitted
*/
func (ref *Reflect) InitMaps() *Reflect {
	ref.Anonym = make(map[int]string)
	ref.Config = make(map[int]string)
	ref.Builder = make(map[int]string)
	ref.Words = make(map[int][]string)
	return ref
}

// reflect fields of interface {}
func (ref *Reflect) Reflection() {
	ref.anonymousFields()
	ref.builderType()
	ref.boolType()
}

func (ref *Reflect) Methods() {
	reflect.TypeOf(ref.IF).NumMethod()
}

// Set temporary variable values to interface {}
// This method is set inside at the end of a build method, where
// fields of a struct got reflected by names
func (ref *Reflect) Set() any {
	ref.Element.Set(ref.Temporary)
	return ref.Element.Interface()
}

// anonymousFields takes all compounds which are anonymous
func (ref *Reflect) anonymousFields() {

	fields := reflect.VisibleFields(reflect.TypeOf(ref.IF))
	i := 0
	for _, field := range fields {
		if ref.isAnonymous(field) {
			ref.Anonym[i] = field.Name
			ref.Words[i] = utils.Split(ref.Anonym[i])
			i++
		}
	}
	if i == 0 {
		log.Println("reflect.anonymousFields: main struct contains map of zero anonymous fields")
	}

}

// builderType with three static filter options

func (ref Reflect) builderType() {

	fields := reflect.VisibleFields(reflect.TypeOf(ref.IF))
	log.Printf("reflect.builderType: fields %+v", fields)
	count := 0
	index := len(ref.Words)
	for _, field := range fields {
		if ref.isNotAnonymous(field) && ref.isNotDefinitions(field) && ref.isBpmnReflection(field) {
			ref.Builder[count] = field.Name
			ref.Words[index] = utils.Split(ref.Builder[count])
			count++
			index++
		}
	}
	log.Printf("reflect.builderType: count %d", count)
	if count == 0 {
		log.Println("reflect.builderType: main struct contains map of zero builder fields")
	}

}

// boolType with one static filter option, which must be kind of reflect.Bool
// The bool type of a field in a struct describes mostly configuration settings
// The field must contain a sibling in title case, e.g. IsExecutable
func (ref Reflect) boolType() {

	fields := reflect.VisibleFields(reflect.TypeOf(ref.IF))
	log.Printf("reflect.boolType: fields %+v", fields)
	count := 0
	index := len(ref.Words)
	for _, field := range fields {
		if field.Type.Kind() == reflect.Bool {
			ref.Config[count] = field.Name
			ref.Words[index] = utils.Split(ref.Config[count])
			count++
			index++
		}
	}

}

// countWords ...
func (ref Reflect) countWords() {
	l := 0
	length := len(ref.Words)
	for i := 0; i < length; i++ {
		l += len(ref.Words[i])
	}
	log.Printf("reflect.countWords: %v", l)
}

// isDefinitions ...
func (ref *Reflect) isDefinitions(field reflect.StructField) bool {
	return strings.ToLower(field.Name) == "def" || strings.ToLower(field.Name) == "definitions"
}

// isNotDefinitions ...
func (ref *Reflect) isNotDefinitions(field reflect.StructField) bool {
	return strings.ToLower(field.Name) != "def" && strings.ToLower(field.Name) != "definitions"
}

// isAnonymous ...
func (ref *Reflect) isAnonymous(field reflect.StructField) bool {
	return field.Anonymous
}

// isNotAnonymous ...
func (ref *Reflect) isNotAnonymous(field reflect.StructField) bool {
	return !field.Anonymous
}

// isBpmnBuilder ...
func (ref *Reflect) isBpmnReflection(field reflect.StructField) bool {
	return field.Type.Name() == "Reflection"
}
