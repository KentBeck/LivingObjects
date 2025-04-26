package main

import (
	"fmt"
)

func DemoFactorial() {
	fmt.Println("SmalltalkLSP Bytecode Interpreter Demo")

	// Create a VM
	vm := NewVM()

	integerClass := vm.IntegerClass
	objectClass := vm.ObjectClass

	// No need to get the method dictionary explicitly when using MethodBuilder

	// Create a simple factorial method that returns a hardcoded value
	factorialSelector := NewSymbol("factorial")

	// Create the method using MethodBuilder
	factorialMethod := NewMethodBuilder(integerClass).
		Selector("factorial").
		Go()

	// Create literals for the factorial method
	oneObj := vm.NewInteger(1)
	equalsSymbol := NewSymbol("=")
	minusSymbol := NewSymbol("-")
	timesSymbol := NewSymbol("*")

	// Add literals to the factorial method
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, oneObj)            // Literal 0: 1
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, factorialSelector) // Literal 1: factorial
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, equalsSymbol)      // Literal 2: =
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, minusSymbol)       // Literal 3: -
	factorialMethod.Method.Literals = append(factorialMethod.Method.Literals, timesSymbol)       // Literal 4: *

	// Create bytecodes for factorial:
	// ^ self = 1
	//   ifTrue: [1]
	//   ifFalse: [self * (self - 1) factorial]

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

	// JUMP_IF_FALSE to next comparison
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, JUMP_IF_FALSE)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 7) // Jump 7 bytes ahead

	// Then branch: return 1
	// PUSH_LITERAL 0 (1)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_LITERAL)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// RETURN_STACK_TOP
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, RETURN_STACK_TOP)

	// For any other value, compute factorial recursively
	// PUSH_SELF (for later use in multiplication)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, PUSH_SELF)

	// DUPLICATE (to save a copy for later multiplication)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, DUPLICATE)

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

	// SEND_MESSAGE * with 1 argument
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, SEND_MESSAGE)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 4) // Selector index 4 (*)
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

	// RETURN_STACK_TOP
	factorialMethod.Method.Bytecodes = append(factorialMethod.Method.Bytecodes, RETURN_STACK_TOP)

	// Create a method to compute factorial of 4
	mainMethod := NewMethodBuilder(objectClass).
		Selector("main").
		Go()

	// Add literals to the main method
	fourObj := vm.NewInteger(4)
	factorialSelectorObj := NewSymbol("factorial")                                        // Create the selector for use in literals
	mainMethod.Method.Literals = append(mainMethod.Method.Literals, fourObj)              // Literal 0: 4
	mainMethod.Method.Literals = append(mainMethod.Method.Literals, factorialSelectorObj) // Literal 1: factorial

	// Create bytecodes for main: 4 factorial
	// PUSH_LITERAL 0 (4)
	mainMethod.Method.Bytecodes = append(mainMethod.Method.Bytecodes, PUSH_LITERAL)
	mainMethod.Method.Bytecodes = append(mainMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// SEND_MESSAGE 1 ("factorial") with 0 arguments
	mainMethod.Method.Bytecodes = append(mainMethod.Method.Bytecodes, SEND_MESSAGE)
	mainMethod.Method.Bytecodes = append(mainMethod.Method.Bytecodes, 0, 0, 0, 1) // Selector index 1
	mainMethod.Method.Bytecodes = append(mainMethod.Method.Bytecodes, 0, 0, 0, 0) // 0 arguments

	// RETURN_STACK_TOP
	mainMethod.Method.Bytecodes = append(mainMethod.Method.Bytecodes, RETURN_STACK_TOP)

	// Print the bytecodes for debugging
	fmt.Println("\nFactorial method bytecodes:")
	for i := 0; i < len(factorialMethod.Method.Bytecodes); i++ {
		if i%5 == 0 {
			fmt.Printf("\n%3d: ", i)
		}
		fmt.Printf("%3d ", factorialMethod.Method.Bytecodes[i])
	}

	// Create a context for the main method
	vm.CurrentContext = NewContext(mainMethod, fourObj, []*Object{}, nil)

	// Execute the VM
	result, err := vm.Execute()
	if err != nil {
		fmt.Printf("Error executing: %s\n", err)
		return
	}

	fmt.Printf("Result: %s\n", result)
}
