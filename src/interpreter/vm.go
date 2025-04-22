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
}

// NewVM creates a new virtual machine
func NewVM() *VM {
	vm := &VM{
		Globals:      make(map[string]*Object),
		ObjectMemory: NewObjectMemory(),
	}

	// Initialize special objects
	vm.ObjectClass = vm.NewObjectClass()
	vm.NilClass = NewClass("UndefinedObject", vm.ObjectClass) // ci
	vm.NilObject = NewNil()
	vm.TrueClass = NewClass("True", vm.ObjectClass)
	vm.TrueObject = &Object{Type: OBJ_BOOLEAN, BooleanValue: true, Class: vm.TrueClass}
	vm.FalseClass = NewClass("False", vm.ObjectClass)
	vm.FalseObject = &Object{Type: OBJ_BOOLEAN, BooleanValue: false, Class: vm.FalseClass}
	vm.IntegerClass = vm.NewIntegerClass()

	return vm
}

func (vm *VM) NewObjectClass() *Object {
	result := NewClass("Object", nil) // patch this up later. then even later when we have real images all this initialization can go away

	// Add basicClass method to Object class
	objectMethodDict := result.GetMethodDict()
	basicClassSelector := NewSymbol("basicClass")
	basicClassMethod := NewMethod(basicClassSelector, result)
	basicClassMethod.Method.IsPrimitive = true
	basicClassMethod.Method.PrimitiveIndex = 5 // basicClass primitive
	objectMethodDict.Entries[basicClassSelector.SymbolValue] = basicClassMethod

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
	integerMethodDict.Entries[plusSelector.SymbolValue] = plusMethod

	// - method (subtraction)
	minusSelector := NewSymbol("-")
	minusMethod := NewMethod(minusSelector, result)
	minusMethod.Method.IsPrimitive = true
	minusMethod.Method.PrimitiveIndex = 4 // Subtraction (new primitive)
	integerMethodDict.Entries[minusSelector.SymbolValue] = minusMethod

	// * method
	timesSelector := NewSymbol("*")
	timesMethod := NewMethod(timesSelector, result)
	timesMethod.Method.IsPrimitive = true
	timesMethod.Method.PrimitiveIndex = 2 // Multiplication
	integerMethodDict.Entries[timesSelector.SymbolValue] = timesMethod

	// = method
	equalsSelector := NewSymbol("=")
	equalsMethod := NewMethod(equalsSelector, result)
	equalsMethod.Method.IsPrimitive = true
	equalsMethod.Method.PrimitiveIndex = 3 // Equality
	integerMethodDict.Entries[equalsSelector.SymbolValue] = equalsMethod

	// < method
	lessSelector := NewSymbol("<")
	lessMethod := NewMethod(lessSelector, result)
	lessMethod.Method.IsPrimitive = true
	lessMethod.Method.PrimitiveIndex = 6 // Less than
	integerMethodDict.Entries[lessSelector.SymbolValue] = lessMethod

	// > method
	greaterSelector := NewSymbol(">")
	greaterMethod := NewMethod(greaterSelector, result)
	greaterMethod.Method.IsPrimitive = true
	greaterMethod.Method.PrimitiveIndex = 7 // Greater than
	integerMethodDict.Entries[greaterSelector.SymbolValue] = greaterMethod

	return result
}

// NewInteger creates a new integer object with the specified class or VM's IntegerClass
func (vm *VM) NewInteger(value int64) *Object {
	intClass := vm.IntegerClass
	return &Object{
		Type:         OBJ_INTEGER,
		IntegerValue: value,
		Class:        intClass,
	}
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
		// Check if we need to collect garbage
		if vm.ObjectMemory.ShouldCollect() {
			vm.ObjectMemory.Collect(vm)
		}

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
		if receiver.Type == OBJ_INTEGER && len(args) == 1 && args[0].Type == OBJ_INTEGER {
			result := receiver.IntegerValue + args[0].IntegerValue
			return vm.NewInteger(result)
		}
	case 2: // Multiplication
		if receiver.Type == OBJ_INTEGER && len(args) == 1 && args[0].Type == OBJ_INTEGER {
			result := receiver.IntegerValue * args[0].IntegerValue
			return vm.NewInteger(result)
		}
	case 3: // Equality
		if receiver.Type == OBJ_INTEGER && len(args) == 1 && args[0].Type == OBJ_INTEGER {
			result := receiver.IntegerValue == args[0].IntegerValue
			return NewBoolean(result)
		}
	case 4: // Subtraction
		if receiver.Type == OBJ_INTEGER && len(args) == 1 && args[0].Type == OBJ_INTEGER {
			result := receiver.IntegerValue - args[0].IntegerValue
			return vm.NewInteger(result)
		}
	case 5: // basicClass - return the class of the receiver
		if len(args) == 0 {
			class := vm.GetClass(receiver)
			fmt.Printf("executePrimitive: basicClass returning %v\n", class)
			return class
		}
	case 6: // Less than
		if receiver.Type == OBJ_INTEGER && len(args) == 1 && args[0].Type == OBJ_INTEGER {
			result := receiver.IntegerValue < args[0].IntegerValue
			return NewBoolean(result)
		}
	case 7: // Greater than
		if receiver.Type == OBJ_INTEGER && len(args) == 1 && args[0].Type == OBJ_INTEGER {
			result := receiver.IntegerValue > args[0].IntegerValue
			return NewBoolean(result)
		}
	default:
		panic("executePrimitive: unknown primitive index\n")
	}
	return nil // Fall through to method
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

	// Get the class of the receiver using the central GetClass function
	class := vm.GetClass(receiver)

	// Check for nil class
	if class == nil {
		panic("lookupMethod: nil class\n")
	}

	// Look up the method in the class hierarchy
	for class != nil {
		// Check if the class has a method dictionary
		methodDict := class.GetMethodDict()
		if methodDict != nil && methodDict.Type == OBJ_DICTIONARY && methodDict.Entries != nil {
			// Check if the method dictionary has the selector
			if method, ok := methodDict.Entries[selector.SymbolValue]; ok {
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
