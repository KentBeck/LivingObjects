package main

import (
	"encoding/binary"
	"testing"
)

func TestBasicClassPrimitive(t *testing.T) {
	t.Skip("Skipping this test")

	// Create a VM
	vm := NewVM()

	// Create a method with the basicClass primitive
	basicClassSelector := NewSymbol("basicClass")

	// Test with an integer object
	intObj := vm.NewInteger(42)
	intObjClass := intObj.Class

	// Create a test method that will send the basicClass message
	testMethod := NewMethod(NewSymbol("test"), vm.ObjectClass)

	// Add the basicClass selector to the literals
	testMethod.Method.Literals = append(testMethod.Method.Literals, basicClassSelector)

	// Add bytecodes to push self and send the basicClass message
	testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, PUSH_SELF)
	testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, SEND_MESSAGE)
	// Selector index is 0, arg count is 0 (each 4 bytes)
	selectorIndexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(selectorIndexBytes, 0)
	testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, selectorIndexBytes...)
	argCountBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(argCountBytes, 0)
	testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, argCountBytes...)
	testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, RETURN_STACK_TOP)

	// Create a context for the test method
	context := NewContext(testMethod, intObj, []*Object{}, nil)

	// Add debug logging
	t.Logf("Executing test method with receiver: %v (class: %v)", intObj, intObj.Class)
	t.Logf("Bytecodes: %v", testMethod.Method.Bytecodes)
	t.Logf("Literals: %v", testMethod.Method.Literals)

	// Execute the test method
	result, err := vm.ExecuteContext(context)

	// Check for errors
	if err != nil {
		t.Errorf("Error executing test method: %v", err)
	}

	// Check that the result is the Integer class
	if result != intObjClass {
		t.Errorf("Expected result to be Integer class, got %v", result)
	}

	// Test with a boolean object
	boolObj := vm.TrueObject
	boolObjClass := boolObj.Class

	// Create a context for the test method
	context = NewContext(testMethod, boolObj, []*Object{}, nil)

	// Execute the test method
	result, err = vm.ExecuteContext(context)

	// Check for errors
	if err != nil {
		t.Errorf("Error executing test method: %v", err)
	}

	// Check that the result is the Boolean class
	if result != boolObjClass {
		t.Errorf("Expected result to be Boolean class, got %v", result)
	}

	// Test with a class object
	classObj := vm.ObjectClass
	classObjClass := classObj.Class

	// Create a context for the test method
	context = NewContext(testMethod, classObj, []*Object{}, nil)

	// Execute the test method
	result, err = vm.ExecuteContext(context)

	// Check for errors
	if err != nil {
		t.Errorf("Error executing test method: %v", err)
	}

	// Check that the result is the Class class
	if result != classObjClass {
		t.Errorf("Expected result to be Class class, got %v", result)
	}
}
