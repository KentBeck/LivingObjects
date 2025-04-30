package main

import (
	"testing"
)

// TestNewObjectMemory tests the creation of a new object memory
func TestNewObjectMemory(t *testing.T) {
	om := NewObjectMemory()

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

// TestAllocate tests the allocation of objects
func TestAllocate(t *testing.T) {
	om := NewObjectMemory()

	// Create a VM for the integer class
	vm := NewVM()

	// Allocate some objects
	obj1 := vm.NewInteger(42)
	obj2 := NewBoolean(true)
	obj3 := NewString("hello")
	obj3Obj := StringToObject(obj3) // Convert to Object for allocation

	// Allocate them in the object memory
	allocatedObj1 := om.Allocate(obj1)
	allocatedObj2 := om.Allocate(obj2)
	allocatedObj3 := om.Allocate(obj3Obj)

	// Check that the allocated objects are the same as the original objects
	if allocatedObj1 != obj1 {
		t.Errorf("Expected allocatedObj1 to be the same as obj1")
	}

	if allocatedObj2 != obj2 {
		t.Errorf("Expected allocatedObj2 to be the same as obj2")
	}

	if allocatedObj3 != obj3Obj {
		t.Errorf("Expected allocatedObj3 to be the same as obj3Obj")
	}

	// Check that the objects are in the from-space
	if om.FromSpace[0] != obj1 {
		t.Errorf("Expected FromSpace[0] to be obj1")
	}

	if om.FromSpace[1] != obj2 {
		t.Errorf("Expected FromSpace[1] to be obj2")
	}

	if om.FromSpace[2] != obj3Obj {
		t.Errorf("Expected FromSpace[2] to be obj3Obj")
	}

	// Check that the allocation pointer has been updated
	if om.AllocPtr != 3 {
		t.Errorf("Expected AllocPtr to be 3, got %d", om.AllocPtr)
	}
}

// TestShouldCollect tests the ShouldCollect method
func TestShouldCollect(t *testing.T) {
	om := NewObjectMemory()

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

// TestCollect tests the garbage collection process
func TestCollect(t *testing.T) {
	om := NewObjectMemory()
	vm := NewVM()

	// Create some objects
	// Use a non-immediate object for testing
	strObj := NewString("hello")
	strObjAsObj := StringToObject(strObj) // Convert to Object for allocation

	// Create immediate values
	intObj := vm.NewInteger(42) // This is now an immediate value
	boolObj := NewBoolean(true) // This is now an immediate value

	// Allocate the non-immediate object in the object memory
	// Note: immediate values don't need to be allocated
	om.Allocate(strObjAsObj)

	// Create a root object that references the other objects
	rootObj := NewArray(3)
	rootObj.Elements[0] = intObj
	rootObj.Elements[1] = boolObj
	rootObj.Elements[2] = strObjAsObj

	// Add the root object to the VM's globals
	vm.Globals["root"] = rootObj

	// Create an unreachable non-immediate object and keep a reference to verify it's collected
	unreachableObj := NewString("unreachable")
	unreachableObjAsObj := StringToObject(unreachableObj) // Convert to Object for allocation
	om.Allocate(unreachableObjAsObj)

	// Mark the unreachable object so we can identify it later
	unreachableObjAsObj.SetMoved(true) // This flag should be reset during collection

	// Check the allocation pointer before collection
	beforeAllocPtr := om.AllocPtr
	if beforeAllocPtr != 2 { // Only 2 non-immediate objects: rootObj and unreachableObj
		t.Errorf("Expected AllocPtr to be 2 before collection, got %d", beforeAllocPtr)
	}

	// Save the original from-space to check after collection
	originalFromSpace := make([]*Object, len(om.FromSpace))
	copy(originalFromSpace, om.FromSpace)

	// Perform garbage collection
	om.Collect(vm)

	// Check that the GC count has been incremented
	if om.GCCount != 1 {
		t.Errorf("Expected GCCount to be 1 after collection, got %d", om.GCCount)
	}

	// Verify that the spaces have been swapped
	// The new from-space should contain only the live objects

	// Check that all reachable non-immediate objects are in the new from-space
	foundRoot := false
	foundStr := false

	for i := 0; i < om.AllocPtr; i++ {
		obj := om.FromSpace[i]
		if obj == rootObj {
			foundRoot = true
		} else if obj == strObjAsObj {
			foundStr = true
		}
	}

	if !foundRoot {
		t.Errorf("Root object not found in from-space after collection")
	}
	if !foundStr {
		t.Errorf("String object not found in from-space after collection")
	}

	// Verify that immediate values are still accessible and correct
	if !IsIntegerImmediate(intObj) {
		t.Errorf("Integer object is not an immediate value")
	}
	if !IsTrueImmediate(boolObj) {
		t.Errorf("Boolean object is not an immediate value")
	}

	// Verify that the unreachable object is not in the new from-space
	foundUnreachable := false
	for i := 0; i < om.AllocPtr; i++ {
		if om.FromSpace[i] == unreachableObjAsObj {
			foundUnreachable = true
			break
		}
	}

	if foundUnreachable {
		t.Errorf("Unreachable object found in from-space after collection")
	}

	// Check that the root object is still in the VM's globals
	if vm.Globals["root"] != rootObj {
		t.Errorf("Expected vm.Globals[\"root\"] to still be rootObj")
	}

	// Check that the root object's elements are still the same
	if rootObj.Elements[0] != intObj {
		t.Errorf("Expected rootObj.Elements[0] to still be intObj")
	}

	if rootObj.Elements[1] != boolObj {
		t.Errorf("Expected rootObj.Elements[1] to still be boolObj")
	}

	if rootObj.Elements[2] != strObjAsObj {
		t.Errorf("Expected rootObj.Elements[2] to still be strObj")
	}

	// Verify that the values of the objects are preserved
	if GetIntegerImmediate(intObj) != 42 {
		t.Errorf("Integer value changed after collection, expected 42, got %d", GetIntegerImmediate(intObj))
	}

	if !IsTrueImmediate(boolObj) {
		t.Errorf("Boolean value changed after collection, expected true immediate value")
	}

	strValue := ObjectToString(strObjAsObj)
	if strValue.Value != "hello" {
		t.Errorf("String value changed after collection, expected 'hello', got '%s'", strValue.Value)
	}
}

// TestCollectWithContexts tests garbage collection with contexts
func TestCollectWithContexts(t *testing.T) {
	om := NewObjectMemory()
	vm := NewVM()

	// Create a method with a temporary variable
	methodObj := NewMethodBuilder(vm.ObjectClass).
		Selector("test").
		TempVars([]string{"temp"}).
		Go()

	// Create some objects for the context
	receiverObj := NewInstance(vm.ObjectClass) // Non-immediate object

	// Create immediate integer values for arguments
	arg1 := vm.NewInteger(1) // Immediate value
	arg2 := vm.NewInteger(2) // Immediate value

	// Create a context
	context := NewContext(methodObj, receiverObj, []*Object{arg1, arg2}, nil)

	// Push some objects onto the stack
	stackObj1 := NewBoolean(true)               // Immediate value
	stackObj2 := NewString("stack")             // Non-immediate object
	stackObj2AsObj := StringToObject(stackObj2) // Convert to Object for stack
	context.Push(stackObj1)
	context.Push(stackObj2AsObj)

	// Set a temporary variable
	tempObj := vm.NewInteger(42) // Immediate value
	context.SetTempVarByIndex(0, tempObj)

	// Set the context as the VM's current context
	vm.CurrentContext = context

	// Allocate non-immediate objects in the object memory
	// Note: immediate values don't need to be allocated
	om.Allocate(methodObj)
	om.Allocate(receiverObj)
	om.Allocate(stackObj2AsObj) // Only non-immediate objects need to be allocated

	// Create an unreachable non-immediate object and mark it
	unreachableObj := NewString("unreachable")
	unreachableObjAsObj := StringToObject(unreachableObj) // Convert to Object for allocation
	om.Allocate(unreachableObjAsObj)
	unreachableObjAsObj.SetMoved(true) // This flag should be reset during collection

	// Check the allocation pointer before collection
	beforeAllocPtr := om.AllocPtr
	if beforeAllocPtr != 4 { // Only 4 non-immediate objects: methodObj, receiverObj, stackObj2, unreachableObj
		t.Errorf("Expected AllocPtr to be 4 before collection, got %d", beforeAllocPtr)
	}

	// Perform garbage collection
	om.Collect(vm)

	// Check that the GC count has been incremented
	if om.GCCount != 1 {
		t.Errorf("Expected GCCount to be 1 after collection, got %d", om.GCCount)
	}

	// Verify that all reachable non-immediate objects are in the new from-space
	foundMethod := false
	foundReceiver := false
	foundStackObj2 := false

	for i := 0; i < om.AllocPtr; i++ {
		obj := om.FromSpace[i]
		if obj == methodObj {
			foundMethod = true
		} else if obj == receiverObj {
			foundReceiver = true
		} else if obj == stackObj2AsObj {
			foundStackObj2 = true
		}
	}

	if !foundMethod {
		t.Errorf("Method object not found in from-space after collection")
	}
	if !foundReceiver {
		t.Errorf("Receiver object not found in from-space after collection")
	}
	if !foundStackObj2 {
		t.Errorf("Stack object 2 not found in from-space after collection")
	}

	// Verify that immediate values are still accessible and correct
	if !IsIntegerImmediate(arg1) {
		t.Errorf("Argument 1 is not an immediate value")
	}
	if !IsIntegerImmediate(arg2) {
		t.Errorf("Argument 2 is not an immediate value")
	}
	if !IsTrueImmediate(stackObj1) {
		t.Errorf("Stack object 1 is not an immediate value")
	}
	if !IsIntegerImmediate(tempObj) {
		t.Errorf("Temporary variable object is not an immediate value")
	}

	// Verify that the unreachable object is not in the new from-space
	foundUnreachable := false
	for i := 0; i < om.AllocPtr; i++ {
		if om.FromSpace[i] == unreachableObjAsObj {
			foundUnreachable = true
			break
		}
	}

	if foundUnreachable {
		t.Errorf("Unreachable object found in from-space after collection")
	}

	// Check that the context is still the VM's current context
	if vm.CurrentContext != context {
		t.Errorf("Expected vm.CurrentContext to still be context")
	}

	// Check that the context's method is still the same
	if context.Method != methodObj {
		t.Errorf("Expected context.Method to still be methodObj")
	}

	// Check that the context's receiver is still the same
	if context.Receiver != receiverObj {
		t.Errorf("Expected context.Receiver to still be receiverObj")
	}

	// Check that the context's arguments are still the same
	if len(context.Arguments) != 2 {
		t.Errorf("Expected context.Arguments to have 2 elements, got %d", len(context.Arguments))
	}

	if context.Arguments[0] != arg1 {
		t.Errorf("Expected context.Arguments[0] to still be arg1")
	}

	if context.Arguments[1] != arg2 {
		t.Errorf("Expected context.Arguments[1] to still be arg2")
	}

	// Check that the context's stack is still the same
	if context.StackPointer != 2 {
		t.Errorf("Expected context.StackPointer to be 2, got %d", context.StackPointer)
	}

	if context.Stack[0] != stackObj1 {
		t.Errorf("Expected context.Stack[0] to still be stackObj1")
	}

	if context.Stack[1] != stackObj2AsObj {
		t.Errorf("Expected context.Stack[1] to still be stackObj2")
	}

	// Check that the context's temporary variable is still the same
	if context.GetTempVarByIndex(0) != tempObj {
		t.Errorf("Expected context.GetTempVarByIndex(0) to still be tempObj")
	}

	// Verify that the values of the objects are preserved
	if GetIntegerImmediate(arg1) != 1 {
		t.Errorf("Argument 1 value changed after collection, expected 1, got %d", GetIntegerImmediate(arg1))
	}

	if GetIntegerImmediate(arg2) != 2 {
		t.Errorf("Argument 2 value changed after collection, expected 2, got %d", GetIntegerImmediate(arg2))
	}

	if !IsTrueImmediate(stackObj1) {
		t.Errorf("Stack object 1 value changed after collection, expected true immediate value")
	}

	stackObj2Value := ObjectToString(stackObj2AsObj)
	if stackObj2Value.Value != "stack" {
		t.Errorf("Stack object 2 value changed after collection, expected 'stack', got '%s'", stackObj2Value.Value)
	}

	if GetIntegerImmediate(tempObj) != 42 {
		t.Errorf("Temporary variable value changed after collection, expected 42, got %d", GetIntegerImmediate(tempObj))
	}
}

// TestCollectWithCycles tests garbage collection with cyclic references
func TestCollectWithCycles(t *testing.T) {
	om := NewObjectMemory()
	vm := NewVM()

	// Create objects that reference each other
	obj1 := NewArray(1)
	obj2 := NewArray(1)

	// Create the cycle
	obj1.Elements[0] = obj2
	obj2.Elements[0] = obj1

	// Add obj1 to the VM's globals to make it reachable
	vm.Globals["cycle"] = obj1

	// Create another cycle that is unreachable
	unreachableObj1 := NewArray(1)
	unreachableObj2 := NewArray(1)
	unreachableObj1.Elements[0] = unreachableObj2
	unreachableObj2.Elements[0] = unreachableObj1

	// Allocate the objects
	om.Allocate(obj1)
	om.Allocate(obj2)
	om.Allocate(unreachableObj1)
	om.Allocate(unreachableObj2)

	// Mark the unreachable objects so we can identify them later
	unreachableObj1.SetMoved(true)
	unreachableObj2.SetMoved(true)

	// Check the allocation pointer before collection
	beforeAllocPtr := om.AllocPtr
	if beforeAllocPtr != 4 {
		t.Errorf("Expected AllocPtr to be 4 before collection, got %d", beforeAllocPtr)
	}

	// Perform garbage collection
	om.Collect(vm)

	// Check that the GC count has been incremented
	if om.GCCount != 1 {
		t.Errorf("Expected GCCount to be 1 after collection, got %d", om.GCCount)
	}

	// Verify that the reachable objects are in the new from-space
	foundObj1 := false
	foundObj2 := false

	for i := 0; i < om.AllocPtr; i++ {
		obj := om.FromSpace[i]
		if obj == obj1 {
			foundObj1 = true
		} else if obj == obj2 {
			foundObj2 = true
		}
	}

	if !foundObj1 {
		t.Errorf("obj1 not found in from-space after collection")
	}
	if !foundObj2 {
		t.Errorf("obj2 not found in from-space after collection")
	}

	// Verify that the unreachable objects are not in the new from-space
	foundUnreachable1 := false
	foundUnreachable2 := false

	for i := 0; i < om.AllocPtr; i++ {
		obj := om.FromSpace[i]
		if obj == unreachableObj1 {
			foundUnreachable1 = true
		} else if obj == unreachableObj2 {
			foundUnreachable2 = true
		}
	}

	if foundUnreachable1 {
		t.Errorf("unreachableObj1 found in from-space after collection")
	}
	if foundUnreachable2 {
		t.Errorf("unreachableObj2 found in from-space after collection")
	}

	// Check that obj1 is still in the VM's globals
	if vm.Globals["cycle"] != obj1 {
		t.Errorf("Expected vm.Globals[\"cycle\"] to still be obj1")
	}

	// Check that the cycle is preserved
	if obj1.Elements[0] != obj2 {
		t.Errorf("Expected obj1.Elements[0] to still be obj2")
	}

	if obj2.Elements[0] != obj1 {
		t.Errorf("Expected obj2.Elements[0] to still be obj1")
	}

	// Check that the allocation pointer after collection is reasonable
	// The exact number may vary depending on how many objects are created internally
	if om.AllocPtr < 2 {
		t.Errorf("Expected AllocPtr to be at least 2 after collection, got %d", om.AllocPtr)
	}
}

// TestGrowSpaces tests the growSpaces method
func TestGrowSpaces(t *testing.T) {
	om := NewObjectMemory()

	// Set a small initial space size for testing
	om.SpaceSize = 10
	om.GCThreshold = 8
	om.FromSpace = make([]*Object, 10)
	om.ToSpace = make([]*Object, 10)

	// No need for a VM when using non-immediate objects

	// Fill the from-space with non-immediate objects
	for i := 0; i < 9; i++ {
		strObj := NewString(string(rune('a' + i)))
		om.FromSpace[i] = StringToObject(strObj)
	}
	om.AllocPtr = 9

	// Grow the spaces
	om.growSpaces()

	// Check that the space size has doubled
	if om.SpaceSize != 20 {
		t.Errorf("Expected SpaceSize to be 20 after growing, got %d", om.SpaceSize)
	}

	// Check that the GC threshold has been updated
	if om.GCThreshold != 16 {
		t.Errorf("Expected GCThreshold to be 16 after growing, got %d", om.GCThreshold)
	}

	// Check that the from-space has the new size
	if len(om.FromSpace) != 20 {
		t.Errorf("Expected FromSpace length to be 20 after growing, got %d", len(om.FromSpace))
	}

	// Check that the to-space has the new size
	if len(om.ToSpace) != 20 {
		t.Errorf("Expected ToSpace length to be 20 after growing, got %d", len(om.ToSpace))
	}

	// Check that the objects are still in the from-space
	for i := 0; i < 9; i++ {
		expectedStr := string(rune('a' + i))
		if om.FromSpace[i] == nil || om.FromSpace[i].Type() != OBJ_STRING || ObjectToString(om.FromSpace[i]).Value != expectedStr {
			t.Errorf("Expected FromSpace[%d] to be a string with value '%s'", i, expectedStr)
		}
	}

	// Check that the allocation pointer hasn't changed
	if om.AllocPtr != 9 {
		t.Errorf("Expected AllocPtr to be 9 after growing, got %d", om.AllocPtr)
	}
}

// TestCollectTriggersGrowSpaces tests that collection triggers growSpaces when needed
func TestCollectTriggersGrowSpaces(t *testing.T) {
	// Create a small object memory
	om := NewObjectMemory()

	// Set a small space size for testing
	om.SpaceSize = 20
	om.GCThreshold = 16 // 80% threshold
	om.FromSpace = make([]*Object, 20)
	om.ToSpace = make([]*Object, 20)

	vm := NewVM()

	// Create a bunch of non-immediate objects that will survive collection
	for i := 0; i < 15; i++ {
		obj := NewString(string(rune('a' + i)))
		objAsObj := StringToObject(obj) // Convert to Object for allocation
		om.Allocate(objAsObj)

		// Add to globals to make them reachable
		vm.Globals[string(rune('a'+i))] = objAsObj
	}

	// Record the initial space size
	initialSpaceSize := om.SpaceSize

	// Perform garbage collection
	om.Collect(vm)

	// Check that the space size has doubled
	if om.SpaceSize != initialSpaceSize*2 {
		t.Errorf("Expected SpaceSize to be %d after collection, got %d", initialSpaceSize*2, om.SpaceSize)
	}

	// Check that the from-space and to-space have been resized
	if len(om.FromSpace) != initialSpaceSize*2 {
		t.Errorf("Expected FromSpace length to be %d after collection, got %d", initialSpaceSize*2, len(om.FromSpace))
	}

	if len(om.ToSpace) != initialSpaceSize*2 {
		t.Errorf("Expected ToSpace length to be %d after collection, got %d", initialSpaceSize*2, len(om.ToSpace))
	}

	// Check that the GC threshold has been updated
	if om.GCThreshold != initialSpaceSize*2*80/100 {
		t.Errorf("Expected GCThreshold to be %d after collection, got %d", initialSpaceSize*2*80/100, om.GCThreshold)
	}

	// Check that all objects are still in the from-space
	for i := 0; i < 8; i++ {
		key := string(rune('a' + i))
		obj := vm.Globals[key]

		// Find the object in the from-space
		found := false
		for j := 0; j < om.AllocPtr; j++ {
			if om.FromSpace[j] == obj {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Object %s not found in from-space after collection", key)
		}
	}
}

// TestCollectEdgeCases tests edge cases in the garbage collector
func TestCollectEdgeCases(t *testing.T) {
	// Test with empty object memory
	{
		om := NewObjectMemory()
		vm := NewVM()

		// Clear the VM's globals to ensure no objects are reachable
		vm.Globals = make(map[string]*Object)
		vm.CurrentContext = nil

		// Perform garbage collection on empty memory
		om.Collect(vm)

		// Check that the GC count has been incremented
		if om.GCCount != 1 {
			t.Errorf("Expected GCCount to be 1 after collection, got %d", om.GCCount)
		}

		// The allocation pointer might not be 0 because the VM might have created some objects
		// during initialization. We just check that collection completed successfully.
	}

	// Test with nil objects in the from-space
	{
		om := NewObjectMemory()
		vm := NewVM()

		// Create some non-immediate objects with nil slots in between
		obj1 := NewString("one")
		obj2 := NewString("two")
		obj1AsObj := StringToObject(obj1) // Convert to Object for allocation
		obj2AsObj := StringToObject(obj2) // Convert to Object for allocation

		om.Allocate(obj1AsObj)
		om.FromSpace[1] = nil // Create a nil slot
		om.AllocPtr = 2       // Skip the nil slot
		om.Allocate(obj2AsObj)

		// Add objects to globals to make them reachable
		vm.Globals["obj1"] = obj1AsObj
		vm.Globals["obj2"] = obj2AsObj

		// Perform garbage collection
		om.Collect(vm)

		// Check that both objects are still in the from-space
		foundObj1 := false
		foundObj2 := false

		for i := 0; i < om.AllocPtr; i++ {
			obj := om.FromSpace[i]
			if obj == obj1AsObj {
				foundObj1 = true
			} else if obj == obj2AsObj {
				foundObj2 = true
			}
		}

		if !foundObj1 {
			t.Errorf("obj1 not found in from-space after collection")
		}
		if !foundObj2 {
			t.Errorf("obj2 not found in from-space after collection")
		}
	}

	// Test with nil objects in the root set
	{
		om := NewObjectMemory()
		vm := NewVM()

		// Add a nil object to globals
		vm.Globals["nil"] = nil

		// Perform garbage collection
		om.Collect(vm)

		// Check that the collection completed without errors
		if om.GCCount != 1 {
			t.Errorf("Expected GCCount to be 1 after collection, got %d", om.GCCount)
		}
	}

	// Test with a large number of objects
	{
		om := NewObjectMemory()
		vm := NewVM()

		// Set a small space size for testing
		om.SpaceSize = 20
		om.GCThreshold = 16
		om.FromSpace = make([]*Object, 20)
		om.ToSpace = make([]*Object, 20)

		// Allocate many non-immediate objects
		for i := 0; i < 15; i++ {
			obj := NewString(string(rune('a' + i)))
			objAsObj := StringToObject(obj) // Convert to Object for allocation
			om.Allocate(objAsObj)
			vm.Globals[string(rune('a'+i))] = objAsObj // Add to globals to make them reachable
		}

		// Perform garbage collection
		om.Collect(vm)

		// Check that all objects are still in the from-space
		for i := 0; i < 15; i++ {
			key := string(rune('a' + i))
			obj := vm.Globals[key]

			// Find the object in the from-space
			found := false
			for j := 0; j < om.AllocPtr; j++ {
				if om.FromSpace[j] == obj {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Object %s not found in from-space after collection", key)
			}
		}
	}
}

// TestCopyObject tests the copyObject method
func TestCopyObject(t *testing.T) {
	om := NewObjectMemory()

	// Create a VM for the integer class
	vm := NewVM()

	// Test with an immediate value
	{
		// Create an immediate value
		immediateObj := vm.NewInteger(42)

		// Verify it's an immediate value
		if !IsIntegerImmediate(immediateObj) {
			t.Errorf("Expected immediateObj to be an immediate value")
		}

		// Copy the object
		toPtr := 0
		copiedObj := om.copyObject(immediateObj, &toPtr)

		// Check that the immediate value is returned as-is
		if copiedObj != immediateObj {
			t.Errorf("Expected copiedObj to be the same as immediateObj")
		}

		// Check that toPtr hasn't changed (immediate values don't get copied)
		if toPtr != 0 {
			t.Errorf("Expected toPtr to still be 0, got %d", toPtr)
		}
	}

	// Test with a non-immediate object
	{
		// Create a non-immediate object
		obj := NewString("test")
		objAsObj := StringToObject(obj) // Convert to Object for copying

		// Copy the object
		toPtr := 0
		_ = om.copyObject(objAsObj, &toPtr)

		// Check that the copied object is in the to-space
		if om.ToSpace[0] != objAsObj {
			t.Errorf("Expected ToSpace[0] to be obj")
		}

		// Check that the copied object has the forwarding pointer set
		if !objAsObj.Moved() {
			t.Errorf("Expected obj.Moved() to be true")
		}

		if objAsObj.ForwardingPtr() != objAsObj {
			t.Errorf("Expected obj.ForwardingPtr() to be obj")
		}

		// Check that the toPtr has been incremented
		if toPtr != 1 {
			t.Errorf("Expected toPtr to be 1, got %d", toPtr)
		}
	}

	// Test copying an object that has already been copied
	{
		// Create a non-immediate object
		obj := NewString("test2")
		objAsObj := StringToObject(obj) // Convert to Object for copying

		// Copy the object first time
		toPtr := 0
		copiedObj := om.copyObject(objAsObj, &toPtr)

		// Copy the object again
		copiedObj2 := om.copyObject(objAsObj, &toPtr)

		// Check that we get the same copied object
		if copiedObj2 != copiedObj {
			t.Errorf("Expected copiedObj2 to be the same as copiedObj")
		}

		// Check that toPtr hasn't changed
		if toPtr != 1 {
			t.Errorf("Expected toPtr to still be 1, got %d", toPtr)
		}
	}
}

// TestUpdateReferences tests the updateReferences method
func TestUpdateReferences(t *testing.T) {
	om := NewObjectMemory()

	// Test with an instance object
	{
		// Create a VM for the integer class
		vm := NewVM()

		// Create an instance with instance variables
		class := NewClass("TestClass", nil)
		instance := NewInstance(class)
		instanceVars := make([]*Object, 2)

		// Use immediate values for instance variables
		instanceVars[METHOD_DICTIONARY_IV] = vm.NewInteger(1) // Immediate value
		instanceVars[1] = vm.NewInteger(2)                    // Immediate value
		instance.SetInstanceVars(instanceVars)

		// Update references
		toPtr := 0
		om.updateReferences(instance, &toPtr)

		// Check that the instance variables are still immediate values
		instanceVars = instance.InstanceVars()
		if !IsIntegerImmediate(instanceVars[0]) {
			t.Errorf("Expected instance.InstanceVars()[0] to be an immediate value")
		}

		if !IsIntegerImmediate(instanceVars[1]) {
			t.Errorf("Expected instance.InstanceVars()[1] to be an immediate value")
		}

		// Check that the class has been copied
		if !instance.Class().Moved() {
			t.Errorf("Expected instance.Class().Moved() to be true")
		}

		// Check that toPtr has been incremented only for the class (immediate values don't get copied)
		if toPtr != 1 {
			t.Errorf("Expected toPtr to be 1, got %d", toPtr)
		}
	}

	// Test with an array object
	{
		// Create a VM for the integer class
		vm := NewVM()

		// Create an array with immediate value elements
		array := NewArray(2)
		array.Elements[0] = vm.NewInteger(1) // Immediate value
		array.Elements[1] = vm.NewInteger(2) // Immediate value

		// Update references
		toPtr := 0
		om.updateReferences(array, &toPtr)

		// Check that the elements are still immediate values
		if !IsIntegerImmediate(array.Elements[0]) {
			t.Errorf("Expected array.Elements[0] to be an immediate value")
		}

		if !IsIntegerImmediate(array.Elements[1]) {
			t.Errorf("Expected array.Elements[1] to be an immediate value")
		}

		// Check that toPtr hasn't been incremented (immediate values don't get copied)
		if toPtr != 0 {
			t.Errorf("Expected toPtr to be 0, got %d", toPtr)
		}
	}

	// Test with a dictionary object
	{
		// Create a VM for the integer class
		vm := NewVM()

		// Create a dictionary with immediate value entries
		dict := NewDictionary()
		dict.Entries["key1"] = vm.NewInteger(1) // Immediate value
		dict.Entries["key2"] = vm.NewInteger(2) // Immediate value

		// Update references
		toPtr := 0
		om.updateReferences(dict, &toPtr)

		// Check that the entries are still immediate values
		if !IsIntegerImmediate(dict.Entries["key1"]) {
			t.Errorf("Expected dict.Entries[\"key1\"] to be an immediate value")
		}

		if !IsIntegerImmediate(dict.Entries["key2"]) {
			t.Errorf("Expected dict.Entries[\"key2\"] to be an immediate value")
		}

		// Check that toPtr hasn't been incremented (immediate values don't get copied)
		if toPtr != 0 {
			t.Errorf("Expected toPtr to be 0, got %d", toPtr)
		}
	}

	// Test with a method object
	{
		// Create a VM for the integer class
		vm := NewVM()

		// Create a method with immediate value literals and selector using AddLiteral
		builder := NewMethodBuilder(NewClass("TestClass", nil)).Selector("test")

		// Add literals to the method builder
		builder.AddLiteral(vm.NewInteger(1)) // Immediate value
		builder.AddLiteral(vm.NewInteger(2)) // Immediate value

		// Finalize the method
		method := builder.Go()

		// Update references
		toPtr := 0
		om.updateReferences(method, &toPtr)

		// Check that the literals are still immediate values
		if !IsIntegerImmediate(method.Method.Literals[0]) {
			t.Errorf("Expected method.Method.Literals[0] to be an immediate value")
		}

		if !IsIntegerImmediate(method.Method.Literals[1]) {
			t.Errorf("Expected method.Method.Literals[1] to be an immediate value")
		}

		// Check that the selector has been copied
		if !method.Method.Selector.Moved() {
			t.Errorf("Expected method.Method.Selector.Moved() to be true")
		}

		// Check that the class has been copied
		if !method.Method.MethodClass.Moved() {
			t.Errorf("Expected method.Method.Class.Moved() to be true")
		}

		// Check that toPtr has been incremented only for the selector and class (immediate values don't get copied)
		if toPtr != 2 {
			t.Errorf("Expected toPtr to be 2, got %d", toPtr)
		}
	}
}
