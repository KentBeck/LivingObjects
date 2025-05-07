package core_test

import (
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/vm"
)

func TestObjectString(t *testing.T) {
	// Create a VM for testing
	virtualMachine := vm.NewVM()

	tests := []struct {
		name     string
		obj      core.ObjectInterface
		expected string
	}{
		{
			name:     "Integer",
			obj:      virtualMachine.NewInteger(42),
			expected: "42",
		},
		{
			name:     "Boolean true",
			obj:      core.NewBoolean(true),
			expected: "true",
		},
		{
			name:     "Boolean false",
			obj:      core.NewBoolean(false),
			expected: "false",
		},
		{
			name:     "Nil",
			obj:      core.NewNil(),
			expected: "nil",
		},
		{
			name:     "String",
			obj:      classes.StringToObject(classes.NewString("hello")),
			expected: "'hello'",
		},
		{
			name:     "Symbol",
			obj:      classes.NewSymbol("test"),
			expected: "#test",
		},
		{
			name:     "Array",
			obj:      classes.ArrayToObject(classes.NewArray(3)),
			expected: "Array(3)",
		},
		{
			name:     "Dictionary",
			obj:      classes.NewDictionary(),
			expected: "Dictionary(0)",
		},
		{
			name:     "Instance with class",
			obj:      core.NewInstance((*core.Class)(unsafe.Pointer(virtualMachine.ObjectClass))),
			expected: "a Object",
		},
		{
			name: "Instance without class", // This should panic
			obj: &core.Object{
				TypeField: core.OBJ_INSTANCE,
			},
			expected: "an Object",
		},
		{
			name:     "Class",
			obj:      classes.ClassToObject(virtualMachine.ObjectClass),
			expected: "Class Object",
		},
		{
			name:     "Method with selector",
			obj:      compiler.NewMethodBuilder(virtualMachine.ObjectClass).Selector("test").Go(),
			expected: "Method test",
		},
		{
			name: "Method without selector",
			obj: classes.MethodToObject(&classes.Method{
				Object: core.Object{
					TypeField: core.OBJ_METHOD,
				},
			}),
			expected: "a Method",
		},
		{
			name: "Unknown object type",
			obj: &core.Object{
				TypeField: 255, // Invalid type
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
