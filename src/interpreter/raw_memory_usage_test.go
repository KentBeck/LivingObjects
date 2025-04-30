package main

import (
	"testing"
	"unsafe"
)

// TestRawMemorySimpleObjectUsage tests allocating and using a simple object in raw memory
func TestRawMemorySimpleObjectUsage(t *testing.T) {
	// Create a new raw memory
	rawMem, err := NewRawMemory(1024*1024, 1024*1024)
	if err != nil {
		t.Fatalf("Failed to create raw memory: %v", err)
	}
	defer rawMem.Close()

	// 1. Allocate a simple object
	obj := rawMem.AllocateObject()
	if obj == nil {
		t.Fatalf("Failed to allocate object")
	}

	// 2. Initialize the object with some values
	obj.SetType(OBJ_INSTANCE)

	// 3. Allocate instance variables array
	instanceVars := rawMem.AllocateObjectArray(3)
	if instanceVars == nil {
		t.Fatalf("Failed to allocate instance variables array")
	}
	obj.SetInstanceVars(instanceVars)

	// 4. Create some values to store in the instance variables
	// For immediate values, we don't need to allocate memory
	intValue := MakeIntegerImmediate(42)
	boolValue := MakeTrueImmediate()

	// For a string, we need to allocate memory
	strObj := rawMem.AllocateString()
	if strObj == nil {
		t.Fatalf("Failed to allocate string")
	}
	strObj.SetType(OBJ_STRING)
	strObj.Value = "Hello, Raw Memory!"
	stringValue := StringToObject(strObj)

	// 5. Store values in instance variables
	instanceVars[0] = intValue
	instanceVars[1] = boolValue
	instanceVars[2] = stringValue

	// 6. Verify the object is in from-space
	if !rawMem.IsPointerInFromSpace(unsafe.Pointer(obj)) {
		t.Errorf("Expected object to be in from-space")
	}

	// 7. Verify instance variables
	objVars := obj.InstanceVars()
	if !IsIntegerImmediate(objVars[0]) {
		t.Errorf("Expected instance variable 0 to be an integer immediate")
	}
	intVal := GetIntegerImmediate(objVars[0])
	if intVal != 42 {
		t.Errorf("Expected instance variable 0 to be 42, got %d", intVal)
	}

	if !IsTrueImmediate(objVars[1]) {
		t.Errorf("Expected instance variable 1 to be true immediate")
	}

	if objVars[2].Type() != OBJ_STRING {
		t.Errorf("Expected instance variable 2 to be a string, got %d", objVars[2].Type())
	}
	strVal := GetStringValue(objVars[2])
	if strVal != "Hello, Raw Memory!" {
		t.Errorf("Expected instance variable 2 to be 'Hello, Raw Memory!', got '%s'", strVal)
	}

	// 8. Test garbage collection
	// First, set up forwarding pointers to nil
	obj.SetMoved(false)
	obj.SetForwardingPtr(nil)
	strObj.SetMoved(false)
	strObj.SetForwardingPtr(nil)

	// Copy the object to to-space
	copiedObj := rawMem.CopyObject(obj)
	if copiedObj == nil {
		t.Fatalf("Failed to copy object")
	}

	// Verify the object was copied
	if !rawMem.IsPointerInToSpace(unsafe.Pointer(copiedObj)) {
		t.Errorf("Expected copied object to be in to-space")
	}

	// Verify the original object has a forwarding pointer
	if !obj.Moved() {
		t.Errorf("Expected original object to be marked as moved")
	}
	if obj.ForwardingPtr() != copiedObj {
		t.Errorf("Expected original object's forwarding pointer to point to copied object")
	}

	// Verify the copied object has the same values
	if copiedObj.Type() != OBJ_INSTANCE {
		t.Errorf("Expected copied object type to be OBJ_INSTANCE, got %d", copiedObj.Type())
	}

	copiedVars := copiedObj.InstanceVars()
	if len(copiedVars) != 3 {
		t.Errorf("Expected copied object to have 3 instance variables, got %d", len(copiedVars))
	}

	// Verify instance variables in copied object
	if !IsIntegerImmediate(copiedVars[0]) {
		t.Errorf("Expected copied instance variable 0 to be an integer immediate")
	}
	copiedIntVal := GetIntegerImmediate(copiedVars[0])
	if copiedIntVal != 42 {
		t.Errorf("Expected copied instance variable 0 to be 42, got %d", copiedIntVal)
	}

	if !IsTrueImmediate(copiedVars[1]) {
		t.Errorf("Expected copied instance variable 1 to be true immediate")
	}

	// The string should also be copied
	copiedStrObj := ObjectToString(copiedVars[2])
	if copiedStrObj.Value != "Hello, Raw Memory!" {
		t.Errorf("Expected copied instance variable 2 to be 'Hello, Raw Memory!', got '%s'", copiedStrObj.Value)
	}
}

