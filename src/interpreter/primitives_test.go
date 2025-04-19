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

	// Create the Integer class
	integerClass := NewClass("Integer", vm.ObjectClass)

	// Create the method dictionary for the Integer class
	integerMethodDict := integerClass.GetMethodDict()

	// Create the subtraction method
	minusSelector := NewSymbol("-")
	minusMethod := NewMethod(minusSelector, integerClass)
	minusMethod.Method.IsPrimitive = true
	minusMethod.Method.PrimitiveIndex = 4 // Subtraction
	integerMethodDict.Entries[minusSelector.SymbolValue] = minusMethod

	// Create two integer objects
	five := NewInteger(5)
	five.Class = integerClass

	two := NewInteger(2)
	two.Class = integerClass

	// Execute the primitive
	result := vm.executePrimitive(five, minusSelector, []*Object{two})

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
	if result.Class != integerClass {
		t.Errorf("Expected result class to be Integer, got %v", result.Class)
	}
}

func testMultiplicationPrimitive(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create the Integer class
	integerClass := NewClass("Integer", vm.ObjectClass)

	// Create the method dictionary for the Integer class
	integerMethodDict := integerClass.GetMethodDict()

	// Create the multiplication method
	timesSelector := NewSymbol("*")
	timesMethod := NewMethod(timesSelector, integerClass)
	timesMethod.Method.IsPrimitive = true
	timesMethod.Method.PrimitiveIndex = 2 // Multiplication
	integerMethodDict.Entries[timesSelector.SymbolValue] = timesMethod

	// Create two integer objects
	five := NewInteger(5)
	five.Class = integerClass

	two := NewInteger(2)
	two.Class = integerClass

	// Execute the primitive
	result := vm.executePrimitive(five, timesSelector, []*Object{two})

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
	if result.Class != integerClass {
		t.Errorf("Expected result class to be Integer, got %v", result.Class)
	}
}
