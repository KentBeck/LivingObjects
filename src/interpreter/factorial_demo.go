package main

import (
	"fmt"
)

func DemoFactorial() {
	fmt.Println("SmalltalkLSP Bytecode Interpreter Demo")

	vm := NewVM()

	integerClass := vm.IntegerClass
	objectClass := vm.ObjectClass

	factorialSelector := NewSymbol("factorial")

	// Create literals for the factorial method
	oneObj := vm.NewInteger(1)
	equalsSymbol := NewSymbol("=")
	minusSymbol := NewSymbol("-")
	timesSymbol := NewSymbol("*")

	// Create the method using MethodBuilder with AddLiteral
	builder := NewMethodBuilder(integerClass).Selector("factorial")

	// Add literals to the method builder
	oneIndex, _ := builder.AddLiteral(oneObj)                  // Literal 0: 1
	factorialIndex, _ := builder.AddLiteral(factorialSelector) // Literal 1: factorial
	equalsIndex, _ := builder.AddLiteral(equalsSymbol)         // Literal 2: =
	minusIndex, _ := builder.AddLiteral(minusSymbol)           // Literal 3: -
	timesIndex, _ := builder.AddLiteral(timesSymbol)           // Literal 4: *

	// Finalize the method
	factorialMethod := builder.Go()

	// Create bytecodes for factorial:
	// ^ self = 1
	//   ifTrue: [1]
	//   ifFalse: [self * (self - 1) factorial]

	// Create bytecodes array
	bytecodes := make([]byte, 0, 100) // Pre-allocate space for efficiency

	// Check if self = 1
	// PUSH_SELF
	bytecodes = append(bytecodes, PUSH_SELF)

	// PUSH_LITERAL oneIndex (1)
	bytecodes = append(bytecodes, PUSH_LITERAL)
	bytecodes = append(bytecodes, 0, 0, 0, byte(oneIndex)) // Use the index from AddLiteral

	// SEND_MESSAGE = with 1 argument
	bytecodes = append(bytecodes, SEND_MESSAGE)
	bytecodes = append(bytecodes, 0, 0, 0, byte(equalsIndex)) // Use the index from AddLiteral
	bytecodes = append(bytecodes, 0, 0, 0, 1)                 // 1 argument

	// JUMP_IF_FALSE to next comparison
	bytecodes = append(bytecodes, JUMP_IF_FALSE)
	bytecodes = append(bytecodes, 0, 0, 0, 7) // Jump 7 bytes ahead

	// Then branch: return 1
	// PUSH_LITERAL oneIndex (1)
	bytecodes = append(bytecodes, PUSH_LITERAL)
	bytecodes = append(bytecodes, 0, 0, 0, byte(oneIndex)) // Use the index from AddLiteral

	// RETURN_STACK_TOP
	bytecodes = append(bytecodes, RETURN_STACK_TOP)

	// For any other value, compute factorial recursively
	// PUSH_SELF (for later use in multiplication)
	bytecodes = append(bytecodes, PUSH_SELF)

	// DUPLICATE (to save a copy for later multiplication)
	bytecodes = append(bytecodes, DUPLICATE)

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

	// SEND_MESSAGE * with 1 argument
	bytecodes = append(bytecodes, SEND_MESSAGE)
	bytecodes = append(bytecodes, 0, 0, 0, byte(timesIndex)) // Use the index from AddLiteral
	bytecodes = append(bytecodes, 0, 0, 0, 1)                // 1 argument

	// RETURN_STACK_TOP
	bytecodes = append(bytecodes, RETURN_STACK_TOP)

	// Add the bytecodes to the method
	factorialMethod.Method.Bytecodes = bytecodes

	// Create literals for the main method
	fourObj := vm.NewInteger(4)
	factorialSelectorObj := NewSymbol("factorial") // Create the selector for use in literals

	// Create a method to compute factorial of 4 using AddLiteral
	mainBuilder := NewMethodBuilder(objectClass).Selector("main")

	// Add literals to the method builder
	fourIndex, _ := mainBuilder.AddLiteral(fourObj)                           // Literal 0: 4
	factorialSelectorIndex, _ := mainBuilder.AddLiteral(factorialSelectorObj) // Literal 1: factorial

	// Create bytecodes for main: 4 factorial
	mainBytecodes := make([]byte, 0, 20)

	// PUSH_LITERAL fourIndex (4)
	mainBytecodes = append(mainBytecodes, PUSH_LITERAL)
	mainBytecodes = append(mainBytecodes, 0, 0, 0, byte(fourIndex))

	// SEND_MESSAGE factorialSelectorIndex ("factorial") with 0 arguments
	mainBytecodes = append(mainBytecodes, SEND_MESSAGE)
	mainBytecodes = append(mainBytecodes, 0, 0, 0, byte(factorialSelectorIndex))
	mainBytecodes = append(mainBytecodes, 0, 0, 0, 0) // 0 arguments

	// RETURN_STACK_TOP
	mainBytecodes = append(mainBytecodes, RETURN_STACK_TOP)

	// Finalize the method
	mainMethod := mainBuilder.Bytecodes(mainBytecodes).Go()

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
