package pile

import (
	"fmt"
	"unsafe"
)

// String represents a Smalltalk string object
type String struct {
	Object
	Value string
}

// NewString creates a new string object (deprecated - use vm.NewString instead)
func NewString(value string) *String {
	return &String{Object: Object{TypeField: OBJ_STRING}, Value: value}
}

// StringToObject converts a String to an Object
func StringToObject(s *String) *Object {
	return (*Object)(unsafe.Pointer(s))
}

// ObjectToString converts an Object to a String
func ObjectToString(o ObjectInterface) *String {
	return (*String)(unsafe.Pointer(o.(*Object)))
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

// GetStringValue gets the string value of a string
// Panics if the object is not a string
func GetStringValue(obj *Object) string {
	// Check if it's an immediate value
	if IsImmediate(obj) {
		panic("GetStringValue: expected a string object, got an immediate value")
	}

	// Check if it's a string object
	if obj.Type() != OBJ_STRING {
		panic("GetStringValue: expected a string object, got a different type")
	}

	return ObjectToString(obj).GetValue()
}
