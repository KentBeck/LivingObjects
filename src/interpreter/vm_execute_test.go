package main

import (
	"testing"
)

func TestExecuteContextEmptyMethod(t *testing.T) {
	// Create a VM for testing
	vm := NewVM()

	// Create a method with no bytecodes
	methodObj := NewMethod(NewSymbol("emptyMethod"), vm.ObjectClass)

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

	// Create a method that pushes a value onto the stack
	methodObj := NewMethod(NewSymbol("pushMethod"), vm.ObjectClass)
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, PUSH_LITERAL)
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// Add a literal
	methodObj.Method.Literals = append(methodObj.Method.Literals, vm.NewInteger(42))

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
	methodObj := NewMethod(NewSymbol("errorMethod"), vm.ObjectClass)

	// Add an invalid bytecode
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, 255) // Invalid bytecode

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Execute the context
	_, err := vm.ExecuteContext(context)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}
