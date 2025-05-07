package core

import (
	"unsafe"
)

// ObjectMemory manages the Smalltalk object memory with stop & copy garbage collection
type ObjectMemory struct {
	FromSpace   []*Object
	ToSpace     []*Object
	AllocPtr    int
	SpaceSize   int
	GCThreshold int
	GCCount     int
}

// VM is a forward declaration to avoid circular imports
type VM interface {
	GetGlobals() []*Object
	GetCurrentContext() interface{}
	GetObjectClass() *Class
}

// ExecutionContext is a forward declaration to avoid circular imports and name conflicts
type ExecutionContext interface {
	GetMethod() *Object
	GetReceiver() *Object
	GetArguments() []*Object
	GetTempVars() []*Object
	GetStack() []*Object
	GetStackPointer() int
	GetSender() interface{}
}

// NewObjectMemory creates a new object memory
func NewObjectMemory() *ObjectMemory {
	spaceSize := 10000 // Initial space size
	return &ObjectMemory{
		FromSpace:   make([]*Object, spaceSize),
		ToSpace:     make([]*Object, spaceSize),
		AllocPtr:    0,
		SpaceSize:   spaceSize,
		GCThreshold: spaceSize * 80 / 100, // 80% threshold
		GCCount:     0,
	}
}

// Allocate allocates a new object
func (om *ObjectMemory) Allocate(obj *Object) *Object {
	// Check if we need to collect garbage
	if om.ShouldCollect() {
		// We'll let the VM handle collection
		return obj
	}

	// Allocate the object in the from-space
	om.FromSpace[om.AllocPtr] = obj
	om.AllocPtr++

	return obj
}

// ShouldCollect returns true if garbage collection should be triggered
func (om *ObjectMemory) ShouldCollect() bool {
	return om.AllocPtr >= om.GCThreshold
}

// Collect performs a stop & copy garbage collection
func (om *ObjectMemory) Collect(vm VM) {
	om.GCCount++

	// Reset the to-space
	for i := range om.ToSpace {
		om.ToSpace[i] = nil
	}

	// Reset forwarding pointers
	for i := 0; i < om.AllocPtr; i++ {
		if om.FromSpace[i] != nil {
			om.FromSpace[i].SetMoved(false)
			om.FromSpace[i].SetForwardingPtr(nil)
		}
	}

	// Start with the root set
	toPtr := 0

	// Add globals to the root set
	for _, obj := range vm.GetGlobals() {
		if obj != nil {
			copied := om.copyObject(obj, &toPtr)
			// Update the global reference
			*obj = *copied
		}
	}

	// Add the current context and its chain to the root set
	context := vm.GetCurrentContext().(ExecutionContext)
	for context != nil {
		// Copy the context's method
		if context.GetMethod() != nil {
			// In a real implementation, this would update the method reference in the context
			_ = om.copyObject(context.GetMethod(), &toPtr)
		}

		// Copy the context's receiver
		if context.GetReceiver() != nil {
			// In a real implementation, this would update the receiver reference in the context
			_ = om.copyObject(context.GetReceiver(), &toPtr)
		}

		// Copy the context's arguments
		for _, arg := range context.GetArguments() {
			if arg != nil {
				// In a real implementation, this would update the argument reference in the context
				_ = om.copyObject(arg, &toPtr)
			}
		}

		// Copy the context's temporary variables
		for _, obj := range context.GetTempVars() {
			if obj != nil {
				// In a real implementation, this would update the temp var reference in the context
				_ = om.copyObject(obj, &toPtr)
			}
		}

		// Copy the context's stack
		for i := 0; i < context.GetStackPointer(); i++ {
			stack := context.GetStack()
			if stack[i] != nil {
				// In a real implementation, this would update the stack reference in the context
				_ = om.copyObject(stack[i], &toPtr)
			}
		}

		// Move to the sender context
		context = context.GetSender().(ExecutionContext)
	}

	// Copy special objects
	// Note: NilObject, TrueObject, and FalseObject are now immediate values, so we don't need to copy them
	if vm.GetObjectClass() != nil {
		// In a real implementation, this would update the object class reference in the VM
		_ = om.copyObject((*Object)(unsafe.Pointer(vm.GetObjectClass())), &toPtr)
	}

	// Scan the to-space for references
	for i := 0; i < toPtr; i++ {
		obj := om.ToSpace[i]
		if obj == nil {
			continue
		}

		// Update references in the object
		om.updateReferences(obj, &toPtr)
	}

	// Swap the spaces
	om.FromSpace, om.ToSpace = om.ToSpace, om.FromSpace
	om.AllocPtr = toPtr

	// Grow the spaces if needed
	if toPtr > om.SpaceSize*70/100 { // If we're using more than 70% after GC
		om.growSpaces()
	}
}

