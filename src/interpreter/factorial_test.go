package main

import (
	"testing"
)

func TestFactorial(t *testing.T) {
	vm := NewVM()

	// We'll use the VM's Integer class
	integerClass := vm.IntegerClass

	// Create selectors for use in literals
	minusSelector := NewSymbol("-")
	timesSelector := NewSymbol("*")
	equalsSelector := NewSymbol("=")
	factorialSelector := NewSymbol("factorial")

	// Create literals for the factorial method
	oneObj := vm.NewInteger(1)

	// Create a simple factorial method using MethodBuilder with AddLiteral
	builder := NewMethodBuilder(integerClass).Selector("factorial")

	// Add literals to the method builder
	oneIndex, builder := builder.AddLiteral(oneObj)                  // Literal 0: 1
	factorialIndex, builder := builder.AddLiteral(factorialSelector) // Literal 1: factorial
	equalsIndex, builder := builder.AddLiteral(equalsSelector)       // Literal 2: =
	minusIndex, builder := builder.AddLiteral(minusSelector)         // Literal 3: -
	timesIndex, builder := builder.AddLiteral(timesSelector)         // Literal 4: *

	// Create bytecodes for factorial:
	// ^ self = 1
	//   ifTrue: [1]
	//   ifFalse: [self * (self - 1) factorial]

	// First, compute self = 1
	builder.PushSelf()
	builder.PushLiteral(oneIndex)
	builder.SendMessage(equalsIndex, 1)

	// Now we have a boolean on the stack
	// We need to duplicate it for the two branches
	builder.Duplicate()

	// JUMP_IF_FALSE to the false branch
	builder.JumpIfFalse(12) // Jump past the true branch

	// True branch: [1]
	// POP the boolean (we don't need it anymore)
	builder.Pop()

	// PUSH_LITERAL oneIndex (1)
	builder.PushLiteral(oneIndex)

	// JUMP past the false branch to the return
	builder.Jump(35) // Jump to the return

	// False branch: [self * (self - 1) factorial]
	// POP the boolean (we don't need it anymore)
	builder.Pop()

	// PUSH_SELF (for later use in multiplication)
	builder.PushSelf()

	// Compute (self - 1) factorial
	// PUSH_SELF (for subtraction)
	builder.PushSelf()

	// PUSH_LITERAL oneIndex (1)
	builder.PushLiteral(oneIndex)

	// SEND_MESSAGE - with 1 argument
	builder.SendMessage(minusIndex, 1)

	// SEND_MESSAGE factorial with 0 arguments
	builder.SendMessage(factorialIndex, 0)

	// SEND_MESSAGE * with 1 argument (the factorial result is already on the stack)
	builder.SendMessage(timesIndex, 1)

	// Return the result
	builder.ReturnStackTop()

	// Finalize the method
	factorialMethod := builder.Go()

	// Test factorial of 1
	t.Run("Factorial of 1", func(t *testing.T) {
		// Create a context for the factorial method
		oneObj := vm.NewInteger(1)
		context := NewContext(factorialMethod, oneObj, []*Object{}, nil)

		// Execute the context
		result, err := vm.ExecuteContext(context)
		if err != nil {
			t.Errorf("Error executing factorial of 1: %v", err)
			return
		}

		// Check the result
		if IsIntegerImmediate(result) {
			intValue := GetIntegerImmediate(result)
			if intValue != 1 {
				t.Errorf("Expected result to be 1, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %v", result)
		}
	})

	// Test factorial of 4
	t.Run("Factorial of 4", func(t *testing.T) {
		// Create a context for the factorial method
		fourObj := vm.NewInteger(4)
		context := NewContext(factorialMethod, fourObj, []*Object{}, nil)

		// Execute the context
		result, err := vm.ExecuteContext(context)
		if err != nil {
			t.Errorf("Error executing factorial of 4: %v", err)
			return
		}

		// Check the result
		if IsIntegerImmediate(result) {
			intValue := GetIntegerImmediate(result)
			if intValue != 24 {
				t.Errorf("Expected result to be 24, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %v", result)
		}
	})
}
