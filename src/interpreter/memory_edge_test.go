package main

import (
	"testing"
)

// TestAllocateWithGCNeeded tests the Allocate method when garbage collection is needed
func TestAllocateWithGCNeeded(t *testing.T) {
	om := NewObjectMemory()

	// Set the allocation pointer to the GC threshold
	om.AllocPtr = om.GCThreshold

	// Create a VM for the integer class
	vm := NewVM()

	// Allocate an object
	obj := vm.NewInteger(42)
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

// TestCollectWithSpecialObjects tests the Collect method with special VM objects
func TestCollectWithSpecialObjects(t *testing.T) {
	om := NewObjectMemory()
	vm := NewVM()

	// Store the immediate values for nil, true, and false
	nilImmediate := vm.NilObject
	trueImmediate := vm.TrueObject
	falseImmediate := vm.FalseObject

	// Create a non-immediate object
	objectClass := NewClass("Object", nil)

	// Set the VM's object class
	vm.ObjectClass = objectClass

	// Allocate the object class
	om.Allocate(ClassToObject(objectClass)) // Convert to Object for allocation

	// Perform garbage collection
	om.Collect(vm)

	// Check that the special immediate objects are still the same
	if vm.NilObject != nilImmediate {
		t.Errorf("Expected vm.NilObject to still be nilImmediate")
	}

	if vm.TrueObject != trueImmediate {
		t.Errorf("Expected vm.TrueObject to still be trueImmediate")
	}

	if vm.FalseObject != falseImmediate {
		t.Errorf("Expected vm.FalseObject to still be falseImmediate")
	}

	if vm.ObjectClass != objectClass {
		t.Errorf("Expected vm.ObjectClass to still be objectClass")
	}

	// Verify that the object class is in the new from-space
	foundObjectClass := false

	for i := 0; i < om.AllocPtr; i++ {
		obj := om.FromSpace[i]
		if obj == ClassToObject(objectClass) {
			foundObjectClass = true
		}
	}

	// Check that special objects are immediate values
	if !IsNilImmediate(vm.NilObject) {
		t.Errorf("Expected NilObject to be an immediate value")
	}

	if !IsTrueImmediate(vm.TrueObject) {
		t.Errorf("Expected TrueObject to be an immediate value")
	}

	if !IsFalseImmediate(vm.FalseObject) {
		t.Errorf("Expected FalseObject to be an immediate value")
	}

	if !foundObjectClass {
		t.Errorf("ObjectClass not found in from-space after collection")
	}
}

// TestUpdateReferencesEdgeCases tests edge cases in the updateReferences method
func TestUpdateReferencesEdgeCases(t *testing.T) {
	om := NewObjectMemory()

	// Test with nil array elements
	{
		// Create an array with nil elements
		array := NewArray(2).(*Array)
		array.Elements[0] = NewNil().(*Object)
		array.Elements[1] = NewNil().(*Object)

		// Update references
		toPtr := 0
		om.updateReferences(&array.Object, &toPtr)

		// Check that toPtr hasn't changed
		if toPtr != 0 {
			t.Errorf("Expected toPtr to be 0, got %d", toPtr)
		}
	}

	// Test with nil dictionary entries
	{
		// Create a dictionary with nil entries
		dict := NewDictionary()
		dict.Entries["key1"] = NewNil().(*Object)
		dict.Entries["key2"] = NewNil().(*Object)

		// Update references
		toPtr := 0
		om.updateReferences(dict, &toPtr)

		// Check that toPtr hasn't changed
		if toPtr != 0 {
			t.Errorf("Expected toPtr to be 0, got %d", toPtr)
		}
	}

	// Test with nil method literals, selector, and class
	{
		// Create a method with nil literals
		method := &Object{
			type1:  OBJ_METHOD,
			Method: &Method{},
		}
		method.Method.Literals = make([]*Object, 2)
		method.Method.Literals[0] = nil
		method.Method.Literals[1] = nil

		// Set nil selector and class
		method.Method.Selector = nil
		method.Method.MethodClass = nil

		// Update references
		toPtr := 0
		om.updateReferences(method, &toPtr)

		// Check that toPtr hasn't changed
		if toPtr != 0 {
			t.Errorf("Expected toPtr to be 0, got %d", toPtr)
		}
	}
}

// TestCollectWithNilSpecialObjects tests the Collect method with nil special VM objects
func TestCollectWithNilSpecialObjects(t *testing.T) {
	om := NewObjectMemory()
	vm := NewVM()

	// Set the VM's special objects to nil
	vm.NilObject = nil
	vm.TrueObject = nil
	vm.FalseObject = nil
	vm.ObjectClass = nil

	// Perform garbage collection
	om.Collect(vm)

	// Check that the collection completed without errors
	if om.GCCount != 1 {
		t.Errorf("Expected GCCount to be 1 after collection, got %d", om.GCCount)
	}
}

// TestCollectWithNilInToSpace tests the Collect method with nil objects in the to-space
func TestCollectWithNilInToSpace(t *testing.T) {
	om := NewObjectMemory()
	vm := NewVM()

	// Create and allocate a non-immediate object
	obj := NewString("test")
	objAsObj := StringToObject(obj) // Convert to Object for allocation
	om.Allocate(objAsObj)

	// Add the object to globals to make it reachable
	vm.Globals["obj"] = objAsObj

	// Create a nil slot in the to-space
	om.ToSpace[0] = nil

	// Perform garbage collection
	om.Collect(vm)

	// Check that the collection completed without errors
	if om.GCCount != 1 {
		t.Errorf("Expected GCCount to be 1 after collection, got %d", om.GCCount)
	}

	// Check that the object is still in the VM's globals
	if vm.Globals["obj"] != objAsObj {
		t.Errorf("Expected vm.Globals[\"obj\"] to still be obj")
	}
}
