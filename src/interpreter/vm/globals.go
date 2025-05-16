package vm

import (
	"smalltalklsp/interpreter/pile"
)

// GetGlobal returns a global variable by name
func (vm *VM) GetGlobal(name string) *pile.Object {
	// Look up the global in the globals map
	if obj, ok := vm.Globals[name]; ok {
		return obj
	}
	
	// Return nil if the global is not found
	return pile.MakeNilImmediate()
}
