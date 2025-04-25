package main

import (
	"fmt"
)

// VM represents the Smalltalk virtual machine
type VM struct {
	Globals        map[string]*Object
	CurrentContext *Context
	ObjectMemory   *ObjectMemory

	// Special objects
	NilObject    *Object
	NilClass     *Object
	TrueObject   *Object
	TrueClass    *Object
	FalseObject  *Object
	FalseClass   *Object
	ObjectClass  *Object
	IntegerClass *Object
	FloatClass   *Object
}

// NewVM creates a new virtual machine
func NewVM() *VM {
	vm := &VM{
		Globals:      make(map[string]*Object),
		ObjectMemory: NewObjectMemory(),
	}

	// Initialize special objects
	vm.ObjectClass = vm.NewObjectClass()
	vm.NilClass = NewClass("UndefinedObject", vm.ObjectClass)
	vm.NilObject = MakeNilImmediate()
	vm.TrueClass = NewClass("True", vm.ObjectClass)
	vm.TrueObject = MakeTrueImmediate()
	vm.FalseClass = NewClass("False", vm.ObjectClass)
	vm.FalseObject = MakeFalseImmediate()
	vm.IntegerClass = vm.NewIntegerClass()
	vm.FloatClass = vm.NewFloatClass()

	return vm
}

func (vm *VM) NewObjectClass() *Object {
	result := NewClass("Object", nil) // patch this up later. then even later when we have real images all this initialization can go away
	builder := NewMethodBuilder()
	builder.AddPrimitive(5)
	builder.Install(result, ObjectToSymbol(NewSymbol("basicClass")))
	return result
}

func (vm *VM) NewIntegerClass() *Object {
	result := NewClass("Integer", vm.ObjectClass)
	// Add primitive methods to the Integer class
	// + method
	plusSelector := NewSymbol("+")
	plusMethod := NewMethod(plusSelector, result)
	plusMethod.Method.IsPrimitive = true
	plusMethod.Method.PrimitiveIndex = 1 // Addition
	integerMethodDict := result.GetMethodDict()
	integerMethodDict.Entries[GetSymbolValue(plusSelector)] = plusMethod

	// - method (subtraction)
	minusSelector := NewSymbol("-")
	minusMethod := NewMethod(minusSelector, result)
	minusMethod.Method.IsPrimitive = true
	minusMethod.Method.PrimitiveIndex = 4 // Subtraction (new primitive)
	integerMethodDict.Entries[GetSymbolValue(minusSelector)] = minusMethod

	// * method
	timesSelector := NewSymbol("*")
	timesMethod := NewMethod(timesSelector, result)
	timesMethod.Method.IsPrimitive = true
	timesMethod.Method.PrimitiveIndex = 2 // Multiplication
	integerMethodDict.Entries[GetSymbolValue(timesSelector)] = timesMethod

	// = method
	equalsSelector := NewSymbol("=")
	equalsMethod := NewMethod(equalsSelector, result)
	equalsMethod.Method.IsPrimitive = true
	equalsMethod.Method.PrimitiveIndex = 3 // Equality
	integerMethodDict.Entries[GetSymbolValue(equalsSelector)] = equalsMethod

	// < method
	lessSelector := NewSymbol("<")
	lessMethod := NewMethod(lessSelector, result)
	lessMethod.Method.IsPrimitive = true
	lessMethod.Method.PrimitiveIndex = 6 // Less than
	integerMethodDict.Entries[GetSymbolValue(lessSelector)] = lessMethod

	// > method
	greaterSelector := NewSymbol(">")
	greaterMethod := NewMethod(greaterSelector, result)
	greaterMethod.Method.IsPrimitive = true
	greaterMethod.Method.PrimitiveIndex = 7 // Greater than
	integerMethodDict.Entries[GetSymbolValue(greaterSelector)] = greaterMethod

	return result
}

