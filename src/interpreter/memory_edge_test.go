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
	obj := vm.NewIntegerWithClass(42, vm.IntegerClass)
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

	// Create special objects
	nilObj := NewNil()
	trueObj := NewBoolean(true)
	falseObj := NewBoolean(false)
	objectClass := NewClass("Object", nil)

	// Set the VM's special objects
	vm.NilObject = nilObj
	vm.TrueObject = trueObj
	vm.FalseObject = falseObj
	vm.ObjectClass = objectClass

	// Allocate the objects
	om.Allocate(nilObj)
	om.Allocate(trueObj)
	om.Allocate(falseObj)
	om.Allocate(objectClass)

	// Perform garbage collection
	om.Collect(vm)

	// Check that the special objects are still in the VM
	if vm.NilObject != nilObj {
		t.Errorf("Expected vm.NilObject to still be nilObj")
	}

	if vm.TrueObject != trueObj {
		t.Errorf("Expected vm.TrueObject to still be trueObj")
	}

	if vm.FalseObject != falseObj {
		t.Errorf("Expected vm.FalseObject to still be falseObj")
	}

	if vm.ObjectClass != objectClass {
		t.Errorf("Expected vm.ObjectClass to still be objectClass")
	}

	// Verify that all special objects are in the new from-space
	foundNil := false
	foundTrue := false
	foundFalse := false
	foundObjectClass := false

	for i := 0; i < om.AllocPtr; i++ {
		obj := om.FromSpace[i]
		if obj == nilObj {
			foundNil = true
		} else if obj == trueObj {
			foundTrue = true
		} else if obj == falseObj {
			foundFalse = true
		} else if obj == objectClass {
			foundObjectClass = true
		}
	}

	if !foundNil {
		t.Errorf("NilObject not found in from-space after collection")
	}

	if !foundTrue {
		t.Errorf("TrueObject not found in from-space after collection")
	}

	if !foundFalse {
		t.Errorf("FalseObject not found in from-space after collection")
	}

	if !foundObjectClass {
		t.Errorf("ObjectClass not found in from-space after collection")
	}
}

// TestUpdateReferencesEdgeCases tests edge cases in the updateReferences method
func TestUpdateReferencesEdgeCases(t *testing.T) {
	om := NewObjectMemory()

	// Test with nil references
	{
		// Create an instance with nil instance variables
		class := NewClass("TestClass", nil)
		instance := NewInstance(class)
		instance.InstanceVars = make([]*Object, 2)
		instance.InstanceVars[METHOD_DICTIONARY_IV] = nil
		instance.InstanceVars[1] = nil

		// Set nil superclass and class
		instance.SuperClass = nil
		instance.Class = nil

		// Update references
		toPtr := 0
		om.updateReferences(instance, &toPtr)

		// Check that toPtr hasn't changed
		if toPtr != 0 {
			t.Errorf("Expected toPtr to be 0, got %d", toPtr)
		}
	}

	// Test with nil array elements
	{
		// Create an array with nil elements
		array := NewArray(2)
		array.Elements[0] = nil
		array.Elements[1] = nil

		// Update references
		toPtr := 0
		om.updateReferences(array, &toPtr)

		// Check that toPtr hasn't changed
		if toPtr != 0 {
			t.Errorf("Expected toPtr to be 0, got %d", toPtr)
		}
	}

	// Test with nil dictionary entries
	{
		// Create a dictionary with nil entries
		dict := NewDictionary()
		dict.Entries["key1"] = nil
		dict.Entries["key2"] = nil

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
		method := NewMethod(nil, nil)
		method.Method.Literals = make([]*Object, 2)
		method.Method.Literals[0] = nil
		method.Method.Literals[1] = nil

		// Set nil selector and class
		method.Method.Selector = nil
		method.Method.Class = nil

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

// TestCollectWithNilInRootSet tests the Collect method with nil objects in the root set
func TestCollectWithNilInRootSet(t *testing.T) {
	om := NewObjectMemory()
	vm := NewVM()

	// Add nil objects to globals
	vm.Globals["nil1"] = nil
	vm.Globals["nil2"] = nil

	// Create a context with nil references
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)
	context := NewContext(methodObj, nil, nil, nil)

	// Set nil stack elements
	context.Stack[0] = nil
	context.Stack[1] = nil
	context.StackPointer = 2

	// Set the context as the VM's current context
	vm.CurrentContext = context

	// Perform garbage collection
	om.Collect(vm)

	// Check that the collection completed without errors
	if om.GCCount != 1 {
		t.Errorf("Expected GCCount to be 1 after collection, got %d", om.GCCount)
	}

	// Check that the context is still the VM's current context
	if vm.CurrentContext != context {
		t.Errorf("Expected vm.CurrentContext to still be context")
	}
}

// TestCollectWithNilInToSpace tests the Collect method with nil objects in the to-space
func TestCollectWithNilInToSpace(t *testing.T) {
	om := NewObjectMemory()
	vm := NewVM()

	// Create and allocate an object
	obj := vm.NewIntegerWithClass(42, vm.IntegerClass)
	om.Allocate(obj)

	// Add the object to globals to make it reachable
	vm.Globals["obj"] = obj

	// Create a nil slot in the to-space
	om.ToSpace[0] = nil

	// Perform garbage collection
	om.Collect(vm)

	// Check that the collection completed without errors
	if om.GCCount != 1 {
		t.Errorf("Expected GCCount to be 1 after collection, got %d", om.GCCount)
	}

	// Check that the object is still in the VM's globals
	if vm.Globals["obj"] != obj {
		t.Errorf("Expected vm.Globals[\"obj\"] to still be obj")
	}
}
