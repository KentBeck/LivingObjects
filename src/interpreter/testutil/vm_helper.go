package testutil

import (
	"sync"
	"testing"

	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/vm"
)

var (
	testVM     *vm.VM
	testVMOnce sync.Once
)

// GetTestVM returns a singleton VM instance for tests
func GetTestVM(t *testing.T) *vm.VM {
	testVMOnce.Do(func() {
		testVM = vm.NewVM()
	})
	return testVM
}

// Create object helper functions for tests that use the VM's factories
// These should eventually replace the direct object creation functions

// NewMethod creates a method for tests
func NewMethod(t *testing.T, selector *core.Object, class *core.Class) *core.Object {
	vm := GetTestVM(t)
	return vm.NewMethod(selector, class)
}

// NewSymbol creates a symbol for tests
func NewSymbol(t *testing.T, value string) *core.Object {
	vm := GetTestVM(t)
	return vm.NewSymbol(value)
}

// NewString creates a string for tests
func NewString(t *testing.T, value string) *core.Object {
	vm := GetTestVM(t)
	return vm.NewString(value)
}

// NewArray creates an array for tests
func NewArray(t *testing.T, size int) *core.Object {
	vm := GetTestVM(t)
	return vm.NewArray(size)
}

// NewDictionary creates a dictionary for tests
func NewDictionary(t *testing.T) *core.Object {
	vm := GetTestVM(t)
	return vm.NewDictionary()
}

// NewBlock creates a block for tests
func NewBlock(t *testing.T, outerContext interface{}) *core.Object {
	vm := GetTestVM(t)
	return vm.NewBlock(outerContext)
}

// NewByteArray creates a byte array for tests
func NewByteArray(t *testing.T, size int) *core.Object {
	vm := GetTestVM(t)
	return vm.NewByteArray(size)
}