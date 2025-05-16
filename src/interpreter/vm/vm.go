package vm

import (
	"fmt"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/pile"
)

// VM represents the Smalltalk virtual machine
type VM struct {
	// Globals map holds all global variables including classes
	Globals      map[string]*pile.Object
	ObjectMemory *pile.ObjectMemory
	Executor     *Executor

	// Special objects
	NilObject   pile.ObjectInterface
	TrueObject  pile.ObjectInterface
	FalseObject pile.ObjectInterface
}

// NewVM creates a new virtual machine
func NewVM() *VM {
	vm := &VM{
		Globals:      make(map[string]*pile.Object),
		ObjectMemory: pile.NewObjectMemory(),
	}

	// Initialize special immediate objects
	vm.NilObject = pile.MakeNilImmediate()
	vm.TrueObject = pile.MakeTrueImmediate()
	vm.FalseObject = pile.MakeFalseImmediate()

	// Register the VM as the default factory for creating objects
	pile.RegisterFactory(vm)

	// Initialize core classes
	objectClass := vm.NewObjectClass()
	vm.Globals["Object"] = pile.ClassToObject(objectClass)

	classClass := vm.NewClassClass()
	vm.Globals["Class"] = pile.ClassToObject(classClass)

	nilClass := pile.NewClass("UndefinedObject", objectClass)
	vm.Globals["UndefinedObject"] = pile.ClassToObject(nilClass)

	trueClass := vm.NewTrueClass()
	vm.Globals["True"] = pile.ClassToObject(trueClass)

	falseClass := vm.NewFalseClass()
	vm.Globals["False"] = pile.ClassToObject(falseClass)

	integerClass := vm.NewIntegerClass()
	vm.Globals["Integer"] = pile.ClassToObject(integerClass)

	floatClass := vm.NewFloatClass()
	vm.Globals["Float"] = pile.ClassToObject(floatClass)

	stringClass := vm.NewStringClass()
	vm.Globals["String"] = pile.ClassToObject(stringClass)

	blockClass := vm.NewBlockClass()
	vm.Globals["Block"] = pile.ClassToObject(blockClass)

	arrayClass := vm.NewArrayClass()
	vm.Globals["Array"] = pile.ClassToObject(arrayClass)

	byteArrayClass := vm.NewByteArrayClass()
	vm.Globals["ByteArray"] = pile.ClassToObject(byteArrayClass)

	// Initialize the executor
	vm.Executor = NewExecutor(vm)

	// Register the VM as a block executor
	vm.RegisterAsBlockExecutor()

	return vm
}

func (vm *VM) NewObjectClass() *pile.Class {
	result := pile.NewClass("Object", nil) // patch this up later. then even later when we have real images all this initialization can go away

	// Add methods to Object class
	builder := compiler.NewMethodBuilder(result)

	// basicClass method
	builder.Selector("basicClass").
		Primitive(5). // basicClass primitive
		Go()

	// should not be here
	// new method (creates a new instance of the class)
	builder.Selector("new").
		Primitive(60). // new primitive
		Go()

	// class method - a more user-friendly name for accessing an object's class
	// class implementation: ^self basicClass
	classSelectorIndex, builder := builder.AddLiteral(pile.NewSymbol("basicClass"))
	builder.Selector("class").
		PushSelf().                         // self
		SendMessage(classSelectorIndex, 0). // self basicClass
		ReturnStackTop().                   // ^
		Go()

	return result
}

