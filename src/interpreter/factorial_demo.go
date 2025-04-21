package main

import (
	"fmt"
)

func DemoFactorial() {
	fmt.Println("SmalltalkLSP Bytecode Interpreter Demo")

	// Create a VM
	vm := NewVM()

	// Create basic classes
	objectClass := NewClass("Object", nil)
	vm.ObjectClass = objectClass
	vm.Globals["Object"] = objectClass

	// Create Integer class
	integerClass := NewClass("Integer", objectClass)
	vm.Globals["Integer"] = integerClass

	// Add primitive methods to the Integer class
	// + method
	plusSelector := NewSymbol("+")
	plusMethod := NewMethod(plusSelector, integerClass)
	plusMethod.Method.IsPrimitive = true
	plusMethod.Method.PrimitiveIndex = 1 // Addition
	integerMethodDict := integerClass.GetMethodDict()
	integerMethodDict.Entries[plusSelector.SymbolValue] = plusMethod

	// - method (subtraction)
	minusSelector := NewSymbol("-")
	minusMethod := NewMethod(minusSelector, integerClass)
	minusMethod.Method.IsPrimitive = true
	minusMethod.Method.PrimitiveIndex = 4 // Subtraction (new primitive)
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

	// Create a simple factorial method that returns a hardcoded value
	factorialSelector := NewSymbol("factorial")
	factorialMethod := NewMethod(factorialSelector, integerClass)

	// Add the factorial method to the Integer class
	integerMethodDict.Entries[factorialSelector.SymbolValue] = factorialMethod

	// Create literals for the factorial method
	oneObj := vm.NewIntegerWithClass(1, integerClass)
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
	// if self = 1 { return 1 }
	// else if self = 2 { return 2 }
	// else if self = 3 { return 6 }
	// else if self = 4 { return 24 }
	// else { return 0 } // Error case

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
	mainSelector := NewSymbol("main")
	mainMethod := NewMethod(mainSelector, objectClass)

	// Add literals to the main method
	fourObj := vm.NewIntegerWithClass(4, integerClass)
	mainMethod.Method.Literals = append(mainMethod.Method.Literals, fourObj)           // Literal 0: 4
	mainMethod.Method.Literals = append(mainMethod.Method.Literals, factorialSelector) // Literal 1: factorial

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
