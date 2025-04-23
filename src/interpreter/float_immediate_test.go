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
	vm := NewVM()

	// Test addition
	{
		float1 := vm.NewFloat(3.14)
		float2 := vm.NewFloat(2.71)
		plusSelector := NewSymbol("+")
		plusMethod := vm.FloatClass.GetMethodDict().Entries[plusSelector.SymbolValue]
		result := vm.executePrimitive(float1, plusSelector, []*Object{float2}, plusMethod)
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
		float1 := vm.NewFloat(3.14)
		float2 := vm.NewFloat(2.71)
		minusSelector := NewSymbol("-")
		minusMethod := vm.FloatClass.GetMethodDict().Entries[minusSelector.SymbolValue]
		result := vm.executePrimitive(float1, minusSelector, []*Object{float2}, minusMethod)
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
		float1 := vm.NewFloat(3.14)
		float2 := vm.NewFloat(2.71)
		timesSelector := NewSymbol("*")
		timesMethod := vm.FloatClass.GetMethodDict().Entries[timesSelector.SymbolValue]
		result := vm.executePrimitive(float1, timesSelector, []*Object{float2}, timesMethod)
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
		float1 := vm.NewFloat(3.14)
		float2 := vm.NewFloat(2.71)
		divideSelector := NewSymbol("/")
		divideMethod := vm.FloatClass.GetMethodDict().Entries[divideSelector.SymbolValue]
		result := vm.executePrimitive(float1, divideSelector, []*Object{float2}, divideMethod)
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
		float1 := vm.NewFloat(3.14)
		float2 := vm.NewFloat(3.14)
		equalsSelector := NewSymbol("=")
		equalsMethod := vm.FloatClass.GetMethodDict().Entries[equalsSelector.SymbolValue]
		result := vm.executePrimitive(float1, equalsSelector, []*Object{float2}, equalsMethod)
		if !IsTrueImmediate(result) {
			t.Errorf("Expected %f = %f to be true", 3.14, 3.14)
		}

		float3 := vm.NewFloat(2.71)
		result = vm.executePrimitive(float1, equalsSelector, []*Object{float3}, equalsMethod)
		if !IsFalseImmediate(result) {
			t.Errorf("Expected %f = %f to be false", 3.14, 2.71)
		}
	}

	// Test less than
	{
		float1 := vm.NewFloat(2.71)
		float2 := vm.NewFloat(3.14)
		lessSelector := NewSymbol("<")
		lessMethod := vm.FloatClass.GetMethodDict().Entries[lessSelector.SymbolValue]
		result := vm.executePrimitive(float1, lessSelector, []*Object{float2}, lessMethod)
		if !IsTrueImmediate(result) {
			t.Errorf("Expected %f < %f to be true", 2.71, 3.14)
		}

		result = vm.executePrimitive(float2, lessSelector, []*Object{float1}, lessMethod)
		if !IsFalseImmediate(result) {
			t.Errorf("Expected %f < %f to be false", 3.14, 2.71)
		}
	}

	// Test greater than
	{
		float1 := vm.NewFloat(3.14)
		float2 := vm.NewFloat(2.71)
		greaterSelector := NewSymbol(">")
		greaterMethod := vm.FloatClass.GetMethodDict().Entries[greaterSelector.SymbolValue]
		result := vm.executePrimitive(float1, greaterSelector, []*Object{float2}, greaterMethod)
		if !IsTrueImmediate(result) {
			t.Errorf("Expected %f > %f to be true", 3.14, 2.71)
		}

		result = vm.executePrimitive(float2, greaterSelector, []*Object{float1}, greaterMethod)
		if !IsFalseImmediate(result) {
			t.Errorf("Expected %f > %f to be false", 2.71, 3.14)
		}
	}

	// Test mixed operations (float and integer)
	{
		float1 := vm.NewFloat(3.14)
		int1 := vm.NewInteger(2)
		plusSelector := NewSymbol("+")
		floatPlusMethod := vm.FloatClass.GetMethodDict().Entries[plusSelector.SymbolValue]
		result := vm.executePrimitive(float1, plusSelector, []*Object{int1}, floatPlusMethod)
		if !IsFloatImmediate(result) {
			t.Errorf("Expected result to be a float immediate")
		}
		value := GetFloatImmediate(result)
		expected := 3.14 + 2.0
		if math.Abs(value-expected) > 1e-10 {
			t.Errorf("Expected %f + %d = %f, got %f", 3.14, 2, expected, value)
		}

		intPlusMethod := vm.IntegerClass.GetMethodDict().Entries[plusSelector.SymbolValue]
		result = vm.executePrimitive(int1, plusSelector, []*Object{float1}, intPlusMethod)
		if !IsFloatImmediate(result) {
			t.Errorf("Expected result to be a float immediate")
		}
		value = GetFloatImmediate(result)
		expected = 2.0 + 3.14
		if math.Abs(value-expected) > 1e-10 {
			t.Errorf("Expected %d + %f = %f, got %f", 2, 3.14, expected, value)
		}
	}
}
