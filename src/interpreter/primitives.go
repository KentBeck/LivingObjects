package main

// primitives.go contains the implementation of primitive methods for the Smalltalk interpreter

// executePrimitive executes a primitive method
func (vm *VM) executePrimitive(receiver *Object, selector *Object, args []*Object, method *Object) *Object {
	if receiver == nil {
		panic("executePrimitive: nil receiver\n")
	}
	if selector == nil {
		panic("executePrimitive: nil selector\n")
	}
	if method == nil {
		panic("executePrimitive: nil method\n")
	}
	if method.Type != OBJ_METHOD {
		panic("executePrimitive: method is not a method\n")
	}
	if !method.Method.IsPrimitive {
		return nil
	}

	// Execute the primitive based on its index
	switch method.Method.PrimitiveIndex {
	case 1: // Addition
		// Handle immediate integers
		if IsIntegerImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetIntegerImmediate(receiver)
			val2 := GetIntegerImmediate(args[0])
			result := val1 + val2
			return vm.NewInteger(result)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type == OBJ_INTEGER || (len(args) > 0 && args[0].Type == OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
		// Handle integer + float
		if IsIntegerImmediate(receiver) && len(args) == 1 && IsFloatImmediate(args[0]) {
			val1 := float64(GetIntegerImmediate(receiver))
			val2 := GetFloatImmediate(args[0])
			result := val1 + val2
			return vm.NewFloat(result)
		}
	case 2: // Multiplication
		// Handle immediate integers
		if IsIntegerImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetIntegerImmediate(receiver)
			val2 := GetIntegerImmediate(args[0])
			result := val1 * val2
			return vm.NewInteger(result)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type == OBJ_INTEGER || (len(args) > 0 && args[0].Type == OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 3: // Equality
		// Handle immediate integers
		if IsIntegerImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetIntegerImmediate(receiver)
			val2 := GetIntegerImmediate(args[0])
			result := val1 == val2
			return NewBoolean(result)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type == OBJ_INTEGER || (len(args) > 0 && args[0].Type == OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 4: // Subtraction
		// Handle immediate integers
		if IsIntegerImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetIntegerImmediate(receiver)
			val2 := GetIntegerImmediate(args[0])
			result := val1 - val2
			return vm.NewInteger(result)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type == OBJ_INTEGER || (len(args) > 0 && args[0].Type == OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 5: // basicClass - return the class of the receiver
		if len(args) == 0 {
			class := vm.GetClass(receiver)
			return class
		}
	case 6: // Less than
		// Handle immediate integers
		if IsIntegerImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetIntegerImmediate(receiver)
			val2 := GetIntegerImmediate(args[0])
			result := val1 < val2
			return NewBoolean(result)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type == OBJ_INTEGER || (len(args) > 0 && args[0].Type == OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 7: // Greater than
		// Handle immediate integers
		if IsIntegerImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetIntegerImmediate(receiver)
			val2 := GetIntegerImmediate(args[0])
			result := val1 > val2
			return NewBoolean(result)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type == OBJ_INTEGER || (len(args) > 0 && args[0].Type == OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 10: // Float addition
		// Handle float + float
		if IsFloatImmediate(receiver) && len(args) == 1 && IsFloatImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := GetFloatImmediate(args[0])
			result := val1 + val2
			return vm.NewFloat(result)
		}
		// Handle float + integer
		if IsFloatImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := float64(GetIntegerImmediate(args[0]))
			result := val1 + val2
			return vm.NewFloat(result)
		}
	case 11: // Float subtraction
		// Handle float - float
		if IsFloatImmediate(receiver) && len(args) == 1 && IsFloatImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := GetFloatImmediate(args[0])
			result := val1 - val2
			return vm.NewFloat(result)
		}
		// Handle float - integer
		if IsFloatImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := float64(GetIntegerImmediate(args[0]))
			result := val1 - val2
			return vm.NewFloat(result)
		}
	case 12: // Float multiplication
		// Handle float * float
		if IsFloatImmediate(receiver) && len(args) == 1 && IsFloatImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := GetFloatImmediate(args[0])
			result := val1 * val2
			return vm.NewFloat(result)
		}
		// Handle float * integer
		if IsFloatImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := float64(GetIntegerImmediate(args[0]))
			result := val1 * val2
			return vm.NewFloat(result)
		}
	case 13: // Float division
		// Handle float / float
		if IsFloatImmediate(receiver) && len(args) == 1 && IsFloatImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := GetFloatImmediate(args[0])
			result := val1 / val2
			return vm.NewFloat(result)
		}
		// Handle float / integer
		if IsFloatImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := float64(GetIntegerImmediate(args[0]))
			result := val1 / val2
			return vm.NewFloat(result)
		}
	case 14: // Float equality
		// Handle float = float
		if IsFloatImmediate(receiver) && len(args) == 1 && IsFloatImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := GetFloatImmediate(args[0])
			result := val1 == val2
			return NewBoolean(result)
		}
		// Handle float = integer
		if IsFloatImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := float64(GetIntegerImmediate(args[0]))
			result := val1 == val2
			return NewBoolean(result)
		}
	case 15: // Float less than
		// Handle float < float
		if IsFloatImmediate(receiver) && len(args) == 1 && IsFloatImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := GetFloatImmediate(args[0])
			result := val1 < val2
			return NewBoolean(result)
		}
		// Handle float < integer
		if IsFloatImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := float64(GetIntegerImmediate(args[0]))
			result := val1 < val2
			return NewBoolean(result)
		}
	case 16: // Float greater than
		// Handle float > float
		if IsFloatImmediate(receiver) && len(args) == 1 && IsFloatImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := GetFloatImmediate(args[0])
			result := val1 > val2
			return NewBoolean(result)
		}
		// Handle float > integer
		if IsFloatImmediate(receiver) && len(args) == 1 && IsIntegerImmediate(args[0]) {
			val1 := GetFloatImmediate(receiver)
			val2 := float64(GetIntegerImmediate(args[0]))
			result := val1 > val2
			return NewBoolean(result)
		}
	default:
		panic("executePrimitive: unknown primitive index\n")
	}
	return nil // Fall through to method
}
