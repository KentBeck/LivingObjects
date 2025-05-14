package pile

import (
	"fmt"
	"unsafe"
)

// Array represents a Smalltalk array object
type Array struct {
	Object
	Elements []*Object
}

// NewArray creates a new array object (deprecated - use vm.NewArray instead)
func NewArray(size int) *Array {
	return &Array{Object: Object{TypeField: OBJ_ARRAY}, Elements: make([]*Object, size)}
}

// ArrayToObject converts an Array to an Object
func ArrayToObject(a *Array) *Object {
	return (*Object)(unsafe.Pointer(a))
}

// ObjectToArray converts an Object to an Array
func ObjectToArray(o *Object) *Array {
	return (*Array)(unsafe.Pointer(o))
}

// String returns a string representation of the array object
func (a *Array) String() string {
	return fmt.Sprintf("Array(%d)", len(a.Elements))
}

// Size returns the size of the array
func (a *Array) Size() int {
	return len(a.Elements)
}

// At returns the element at the given index
func (a *Array) At(index int) *Object {
	if index < 0 || index >= len(a.Elements) {
		panic("index out of bounds")
	}
	return a.Elements[index]
}

// AtPut sets the element at the given index
func (a *Array) AtPut(index int, value *Object) {
	if index < 0 || index >= len(a.Elements) {
		panic("index out of bounds")
	}
	a.Elements[index] = value
}