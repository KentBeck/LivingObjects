package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/bytecode"
	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

// TestExecuteSendMessageExtended tests the ExecuteSendMessage function with more complex scenarios
func TestExecuteSendMessageExtended(t *testing.T) {
	t.Skip("Skipping extended send message tests while fixing pile package integration")
	virtualMachine := vm.NewVM()

	// Test cases
	t.Run("primitive method", func(t *testing.T) {
		// Add primitive methods to the Integer class
		integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])
		plusSymbol := pile.NewSymbol("+")
		compiler.NewMethodBuilder(integerClass).
			Primitive(1). // Addition primitive
			Go("+")

		// Create literals
		twoObj := virtualMachine.NewInteger(2)
		threeObj := virtualMachine.NewInteger(3)

		// Create a method with a SEND_MESSAGE bytecode for addition using AddLiteral
		builder := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"]))

		// Add literals to the method builder
		twoIndex, builder := builder.AddLiteral(twoObj)      // Index 0
		threeIndex, builder := builder.AddLiteral(threeObj)  // Index 1
		plusIndex, builder := builder.AddLiteral(plusSymbol) // Index 2

		// Create bytecodes for: 2 + 3
		builder.PushLiteral(twoIndex)
		builder.PushLiteral(threeIndex)
		builder.SendMessage(plusIndex, 1)

		// Finalize the method
		method := builder.Go("test")

		// Create a context
		context := vm.NewContext(method, pile.ObjectToClass(virtualMachine.Globals["Object"]), []*pile.Object{}, nil)

		// Execute the PUSH_LITERAL bytecodes to set up the stack
		context.PC = 0
		if err := virtualMachine.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += bytecode.InstructionSize(bytecode.PUSH_LITERAL)

		context.PC = 5
		if err := virtualMachine.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += bytecode.InstructionSize(bytecode.PUSH_LITERAL)

		// Execute the SEND_MESSAGE bytecode
		context.PC = 10
		result, err := virtualMachine.ExecuteSendMessage(context)
		if err != nil {
			t.Fatalf("Error executing SEND_MESSAGE: %s", err)
		}

		// Check the result
		if pile.IsIntegerImmediate(result) {
			intValue := pile.GetIntegerImmediate(result)
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
		if pile.IsIntegerImmediate(stackTop) {
			intValue := pile.GetIntegerImmediate(stackTop)
			if intValue != 5 {
				t.Errorf("Expected stack top to be 5, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %s", stackTop)
		}
	})

	t.Run("non-primitive method", func(t *testing.T) {
		// We'll use the VM's Object and Integer classes
		objectClass := pile.ObjectToClass(virtualMachine.Globals["Object"])
		integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])

		// Create literals
		factorialSelector := pile.NewSymbol("factorial")
		oneObj := virtualMachine.NewInteger(1)
		fiveObj := virtualMachine.NewInteger(5)

		// Create a factorial method for Integer using AddLiteral
		factorialBuilder := compiler.NewMethodBuilder(integerClass)

		// Add literals to the factorial method builder
		oneIndex, factorialBuilder := factorialBuilder.AddLiteral(oneObj) // Index 0

		// Create bytecodes for factorial: if self = 1 { return 1 } else { ... }
		// For simplicity, we'll just return 1 for any input
		factorialBuilder.PushLiteral(oneIndex)
		factorialBuilder.ReturnStackTop()

		// Finalize the factorial method
		factorialBuilder.Go("factorial")

		// Create a method that calls factorial using AddLiteral
		testBuilder := compiler.NewMethodBuilder(objectClass)

		// Add literals to the test method builder
		fiveIndex, testBuilder := testBuilder.AddLiteral(fiveObj)                        // Index 0
		factorialSelectorIndex, testBuilder := testBuilder.AddLiteral(factorialSelector) // Index 1

		// Create bytecodes for: 5 factorial
		testBuilder.PushLiteral(fiveIndex)
		testBuilder.SendMessage(factorialSelectorIndex, 0)

		// Finalize the test method
		testMethod := testBuilder.Go("test")

		// Create a context
		context := vm.NewContext(testMethod, objectClass, []*pile.Object{}, nil)

		// Set the VM's object class
		// REMOVED: virtualMachine.Classes.Register(vm.Object, objectClass)
		virtualMachine.Globals["Object"] = pile.ClassToObject(objectClass)
		virtualMachine.Globals["Integer"] = pile.ClassToObject(integerClass)

		// Execute the PUSH_LITERAL bytecode to set up the stack
		context.PC = 0
		if err := virtualMachine.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += bytecode.InstructionSize(bytecode.PUSH_LITERAL)

		// Execute the SEND_MESSAGE bytecode
		context.PC = 5
		result, err := virtualMachine.ExecuteSendMessage(context)
		if err != nil {
			t.Fatalf("Error executing SEND_MESSAGE: %s", err)
		}

		// Check the result
		if pile.IsIntegerImmediate(result) {
			intValue := pile.GetIntegerImmediate(result)
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
		if pile.IsIntegerImmediate(stackTop) {
			intValue := pile.GetIntegerImmediate(stackTop)
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
		receiver := virtualMachine.NewInteger(2)
		unknownSelector := pile.NewSymbol("unknown")

		// Create a method with a SEND_MESSAGE bytecode for an unknown method using AddLiteral
		builder := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"]))

		// Add literals to the method builder
		receiverIndex, builder := builder.AddLiteral(receiver)               // Index 0
		unknownSelectorIndex, builder := builder.AddLiteral(unknownSelector) // Index 1

		// Create bytecodes for: 2 unknown
		builder.PushLiteral(receiverIndex)
		builder.SendMessage(unknownSelectorIndex, 0)

		// Finalize the method
		method := builder.Go("test")

		// Create a context
		context := vm.NewContext(method, pile.ObjectToClass(virtualMachine.Globals["Object"]), []*pile.Object{}, nil)

		// Execute the PUSH_LITERAL bytecode to set up the stack
		context.PC = 0
		if err := virtualMachine.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += bytecode.InstructionSize(bytecode.PUSH_LITERAL)

		// Execute the SEND_MESSAGE bytecode
		context.PC = 5
		_, err := virtualMachine.ExecuteSendMessage(context)

		// Check that we got an error
		if err == nil {
			t.Errorf("Expected error for unknown method, got nil")
		}
	})
}

