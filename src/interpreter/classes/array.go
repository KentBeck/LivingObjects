package classes

import (
	"fmt"
	"unsafe"

	"smalltalklsp/interpreter/core"
)

// Array represents a Smalltalk array object
type Array struct {
	core.Object
	Elements []*core.Object
}

// NewArray creates a new array object
func NewArray(size int) *Array {
	obj := &Array{
		Object: core.Object{
			TypeField: core.OBJ_ARRAY,
		},
		Elements: make([]*core.Object, size),
	}
	return obj
}

// ArrayToObject converts an Array to an Object
func ArrayToObject(a *Array) *core.Object {
	return (*core.Object)(unsafe.Pointer(a))
}

// ObjectToArray converts an Object to an Array
func ObjectToArray(o *core.Object) *Array {
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
func (a *Array) At(index int) *core.Object {
	if index < 0 || index >= len(a.Elements) {
		panic("index out of bounds")
	}
	return a.Elements[index]
}

// AtPut sets the element at the given index
func (a *Array) AtPut(index int, value *core.Object) {
	if index < 0 || index >= len(a.Elements) {
		panic("index out of bounds")
	}
	a.Elements[index] = value
}

// Copy returns a copy of the array
func (a *Array) Copy() *Array {
	newArray := NewArray(len(a.Elements))
	copy(newArray.Elements, a.Elements)
	return newArray
}

// Collect applies a function to each element of the array and returns a new array
func (a *Array) Collect(fn func(*core.Object) *core.Object) *Array {
	newArray := NewArray(len(a.Elements))
	for i, elem := range a.Elements {
		newArray.Elements[i] = fn(elem)
	}
	return newArray
}

// Select returns a new array with elements that satisfy the predicate
func (a *Array) Select(fn func(*core.Object) bool) *Array {
	// First pass: count elements that satisfy the predicate
	count := 0
	for _, elem := range a.Elements {
		if fn(elem) {
			count++
		}
	}

	// Second pass: create a new array with selected elements
	newArray := NewArray(count)
	index := 0
	for _, elem := range a.Elements {
		if fn(elem) {
			newArray.Elements[index] = elem
			index++
		}
	}
	return newArray
}

// Reject returns a new array with elements that do not satisfy the predicate
func (a *Array) Reject(fn func(*core.Object) bool) *Array {
	return a.Select(func(obj *core.Object) bool {
		return !fn(obj)
	})
}

// Do applies a function to each element of the array
func (a *Array) Do(fn func(*core.Object)) {
	for _, elem := range a.Elements {
		fn(elem)
	}
}

// WithIndexDo applies a function to each element of the array along with its index
func (a *Array) WithIndexDo(fn func(int, *core.Object)) {
	for i, elem := range a.Elements {
		fn(i, elem)
	}
}
