package main

import (
	"testing"
)

func TestBasicClassPrimitive(t *testing.T) {
	vm := NewVM()

	EnsureObjectIsClass(t, vm, vm.NewInteger(42), vm.IntegerClass)
	EnsureObjectIsClass(t, vm, NewNil(), vm.NilClass)
	EnsureObjectIsClass(t, vm, vm.TrueObject, vm.TrueClass)
	EnsureObjectIsClass(t, vm, vm.FalseObject, vm.FalseClass)
	EnsureObjectIsClass(t, vm, vm.NewFloat(3.14), vm.FloatClass)
}

func EnsureObjectIsClass(t *testing.T, vm *VM, object *Object, class *Object) {

	basicClassSelector := NewSymbol("basicClass")

	builder := NewMethodBuilder(vm.ObjectClass).Selector("test")
	selectorIndex, builder := builder.AddLiteral(basicClassSelector)

	// Create bytecodes for pushing self, sending basicClass message, and returning
	builder.PushSelf()
	builder.SendMessage(selectorIndex, 0)
	builder.ReturnStackTop()

	testMethod := builder.Go()

	// Create a context for the test method
	context := NewContext(testMethod, object, []*Object{}, nil)

	// Execute the test method
	result, err := vm.ExecuteContext(context)

	// Check for errors
	if err != nil {
		t.Errorf("Error executing test method: %v", err)
	}

	// Check that the result is the expected class
	if result != class {
		t.Errorf("Expected result to be Integer class, got %v", result)
	}
}
