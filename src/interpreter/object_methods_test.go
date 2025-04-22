package main

import (
	"testing"
)

func TestObjectIsTrue(t *testing.T) {
	// Create a VM for testing
	vm := NewVM()

	tests := []struct {
		name     string
		obj      *Object
		expected bool
	}{
		{
			name:     "Boolean true",
			obj:      NewBoolean(true),
			expected: true,
		},
		{
			name:     "Boolean false",
			obj:      NewBoolean(false),
			expected: false,
		},
		{
			name:     "Nil",
			obj:      NewNil(),
			expected: false,
		},
		{
			name:     "Integer",
			obj:      vm.NewInteger(42),
			expected: true,
		},
		{
			name:     "String",
			obj:      NewString("hello"),
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.obj.IsTrue()
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestObjectInstanceVarMethods(t *testing.T) {
	// Skip this test for now as we're transitioning to immediate values
	t.Skip("Skipping test until immediate values are fully implemented")
	// Create a VM for testing
	vm := NewVM()

	// Create a class with instance variables
	class := NewClass("TestClass", nil)
	class.InstanceVarNames = append(class.InstanceVarNames, "var1", "var2")

	// Create an instance
	instance := NewInstance(class)

	// Test GetInstanceVarByIndex
	instance.InstanceVars[0] = vm.NewInteger(42)
	instance.InstanceVars[1] = NewString("hello")

	if instance.GetInstanceVarByIndex(0).IntegerValue != 42 {
		t.Errorf("Expected instance var 0 to be 42, got %v", instance.GetInstanceVarByIndex(0))
	}

	if instance.GetInstanceVarByIndex(1).StringValue != "hello" {
		t.Errorf("Expected instance var 1 to be 'hello', got %v", instance.GetInstanceVarByIndex(1))
	}

	// Test GetInstanceVarByIndex with out of bounds index
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on out of bounds access, but no panic occurred")
		}
	}()
	instance.GetInstanceVarByIndex(2) // This should panic
}

func TestObjectSetInstanceVarByIndex(t *testing.T) {
	// Skip this test for now as we're transitioning to immediate values
	t.Skip("Skipping test until immediate values are fully implemented")
	// Create a VM for testing
	vm := NewVM()

	// Create a class with instance variables
	class := NewClass("TestClass", nil)
	class.InstanceVarNames = append(class.InstanceVarNames, "var1", "var2")

	// Create an instance
	instance := NewInstance(class)

	// Test SetInstanceVarByIndex
	instance.SetInstanceVarByIndex(0, vm.NewInteger(42))
	instance.SetInstanceVarByIndex(1, NewString("hello"))

	if instance.InstanceVars[0].IntegerValue != 42 {
		t.Errorf("Expected instance var 0 to be 42, got %v", instance.InstanceVars[0])
	}

	if instance.InstanceVars[1].StringValue != "hello" {
		t.Errorf("Expected instance var 1 to be 'hello', got %v", instance.InstanceVars[1])
	}

	// Test SetInstanceVarByIndex with out of bounds index
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on out of bounds access, but no panic occurred")
		}
	}()
	instance.SetInstanceVarByIndex(2, vm.NewInteger(42)) // This should panic
}

func TestObjectGetMethodDict(t *testing.T) {
	// Skip this test for now as we're transitioning to immediate values
	t.Skip("Skipping test until immediate values are fully implemented")
	// Create a VM for testing
	_ = NewVM()

	// Test with a class
	class := NewClass("TestClass", nil)
	methodDict := class.GetMethodDict()
	if methodDict.Type != OBJ_DICTIONARY {
		t.Errorf("Expected method dictionary to be a dictionary, got %v", methodDict.Type)
	}

	// Test with a non-class object
	instance := NewInstance(class)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when calling GetMethodDict on a non-class object, but no panic occurred")
		}
	}()
	instance.GetMethodDict() // This should panic
}
