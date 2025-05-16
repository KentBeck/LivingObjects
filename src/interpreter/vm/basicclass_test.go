package vm_test

import (
	"smalltalklsp/interpreter/pile"
	"testing"

	"smalltalklsp/interpreter/vm"
)

func TestBasicClassPrimitive(t *testing.T) {
	virtualMachine := vm.NewVM()

	EnsureObjectIsClass(t, virtualMachine, virtualMachine.NewInteger(42), pile.ObjectToClass(virtualMachine.Globals["Integer"]))
	EnsureObjectIsClass(t, virtualMachine, pile.NewNil(), pile.ObjectToClass(virtualMachine.Globals["UndefinedObject"]))
	EnsureObjectIsClass(t, virtualMachine, virtualMachine.TrueObject, pile.ObjectToClass(virtualMachine.Globals["True"]))
	EnsureObjectIsClass(t, virtualMachine, virtualMachine.FalseObject, pile.ObjectToClass(virtualMachine.Globals["False"]))
	EnsureObjectIsClass(t, virtualMachine, virtualMachine.NewFloat(3.14), pile.ObjectToClass(virtualMachine.Globals["Float"]))
}

func EnsureObjectIsClass(t *testing.T, virtualMachine *vm.VM, object pile.ObjectInterface, expected interface{}) {
	// Get the actual class of the object
	actual := virtualMachine.GetClass(object.(*pile.Object))

	// Log information for debugging
	t.Log(pile.GetClassName(actual))
	t.Log(object.String())
	t.Log(actual)

	// Check that the result is the expected class
	if actual != expected {
		t.Errorf("Expected result to be %v, got %v", expected, actual)
	}
}