package main

import (
	"testing"
)

func TestFactorial(t *testing.T) {
	// Create a VM
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
	oneIndex, _ := builder.AddLiteral(oneObj)                  // Literal 0: 1
	factorialIndex, _ := builder.AddLiteral(factorialSelector) // Literal 1: factorial
	equalsIndex, _ := builder.AddLiteral(equalsSelector)       // Literal 2: =
	minusIndex, _ := builder.AddLiteral(minusSelector)         // Literal 3: -
	timesIndex, _ := builder.AddLiteral(timesSelector)         // Literal 4: *

	// Finalize the method
	factorialMethod := builder.Go()

	// Create bytecodes for factorial:
	// ^ self = 1
	//   ifTrue: [1]
	//   ifFalse: [self * (self - 1) factorial]

	// Create bytecodes array
	bytecodes := make([]byte, 0, 100) // Pre-allocate space for efficiency

	// First, compute self = 1
	// PUSH_SELF
	bytecodes = append(bytecodes, PUSH_SELF)

	// PUSH_LITERAL oneIndex (1)
	bytecodes = append(bytecodes, PUSH_LITERAL)
	bytecodes = append(bytecodes, 0, 0, 0, byte(oneIndex)) // Use the index from AddLiteral

	// SEND_MESSAGE = with 1 argument
	bytecodes = append(bytecodes, SEND_MESSAGE)
	bytecodes = append(bytecodes, 0, 0, 0, byte(equalsIndex)) // Use the index from AddLiteral
	bytecodes = append(bytecodes, 0, 0, 0, 1)                 // 1 argument

	// Now we have a boolean on the stack
	// We need to duplicate it for the two branches
	bytecodes = append(bytecodes, DUPLICATE)

	// JUMP_IF_FALSE to the false branch
	bytecodes = append(bytecodes, JUMP_IF_FALSE)
	bytecodes = append(bytecodes, 0, 0, 0, 12) // Jump past the true branch

	// True branch: [1]
	// POP the boolean (we don't need it anymore)
	bytecodes = append(bytecodes, POP)

	// PUSH_LITERAL oneIndex (1)
	bytecodes = append(bytecodes, PUSH_LITERAL)
	bytecodes = append(bytecodes, 0, 0, 0, byte(oneIndex)) // Use the index from AddLiteral

	// JUMP past the false branch to the return
	bytecodes = append(bytecodes, JUMP)
	bytecodes = append(bytecodes, 0, 0, 0, 35) // Jump to the return

	// False branch: [self * (self - 1) factorial]
	// POP the boolean (we don't need it anymore)
	bytecodes = append(bytecodes, POP)

	// PUSH_SELF (for later use in multiplication)
	bytecodes = append(bytecodes, PUSH_SELF)

	// Compute (self - 1) factorial
	// PUSH_SELF (for subtraction)
	bytecodes = append(bytecodes, PUSH_SELF)

	// PUSH_LITERAL oneIndex (1)
	bytecodes = append(bytecodes, PUSH_LITERAL)
	bytecodes = append(bytecodes, 0, 0, 0, byte(oneIndex)) // Use the index from AddLiteral

	// SEND_MESSAGE - with 1 argument
	bytecodes = append(bytecodes, SEND_MESSAGE)
	bytecodes = append(bytecodes, 0, 0, 0, byte(minusIndex)) // Use the index from AddLiteral
	bytecodes = append(bytecodes, 0, 0, 0, 1)                // 1 argument

	// SEND_MESSAGE factorial with 0 arguments
	bytecodes = append(bytecodes, SEND_MESSAGE)
	bytecodes = append(bytecodes, 0, 0, 0, byte(factorialIndex)) // Use the index from AddLiteral
	bytecodes = append(bytecodes, 0, 0, 0, 0)                    // 0 arguments

	// SEND_MESSAGE * with 1 argument (the factorial result is already on the stack)
	bytecodes = append(bytecodes, SEND_MESSAGE)
	bytecodes = append(bytecodes, 0, 0, 0, byte(timesIndex)) // Use the index from AddLiteral
	bytecodes = append(bytecodes, 0, 0, 0, 1)                // 1 argument

	// Return the result
	bytecodes = append(bytecodes, RETURN_STACK_TOP)

	// Add the bytecodes to the method
	factorialMethod.Method.Bytecodes = bytecodes

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
