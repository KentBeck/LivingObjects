package main

import (
	"testing"
)

// TestExecuteSendMessageExtended tests the ExecuteSendMessage function with more complex scenarios
func TestExecuteSendMessageExtended(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Test cases
	t.Run("primitive method", func(t *testing.T) {
		// Create a method with a SEND_MESSAGE bytecode for addition
		method := NewMethod(NewSymbol("test"), vm.ObjectClass)

		// Create literals
		twoObj := vm.NewInteger(2)
		threeObj := vm.NewInteger(3)
		plusSymbol := NewSymbol("+")

		// Add literals to the method
		method.Method.Literals = append(method.Method.Literals, twoObj)     // Index 0
		method.Method.Literals = append(method.Method.Literals, threeObj)   // Index 1
		method.Method.Literals = append(method.Method.Literals, plusSymbol) // Index 2

		// Create bytecodes for: 2 + 3
		// PUSH_LITERAL 0 (2)
		method.Method.Bytecodes = append(method.Method.Bytecodes, PUSH_LITERAL)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 0) // Index 0

		// PUSH_LITERAL 1 (3)
		method.Method.Bytecodes = append(method.Method.Bytecodes, PUSH_LITERAL)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 1) // Index 1

		// SEND_MESSAGE 2 ("+") with 1 argument
		method.Method.Bytecodes = append(method.Method.Bytecodes, SEND_MESSAGE)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 2) // Selector index 2
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

		// Create a context
		context := NewContext(method, vm.ObjectClass, []*Object{}, nil)

		// Execute the PUSH_LITERAL bytecodes to set up the stack
		context.PC = 0
		if err := vm.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += InstructionSize(PUSH_LITERAL)

		context.PC = 5
		if err := vm.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += InstructionSize(PUSH_LITERAL)

		// Execute the SEND_MESSAGE bytecode
		context.PC = 10
		result, err := vm.ExecuteSendMessage(context)
		if err != nil {
			t.Fatalf("Error executing SEND_MESSAGE: %s", err)
		}

		// Check the result
		if IsIntegerImmediate(result) {
			intValue := GetIntegerImmediate(result)
			if intValue != 5 {
				t.Errorf("Expected 5, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %s", result)
		}

		// Check the stack
		if context.StackPointer != 1 {
			t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
		}

		// Check the stack top
		stackTop := context.Stack[0]
		if IsIntegerImmediate(stackTop) {
			intValue := GetIntegerImmediate(stackTop)
			if intValue != 5 {
				t.Errorf("Expected stack top to be 5, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %s", stackTop)
		}
	})

	t.Run("non-primitive method", func(t *testing.T) {
		// We'll use the VM's Object and Integer classes
		objectClass := vm.ObjectClass
		integerClass := vm.IntegerClass

		// Create a factorial method for Integer
		factorialSelector := NewSymbol("factorial")
		factorialMethod := NewMethod(factorialSelector, integerClass)

		// Add the method to the Integer class
		methodDict := integerClass.GetMethodDict()
		methodDict.Entries[factorialSelector.SymbolValue] = factorialMethod

		// Create literals for the factorial method
		oneObj := vm.NewInteger(1)

		// Add literals to the factorial method
		factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, oneObj) // Index 0

		// Create bytecodes for factorial: if self = 1 { return 1 } else { ... }
		// For simplicity, we'll just return 1 for any input
		// PUSH_LITERAL 0 (1)
		factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_LITERAL)
		factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

		// RETURN_STACK_TOP
		factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, RETURN_STACK_TOP)

		// Create a method that calls factorial
		testMethod := NewMethod(NewSymbol("test"), objectClass)

		// Create literals for the test method
		fiveObj := vm.NewInteger(5)

		// Add literals to the test method
		testMethod.Method.Literals = append(testMethod.Method.Literals, fiveObj)           // Index 0
		testMethod.Method.Literals = append(testMethod.Method.Literals, factorialSelector) // Index 1

		// Create bytecodes for: 5 factorial
		// PUSH_LITERAL 0 (5)
		testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, PUSH_LITERAL)
		testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

		// SEND_MESSAGE 1 ("factorial") with 0 arguments
		testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, SEND_MESSAGE)
		testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 1) // Selector index 1
		testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 0) // 0 arguments

		// Create a context
		context := NewContext(testMethod, objectClass, []*Object{}, nil)

		// Set the VM's object class
		vm.ObjectClass = objectClass
		vm.Globals["Object"] = objectClass
		vm.Globals["Integer"] = integerClass

		// Execute the PUSH_LITERAL bytecode to set up the stack
		context.PC = 0
		if err := vm.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += InstructionSize(PUSH_LITERAL)

		// Execute the SEND_MESSAGE bytecode
		context.PC = 5
		result, err := vm.ExecuteSendMessage(context)
		if err != nil {
			t.Fatalf("Error executing SEND_MESSAGE: %s", err)
		}

		// Check the result
		if IsIntegerImmediate(result) {
			intValue := GetIntegerImmediate(result)
			if intValue != 1 {
				t.Errorf("Expected 1, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %s", result)
		}

		// Check the stack
		if context.StackPointer != 1 {
			t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
		}

		// Check the stack top
		stackTop := context.Stack[0]
		if IsIntegerImmediate(stackTop) {
			intValue := GetIntegerImmediate(stackTop)
			if intValue != 1 {
				t.Errorf("Expected stack top to be 1, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %s", stackTop)
		}
	})

	t.Run("method not found", func(t *testing.T) {
		t.Skip("Implement message not understood later")

		// Create a method with a SEND_MESSAGE bytecode for an unknown method
		method := NewMethod(NewSymbol("test"), vm.ObjectClass)

		// Create literals
		receiver := vm.NewInteger(2)
		unknownSelector := NewSymbol("unknown")

		// Add literals to the method
		method.Method.Literals = append(method.Method.Literals, receiver)        // Index 0
		method.Method.Literals = append(method.Method.Literals, unknownSelector) // Index 1

		// Create bytecodes for: 2 unknown
		// PUSH_LITERAL 0 (2)
		method.Method.Bytecodes = append(method.Method.Bytecodes, PUSH_LITERAL)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 0) // Index 0

		// SEND_MESSAGE 1 ("unknown") with 0 arguments
		method.Method.Bytecodes = append(method.Method.Bytecodes, SEND_MESSAGE)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 1) // Selector index 1
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 0) // 0 arguments

		// Create a context
		context := NewContext(method, vm.ObjectClass, []*Object{}, nil)

		// Execute the PUSH_LITERAL bytecode to set up the stack
		context.PC = 0
		if err := vm.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += InstructionSize(PUSH_LITERAL)

		// Execute the SEND_MESSAGE bytecode
		context.PC = 5
		_, err := vm.ExecuteSendMessage(context)

		// Check that we got an error
		if err == nil {
			t.Errorf("Expected error for unknown method, got nil")
		}
	})

	t.Run("nil receiver", func(t *testing.T) {
		t.Skip("Implement test for panic later")

		// Create a method with a SEND_MESSAGE bytecode for a nil receiver
		method := NewMethod(NewSymbol("test"), vm.ObjectClass)

		// Create literals
		nilObj := NewNil()
		selector := NewSymbol("test")

		// Add literals to the method
		method.Method.Literals = append(method.Method.Literals, nilObj)   // Index 0
		method.Method.Literals = append(method.Method.Literals, selector) // Index 1

		// Create bytecodes for: nil test
		// PUSH_LITERAL 0 (nil)
		method.Method.Bytecodes = append(method.Method.Bytecodes, PUSH_LITERAL)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 0) // Index 0

		// SEND_MESSAGE 1 ("test") with 0 arguments
		method.Method.Bytecodes = append(method.Method.Bytecodes, SEND_MESSAGE)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 1) // Selector index 1
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 0) // 0 arguments

		// Create a context
		context := NewContext(method, vm.ObjectClass, []*Object{}, nil)

		// Execute the PUSH_LITERAL bytecode to set up the stack
		context.PC = 0
		if err := vm.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += InstructionSize(PUSH_LITERAL)

		// Execute the SEND_MESSAGE bytecode
		context.PC = 5
		_, err := vm.ExecuteSendMessage(context)

		// Check that we got an error
		if err == nil {
			t.Errorf("Expected error for nil receiver, got nil")
		}
	})

	t.Run("invalid selector index", func(t *testing.T) {
		// Create a method with a SEND_MESSAGE bytecode with an invalid selector index
		method := NewMethod(NewSymbol("test"), vm.ObjectClass)

		// Create bytecodes for an invalid selector index
		// SEND_MESSAGE 999 with 0 arguments
		method.Method.Bytecodes = append(method.Method.Bytecodes, SEND_MESSAGE)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 3, 231) // Selector index 999
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 0)   // 0 arguments

		// Create a context
		context := NewContext(method, vm.ObjectClass, []*Object{}, nil)

		// Push a receiver onto the stack
		context.Push(vm.NewInteger(2))

		// Execute the SEND_MESSAGE bytecode
		context.PC = 0
		_, err := vm.ExecuteSendMessage(context)

		// Check that we got an error
		if err == nil {
			t.Errorf("Expected error for invalid selector index, got nil")
		}
	})

	t.Run("non-symbol selector", func(t *testing.T) {
		t.Skip("Skipping test until immediate values are fully implemented")
		// Create a method with a SEND_MESSAGE bytecode with a non-symbol selector
		method := NewMethod(NewSymbol("test"), vm.ObjectClass)

		// Create literals
		receiver := vm.NewInteger(2)
		nonSymbolSelector := vm.NewInteger(42)

		// Add literals to the method
		method.Method.Literals = append(method.Method.Literals, receiver)          // Index 0
		method.Method.Literals = append(method.Method.Literals, nonSymbolSelector) // Index 1

		// Create bytecodes for: 2 42 (where 42 is not a symbol)
		// PUSH_LITERAL 0 (2)
		method.Method.Bytecodes = append(method.Method.Bytecodes, PUSH_LITERAL)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 0) // Index 0

		// SEND_MESSAGE 1 (42) with 0 arguments
		method.Method.Bytecodes = append(method.Method.Bytecodes, SEND_MESSAGE)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 1) // Selector index 1
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 0) // 0 arguments

		// Create a context
		context := NewContext(method, vm.ObjectClass, []*Object{}, nil)

		// Execute the PUSH_LITERAL bytecode to set up the stack
		context.PC = 0
		if err := vm.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += InstructionSize(PUSH_LITERAL)

		// Execute the SEND_MESSAGE bytecode
		context.PC = 5
		_, err := vm.ExecuteSendMessage(context)

		// Check that we got an error
		if err == nil {
			t.Errorf("Expected error for non-symbol selector, got nil")
		}
	})
}

