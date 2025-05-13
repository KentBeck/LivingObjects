package vm

import (
	"smalltalklsp/interpreter/core"
)

// For backward compatibility
// These accessor fields allow existing code to access classes through
// the old field names while we migrate to using ClassRegistry

// ObjectClass field (struct member, not a method)
// Deprecated: Use vm.Classes.Get(Object) instead
func (vm *VM) GetObjectClass() *core.Class {
	return vm.Classes.Get(Object)
}