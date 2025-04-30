package main

import (
	"fmt"
	"unsafe"
)

// RawMemoryManager manages the Smalltalk object memory using raw memory allocation
type RawMemoryManager struct {
	RawMem  *RawMemory
	GCCount int
}

// NewRawMemoryManager creates a new raw memory manager
func NewRawMemoryManager() (*RawMemoryManager, error) {
	// Create a new raw memory with default space sizes (1MB each)
	rawMem, err := NewRawMemory(DefaultSpaceSize, DefaultSpaceSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create raw memory: %w", err)
	}

	return &RawMemoryManager{
		RawMem:  rawMem,
		GCCount: 0,
	}, nil
}

// Close releases all allocated memory
func (rmm *RawMemoryManager) Close() error {
	if rmm.RawMem != nil {
		return rmm.RawMem.Close()
	}
	return nil
}

// ShouldCollect returns true if garbage collection should be triggered
func (rmm *RawMemoryManager) ShouldCollect() bool {
	return rmm.RawMem.ShouldCollect()
}

// AllocateObject allocates a new Object
func (rmm *RawMemoryManager) AllocateObject() *Object {
	// Check if we need to collect garbage
	if rmm.ShouldCollect() {
		// We'll let the VM handle collection
		return nil
	}

	return rmm.RawMem.AllocateObject()
}

// AllocateString allocates a new String
func (rmm *RawMemoryManager) AllocateString() *String {
	// Check if we need to collect garbage
	if rmm.ShouldCollect() {
		// We'll let the VM handle collection
		return nil
	}

	return rmm.RawMem.AllocateString()
}

// AllocateSymbol allocates a new Symbol
func (rmm *RawMemoryManager) AllocateSymbol() *Symbol {
	// Check if we need to collect garbage
	if rmm.ShouldCollect() {
		// We'll let the VM handle collection
		return nil
	}

	return rmm.RawMem.AllocateSymbol()
}

// AllocateClass allocates a new Class
func (rmm *RawMemoryManager) AllocateClass() *Class {
	// Check if we need to collect garbage
	if rmm.ShouldCollect() {
		// We'll let the VM handle collection
		return nil
	}

	return rmm.RawMem.AllocateClass()
}

// AllocateMethod allocates a new Method
func (rmm *RawMemoryManager) AllocateMethod() *Method {
	// Check if we need to collect garbage
	if rmm.ShouldCollect() {
		// We'll let the VM handle collection
		return nil
	}

	return rmm.RawMem.AllocateMethod()
}

// AllocateObjectArray allocates a new array of Object pointers
func (rmm *RawMemoryManager) AllocateObjectArray(size int) []*Object {
	// Check if we need to collect garbage
	if rmm.ShouldCollect() {
		// We'll let the VM handle collection
		return nil
	}

	return rmm.RawMem.AllocateObjectArray(size)
}

// AllocateByteArray allocates a new byte array
func (rmm *RawMemoryManager) AllocateByteArray(size int) []byte {
	// Check if we need to collect garbage
	if rmm.ShouldCollect() {
		// We'll let the VM handle collection
		return nil
	}

	return rmm.RawMem.AllocateByteArray(size)
}

// AllocateStringMap allocates a new string map
func (rmm *RawMemoryManager) AllocateStringMap() map[string]*Object {
	// Check if we need to collect garbage
	if rmm.ShouldCollect() {
		// We'll let the VM handle collection
		return nil
	}

	return rmm.RawMem.AllocateStringMap()
}

