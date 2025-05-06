package main

import (
	"math"
	"testing"
)

// TestFloatImmediate tests the immediate float implementation
func TestFloatImmediate(t *testing.T) {
	// Test MakeFloatImmediate and GetFloatImmediate
	value := 3.14159
	obj := MakeFloatImmediate(value)

	// Check that it's an immediate value
	if !IsImmediate(obj) {
		t.Errorf("Expected obj to be an immediate value")
	}

	// Check that it's a float immediate
	if !IsFloatImmediate(obj) {
		t.Errorf("Expected obj to be a float immediate")
	}

	// Check that the tag is correct
	if GetTag(obj) != TAG_FLOAT {
		t.Errorf("Expected tag to be TAG_FLOAT, got %d", GetTag(obj))
	}

	// Check that we can get the value back
	retrievedValue := GetFloatImmediate(obj)
	if math.Abs(retrievedValue-value) > 1e-10 {
		t.Errorf("Expected to get back %f, got %f", value, retrievedValue)
	}
}

// TestFloatPrimitives tests the float primitive operations
func TestFloatPrimitives(t *testing.T) {
	// Create a new VM
	vm := NewVM()

	// Make sure the FloatClass is properly initialized
	if vm.FloatClass == nil {
		t.Fatalf("FloatClass is nil")
	}

	// Check if the method dictionary is properly initialized
	floatMethodDict := vm.FloatClass.GetMethodDict()
	if floatMethodDict == nil {
		t.Fatalf("FloatClass method dictionary is nil")
	}

	// Check if the method dictionary has entries
	// Convert to Dictionary to access entries
	dict := ObjectToDictionary(floatMethodDict)
	if dict.Entries == nil {
		t.Fatalf("FloatClass method dictionary entries is nil")
	}

	// Test addition
	{
		// Create a simple test for float addition
		float1 := vm.NewFloat(3.14)
		float2 := vm.NewFloat(2.71)

		// Create a new VM for this test
		testVM := NewVM()

		// Get the + method
		plusSelector := NewSymbol("+")
		methodDict := testVM.FloatClass.GetMethodDict()
		dict := ObjectToDictionary(methodDict)
		plusMethod := dict.Entries[GetSymbolValue(plusSelector)]

		// Execute the primitive
		result := testVM.executePrimitive(float1, plusSelector, []*Object{float2}, plusMethod)

		// Check the result
		if !IsFloatImmediate(result) {
			t.Errorf("Expected result to be a float immediate")
		}
		value := GetFloatImmediate(result)
		expected := 3.14 + 2.71
		if math.Abs(value-expected) > 1e-10 {
			t.Errorf("Expected %f + %f = %f, got %f", 3.14, 2.71, expected, value)
		}
	}

	// Test subtraction
	{
		// Create a simple test for float subtraction
		float1 := vm.NewFloat(3.14)
		float2 := vm.NewFloat(2.71)

		// Create a new VM for this test
		testVM := NewVM()

		// Get the - method
		minusSelector := NewSymbol("-")
		methodDict := testVM.FloatClass.GetMethodDict()
		dict := ObjectToDictionary(methodDict)
		minusMethod := dict.Entries[GetSymbolValue(minusSelector)]

		// Execute the primitive
		result := testVM.executePrimitive(float1, minusSelector, []*Object{float2}, minusMethod)

		// Check the result
		if !IsFloatImmediate(result) {
			t.Errorf("Expected result to be a float immediate")
		}
		value := GetFloatImmediate(result)
		expected := 3.14 - 2.71
		if math.Abs(value-expected) > 1e-10 {
			t.Errorf("Expected %f - %f = %f, got %f", 3.14, 2.71, expected, value)
		}
	}

	// Test multiplication
	{
		// Create a simple test for float multiplication
		float1 := vm.NewFloat(3.14)
		float2 := vm.NewFloat(2.71)

		// Create a new VM for this test
		testVM := NewVM()

		// Get the * method
		timesSelector := NewSymbol("*")
		methodDict := testVM.FloatClass.GetMethodDict()
		dict := ObjectToDictionary(methodDict)
		timesMethod := dict.Entries[GetSymbolValue(timesSelector)]

		// Execute the primitive
		result := testVM.executePrimitive(float1, timesSelector, []*Object{float2}, timesMethod)

		// Check the result
		if !IsFloatImmediate(result) {
			t.Errorf("Expected result to be a float immediate")
		}
		value := GetFloatImmediate(result)
		expected := 3.14 * 2.71
		if math.Abs(value-expected) > 1e-10 {
			t.Errorf("Expected %f * %f = %f, got %f", 3.14, 2.71, expected, value)
		}
	}

	// Test division
	{
		// Create a simple test for float division
		float1 := vm.NewFloat(3.14)
		float2 := vm.NewFloat(2.71)

		// Create a new VM for this test
		testVM := NewVM()

		// Get the / method
		divideSelector := NewSymbol("/")
		methodDict := testVM.FloatClass.GetMethodDict()
		dict := ObjectToDictionary(methodDict)
		divideMethod := dict.Entries[GetSymbolValue(divideSelector)]

		// Execute the primitive
		result := testVM.executePrimitive(float1, divideSelector, []*Object{float2}, divideMethod)

		// Check the result
		if !IsFloatImmediate(result) {
			t.Errorf("Expected result to be a float immediate")
		}
		value := GetFloatImmediate(result)
		expected := 3.14 / 2.71
		if math.Abs(value-expected) > 1e-10 {
			t.Errorf("Expected %f / %f = %f, got %f", 3.14, 2.71, expected, value)
		}
	}

	// Test equality
	{
		// Create a simple test for float equality
		float1 := vm.NewFloat(3.14)
		float2 := vm.NewFloat(3.14)
		float3 := vm.NewFloat(2.71)

		// Create a new VM for this test
		testVM := NewVM()

		// Get the = method
		equalsSelector := NewSymbol("=")
		methodDict := testVM.FloatClass.GetMethodDict()
		dict := ObjectToDictionary(methodDict)
		equalsMethod := dict.Entries[GetSymbolValue(equalsSelector)]

		// Test equality with equal values
		result := testVM.executePrimitive(float1, equalsSelector, []*Object{float2}, equalsMethod)
		if !IsTrueImmediate(result) {
			t.Errorf("Expected %f = %f to be true", 3.14, 3.14)
		}

		// Test equality with different values
		result = testVM.executePrimitive(float1, equalsSelector, []*Object{float3}, equalsMethod)
		if !IsFalseImmediate(result) {
			t.Errorf("Expected %f = %f to be false", 3.14, 2.71)
		}
	}

	// Test less than
	{
		// Create a simple test for float less than
		float1 := vm.NewFloat(2.71)
		float2 := vm.NewFloat(3.14)

		// Create a new VM for this test
		testVM := NewVM()

		// Get the < method
		lessSelector := NewSymbol("<")
		methodDict := testVM.FloatClass.GetMethodDict()
		dict := ObjectToDictionary(methodDict)
		lessMethod := dict.Entries[GetSymbolValue(lessSelector)]

		// Test less than with smaller value first
		result := testVM.executePrimitive(float1, lessSelector, []*Object{float2}, lessMethod)
		if !IsTrueImmediate(result) {
			t.Errorf("Expected %f < %f to be true", 2.71, 3.14)
		}

		// Test less than with larger value first
		result = testVM.executePrimitive(float2, lessSelector, []*Object{float1}, lessMethod)
		if !IsFalseImmediate(result) {
			t.Errorf("Expected %f < %f to be false", 3.14, 2.71)
		}
	}

	// Test greater than
	{
		// Create a simple test for float greater than
		float1 := vm.NewFloat(3.14)
		float2 := vm.NewFloat(2.71)

		// Create a new VM for this test
		testVM := NewVM()

		// Get the > method
		greaterSelector := NewSymbol(">")
		methodDict := testVM.FloatClass.GetMethodDict()
		dict := ObjectToDictionary(methodDict)
		greaterMethod := dict.Entries[GetSymbolValue(greaterSelector)]

		// Test greater than with larger value first
		result := testVM.executePrimitive(float1, greaterSelector, []*Object{float2}, greaterMethod)
		if !IsTrueImmediate(result) {
			t.Errorf("Expected %f > %f to be true", 3.14, 2.71)
		}

		// Test greater than with smaller value first
		result = testVM.executePrimitive(float2, greaterSelector, []*Object{float1}, greaterMethod)
		if !IsFalseImmediate(result) {
			t.Errorf("Expected %f > %f to be false", 2.71, 3.14)
		}
	}

	// Test mixed operations (float and integer)
	{
		// Create a simple test for mixed float and integer operations
		float1 := vm.NewFloat(3.14)
		int1 := vm.NewInteger(2)

		// Create a new VM for this test
		testVM := NewVM()

		// Test float + integer
		plusSelector := NewSymbol("+")
		methodDict := testVM.FloatClass.GetMethodDict()
		dict := ObjectToDictionary(methodDict)
		floatPlusMethod := dict.Entries[GetSymbolValue(plusSelector)]
		result := testVM.executePrimitive(float1, plusSelector, []*Object{int1}, floatPlusMethod)
		if !IsFloatImmediate(result) {
			t.Errorf("Expected result to be a float immediate")
		}
		value := GetFloatImmediate(result)
		expected := 3.14 + 2.0
		if math.Abs(value-expected) > 1e-10 {
			t.Errorf("Expected %f + %d = %f, got %f", 3.14, 2, expected, value)
		}

		// Test integer + float
		// Skip this part for now as it requires updating the integer primitives
		// to handle float arguments
	}
}
