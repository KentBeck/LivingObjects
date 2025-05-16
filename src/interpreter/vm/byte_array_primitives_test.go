package vm_test

import (
	"smalltalklsp/interpreter/pile"
	"testing"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/vm"
)

// TestByteArrayAtPrimitive tests the ByteArray at: primitive
func TestByteArrayAtPrimitive(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Add primitive methods to the ByteArray class
	byteArrayClass := pile.ObjectToClass(virtualMachine.Globals["ByteArray"])
	atSelector := pile.NewSymbol("at:")
	atMethod := compiler.NewMethodBuilder(byteArrayClass).
		Selector("at:").
		Primitive(50). // ByteArray at: primitive
		Go()

	// Create a test byte array with 3 elements
	byteArray := virtualMachine.NewByteArray(3)
	byteArrayObj := pile.ObjectToByteArray(byteArray)
	
	// Fill the byte array with values
	byteArrayObj.AtPut(0, 10)
	byteArrayObj.AtPut(1, 20)
	byteArrayObj.AtPut(2, 30)
	
	// Create the index argument (1-based in Smalltalk)
	indexArg := virtualMachine.NewInteger(2)
	
	// Execute the primitive
	result := virtualMachine.ExecutePrimitive(byteArray, atSelector, []*pile.Object{indexArg}, atMethod)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("ByteArray at: primitive returned nil")
		return
	}

	// Check that the result is an integer
	if !pile.IsIntegerImmediate(result) {
		t.Errorf("Expected integer result, got %v", result.Type())
		return
	}

	// Check that the result is the expected value
	value := pile.GetIntegerImmediate(result)
	if value != 20 {
		t.Errorf("Expected 20, got %d", value)
	}

	// Test out of bounds access
	indexOutOfBounds := virtualMachine.NewInteger(4)
	
	// Execute the primitive with an out of bounds index
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for index out of bounds, but got none")
			}
		}()
		virtualMachine.ExecutePrimitive(byteArray, atSelector, []*pile.Object{indexOutOfBounds}, atMethod)
	}()
}

// TestByteArrayAtPutPrimitive tests the ByteArray at:put: primitive
func TestByteArrayAtPutPrimitive(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Add primitive methods to the ByteArray class
	byteArrayClass := pile.ObjectToClass(virtualMachine.Globals["ByteArray"])
	atPutSelector := pile.NewSymbol("at:put:")
	atPutMethod := compiler.NewMethodBuilder(byteArrayClass).
		Selector("at:put:").
		Primitive(51). // ByteArray at:put: primitive
		Go()

	// Create a test byte array with 3 elements
	byteArray := virtualMachine.NewByteArray(3)
	
	// Create the index argument (1-based in Smalltalk)
	indexArg := virtualMachine.NewInteger(2)
	
	// Create the value argument
	valueArg := virtualMachine.NewInteger(42)
	
	// Execute the primitive
	result := virtualMachine.ExecutePrimitive(byteArray, atPutSelector, []*pile.Object{indexArg, valueArg}, atPutMethod)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("ByteArray at:put: primitive returned nil")
		return
	}

	// Check that the result is the value argument
	if result != valueArg {
		t.Errorf("Expected result to be the value argument")
	}

	// Check that the byte array was updated
	byteArrayObj := pile.ObjectToByteArray(byteArray)
	value := byteArrayObj.At(1) // 0-based in Go
	if value != 42 {
		t.Errorf("Expected byte array value at index 1 to be 42, got %d", value)
	}

	// Test out of bounds access
	indexOutOfBounds := virtualMachine.NewInteger(4)
	
	// Execute the primitive with an out of bounds index
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for index out of bounds, but got none")
			}
		}()
		virtualMachine.ExecutePrimitive(byteArray, atPutSelector, []*pile.Object{indexOutOfBounds, valueArg}, atPutMethod)
	}()

	// Test value out of range
	valueOutOfRange := virtualMachine.NewInteger(256)
	
	// Execute the primitive with a value out of range
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for value out of range, but got none")
			}
		}()
		virtualMachine.ExecutePrimitive(byteArray, atPutSelector, []*pile.Object{indexArg, valueOutOfRange}, atPutMethod)
	}()
}
