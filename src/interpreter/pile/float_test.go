package pile_test

import (
	"math"
	"testing"

	"smalltalklsp/interpreter/pile"
)

// TestFloatImmediate tests the immediate float implementation
func TestFloatImmediate(t *testing.T) {
	// Test MakeFloatImmediate and GetFloatImmediate
	value := 3.14159
	obj := pile.MakeFloatImmediate(value)

	// Check that it's an immediate value
	if !pile.IsImmediate(obj) {
		t.Errorf("Expected obj to be an immediate value")
	}

	// Check that it's a float immediate
	if !pile.IsFloatImmediate(obj) {
		t.Errorf("Expected obj to be a float immediate")
	}

	// Check that the tag is correct
	if pile.GetTag(obj) != pile.TAG_FLOAT {
		t.Errorf("Expected tag to be TAG_FLOAT, got %d", pile.GetTag(obj))
	}

	// Check that we can get the value back
	retrievedValue := pile.GetFloatImmediate(obj)
	if math.Abs(retrievedValue-value) > 1e-10 {
		t.Errorf("Expected to get back %f, got %f", value, retrievedValue)
	}
}