package vm_test

import (
	"smalltalklsp/interpreter/pile"
	"testing"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/vm"
)

func TestIntegerPrimitives(t *testing.T) {
	t.Run("Addition", testAdditionPrimitive)
	t.Run("Subtraction", testSubtractionPrimitive)
	t.Run("Multiplication", testMultiplicationPrimitive)
	t.Run("LessThan", testLessThanPrimitive)
	t.Run("GreaterThan", testGreaterThanPrimitive)
}

func testSubtractionPrimitive(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Add primitive methods to the Integer class
	integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])
	minusSelector := pile.NewSymbol("-")
	minusMethod := compiler.NewMethodBuilder(integerClass).
		Selector("-").
		Primitive(4). // Subtraction primitive
		Go()

	five := virtualMachine.NewInteger(5)
	two := virtualMachine.NewInteger(2)
	method := minusMethod

	// Execute the primitive
	result := virtualMachine.ExecutePrimitive(five, minusSelector, []*pile.Object{two}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Subtraction primitive returned nil")
		return
	}

	// Check that the result is correct
	if pile.IsIntegerImmediate(result) {
		intValue := pile.GetIntegerImmediate(result)
		if intValue != 3 {
			t.Errorf("Expected result to be 3, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}

	// For immediate values, we don't check the class as it's encoded in the tag bits
	if !pile.IsIntegerImmediate(result) && result.Class() != pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Integer"])) {
		t.Errorf("Expected result class to be Integer, got %v", result.Class())
	}
}

func testMultiplicationPrimitive(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Add primitive methods to the Integer class
	integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])
	timesSelector := pile.NewSymbol("*")
	timesMethod := compiler.NewMethodBuilder(integerClass).
		Selector("*").
		Primitive(2). // Multiplication primitive
		Go()

	five := virtualMachine.NewInteger(5)
	two := virtualMachine.NewInteger(2)
	method := timesMethod

	// Execute the primitive
	result := virtualMachine.ExecutePrimitive(five, timesSelector, []*pile.Object{two}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Multiplication primitive returned nil")
		return
	}

	// Check that the result is correct
	if pile.IsIntegerImmediate(result) {
		intValue := pile.GetIntegerImmediate(result)
		if intValue != 10 {
			t.Errorf("Expected result to be 10, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}

	// For immediate values, we don't check the class as it's encoded in the tag bits
	if !pile.IsIntegerImmediate(result) && result.Class() != pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Integer"])) {
		t.Errorf("Expected result class to be Integer, got %v", result.Class())
	}
}

func testAdditionPrimitive(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Add primitive methods to the Integer class
	integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])
	plusSelector := pile.NewSymbol("+")
	plusMethod := compiler.NewMethodBuilder(integerClass).
		Selector("+").
		Primitive(1). // Addition primitive
		Go()

	three := virtualMachine.NewInteger(3)
	four := virtualMachine.NewInteger(4)
	method := plusMethod

	// Execute the primitive
	result := virtualMachine.ExecutePrimitive(three, plusSelector, []*pile.Object{four}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Addition primitive returned nil")
		return
	}

	// Check that the result is correct
	if pile.IsIntegerImmediate(result) {
		intValue := pile.GetIntegerImmediate(result)
		if intValue != 7 {
			t.Errorf("Expected result to be 7, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}

	// For immediate values, we don't check the class as it's encoded in the tag bits
	if !pile.IsIntegerImmediate(result) && result.Class() != pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Integer"])) {
		t.Errorf("Expected result class to be Integer, got %v", result.Class())
	}
}

func testLessThanPrimitive(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Add primitive methods to the Integer class
	integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])
	lessSelector := pile.NewSymbol("<")
	lessMethod := compiler.NewMethodBuilder(integerClass).
		Selector("<").
		Primitive(6). // Less than primitive
		Go()

	two := virtualMachine.NewInteger(2)
	five := virtualMachine.NewInteger(5)
	method := lessMethod

	// Execute the primitive
	result := virtualMachine.ExecutePrimitive(two, lessSelector, []*pile.Object{five}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Less than primitive returned nil")
		return
	}

	// Check that the result is correct
	// For immediate values, we can't access the Type field directly
	if !pile.IsTrueImmediate(result) && !pile.IsFalseImmediate(result) {
		t.Errorf("Expected result to be a boolean immediate value")
	}

	if !result.IsTrue() {
		t.Errorf("Expected result to be true, got false")
	}

	// Test the opposite case
	result = virtualMachine.ExecutePrimitive(five, lessSelector, []*pile.Object{two}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Less than primitive returned nil")
		return
	}

	// Check that the result is correct
	// For immediate values, we can't access the Type field directly
	if !pile.IsTrueImmediate(result) && !pile.IsFalseImmediate(result) {
		t.Errorf("Expected result to be a boolean immediate value")
	}

	if result.IsTrue() {
		t.Errorf("Expected result to be false, got true")
	}
}

func testGreaterThanPrimitive(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Add primitive methods to the Integer class
	integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])
	greaterSelector := pile.NewSymbol(">")
	greaterMethod := compiler.NewMethodBuilder(integerClass).
		Selector(">").
		Primitive(7). // Greater than primitive
		Go()

	five := virtualMachine.NewInteger(5)
	two := virtualMachine.NewInteger(2)
	method := greaterMethod

	// Execute the primitive
	result := virtualMachine.ExecutePrimitive(five, greaterSelector, []*pile.Object{two}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Greater than primitive returned nil")
		return
	}

	// Check that the result is correct
	// For immediate values, we can't access the Type field directly
	if !pile.IsTrueImmediate(result) && !pile.IsFalseImmediate(result) {
		t.Errorf("Expected result to be a boolean immediate value")
	}

	if !result.IsTrue() {
		t.Errorf("Expected result to be true, got false")
	}

	// Test the opposite case
	result = virtualMachine.ExecutePrimitive(two, greaterSelector, []*pile.Object{five}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Greater than primitive returned nil")
		return
	}

	// Check that the result is correct
	// For immediate values, we can't access the Type field directly
	if !pile.IsTrueImmediate(result) && !pile.IsFalseImmediate(result) {
		t.Errorf("Expected result to be a boolean immediate value")
	}

	if result.IsTrue() {
		t.Errorf("Expected result to be false, got true")
	}
}