func (vm *VM) NewFloatClass() *Object {
	result := NewClass("Float", vm.ObjectClass) // patch this up later. then even later when we have real images all this initialization can go away

	// Get the method dictionary for the Float class
	floatMethodDict := result.GetMethodDict()

	// + method (addition)
	plusSelector := NewSymbol("+")
	plusMethod := NewMethod(plusSelector, result)
	plusMethod.Method.IsPrimitive = true
	plusMethod.Method.PrimitiveIndex = 10 // Float addition
	floatMethodDict.Entries[GetSymbolValue(plusSelector)] = plusMethod

	// - method (subtraction)
	minusSelector := NewSymbol("-")
	minusMethod := NewMethod(minusSelector, result)
	minusMethod.Method.IsPrimitive = true
	minusMethod.Method.PrimitiveIndex = 11 // Float subtraction
	floatMethodDict.Entries[GetSymbolValue(minusSelector)] = minusMethod

	// * method (multiplication)
	timesSelector := NewSymbol("*")
	timesMethod := NewMethod(timesSelector, result)
	timesMethod.Method.IsPrimitive = true
	timesMethod.Method.PrimitiveIndex = 12 // Float multiplication
	floatMethodDict.Entries[GetSymbolValue(timesSelector)] = timesMethod

	// / method (division)
	divideSelector := NewSymbol("/")
	divideMethod := NewMethod(divideSelector, result)
	divideMethod.Method.IsPrimitive = true
	divideMethod.Method.PrimitiveIndex = 13 // Float division
	floatMethodDict.Entries[GetSymbolValue(divideSelector)] = divideMethod

	// = method (equality)
	equalsSelector := NewSymbol("=")
	equalsMethod := NewMethod(equalsSelector, result)
	equalsMethod.Method.IsPrimitive = true
	equalsMethod.Method.PrimitiveIndex = 14 // Float equality
	floatMethodDict.Entries[GetSymbolValue(equalsSelector)] = equalsMethod

	// < method (less than)
	lessSelector := NewSymbol("<")
	lessMethod := NewMethod(lessSelector, result)
	lessMethod.Method.IsPrimitive = true
	lessMethod.Method.PrimitiveIndex = 15 // Float less than
	floatMethodDict.Entries[GetSymbolValue(lessSelector)] = lessMethod

	// > method (greater than)
	greaterSelector := NewSymbol(">")
	greaterMethod := NewMethod(greaterSelector, result)
	greaterMethod.Method.IsPrimitive = true
	greaterMethod.Method.PrimitiveIndex = 16 // Float greater than
	floatMethodDict.Entries[GetSymbolValue(greaterSelector)] = greaterMethod

	return result
}

// NewInteger creates a new integer object
// This returns an immediate value for integers
func (vm *VM) NewInteger(value int64) *Object {
	// Check if the value fits in 62 bits
	if value <= 0x1FFFFFFFFFFFFFFF && value >= -0x2000000000000000 {
		// Use immediate integer
		return MakeIntegerImmediate(value)
	}

	// Panic for large values that don't fit in 62 bits
	panic("Integer value too large for immediate representation")
}

func (vm *VM) NewFloat(value float64) *Object {
	return MakeFloatImmediate(value)
}

// LoadImage loads a Smalltalk image from a file
func (vm *VM) LoadImage(path string) error {
	vm.Globals["Object"] = vm.ObjectClass

	return nil
}

// Execute executes the current context
func (vm *VM) Execute() (*Object, error) {
	var finalResult *Object

	for vm.CurrentContext != nil {
		// Execute the current context
		result, err := vm.ExecuteContext(vm.CurrentContext)
		if err != nil {
			return nil, err
		}

		// Save the result if this is the top-level context
		if vm.CurrentContext.Sender == nil {
			finalResult = result
		}

		// Move to the sender context
		vm.CurrentContext = vm.CurrentContext.Sender

		// If we have a sender, push the result onto its stack
		if vm.CurrentContext != nil {
			vm.CurrentContext.Push(result)
		}
	}

	return finalResult, nil
}

