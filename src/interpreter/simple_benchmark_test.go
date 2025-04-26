package main

import (
	"testing"
)

// BenchmarkSimpleMessageSend measures the performance of simple message sends
func BenchmarkSimpleMessageSend(b *testing.B) {
	b.ReportAllocs()
	// Create a VM
	vm := NewVM()

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

	// Reset the timer before the benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Create a context for the method
		context := NewContext(returnValueMethod, receiver, []*Object{}, nil)

		// Execute the method
		result, err := vm.ExecuteContext(context)
		if err != nil {
			b.Fatalf("Error executing method: %v", err)
		}

		// Verify the result
		if IsIntegerImmediate(result) {
			intValue := GetIntegerImmediate(result)
			if intValue != 42 {
				b.Fatalf("Expected 42, got %d", intValue)
			}
		} else {
			b.Fatalf("Expected an immediate integer, got %v", result)
		}
	}
}

// BenchmarkAddition measures the performance of integer addition
func BenchmarkAddition(b *testing.B) {
	b.ReportAllocs()
	// Create a VM
	vm := NewVM()

	// We'll use the VM's Integer class
	integerClass := vm.IntegerClass

	// Create a simple method that adds two numbers
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
			if intValue != 15 {
				b.Fatalf("Expected 15, got %d", intValue)
			}
		} else {
			b.Fatalf("Expected an immediate integer, got %v", result)
		}
	}
}
