package main

import (
	"testing"
)

func TestFactorial(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// We'll use the VM's Integer class
	integerClass := vm.IntegerClass

	// Create the method dictionary for the Integer class
	integerMethodDict := integerClass.GetMethodDict()

	minusSelector := NewSymbol("-")
	timesSelector := NewSymbol("*")
	equalsSelector := NewSymbol("=")

	// Create a simple factorial method
	factorialSelector := NewSymbol("factorial")
	factorialMethod := NewMethod(factorialSelector, integerClass)

	// Add the factorial method to the Integer class
	integerMethodDict.Entries[factorialSelector.SymbolValue] = factorialMethod

	// Create literals for the factorial method
	oneObj := vm.NewInteger(1)

	// Add literals to the factorial method
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, oneObj)            // Literal 0: 1
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, factorialSelector) // Literal 1: factorial
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, equalsSelector)    // Literal 2: =
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, minusSelector)     // Literal 3: -
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, timesSelector)     // Literal 4: *

	// Create bytecodes for factorial:
	// ^ self = 1
	//   ifTrue: [1]
	//   ifFalse: [self * (self - 1) factorial]

	// First, compute self = 1
	// PUSH_SELF
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_SELF)

	// PUSH_LITERAL 0 (1)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_LITERAL)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// SEND_MESSAGE = with 1 argument
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, SEND_MESSAGE)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 2) // Selector index 2 (=)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

	// Now we have a boolean on the stack
	// We need to duplicate it for the two branches
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, DUPLICATE)

	// JUMP_IF_FALSE to the false branch
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, JUMP_IF_FALSE)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 12) // Jump past the true branch

	// True branch: [1]
	// POP the boolean (we don't need it anymore)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, POP)

	// PUSH_LITERAL 0 (1)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_LITERAL)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// JUMP past the false branch to the return
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, JUMP)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 35) // Jump to the return

	// False branch: [self * (self - 1) factorial]
	// POP the boolean (we don't need it anymore)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, POP)

	// PUSH_SELF (for later use in multiplication)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_SELF)

	// Compute (self - 1) factorial
	// PUSH_SELF (for subtraction)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_SELF)

	// PUSH_LITERAL 0 (1)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_LITERAL)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// SEND_MESSAGE - with 1 argument
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, SEND_MESSAGE)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 3) // Selector index 3 (-)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

	// SEND_MESSAGE factorial with 0 arguments
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, SEND_MESSAGE)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 1) // Selector index 1 (factorial)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 0) // 0 arguments

	// SEND_MESSAGE * with 1 argument (the factorial result is already on the stack)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, SEND_MESSAGE)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 4) // Selector index 4 (*)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

	// Return the result
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, RETURN_STACK_TOP)

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
		} else if result.Type != OBJ_INTEGER {
			t.Errorf("Expected result to be an integer, got %v", result.Type)
		} else if result.IntegerValue != 1 {
			t.Errorf("Expected result to be 1, got %d", result.IntegerValue)
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
		} else if result.Type != OBJ_INTEGER {
			t.Errorf("Expected result to be an integer, got %v", result.Type)
		} else if result.IntegerValue != 24 {
			t.Errorf("Expected result to be 24, got %d", result.IntegerValue)
		}
	})
}
