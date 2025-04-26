package main

import (
	"testing"
	"time"
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
func setupFactorialMethod(vm *VM) *Object {
	// We'll use the VM's Integer class
	integerClass := vm.IntegerClass

	// Create selectors for use in literals
	minusSelector := NewSymbol("-")
	timesSelector := NewSymbol("*")
	equalsSelector := NewSymbol("=")
	factorialSelector := NewSymbol("factorial")

	// Create literals for the factorial method
	oneObj := vm.NewInteger(1)

	// Create a simple factorial method using MethodBuilder with AddLiteral
	builder := NewMethodBuilder(integerClass).Selector("factorial")

	// Add literals to the method builder
	oneIndex, _ := builder.AddLiteral(oneObj)                  // Literal 0: 1
	factorialIndex, _ := builder.AddLiteral(factorialSelector) // Literal 1: factorial
	equalsIndex, _ := builder.AddLiteral(equalsSelector)       // Literal 2: =
	minusIndex, _ := builder.AddLiteral(minusSelector)         // Literal 3: -
	timesIndex, _ := builder.AddLiteral(timesSelector)         // Literal 4: *

	// Finalize the method
	factorialMethod := builder.Go()

	// Create bytecodes for factorial:
	// ^ self = 1
	//   ifTrue: [1]
	//   ifFalse: [self * (self - 1) factorial]

	// Create bytecodes array
	bytecodes := make([]byte, 0, 100) // Pre-allocate space for efficiency

	// First, compute self = 1
	// PUSH_SELF
	bytecodes = append(bytecodes, PUSH_SELF)

	// PUSH_LITERAL oneIndex (1)
	bytecodes = append(bytecodes, PUSH_LITERAL)
	bytecodes = append(bytecodes, 0, 0, 0, byte(oneIndex)) // Use the index from AddLiteral

	// SEND_MESSAGE = with 1 argument
	bytecodes = append(bytecodes, SEND_MESSAGE)
	bytecodes = append(bytecodes, 0, 0, 0, byte(equalsIndex)) // Use the index from AddLiteral
	bytecodes = append(bytecodes, 0, 0, 0, 1)                 // 1 argument

	// Now we have a boolean on the stack
	// We need to duplicate it for the two branches
	bytecodes = append(bytecodes, DUPLICATE)

	// JUMP_IF_FALSE to the false branch
	bytecodes = append(bytecodes, JUMP_IF_FALSE)
	bytecodes = append(bytecodes, 0, 0, 0, 12) // Jump past the true branch

	// True branch: [1]
	// POP the boolean (we don't need it anymore)
	bytecodes = append(bytecodes, POP)

	// PUSH_LITERAL oneIndex (1)
	bytecodes = append(bytecodes, PUSH_LITERAL)
	bytecodes = append(bytecodes, 0, 0, 0, byte(oneIndex)) // Use the index from AddLiteral

	// JUMP past the false branch to the return
	bytecodes = append(bytecodes, JUMP)
	bytecodes = append(bytecodes, 0, 0, 0, 35) // Jump to the return

	// False branch: [self * (self - 1) factorial]
	// POP the boolean (we don't need it anymore)
	bytecodes = append(bytecodes, POP)

	// PUSH_SELF (for later use in multiplication)
	bytecodes = append(bytecodes, PUSH_SELF)

	// Compute (self - 1) factorial
	// PUSH_SELF (for subtraction)
	bytecodes = append(bytecodes, PUSH_SELF)

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

	// SEND_MESSAGE * with 1 argument (the factorial result is already on the stack)
	bytecodes = append(bytecodes, SEND_MESSAGE)
	bytecodes = append(bytecodes, 0, 0, 0, byte(timesIndex)) // Use the index from AddLiteral
	bytecodes = append(bytecodes, 0, 0, 0, 1)                // 1 argument

	// Return the result
	bytecodes = append(bytecodes, RETURN_STACK_TOP)

	// Add the bytecodes to the method
	factorialMethod.Method.Bytecodes = bytecodes

	return factorialMethod
}

// BenchmarkFactorial is a parameterized benchmark for factorial calculation
func BenchmarkFactorial(b *testing.B) {
	for _, tc := range factorialTestCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			// Create a VM
			vm := NewVM()

			// Setup the factorial method
			factorialMethod := setupFactorialMethod(vm)

			// Create the argument
			argObj := vm.NewInteger(tc.input)

			// Reset the timer before the benchmark
			b.ResetTimer()

			// Start time for message sends calculation
			startTime := time.Now()

			// Number of message sends
			var messageSends int64

			// Run the benchmark
			for i := 0; i < b.N; i++ {
				// Create a context for the factorial method
				context := NewContext(factorialMethod, argObj, []*Object{}, nil)

				// Execute the factorial method
				result, err := vm.ExecuteContext(context)
				if err != nil {
					b.Fatalf("Error executing factorial method: %v", err)
				}

				// Verify the result
				if IsIntegerImmediate(result) {
					intValue := GetIntegerImmediate(result)
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
	setup    func(*VM) (*Object, *Object)
	expected int64
}{
	{
		name: "SimpleReturn",
		setup: func(vm *VM) (*Object, *Object) {
			// We'll use the VM's Integer class
			integerClass := vm.IntegerClass

			// Create a literal for the method
			valueObj := vm.NewInteger(42)

			// Create a simple method that returns a value using AddLiteral
			builder := NewMethodBuilder(integerClass).Selector("returnValue")
			valueIndex, _ := builder.AddLiteral(valueObj) // Literal 0: 42

			// Create bytecodes for the method: just return 42
			bytecodes := make([]byte, 0, 10)

			// PUSH_LITERAL valueIndex (42)
			bytecodes = append(bytecodes, PUSH_LITERAL)
			bytecodes = append(bytecodes, 0, 0, 0, byte(valueIndex))

			// RETURN_STACK_TOP
			bytecodes = append(bytecodes, RETURN_STACK_TOP)

			// Finalize the method
			returnValueMethod := builder.Bytecodes(bytecodes).Go()

			// Create a receiver for the method
			receiver := vm.NewInteger(5)

			return returnValueMethod, receiver
		},
		expected: 42,
	},
	{
		name: "Addition",
		setup: func(vm *VM) (*Object, *Object) {
			// We'll use the VM's Integer class
			integerClass := vm.IntegerClass

			// Create a simple addition method
			NewMethodBuilder(integerClass).
				Selector("+").
				Primitive(1). // Addition
				Go()

			// Create literals for the test method
			fiveObj := vm.NewInteger(5)
			tenObj := vm.NewInteger(10)
			plusSymbol := NewSymbol("+")

			// Create a method that calls the addition method using AddLiteral
			builder := NewMethodBuilder(integerClass).Selector("test")

			// Add literals to the method builder
			fiveIndex, _ := builder.AddLiteral(fiveObj)    // Literal 0: 5
			tenIndex, _ := builder.AddLiteral(tenObj)      // Literal 1: 10
			plusIndex, _ := builder.AddLiteral(plusSymbol) // Literal 2: +

			// Create bytecodes for the test method: 5 + 10
			bytecodes := make([]byte, 0, 20)

			// PUSH_LITERAL fiveIndex (5)
			bytecodes = append(bytecodes, PUSH_LITERAL)
			bytecodes = append(bytecodes, 0, 0, 0, byte(fiveIndex))

			// PUSH_LITERAL tenIndex (10)
			bytecodes = append(bytecodes, PUSH_LITERAL)
			bytecodes = append(bytecodes, 0, 0, 0, byte(tenIndex))

			// SEND_MESSAGE plusIndex ("+") with 1 argument
			bytecodes = append(bytecodes, SEND_MESSAGE)
			bytecodes = append(bytecodes, 0, 0, 0, byte(plusIndex))
			bytecodes = append(bytecodes, 0, 0, 0, 1) // 1 argument

			// RETURN_STACK_TOP
			bytecodes = append(bytecodes, RETURN_STACK_TOP)

			// Finalize the method
			testMethod := builder.Bytecodes(bytecodes).Go()

			// Create a receiver for the test method
			receiver := vm.NewInteger(0)

			return testMethod, receiver
		},
		expected: 15,
	},
	{
		name: "MultipleAdditions",
		setup: func(vm *VM) (*Object, *Object) {
			// We'll use the VM's Integer class
			integerClass := vm.IntegerClass

			// Create a simple addition method
			NewMethodBuilder(integerClass).
				Selector("+").
				Primitive(1). // Addition
				Go()

			// Create literals for the test method
			oneObj := vm.NewInteger(1)
			twoObj := vm.NewInteger(2)
			threeObj := vm.NewInteger(3)
			fourObj := vm.NewInteger(4)
			fiveObj := vm.NewInteger(5)
			plusSymbol := NewSymbol("+")

			// Create a method that adds multiple numbers using AddLiteral
			builder := NewMethodBuilder(integerClass).Selector("test")

			// Add literals to the method builder
			oneIndex, _ := builder.AddLiteral(oneObj)      // Literal 0: 1
			twoIndex, _ := builder.AddLiteral(twoObj)      // Literal 1: 2
			threeIndex, _ := builder.AddLiteral(threeObj)  // Literal 2: 3
			fourIndex, _ := builder.AddLiteral(fourObj)    // Literal 3: 4
			fiveIndex, _ := builder.AddLiteral(fiveObj)    // Literal 4: 5
			plusIndex, _ := builder.AddLiteral(plusSymbol) // Literal 5: +

			// Create bytecodes for the test method: 1 + 2 + 3 + 4 + 5
			bytecodes := make([]byte, 0, 50)

			// PUSH_LITERAL oneIndex (1)
			bytecodes = append(bytecodes, PUSH_LITERAL)
			bytecodes = append(bytecodes, 0, 0, 0, byte(oneIndex))

			// PUSH_LITERAL twoIndex (2)
			bytecodes = append(bytecodes, PUSH_LITERAL)
			bytecodes = append(bytecodes, 0, 0, 0, byte(twoIndex))

			// SEND_MESSAGE + with 1 argument
			bytecodes = append(bytecodes, SEND_MESSAGE)
			bytecodes = append(bytecodes, 0, 0, 0, byte(plusIndex))
			bytecodes = append(bytecodes, 0, 0, 0, 1) // 1 argument

			// PUSH_LITERAL threeIndex (3)
			bytecodes = append(bytecodes, PUSH_LITERAL)
			bytecodes = append(bytecodes, 0, 0, 0, byte(threeIndex))

			// SEND_MESSAGE + with 1 argument
			bytecodes = append(bytecodes, SEND_MESSAGE)
			bytecodes = append(bytecodes, 0, 0, 0, byte(plusIndex))
			bytecodes = append(bytecodes, 0, 0, 0, 1) // 1 argument

			// PUSH_LITERAL fourIndex (4)
			bytecodes = append(bytecodes, PUSH_LITERAL)
			bytecodes = append(bytecodes, 0, 0, 0, byte(fourIndex))

			// SEND_MESSAGE + with 1 argument
			bytecodes = append(bytecodes, SEND_MESSAGE)
			bytecodes = append(bytecodes, 0, 0, 0, byte(plusIndex))
			bytecodes = append(bytecodes, 0, 0, 0, 1) // 1 argument

			// PUSH_LITERAL fiveIndex (5)
			bytecodes = append(bytecodes, PUSH_LITERAL)
			bytecodes = append(bytecodes, 0, 0, 0, byte(fiveIndex))

			// SEND_MESSAGE + with 1 argument
			bytecodes = append(bytecodes, SEND_MESSAGE)
			bytecodes = append(bytecodes, 0, 0, 0, byte(plusIndex))
			bytecodes = append(bytecodes, 0, 0, 0, 1) // 1 argument

			// RETURN_STACK_TOP
			bytecodes = append(bytecodes, RETURN_STACK_TOP)

			// Finalize the method
			testMethod := builder.Bytecodes(bytecodes).Go()

			// Create a receiver for the test method
			receiver := vm.NewInteger(0)

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
			// Create a VM
			vm := NewVM()

			// Setup the test method and receiver
			testMethod, receiver := tc.setup(vm)

			// Reset the timer before the benchmark
			b.ResetTimer()

			// Run the benchmark
			for i := 0; i < b.N; i++ {
				// Create a context for the test method
				context := NewContext(testMethod, receiver, []*Object{}, nil)

				// Execute the test method
				result, err := vm.ExecuteContext(context)
				if err != nil {
					b.Fatalf("Error executing test method: %v", err)
				}

				// Verify the result
				if IsIntegerImmediate(result) {
					intValue := GetIntegerImmediate(result)
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
