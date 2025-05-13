package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/vm"
)

func TestBasicClassPrimitive(t *testing.T) {
	virtualMachine := vm.NewVM()

	EnsureObjectIsClass(t, virtualMachine, virtualMachine.NewInteger(42), virtualMachine.Classes.Get(vm.Integer))
	EnsureObjectIsClass(t, virtualMachine, core.NewNil(), virtualMachine.Classes.Get(vm.UndefinedObject))
	EnsureObjectIsClass(t, virtualMachine, virtualMachine.TrueObject, virtualMachine.Classes.Get(vm.True))
	EnsureObjectIsClass(t, virtualMachine, virtualMachine.FalseObject, virtualMachine.Classes.Get(vm.False))
	EnsureObjectIsClass(t, virtualMachine, virtualMachine.NewFloat(3.14), virtualMachine.Classes.Get(vm.Float))
}

func EnsureObjectIsClass(t *testing.T, virtualMachine *vm.VM, object core.ObjectInterface, expected interface{}) {
	// Get the actual class of the object
	actual := virtualMachine.GetClass(object.(*core.Object))

	// Log information for debugging
	t.Log(classes.GetClassName(actual))
	t.Log(object.String())
	t.Log(actual)

	// Check that the result is the expected class
	if actual != expected {
		t.Errorf("Expected result to be %v, got %v", expected, actual)
	}
}
