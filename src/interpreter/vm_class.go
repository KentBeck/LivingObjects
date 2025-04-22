package main

// GetClass returns the class of an object
// This is the single function that should be used to get the class of an object
func (vm *VM) GetClass(obj *Object) *Object {
	if obj == nil {
		panic("GetClass: nil object")
	}

	// If the object is a class, return itself
	if obj.Type == OBJ_CLASS {
		return obj
	}

	// Special case for nil object
	if obj.Type == OBJ_NIL {
		return vm.NilClass
	}

	// Otherwise, return the class field
	if obj.Class == nil {
		panic("GetClass: object has nil class")
	}

	return obj.Class
}
