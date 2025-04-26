package main

import (
	"testing"
)

func TestSendMessageStackManagement(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// We'll use the VM's Object and Integer classes
	integerClass := vm.IntegerClass

	// Create literals
	returnValueSelector := NewSymbol("returnValue")
	valueObj := vm.NewInteger(42)
	receiverObj := vm.NewInteger(10)

	// Create a simple method that returns a value using AddLiteral
	returnValueBuilder := NewMethodBuilder(integerClass).Selector("returnValue")

	// Add literals to the method builder
	valueIndex, _ := returnValueBuilder.AddLiteral(valueObj) // Literal 0: 42

	// Create bytecodes for the method: just return 42
	returnValueBytecodes := make([]byte, 0, 10)

	// PUSH_LITERAL valueIndex (42)
	returnValueBytecodes = append(returnValueBytecodes, PUSH_LITERAL)
	returnValueBytecodes = append(returnValueBytecodes, 0, 0, 0, byte(valueIndex))

	// RETURN_STACK_TOP
	returnValueBytecodes = append(returnValueBytecodes, RETURN_STACK_TOP)

	// Finalize the method
	returnValueBuilder.Bytecodes(returnValueBytecodes).Go()

	// Create a caller method that will call returnValue and then use the result using AddLiteral
	callerBuilder := NewMethodBuilder(integerClass).Selector("caller")

	// Add literals to the caller method builder
	receiverIndex, _ := callerBuilder.AddLiteral(receiverObj)                    // Literal 0: 10
	returnValueSelectorIndex, _ := callerBuilder.AddLiteral(returnValueSelector) // Literal 1: returnValue

	// Create bytecodes for the caller method:
	// 1. Push a value onto the stack that should be preserved
	// 2. Send returnValue message to receiver
	// 3. Check that both the original value and the result are on the stack
	callerBytecodes := make([]byte, 0, 20)

	// PUSH_LITERAL receiverIndex (10) - this is a value we want to preserve across the method call
	callerBytecodes = append(callerBytecodes, PUSH_LITERAL)
	callerBytecodes = append(callerBytecodes, 0, 0, 0, byte(receiverIndex))

	// PUSH_SELF - this will be the receiver of the returnValue message
	callerBytecodes = append(callerBytecodes, PUSH_SELF)

	// SEND_MESSAGE returnValue with 0 arguments
	callerBytecodes = append(callerBytecodes, SEND_MESSAGE)
	callerBytecodes = append(callerBytecodes, 0, 0, 0, byte(returnValueSelectorIndex))
	callerBytecodes = append(callerBytecodes, 0, 0, 0, 0) // 0 arguments

	// At this point, the stack should have two values:
	// 1. The original value (10)
	// 2. The result of the returnValue method (42)

	// RETURN_STACK_TOP - just return the top of the stack (which should be 42)
	callerBytecodes = append(callerBytecodes, RETURN_STACK_TOP)

	// Finalize the method
	callerMethod := callerBuilder.Bytecodes(callerBytecodes).Go()

	// Create a receiver for the caller method
	receiver := vm.NewInteger(5)

	// Create a context for the caller method
	context := NewContext(callerMethod, receiver, []*Object{}, nil)

	// Execute the context
	result, err := vm.ExecuteContext(context)
	if err != nil {
		t.Errorf("Error executing caller method: %v", err)
		return
	}

	// Check that the result is 42 (the value returned by the returnValue method)
	if IsIntegerImmediate(result) {
		intValue := GetIntegerImmediate(result)
		if intValue != 42 {
			t.Errorf("Expected result to be 42, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}

}

func TestSendMessageWithMultiplication(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// We'll use the VM's Object and Integer classes
	integerClass := vm.IntegerClass

	// Create literals
	returnValueSelector := NewSymbol("returnValue")
	valueObj := vm.NewInteger(42)
	timesSelector := NewSymbol("*")

	// Create a simple method that returns a value using AddLiteral
	returnValueBuilder := NewMethodBuilder(integerClass).Selector("returnValue")

	// Add literals to the method builder
	valueIndex, _ := returnValueBuilder.AddLiteral(valueObj) // Literal 0: 42

	// Create bytecodes for the method: just return 42
	returnValueBytecodes := make([]byte, 0, 10)

	// PUSH_LITERAL valueIndex (42)
	returnValueBytecodes = append(returnValueBytecodes, PUSH_LITERAL)
	returnValueBytecodes = append(returnValueBytecodes, 0, 0, 0, byte(valueIndex))

	// RETURN_STACK_TOP
	returnValueBytecodes = append(returnValueBytecodes, RETURN_STACK_TOP)

	// Finalize the method
	returnValueBuilder.Bytecodes(returnValueBytecodes).Go()

	// Create the multiplication method
	NewMethodBuilder(integerClass).
		Selector("*").
		Primitive(2). // Multiplication
		Go()

	// Create a method that will call returnValue and then use the result for multiplication using AddLiteral
	multiplyBuilder := NewMethodBuilder(integerClass).Selector("multiply")

	// Add literals to the multiply method builder
	returnValueSelectorIndex, _ := multiplyBuilder.AddLiteral(returnValueSelector) // Literal 0: returnValue
	timesSelectorIndex, _ := multiplyBuilder.AddLiteral(timesSelector)             // Literal 1: *

	// Create bytecodes for the multiply method:
	// 1. Push self (for later use in multiplication)
	// 2. Send returnValue message to self
	// 3. Multiply self by the result
	multiplyBytecodes := make([]byte, 0, 20)

	// PUSH_SELF (for later use in multiplication)
	multiplyBytecodes = append(multiplyBytecodes, PUSH_SELF)

	// DUPLICATE (to save a copy for later multiplication)
	multiplyBytecodes = append(multiplyBytecodes, PUSH_SELF)

	// SEND_MESSAGE returnValue with 0 arguments
	multiplyBytecodes = append(multiplyBytecodes, SEND_MESSAGE)
	multiplyBytecodes = append(multiplyBytecodes, 0, 0, 0, byte(returnValueSelectorIndex))
	multiplyBytecodes = append(multiplyBytecodes, 0, 0, 0, 0) // 0 arguments

	// SEND_MESSAGE * with 1 argument
	multiplyBytecodes = append(multiplyBytecodes, SEND_MESSAGE)
	multiplyBytecodes = append(multiplyBytecodes, 0, 0, 0, byte(timesSelectorIndex))
	multiplyBytecodes = append(multiplyBytecodes, 0, 0, 0, 1) // 1 argument

	// RETURN_STACK_TOP
	multiplyBytecodes = append(multiplyBytecodes, RETURN_STACK_TOP)

	// Finalize the method
	multiplyMethod := multiplyBuilder.Bytecodes(multiplyBytecodes).Go()

	// Create a receiver for the multiply method
	multiplyReceiver := vm.NewInteger(5)

	// Create a context for the multiply method
	multiplyContext := NewContext(multiplyMethod, multiplyReceiver, []*Object{}, nil)

	// Execute the context
	multiplyResult, err := vm.ExecuteContext(multiplyContext)
	if err != nil {
		t.Errorf("Error executing multiply method: %v", err)
		return
	}

	// Check that the result is 5 * 42 = 210
	if IsIntegerImmediate(multiplyResult) {
		intValue := GetIntegerImmediate(multiplyResult)
		if intValue != 210 {
			t.Errorf("Expected result to be 210 (5 * 42), got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", multiplyResult)
	}
}
