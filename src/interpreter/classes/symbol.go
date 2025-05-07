package classes

import (
	"fmt"
	"unsafe"

	"smalltalklsp/interpreter/core"
)

// Symbol represents a Smalltalk symbol object
type Symbol struct {
	core.Object
	Value string
}

// NewSymbol creates a new symbol object
func NewSymbol(value string) *core.Object {
	sym := &Symbol{
		Object: core.Object{
			TypeField: core.OBJ_SYMBOL,
		},
		Value: value,
	}
	return SymbolToObject(sym)
}

// SymbolToObject converts a Symbol to an Object
func SymbolToObject(s *Symbol) *core.Object {
	return (*core.Object)(unsafe.Pointer(s))
}

// ObjectToSymbol converts an Object to a Symbol
func ObjectToSymbol(o core.ObjectInterface) *Symbol {
	return (*Symbol)(unsafe.Pointer(o.(*core.Object)))
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
func GetSymbolValue(o core.ObjectInterface) string {
	if o.Type() == core.OBJ_SYMBOL {
		return ObjectToSymbol(o).Value
	}
	panic("GetSymbolValue: not a symbol")
}
