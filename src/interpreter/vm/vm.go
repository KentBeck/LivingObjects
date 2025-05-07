package vm

import (
	"fmt"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// VM represents the Smalltalk virtual machine
type VM struct {
	Globals        map[string]*core.Object
	CurrentContext *Context
	ObjectMemory   *core.ObjectMemory

	// Special objects
	NilObject    core.ObjectInterface
	NilClass     *classes.Class
	TrueObject   core.ObjectInterface
	TrueClass    *classes.Class
	FalseObject  core.ObjectInterface
	FalseClass   *classes.Class
	ObjectClass  *classes.Class
	IntegerClass *classes.Class
	FloatClass   *classes.Class
	StringClass  *classes.Class
	BlockClass   *classes.Class
}

// NewVM creates a new virtual machine
func NewVM() *VM {
	vm := &VM{
		Globals:      make(map[string]*core.Object),
		ObjectMemory: core.NewObjectMemory(),
	}

	// Initialize special objects
	vm.ObjectClass = vm.NewObjectClass()
	vm.NilClass = classes.NewClass("UndefinedObject", vm.ObjectClass)
	vm.NilObject = core.MakeNilImmediate()
	vm.TrueClass = classes.NewClass("True", vm.ObjectClass)
	vm.TrueObject = core.MakeTrueImmediate()
	vm.FalseClass = classes.NewClass("False", vm.ObjectClass)
	vm.FalseObject = core.MakeFalseImmediate()
	vm.IntegerClass = vm.NewIntegerClass()
	vm.FloatClass = vm.NewFloatClass()
	vm.StringClass = vm.NewStringClass()
	vm.BlockClass = vm.NewBlockClass()

	// Register the VM as a block executor
	vm.RegisterAsBlockExecutor()

	return vm
}

func (vm *VM) NewObjectClass() *classes.Class {
	result := classes.NewClass("Object", nil) // patch this up later. then even later when we have real images all this initialization can go away

	// Add basicClass method to Object class
	// TODO: Implement method builder in compiler package
	// NewMethodBuilder(result).
	// 	Selector("basicClass").
	// 	Primitive(5). // basicClass primitive
	// 	Go()

	return result
}

func (vm *VM) NewIntegerClass() *classes.Class {
	result := classes.NewClass("Integer", vm.ObjectClass)

	// Add primitive methods to the Integer class
	// TODO: Implement method builder in compiler package
	// builder := NewMethodBuilder(result)

	// // + method (addition)
	// builder.Selector("+").Primitive(1).Go()

	// // - method (subtraction)
	// builder.Selector("-").Primitive(4).Go()

	// // * method (multiplication)
	// builder.Selector("*").Primitive(2).Go()

	// // = method (equality)
	// builder.Selector("=").Primitive(3).Go()

	// // < method (less than)
	// builder.Selector("<").Primitive(6).Go()

	// // > method (greater than)
	// builder.Selector(">").Primitive(7).Go()

	return result
}

func (vm *VM) NewFloatClass() *classes.Class {
	result := classes.NewClass("Float", vm.ObjectClass) // patch this up later. then even later when we have real images all this initialization can go away

	// Add primitive methods to the Float class
	// TODO: Implement method builder in compiler package
	// builder := NewMethodBuilder(result)

	// // + method (addition)
	// builder.Selector("+").Primitive(10).Go()

	// // - method (subtraction)
	// builder.Selector("-").Primitive(11).Go()

	// // * method (multiplication)
	// builder.Selector("*").Primitive(12).Go()

	// // / method (division)
	// builder.Selector("/").Primitive(13).Go()

	// // = method (equality)
	// builder.Selector("=").Primitive(14).Go()

	// // < method (less than)
	// builder.Selector("<").Primitive(15).Go()

	// // > method (greater than)
	// builder.Selector(">").Primitive(16).Go()

	return result
}

// NewInteger creates a new integer object
// This returns an immediate value for integers
func (vm *VM) NewInteger(value int64) *core.Object {
	// Check if the value fits in 62 bits
	if value <= 0x1FFFFFFFFFFFFFFF && value >= -0x2000000000000000 {
		// Use immediate integer
		return core.MakeIntegerImmediate(value)
	}

	// Panic for large values that don't fit in 62 bits
	panic("Integer value too large for immediate representation")
}

func (vm *VM) NewFloat(value float64) *core.Object {
	return core.MakeFloatImmediate(value)
}

func (vm *VM) NewStringClass() *classes.Class {
	result := classes.NewClass("String", vm.ObjectClass)

	// Add primitive methods to the String class
	// Add the , method (concatenation)
	commaMethod := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes:      []byte{},
		Literals:       []*core.Object{},
		TempVarNames:   []string{},
		IsPrimitive:    true,
		PrimitiveIndex: 30, // Primitive index for string concatenation
	}
	result.AddMethod(classes.NewSymbol(","), classes.MethodToObject(commaMethod))

	return result
}

