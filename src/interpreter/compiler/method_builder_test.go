package compiler_test

import (
	"testing"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/pile"
)

func TestMethodBuilderBasic(t *testing.T) {
	// Create a class for testing
	testClass := pile.NewClass("TestClass", nil)

	// Create a simple method
	method := compiler.NewMethodBuilder(testClass).
		Go("testMethod")

	// Check that the method is not nil
	if method == nil {
		t.Fatal("Method should not be nil")
	}

	// Check that the method has the correct selector
	methodObj := pile.ObjectToMethod(method)
	if methodObj.Selector == nil {
		t.Fatal("Method selector should not be nil")
	}

	symbol := pile.ObjectToSymbol(methodObj.Selector)
	if symbol.GetValue() != "testMethod" {
		t.Errorf("Expected method selector to be 'testMethod', got '%s'", symbol.GetValue())
	}
}

func TestMethodBuilderWithLiterals(t *testing.T) {
	// Create a class for testing
	testClass := pile.NewClass("TestClass", nil)

	// Create a method with literals
	literal := pile.MakeIntegerImmediate(42)
	builder := compiler.NewMethodBuilder(testClass)
	literalIndex, builder := builder.AddLiteral(literal)
	method := builder.Go("testWithLiterals")

	// Check that the method is not nil
	if method == nil {
		t.Fatal("Method should not be nil")
	}

	// Check that the method has the correct literal
	methodObj := pile.ObjectToMethod(method)
	if len(methodObj.Literals) != 1 {
		t.Errorf("Expected 1 literal, got %d", len(methodObj.Literals))
	}

	if methodObj.Literals[literalIndex] != literal {
		t.Errorf("Expected literal at index %d to be %v, got %v", literalIndex, literal, methodObj.Literals[literalIndex])
	}
}

func TestMethodBuilderWithTempVars(t *testing.T) {
	// Create a class for testing
	testClass := pile.NewClass("TestClass", nil)

	// Create a method with temporary variables
	method := compiler.NewMethodBuilder(testClass).
		TempVars([]string{"temp1", "temp2"}).
		Go("testWithTempVars")

	// Check that the method is not nil
	if method == nil {
		t.Fatal("Method should not be nil")
	}

	// Check that the method has the correct temporary variable names
	methodObj := pile.ObjectToMethod(method)
	if len(methodObj.TempVarNames) != 2 {
		t.Errorf("Expected 2 temporary variable names, got %d", len(methodObj.TempVarNames))
	}

	if methodObj.TempVarNames[0] != "temp1" {
		t.Errorf("Expected temporary variable name 'temp1', got '%s'", methodObj.TempVarNames[0])
	}

	if methodObj.TempVarNames[1] != "temp2" {
		t.Errorf("Expected temporary variable name 'temp2', got '%s'", methodObj.TempVarNames[1])
	}
}