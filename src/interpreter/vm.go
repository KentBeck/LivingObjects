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
	TrueObject   *Object
	FalseObject  *Object
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
	vm.NilObject = NewNil()
	vm.TrueObject = NewBoolean(true)
	vm.FalseObject = NewBoolean(false)
	vm.ObjectClass = NewClass("Object", vm.NilObject)
	vm.IntegerClass = vm.NewIntegerClass()

	return vm
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
	// For now, we'll just create a simple test image in memory
	// In a real implementation, this would load from a file

	// Create basic classes
	vm.ObjectClass = NewClass("Object", nil)
	vm.Globals["Object"] = vm.ObjectClass

	// The method dictionary is already created in NewClass at index 0

	// Create a simple test method: 2 + 3
	twoObj := vm.NewInteger(2)
	threeObj := vm.NewInteger(3)
	plusSymbol := NewSymbol("+")

	// Create a method that adds 2 and 3
	methodObj := NewMethod(plusSymbol, vm.ObjectClass)

	// Add literals
	methodObj.Method.Literals = append(methodObj.Method.Literals, twoObj)
	methodObj.Method.Literals = append(methodObj.Method.Literals, threeObj)
	methodObj.Method.Literals = append(methodObj.Method.Literals, plusSymbol)

	// Create bytecodes for: 2 + 3
	// PUSH_LITERAL 0 (2)
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, PUSH_LITERAL)
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// PUSH_LITERAL 1 (3)
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, PUSH_LITERAL)
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, 0, 0, 0, 1) // Index 1

	// SEND_MESSAGE 2 ("+") with 1 argument
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, SEND_MESSAGE)
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, 0, 0, 0, 2) // Selector index 2
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

	// RETURN_STACK_TOP
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, RETURN_STACK_TOP)

	// Create a context for this method
	vm.CurrentContext = NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

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

		case SET_CLASS:
			err = vm.ExecuteSetClass(context)

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
func (vm *VM) executePrimitive(receiver *Object, selector *Object, args []*Object) *Object {
	// Check for nil receiver or selector
	if receiver == nil || selector == nil {
		return nil
	}

	// First, check if the method is explicitly marked as a primitive
	method := vm.lookupMethod(receiver, selector)
	if method != nil && method.Type == OBJ_METHOD && method.Method.IsPrimitive {
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
		}
	}

	// If not a primitive method, handle built-in operations
	if receiver.Type == OBJ_INTEGER {
		switch selector.SymbolValue {
		case "+":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER {
				return vm.NewInteger(receiver.IntegerValue + args[0].IntegerValue)
			}
		case "-":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER {
				return vm.NewInteger(receiver.IntegerValue - args[0].IntegerValue)
			}
		case "*":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER {
				return vm.NewInteger(receiver.IntegerValue * args[0].IntegerValue)
			}
		case "/":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER && args[0].IntegerValue != 0 {
				return vm.NewInteger(receiver.IntegerValue / args[0].IntegerValue)
			}
		case "=":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER {
				return NewBoolean(receiver.IntegerValue == args[0].IntegerValue)
			}
		case "<":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER {
				return NewBoolean(receiver.IntegerValue < args[0].IntegerValue)
			}
		case ">":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER {
				return NewBoolean(receiver.IntegerValue > args[0].IntegerValue)
			}
		}
	}

	// No primitive method found
	return nil
}

// lookupMethod looks up a method in a class hierarchy
func (vm *VM) lookupMethod(receiver *Object, selector *Object) *Object {
	// Check for nil receiver or selector
	if receiver == nil || selector == nil {
		return nil
	}

	// Get the class of the receiver
	class := receiver
	if receiver.Type != OBJ_CLASS {
		class = receiver.Class
	}

	// Check for nil class
	if class == nil {
		return nil
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
		}

		// Move up the class hierarchy
		class = class.SuperClass
	}

	// Method not found
	return nil
}
