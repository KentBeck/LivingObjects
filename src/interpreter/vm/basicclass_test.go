package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/vm"
)

func TestBasicClassPrimitive(t *testing.T) {
	virtualMachine := vm.NewVM()

	EnsureObjectIsClass(t, virtualMachine, virtualMachine.NewInteger(42), virtualMachine.IntegerClass)
	EnsureObjectIsClass(t, virtualMachine, core.NewNil(), virtualMachine.NilClass)
	EnsureObjectIsClass(t, virtualMachine, virtualMachine.TrueObject, virtualMachine.TrueClass)
	EnsureObjectIsClass(t, virtualMachine, virtualMachine.FalseObject, virtualMachine.FalseClass)
	EnsureObjectIsClass(t, virtualMachine, virtualMachine.NewFloat(3.14), virtualMachine.FloatClass)
}

func EnsureObjectIsClass(t *testing.T, virtualMachine *vm.VM, object core.ObjectInterface, expected interface{}) {
	// Get the actual class of the object
	actual := virtualMachine.GetClass(object.(*core.Object))

	// Log information for debugging
	t.Log(actual.GetName())
	t.Log(object.String())
	t.Log(actual)

	// Check that the result is the expected class
	if actual != expected {
		t.Errorf("Expected result to be %v, got %v", expected, actual)
	}
}
