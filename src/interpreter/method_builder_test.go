package main

import (
	"testing"
)

func TestMethodBuilder(t *testing.T) {
	// Create a class for testing
	testClass := NewClass("TestClass", nil)

	t.Run("PrimitiveMethod", func(t *testing.T) {
		// Create a primitive method using MethodBuilder
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
		selectorValue := "testPrimitive"
		methodInDict := methodDict.Entries[selectorValue]
		
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
		bytecodes := []byte{PUSH_SELF, RETURN_STACK_TOP}
		method := NewMethodBuilder(testClass).
			Selector("testMethod").
			Bytecodes(bytecodes).
			Go()

		// Verify the method was created correctly
		if method == nil {
			t.Fatal("Method should not be nil")
		}

		// Check that the method is in the method dictionary
		methodDict := testClass.GetMethodDict()
		selectorValue := "testMethod"
		methodInDict := methodDict.Entries[selectorValue]
		
		if methodInDict == nil {
			t.Fatalf("Method not found in dictionary for selector %q", selectorValue)
		}

		// Check the method properties
		if len(methodInDict.Method.Bytecodes) != len(bytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(bytecodes), len(methodInDict.Method.Bytecodes))
		}

		for i, b := range bytecodes {
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
			Literals(literals).
			Go()

		// Verify the method was created correctly
		if method == nil {
			t.Fatal("Method should not be nil")
		}

		// Check that the method is in the method dictionary
		methodDict := testClass.GetMethodDict()
		selectorValue := "testLiterals"
		methodInDict := methodDict.Entries[selectorValue]
		
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
		bytecodes := []byte{PUSH_SELF, RETURN_STACK_TOP}
		literals := []*Object{MakeIntegerImmediate(123), NewSymbol("test")}
		tempVars := []string{"temp1", "temp2"}
		
		method := NewMethodBuilder(testClass).
			Selector("completeMethod").
			Primitive(42).
			Bytecodes(bytecodes).
			Literals(literals).
			TempVars(tempVars).
			Go()

		// Verify the method was created correctly
		if method == nil {
			t.Fatal("Method should not be nil")
		}

		// Check that the method is in the method dictionary
		methodDict := testClass.GetMethodDict()
		selectorValue := "completeMethod"
		methodInDict := methodDict.Entries[selectorValue]
		
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

		if len(methodInDict.Method.Bytecodes) != len(bytecodes) {
			t.Errorf("Expected %d bytecodes, got %d", len(bytecodes), len(methodInDict.Method.Bytecodes))
		}

		if len(methodInDict.Method.Literals) != len(literals) {
			t.Errorf("Expected %d literals, got %d", len(literals), len(methodInDict.Method.Literals))
		}

		if len(methodInDict.Method.TempVarNames) != len(tempVars) {
			t.Errorf("Expected %d temp vars, got %d", len(tempVars), len(methodInDict.Method.TempVarNames))
		}
	})
}
