package main

import (
	"testing"
)

func TestNilClassPanic(t *testing.T) {
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

	// Create a test method that will send the basicClass message
	builder := NewMethodBuilder(vm.ObjectClass).Selector("test")
	selectorIndex, builder := builder.AddLiteral(basicClassSelector)

	// Create bytecodes for the test method
	builder.PushSelf()
	builder.SendMessage(selectorIndex, 0)
	builder.ReturnStackTop()

	testMethod := builder.Go()

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
