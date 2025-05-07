package core_test

import (
	"testing"

	"smalltalklsp/interpreter/core"
)

// TestAllocateWithGCNeeded tests the Allocate method when garbage collection is needed
func TestAllocateWithGCNeeded(t *testing.T) {
	om := core.NewObjectMemory()

	// Set the allocation pointer to the GC threshold
	om.AllocPtr = om.GCThreshold

	// Allocate an object
	obj := core.MakeIntegerImmediate(42)
	result := om.Allocate(obj)

	// Check that the object is returned as-is
	if result != obj {
		t.Errorf("Expected Allocate to return the original object when GC is needed")
	}

	// Check that the object was not allocated in the from-space
	// (since we're returning early to let the VM handle collection)
	if om.AllocPtr > om.GCThreshold {
		t.Errorf("Expected AllocPtr to not be incremented when GC is needed")
	}
}
