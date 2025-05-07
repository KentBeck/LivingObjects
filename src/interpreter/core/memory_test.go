package core_test

import (
	"testing"

	"smalltalklsp/interpreter/core"
)

// TestNewObjectMemory tests the creation of a new object memory
func TestNewObjectMemory(t *testing.T) {
	om := core.NewObjectMemory()

	// Check initial values
	if om.AllocPtr != 0 {
		t.Errorf("Expected AllocPtr to be 0, got %d", om.AllocPtr)
	}

	if om.SpaceSize != 10000 {
		t.Errorf("Expected SpaceSize to be 10000, got %d", om.SpaceSize)
	}

	if om.GCThreshold != 8000 {
		t.Errorf("Expected GCThreshold to be 8000, got %d", om.GCThreshold)
	}

	if om.GCCount != 0 {
		t.Errorf("Expected GCCount to be 0, got %d", om.GCCount)
	}

	if len(om.FromSpace) != 10000 {
		t.Errorf("Expected FromSpace length to be 10000, got %d", len(om.FromSpace))
	}

	if len(om.ToSpace) != 10000 {
		t.Errorf("Expected ToSpace length to be 10000, got %d", len(om.ToSpace))
	}
}

// TestShouldCollect tests the ShouldCollect method
func TestShouldCollect(t *testing.T) {
	om := core.NewObjectMemory()

	// Initially, we shouldn't need to collect
	if om.ShouldCollect() {
		t.Errorf("Expected ShouldCollect to return false initially")
	}

	// Set the allocation pointer to just below the threshold
	om.AllocPtr = om.GCThreshold - 1

	// We still shouldn't need to collect
	if om.ShouldCollect() {
		t.Errorf("Expected ShouldCollect to return false when AllocPtr < GCThreshold")
	}

	// Set the allocation pointer to the threshold
	om.AllocPtr = om.GCThreshold

	// Now we should need to collect
	if !om.ShouldCollect() {
		t.Errorf("Expected ShouldCollect to return true when AllocPtr >= GCThreshold")
	}

	// Set the allocation pointer above the threshold
	om.AllocPtr = om.GCThreshold + 1

	// We should still need to collect
	if !om.ShouldCollect() {
		t.Errorf("Expected ShouldCollect to return true when AllocPtr > GCThreshold")
	}
}
