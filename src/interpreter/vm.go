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
	BlockClass   *Object
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
	vm.BlockClass = vm.NewBlockClass()

	return vm
}

func (vm *VM) NewObjectClass() *Object {
	result := NewClass("Object", nil) // patch this up later. then even later when we have real images all this initialization can go away

	// Add basicClass method to Object class
	NewMethodBuilder(result).
		Selector("basicClass").
		Primitive(5). // basicClass primitive
		Go()

	return result
}

func (vm *VM) NewIntegerClass() *Object {
	result := NewClass("Integer", vm.ObjectClass)

	// Add primitive methods to the Integer class
	builder := NewMethodBuilder(result)

	// + method (addition)
	builder.Selector("+").Primitive(1).Go()

	// - method (subtraction)
	builder.Selector("-").Primitive(4).Go()

	// * method (multiplication)
	builder.Selector("*").Primitive(2).Go()

	// = method (equality)
	builder.Selector("=").Primitive(3).Go()

	// < method (less than)
	builder.Selector("<").Primitive(6).Go()

	// > method (greater than)
	builder.Selector(">").Primitive(7).Go()

	return result
}

func (vm *VM) NewFloatClass() *Object {
	result := NewClass("Float", vm.ObjectClass) // patch this up later. then even later when we have real images all this initialization can go away

	// Add primitive methods to the Float class
	builder := NewMethodBuilder(result)

	// + method (addition)
	builder.Selector("+").Primitive(10).Go()

	// - method (subtraction)
	builder.Selector("-").Primitive(11).Go()

	// * method (multiplication)
	builder.Selector("*").Primitive(12).Go()

	// / method (division)
	builder.Selector("/").Primitive(13).Go()

	// = method (equality)
	builder.Selector("=").Primitive(14).Go()

	// < method (less than)
	builder.Selector("<").Primitive(15).Go()

	// > method (greater than)
	builder.Selector(">").Primitive(16).Go()

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

func (vm *VM) NewBlockClass() *Object {
	result := NewClass("Block", vm.ObjectClass)

	// Add primitive methods to the Block class
	builder := NewMethodBuilder(result)

	// new method (creates a new block instance)
	builder.Selector("new").Primitive(20).Go()

	// value method (executes the block with no arguments)
	builder.Selector("value").Primitive(21).Go()

	// value: method (executes the block with one argument)
	builder.Selector("value:").Primitive(22).Go()

	return result
}

// LoadImage loads a Smalltalk image from a file
func (vm *VM) LoadImage(path string) error {
	vm.Globals["Object"] = vm.ObjectClass

	return nil
}

// Execute executes the current context
func (vm *VM) Execute() (ObjectInterface, error) {
	var finalResult ObjectInterface

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
			returnValue, err := vm.ExecuteSendMessage(context)
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
			returnValue, err := vm.ExecuteReturnStackTop(context)
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
		if methodDict != nil && methodDict.Type() == OBJ_DICTIONARY && methodDict.Entries != nil {
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
