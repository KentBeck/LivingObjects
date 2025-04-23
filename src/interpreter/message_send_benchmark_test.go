package main

import (
	"testing"
)

// BenchmarkMessageSend measures the performance of a simple message send
func BenchmarkMessageSend(b *testing.B) {
	b.ReportAllocs()
	// Create a VM
	vm := NewVM()

	// We'll use the VM's Integer class
	integerClass := vm.IntegerClass

	// Create the method dictionary for the Integer class
	integerMethodDict := integerClass.GetMethodDict()

	// Create a simple addition method
	addSelector := NewSymbol("+")
	addMethod := NewMethod(addSelector, integerClass)
	addMethod.Method.IsPrimitive = true
	addMethod.Method.PrimitiveIndex = 1 // Addition

	// Add the method to the Integer class
	integerMethodDict.Entries[GetSymbolValue(addSelector)] = addMethod

	// Create a simple method that adds two numbers
	testSelector := NewSymbol("test")
	testMethod := NewMethod(testSelector, integerClass)

	// Create literals for the test method
	oneObj := vm.NewInteger(1)
	twoObj := vm.NewInteger(2)

	// Add literals to the test method
	testMethod.Method.Literals = append(testMethod.Method.Literals, oneObj)      // Literal 0: 1
	testMethod.Method.Literals = append(testMethod.Method.Literals, twoObj)      // Literal 1: 2
	testMethod.Method.Literals = append(testMethod.Method.Literals, addSelector) // Literal 2: +

	// Create bytecodes for the test method: 1 + 2
	// PUSH_LITERAL 0 (1)
	testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, PUSH_LITERAL)
	testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// PUSH_LITERAL 1 (2)
	testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, PUSH_LITERAL)
	testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 1) // Index 1

	// SEND_MESSAGE + with 1 argument
	testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, SEND_MESSAGE)
	testMethod.Method.Bytecodes = append(testMethod.Method.Bytecodes, 0, 0, 0, 2) // Selector index 2 (+)
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

		// Verify the result (1 + 2 = 3)
		if IsIntegerImmediate(result) {
			intValue := GetIntegerImmediate(result)
			if intValue != 3 {
				b.Fatalf("Expected 3, got %d", intValue)
			}
		} else {
			b.Fatalf("Returned non-immediate value: %v", result)
		}
	}
}

// BenchmarkMultipleMessageSends measures the performance of multiple message sends
func BenchmarkMultipleMessageSends(b *testing.B) {
	b.ReportAllocs()
	// Create a VM
	vm := NewVM()

	// We'll use the VM's Integer class
	integerClass := vm.IntegerClass

	// Create the method dictionary for the Integer class
	integerMethodDict := integerClass.GetMethodDict()

	// Create a simple addition method
	addSelector := NewSymbol("+")
	addMethod := NewMethod(addSelector, integerClass)
	addMethod.Method.IsPrimitive = true
	addMethod.Method.PrimitiveIndex = 1 // Addition

	// Add the method to the Integer class
	integerMethodDict.Entries[GetSymbolValue(addSelector)] = addMethod

	// Create a simple method that adds multiple numbers
	testSelector := NewSymbol("test")
	testMethod := NewMethod(testSelector, integerClass)

	// Create literals for the test method
	oneObj := vm.NewInteger(1)
	twoObj := vm.NewInteger(2)
	threeObj := vm.NewInteger(3)
	fourObj := vm.NewInteger(4)
	fiveObj := vm.NewInteger(5)

	// Add literals to the test method
	testMethod.Method.Literals = append(testMethod.Method.Literals, oneObj)      // Literal 0: 1
	testMethod.Method.Literals = append(testMethod.Method.Literals, twoObj)      // Literal 1: 2
	testMethod.Method.Literals = append(testMethod.Method.Literals, threeObj)    // Literal 2: 3
	testMethod.Method.Literals = append(testMethod.Method.Literals, fourObj)     // Literal 3: 4
	testMethod.Method.Literals = append(testMethod.Method.Literals, fiveObj)     // Literal 4: 5
	testMethod.Method.Literals = append(testMethod.Method.Literals, addSelector) // Literal 5: +

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

		// Verify the result (1 + 2 + 3 + 4 + 5 = 15)
		if IsIntegerImmediate(result) {
			intValue := GetIntegerImmediate(result)
			if intValue != 15 {
				b.Fatalf("Expected 15, got %d", intValue)
			}
		} else {
			b.Fatalf("Returned non-immediate value: %v", result)
		}
	}
}
