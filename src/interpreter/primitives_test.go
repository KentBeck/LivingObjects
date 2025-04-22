package main

import (
	"testing"
)

func TestPrimitives(t *testing.T) {
	t.Run("Subtraction", testSubtractionPrimitive)
	t.Run("Multiplication", testMultiplicationPrimitive)
}

func testSubtractionPrimitive(t *testing.T) {
	// Create a VM
	vm := NewVM()

	minusSelector := NewSymbol("-")
	method := vm.lookupMethod(vm.IntegerClass, minusSelector)

	// Create two integer objects
	five := vm.NewInteger(5)
	two := vm.NewInteger(2)

	// Execute the primitive
	result := vm.executePrimitive(five, minusSelector, []*Object{two}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Subtraction primitive returned nil")
		return
	}

	// Check that the result is correct
	if result.Type != OBJ_INTEGER {
		t.Errorf("Expected result to be an integer, got %v", result.Type)
	}

	if result.IntegerValue != 3 {
		t.Errorf("Expected result to be 3, got %d", result.IntegerValue)
	}

	// Check that the class of the result is set correctly
	if result.Class != vm.IntegerClass {
		t.Errorf("Expected result class to be Integer, got %v", result.Class)
	}
}

func testMultiplicationPrimitive(t *testing.T) {
	// Create a VM
	vm := NewVM()

	timesSelector := NewSymbol("*")
	method := vm.lookupMethod(vm.IntegerClass, timesSelector)

	// Create two integer objects
	five := vm.NewInteger(5)
	two := vm.NewInteger(2)

	// Execute the primitive
	result := vm.executePrimitive(five, timesSelector, []*Object{two}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Multiplication primitive returned nil")
		return
	}

	// Check that the result is correct
	if result.Type != OBJ_INTEGER {
		t.Errorf("Expected result to be an integer, got %v", result.Type)
	}

	if result.IntegerValue != 10 {
		t.Errorf("Expected result to be 10, got %d", result.IntegerValue)
	}

	// Check that the class of the result is set correctly
	if result.Class != vm.IntegerClass {
		t.Errorf("Expected result class to be Integer, got %v", result.Class)
	}
}
