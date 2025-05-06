package main

import (
	"testing"
)

func TestBasicClassPrimitive(t *testing.T) {
	vm := NewVM()

	EnsureObjectIsClass(t, vm, vm.NewInteger(42), vm.IntegerClass)
	EnsureObjectIsClass(t, vm, NewNil(), vm.NilClass)
	EnsureObjectIsClass(t, vm, vm.TrueObject, vm.TrueClass)
	EnsureObjectIsClass(t, vm, vm.FalseObject, vm.FalseClass)
	EnsureObjectIsClass(t, vm, vm.NewFloat(3.14), vm.FloatClass)
}

func EnsureObjectIsClass(t *testing.T, vm *VM, object ObjectInterface, expected ObjectInterface) {
	// Get the actual class of the object
	actual := vm.GetClass(object.(*Object))

	// Log information for debugging
	t.Log(expected.(*Class).Name)
	t.Log(object.String())
	t.Log(actual)

	// Check that the result is the expected class
	if actual != expected {
		t.Errorf("Expected result to be %v, got %v", expected, actual)
	}
}
