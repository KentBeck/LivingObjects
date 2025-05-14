package vm

import (
	"smalltalklsp/interpreter/pile"
)

// NewSymbol creates a symbol object with proper class field
func (vm *VM) NewSymbol(value string) *pile.Object {
	sym := pile.NewSymbolInternal(value)
	symObj := pile.SymbolToObject(sym)
	symObj.SetClass(pile.ClassToObject(vm.Classes.Get(Symbol))) // Symbols are instances of the Symbol class
	return symObj
}