// TestExecuteSendMessageWithMultipleArguments tests the ExecuteSendMessage function with multiple arguments
func TestExecuteSendMessageWithMultipleArguments(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Test case for a method with multiple arguments
	t.Run("direct primitive call with multiple arguments", func(t *testing.T) {
		// Create a method with a SEND_MESSAGE bytecode for addition
		method := NewMethod(NewSymbol("test"), vm.ObjectClass)

		// Create literals
		twoObj := vm.NewInteger(2)
		threeObj := vm.NewInteger(3)
		fourObj := vm.NewInteger(4)
		plusSymbol := NewSymbol("+")

		// Add literals to the method
		method.Method.Literals = append(method.Method.Literals, twoObj)     // Index 0
		method.Method.Literals = append(method.Method.Literals, threeObj)   // Index 1
		method.Method.Literals = append(method.Method.Literals, fourObj)    // Index 2
		method.Method.Literals = append(method.Method.Literals, plusSymbol) // Index 3

		// Create bytecodes for: 2 + 3 + 4
		// PUSH_LITERAL 0 (2)
		method.Method.Bytecodes = append(method.Method.Bytecodes, PUSH_LITERAL)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 0) // Index 0

		// PUSH_LITERAL 1 (3)
		method.Method.Bytecodes = append(method.Method.Bytecodes, PUSH_LITERAL)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 1) // Index 1

		// SEND_MESSAGE 3 ("+") with 1 argument
		method.Method.Bytecodes = append(method.Method.Bytecodes, SEND_MESSAGE)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 3) // Selector index 3
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

		// PUSH_LITERAL 2 (4)
		method.Method.Bytecodes = append(method.Method.Bytecodes, PUSH_LITERAL)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 2) // Index 2

		// SEND_MESSAGE 3 ("+") with 1 argument
		method.Method.Bytecodes = append(method.Method.Bytecodes, SEND_MESSAGE)
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 3) // Selector index 3
		method.Method.Bytecodes = append(method.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

		// Create a context
		context := NewContext(method, vm.ObjectClass, []*Object{}, nil)

		// Execute the first PUSH_LITERAL bytecode
		context.PC = 0
		if err := vm.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += InstructionSize(PUSH_LITERAL)

		// Execute the second PUSH_LITERAL bytecode
		if err := vm.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += InstructionSize(PUSH_LITERAL)

		// Execute the first SEND_MESSAGE bytecode (2 + 3)
		result, err := vm.ExecuteSendMessage(context)
		if err != nil {
			t.Fatalf("Error executing SEND_MESSAGE: %s", err)
		}
		context.PC += InstructionSize(SEND_MESSAGE)

		// Check the intermediate result
		if IsIntegerImmediate(result) {
			intValue := GetIntegerImmediate(result)
			if intValue != 5 {
				t.Errorf("Expected 5, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %s", result)
		}

		// Execute the third PUSH_LITERAL bytecode
		if err := vm.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += InstructionSize(PUSH_LITERAL)

		// Execute the second SEND_MESSAGE bytecode (5 + 4)
		result, err = vm.ExecuteSendMessage(context)
		if err != nil {
			t.Fatalf("Error executing SEND_MESSAGE: %s", err)
		}

		// Check the final result
		if IsIntegerImmediate(result) {
			intValue := GetIntegerImmediate(result)
			if intValue != 9 {
				t.Errorf("Expected 9, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %s", result)
		}

		// Check the stack
		if context.StackPointer != 1 {
			t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
		}

		// Check the stack top
		stackTop := context.Stack[0]
		if IsIntegerImmediate(stackTop) {
			intValue := GetIntegerImmediate(stackTop)
			if intValue != 9 {
				t.Errorf("Expected stack top to be 9, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %s", stackTop)
		}
	})
}