func (vm *VM) NewBlockClass() *classes.Class {
	result := classes.NewClass("Block", vm.ObjectClass)

	// Add primitive methods to the Block class
	// TODO: Implement method builder in compiler package
	// builder := NewMethodBuilder(result)

	// // new method (creates a new block instance)
	// // fixme sketchy
	// builder.Selector("new").Primitive(20).Go()

	// // value method (executes the block with no arguments)
	// builder.Selector("value").Primitive(21).Go()

	// // value: method (executes the block with one argument)
	// builder.Selector("value:").Primitive(22).Go()

	return result
}

// LoadImage loads a Smalltalk image from a file
func (vm *VM) LoadImage(path string) error {
	vm.Globals["Object"] = classes.ClassToObject(vm.ObjectClass)

	return nil
}

// Execute executes the current context
func (vm *VM) Execute() (core.ObjectInterface, error) {
	var finalResult core.ObjectInterface

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
func (vm *VM) ExecuteContext(context *Context) (core.ObjectInterface, error) {
	// Execute the context

	for {
		// Get the method
		method := classes.ObjectToMethod(context.Method)

		// Check if we've reached the end of the method
		if context.PC >= len(method.GetBytecodes()) {
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
		bytecode := method.GetBytecodes()[context.PC]

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

		case CREATE_BLOCK:
			err = vm.ExecuteCreateBlock(context)

		case EXECUTE_BLOCK:
			returnValue, err := vm.ExecuteExecuteBlock(context)
			if err == nil {
				if returnValue != nil {
					// We got a result from executing the block
					// Continue execution in the current context
					context.PC += size
					continue
				} else {
					// A nil return value with no error means we've started a new context
					return vm.NilObject, nil
				}
			}

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

// GetClass returns the class of an object
// This is the single function that should be used to get the class of an object
func (vm *VM) GetClass(obj *core.Object) *classes.Class {
	if obj == nil {
		panic("GetClass: nil object")
	}

	// Check if it's an immediate value
	if core.IsImmediate(obj) {
		// Handle immediate nil
		if core.IsNilImmediate(obj) {
			return vm.NilClass
		}
		// Handle immediate true
		if core.IsTrueImmediate(obj) {
			return vm.TrueClass
		}
		// Handle immediate false
		if core.IsFalseImmediate(obj) {
			return vm.FalseClass
		}
		// Handle immediate integer
		if core.IsIntegerImmediate(obj) {
			return vm.IntegerClass
		}
		// Handle immediate float
		if core.IsFloatImmediate(obj) {
			return vm.FloatClass
		}
		// Other immediate types will be added later
		panic("GetClass: unknown immediate type")
	}

	// If it's a regular object, proceed as before

	// If the object is a class, return itself
	if obj.Type() == core.OBJ_CLASS {
		return classes.ObjectToClass(obj) // Later Metaclass
	}

	// Special case for nil object (legacy non-immediate nil)
	if obj.Type() == core.OBJ_NIL {
		return nil
	}

	// Otherwise, return the class field
	if obj.Class() == nil {
		panic("GetClass: object has nil class")
	}

	return classes.ObjectToClass(obj.Class())
}

// LookupMethod looks up a method in a class hierarchy
func (vm *VM) LookupMethod(receiver *core.Object, selector core.ObjectInterface) *core.Object {
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
		methodDict := class.GetMethodDictionary()
		if methodDict != nil && methodDict.GetEntryCount() > 0 {
			// Check if the method dictionary has the selector
			selectorSymbol := classes.ObjectToSymbol(selector.(*core.Object))
			if method := methodDict.GetEntry(selectorSymbol.GetValue()); method != nil {
				return method
			}
		}

		// Move up the class hierarchy
		class = classes.ObjectToClass(class.GetSuperClass())
	}

	// Method not found
	return nil
}

// ExecutePrimitive executes a primitive method
func (vm *VM) ExecutePrimitive(receiver *core.Object, selector *core.Object, args []*core.Object, method *core.Object) *core.Object {
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
	default:
		panic("executePrimitive: unknown primitive index\n")
	}
	return nil // Fall through to method
}

// GetGlobals returns the globals map
func (vm *VM) GetGlobals() []*core.Object {
	// Convert map to slice for memory management
	globals := make([]*core.Object, 0, len(vm.Globals))
	for _, obj := range vm.Globals {
		globals = append(globals, obj)
	}
	return globals
}

// GetCurrentContext returns the current context
func (vm *VM) GetCurrentContext() interface{} {
	return vm.CurrentContext
}

// GetObjectClass returns the object class
func (vm *VM) GetObjectClass() *classes.Class {
	return vm.ObjectClass
}
