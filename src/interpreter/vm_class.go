package main

// GetClass returns the class of an object
// This is the single function that should be used to get the class of an object
func (vm *VM) GetClass(obj *Object) *Object {
	if obj == nil {
		panic("GetClass: nil object")
	}

	// Check if it's an immediate value
	if IsImmediate(obj) {
		// Handle immediate nil
		if IsNilImmediate(obj) {
			return vm.NilClass
		}
		// Handle immediate true
		if IsTrueImmediate(obj) {
			return vm.TrueClass
		}
		// Handle immediate false
		if IsFalseImmediate(obj) {
			return vm.FalseClass
		}
		// Handle immediate integer
		if IsIntegerImmediate(obj) {
			return vm.IntegerClass
		}
		// Other immediate types will be added later
		panic("GetClass: unknown immediate type")
	}

	// If it's a regular object, proceed as before

	// If the object is a class, return itself
	if obj.Type == OBJ_CLASS {
		return obj
	}

	// Special case for nil object (legacy non-immediate nil)
	if obj.Type == OBJ_NIL {
		return vm.NilClass
	}

	// Otherwise, return the class field
	if obj.Class == nil {
		panic("GetClass: object has nil class")
	}

	return obj.Class
}