func (vm *VM) NewIntegerClass() *pile.Class {
	objectClass := pile.ObjectToClass(vm.Globals["Object"])
	result := pile.NewClass("Integer", objectClass)

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

func (vm *VM) NewFloatClass() *pile.Class {
	objectClass := pile.ObjectToClass(vm.Globals["Object"])
	result := pile.NewClass("Float", objectClass) // then even later when we have real images all this initialization can go away

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
func (vm *VM) NewInteger(value int64) *pile.Object {
	// Check if the value fits in 62 bits
	if value <= 0x1FFFFFFFFFFFFFFF && value >= -0x2000000000000000 {
		// Use immediate integer
		return pile.MakeIntegerImmediate(value)
	}

	// Panic for large values that don't fit in 62 bits
	panic("Integer value too large for immediate representation")
}

func (vm *VM) NewFloat(value float64) *pile.Object {
	return pile.MakeFloatImmediate(value)
}

// NewString creates a new string object
func (vm *VM) NewString(value string) *pile.Object {
	str := &pile.String{Object: pile.Object{TypeField: pile.OBJ_STRING}, Value: value}
	strObj := pile.StringToObject(str)
	strObj.SetClass(vm.Globals["String"])
	return strObj
}

// NewArray creates a new array object
func (vm *VM) NewArray(size int) *pile.Object {
	array := &pile.Array{Object: pile.Object{TypeField: pile.OBJ_ARRAY}, Elements: make([]*pile.Object, size)}
	arrayObj := pile.ArrayToObject(array)
	arrayObj.SetClass(vm.Globals["Array"])
	return arrayObj
}

// NewTrue returns the true object
func (vm *VM) NewTrue() *pile.Object {
	return pile.MakeTrueImmediate()
}

// NewFalse returns the false object
func (vm *VM) NewFalse() *pile.Object {
	return pile.MakeFalseImmediate()
}

// NewNil returns the nil object
func (vm *VM) NewNil() *pile.Object {
	return pile.MakeNilImmediate()
}

func (vm *VM) NewTrueClass() *pile.Class {
	objectClass := pile.ObjectToClass(vm.Globals["Object"])
	result := pile.NewClass("True", objectClass)

	// Add methods to the True class
	builder := compiler.NewMethodBuilder(result)

	// Create a literal for false
	falseIndex, builder := builder.AddLiteral(pile.MakeFalseImmediate())

	// not method (returns false)
	builder.Selector("not").
		PushLiteral(falseIndex).
		ReturnStackTop().
		Go()

	return result
}

func (vm *VM) NewFalseClass() *pile.Class {
	objectClass := pile.ObjectToClass(vm.Globals["Object"])
	result := pile.NewClass("False", objectClass)

	// Add methods to the False class
	builder := compiler.NewMethodBuilder(result)

	// Create a literal for true
	trueIndex, builder := builder.AddLiteral(pile.MakeTrueImmediate())

	// not method (returns true)
	builder.Selector("not").
		PushLiteral(trueIndex).
		ReturnStackTop().
		Go()

	return result
}

func (vm *VM) NewStringClass() *pile.Class {
	objectClass := pile.ObjectToClass(vm.Globals["Object"])
	result := pile.NewClass("String", objectClass)

	// Add primitive methods to the String class
	builder := compiler.NewMethodBuilder(result)

	// size method (returns the length of the string)
	builder.Selector("size").Primitive(30).Go()

	return result
}

func (vm *VM) NewArrayClass() *pile.Class {
	objectClass := pile.ObjectToClass(vm.Globals["Object"])
	result := pile.NewClass("Array", objectClass)

	// Add primitive methods to the Array class
	builder := compiler.NewMethodBuilder(result)

	// at: method (returns the element at the given index)
	builder.Selector("at:").Primitive(40).Go()

	return result
}

func (vm *VM) NewBlockClass() *pile.Class {
	objectClass := pile.ObjectToClass(vm.Globals["Object"])
	result := pile.NewClass("Block", objectClass)

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

func (vm *VM) NewClassClass() *pile.Class {
	objectClass := pile.ObjectToClass(vm.Globals["Object"])
	result := pile.NewClass("Class", objectClass)

	// Add primitive methods to the Class class
	builder := compiler.NewMethodBuilder(result)

	// new method (creates a new instance of the class)
	builder.Selector("new").Primitive(60).Go()

	return result
}

// LoadImage loads a Smalltalk image from a file
func (vm *VM) LoadImage(path string) error {
	// Currently just returns nil as globals are already set up during initialization
	return nil
}

// Execute executes the current context
func (vm *VM) Execute() (pile.ObjectInterface, error) {
	// Execute using the executor
	return vm.Executor.Execute()
}

// ExecuteContext executes a single context until it returns
func (vm *VM) ExecuteContext(context *Context) (pile.ObjectInterface, error) {
	// Set the context in the Executor
	vm.Executor.CurrentContext = context

	// Execute
	return vm.Executor.ExecuteContext(context)
}

// GetClass returns the class of an object
// This is the single function that should be used to get the class of an object
func (vm *VM) GetClass(obj *pile.Object) *pile.Class {
	if obj == nil {
		panic("GetClass: nil object")
	}

	// Check if it's an immediate value
	if pile.IsImmediate(obj) {
		// Handle immediate nil
		if pile.IsNilImmediate(obj) {
			if classObj, ok := vm.Globals["UndefinedObject"]; ok {
				return pile.ObjectToClass(classObj)
			}
			panic("GetClass: UndefinedObject class not found in globals")
		}
		// Handle immediate true
		if pile.IsTrueImmediate(obj) {
			if classObj, ok := vm.Globals["True"]; ok {
				return pile.ObjectToClass(classObj)
			}
			panic("GetClass: True class not found in globals")
		}
		// Handle immediate false
		if pile.IsFalseImmediate(obj) {
			if classObj, ok := vm.Globals["False"]; ok {
				return pile.ObjectToClass(classObj)
			}
			panic("GetClass: False class not found in globals")
		}
		// Handle immediate integer
		if pile.IsIntegerImmediate(obj) {
			if classObj, ok := vm.Globals["Integer"]; ok {
				return pile.ObjectToClass(classObj)
			}
			panic("GetClass: Integer class not found in globals")
		}
		// Handle immediate float
		if pile.IsFloatImmediate(obj) {
			if classObj, ok := vm.Globals["Float"]; ok {
				return pile.ObjectToClass(classObj)
			}
			panic("GetClass: Float class not found in globals")
		}
		// Other immediate types will be added later
		panic("GetClass: unknown immediate type")
	}

	// If it's a regular object, proceed as before

	// If the object is a class, return itself
	if obj.Type() == pile.OBJ_CLASS {
		return pile.ObjectToClass(obj) // Later Metaclass
	}

	// Special case for nil object (legacy non-immediate nil)
	if obj.Type() == pile.OBJ_NIL {
		return nil
	}

	// Otherwise, return the class field
	if obj.Class() == nil {
		panic(fmt.Sprintf("GetClass: object has nil class %s\n", obj.String()))
	}

	return pile.ObjectToClass(obj.Class())
}

// LookupMethod looks up a method in a class hierarchy
func (vm *VM) LookupMethod(receiver *pile.Object, selector pile.ObjectInterface) *pile.Object {
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
		methodDict := pile.ObjectToDictionary(class.MethodDictionary)
		if methodDict != nil && methodDict.GetEntryCount() > 0 {
			// Check if the method dictionary has the selector
			selectorSymbol := pile.ObjectToSymbol(selector.(*pile.Object))
			if method := methodDict.GetEntry(selectorSymbol.GetValue()); method != nil {
				return method
			}
		}

		// Move up the class hierarchy
		class = pile.ObjectToClass(class.SuperClass)
	}

	// Method not found
	return nil
}

// ExecutePrimitive executes a primitive method
func (vm *VM) ExecutePrimitive(receiver *pile.Object, selector *pile.Object, args []*pile.Object, method *pile.Object) *pile.Object {
	if receiver == nil {
		panic("executePrimitive: nil receiver\n")
	}
	if selector == nil {
		panic("executePrimitive: nil selector\n")
	}
	if method == nil {
		panic("executePrimitive: nil method\n")
	}
	if method.Type() != pile.OBJ_METHOD {
		return nil
	}
	methodObj := pile.ObjectToMethod(method)
	if !methodObj.IsPrimitiveMethod() {
		return nil
	}

	// Execute the primitive based on its index
	switch methodObj.GetPrimitiveIndex() {
	case 1: // Addition
		// Handle immediate integers
		if pile.IsIntegerImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetIntegerImmediate(receiver)
			val2 := pile.GetIntegerImmediate(args[0])
			result := val1 + val2
			return vm.NewInteger(result)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type() == pile.OBJ_INTEGER || (len(args) > 0 && args[0].Type() == pile.OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
		// Handle integer + float
		if pile.IsIntegerImmediate(receiver) && len(args) == 1 && pile.IsFloatImmediate(args[0]) {
			val1 := float64(pile.GetIntegerImmediate(receiver))
			val2 := pile.GetFloatImmediate(args[0])
			result := val1 + val2
			return vm.NewFloat(result)
		}
	case 2: // Multiplication
		// Handle immediate integers
		if pile.IsIntegerImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetIntegerImmediate(receiver)
			val2 := pile.GetIntegerImmediate(args[0])
			result := val1 * val2
			return vm.NewInteger(result)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type() == pile.OBJ_INTEGER || (len(args) > 0 && args[0].Type() == pile.OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 3: // Equality
		// Handle immediate integers
		if pile.IsIntegerImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetIntegerImmediate(receiver)
			val2 := pile.GetIntegerImmediate(args[0])
			result := val1 == val2
			return pile.NewBoolean(result).(*pile.Object)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type() == pile.OBJ_INTEGER || (len(args) > 0 && args[0].Type() == pile.OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 4: // Subtraction
		// Handle immediate integers
		if pile.IsIntegerImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetIntegerImmediate(receiver)
			val2 := pile.GetIntegerImmediate(args[0])
			result := val1 - val2
			return vm.NewInteger(result)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type() == pile.OBJ_INTEGER || (len(args) > 0 && args[0].Type() == pile.OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 5: // basicClass - return the class of the receiver
		class := vm.GetClass(receiver)
		return pile.ClassToObject(class)
	case 6: // Less than
		// Handle immediate integers
		if pile.IsIntegerImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetIntegerImmediate(receiver)
			val2 := pile.GetIntegerImmediate(args[0])
			result := val1 < val2
			return pile.NewBoolean(result).(*pile.Object)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type() == pile.OBJ_INTEGER || (len(args) > 0 && args[0].Type() == pile.OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 7: // Greater than
		// Handle immediate integers
		if pile.IsIntegerImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetIntegerImmediate(receiver)
			val2 := pile.GetIntegerImmediate(args[0])
			result := val1 > val2
			return pile.NewBoolean(result).(*pile.Object)
		}
		// Handle non-immediate integers - should panic
		if receiver.Type() == pile.OBJ_INTEGER || (len(args) > 0 && args[0].Type() == pile.OBJ_INTEGER) {
			panic("Non-immediate integer encountered")
		}
	case 10: // Float addition
		// Handle float + float
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsFloatImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := pile.GetFloatImmediate(args[0])
			result := val1 + val2
			return vm.NewFloat(result)
		}
		// Handle float + integer
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := float64(pile.GetIntegerImmediate(args[0]))
			result := val1 + val2
			return vm.NewFloat(result)
		}
	case 11: // Float subtraction
		// Handle float - float
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsFloatImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := pile.GetFloatImmediate(args[0])
			result := val1 - val2
			return vm.NewFloat(result)
		}
		// Handle float - integer
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := float64(pile.GetIntegerImmediate(args[0]))
			result := val1 - val2
			return vm.NewFloat(result)
		}
	case 12: // Float multiplication
		// Handle float * float
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsFloatImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := pile.GetFloatImmediate(args[0])
			result := val1 * val2
			return vm.NewFloat(result)
		}
		// Handle float * integer
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := float64(pile.GetIntegerImmediate(args[0]))
			result := val1 * val2
			return vm.NewFloat(result)
		}
	case 13: // Float division
		// Handle float / float
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsFloatImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := pile.GetFloatImmediate(args[0])
			result := val1 / val2
			return vm.NewFloat(result)
		}
		// Handle float / integer
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := float64(pile.GetIntegerImmediate(args[0]))
			result := val1 / val2
			return vm.NewFloat(result)
		}
	case 14: // Float equality
		// Handle float = float
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsFloatImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := pile.GetFloatImmediate(args[0])
			result := val1 == val2
			return pile.NewBoolean(result).(*pile.Object)
		}
		// Handle float = integer
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := float64(pile.GetIntegerImmediate(args[0]))
			result := val1 == val2
			return pile.NewBoolean(result).(*pile.Object)
		}
	case 15: // Float less than
		// Handle float < float
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsFloatImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := pile.GetFloatImmediate(args[0])
			result := val1 < val2
			return pile.NewBoolean(result).(*pile.Object)
		}
		// Handle float < integer
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := float64(pile.GetIntegerImmediate(args[0]))
			result := val1 < val2
			return pile.NewBoolean(result).(*pile.Object)
		}
	case 16: // Float greater than
		// Handle float > float
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsFloatImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := pile.GetFloatImmediate(args[0])
			result := val1 > val2
			return pile.NewBoolean(result).(*pile.Object)
		}
		// Handle float > integer
		if pile.IsFloatImmediate(receiver) && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			val1 := pile.GetFloatImmediate(receiver)
			val2 := float64(pile.GetIntegerImmediate(args[0]))
			result := val1 > val2
			return pile.NewBoolean(result).(*pile.Object)
		}
	case 20: // Block new - create a new block instance
		if receiver.Type() == pile.OBJ_CLASS && receiver == vm.Globals["Block"] {
			// Create a new block instance with proper class field
			blockInstance := vm.NewBlock(vm.Executor.CurrentContext)
			return blockInstance
		}
	case 21: // Block value - execute a block with no arguments
		if receiver.Type() == pile.OBJ_BLOCK {
			// Get the block
			block := pile.ObjectToBlock(receiver)

			// Create a method object for the block's bytecodes
			method := &pile.Method{
				Object: pile.Object{
					TypeField: pile.OBJ_METHOD,
				},
				Bytecodes: block.GetBytecodes(),
				Literals:  block.GetLiterals(),
			}
			methodObj := pile.MethodToObject(method)

			// Create a new context for the block execution
			blockContext := NewContext(methodObj, receiver, []*pile.Object{}, block.GetOuterContext().(*Context))

			// Execute the block's bytecodes
			result, err := vm.ExecuteContext(blockContext)
			if err != nil {
				panic(fmt.Sprintf("Error executing block: %v", err))
			}
			return result.(*pile.Object)
		}
	case 22: // Block value: - execute a block with one argument
		if receiver.Type() == pile.OBJ_BLOCK && len(args) == 1 {
			// Get the block
			block := pile.ObjectToBlock(receiver)

			// Create a method object for the block's bytecodes
			method := &pile.Method{
				Object: pile.Object{
					TypeField: pile.OBJ_METHOD,
				},
				Bytecodes: block.GetBytecodes(),
				Literals:  block.GetLiterals(),
			}
			methodObj := pile.MethodToObject(method)

			// Create a new context for the block execution
			blockContext := NewContext(methodObj, receiver, args, block.GetOuterContext().(*Context))

			// Execute the block's bytecodes
			result, err := vm.ExecuteContext(blockContext)
			if err != nil {
				panic(fmt.Sprintf("Error executing block: %v", err))
			}
			return result.(*pile.Object)
		}
	case 30: // String size - return the length of the string
		if receiver.Type() == pile.OBJ_STRING {
			// Get the string
			str := pile.ObjectToString(receiver)

			// Get the length of the string
			length := str.Length()

			// Return the length as an integer
			return vm.NewInteger(int64(length))
		}
	case 40: // Array at: - return the element at the given index
		if receiver.Type() == pile.OBJ_ARRAY && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			// Get the array
			array := pile.ObjectToArray(receiver)

			// Get the index (1-based in Smalltalk, 0-based in Go)
			index := pile.GetIntegerImmediate(args[0]) - 1

			// Check bounds
			if index < 0 || int(index) >= array.Size() {
				panic(fmt.Sprintf("Array index out of bounds: %d", index+1))
			}

			// Return the element at the given index
			return array.At(int(index))
		}
	case 50: // ByteArray at: - return the byte at the given index
		if receiver.Type() == pile.OBJ_BYTE_ARRAY && len(args) == 1 && pile.IsIntegerImmediate(args[0]) {
			// Get the byte array
			byteArray := pile.ObjectToByteArray(receiver)

			// Get the index (1-based in Smalltalk, 0-based in Go)
			index := pile.GetIntegerImmediate(args[0]) - 1

			// Check bounds
			if index < 0 || int(index) >= byteArray.Size() {
				panic(fmt.Sprintf("ByteArray index out of bounds: %d", index+1))
			}

			// Return the byte at the given index as an integer
			return vm.NewInteger(int64(byteArray.At(int(index))))
		}
	case 51: // ByteArray at:put: - set the byte at the given index
		if receiver.Type() == pile.OBJ_BYTE_ARRAY && len(args) == 2 &&
			pile.IsIntegerImmediate(args[0]) && pile.IsIntegerImmediate(args[1]) {
			// Get the byte array
			byteArray := pile.ObjectToByteArray(receiver)

			// Get the index (1-based in Smalltalk, 0-based in Go)
			index := pile.GetIntegerImmediate(args[0]) - 1

			// Get the value
			value := pile.GetIntegerImmediate(args[1])

			// Check bounds
			if index < 0 || int(index) >= byteArray.Size() {
				panic(fmt.Sprintf("ByteArray index out of bounds: %d", index+1))
			}

			// Check value range (0-255)
			if value < 0 || value > 255 {
				panic(fmt.Sprintf("ByteArray value out of range (0-255): %d", value))
			}

			// Set the byte at the given index
			byteArray.AtPut(int(index), byte(value))

			// Return the value
			return args[1]
		}
	case 60: // Class new - create a new instance of a class
		if receiver.Type() == pile.OBJ_CLASS {
			// Get the class
			class := pile.ObjectToClass(receiver)

			// Create a new instance of the class
			instance := pile.NewInstance(class)

			// We need to explicitly set the class of the instance
			instance.SetClass(receiver)

			return instance
		}
	default:
		panic("executePrimitive: unknown primitive index\n")
	}
	return nil // Fall through to method
}

// GetGlobals returns the globals map as a slice
func (vm *VM) GetGlobals() []*pile.Object {
	// Convert map to slice for memory management
	globals := make([]*pile.Object, 0, len(vm.Globals))
	for _, obj := range vm.Globals {
		globals = append(globals, obj)
	}
	return globals
}
