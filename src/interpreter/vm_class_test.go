package main

import (
	"testing"
)

func TestGetClass(t *testing.T) {
	vm := NewVM()

	// Test cases
	tests := []struct {
		name     string
		obj      *Object
		expected *Object
	}{
		{
			name:     "Integer",
			obj:      vm.NewInteger(42),
			expected: vm.IntegerClass,
		},
		{
			name:     "Boolean true",
			obj:      vm.TrueObject,
			expected: vm.TrueClass,
		},
		{
			name:     "Boolean false",
			obj:      vm.FalseObject,
			expected: vm.FalseClass,
		},
		{
			name:     "Nil",
			obj:      vm.NilObject,
			expected: vm.NilClass,
		},
		{
			name:     "Class",
			obj:      vm.ObjectClass,
			expected: vm.ObjectClass, // A class is its own class
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := vm.GetClass(test.obj)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestGetClassPanics(t *testing.T) {
	vm := NewVM()

	// Test with nil object
	t.Run("Nil object", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic with nil object, but no panic occurred")
			}
		}()
		vm.GetClass(nil)
	})

	// Test with object that has nil class
	t.Run("Object with nil class", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic with object that has nil class, but no panic occurred")
			}
		}()
		objWithNilClass := &Object{
			type1: OBJ_METHOD,
			Class: nil, // Explicitly set class to nil
		}
		vm.GetClass(objWithNilClass)
	})
}
