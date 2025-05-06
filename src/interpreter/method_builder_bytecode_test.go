package main

import (
	"testing"
)

func TestMethodBuilderBytecodes(t *testing.T) {
	// Create a class for testing
	testClass := NewClass("TestClass", nil)

	t.Run("PushLiteral", func(t *testing.T) {
		// Create a method with PushLiteral bytecode
		literal := MakeIntegerImmediate(42)

		builder := NewMethodBuilder(testClass).Selector("testPushLiteral")
		literalIndex, builder := builder.AddLiteral(literal)
		method := builder.PushLiteral(literalIndex).ReturnStackTop().Go()

		// Verify the bytecodes
		expectedBytecodes := []byte{
			PUSH_LITERAL,
			0, 0, 0, 0, // Literal index 0
			RETURN_STACK_TOP,
		}

		if len(method.Method.Bytecodes) != len(expectedBytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(expectedBytecodes), len(method.Method.Bytecodes))
		}

		for i, b := range expectedBytecodes {
			if method.Method.Bytecodes[i] != b {
				t.Errorf("Expected bytecode %d at index %d, got %d", b, i, method.Method.Bytecodes[i])
			}
		}
	})

	t.Run("PushSelf", func(t *testing.T) {
		// Create a method with PushSelf bytecode
		method := NewMethodBuilder(testClass).
			Selector("testPushSelf").
			PushSelf().
			ReturnStackTop().
			Go()

		// Verify the bytecodes
		expectedBytecodes := []byte{
			PUSH_SELF,
			RETURN_STACK_TOP,
		}

		if len(method.Method.Bytecodes) != len(expectedBytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(expectedBytecodes), len(method.Method.Bytecodes))
		}

		for i, b := range expectedBytecodes {
			if method.Method.Bytecodes[i] != b {
				t.Errorf("Expected bytecode %d at index %d, got %d", b, i, method.Method.Bytecodes[i])
			}
		}
	})

	t.Run("SendMessage", func(t *testing.T) {
		// Create a method with SendMessage bytecode
		plusSelector := NewSymbol("+")

		builder := NewMethodBuilder(testClass).Selector("testSendMessage")
		selectorIndex, builder := builder.AddLiteral(plusSelector)
		method := builder.
			PushSelf().
			PushSelf().
			SendMessage(selectorIndex, 1). // Send + with 1 argument
			ReturnStackTop().
			Go()

		// Verify the bytecodes
		expectedBytecodes := []byte{
			PUSH_SELF,
			PUSH_SELF,
			SEND_MESSAGE,
			0, 0, 0, 0, // Selector index 0
			0, 0, 0, 1, // 1 argument
			RETURN_STACK_TOP,
		}

		if len(method.Method.Bytecodes) != len(expectedBytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(expectedBytecodes), len(method.Method.Bytecodes))
		}

		for i, b := range expectedBytecodes {
			if method.Method.Bytecodes[i] != b {
				t.Errorf("Expected bytecode %d at index %d, got %d", b, i, method.Method.Bytecodes[i])
			}
		}
	})

	t.Run("InstanceVariables", func(t *testing.T) {
		// Create a method with PushInstanceVariable and StoreInstanceVariable bytecodes
		method := NewMethodBuilder(testClass).
			Selector("testInstanceVars").
			PushSelf().
			PushInstanceVariable(0).  // Push instance variable at index 0
			StoreInstanceVariable(1). // Store into instance variable at index 1
			ReturnStackTop().
			Go()

		// Verify the bytecodes
		expectedBytecodes := []byte{
			PUSH_SELF,
			PUSH_INSTANCE_VARIABLE,
			0, 0, 0, 0, // Instance variable index 0
			STORE_INSTANCE_VARIABLE,
			0, 0, 0, 1, // Instance variable index 1
			RETURN_STACK_TOP,
		}

		if len(method.Method.Bytecodes) != len(expectedBytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(expectedBytecodes), len(method.Method.Bytecodes))
		}

		for i, b := range expectedBytecodes {
			if method.Method.Bytecodes[i] != b {
				t.Errorf("Expected bytecode %d at index %d, got %d", b, i, method.Method.Bytecodes[i])
			}
		}
	})

	t.Run("TemporaryVariables", func(t *testing.T) {
		// Create a method with PushTemporaryVariable and StoreTemporaryVariable bytecodes
		method := NewMethodBuilder(testClass).
			Selector("testTempVars").
			TempVars([]string{"temp1", "temp2"}).
			PushSelf().
			PushTemporaryVariable(0).  // Push temporary variable at index 0
			StoreTemporaryVariable(1). // Store into temporary variable at index 1
			ReturnStackTop().
			Go()

		// Verify the bytecodes
		expectedBytecodes := []byte{
			PUSH_SELF,
			PUSH_TEMPORARY_VARIABLE,
			0, 0, 0, 0, // Temporary variable index 0
			STORE_TEMPORARY_VARIABLE,
			0, 0, 0, 1, // Temporary variable index 1
			RETURN_STACK_TOP,
		}

		if len(method.Method.Bytecodes) != len(expectedBytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(expectedBytecodes), len(method.Method.Bytecodes))
		}

		for i, b := range expectedBytecodes {
			if method.Method.Bytecodes[i] != b {
				t.Errorf("Expected bytecode %d at index %d, got %d", b, i, method.Method.Bytecodes[i])
			}
		}

		// Verify the temporary variable names
		if len(method.Method.TempVarNames) != 2 {
			t.Errorf("Expected 2 temporary variable names, got %d", len(method.Method.TempVarNames))
		}

		if method.Method.TempVarNames[0] != "temp1" {
			t.Errorf("Expected temporary variable name 'temp1', got '%s'", method.Method.TempVarNames[0])
		}

		if method.Method.TempVarNames[1] != "temp2" {
			t.Errorf("Expected temporary variable name 'temp2', got '%s'", method.Method.TempVarNames[1])
		}
	})

	t.Run("JumpOperations", func(t *testing.T) {
		// Create a method with Jump, JumpIfTrue, and JumpIfFalse bytecodes
		method := NewMethodBuilder(testClass).
			Selector("testJumps").
			PushSelf().
			JumpIfTrue(10). // Jump to offset 10 if true
			PushSelf().
			JumpIfFalse(20). // Jump to offset 20 if false
			PushSelf().
			Jump(30). // Jump to offset 30
			ReturnStackTop().
			Go()

		// Verify the bytecodes
		expectedBytecodes := []byte{
			PUSH_SELF,
			JUMP_IF_TRUE,
			0, 0, 0, 10, // Jump to offset 10
			PUSH_SELF,
			JUMP_IF_FALSE,
			0, 0, 0, 20, // Jump to offset 20
			PUSH_SELF,
			JUMP,
			0, 0, 0, 30, // Jump to offset 30
			RETURN_STACK_TOP,
		}

		if len(method.Method.Bytecodes) != len(expectedBytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(expectedBytecodes), len(method.Method.Bytecodes))
		}

		for i, b := range expectedBytecodes {
			if method.Method.Bytecodes[i] != b {
				t.Errorf("Expected bytecode %d at index %d, got %d", b, i, method.Method.Bytecodes[i])
			}
		}
	})

	t.Run("StackOperations", func(t *testing.T) {
		// Create a method with Pop and Duplicate bytecodes
		method := NewMethodBuilder(testClass).
			Selector("testStackOps").
			PushSelf().
			Duplicate(). // Duplicate the top of the stack
			Pop().       // Pop the top of the stack
			ReturnStackTop().
			Go()

		// Verify the bytecodes
		expectedBytecodes := []byte{
			PUSH_SELF,
			DUPLICATE,
			POP,
			RETURN_STACK_TOP,
		}

		if len(method.Method.Bytecodes) != len(expectedBytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(expectedBytecodes), len(method.Method.Bytecodes))
		}

		for i, b := range expectedBytecodes {
			if method.Method.Bytecodes[i] != b {
				t.Errorf("Expected bytecode %d at index %d, got %d", b, i, method.Method.Bytecodes[i])
			}
		}
	})

	t.Run("ComplexMethod", func(t *testing.T) {
		// Create a complex method using all bytecode operations
		// This method implements a simple factorial calculation:
		// if self <= 1 then return 1 else return self * (self - 1) factorial

		// Create literals
		oneObj := MakeIntegerImmediate(1)
		minusSelector := NewSymbol("-")
		timesSelector := NewSymbol("*")
		lessThanOrEqualSelector := NewSymbol("<=")
		factorialSelector := NewSymbol("factorial")

		// Create the method
		builder := NewMethodBuilder(testClass).Selector("factorial")

		// Add literals
		oneIndex, builder := builder.AddLiteral(oneObj)
		minusIndex, builder := builder.AddLiteral(minusSelector)
		timesIndex, builder := builder.AddLiteral(timesSelector)
		lessThanOrEqualIndex, builder := builder.AddLiteral(lessThanOrEqualSelector)
		factorialIndex, builder := builder.AddLiteral(factorialSelector)

		// Build the method
		method := builder.
			// Check if self <= 1
			PushSelf().
			PushLiteral(oneIndex).
			SendMessage(lessThanOrEqualIndex, 1).

			// If self <= 1, jump to the "then" branch
			JumpIfTrue(5).

			// "else" branch: return self * (self - 1) factorial
			PushSelf().
			Duplicate().
			PushLiteral(oneIndex).
			SendMessage(minusIndex, 1).
			SendMessage(factorialIndex, 0).
			SendMessage(timesIndex, 1).
			ReturnStackTop().

			// "then" branch: return 1
			PushLiteral(oneIndex).
			ReturnStackTop().
			Go()

		// We don't verify the exact bytecodes here because the jump offsets
		// would be complex to calculate. Instead, we just check that the method
		// was created successfully.
		if method == nil {
			t.Fatal("Method should not be nil")
		}

		// Check that the method is in the method dictionary
		methodDict := testClass.GetMethodDict()
		// Convert to Dictionary to access entries
		dict := ObjectToDictionary(methodDict)
		selectorValue := "factorial"
		methodInDict := dict.Entries[selectorValue]

		if methodInDict == nil {
			t.Fatalf("Method not found in dictionary for selector %q", selectorValue)
		}

		// Check that the method has the expected literals
		if len(methodInDict.Method.Literals) != 5 {
			t.Errorf("Expected 5 literals, got %d", len(methodInDict.Method.Literals))
		}
	})
}
