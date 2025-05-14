package demo

import (
	"fmt"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

func runDemo() {
	fmt.Println("SmalltalkLSP Bytecode Interpreter Demo")

	virtualMachine := vm.NewVM()

	integerClass := virtualMachine.Classes.Get(vm.Integer)
	objectClass := virtualMachine.Classes.Get(vm.Object)

	factorialSelector := pile.NewSymbol("factorial")

	// Create literals for the factorial method
	oneObj := virtualMachine.NewInteger(1)
	equalsSymbol := pile.NewSymbol("=")
	minusSymbol := pile.NewSymbol("-")
	timesSymbol := pile.NewSymbol("*")

	// Create the method using MethodBuilder with AddLiteral
	builder := compiler.NewMethodBuilder(integerClass).Selector("factorial")

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

	// Create literals for the main method
	fourObj := virtualMachine.NewInteger(4)
	factorialSelectorObj := pile.NewSymbol("factorial") // Create the selector for use in literals

	// Create a method to compute factorial of 4 using AddLiteral
	mainBuilder := compiler.NewMethodBuilder(objectClass).Selector("main")

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
	method := pile.ObjectToMethod(factorialMethod)
	for i := 0; i < len(method.Bytecodes); i++ {
		if i%5 == 0 {
			fmt.Printf("\n%3d: ", i)
		}
		fmt.Printf("%3d ", method.Bytecodes[i])
	}

	context := vm.NewContext(mainMethod, fourObj, []*core.Object{}, nil)

	result, err := virtualMachine.ExecuteContext(context)
	if err != nil {
		fmt.Printf("Error executing: %s\n", err)
		return
	}

	fmt.Printf("Result: %s\n", result)
}

// RunFactorialDemo runs a factorial calculation on the Smalltalk VM
func RunFactorialDemo() {
	runDemo()
}