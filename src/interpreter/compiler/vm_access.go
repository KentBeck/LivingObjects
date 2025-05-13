package compiler

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// VMAccess provides an interface for accessing VM factory methods
type VMAccess interface {
	NewSymbol(value string) *core.Object
	NewMethod(selector *core.Object, class *classes.Class) *core.Object
}

// DefaultVMAccess is the global VM access instance
// This should be set by the VM during initialization
var DefaultVMAccess VMAccess

// RegisterVMAccess registers a VM instance for compiler access
func RegisterVMAccess(vm VMAccess) {
	DefaultVMAccess = vm
}