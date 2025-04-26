package main

import (
	"testing"
)

func TestExecuteContextEmptyMethod(t *testing.T) {
	// Create a VM for testing
	vm := NewVM()

	// Create a method with no bytecodes using MethodBuilder
	methodObj := NewMethodBuilder(vm.ObjectClass).
		Selector("emptyMethod").
		Go()

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Execute the context
	result, err := vm.ExecuteContext(context)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Empty method should return nil
	if result != vm.NilObject {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestExecuteContextWithStackValue(t *testing.T) {
	// Create a VM for testing
	vm := NewVM()

	// Create a method that pushes a value onto the stack using MethodBuilder
	builder := NewMethodBuilder(vm.ObjectClass).Selector("pushMethod")
	literalIndex, builder := builder.AddLiteral(vm.NewInteger(42))
	methodObj := builder.PushLiteral(literalIndex).Go()

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Execute the context
	result, err := vm.ExecuteContext(context)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Method should return the value on the stack
	if IsIntegerImmediate(result) {
		intValue := GetIntegerImmediate(result)
		if intValue != 42 {
			t.Errorf("Expected 42, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}
}

func TestExecuteContextWithError(t *testing.T) {
	// Create a VM for testing
	vm := NewVM()

	// Create a method with an invalid bytecode
	// Since we can't use the fluent API for invalid bytecodes, we'll create the method
	// and then manually set the bytecodes
	methodObj := NewMethodBuilder(vm.ObjectClass).
		Selector("errorMethod").
		Go()

	// Set invalid bytecode manually
	methodObj.Method.Bytecodes = []byte{255} // Invalid bytecode

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Execute the context
	_, err := vm.ExecuteContext(context)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}
