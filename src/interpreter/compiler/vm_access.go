package compiler

import (
	"smalltalklsp/interpreter/pile"
)

// VMAccess provides an interface for accessing VM factory methods
type VMAccess interface {
	NewSymbol(value string) *pile.Object
	NewMethod(selector *pile.Object, class *pile.Class) *pile.Object
}

// DefaultVMAccess is the global VM access instance
// This should be set by the VM during initialization
var DefaultVMAccess VMAccess

// RegisterVMAccess registers a VM instance for compiler access
func RegisterVMAccess(vm VMAccess) {
	DefaultVMAccess = vm
}