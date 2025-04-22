package main

import (
	"testing"
	"time"
)

// BenchmarkFactorial10 measures the performance of factorial calculation for 10
func BenchmarkFactorial10(b *testing.B) {
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

	// Create the argument (10)
	tenObj := vm.NewInteger(10)

	// Reset the timer before the benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create a context for the factorial method
		context := NewContext(factorialMethod, tenObj, []*Object{}, nil)

		// Execute the factorial method
		result, err := vm.ExecuteContext(context)
		if err != nil {
			b.Fatalf("Error executing factorial method: %v", err)
		}

		// Verify the result (factorial of 10 = 3628800)
		if IsIntegerImmediate(result) {
			intValue := GetIntegerImmediate(result)
			if intValue != 3628800 {
				b.Fatalf("Expected factorial of 10 to be 3628800, got %d", intValue)
			}
		} else if result.Type != OBJ_INTEGER || result.IntegerValue != 3628800 {
			b.Fatalf("Expected factorial of 10 to be 3628800, got %v", result)
		}
	}
}

// BenchmarkFactorial19 measures the performance of factorial calculation for 19
func BenchmarkFactorial19(b *testing.B) {
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

	// Create the argument (19)
	nineteenObj := vm.NewInteger(19)

	// Reset the timer before the benchmark
	b.ResetTimer()

	// Start time
	startTime := time.Now()

	// Number of message sends
	var messageSends int64

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create a context for the factorial method
		context := NewContext(factorialMethod, nineteenObj, []*Object{}, nil)

		// Execute the factorial method
		result, err := vm.ExecuteContext(context)
		if err != nil {
			b.Fatalf("Error executing factorial method: %v", err)
		}

		// Verify the result (factorial of 19 = 121645100408832000)
		if IsIntegerImmediate(result) {
			intValue := GetIntegerImmediate(result)
			if intValue != 121645100408832000 {
				b.Fatalf("Expected factorial of 19 to be 121645100408832000, got %d", intValue)
			}
		} else if result.Type != OBJ_INTEGER || result.IntegerValue != 121645100408832000 {
			b.Fatalf("Expected factorial of 19 to be 121645100408832000, got %v", result)
		}

		// Count the number of message sends
		// For factorial of 19, we have:
		// - 1 call to factorial
		// - 19 recursive calls (including the base case)
		// - Each call involves 3 message sends (=, -, and *) except the base case which only uses =
		// So total message sends = 1 + 19 + (19 * 3 - 2) = 76
		messageSends += 76
	}

	// End time
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// Calculate message sends per second
	messagesSendsPerSecond := float64(messageSends) / duration.Seconds()

	// Report the result
	b.ReportMetric(messagesSendsPerSecond, "sends/sec")
}
