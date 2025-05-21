package vm_test

import (
	"smalltalklsp/interpreter/pile"
	"testing"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/vm"
)

// TestStringSizePrimitive tests the String size primitive
func TestStringSizePrimitive(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Add primitive methods to the String class
	stringClass := pile.ObjectToClass(virtualMachine.Globals["String"])
	sizeMethod := compiler.NewMethodBuilder(stringClass).
		Primitive(30). // String size primitive
		Go("size")

	// Create a test string
	testString := virtualMachine.NewString("hello")
	method := sizeMethod

	// Create a selector object
	selectorObj := pile.NewSymbol("size")

	// Execute the primitive
	result := virtualMachine.ExecutePrimitive(testString, selectorObj, []*pile.Object{}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("String size primitive returned nil")
		return
	}

	// Check that the result is an integer
	if !pile.IsIntegerImmediate(result) {
		t.Errorf("Expected integer result, got %v", result.Type())
		return
	}

	// Check that the result is 5 (length of "hello")
	value := pile.GetIntegerImmediate(result)
	if value != 5 {
		t.Errorf("Expected size 5, got %d", value)
	}

	// Test with an empty string
	emptyString := virtualMachine.NewString("")
	result = virtualMachine.ExecutePrimitive(emptyString, selectorObj, []*pile.Object{}, method)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("String size primitive returned nil for empty string")
		return
	}

	// Check that the result is an integer
	if !pile.IsIntegerImmediate(result) {
		t.Errorf("Expected integer result for empty string, got %v", result.Type())
		return
	}

	// Check that the result is 0 (length of "")
	value = pile.GetIntegerImmediate(result)
	if value != 0 {
		t.Errorf("Expected size 0 for empty string, got %d", value)
	}
}
