package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// NewSymbol creates a symbol object with proper class field
func (vm *VM) NewSymbol(value string) *core.Object {
	sym := classes.NewSymbolInternal(value)
	symObj := classes.SymbolToObject(sym)
	symObj.SetClass(classes.ClassToObject(vm.Classes.Get(Object))) // Symbols are instances of the Object class for now
	return symObj
}