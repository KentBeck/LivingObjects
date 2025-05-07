package vm

import (
	"fmt"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// primitives.go contains the implementation of primitive methods for the Smalltalk interpreter

// executePrimitive executes a primitive method
func (vm *VM) executePrimitive(receiver *core.Object, selector *core.Object, args []*core.Object, method *core.Object) *core.Object {
	if receiver == nil {
		panic("executePrimitive: nil receiver\n")
	}
	if selector == nil {
		panic("executePrimitive: nil selector\n")
	}
	if method == nil {
		panic("executePrimitive: nil method\n")
	}
	if method.Type() != core.OBJ_METHOD {
		return nil
	}
	methodObj := classes.ObjectToMethod(method)
	if !methodObj.IsPrimitiveMethod() {
		return nil
	}

	// Execute the primitive based on its index
	switch methodObj.GetPrimitiveIndex() {
	case 1: // Addition
		// Handle immediate integers
		if core.IsIntegerImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetIntegerImmediate(receiver)
			val2 := core.GetIntegerImmediate(args[0])
			result := val1 + val2
			return vm.NewInteger(result)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type() == core.OBJ_INTEGER || (len(args) > 0 && args[0].Type() == core.OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
		// Handle integer + float
		if core.IsIntegerImmediate(receiver) && len(args) == 1 && core.IsFloatImmediate(args[0]) {
			val1 := float64(core.GetIntegerImmediate(receiver))
			val2 := core.GetFloatImmediate(args[0])
			result := val1 + val2
			return vm.NewFloat(result)
		}
	case 2: // Multiplication
		// Handle immediate integers
		if core.IsIntegerImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetIntegerImmediate(receiver)
			val2 := core.GetIntegerImmediate(args[0])
			result := val1 * val2
			return vm.NewInteger(result)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type() == core.OBJ_INTEGER || (len(args) > 0 && args[0].Type() == core.OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 3: // Equality
		// Handle immediate integers
		if core.IsIntegerImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetIntegerImmediate(receiver)
			val2 := core.GetIntegerImmediate(args[0])
			result := val1 == val2
			return core.NewBoolean(result).(*core.Object)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type() == core.OBJ_INTEGER || (len(args) > 0 && args[0].Type() == core.OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 4: // Subtraction
		// Handle immediate integers
		if core.IsIntegerImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetIntegerImmediate(receiver)
			val2 := core.GetIntegerImmediate(args[0])
			result := val1 - val2
			return vm.NewInteger(result)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type() == core.OBJ_INTEGER || (len(args) > 0 && args[0].Type() == core.OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 5: // basicClass - return the class of the receiver
		class := vm.GetClass(receiver)
		return classes.ClassToObject(class)
	case 6: // Less than
		// Handle immediate integers
		if core.IsIntegerImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetIntegerImmediate(receiver)
			val2 := core.GetIntegerImmediate(args[0])
			result := val1 < val2
			return core.NewBoolean(result).(*core.Object)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type() == core.OBJ_INTEGER || (len(args) > 0 && args[0].Type() == core.OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 7: // Greater than
		// Handle immediate integers
		if core.IsIntegerImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetIntegerImmediate(receiver)
			val2 := core.GetIntegerImmediate(args[0])
			result := val1 > val2
			return core.NewBoolean(result).(*core.Object)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type() == core.OBJ_INTEGER || (len(args) > 0 && args[0].Type() == core.OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 10: // Float addition
		// Handle float + float
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsFloatImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := core.GetFloatImmediate(args[0])
			result := val1 + val2
			return vm.NewFloat(result)
		}
		// Handle float + integer
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := float64(core.GetIntegerImmediate(args[0]))
			result := val1 + val2
			return vm.NewFloat(result)
		}
	case 11: // Float subtraction
		// Handle float - float
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsFloatImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := core.GetFloatImmediate(args[0])
			result := val1 - val2
			return vm.NewFloat(result)
		}
		// Handle float - integer
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := float64(core.GetIntegerImmediate(args[0]))
			result := val1 - val2
			return vm.NewFloat(result)
		}
	case 12: // Float multiplication
		// Handle float * float
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsFloatImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := core.GetFloatImmediate(args[0])
			result := val1 * val2
			return vm.NewFloat(result)
		}
		// Handle float * integer
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := float64(core.GetIntegerImmediate(args[0]))
			result := val1 * val2
			return vm.NewFloat(result)
		}
	case 13: // Float division
		// Handle float / float
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsFloatImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := core.GetFloatImmediate(args[0])
			result := val1 / val2
			return vm.NewFloat(result)
		}
		// Handle float / integer
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := float64(core.GetIntegerImmediate(args[0]))
			result := val1 / val2
			return vm.NewFloat(result)
		}
	case 14: // Float equality
		// Handle float = float
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsFloatImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := core.GetFloatImmediate(args[0])
			result := val1 == val2
			return core.NewBoolean(result).(*core.Object)
		}
		// Handle float = integer
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := float64(core.GetIntegerImmediate(args[0]))
			result := val1 == val2
			return core.NewBoolean(result).(*core.Object)
		}
	case 15: // Float less than
		// Handle float < float
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsFloatImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := core.GetFloatImmediate(args[0])
			result := val1 < val2
			return core.NewBoolean(result).(*core.Object)
		}
		// Handle float < integer
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := float64(core.GetIntegerImmediate(args[0]))
			result := val1 < val2
			return core.NewBoolean(result).(*core.Object)
		}
	case 16: // Float greater than
		// Handle float > float
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsFloatImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := core.GetFloatImmediate(args[0])
			result := val1 > val2
			return core.NewBoolean(result).(*core.Object)
		}
		// Handle float > integer
		if core.IsFloatImmediate(receiver) && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			val1 := core.GetFloatImmediate(receiver)
			val2 := float64(core.GetIntegerImmediate(args[0]))
			result := val1 > val2
			return core.NewBoolean(result).(*core.Object)
		}
	case 20: // Block new - create a new block instance
		if receiver.Type() == core.OBJ_CLASS && receiver == classes.ClassToObject(vm.BlockClass) {
			// Create a new block instance
			blockInstance := classes.NewBlock(vm.CurrentContext)
			blockInstance.SetClass(classes.ClassToObject(vm.BlockClass))
			return blockInstance
		}
	case 21: // Block value - execute a block with no arguments
		if receiver.Type() == core.OBJ_BLOCK {
			// Get the block
			block := classes.ObjectToBlock(receiver)

			// Create a method object for the block's bytecodes
			method := &classes.Method{
				Object: core.Object{
					TypeField: core.OBJ_METHOD,
				},
				Bytecodes: block.GetBytecodes(),
				Literals:  block.GetLiterals(),
			}
			methodObj := classes.MethodToObject(method)

			// Create a new context for the block execution
			blockContext := NewContext(methodObj, receiver, []*core.Object{}, block.GetOuterContext().(*Context))

			// Execute the block's bytecodes
			result, err := vm.ExecuteContext(blockContext)
			if err != nil {
				panic(fmt.Sprintf("Error executing block: %v", err))
			}
			return result.(*core.Object)
		}
	case 22: // Block value: - execute a block with one argument
		if receiver.Type() == core.OBJ_BLOCK && len(args) == 1 {
			// Get the block
			block := classes.ObjectToBlock(receiver)

			// Create a method object for the block's bytecodes
			method := &classes.Method{
				Object: core.Object{
					TypeField: core.OBJ_METHOD,
				},
				Bytecodes: block.GetBytecodes(),
				Literals:  block.GetLiterals(),
			}
			methodObj := classes.MethodToObject(method)

			// Create a new context for the block execution
			blockContext := NewContext(methodObj, receiver, args, block.GetOuterContext().(*Context))

			// Execute the block's bytecodes
			result, err := vm.ExecuteContext(blockContext)
			if err != nil {
				panic(fmt.Sprintf("Error executing block: %v", err))
			}
			return result.(*core.Object)
		}
	case 30: // String concatenation (,)
		if receiver.Type() == core.OBJ_STRING && len(args) == 1 && args[0].Type() == core.OBJ_STRING {
			// Get the string values
			str1 := classes.ObjectToString(receiver)
			str2 := classes.ObjectToString(args[0])

			// Concatenate the strings
			result := str1.Concat(str2)

			// Return the result
			return classes.StringToObject(result)
		}
	default:
		panic("executePrimitive: unknown primitive index\n")
	}
	return nil // Fall through to method
}
