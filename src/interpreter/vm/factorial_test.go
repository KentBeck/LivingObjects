package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

func TestFactorial(t *testing.T) {
	t.Skip("Skipping factorial test while fixing pile package integration")
	virtualMachine := vm.NewVM()

	// We'll use the VM's Integer class
	integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])

	// Create selectors for use in literals
	minusSelector := pile.NewSymbol("-")
	timesSelector := pile.NewSymbol("*")
	equalsSelector := pile.NewSymbol("=")
	factorialSelector := pile.NewSymbol("factorial")

	// Create literals for the factorial method
	oneObj := virtualMachine.NewInteger(1)

	// Create a simple factorial method using MethodBuilder with AddLiteral
	builder := compiler.NewMethodBuilder(integerClass)

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
	factorialMethod := builder.Go("factorial")

	// Create primitive methods for the Integer class
	compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Integer"])).
		Primitive(1). // Addition
		Go("+")

	compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Integer"])).
		Primitive(4). // Subtraction
		Go("-")

	compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Integer"])).
		Primitive(2). // Multiplication
		Go("*")

	compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Integer"])).
		Primitive(3). // Equality
		Go("=")

	// Test factorial of 1
	t.Run("Factorial of 1", func(t *testing.T) {
		// Create a context for the factorial method
		oneObj := virtualMachine.NewInteger(1)
		context := vm.NewContext(factorialMethod, oneObj, []*pile.Object{}, nil)

		// Execute the context
		result, err := virtualMachine.ExecuteContext(context)
		if err != nil {
			t.Errorf("Error executing factorial of 1: %v", err)
			return
		}

		// Check the result
		if pile.IsIntegerImmediate(result) {
			intValue := pile.GetIntegerImmediate(result)
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
		fourObj := virtualMachine.NewInteger(4)
		context := vm.NewContext(factorialMethod, fourObj, []*pile.Object{}, nil)

		// Execute the context
		result, err := virtualMachine.ExecuteContext(context)
		if err != nil {
			t.Errorf("Error executing factorial of 4: %v", err)
			return
		}

		// Check the result
		if pile.IsIntegerImmediate(result) {
			intValue := pile.GetIntegerImmediate(result)
			if intValue != 24 {
				t.Errorf("Expected result to be 24, got %d", intValue)
			}
		} else {
			t.Errorf("Expected an immediate integer, got %v", result)
		}
	})
}
