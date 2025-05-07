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
		methodInDict := dict.GetEntry(selectorValue)

		if methodInDict == nil {
			t.Fatalf("Method not found in dictionary for selector %q", selectorValue)
		}

		// Check the method properties
		methodObj := ObjectToMethod(methodInDict)
		if !methodObj.IsPrimitive {
			t.Error("Method should be marked as primitive")
		}

		if methodObj.PrimitiveIndex != 42 {
			t.Errorf("Expected primitive index 42, got %d", methodObj.PrimitiveIndex)
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
		methodInDict := dict.GetEntry(selectorValue)

		if methodInDict == nil {
			t.Fatalf("Method not found in dictionary for selector %q", selectorValue)
		}

		// Check the method properties
		methodObj := ObjectToMethod(methodInDict)
		expectedBytecodes := []byte{PUSH_SELF, RETURN_STACK_TOP}
		if len(methodObj.Bytecodes) != len(expectedBytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(expectedBytecodes), len(methodObj.Bytecodes))
		}

		for i, b := range expectedBytecodes {
			if methodObj.Bytecodes[i] != b {
				t.Errorf("Expected bytecode %d at index %d, got %d", b, i, methodObj.Bytecodes[i])
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
		methodInDict := dict.GetEntry(selectorValue)

		if methodInDict == nil {
			t.Fatalf("Method not found in dictionary for selector %q", selectorValue)
		}

		// Check the method properties
		methodObj := ObjectToMethod(methodInDict)
		if len(methodObj.Literals) != len(literals) {
			t.Errorf("Expected %d literals, got %d", len(literals), len(methodObj.Literals))
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
		methodInDict := dict.GetEntry(selectorValue)

		if methodInDict == nil {
			t.Fatalf("Method not found in dictionary for selector %q", selectorValue)
		}

		// Check the method properties
		methodObj := ObjectToMethod(methodInDict)
		if !methodObj.IsPrimitive {
			t.Error("Method should be marked as primitive")
		}

		if methodObj.PrimitiveIndex != 42 {
			t.Errorf("Expected primitive index 42, got %d", methodObj.PrimitiveIndex)
		}

		expectedBytecodes := []byte{PUSH_SELF, RETURN_STACK_TOP}
		if len(methodObj.Bytecodes) != len(expectedBytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(expectedBytecodes), len(methodObj.Bytecodes))
		}

		if len(methodObj.Literals) != len(literals) {
			t.Errorf("Expected %d literals, got %d", len(literals), len(methodObj.Literals))
		}

		if len(methodObj.TempVarNames) != len(tempVars) {
			t.Errorf("Expected %d temp vars, got %d", len(tempVars), len(methodObj.TempVarNames))
		}
	})
}
