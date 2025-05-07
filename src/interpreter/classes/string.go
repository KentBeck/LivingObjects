package classes

import (
	"fmt"
	"unsafe"

	"smalltalklsp/interpreter/core"
)

// String represents a Smalltalk string object
type String struct {
	core.Object
	Value string
}

// NewString creates a new string object
func NewString(value string) *String {
	str := &String{
		Object: core.Object{
			TypeField: core.OBJ_STRING,
		},
		Value: value,
	}
	return str
}

// StringToObject converts a String to an Object
func StringToObject(s *String) *core.Object {
	return (*core.Object)(unsafe.Pointer(s))
}

// ObjectToString converts an Object to a String
func ObjectToString(o core.ObjectInterface) *String {
	return (*String)(unsafe.Pointer(o.(*core.Object)))
}

// String returns a string representation of the string object
func (s *String) String() string {
	return fmt.Sprintf("'%s'", s.Value)
}

// GetValue returns the string value
func (s *String) GetValue() string {
	return s.Value
}

// SetValue sets the string value
func (s *String) SetValue(value string) {
	s.Value = value
}

// Length returns the length of the string
func (s *String) Length() int {
	return len(s.Value)
}

// CharAt returns the character at the given index
func (s *String) CharAt(index int) byte {
	if index < 0 || index >= len(s.Value) {
		panic("index out of bounds")
	}
	return s.Value[index]
}

// Substring returns a substring of the string
func (s *String) Substring(start, end int) *String {
	if start < 0 || start >= len(s.Value) || end < 0 || end > len(s.Value) || start > end {
		panic("invalid substring range")
	}
	return NewString(s.Value[start:end])
}

// Concat concatenates this string with another string
func (s *String) Concat(other *String) *String {
	return NewString(s.Value + other.Value)
}

// Equal returns true if this string is equal to another string
func (s *String) Equal(other *String) bool {
	return s.Value == other.Value
}
