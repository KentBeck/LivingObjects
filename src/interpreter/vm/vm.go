package vm

import (
	"fmt"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/core"
)

// VM represents the Smalltalk virtual machine
type VM struct {
	Globals      map[string]*core.Object
	ObjectMemory *core.ObjectMemory
	Executor     *Executor

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
	ArrayClass   *classes.Class
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
	vm.TrueClass = vm.NewTrueClass()
	vm.TrueObject = core.MakeTrueImmediate()
	vm.FalseClass = vm.NewFalseClass()
	vm.FalseObject = core.MakeFalseImmediate()
	vm.IntegerClass = vm.NewIntegerClass()
	vm.FloatClass = vm.NewFloatClass()
	vm.StringClass = vm.NewStringClass()
	vm.BlockClass = vm.NewBlockClass()
	vm.ArrayClass = vm.NewArrayClass()

	// Initialize the executor
	vm.Executor = NewExecutor(vm)

	// Register the VM as a block executor
	vm.RegisterAsBlockExecutor()

	return vm
}

func (vm *VM) NewObjectClass() *classes.Class {
	result := classes.NewClass("Object", nil) // patch this up later. then even later when we have real images all this initialization can go away

	// Add basicClass method to Object class
	compiler.NewMethodBuilder(result).
		Selector("basicClass").
		Primitive(5). // basicClass primitive
		Go()

	return result
}

func (vm *VM) NewIntegerClass() *classes.Class {
	result := classes.NewClass("Integer", vm.ObjectClass)

	// Add primitive methods to the Integer class
	builder := compiler.NewMethodBuilder(result)

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

func (vm *VM) NewFloatClass() *classes.Class {
	result := classes.NewClass("Float", vm.ObjectClass) // patch this up later. then even later when we have real images all this initialization can go away

	// Add primitive methods to the Float class
	builder := compiler.NewMethodBuilder(result)

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

// NewString creates a new string object
func (vm *VM) NewString(value string) *core.Object {
	str := classes.NewString(value)
	strObj := classes.StringToObject(str)
	strObj.SetClass(classes.ClassToObject(vm.StringClass))
	return strObj
}

// NewArray creates a new array object
func (vm *VM) NewArray(size int) *core.Object {
	array := classes.NewArray(size)
	arrayObj := classes.ArrayToObject(array)
	arrayObj.SetClass(classes.ClassToObject(vm.ArrayClass))
	return arrayObj
}

func (vm *VM) NewTrueClass() *classes.Class {
	result := classes.NewClass("True", vm.ObjectClass)

	// Add primitive methods to the True class
	builder := compiler.NewMethodBuilder(result)

	// not method (returns false)
	builder.Selector("not").Primitive(50).Go()

	return result
}

func (vm *VM) NewFalseClass() *classes.Class {
	result := classes.NewClass("False", vm.ObjectClass)

	// Add primitive methods to the False class
	builder := compiler.NewMethodBuilder(result)

	// not method (returns true)
	builder.Selector("not").Primitive(51).Go()

	return result
}

func (vm *VM) NewStringClass() *classes.Class {
	result := classes.NewClass("String", vm.ObjectClass)

	// Add primitive methods to the String class
	builder := compiler.NewMethodBuilder(result)

	// size method (returns the length of the string)
	builder.Selector("size").Primitive(30).Go()

	return result
}

func (vm *VM) NewArrayClass() *classes.Class {
	result := classes.NewClass("Array", vm.ObjectClass)

	// Add primitive methods to the Array class
	builder := compiler.NewMethodBuilder(result)

	// at: method (returns the element at the given index)
	builder.Selector("at:").Primitive(40).Go()

	return result
}

func (vm *VM) NewBlockClass() *classes.Class {
	result := classes.NewClass("Block", vm.ObjectClass)

	// Add primitive methods to the Block class
	builder := compiler.NewMethodBuilder(result)

	// new method (creates a new block instance)
	// fixme sketchy
	builder.Selector("new").Primitive(20).Go()

	// value method (executes the block with no arguments)
	builder.Selector("value").Primitive(21).Go()

	// value: method (executes the block with one argument)
	builder.Selector("value:").Primitive(22).Go()

	return result
}

// LoadImage loads a Smalltalk image from a file
func (vm *VM) LoadImage(path string) error {
	vm.Globals["Object"] = classes.ClassToObject(vm.ObjectClass)

	return nil
}

// Execute executes the current context
func (vm *VM) Execute() (core.ObjectInterface, error) {
	// Execute using the executor
	return vm.Executor.Execute()
}

// ExecuteContext executes a single context until it returns
func (vm *VM) ExecuteContext(context *Context) (core.ObjectInterface, error) {
	// Set the context in the Executor
	vm.Executor.CurrentContext = context

	// Execute
	return vm.Executor.ExecuteContext(context)
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
		panic(fmt.Sprintf("GetClass: object has nil class %s\n", obj.String()))
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
			blockInstance := classes.NewBlock(vm.Executor.CurrentContext)
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
	case 30: // String size - return the length of the string
		if receiver.Type() == core.OBJ_STRING {
			// Get the string
			str := classes.ObjectToString(receiver)

			// Get the length of the string
			length := str.Length()

			// Return the length as an integer
			return vm.NewInteger(int64(length))
		}
	case 40: // Array at: - return the element at the given index
		if receiver.Type() == core.OBJ_ARRAY && len(args) == 1 && core.IsIntegerImmediate(args[0]) {
			// Get the array
			array := classes.ObjectToArray(receiver)

			// Get the index (1-based in Smalltalk, 0-based in Go)
			index := core.GetIntegerImmediate(args[0]) - 1

			// Check bounds
			if index < 0 || int(index) >= array.Size() {
				panic(fmt.Sprintf("Array index out of bounds: %d", index+1))
			}

			// Return the element at the given index
			return array.At(int(index))
		}
	case 50: // True not - return false
		if core.IsTrueImmediate(receiver) {
			return core.MakeFalseImmediate()
		}
	case 51: // False not - return true
		if core.IsFalseImmediate(receiver) {
			return core.MakeTrueImmediate()
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

// GetObjectClass returns the object class
func (vm *VM) GetObjectClass() *classes.Class {
	return vm.ObjectClass
}
