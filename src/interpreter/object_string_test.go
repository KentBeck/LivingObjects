package main

import (
	"testing"
)

func TestObjectString(t *testing.T) {
	// Create a VM for testing
	vm := NewVM()

	tests := []struct {
		name     string
		obj      *Object
		expected string
	}{
		{
			name:     "Integer",
			obj:      vm.NewInteger(42),
			expected: "42",
		},
		{
			name:     "Boolean true",
			obj:      NewBoolean(true),
			expected: "true",
		},
		{
			name:     "Boolean false",
			obj:      NewBoolean(false),
			expected: "false",
		},
		{
			name:     "Nil",
			obj:      NewNil(),
			expected: "nil",
		},
		{
			name:     "String",
			obj:      StringToObject(NewString("hello")),
			expected: "'hello'",
		},
		{
			name:     "Symbol",
			obj:      NewSymbol("test"),
			expected: "#test",
		},
		{
			name:     "Array",
			obj:      NewArray(3),
			expected: "Array(3)",
		},
		{
			name:     "Dictionary",
			obj:      NewDictionary(),
			expected: "Dictionary(0)",
		},
		{
			name:     "Instance with class",
			obj:      NewInstance(vm.ObjectClass),
			expected: "a Object",
		},
		{
			name: "Instance without class",
			obj: &Object{
				type1: OBJ_INSTANCE,
			},
			expected: "an Object",
		},
		{
			name:     "Class",
			obj:      vm.ObjectClass,
			expected: "Class Object",
		},
		{
			name:     "Method with selector",
			obj:      NewMethodBuilder(vm.ObjectClass).Selector("test").Go(),
			expected: "Method test",
		},
		{
			name: "Method without selector",
			obj: &Object{
				type1:  OBJ_METHOD,
				Method: &Method{},
			},
			expected: "a Method",
		},
		{
			name: "Unknown object type",
			obj: &Object{
				type1: 255, // Invalid type
			},
			expected: "Unknown object",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.obj.String()
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}
