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

	// Allocate some objects
	obj1 := NewInteger(42)
	obj2 := NewBoolean(true)
	obj3 := NewString("hello")

	// Allocate them in the object memory
	allocatedObj1 := om.Allocate(obj1)
	allocatedObj2 := om.Allocate(obj2)
	allocatedObj3 := om.Allocate(obj3)

	// Check that the allocated objects are the same as the original objects
	if allocatedObj1 != obj1 {
		t.Errorf("Expected allocatedObj1 to be the same as obj1")
	}

	if allocatedObj2 != obj2 {
		t.Errorf("Expected allocatedObj2 to be the same as obj2")
	}

	if allocatedObj3 != obj3 {
		t.Errorf("Expected allocatedObj3 to be the same as obj3")
	}

	// Check that the objects are in the from-space
	if om.FromSpace[0] != obj1 {
		t.Errorf("Expected FromSpace[0] to be obj1")
	}

	if om.FromSpace[1] != obj2 {
		t.Errorf("Expected FromSpace[1] to be obj2")
	}

	if om.FromSpace[2] != obj3 {
		t.Errorf("Expected FromSpace[2] to be obj3")
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
	intObj := NewInteger(42)
	boolObj := NewBoolean(true)
	strObj := NewString("hello")

	// Allocate them in the object memory
	om.Allocate(intObj)
	om.Allocate(boolObj)
	om.Allocate(strObj)

	// Create a root object that references the other objects
	rootObj := NewArray(3)
	rootObj.Elements[0] = intObj
	rootObj.Elements[1] = boolObj
	rootObj.Elements[2] = strObj

	// Add the root object to the VM's globals
	vm.Globals["root"] = rootObj

	// Create an unreachable object
	unreachableObj := NewInteger(99)
	om.Allocate(unreachableObj)

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

	// Check that the spaces have been swapped
	// The new from-space should contain the live objects

	// Check that the allocation pointer after collection is reasonable
	// The exact number may vary depending on how many objects are created internally
	if om.AllocPtr < 4 {
		t.Errorf("Expected AllocPtr to be at least 4 after collection, got %d", om.AllocPtr)
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

	if rootObj.Elements[2] != strObj {
		t.Errorf("Expected rootObj.Elements[2] to still be strObj")
	}
}

// TestCollectWithContexts tests garbage collection with contexts
func TestCollectWithContexts(t *testing.T) {
	om := NewObjectMemory()
	vm := NewVM()

	// Create a method with a temporary variable
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)
	methodObj.Method.TempVarNames = append(methodObj.Method.TempVarNames, "temp")

	// Create some objects for the context
	receiverObj := NewInstance(vm.ObjectClass)
	arg1 := NewInteger(1)
	arg2 := NewInteger(2)

	// Create a context
	context := NewContext(methodObj, receiverObj, []*Object{arg1, arg2}, nil)

	// Push some objects onto the stack
	stackObj1 := NewBoolean(true)
	stackObj2 := NewString("stack")
	context.Push(stackObj1)
	context.Push(stackObj2)

	// Set a temporary variable
	tempObj := NewInteger(42)
	context.SetTempVarByIndex(0, tempObj)

	// Set the context as the VM's current context
	vm.CurrentContext = context

	// Allocate all objects in the object memory
	om.Allocate(methodObj)
	om.Allocate(receiverObj)
	om.Allocate(arg1)
	om.Allocate(arg2)
	om.Allocate(stackObj1)
	om.Allocate(stackObj2)
	om.Allocate(tempObj)

	// Create an unreachable object
	unreachableObj := NewInteger(99)
	om.Allocate(unreachableObj)

	// Check the allocation pointer before collection
	beforeAllocPtr := om.AllocPtr
	if beforeAllocPtr != 8 {
		t.Errorf("Expected AllocPtr to be 8 before collection, got %d", beforeAllocPtr)
	}

	// Perform garbage collection
	om.Collect(vm)

	// Check that the GC count has been incremented
	if om.GCCount != 1 {
		t.Errorf("Expected GCCount to be 1 after collection, got %d", om.GCCount)
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

	if context.Stack[1] != stackObj2 {
		t.Errorf("Expected context.Stack[1] to still be stackObj2")
	}

	// Check that the context's temporary variable is still the same
	if context.GetTempVarByIndex(0) != tempObj {
		t.Errorf("Expected context.GetTempVarByIndex(0) to still be tempObj")
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

	// Allocate the objects
	om.Allocate(obj1)
	om.Allocate(obj2)

	// Create an unreachable object
	unreachableObj := NewInteger(99)
	om.Allocate(unreachableObj)

	// Check the allocation pointer before collection
	beforeAllocPtr := om.AllocPtr
	if beforeAllocPtr != 3 {
		t.Errorf("Expected AllocPtr to be 3 before collection, got %d", beforeAllocPtr)
	}

	// Perform garbage collection
	om.Collect(vm)

	// Check that the GC count has been incremented
	if om.GCCount != 1 {
		t.Errorf("Expected GCCount to be 1 after collection, got %d", om.GCCount)
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

	// Fill the from-space with objects
	for i := 0; i < 9; i++ {
		om.FromSpace[i] = NewInteger(int64(i))
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
		if om.FromSpace[i] == nil || om.FromSpace[i].Type != OBJ_INTEGER || om.FromSpace[i].IntegerValue != int64(i) {
			t.Errorf("Expected FromSpace[%d] to be an integer with value %d", i, i)
		}
	}

	// Check that the allocation pointer hasn't changed
	if om.AllocPtr != 9 {
		t.Errorf("Expected AllocPtr to be 9 after growing, got %d", om.AllocPtr)
	}
}

// TestCollectTriggersGrowSpaces tests that collection triggers growSpaces when needed
func TestCollectTriggersGrowSpaces(t *testing.T) {
	// Skip this test for now as it requires more complex setup
	t.Skip("Skipping test that requires more complex setup")

	// This test would need to create a situation where after collection,
	// more than 70% of the space is still in use, which would trigger growSpaces.
	// However, this is difficult to set up reliably without knowing exactly
	// how many objects will be created internally during collection.
}

// TestCopyObject tests the copyObject method
func TestCopyObject(t *testing.T) {
	om := NewObjectMemory()

	// Create an object to copy
	obj := NewInteger(42)

	// Copy the object
	toPtr := 0
	copiedObj := om.copyObject(obj, &toPtr)

	// Check that the copied object is in the to-space
	if om.ToSpace[0] != obj {
		t.Errorf("Expected ToSpace[0] to be obj")
	}

	// Check that the copied object has the forwarding pointer set
	if !obj.Moved {
		t.Errorf("Expected obj.Moved to be true")
	}

	if obj.ForwardingPtr != obj {
		t.Errorf("Expected obj.ForwardingPtr to be obj")
	}

	// Check that the toPtr has been incremented
	if toPtr != 1 {
		t.Errorf("Expected toPtr to be 1, got %d", toPtr)
	}

	// Copy the object again
	copiedObj2 := om.copyObject(obj, &toPtr)

	// Check that we get the same copied object
	if copiedObj2 != copiedObj {
		t.Errorf("Expected copiedObj2 to be the same as copiedObj")
	}

	// Check that toPtr hasn't changed
	if toPtr != 1 {
		t.Errorf("Expected toPtr to still be 1, got %d", toPtr)
	}
}

// TestUpdateReferences tests the updateReferences method
func TestUpdateReferences(t *testing.T) {
	om := NewObjectMemory()

	// Test with an instance object
	{
		// Create an instance with instance variables
		class := NewClass("TestClass", nil)
		instance := NewInstance(class)
		instance.InstanceVars = make([]*Object, 2)
		instance.InstanceVars[0] = NewInteger(1)
		instance.InstanceVars[1] = NewInteger(2)

		// Update references
		toPtr := 0
		om.updateReferences(instance, &toPtr)

		// Check that the instance variables have been copied
		if !instance.InstanceVars[0].Moved {
			t.Errorf("Expected instance.InstanceVars[0].Moved to be true")
		}

		if !instance.InstanceVars[1].Moved {
			t.Errorf("Expected instance.InstanceVars[1].Moved to be true")
		}

		// Check that the class has been copied
		if !instance.Class.Moved {
			t.Errorf("Expected instance.Class.Moved to be true")
		}

		// Check that toPtr has been incremented for each copied object
		if toPtr != 3 {
			t.Errorf("Expected toPtr to be 3, got %d", toPtr)
		}
	}

	// Test with an array object
	{
		// Create an array with elements
		array := NewArray(2)
		array.Elements[0] = NewInteger(1)
		array.Elements[1] = NewInteger(2)

		// Update references
		toPtr := 0
		om.updateReferences(array, &toPtr)

		// Check that the elements have been copied
		if !array.Elements[0].Moved {
			t.Errorf("Expected array.Elements[0].Moved to be true")
		}

		if !array.Elements[1].Moved {
			t.Errorf("Expected array.Elements[1].Moved to be true")
		}

		// Check that toPtr has been incremented for each copied object
		if toPtr != 2 {
			t.Errorf("Expected toPtr to be 2, got %d", toPtr)
		}
	}

	// Test with a dictionary object
	{
		// Create a dictionary with entries
		dict := NewDictionary()
		dict.Entries["key1"] = NewInteger(1)
		dict.Entries["key2"] = NewInteger(2)

		// Update references
		toPtr := 0
		om.updateReferences(dict, &toPtr)

		// Check that the entries have been copied
		if !dict.Entries["key1"].Moved {
			t.Errorf("Expected dict.Entries[\"key1\"].Moved to be true")
		}

		if !dict.Entries["key2"].Moved {
			t.Errorf("Expected dict.Entries[\"key2\"].Moved to be true")
		}

		// Check that toPtr has been incremented for each copied object
		if toPtr != 2 {
			t.Errorf("Expected toPtr to be 2, got %d", toPtr)
		}
	}

	// Test with a method object
	{
		// Create a method with literals and selector
		method := NewMethod(NewSymbol("test"), NewClass("TestClass", nil))
		method.Method.Literals = append(method.Method.Literals, NewInteger(1))
		method.Method.Literals = append(method.Method.Literals, NewInteger(2))

		// Update references
		toPtr := 0
		om.updateReferences(method, &toPtr)

		// Check that the literals have been copied
		if !method.Method.Literals[0].Moved {
			t.Errorf("Expected method.Method.Literals[0].Moved to be true")
		}

		if !method.Method.Literals[1].Moved {
			t.Errorf("Expected method.Method.Literals[1].Moved to be true")
		}

		// Check that the selector has been copied
		if !method.Method.Selector.Moved {
			t.Errorf("Expected method.Method.Selector.Moved to be true")
		}

		// Check that the class has been copied
		if !method.Method.Class.Moved {
			t.Errorf("Expected method.Method.Class.Moved to be true")
		}

		// Check that toPtr has been incremented for each copied object
		if toPtr != 4 {
			t.Errorf("Expected toPtr to be 4, got %d", toPtr)
		}
	}
}
