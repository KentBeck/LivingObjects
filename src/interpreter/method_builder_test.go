package main

import (
	"testing"
)

func TestMethodBuilder(t *testing.T) {
	// Create a class for testing
	testClass := NewClass("TestClass", nil)

	t.Run("PrimitiveMethod", func(t *testing.T) {
		method := NewMethodBuilder(testClass).
			Selector("testPrimitive").
			Primitive(42).
			Go()

		// Verify the method was created correctly
		if method == nil {
			t.Fatal("Method should not be nil")
		}

		// Check that the method is in the method dictionary
		methodDict := testClass.GetMethodDict()
		// Convert to Dictionary to access entries
		dict := ObjectToDictionary(methodDict)
		selectorValue := "testPrimitive"
		methodInDict := dict.Entries[selectorValue]

		if methodInDict == nil {
			t.Fatalf("Method not found in dictionary for selector %q", selectorValue)
		}

		// Check the method properties
		if !methodInDict.Method.IsPrimitive {
			t.Error("Method should be marked as primitive")
		}

		if methodInDict.Method.PrimitiveIndex != 42 {
			t.Errorf("Expected primitive index 42, got %d", methodInDict.Method.PrimitiveIndex)
		}
	})

	t.Run("MethodWithBytecodes", func(t *testing.T) {
		// Create a method with bytecodes using MethodBuilder
		method := NewMethodBuilder(testClass).
			Selector("testMethod").
			PushSelf().
			ReturnStackTop().
			Go()

		// Verify the method was created correctly
		if method == nil {
			t.Fatal("Method should not be nil")
		}

		// Check that the method is in the method dictionary
		methodDict := testClass.GetMethodDict()
		// Convert to Dictionary to access entries
		dict := ObjectToDictionary(methodDict)
		selectorValue := "testMethod"
		methodInDict := dict.Entries[selectorValue]

		if methodInDict == nil {
			t.Fatalf("Method not found in dictionary for selector %q", selectorValue)
		}

		// Check the method properties
		expectedBytecodes := []byte{PUSH_SELF, RETURN_STACK_TOP}
		if len(methodInDict.Method.Bytecodes) != len(expectedBytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(expectedBytecodes), len(methodInDict.Method.Bytecodes))
		}

		for i, b := range expectedBytecodes {
			if methodInDict.Method.Bytecodes[i] != b {
				t.Errorf("Expected bytecode %d at index %d, got %d", b, i, methodInDict.Method.Bytecodes[i])
			}
		}
	})

	t.Run("MethodWithLiterals", func(t *testing.T) {
		// Create a method with literals using MethodBuilder
		literals := []*Object{MakeIntegerImmediate(123), NewSymbol("test")}
		method := NewMethodBuilder(testClass).
			Selector("testLiterals").
			AddLiterals(literals).
			Go()

		// Verify the method was created correctly
		if method == nil {
			t.Fatal("Method should not be nil")
		}

		// Check that the method is in the method dictionary
		methodDict := testClass.GetMethodDict()
		// Convert to Dictionary to access entries
		dict := ObjectToDictionary(methodDict)
		selectorValue := "testLiterals"
		methodInDict := dict.Entries[selectorValue]

		if methodInDict == nil {
			t.Fatalf("Method not found in dictionary for selector %q", selectorValue)
		}

		// Check the method properties
		if len(methodInDict.Method.Literals) != len(literals) {
			t.Errorf("Expected %d literals, got %d", len(literals), len(methodInDict.Method.Literals))
		}
	})

	t.Run("CompleteMethod", func(t *testing.T) {
		// Create a complete method using MethodBuilder
		literals := []*Object{MakeIntegerImmediate(123), NewSymbol("test")}
		tempVars := []string{"temp1", "temp2"}

		builder := NewMethodBuilder(testClass).
			Selector("completeMethod").
			Primitive(42).
			AddLiterals(literals).
			TempVars(tempVars)

		// Add bytecodes using the fluent API
		method := builder.
			PushSelf().
			ReturnStackTop().
			Go()

		// Verify the method was created correctly
		if method == nil {
			t.Fatal("Method should not be nil")
		}

		// Check that the method is in the method dictionary
		methodDict := testClass.GetMethodDict()
		// Convert to Dictionary to access entries
		dict := ObjectToDictionary(methodDict)
		selectorValue := "completeMethod"
		methodInDict := dict.Entries[selectorValue]

		if methodInDict == nil {
			t.Fatalf("Method not found in dictionary for selector %q", selectorValue)
		}

		// Check the method properties
		if !methodInDict.Method.IsPrimitive {
			t.Error("Method should be marked as primitive")
		}

		if methodInDict.Method.PrimitiveIndex != 42 {
			t.Errorf("Expected primitive index 42, got %d", methodInDict.Method.PrimitiveIndex)
		}

		expectedBytecodes := []byte{PUSH_SELF, RETURN_STACK_TOP}
		if len(methodInDict.Method.Bytecodes) != len(expectedBytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(expectedBytecodes), len(methodInDict.Method.Bytecodes))
		}

		if len(methodInDict.Method.Literals) != len(literals) {
			t.Errorf("Expected %d literals, got %d", len(literals), len(methodInDict.Method.Literals))
		}

		if len(methodInDict.Method.TempVarNames) != len(tempVars) {
			t.Errorf("Expected %d temp vars, got %d", len(tempVars), len(methodInDict.Method.TempVarNames))
		}
	})
}