// TestRawMemoryManagerObjectUsage tests using the RawMemoryManager to allocate and use objects
func TestRawMemoryManagerObjectUsage(t *testing.T) {
	// Create a new raw memory manager
	memoryManager, err := NewRawMemoryManager()
	if err != nil {
		t.Fatalf("Failed to create raw memory manager: %v", err)
	}
	defer memoryManager.Close()

	// 1. Allocate a simple object
	obj := memoryManager.AllocateObject()
	if obj == nil {
		t.Fatalf("Failed to allocate object")
	}

	// 2. Initialize the object with some values
	obj.SetType(OBJ_INSTANCE)

	// 3. Allocate instance variables array
	instanceVars := memoryManager.AllocateObjectArray(3)
	if instanceVars == nil {
		t.Fatalf("Failed to allocate instance variables array")
	}
	obj.SetInstanceVars(instanceVars)

	// 4. Create some values to store in the instance variables
	// For immediate values, we don't need to allocate memory
	intValue := MakeIntegerImmediate(42)
	boolValue := MakeTrueImmediate()

	// For a string, we need to allocate memory
	strObj := memoryManager.AllocateString()
	if strObj == nil {
		t.Fatalf("Failed to allocate string")
	}
	strObj.SetType(OBJ_STRING)
	strObj.Value = "Hello, Memory Manager!"
	stringValue := StringToObject(strObj)

	// 5. Store values in instance variables
	instanceVars[0] = intValue
	instanceVars[1] = boolValue
	instanceVars[2] = stringValue

	// 6. Verify instance variables
	objVars := obj.InstanceVars()
	if !IsIntegerImmediate(objVars[0]) {
		t.Errorf("Expected instance variable 0 to be an integer immediate")
	}
	intVal := GetIntegerImmediate(objVars[0])
	if intVal != 42 {
		t.Errorf("Expected instance variable 0 to be 42, got %d", intVal)
	}

	if !IsTrueImmediate(objVars[1]) {
		t.Errorf("Expected instance variable 1 to be true immediate")
	}

	if objVars[2].Type() != OBJ_STRING {
		t.Errorf("Expected instance variable 2 to be a string, got %d", objVars[2].Type())
	}
	strVal := GetStringValue(objVars[2])
	if strVal != "Hello, Memory Manager!" {
		t.Errorf("Expected instance variable 2 to be 'Hello, Memory Manager!', got '%s'", strVal)
	}

	// 7. Verify the object is in the raw memory's from-space
	if !memoryManager.RawMem.IsPointerInFromSpace(unsafe.Pointer(obj)) {
		t.Errorf("Expected object to be in from-space")
	}
}

// TestVMWithRawMemoryObjectUsage tests using the VMWithRawMemory to create and use objects
func TestVMWithRawMemoryObjectUsage(t *testing.T) {
	// Skip this test for now as it requires more implementation
	t.Skip("Skipping test until VMWithRawMemory is fully implemented")

	// The test would look like this:
	/*
		// Create a new VM with raw memory
		vm, err := NewVMWithRawMemory()
		if err != nil {
			t.Fatalf("Failed to create VM with raw memory: %v", err)
		}
		defer vm.Close()

		// Create objects and verify they're in raw memory
		// ...
	*/
}
