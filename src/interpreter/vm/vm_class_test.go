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
			expected: virtualMachine.Classes.Get(vm.Integer),
		},
		{
			name:     "Boolean true",
			obj:      virtualMachine.TrueObject,
			expected: virtualMachine.Classes.Get(vm.True),
		},
		{
			name:     "Boolean false",
			obj:      virtualMachine.FalseObject,
			expected: virtualMachine.Classes.Get(vm.False),
		},
		{
			name:     "Nil",
			obj:      virtualMachine.NilObject,
			expected: virtualMachine.Classes.Get(vm.UndefinedObject),
		},
		{
			name:     "Class",
			obj:      pile.ClassToObject(virtualMachine.Classes.Get(vm.Object)),
			expected: virtualMachine.Classes.Get(vm.Object), // A class is its own class
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