package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/vm"
)

func TestGetClass(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Test cases
	tests := []struct {
		name     string
		obj      core.ObjectInterface
		expected core.ObjectInterface
	}{
		{
			name:     "Integer",
			obj:      virtualMachine.NewInteger(42),
			expected: virtualMachine.IntegerClass,
		},
		{
			name:     "Boolean true",
			obj:      virtualMachine.TrueObject,
			expected: virtualMachine.TrueClass,
		},
		{
			name:     "Boolean false",
			obj:      virtualMachine.FalseObject,
			expected: virtualMachine.FalseClass,
		},
		{
			name:     "Nil",
			obj:      virtualMachine.NilObject,
			expected: virtualMachine.NilClass,
		},
		{
			name:     "Class",
			obj:      classes.ClassToObject(virtualMachine.ObjectClass),
			expected: virtualMachine.ObjectClass, // A class is its own class
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := virtualMachine.GetClass(test.obj.(*core.Object))
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
