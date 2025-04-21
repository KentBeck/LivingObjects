package main

import (
	"fmt"
	"testing"
)

func TestFactorial(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// We'll use the VM's Integer class
	integerClass := vm.IntegerClass

	// Create the method dictionary for the Integer class
	integerMethodDict := integerClass.GetMethodDict()

	// + method
	plusSelector := NewSymbol("+")
	plusMethod := NewMethod(plusSelector, integerClass)
	plusMethod.Method.IsPrimitive = true
	plusMethod.Method.PrimitiveIndex = 1 // Addition
	integerMethodDict.Entries[plusSelector.SymbolValue] = plusMethod

	// - method
	minusSelector := NewSymbol("-")
	minusMethod := NewMethod(minusSelector, integerClass)
	minusMethod.Method.IsPrimitive = true
	minusMethod.Method.PrimitiveIndex = 4 // Subtraction
	integerMethodDict.Entries[minusSelector.SymbolValue] = minusMethod

	// * method
	timesSelector := NewSymbol("*")
	timesMethod := NewMethod(timesSelector, integerClass)
	timesMethod.Method.IsPrimitive = true
	timesMethod.Method.PrimitiveIndex = 2 // Multiplication
	integerMethodDict.Entries[timesSelector.SymbolValue] = timesMethod

	// = method
	equalsSelector := NewSymbol("=")
	equalsMethod := NewMethod(equalsSelector, integerClass)
	equalsMethod.Method.IsPrimitive = true
	equalsMethod.Method.PrimitiveIndex = 3 // Equality
	integerMethodDict.Entries[equalsSelector.SymbolValue] = equalsMethod

	// Create a simple factorial method
	factorialSelector := NewSymbol("factorial")
	factorialMethod := NewMethod(factorialSelector, integerClass)

	// Add the factorial method to the Integer class
	integerMethodDict.Entries[factorialSelector.SymbolValue] = factorialMethod

	// Create literals for the factorial method
	oneObj := vm.NewIntegerWithClass(1)

	// Add literals to the factorial method
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, oneObj)            // Literal 0: 1
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, factorialSelector) // Literal 1: factorial
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, equalsSelector)    // Literal 2: =
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, minusSelector)     // Literal 3: -
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, timesSelector)     // Literal 4: *

	// Add debug output
	fmt.Printf("Literals: %v\n", factorialMethod.Method.Literals)

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

	// Add a debug message
	fmt.Printf("Added JUMP with offset 35 at PC %d\n", len(factorialMethod.Method.Bytecodes)-5)

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

	// Add a debug message
	fmt.Printf("Added factorial message with selector index 1\n")

	// SEND_MESSAGE * with 1 argument (the factorial result is already on the stack)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, SEND_MESSAGE)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 4) // Selector index 4 (*)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

	// Add a debug message
	fmt.Printf("Added * message with selector index 4\n")

	// Return the result
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, RETURN_STACK_TOP)

	// Add a debug message
	fmt.Printf("Added RETURN_STACK_TOP at PC %d\n", len(factorialMethod.Method.Bytecodes)-1)

	// Test factorial of 1
	t.Run("Factorial of 1", func(t *testing.T) {
		// Create a context for the factorial method
		oneObj := vm.NewIntegerWithClass(1)
		context := NewContext(factorialMethod, oneObj, []*Object{}, nil)

		// Execute the context
		result, err := vm.ExecuteContext(context)
		if err != nil {
			t.Errorf("Error executing factorial of 1: %v", err)
			return
		}

		// Check the result
		if result.Type != OBJ_INTEGER {
			t.Errorf("Expected result to be an integer, got %v", result.Type)
		}

		if result.IntegerValue != 1 {
			t.Errorf("Expected result to be 1, got %d", result.IntegerValue)
		}
	})

	// Test factorial of 4
	t.Run("Factorial of 4", func(t *testing.T) {
		// Create a context for the factorial method
		fourObj := vm.NewIntegerWithClass(4)
		context := NewContext(factorialMethod, fourObj, []*Object{}, nil)

		// Execute the context
		fmt.Printf("Executing factorial of 4...\n")
		result, err := vm.ExecuteContext(context)
		if err != nil {
			t.Errorf("Error executing factorial of 4: %v", err)
			return
		}

		// Check the result
		if result.Type != OBJ_INTEGER {
			t.Errorf("Expected result to be an integer, got %v", result.Type)
		}

		if result.IntegerValue != 24 {
			t.Errorf("Expected result to be 24, got %d", result.IntegerValue)
		}
	})
}
