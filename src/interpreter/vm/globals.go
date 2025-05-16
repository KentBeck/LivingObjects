package vm

import (
	"smalltalklsp/interpreter/pile"
)

// GetGlobal returns a global variable by name
func (vm *VM) GetGlobal(name string) *pile.Object {
	if obj, ok := vm.Globals[name]; ok {
		return obj
	}
	return nil
}
