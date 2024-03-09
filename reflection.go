package gobpmn_reflection

import (
	"crypto/rand"
	"fmt"
	"hash/fnv"
	"log"
	"reflect"
	"strings"

	"github.com/deemount/gobpmnModels/pkg/core"
	"github.com/deemount/gobpmnReflection/internals/utils"
)

// Def ...
type Def core.DefinitionsRepository

// Reflection ...
type Reflection struct {
	Reflect
	Map       map[string][]interface{}
	Suffix    string
	HashTable []string
}

// Hash ...
func (r *Reflection) Hash() string {
	if r.Suffix == "" {
		result, _ := r.hash()
		r.Suffix = result.Suffix
	}
	return r.Suffix
}

// Inject all anonymous fields with a hash value by fields with type Builder
// It also setup reflected fields with boolean type and checks out the configuration by wording
// usage:
// e.g.
// var p CProcess
// p = r.inject(CProcess{}).(CProcess)
func (r *Reflection) Inject(p interface{}) interface{} { return r.inject(p) }

// Create receives the definitions repository by the app in p argument
// and calls the main elements to set the maps, including process parameters
// n of process.
func (r *Reflection) Create(p interface{}) { r.create(p) }

/*
 * @private
 */

// hash generates a hash value by using the crypto/rand package
// and the hash/fnv package to generate a 32-bit FNV-1a hash.
// If the error is not nil, it means that the hash value could not be generated.
// The suffix is used to generate a unique ID for each element of a process.
func (r Reflection) hash() (Reflection, error) {

	n := 8
	b := make([]byte, n)
	c := fnv.New32a()

	if _, err := rand.Read(b); err != nil {
		return Reflection{}, err
	}
	s := fmt.Sprintf("%x", b)

	if _, err := c.Write([]byte(s)); err != nil {
		return Reflection{}, err
	}
	defer c.Reset()

	result := Reflection{Suffix: fmt.Sprintf("%x", string(c.Sum(nil)))}

	return result, nil
}

// inject itself reflects a given struct and inject
// signed fields with hash values.
// There are two conditions to assign fields of a struct:
// a) The struct has anonymous fields
// b) The struct has no anymous fields
// It also counts the element in their specification to know
// how much elements of each package needs to be mapped later then.
func (r *Reflection) inject(p interface{}) interface{} {

	c := Count{}

	ref := NewReflect(p)
	ref.Interface().New().InitMaps().Reflection()

	// anonymous field are reflected
	if len(ref.Anonym) > 0 {

		length := len(ref.Anonym)

		// create anonymMap and hashMap
		anonymMap := make(map[int][]interface{}, length)
		hashMap := make(map[string][]interface{}, length)

		// walk through the map with names of anonymous fields
		// e.g. starts with customer support, customer, ...
		for index, field := range ref.Anonym {

			// get the reflected value of field
			n := ref.Temporary.FieldByName(field)

			// create the field map and the hash slice
			fieldMap := make(map[int][]interface{}, n.NumField())
			hashSlice := make([]interface{}, n.NumField())

			// append to anonymMap the name of anonymous field
			anonymMap[index] = append(anonymMap[index], n.Type().Name())

			// walk through the values of fields assigned to the interface {}
			for i := 0; i < n.NumField(); i++ {

				// get the name of field and append to fieldMap the name of field
				name := n.Type().Field(i).Name
				fieldMap[i] = append(fieldMap[i], name)

				// set by kind of reflected value above
				switch n.Field(i).Kind() {

				// kind is a bool
				case reflect.Bool:

					// only the first field, which IsExecutable, is set to true.
					// means, only one process in a collaboration can be executed at runtime
					// this can be changed in the future, if the engine fits for more execution
					// options
					if strings.Contains(name, "IsExecutable") && i == 0 {
						n.Field(0).SetBool(true)
						hashSlice[i] = bool(true)
					} else {
						hashSlice[i] = bool(false)
					}

				// kind is a struct
				case reflect.Struct:

					c.countPool(field, name)    // counts processes, participants and their shapes
					c.countMessage(field, name) // counts message flows and their edges
					c.countElements(name)       // counts elements

					// if the field Suffix is empty, generate hash value and
					// start to inject by each index of the given structs. Then,
					// check next element, generate hash value and inject the field Suffix
					r.injectCurrentField(i, hashSlice, n)
					r.injectNextField(i, hashSlice, n)

				}

			}

			// merge the hashSlice with the hashMap
			utils.MergeStringSliceToMap(hashMap, n.Type().Name(), hashSlice)

		}

	}

	// anonymous field aren't reflected
	if len(ref.Anonym) == 0 {

		// walk through the map with names of builder fields
		for _, builderField := range ref.Builder {

			// get the reflected name of builderField
			nonAnonymBuilderField := ref.Temporary.FieldByName(builderField)

			c.countProcess(builderField)  // count processes
			c.countElements(builderField) // counts elements

			hash, _ := r.hash()                              // generate hash value
			nonAnonymBuilderField.Set(reflect.ValueOf(hash)) // inject the field

			log.Printf("reflection.inject: inject struct field %v", builderField)

		}

		// walk through the map with names of boolean fields
		for _, configField := range ref.Config {
			// get the reflected value of field
			nonAnonymConfigField := ref.Temporary.FieldByName(configField)
			// only the first field, which IsExecutable is set to true
			nonAnonymConfigField.SetBool(true)
			log.Printf("reflection.inject: inject config field %v", configField)
		}

	}

	p = ref.Set()

	utils.MergeStructs(p, &c)

	return p

}

// injectCurrentField injects the current field with a hash value
func (r *Reflection) injectCurrentField(index int, slice []interface{}, field reflect.Value) {
	strHash := fmt.Sprintf("%s", field.Field(index).FieldByName("Suffix"))
	if strHash == "" {
		hash, _ := r.hash()
		slice[index] = hash.Suffix
		field.Field(index).Set(reflect.ValueOf(hash))
	}
}

// injectNextField injects the next field with a hash value
func (r *Reflection) injectNextField(index int, slice []interface{}, field reflect.Value) {
	if index+1 < field.NumField() {
		nexthash, _ := r.hash()
		slice[index+1] = nexthash.Suffix
		field.Field(index + 1).Set(reflect.ValueOf(nexthash))
	}
}

// create contains the reflected process definition (p interface{})
// and calls it by the reflected method name.
// This method hides specific setters (SetProcess, SetCollaboration, SetDiagram)
// in the example process by building the model with reflection.
func (r *Reflection) create(p interface{}) {

	// el is the interface {}
	el := reflect.ValueOf(&p).Elem()

	// Allocate a temporary variable with type of the struct.
	// el.Elem() is the value contained in the interface
	definitions := reflect.New(el.Elem().Type()).Elem() // *core.Definitions
	definitions.Set(el.Elem())                          // reflected process definitions el will be assigned to the core definitions

	// set collaboration, process and diagram
	collaboration := definitions.MethodByName("SetCollaboration")
	collaboration.Call([]reflect.Value{})

	process := definitions.MethodByName("SetProcess")
	process.Call([]reflect.Value{reflect.ValueOf(2)}) // r.Process represents number of processes

	diagram := definitions.MethodByName("SetDiagram")
	diagram.Call([]reflect.Value{reflect.ValueOf(1)}) // 1 represents number of diagrams

}
