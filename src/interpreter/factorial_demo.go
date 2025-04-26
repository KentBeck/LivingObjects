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
	oneIndex, builder := builder.AddLiteral(oneObj)                  // Literal 0: 1
	factorialIndex, builder := builder.AddLiteral(factorialSelector) // Literal 1: factorial
	equalsIndex, builder := builder.AddLiteral(equalsSymbol)         // Literal 2: =
	minusIndex, builder := builder.AddLiteral(minusSymbol)           // Literal 3: -
	timesIndex, builder := builder.AddLiteral(timesSymbol)           // Literal 4: *

	// Create bytecodes for factorial:
	// ^ self = 1
	//   ifTrue: [1]
	//   ifFalse: [self * (self - 1) factorial]

	// Check if self = 1
	builder.PushSelf()
	builder.PushLiteral(oneIndex)
	builder.SendMessage(equalsIndex, 1)

	// Jump to else branch if false
	builder.JumpIfFalse(7)

	// Then branch: return 1
	builder.PushLiteral(oneIndex)
	builder.ReturnStackTop()

	// Else branch: self * (self - 1) factorial
	builder.PushSelf()
	builder.Duplicate()
	builder.PushLiteral(oneIndex)
	builder.SendMessage(minusIndex, 1)
	builder.SendMessage(factorialIndex, 0)
	builder.SendMessage(timesIndex, 1)
	builder.ReturnStackTop()

	// Finalize the method
	factorialMethod := builder.Go()

	// Create literals for the main method
	fourObj := vm.NewInteger(4)
	factorialSelectorObj := NewSymbol("factorial") // Create the selector for use in literals

	// Create a method to compute factorial of 4 using AddLiteral
	mainBuilder := NewMethodBuilder(objectClass).Selector("main")

	// Add literals to the method builder
	fourIndex, mainBuilder := mainBuilder.AddLiteral(fourObj)                           // Literal 0: 4
	factorialSelectorIndex, mainBuilder := mainBuilder.AddLiteral(factorialSelectorObj) // Literal 1: factorial

	// Create bytecodes for main: 4 factorial
	mainBuilder.PushLiteral(fourIndex)
	mainBuilder.SendMessage(factorialSelectorIndex, 0)
	mainBuilder.ReturnStackTop()

	// Finalize the method
	mainMethod := mainBuilder.Go()

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
