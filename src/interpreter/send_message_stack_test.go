package main

import (
	"testing"
)

func TestSendMessageStackManagement(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create the Object class
	objectClass := vm.ObjectClass

	// Create the Integer class
	integerClass := NewClass("Integer", objectClass)

	// Create the method dictionary for the Integer class
	integerMethodDict := integerClass.GetMethodDict()

	// Create a simple method that returns a value
	returnValueSelector := NewSymbol("returnValue")
	returnValueMethod := NewMethod(returnValueSelector, integerClass)

	// Add the method to the Integer class
	integerMethodDict.Entries[returnValueSelector.SymbolValue] = returnValueMethod

	// Create a literal for the method
	valueObj := vm.NewIntegerWithClass(42, integerClass)

	// Add the literal to the method
	returnValueMethod.Method.Literals = append(returnValueMethod.Method.Literals, valueObj) // Literal 0: 42

	// Create bytecodes for the method: just return 42
	// PUSH_LITERAL 0 (42)
	returnValueMethod.Method.Bytecodes = append(returnValueMethod.Method.Bytecodes, PUSH_LITERAL)
	returnValueMethod.Method.Bytecodes = append(returnValueMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// RETURN_STACK_TOP
	returnValueMethod.Method.Bytecodes = append(returnValueMethod.Method.Bytecodes, RETURN_STACK_TOP)

	// Create a caller method that will call returnValue and then use the result
	callerSelector := NewSymbol("caller")
	callerMethod := NewMethod(callerSelector, integerClass)

	// Add the caller method to the Integer class
	integerMethodDict.Entries[callerSelector.SymbolValue] = callerMethod

	// Create literals for the caller method
	receiverObj := vm.NewIntegerWithClass(10, integerClass)

	// Add literals to the caller method
	callerMethod.Method.Literals = append(callerMethod.Method.Literals, receiverObj)         // Literal 0: 10
	callerMethod.Method.Literals = append(callerMethod.Method.Literals, returnValueSelector) // Literal 1: returnValue

	// Create bytecodes for the caller method:
	// 1. Push a value onto the stack that should be preserved
	// 2. Send returnValue message to receiver
	// 3. Check that both the original value and the result are on the stack

	// PUSH_LITERAL 0 (10) - this is a value we want to preserve across the method call
	callerMethod.Method.Bytecodes = append(callerMethod.Method.Bytecodes, PUSH_LITERAL)
	callerMethod.Method.Bytecodes = append(callerMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// PUSH_SELF - this will be the receiver of the returnValue message
	callerMethod.Method.Bytecodes = append(callerMethod.Method.Bytecodes, PUSH_SELF)

	// SEND_MESSAGE returnValue with 0 arguments
	callerMethod.Method.Bytecodes = append(callerMethod.Method.Bytecodes, SEND_MESSAGE)
	callerMethod.Method.Bytecodes = append(callerMethod.Method.Bytecodes, 0, 0, 0, 1) // Selector index 1 (returnValue)
	callerMethod.Method.Bytecodes = append(callerMethod.Method.Bytecodes, 0, 0, 0, 0) // 0 arguments

	// At this point, the stack should have two values:
	// 1. The original value (10)
	// 2. The result of the returnValue method (42)

	// RETURN_STACK_TOP - just return the top of the stack (which should be 42)
	callerMethod.Method.Bytecodes = append(callerMethod.Method.Bytecodes, RETURN_STACK_TOP)

	// Create a receiver for the caller method
	receiver := vm.NewIntegerWithClass(5, integerClass)

	// Create a context for the caller method
	context := NewContext(callerMethod, receiver, []*Object{}, nil)

	// Execute the context
	result, err := vm.ExecuteContext(context)
	if err != nil {
		t.Errorf("Error executing caller method: %v", err)
		return
	}

	// Check that the result is 42 (the value returned by the returnValue method)
	if result.Type != OBJ_INTEGER {
		t.Errorf("Expected result to be an integer, got %v", result.Type)
	}

	if result.IntegerValue != 42 {
		t.Errorf("Expected result to be 42, got %d", result.IntegerValue)
	}

}

func TestSendMessageWithMultiplication(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create the Object class
	objectClass := vm.ObjectClass

	// Create the Integer class
	integerClass := NewClass("Integer", objectClass)

	// Create the method dictionary for the Integer class
	integerMethodDict := integerClass.GetMethodDict()

	// Create a simple method that returns a value
	returnValueSelector := NewSymbol("returnValue")
	returnValueMethod := NewMethod(returnValueSelector, integerClass)

	// Add the method to the Integer class
	integerMethodDict.Entries[returnValueSelector.SymbolValue] = returnValueMethod

	// Create a literal for the method
	valueObj := vm.NewIntegerWithClass(42, integerClass)

	// Add the literal to the method
	returnValueMethod.Method.Literals = append(returnValueMethod.Method.Literals, valueObj) // Literal 0: 42

	// Create bytecodes for the method: just return 42
	// PUSH_LITERAL 0 (42)
	returnValueMethod.Method.Bytecodes = append(returnValueMethod.Method.Bytecodes, PUSH_LITERAL)
	returnValueMethod.Method.Bytecodes = append(returnValueMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// RETURN_STACK_TOP
	returnValueMethod.Method.Bytecodes = append(returnValueMethod.Method.Bytecodes, RETURN_STACK_TOP)

	// Create a method that will call returnValue and then use the result for multiplication
	multiplySelector := NewSymbol("multiply")
	multiplyMethod := NewMethod(multiplySelector, integerClass)

	// Add the multiply method to the Integer class
	integerMethodDict.Entries[multiplySelector.SymbolValue] = multiplyMethod

	// Create the multiplication selector
	timesSelector := NewSymbol("*")
	timesMethod := NewMethod(timesSelector, integerClass)
	timesMethod.Method.IsPrimitive = true
	timesMethod.Method.PrimitiveIndex = 2 // Multiplication
	integerMethodDict.Entries[timesSelector.SymbolValue] = timesMethod

	// Add literals to the multiply method
	multiplyMethod.Method.Literals = append(multiplyMethod.Method.Literals, returnValueSelector) // Literal 0: returnValue
	multiplyMethod.Method.Literals = append(multiplyMethod.Method.Literals, timesSelector)       // Literal 1: *

	// Create bytecodes for the multiply method:
	// 1. Push self (for later use in multiplication)
	// 2. Send returnValue message to self
	// 3. Multiply self by the result

	// PUSH_SELF (for later use in multiplication)
	multiplyMethod.Method.Bytecodes = append(multiplyMethod.Method.Bytecodes, PUSH_SELF)

	// DUPLICATE (to save a copy for later multiplication)
	multiplyMethod.Method.Bytecodes = append(multiplyMethod.Method.Bytecodes, PUSH_SELF)

	// SEND_MESSAGE returnValue with 0 arguments
	multiplyMethod.Method.Bytecodes = append(multiplyMethod.Method.Bytecodes, SEND_MESSAGE)
	multiplyMethod.Method.Bytecodes = append(multiplyMethod.Method.Bytecodes, 0, 0, 0, 0) // Selector index 0 (returnValue)
	multiplyMethod.Method.Bytecodes = append(multiplyMethod.Method.Bytecodes, 0, 0, 0, 0) // 0 arguments

	// SEND_MESSAGE * with 1 argument
	multiplyMethod.Method.Bytecodes = append(multiplyMethod.Method.Bytecodes, SEND_MESSAGE)
	multiplyMethod.Method.Bytecodes = append(multiplyMethod.Method.Bytecodes, 0, 0, 0, 1) // Selector index 1 (*)
	multiplyMethod.Method.Bytecodes = append(multiplyMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

	// RETURN_STACK_TOP
	multiplyMethod.Method.Bytecodes = append(multiplyMethod.Method.Bytecodes, RETURN_STACK_TOP)

	// Create a receiver for the multiply method
	multiplyReceiver := vm.NewIntegerWithClass(5, integerClass)

	// Create a context for the multiply method
	multiplyContext := NewContext(multiplyMethod, multiplyReceiver, []*Object{}, nil)

	// Execute the context
	multiplyResult, err := vm.ExecuteContext(multiplyContext)
	if err != nil {
		t.Errorf("Error executing multiply method: %v", err)
		return
	}

	// Check that the result is 5 * 42 = 210
	if multiplyResult.Type != OBJ_INTEGER {
		t.Errorf("Expected result to be an integer, got %v", multiplyResult.Type)
	}

	if multiplyResult.IntegerValue != 210 {
		t.Errorf("Expected result to be 210 (5 * 42), got %d", multiplyResult.IntegerValue)
	}
}
