package vm_test

import (
	"smalltalklsp/interpreter/pile"
	"testing"
	"time"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/vm"
)

// factorialTestCases defines the test cases for factorial benchmarks
var factorialTestCases = []struct {
	name     string
	input    int64
	expected int64
}{
	{"Factorial1", 1, 1},
	{"Factorial2", 2, 2},
	{"Factorial4", 4, 24},
	// Larger factorials can be added if needed
}

// setupFactorialMethod creates a factorial method for benchmarking
// This implementation is taken directly from the working factorial_test.go
func setupFactorialMethod(virtualMachine *vm.VM) *pile.Object {
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

	return factorialMethod
}

// BenchmarkFactorial is a parameterized benchmark for factorial calculation
func BenchmarkFactorial(b *testing.B) {
	for _, tc := range factorialTestCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			virtualMachine := vm.NewVM()

			// Setup the factorial method
			factorialMethod := setupFactorialMethod(virtualMachine)

			// Create the argument
			argObj := virtualMachine.NewInteger(tc.input)

			// Reset the timer before the benchmark
			b.ResetTimer()

			// Start time for message sends calculation
			startTime := time.Now()

			// Number of message sends
			var messageSends int64

			// Run the benchmark
			for i := 0; i < b.N; i++ {
				// Create a context for the factorial method
				context := vm.NewContext(factorialMethod, argObj, []*pile.Object{}, nil)

				// Execute the factorial method
				result, err := virtualMachine.ExecuteContext(context)
				if err != nil {
					b.Fatalf("Error executing factorial method: %v", err)
				}

				// Verify the result
				if pile.IsIntegerImmediate(result) {
					intValue := pile.GetIntegerImmediate(result)
					if intValue != tc.expected {
						b.Fatalf("Expected factorial of %d to be %d, got %d", tc.input, tc.expected, intValue)
					}
				} else {
					b.Fatalf("Expected an immediate integer, got %v", result)
				}

				// Count the number of message sends
				// For factorial of n, we have:
				// - 1 call to factorial
				// - n recursive calls (including the base case)
				// - Each call involves 3 message sends (=, -, and *) except the base case which only uses =
				// So total message sends = 1 + n + (n * 3 - 2) = 4*n - 1
				messageSends += 4*tc.input - 1
			}

			// Report message sends per second
			endTime := time.Now()
			duration := endTime.Sub(startTime)
			if duration.Seconds() > 0 {
				messagesSendsPerSecond := float64(messageSends) / duration.Seconds()
				b.ReportMetric(messagesSendsPerSecond, "sends/sec")
			}
		})
	}
}

