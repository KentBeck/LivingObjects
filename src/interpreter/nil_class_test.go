package main

import (
	"encoding/binary"
	"testing"
)

func TestNilClassPanic(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method with the basicClass primitive using MethodBuilder
	basicClassSelector := NewSymbol("basicClass")

	// Create the method using MethodBuilder
	NewMethodBuilder(vm.ObjectClass).
		Selector("basicClass").
		Primitive(5). // basicClass primitive
		Go()

	// Get the method dictionary for later use
	objectMethodDict := vm.ObjectClass.GetMethodDict()
	if objectMethodDict.Entries == nil {
		objectMethodDict.Entries = make(map[string]*Object)
	}

	// Create an object with a nil class
	objWithNilClass := &Object{
		Type:  OBJ_INSTANCE,
		Class: nil, // Explicitly set class to nil
	}

	// Create bytecodes for the test method
	bytecodes := []byte{PUSH_SELF, SEND_MESSAGE}

	// Selector index is 0, arg count is 0 (each 4 bytes)
	selectorIndexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(selectorIndexBytes, 0)
	bytecodes = append(bytecodes, selectorIndexBytes...)

	argCountBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(argCountBytes, 0)
	bytecodes = append(bytecodes, argCountBytes...)

	bytecodes = append(bytecodes, RETURN_STACK_TOP)

	// Create a test method that will send the basicClass message using MethodBuilder
	testMethod := NewMethodBuilder(vm.ObjectClass).
		Selector("test").
		AddLiterals([]*Object{basicClassSelector}).
		Bytecodes(bytecodes).
		Go()

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
