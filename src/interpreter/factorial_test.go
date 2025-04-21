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
	// if self = 1 { return 1 } else { return self * (self - 1) factorial }

	// Check if self = 1
	// PUSH_SELF
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_SELF)

	// PUSH_LITERAL 0 (1)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_LITERAL)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// SEND_MESSAGE = with 1 argument
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, SEND_MESSAGE)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 2) // Selector index 2 (=)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

	// JUMP_IF_FALSE to else branch
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, JUMP_IF_FALSE)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 7) // Jump 7 bytes ahead

	// Then branch: return 1
	// PUSH_LITERAL 0 (1)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_LITERAL)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// RETURN_STACK_TOP
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, RETURN_STACK_TOP)

	// Else branch: return self * (self - 1) factorial
	// PUSH_SELF (for later use in multiplication)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_SELF)

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

	// SEND_MESSAGE * with 1 argument
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, SEND_MESSAGE)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 4) // Selector index 4 (*)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

	// Add a debug message
	fmt.Printf("Added * message with selector index 4\n")

	// RETURN_STACK_TOP
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, RETURN_STACK_TOP)

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
