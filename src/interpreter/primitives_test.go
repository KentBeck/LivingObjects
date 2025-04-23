package main

import (
	"testing"
)

func TestIntegerPrimitives(t *testing.T) {
	t.Run("Addition", testAdditionPrimitive)
	t.Run("Subtraction", testSubtractionPrimitive)
	t.Run("Multiplication", testMultiplicationPrimitive)
	t.Run("LessThan", testLessThanPrimitive)
	t.Run("GreaterThan", testGreaterThanPrimitive)
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
	if IsIntegerImmediate(result) {
		intValue := GetIntegerImmediate(result)
		if intValue != 3 {
			t.Errorf("Expected result to be 3, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}

	// For immediate values, we don't check the class as it's encoded in the tag bits
	if !IsIntegerImmediate(result) && result.Class != vm.IntegerClass {
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
	if IsIntegerImmediate(result) {
		intValue := GetIntegerImmediate(result)
		if intValue != 10 {
			t.Errorf("Expected result to be 10, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}

	// For immediate values, we don't check the class as it's encoded in the tag bits
	if !IsIntegerImmediate(result) && result.Class != vm.IntegerClass {
		t.Errorf("Expected result class to be Integer, got %v", result.Class)
	}
}

func testAdditionPrimitive(t *testing.T) {
	// Create a VM
	vm := NewVM()

	plusSelector := NewSymbol("+")
	method := vm.lookupMethod(vm.IntegerClass, plusSelector)

	// Create two integer objects
	three := vm.NewInteger(3)
	four := vm.NewInteger(4)

	// Execute the primitive
	result := vm.executePrimitive(three, plusSelector, []*Object{four}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Addition primitive returned nil")
		return
	}

	// Check that the result is correct
	if IsIntegerImmediate(result) {
		intValue := GetIntegerImmediate(result)
		if intValue != 7 {
			t.Errorf("Expected result to be 7, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}

	// For immediate values, we don't check the class as it's encoded in the tag bits
	if !IsIntegerImmediate(result) && result.Class != vm.IntegerClass {
		t.Errorf("Expected result class to be Integer, got %v", result.Class)
	}
}

func testLessThanPrimitive(t *testing.T) {
	// Create a VM
	vm := NewVM()

	lessSelector := NewSymbol("<")
	method := vm.lookupMethod(vm.IntegerClass, lessSelector)

	// Create two integer objects
	two := vm.NewInteger(2)
	five := vm.NewInteger(5)

	// Execute the primitive
	result := vm.executePrimitive(two, lessSelector, []*Object{five}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Less than primitive returned nil")
		return
	}

	// Check that the result is correct
	// For immediate values, we can't access the Type field directly
	if !IsTrueImmediate(result) && !IsFalseImmediate(result) {
		t.Errorf("Expected result to be a boolean immediate value")
	}

	if !result.IsTrue() {
		t.Errorf("Expected result to be true, got false")
	}

	// Test the opposite case
	result = vm.executePrimitive(five, lessSelector, []*Object{two}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Less than primitive returned nil")
		return
	}

	// Check that the result is correct
	// For immediate values, we can't access the Type field directly
	if !IsTrueImmediate(result) && !IsFalseImmediate(result) {
		t.Errorf("Expected result to be a boolean immediate value")
	}

	if result.IsTrue() {
		t.Errorf("Expected result to be false, got true")
	}
}

func testGreaterThanPrimitive(t *testing.T) {
	// Create a VM
	vm := NewVM()

	greaterSelector := NewSymbol(">")
	method := vm.lookupMethod(vm.IntegerClass, greaterSelector)

	// Create two integer objects
	five := vm.NewInteger(5)
	two := vm.NewInteger(2)

	// Execute the primitive
	result := vm.executePrimitive(five, greaterSelector, []*Object{two}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Greater than primitive returned nil")
		return
	}

	// Check that the result is correct
	// For immediate values, we can't access the Type field directly
	if !IsTrueImmediate(result) && !IsFalseImmediate(result) {
		t.Errorf("Expected result to be a boolean immediate value")
	}

	if !result.IsTrue() {
		t.Errorf("Expected result to be true, got false")
	}

	// Test the opposite case
	result = vm.executePrimitive(two, greaterSelector, []*Object{five}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Greater than primitive returned nil")
		return
	}

	// Check that the result is correct
	// For immediate values, we can't access the Type field directly
	if !IsTrueImmediate(result) && !IsFalseImmediate(result) {
		t.Errorf("Expected result to be a boolean immediate value")
	}

	if result.IsTrue() {
		t.Errorf("Expected result to be false, got true")
	}
}
