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
		// Create literals
		twoObj := vm.NewInteger(2)
		threeObj := vm.NewInteger(3)
		plusSymbol := NewSymbol("+")

		// Create a method with a SEND_MESSAGE bytecode for addition using AddLiteral
		builder := NewMethodBuilder(vm.ObjectClass).Selector("test")

		// Add literals to the method builder
		twoIndex, _ := builder.AddLiteral(twoObj)      // Index 0
		threeIndex, _ := builder.AddLiteral(threeObj)  // Index 1
		plusIndex, _ := builder.AddLiteral(plusSymbol) // Index 2

		// Create bytecodes for: 2 + 3
		bytecodes := make([]byte, 0, 20)

		// PUSH_LITERAL twoIndex (2)
		bytecodes = append(bytecodes, PUSH_LITERAL)
		bytecodes = append(bytecodes, 0, 0, 0, byte(twoIndex))

		// PUSH_LITERAL threeIndex (3)
		bytecodes = append(bytecodes, PUSH_LITERAL)
		bytecodes = append(bytecodes, 0, 0, 0, byte(threeIndex))

		// SEND_MESSAGE plusIndex ("+") with 1 argument
		bytecodes = append(bytecodes, SEND_MESSAGE)
		bytecodes = append(bytecodes, 0, 0, 0, byte(plusIndex))
		bytecodes = append(bytecodes, 0, 0, 0, 1) // 1 argument

		// Finalize the method
		method := builder.Bytecodes(bytecodes).Go()

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

		// Create literals
		factorialSelector := NewSymbol("factorial")
		oneObj := vm.NewInteger(1)
		fiveObj := vm.NewInteger(5)

		// Create a factorial method for Integer using AddLiteral
		factorialBuilder := NewMethodBuilder(integerClass).Selector("factorial")

		// Add literals to the factorial method builder
		oneIndex, _ := factorialBuilder.AddLiteral(oneObj) // Index 0

		// Create bytecodes for factorial: if self = 1 { return 1 } else { ... }
		// For simplicity, we'll just return 1 for any input
		factorialBytecodes := make([]byte, 0, 10)

		// PUSH_LITERAL oneIndex (1)
		factorialBytecodes = append(factorialBytecodes, PUSH_LITERAL)
		factorialBytecodes = append(factorialBytecodes, 0, 0, 0, byte(oneIndex))

		// RETURN_STACK_TOP
		factorialBytecodes = append(factorialBytecodes, RETURN_STACK_TOP)

		// Finalize the factorial method
		factorialBuilder.Bytecodes(factorialBytecodes).Go()

		// Create a method that calls factorial using AddLiteral
		testBuilder := NewMethodBuilder(objectClass).Selector("test")

		// Add literals to the test method builder
		fiveIndex, _ := testBuilder.AddLiteral(fiveObj)                        // Index 0
		factorialSelectorIndex, _ := testBuilder.AddLiteral(factorialSelector) // Index 1

		// Create bytecodes for: 5 factorial
		testBytecodes := make([]byte, 0, 15)

		// PUSH_LITERAL fiveIndex (5)
		testBytecodes = append(testBytecodes, PUSH_LITERAL)
		testBytecodes = append(testBytecodes, 0, 0, 0, byte(fiveIndex))

		// SEND_MESSAGE factorialSelectorIndex ("factorial") with 0 arguments
		testBytecodes = append(testBytecodes, SEND_MESSAGE)
		testBytecodes = append(testBytecodes, 0, 0, 0, byte(factorialSelectorIndex))
		testBytecodes = append(testBytecodes, 0, 0, 0, 0) // 0 arguments

		// Finalize the test method
		testMethod := testBuilder.Bytecodes(testBytecodes).Go()

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

		// Create literals
		receiver := vm.NewInteger(2)
		unknownSelector := NewSymbol("unknown")

		// Create a method with a SEND_MESSAGE bytecode for an unknown method using AddLiteral
		builder := NewMethodBuilder(vm.ObjectClass).Selector("test")

		// Add literals to the method builder
		receiverIndex, _ := builder.AddLiteral(receiver)               // Index 0
		unknownSelectorIndex, _ := builder.AddLiteral(unknownSelector) // Index 1

		// Create bytecodes for: 2 unknown
		bytecodes := make([]byte, 0, 15)

		// PUSH_LITERAL receiverIndex (2)
		bytecodes = append(bytecodes, PUSH_LITERAL)
		bytecodes = append(bytecodes, 0, 0, 0, byte(receiverIndex))

		// SEND_MESSAGE unknownSelectorIndex ("unknown") with 0 arguments
		bytecodes = append(bytecodes, SEND_MESSAGE)
		bytecodes = append(bytecodes, 0, 0, 0, byte(unknownSelectorIndex))
		bytecodes = append(bytecodes, 0, 0, 0, 0) // 0 arguments

		// Finalize the method
		method := builder.Bytecodes(bytecodes).Go()

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

		// Create literals
		nilObj := NewNil()
		selector := NewSymbol("test")

		// Create a method with a SEND_MESSAGE bytecode for a nil receiver using AddLiteral
		builder := NewMethodBuilder(vm.ObjectClass).Selector("test")

		// Add literals to the method builder
		nilIndex, _ := builder.AddLiteral(nilObj)        // Index 0
		selectorIndex, _ := builder.AddLiteral(selector) // Index 1

		// Create bytecodes for: nil test
		bytecodes := make([]byte, 0, 15)

		// PUSH_LITERAL nilIndex (nil)
		bytecodes = append(bytecodes, PUSH_LITERAL)
		bytecodes = append(bytecodes, 0, 0, 0, byte(nilIndex))

		// SEND_MESSAGE selectorIndex ("test") with 0 arguments
		bytecodes = append(bytecodes, SEND_MESSAGE)
		bytecodes = append(bytecodes, 0, 0, 0, byte(selectorIndex))
		bytecodes = append(bytecodes, 0, 0, 0, 0) // 0 arguments

		// Finalize the method
		method := builder.Bytecodes(bytecodes).Go()

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
		method := NewMethodBuilder(vm.ObjectClass).
			Selector("test").
			Go()

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
		// Create literals
		receiver := vm.NewInteger(2)
		nonSymbolSelector := vm.NewInteger(42)

		// Create a method with a SEND_MESSAGE bytecode with a non-symbol selector using AddLiteral
		builder := NewMethodBuilder(vm.ObjectClass).Selector("test")

		// Add literals to the method builder
		receiverIndex, _ := builder.AddLiteral(receiver)                   // Index 0
		nonSymbolSelectorIndex, _ := builder.AddLiteral(nonSymbolSelector) // Index 1

		// Create bytecodes for: 2 42 (where 42 is not a symbol)
		bytecodes := make([]byte, 0, 15)

		// PUSH_LITERAL receiverIndex (2)
		bytecodes = append(bytecodes, PUSH_LITERAL)
		bytecodes = append(bytecodes, 0, 0, 0, byte(receiverIndex))

		// SEND_MESSAGE nonSymbolSelectorIndex (42) with 0 arguments
		bytecodes = append(bytecodes, SEND_MESSAGE)
		bytecodes = append(bytecodes, 0, 0, 0, byte(nonSymbolSelectorIndex))
		bytecodes = append(bytecodes, 0, 0, 0, 0) // 0 arguments

		// Finalize the method
		method := builder.Bytecodes(bytecodes).Go()

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
		// Create literals
		twoObj := vm.NewInteger(2)
		threeObj := vm.NewInteger(3)
		fourObj := vm.NewInteger(4)
		plusSymbol := NewSymbol("+")

		// Create a method with a SEND_MESSAGE bytecode for addition using AddLiteral
		builder := NewMethodBuilder(vm.ObjectClass).Selector("test")

		// Add literals to the method builder
		twoIndex, _ := builder.AddLiteral(twoObj)      // Index 0
		threeIndex, _ := builder.AddLiteral(threeObj)  // Index 1
		fourIndex, _ := builder.AddLiteral(fourObj)    // Index 2
		plusIndex, _ := builder.AddLiteral(plusSymbol) // Index 3

		// Create bytecodes for: 2 + 3 + 4
		bytecodes := make([]byte, 0, 30)

		// PUSH_LITERAL twoIndex (2)
		bytecodes = append(bytecodes, PUSH_LITERAL)
		bytecodes = append(bytecodes, 0, 0, 0, byte(twoIndex))

		// PUSH_LITERAL threeIndex (3)
		bytecodes = append(bytecodes, PUSH_LITERAL)
		bytecodes = append(bytecodes, 0, 0, 0, byte(threeIndex))

		// SEND_MESSAGE plusIndex ("+") with 1 argument
		bytecodes = append(bytecodes, SEND_MESSAGE)
		bytecodes = append(bytecodes, 0, 0, 0, byte(plusIndex))
		bytecodes = append(bytecodes, 0, 0, 0, 1) // 1 argument

		// PUSH_LITERAL fourIndex (4)
		bytecodes = append(bytecodes, PUSH_LITERAL)
		bytecodes = append(bytecodes, 0, 0, 0, byte(fourIndex))

		// SEND_MESSAGE plusIndex ("+") with 1 argument
		bytecodes = append(bytecodes, SEND_MESSAGE)
		bytecodes = append(bytecodes, 0, 0, 0, byte(plusIndex))
		bytecodes = append(bytecodes, 0, 0, 0, 1) // 1 argument

		// Finalize the method
		method := builder.Bytecodes(bytecodes).Go()

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
