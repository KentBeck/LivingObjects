package main

// ObjectMemory manages the Smalltalk object memory with stop & copy garbage collection
type ObjectMemory struct {
	FromSpace   []*Object
	ToSpace     []*Object
	AllocPtr    int
	SpaceSize   int
	GCThreshold int
	GCCount     int
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
func (om *ObjectMemory) Collect(vm *VM) {
	om.GCCount++

	// Reset the to-space
	for i := range om.ToSpace {
		om.ToSpace[i] = nil
	}

	// Reset forwarding pointers
	for i := 0; i < om.AllocPtr; i++ {
		if om.FromSpace[i] != nil {
			om.FromSpace[i].Moved = false
			om.FromSpace[i].ForwardingPtr = nil
		}
	}

	// Start with the root set
	toPtr := 0

	// Add globals to the root set
	for _, obj := range vm.Globals {
		if obj != nil {
			copied := om.copyObject(obj, &toPtr)
			// Update the global reference
			*obj = *copied
		}
	}

	// Add the current context and its chain to the root set
	context := vm.CurrentContext
	for context != nil {
		// Copy the context's method
		if context.Method != nil {
			copied := om.copyObject(context.Method, &toPtr)
			context.Method = copied
		}

		// Copy the context's receiver
		if context.Receiver != nil {
			copied := om.copyObject(context.Receiver, &toPtr)
			context.Receiver = copied
		}

		// Copy the context's arguments
		for i, arg := range context.Arguments {
			if arg != nil {
				copied := om.copyObject(arg, &toPtr)
				context.Arguments[i] = copied
			}
		}

		// Copy the context's temporary variables
		for i, obj := range context.TempVars {
			if obj != nil {
				copied := om.copyObject(obj, &toPtr)
				context.TempVars[i] = copied
			}
		}

		// Copy the context's stack
		for i := 0; i < context.StackPointer; i++ {
			if context.Stack[i] != nil {
				copied := om.copyObject(context.Stack[i], &toPtr)
				context.Stack[i] = copied
			}
		}

		// Move to the sender context
		context = context.Sender
	}

	// Copy special objects
	// Note: NilObject, TrueObject, and FalseObject are now immediate values, so we don't need to copy them
	// if vm.NilObject != nil {
	// 	vm.NilObject = om.copyObject(vm.NilObject, &toPtr)
	// }
	// if vm.TrueObject != nil {
	// 	vm.TrueObject = om.copyObject(vm.TrueObject, &toPtr)
	// }
	// if vm.FalseObject != nil {
	// 	vm.FalseObject = om.copyObject(vm.FalseObject, &toPtr)
	// }
	if vm.ObjectClass != nil {
		vm.ObjectClass = om.copyObject(vm.ObjectClass, &toPtr)
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
func (om *ObjectMemory) copyObject(obj *Object, toPtr *int) *Object {
	// Check if it's an immediate value
	if IsImmediate(obj) {
		// Immediate values don't need to be copied
		return obj
	}

	// Check if the object has already been moved
	if obj.Moved {
		return obj.ForwardingPtr
	}

	// Copy the object to the to-space
	om.ToSpace[*toPtr] = obj
	obj.Moved = true
	obj.ForwardingPtr = om.ToSpace[*toPtr]
	*toPtr++

	return obj.ForwardingPtr
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
		for i, elem := range obj.Elements {
			if elem != nil {
				obj.Elements[i] = om.copyObject(elem, toPtr)
			}
		}

	case OBJ_DICTIONARY:
		// Update dictionary entries
		for key, value := range obj.Entries {
			if value != nil {
				obj.Entries[key] = om.copyObject(value, toPtr)
			}
		}

	case OBJ_INSTANCE:
		// Update instance variables
		for i, value := range obj.InstanceVars {
			if value != nil {
				obj.InstanceVars[i] = om.copyObject(value, toPtr)
			}
		}

		// Update superclass reference
		if obj.SuperClass != nil {
			obj.SuperClass = om.copyObject(obj.SuperClass, toPtr)
		}

		// Update class reference
		if obj.Class() != nil {
			obj.SetClass(om.copyObject(obj.Class(), toPtr))
		}

	case OBJ_CLASS:
		// Update instance variables (which includes the method dictionary)
		for i, value := range obj.InstanceVars {
			if value != nil {
				obj.InstanceVars[i] = om.copyObject(value, toPtr)
			}
		}

		// Update superclass reference
		if obj.SuperClass != nil {
			obj.SuperClass = om.copyObject(obj.SuperClass, toPtr)
		}

	case OBJ_METHOD:
		// Update method literals
		for i, lit := range obj.Method.Literals {
			if lit != nil {
				obj.Method.Literals[i] = om.copyObject(lit, toPtr)
			}
		}

		// Update method selector
		if obj.Method.Selector != nil {
			obj.Method.Selector = om.copyObject(obj.Method.Selector, toPtr)
		}

		// Update method class
		if obj.Method.MethodClass != nil {
			obj.Method.MethodClass = om.copyObject(obj.Method.MethodClass, toPtr)
		}

	case OBJ_BLOCK:
		// Update block literals
		for i, lit := range obj.Literals {
			if lit != nil {
				obj.Literals[i] = om.copyObject(lit, toPtr)
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