// copyObject copies an object to the to-space
func (om *ObjectMemory) copyObject(obj ObjectInterface, toPtr *int) *Object {
	// Check if it's an immediate value
	if IsImmediate(obj) {
		// Immediate values don't need to be copied
		return obj.(*Object)
	}

	// Check if the object has already been moved
	if obj.Moved() {
		return obj.ForwardingPtr()
	}

	// Copy the object to the to-space
	om.ToSpace[*toPtr] = obj.(*Object)
	obj.SetMoved(true)
	obj.SetForwardingPtr(om.ToSpace[*toPtr])
	*toPtr++

	return obj.ForwardingPtr()
}

// updateReferences updates references in an object
func (om *ObjectMemory) updateReferences(obj *Object, toPtr *int) {
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
		array := (*Array)(unsafe.Pointer(obj))
		for i, elem := range array.Elements {
			if elem != nil {
				array.Elements[i] = om.copyObject(elem, toPtr)
			}
		}

	case OBJ_DICTIONARY:
		// Update dictionary entries
		dict := (*Dictionary)(unsafe.Pointer(obj))
		entries := dict.GetEntries()
		for key, value := range entries {
			if value != nil {
				dict.SetEntry(key, om.copyObject(value, toPtr))
			}
		}

	case OBJ_INSTANCE:
		// Update instance variables
		instanceVars := obj.InstanceVars()
		for i, value := range instanceVars {
			if value != nil {
				instanceVars[i] = om.copyObject(value, toPtr)
			}
		}

		// Update class reference
		if obj.Class() != nil {
			obj.SetClass(om.copyObject(obj.Class(), toPtr))
		}

	case OBJ_CLASS:
		// Update instance variables (which includes the method dictionary)
		instanceVars := obj.InstanceVars()
		for i, value := range instanceVars {
			if value != nil {
				instanceVars[i] = om.copyObject(value, toPtr)
			}
		}

		// Update superclass reference
		class := (*Class)(unsafe.Pointer(obj))
		if class.SuperClass != nil {
			class.SuperClass = om.copyObject(class.SuperClass, toPtr)
		}

	case OBJ_METHOD:
		// Update method literals
		method := (*Method)(unsafe.Pointer(obj))
		for i, lit := range method.Literals {
			if lit != nil {
				method.Literals[i] = om.copyObject(lit, toPtr)
			}
		}

		// Update method selector
		if method.Selector != nil {
			method.Selector = om.copyObject(method.Selector, toPtr)
		}

		// Update method class
		if method.MethodClass != nil {
			method.MethodClass = (*Class)(unsafe.Pointer(om.copyObject((*Object)(unsafe.Pointer(method.MethodClass)), toPtr)))
		}

	case OBJ_BLOCK:
		// Update block literals
		block := (*Block)(unsafe.Pointer(obj))
		for i, lit := range block.Literals {
			if lit != nil {
				block.Literals[i] = om.copyObject(lit, toPtr)
			}
		}
	}
}

// growSpaces grows the from-space and to-space
func (om *ObjectMemory) growSpaces() {
	newSize := om.SpaceSize * 2

	// Create new spaces
	newFromSpace := make([]*Object, newSize)
	newToSpace := make([]*Object, newSize)

	// Copy objects to the new from-space
	copy(newFromSpace, om.FromSpace)

	// Update the spaces
	om.FromSpace = newFromSpace
	om.ToSpace = newToSpace
	om.SpaceSize = newSize
	om.GCThreshold = newSize * 80 / 100 // 80% threshold
}
