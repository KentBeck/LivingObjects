package main

import (
	"testing"
)

func TestContextPush(t *testing.T) {
	// Create a VM for testing
	vm := NewVM()

	methodObj := NewMethodBuilder(vm.ObjectClass).
		Selector("test").
		Go()

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Test pushing an object
	obj := vm.NewInteger(42)
	context.Push(obj)

	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	if context.Stack[0] != obj {
		t.Errorf("Expected stack[0] to be %v, got %v", obj, context.Stack[0])
	}

	// Test pushing nil
	context.Push(nil) // seems like this should be a panic

	if context.StackPointer != 2 {
		t.Errorf("Expected stack pointer to be 2, got %d", context.StackPointer)
	}

	if context.Stack[1] != nil {
		t.Errorf("Expected stack[1] to be nil, got %v", context.Stack[1])
	}

	// Test stack growth
	// Fill the stack to its initial capacity
	initialCapacity := len(context.Stack)
	for i := context.StackPointer; i < initialCapacity; i++ {
		context.Push(obj)
	}

	// Push one more object to trigger stack growth
	context.Push(obj)

	// Verify the stack grew
	if len(context.Stack) <= initialCapacity {
		t.Errorf("Expected stack to grow beyond %d, got %d", initialCapacity, len(context.Stack))
	}
}

func TestContextPop(t *testing.T) {
	// Create a VM for testing
	vm := NewVM()

	methodObj := NewMethodBuilder(vm.ObjectClass).
		Selector("test").
		Go()

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Test popping from an empty stack
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on stack underflow, but no panic occurred")
		}
	}()

	// Push an object and then pop it
	obj := vm.NewInteger(42)
	context.Push(obj)
	popped := context.Pop()

	if popped != obj {
		t.Errorf("Expected popped object to be %v, got %v", obj, popped)
	}

	if context.StackPointer != 0 {
		t.Errorf("Expected stack pointer to be 0, got %d", context.StackPointer)
	}

	// This should panic
	context.Pop()
}

func TestContextTop(t *testing.T) {
	// Create a VM for testing
	vm := NewVM()

	methodObj := NewMethodBuilder(vm.ObjectClass).
		Selector("test").
		Go()

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Test top on an empty stack
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on empty stack, but no panic occurred")
		}
	}()

	// Push an object and then check top
	obj := vm.NewInteger(42)
	context.Push(obj)
	top := context.Top()

	if top != obj {
		t.Errorf("Expected top object to be %v, got %v", obj, top)
	}

	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	// Clear the stack
	context.StackPointer = 0

	// This should panic
	context.Top()
}

func TestContextTempVars(t *testing.T) {
	// Create a VM for testing
	vm := NewVM()

	methodObj := NewMethodBuilder(vm.ObjectClass).
		Selector("test").
		Go()

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Set up temporary variables
	context.TempVars = make([]*Object, 2)

	// Test setting and getting temporary variables
	obj := vm.NewInteger(42)
	context.SetTempVarByIndex(0, obj)

	if context.GetTempVarByIndex(0) != obj {
		t.Errorf("Expected temp var 0 to be %v, got %v", obj, context.GetTempVarByIndex(0))
	}

	// Test out of bounds access for GetTempVarByIndex
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic on out of bounds access for GetTempVarByIndex, but no panic occurred")
			}
		}()
		context.GetTempVarByIndex(2) // This should panic
	}()

	// Test setting nil value
	context.SetTempVarByIndex(1, nil)
	if context.TempVars[1] != nil {
		t.Errorf("Expected temp var 1 to be nil, got %v", context.TempVars[1])
	}
}
