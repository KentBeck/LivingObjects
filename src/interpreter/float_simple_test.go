package main

import (
	"math"
	"testing"
)

// TestFloatImmediateSimple tests the immediate float implementation
func TestFloatImmediateSimple(t *testing.T) {
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