// TestExecuteSendMessageWithMultipleArguments tests the ExecuteSendMessage function with multiple arguments
func TestExecuteSendMessageWithMultipleArguments(t *testing.T) {
	t.Skip("Skipping multiple argument send message tests while fixing pile package integration")
	virtualMachine := vm.NewVM()

	// Test case for a method with multiple arguments
	t.Run("direct primitive call with multiple arguments", func(t *testing.T) {
		// Add primitive methods to the Integer class
		integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])
		plusSymbol := pile.NewSymbol("+")
		compiler.NewMethodBuilder(integerClass).
			Primitive(1). // Addition primitive
			Go("+")

		// Create literals
		twoObj := virtualMachine.NewInteger(2)
		threeObj := virtualMachine.NewInteger(3)
		fourObj := virtualMachine.NewInteger(4)

		// Create a method with a SEND_MESSAGE bytecode for addition using AddLiteral
		builder := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"]))

		// Add literals to the method builder
		twoIndex, builder := builder.AddLiteral(twoObj)      // Index 0
		threeIndex, builder := builder.AddLiteral(threeObj)  // Index 1
		fourIndex, builder := builder.AddLiteral(fourObj)    // Index 2
		plusIndex, builder := builder.AddLiteral(plusSymbol) // Index 3

		// Create bytecodes for: 2 + 3 + 4
		builder.PushLiteral(twoIndex)
		builder.PushLiteral(threeIndex)
		builder.SendMessage(plusIndex, 1)
		builder.PushLiteral(fourIndex)
		builder.SendMessage(plusIndex, 1)

		// Finalize the method
		method := builder.Go("test")

		// Create a context
		context := vm.NewContext(method, pile.ObjectToClass(virtualMachine.Globals["Object"]), []*pile.Object{}, nil)

		// Execute the first PUSH_LITERAL bytecode
		context.PC = 0
		if err := virtualMachine.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += bytecode.InstructionSize(bytecode.PUSH_LITERAL)

		// Execute the second PUSH_LITERAL bytecode
		if err := virtualMachine.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += bytecode.InstructionSize(bytecode.PUSH_LITERAL)

		// Execute the first SEND_MESSAGE bytecode (2 + 3)
		result, err := virtualMachine.ExecuteSendMessage(context)
		if err != nil {
			t.Fatalf("Error executing SEND_MESSAGE: %s", err)
		}
		context.PC += bytecode.InstructionSize(bytecode.SEND_MESSAGE)

		// Check the intermediate result
		if pile.IsIntegerImmediate(result) {
			intValue := pile.GetIntegerImmediate(result)
			if intValue != 5 {
				t.Errorf("Expected 5, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %s", result)
		}

		// Execute the third PUSH_LITERAL bytecode
		if err := virtualMachine.ExecutePushLiteral(context); err != nil {
			t.Fatalf("Error executing PUSH_LITERAL: %s", err)
		}
		context.PC += bytecode.InstructionSize(bytecode.PUSH_LITERAL)

		// Execute the second SEND_MESSAGE bytecode (5 + 4)
		result, err = virtualMachine.ExecuteSendMessage(context)
		if err != nil {
			t.Fatalf("Error executing SEND_MESSAGE: %s", err)
		}

		// Check the final result
		if pile.IsIntegerImmediate(result) {
			intValue := pile.GetIntegerImmediate(result)
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
		if pile.IsIntegerImmediate(stackTop) {
			intValue := pile.GetIntegerImmediate(stackTop)
			if intValue != 9 {
				t.Errorf("Expected stack top to be 9, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %s", stackTop)
		}
	})
}