// Collect performs a stop & copy garbage collection
func (rmm *RawMemoryManager) Collect(vm *VM) {
	rmm.GCCount++

	// Reset forwarding pointers for all objects in from-space
	// This is a bit tricky with raw memory, as we don't have a list of objects
	// We'll rely on the objects' Moved and ForwardingPtr fields

	// Start with the root set
	// Initialize the to-space pointer
	rmm.RawMem.ToSpacePtr = uintptr(unsafe.Pointer(&rmm.RawMem.ToSpace[0]))

	// If VM is nil, we can't collect garbage
	if vm == nil {
		return
	}

	// Add globals to the root set
	for _, obj := range vm.Globals {
		if obj != nil {
			copied := rmm.copyObject(obj)
			// Update the global reference
			*obj = *copied
		}
	}

	// Add the current context and its chain to the root set
	context := vm.CurrentContext
	for context != nil {
		// Copy the context's method
		if context.Method != nil {
			copied := rmm.copyObject(context.Method)
			context.Method = copied
		}

		// Copy the context's receiver
		if context.Receiver != nil {
			copied := rmm.copyObject(context.Receiver)
			context.Receiver = copied
		}

		// Copy the context's arguments
		for i, arg := range context.Arguments {
			if arg != nil {
				copied := rmm.copyObject(arg)
				context.Arguments[i] = copied
			}
		}

		// Copy the context's temporary variables
		for i, obj := range context.TempVars {
			if obj != nil {
				copied := rmm.copyObject(obj)
				context.TempVars[i] = copied
			}
		}

		// Copy the context's stack
		for i := 0; i < context.StackPointer; i++ {
			if context.Stack[i] != nil {
				copied := rmm.copyObject(context.Stack[i])
				context.Stack[i] = copied
			}
		}

		// Move to the sender context
		context = context.Sender
	}

	// Copy special objects
	// Note: NilObject, TrueObject, and FalseObject are now immediate values, so we don't need to copy them
	if vm.ObjectClass != nil {
		vm.ObjectClass = rmm.copyObject(vm.ObjectClass)
	}

	// Scan the to-space for references
	// This is tricky with raw memory, as we don't have a list of objects
	// We'll need to scan the memory from the beginning to the current to-space pointer

	// Calculate how many objects we've copied to to-space
	objectsInToSpace := int((rmm.RawMem.ToSpacePtr - uintptr(unsafe.Pointer(&rmm.RawMem.ToSpace[0]))) / uintptr(unsafe.Sizeof(Object{})))

	// Scan each object in to-space
	for i := 0; i < objectsInToSpace; i++ {
		// Calculate the address of the object
		objAddr := uintptr(unsafe.Pointer(&rmm.RawMem.ToSpace[0])) + uintptr(i)*uintptr(unsafe.Sizeof(Object{}))
		obj := (*Object)(unsafe.Pointer(objAddr))

		// Update references in the object
		rmm.updateReferences(obj)
	}

	// Swap the spaces
	rmm.RawMem.SwapSpaces()

	// No need to grow spaces for now, as we're using fixed-size spaces
	// In a more advanced implementation, we could allocate new, larger spaces
}

// copyObject copies an object to the to-space
func (rmm *RawMemoryManager) copyObject(obj *Object) *Object {
	// Use the raw memory's CopyObject method
	return rmm.RawMem.CopyObject(obj)
}

// updateReferences updates references in an object
func (rmm *RawMemoryManager) updateReferences(obj *Object) {
	// Check if it's an immediate value
	if IsImmediate(obj) {
		// Immediate values don't have references
		return
	}

	switch obj.Type() {
	case OBJ_STRING:
		// String objects don't have references to update
		return

	case OBJ_SYMBOL:
		// Symbol objects don't have references to update
		return

	case OBJ_ARRAY:
		// Update array elements
		for i, elem := range obj.Elements {
			if elem != nil {
				obj.Elements[i] = rmm.copyObject(elem)
			}
		}

	case OBJ_DICTIONARY:
		// Update dictionary entries
		for key, value := range obj.Entries {
			if value != nil {
				obj.Entries[key] = rmm.copyObject(value)
			}
		}

	case OBJ_INSTANCE:
		// Update instance variables
		instanceVars := obj.InstanceVars()
		for i, value := range instanceVars {
			if value != nil {
				instanceVars[i] = rmm.copyObject(value)
			}
		}

		// Update superclass reference
		if obj.SuperClass != nil {
			obj.SuperClass = rmm.copyObject(obj.SuperClass)
		}

		// Update class reference
		if obj.Class() != nil {
			obj.SetClass(rmm.copyObject(obj.Class()))
		}

	case OBJ_CLASS:
		// Update instance variables (which includes the method dictionary)
		instanceVars := obj.InstanceVars()
		for i, value := range instanceVars {
			if value != nil {
				instanceVars[i] = rmm.copyObject(value)
			}
		}

		// Update superclass reference
		if obj.SuperClass != nil {
			obj.SuperClass = rmm.copyObject(obj.SuperClass)
		}

	case OBJ_METHOD:
		// Update method literals
		for i, lit := range obj.Method.Literals {
			if lit != nil {
				obj.Method.Literals[i] = rmm.copyObject(lit)
			}
		}

		// Update method selector
		if obj.Method.Selector != nil {
			obj.Method.Selector = rmm.copyObject(obj.Method.Selector)
		}

		// Update method class
		if obj.Method.MethodClass != nil {
			obj.Method.MethodClass = rmm.copyObject(obj.Method.MethodClass)
		}

	case OBJ_BLOCK:
		// Update block literals
		for i, lit := range obj.Literals {
			if lit != nil {
				obj.Literals[i] = rmm.copyObject(lit)
			}
		}
	}
}
