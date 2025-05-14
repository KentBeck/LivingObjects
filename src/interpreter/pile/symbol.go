package pile

import (
	"fmt"
	"unsafe"
)

// Symbol represents a Smalltalk symbol object
type Symbol struct {
	Object
	Value string
}

// newSymbol creates a new symbol object without setting its class field
// This is a private helper function used by vm.NewSymbol
func NewSymbolInternal(value string) *Symbol {
	return &Symbol{
		Object: Object{
			TypeField: OBJ_SYMBOL,
		},
		Value: value,
	}
}

// NewSymbol creates a new symbol object (deprecated - use vm.NewSymbol instead)
func NewSymbol(value string) *Object {
	sym := NewSymbolInternal(value)
	return SymbolToObject(sym)
}

// SymbolToObject converts a Symbol to an Object
func SymbolToObject(s *Symbol) *Object {
	return (*Object)(unsafe.Pointer(s))
}

// ObjectToSymbol converts an Object to a Symbol
func ObjectToSymbol(o ObjectInterface) *Symbol {
	return (*Symbol)(unsafe.Pointer(o.(*Object)))
}

// String returns a string representation of the symbol object
func (s *Symbol) String() string {
	return fmt.Sprintf("#%s", s.Value)
}

// GetValue returns the symbol value
func (s *Symbol) GetValue() string {
	return s.Value
}

// SetValue sets the symbol value
func (s *Symbol) SetValue(value string) {
	s.Value = value
}

// Length returns the length of the symbol
func (s *Symbol) Length() int {
	return len(s.Value)
}

// Equal returns true if this symbol is equal to another symbol
func (s *Symbol) Equal(other *Symbol) bool {
	return s.Value == other.Value
}

// GetSymbolValue gets the value from a Symbol object
func GetSymbolValue(o ObjectInterface) string {
	if o.Type() == OBJ_SYMBOL {
		return ObjectToSymbol(o).Value
	}
	panic("GetSymbolValue: not a symbol")
}