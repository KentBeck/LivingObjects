package main

import (
	"encoding/binary"
	"testing"
)

func TestNilClassPanic(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method with the basicClass primitive
	basicClassSelector := NewSymbol("basicClass")
	basicClassMethod := NewMethod(basicClassSelector, vm.ObjectClass)
	basicClassMethod.Method.IsPrimitive = true
	basicClassMethod.Method.PrimitiveIndex = 5 // basicClass primitive

	// Add the method to the Object class method dictionary
	objectMethodDict := vm.ObjectClass.GetMethodDict()
	if objectMethodDict.Entries == nil {
		objectMethodDict.Entries = make(map[string]*Object)
	}
	objectMethodDict.Entries[GetSymbolValue(basicClassSelector)] = basicClassMethod

	// Create an object with a nil class
	objWithNilClass := &Object{
		Type:  OBJ_INSTANCE,
		Class: nil, // Explicitly set class to nil
	}

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
	context := NewContext(testMethod, objWithNilClass, []*Object{}, nil)

	// Execute the test method and expect a panic
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Expected a panic when accessing basicClass on an object with nil class, but no panic occurred")
		} else {
			t.Logf("Got expected panic: %v", r)
		}
	}()

	// This should panic
	vm.ExecuteContext(context)
}
