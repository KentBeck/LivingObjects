package main

import (
	"testing"
)

func TestDuplicate(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a context
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Create a test object
	testObj := vm.NewInteger(42)

	// Push the test object onto the stack
	context.Push(testObj)

	// Check the stack pointer
	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	// Execute the DUPLICATE bytecode
	err := vm.ExecuteDuplicate(context)
	if err != nil {
		t.Errorf("ExecuteDuplicate returned an error: %v", err)
	}

	// Check the stack pointer
	if context.StackPointer != 2 {
		t.Errorf("Expected stack pointer to be 2, got %d", context.StackPointer)
	}

	// Check that the top two objects on the stack are the same
	top := context.Pop()
	if top != testObj {
		t.Errorf("Expected top of stack to be %v, got %v", testObj, top)
	}

	next := context.Pop()
	if next != testObj {
		t.Errorf("Expected next on stack to be %v, got %v", testObj, next)
	}
}
