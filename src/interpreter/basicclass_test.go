package main

import (
	"encoding/binary"
	"testing"
)

func TestBasicClassPrimitive(t *testing.T) {
	// Create a VM
	vm := NewVM()

	EnsureObjectIsClass(t, vm, vm.NewInteger(42), vm.IntegerClass)
	// EnsureObjectIsClass(t, vm, NewNil(), vm.NilClass) also Boolean
}

func EnsureObjectIsClass(t *testing.T, vm *VM, intObj *Object, intObjClass *Object) {
	testMethod := NewMethod(NewSymbol("test"), vm.ObjectClass)

	// Add the basicClass selector to the literals
	basicClassSelector := NewSymbol("basicClass")
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
}