// ExecuteContext executes a single context until it returns
func (vm *VM) ExecuteContext(context *Context) (*Object, error) {
	// Execute the context

	for {
		// Check if we've reached the end of the method
		if context.PC >= len(context.Method.Method.Bytecodes) {
			// Reached end of bytecode array

			// If we've reached the end of the method, return the top of the stack
			// This handles the case where we jump to the end of the bytecode array
			if context.StackPointer > 0 {
				returnValue := context.Pop()
				return returnValue, nil
			}
			return vm.NilObject, nil
		}

		// Get the current bytecode
		bytecode := context.Method.Method.Bytecodes[context.PC]

		// Get the instruction size
		size := InstructionSize(bytecode)

		// Execute the bytecode
		var err error
		var skipIncrement bool
		var returnValue *Object

		switch bytecode {
		case PUSH_LITERAL:
			err = vm.ExecutePushLiteral(context)

		case PUSH_INSTANCE_VARIABLE:
			err = vm.ExecutePushInstanceVariable(context)

		case PUSH_TEMPORARY_VARIABLE:
			err = vm.ExecutePushTemporaryVariable(context)

		case PUSH_SELF:
			err = vm.ExecutePushSelf(context)

		case STORE_INSTANCE_VARIABLE:
			err = vm.ExecuteStoreInstanceVariable(context)

		case STORE_TEMPORARY_VARIABLE:
			err = vm.ExecuteStoreTemporaryVariable(context)

		case SEND_MESSAGE:
			returnValue, err = vm.ExecuteSendMessage(context)
			if err == nil {
				if returnValue != nil {
					// We got a result from a primitive method
					// Continue execution in the current context
					context.PC += size
					continue
				} else {
					// A nil return value with no error means we've started a new context
					return vm.NilObject, nil
				}
			}

		case RETURN_STACK_TOP:
			returnValue, err = vm.ExecuteReturnStackTop(context)
			if err == nil {
				return returnValue, nil
			}

		case JUMP:
			skipIncrement, err = vm.ExecuteJump(context)
			if err == nil && skipIncrement {
				continue
			}

		case JUMP_IF_TRUE:
			skipIncrement, err = vm.ExecuteJumpIfTrue(context)
			if err == nil && skipIncrement {
				continue
			}

		case JUMP_IF_FALSE:
			skipIncrement, err = vm.ExecuteJumpIfFalse(context)
			if err == nil && skipIncrement {
				continue
			}

		case POP:
			err = vm.ExecutePop(context)

		case DUPLICATE:
			err = vm.ExecuteDuplicate(context)

		default:
			return nil, fmt.Errorf("unknown bytecode: %d", bytecode)
		}

		// Check for errors
		if err != nil {
			return nil, err
		}

		// Increment the PC
		context.PC += size
	}
}

// lookupMethod looks up a method in a class hierarchy
func (vm *VM) lookupMethod(receiver *Object, selector *Object) *Object {
	// Check for nil receiver or selector
	if receiver == nil {
		panic("lookupMethod: nil receiver\n")
	}

	if selector == nil {
		panic("lookupMethod: nil  selector\n")
	}

	class := vm.GetClass(receiver)
	if class == nil {
		panic("lookupMethod: nil class\n")
	}

	// Look up the method in the class hierarchy
	for class != nil {
		// Check if the class has a method dictionary
		methodDict := class.GetMethodDict()
		if methodDict != nil && methodDict.Type == OBJ_DICTIONARY && methodDict.Entries != nil {
			// Check if the method dictionary has the selector
			if method, ok := methodDict.Entries[GetSymbolValue(selector)]; ok {
				return method
			}
		} else {
			panic("method dictionary is nil or not a dictionary")
		}

		// Move up the class hierarchy
		class = class.SuperClass
	}

	// Method not found
	return nil
}