// messageSendTestCases defines the test cases for message send benchmarks
var messageSendTestCases = []struct {
	name     string
	setup    func(*vm.VM) (*pile.Object, *pile.Object)
	expected int64
}{
	{
		name: "SimpleReturn",
		setup: func(virtualMachine *vm.VM) (*pile.Object, *pile.Object) {
			// We'll use the VM's Integer class
			integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])

			// Create a literal for the method
			valueObj := virtualMachine.NewInteger(42)

			// Create a simple method that returns a value using AddLiteral
			builder := compiler.NewMethodBuilder(integerClass)
			valueIndex, builder := builder.AddLiteral(valueObj) // Literal 0: 42

			// Create bytecodes for the method: just return 42
			builder.PushLiteral(valueIndex)
			builder.ReturnStackTop()

			// Finalize the method
			returnValueMethod := builder.Go("returnValue")

			// Create a receiver for the method
			receiver := virtualMachine.NewInteger(5)

			return returnValueMethod, receiver
		},
		expected: 42,
	},
	{
		name: "Addition",
		setup: func(virtualMachine *vm.VM) (*pile.Object, *pile.Object) {
			// We'll use the VM's Integer class
			integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])

			// Create a simple addition method
			compiler.NewMethodBuilder(integerClass).
				Primitive(1). // Addition
				Go("+")

			// Create literals for the test method
			fiveObj := virtualMachine.NewInteger(5)
			tenObj := virtualMachine.NewInteger(10)
			plusSymbol := pile.NewSymbol("+")

			// Create a method that calls the addition method using AddLiteral
			builder := compiler.NewMethodBuilder(integerClass)

			// Add literals to the method builder
			fiveIndex, builder := builder.AddLiteral(fiveObj)    // Literal 0: 5
			tenIndex, builder := builder.AddLiteral(tenObj)      // Literal 1: 10
			plusIndex, builder := builder.AddLiteral(plusSymbol) // Literal 2: +

			// Create bytecodes for the test method: 5 + 10
			builder.PushLiteral(fiveIndex)
			builder.PushLiteral(tenIndex)
			builder.SendMessage(plusIndex, 1)
			builder.ReturnStackTop()

			// Finalize the method
			testMethod := builder.Go("test")

			// Create a receiver for the test method
			receiver := virtualMachine.NewInteger(0)

			return testMethod, receiver
		},
		expected: 15,
	},
	{
		name: "MultipleAdditions",
		setup: func(virtualMachine *vm.VM) (*pile.Object, *pile.Object) {
			// We'll use the VM's Integer class
			integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])

			// Create a simple addition method
			compiler.NewMethodBuilder(integerClass).
				Primitive(1). // Addition
				Go("+")

			// Create literals for the test method
			oneObj := virtualMachine.NewInteger(1)
			twoObj := virtualMachine.NewInteger(2)
			threeObj := virtualMachine.NewInteger(3)
			fourObj := virtualMachine.NewInteger(4)
			fiveObj := virtualMachine.NewInteger(5)
			plusSymbol := pile.NewSymbol("+")

			// Create a method that adds multiple numbers using AddLiteral
			builder := compiler.NewMethodBuilder(integerClass)

			// Add literals to the method builder
			oneIndex, builder := builder.AddLiteral(oneObj)      // Literal 0: 1
			twoIndex, builder := builder.AddLiteral(twoObj)      // Literal 1: 2
			threeIndex, builder := builder.AddLiteral(threeObj)  // Literal 2: 3
			fourIndex, builder := builder.AddLiteral(fourObj)    // Literal 3: 4
			fiveIndex, builder := builder.AddLiteral(fiveObj)    // Literal 4: 5
			plusIndex, builder := builder.AddLiteral(plusSymbol) // Literal 5: +

			// Create bytecodes for the test method: 1 + 2 + 3 + 4 + 5
			builder.PushLiteral(oneIndex)
			builder.PushLiteral(twoIndex)
			builder.SendMessage(plusIndex, 1)

			builder.PushLiteral(threeIndex)
			builder.SendMessage(plusIndex, 1)

			builder.PushLiteral(fourIndex)
			builder.SendMessage(plusIndex, 1)

			builder.PushLiteral(fiveIndex)
			builder.SendMessage(plusIndex, 1)

			builder.ReturnStackTop()

			// Finalize the method
			testMethod := builder.Go("test")

			// Create a receiver for the test method
			receiver := virtualMachine.NewInteger(0)

			return testMethod, receiver
		},
		expected: 15,
	},
}

// BenchmarkMessageSend is a parameterized benchmark for message sending
func BenchmarkMessageSend(b *testing.B) {
	for _, tc := range messageSendTestCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			virtualMachine := vm.NewVM()

			// Setup the test method and receiver
			testMethod, receiver := tc.setup(virtualMachine)

			// Reset the timer before the benchmark
			b.ResetTimer()

			// Run the benchmark
			for i := 0; i < b.N; i++ {
				// Create a context for the test method
				context := vm.NewContext(testMethod, receiver, []*pile.Object{}, nil)

				// Execute the test method
				result, err := virtualMachine.ExecuteContext(context)
				if err != nil {
					b.Fatalf("Error executing test method: %v", err)
				}

				// Verify the result
				if pile.IsIntegerImmediate(result) {
					intValue := pile.GetIntegerImmediate(result)
					if intValue != tc.expected {
						b.Fatalf("Expected %d, got %d", tc.expected, intValue)
					}
				} else {
					b.Fatalf("Expected an immediate integer, got %v", result)
				}
			}
		})
	}
}
