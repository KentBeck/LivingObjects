package main

import (
	"testing"
)

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

			// Create a simple method that returns a value
			returnValueMethod := NewMethodBuilder(integerClass).
				Selector("returnValue").
				Go()

			// Create a literal for the method
			valueObj := vm.NewInteger(42)

			// Add the literal to the method
			returnValueMethod.Method.Literals = append(returnValueMethod.Method.Literals, valueObj) // Literal 0: 42

			// Create bytecodes for the method: just return 42
			// PUSH_LITERAL 0 (42)
			returnValueMethod.Method.Bytecodes = append(returnValueMethod.Method.Bytecodes, PUSH_LITERAL)
			returnValueMethod.Method.Bytecodes = append(returnValueMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

			// RETURN_STACK_TOP
			returnValueMethod.Method.Bytecodes = append(returnValueMethod.Method.Bytecodes, RETURN_STACK_TOP)

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

			// Create a method that calls the addition method
			testMethod := NewMethodBuilder(integerClass).
				Selector("test").
				Go()

			// Create literals for the test method
			fiveObj := vm.NewInteger(5)
			tenObj := vm.NewInteger(10)

			// Add literals to the test method
			testMethod.Method.Literals = append(testMethod.Method.Literals, fiveObj)        // Literal 0: 5
			testMethod.Method.Literals = append(testMethod.Method.Literals, tenObj)         // Literal 1: 10
			testMethod.Method.Literals = append(testMethod.Method.Literals, NewSymbol("+")) // Literal 2: +

			// Create bytecodes for the test method: 5 + 10
			// PUSH_LITERAL 0 (5)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, PUSH_LITERAL)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

			// PUSH_LITERAL 1 (10)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, PUSH_LITERAL)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 1) // Index 1

			// SEND_MESSAGE 2 ("+") with 1 argument
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, SEND_MESSAGE)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 2) // Selector index 2
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

			// RETURN_STACK_TOP
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, RETURN_STACK_TOP)

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

			// Create a simple method that adds multiple numbers
			testMethod := NewMethodBuilder(integerClass).
				Selector("test").
				Go()

			// Create literals for the test method
			oneObj := vm.NewInteger(1)
			twoObj := vm.NewInteger(2)
			threeObj := vm.NewInteger(3)
			fourObj := vm.NewInteger(4)
			fiveObj := vm.NewInteger(5)

			// Add literals to the test method
			testMethod.Method.Literals = append(testMethod.Method.Literals, oneObj)         // Literal 0: 1
			testMethod.Method.Literals = append(testMethod.Method.Literals, twoObj)         // Literal 1: 2
			testMethod.Method.Literals = append(testMethod.Method.Literals, threeObj)       // Literal 2: 3
			testMethod.Method.Literals = append(testMethod.Method.Literals, fourObj)        // Literal 3: 4
			testMethod.Method.Literals = append(testMethod.Method.Literals, fiveObj)        // Literal 4: 5
			testMethod.Method.Literals = append(testMethod.Method.Literals, NewSymbol("+")) // Literal 5: +

			// Create bytecodes for the test method: 1 + 2 + 3 + 4 + 5
			// PUSH_LITERAL 0 (1)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, PUSH_LITERAL)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

			// PUSH_LITERAL 1 (2)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, PUSH_LITERAL)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 1) // Index 1

			// SEND_MESSAGE + with 1 argument
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, SEND_MESSAGE)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 5) // Selector index 5 (+)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

			// PUSH_LITERAL 2 (3)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, PUSH_LITERAL)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 2) // Index 2

			// SEND_MESSAGE + with 1 argument
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, SEND_MESSAGE)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 5) // Selector index 5 (+)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

			// PUSH_LITERAL 3 (4)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, PUSH_LITERAL)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 3) // Index 3

			// SEND_MESSAGE + with 1 argument
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, SEND_MESSAGE)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 5) // Selector index 5 (+)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

			// PUSH_LITERAL 4 (5)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, PUSH_LITERAL)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 4) // Index 4

			// SEND_MESSAGE + with 1 argument
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, SEND_MESSAGE)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 5) // Selector index 5 (+)
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

			// RETURN_STACK_TOP
			testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, RETURN_STACK_TOP)

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
