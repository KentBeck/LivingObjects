package vm

import (
	"smalltalklsp/interpreter/pile"
)

// For backward compatibility
// These accessor fields allow existing code to access classes through
// the old field names while we migrate to using ClassRegistry

// ObjectClass field (struct member, not a method)
// Deprecated: Use vm.Classes.Get(Object) instead
func (vm *VM) GetObjectClass() *pile.Class {
	return vm.Classes.Get(Object)
}