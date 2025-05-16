package vm_test

import (
	"smalltalklsp/interpreter/pile"
	"testing"

	"smalltalklsp/interpreter/vm"
)

func TestGetClass(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Test cases
	tests := []struct {
		name     string
		obj      pile.ObjectInterface
		expected pile.ObjectInterface
	}{
		{
			name:     "Integer",
			obj:      virtualMachine.NewInteger(42),
			expected: pile.ObjectToClass(virtualMachine.Globals["Integer"]),
		},
		{
			name:     "Boolean true",
			obj:      virtualMachine.TrueObject,
			expected: pile.ObjectToClass(virtualMachine.Globals["True"]),
		},
		{
			name:     "Boolean false",
			obj:      virtualMachine.FalseObject,
			expected: pile.ObjectToClass(virtualMachine.Globals["False"]),
		},
		{
			name:     "Nil",
			obj:      virtualMachine.NilObject,
			expected: pile.ObjectToClass(virtualMachine.Globals["UndefinedObject"]),
		},
		{
			name:     "Class",
			obj:      pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])),
			expected: pile.ObjectToClass(virtualMachine.Globals["Object"]), // A class is its own class
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := virtualMachine.GetClass(test.obj.(*pile.Object))
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestGetClassPanics(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Test with nil object
	t.Run("Nil object", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic with nil object, but no panic occurred")
			}
		}()
		virtualMachine.GetClass(nil)
	})
}