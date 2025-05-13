package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// NewSymbol creates a symbol object with proper class field
func (vm *VM) NewSymbol(value string) *core.Object {
	sym := &classes.Symbol{
		Object: core.Object{
			TypeField: core.OBJ_SYMBOL,
		},
		Value: value,
	}
	
	symObj := classes.SymbolToObject(sym)
	symObj.SetClass(classes.ClassToObject(vm.ObjectClass)) // Symbols are instances of the Object class for now
	return symObj
}