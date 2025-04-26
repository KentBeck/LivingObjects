package main

import (
	"fmt"
)

// VMWithRawMemory represents the Smalltalk virtual machine using raw memory allocation
type VMWithRawMemory struct {
	Globals        map[string]*Object
	CurrentContext *Context
	MemoryManager  *RawMemoryManager

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

// NewVMWithRawMemory creates a new virtual machine with raw memory allocation
func NewVMWithRawMemory() (*VMWithRawMemory, error) {
	// Create a new raw memory manager
	memoryManager, err := NewRawMemoryManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create raw memory manager: %w", err)
	}

	vm := &VMWithRawMemory{
		Globals:       make(map[string]*Object),
		MemoryManager: memoryManager,
	}

	// Initialize special objects
	vm.ObjectClass = vm.NewObjectClass()
	vm.NilClass = vm.NewClass("UndefinedObject", vm.ObjectClass)
	// Use immediate nil value
	vm.NilObject = MakeNilImmediate()
	vm.TrueClass = vm.NewClass("True", vm.ObjectClass)
	// Use immediate true value
	vm.TrueObject = MakeTrueImmediate()
	vm.FalseClass = vm.NewClass("False", vm.ObjectClass)
	// Use immediate false value
	vm.FalseObject = MakeFalseImmediate()
	vm.IntegerClass = vm.NewIntegerClass()
	vm.FloatClass = vm.NewFloatClass()

	return vm, nil
}

// Close releases all allocated memory
func (vm *VMWithRawMemory) Close() error {
	if vm.MemoryManager != nil {
		return vm.MemoryManager.Close()
	}
	return nil
}

// NewObjectClass creates a new Object class
func (vm *VMWithRawMemory) NewObjectClass() *Object {
	// Allocate a new class
	class := vm.MemoryManager.AllocateClass()
	if class == nil {
		// Need to collect garbage
		vm.MemoryManager.Collect((*VM)(nil)) // No VM yet, so pass nil
		class = vm.MemoryManager.AllocateClass()
		if class == nil {
			panic("Failed to allocate Object class")
		}
	}

	// Initialize the class
	class.Type = OBJ_CLASS
	class.Name = "Object"
	class.SuperClass = nil
	class.InstanceVarNames = make([]string, 0)

	// Create a method dictionary
	methodDict := vm.NewDictionary()
	class.InstanceVars = make([]*Object, 1)
	class.InstanceVars[METHOD_DICTIONARY_IV] = methodDict

	// Add basicClass method to Object class using MethodBuilder
	NewMethodBuilder(ClassToObject(class)).
		Selector("basicClass").
		Primitive(5). // basicClass primitive
		Go()

	return ClassToObject(class)
}

// NewClass creates a new class
func (vm *VMWithRawMemory) NewClass(name string, superClass *Object) *Object {
	// Allocate a new class
	class := vm.MemoryManager.AllocateClass()
	if class == nil {
		// Need to collect garbage
		vm.MemoryManager.Collect((*VM)(nil)) // No VM yet, so pass nil
		class = vm.MemoryManager.AllocateClass()
		if class == nil {
			panic("Failed to allocate class")
		}
	}

	// Initialize the class
	class.Type = OBJ_CLASS
	class.Name = name
	class.SuperClass = superClass
	class.InstanceVarNames = make([]string, 0)

	// Create a method dictionary
	methodDict := vm.NewDictionary()
	class.InstanceVars = make([]*Object, 1)
	class.InstanceVars[METHOD_DICTIONARY_IV] = methodDict

	return ClassToObject(class)
}

// NewIntegerClass creates a new Integer class
func (vm *VMWithRawMemory) NewIntegerClass() *Object {
	result := vm.NewClass("Integer", vm.ObjectClass)

	// Add primitive methods to the Integer class using MethodBuilder
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

// NewFloatClass creates a new Float class
func (vm *VMWithRawMemory) NewFloatClass() *Object {
	result := vm.NewClass("Float", vm.ObjectClass)

	// Add primitive methods to the Float class using MethodBuilder
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
func (vm *VMWithRawMemory) NewInteger(value int64) *Object {
	// Check if the value fits in 62 bits
	if value <= 0x1FFFFFFFFFFFFFFF && value >= -0x2000000000000000 {
		// Use immediate integer
		return MakeIntegerImmediate(value)
	}

	// Panic for large values that don't fit in 62 bits
	panic("Integer value too large for immediate representation")
}

// NewFloat creates a new float object
func (vm *VMWithRawMemory) NewFloat(value float64) *Object {
	return MakeFloatImmediate(value)
}

// NewString creates a new string object
func (vm *VMWithRawMemory) NewString(value string) *Object {
	// Allocate a new string
	str := vm.MemoryManager.AllocateString()
	if str == nil {
		// Need to collect garbage
		vm.MemoryManager.Collect((*VM)(nil)) // No VM yet, so pass nil
		str = vm.MemoryManager.AllocateString()
		if str == nil {
			panic("Failed to allocate string")
		}
	}

	// Initialize the string
	str.Type = OBJ_STRING
	str.Value = value

	return StringToObject(str)
}

// NewSymbol creates a new symbol object
func (vm *VMWithRawMemory) NewSymbol(value string) *Object {
	// Allocate a new symbol
	sym := vm.MemoryManager.AllocateSymbol()
	if sym == nil {
		// Need to collect garbage
		vm.MemoryManager.Collect((*VM)(nil)) // No VM yet, so pass nil
		sym = vm.MemoryManager.AllocateSymbol()
		if sym == nil {
			panic("Failed to allocate symbol")
		}
	}

	// Initialize the symbol
	sym.Type = OBJ_SYMBOL
	sym.Value = value

	return SymbolToObject(sym)
}

// NewArray creates a new array object
func (vm *VMWithRawMemory) NewArray(size int) *Object {
	// Allocate a new object
	obj := vm.MemoryManager.AllocateObject()
	if obj == nil {
		// Need to collect garbage
		vm.MemoryManager.Collect((*VM)(nil)) // No VM yet, so pass nil
		obj = vm.MemoryManager.AllocateObject()
		if obj == nil {
			panic("Failed to allocate array")
		}
	}

	// Allocate the elements array
	elements := vm.MemoryManager.AllocateObjectArray(size)
	if elements == nil {
		// Need to collect garbage
		vm.MemoryManager.Collect((*VM)(nil)) // No VM yet, so pass nil
		elements = vm.MemoryManager.AllocateObjectArray(size)
		if elements == nil {
			panic("Failed to allocate array elements")
		}
	}

	// Initialize the array
	obj.Type = OBJ_ARRAY
	obj.Elements = elements

	return obj
}

// NewDictionary creates a new dictionary object
func (vm *VMWithRawMemory) NewDictionary() *Object {
	// Allocate a new object
	obj := vm.MemoryManager.AllocateObject()
	if obj == nil {
		// Need to collect garbage
		vm.MemoryManager.Collect((*VM)(nil)) // No VM yet, so pass nil
		obj = vm.MemoryManager.AllocateObject()
		if obj == nil {
			panic("Failed to allocate dictionary")
		}
	}

	// Allocate the entries map
	entries := vm.MemoryManager.AllocateStringMap()
	if entries == nil {
		// Need to collect garbage
		vm.MemoryManager.Collect((*VM)(nil)) // No VM yet, so pass nil
		entries = vm.MemoryManager.AllocateStringMap()
		if entries == nil {
			panic("Failed to allocate dictionary entries")
		}
	}

	// Initialize the dictionary
	obj.Type = OBJ_DICTIONARY
	obj.Entries = entries

	return obj
}

// NewInstance creates a new instance of a class
func (vm *VMWithRawMemory) NewInstance(class *Object) *Object {
	// Allocate a new object
	obj := vm.MemoryManager.AllocateObject()
	if obj == nil {
		// Need to collect garbage
		vm.MemoryManager.Collect((*VM)(nil)) // No VM yet, so pass nil
		obj = vm.MemoryManager.AllocateObject()
		if obj == nil {
			panic("Failed to allocate instance")
		}
	}

	// Initialize instance variables array with nil values
	instVarsSize := 0
	if class != nil && len(class.InstanceVarNames) > 0 {
		instVarsSize = len(class.InstanceVarNames)
	}

	// Allocate the instance variables array
	instVars := vm.MemoryManager.AllocateObjectArray(instVarsSize)
	if instVars == nil {
		// Need to collect garbage
		vm.MemoryManager.Collect((*VM)(nil)) // No VM yet, so pass nil
		instVars = vm.MemoryManager.AllocateObjectArray(instVarsSize)
		if instVars == nil {
			panic("Failed to allocate instance variables")
		}
	}

	// Initialize the instance variables with nil
	for i := range instVars {
		instVars[i] = vm.NilObject
	}

	// Initialize the instance
	obj.Type = OBJ_INSTANCE
	obj.Class = class
	obj.InstanceVars = instVars

	return obj
}

// NewMethod creates a new method object
func (vm *VMWithRawMemory) NewMethod(selector *Object, class *Object) *Object {
	// Allocate a new method
	method := vm.MemoryManager.AllocateMethod()
	if method == nil {
		// Need to collect garbage
		vm.MemoryManager.Collect((*VM)(nil)) // No VM yet, so pass nil
		method = vm.MemoryManager.AllocateMethod()
		if method == nil {
			panic("Failed to allocate method")
		}
	}

	// Initialize the method
	method.Type = OBJ_METHOD
	method.Bytecodes = make([]byte, 0)
	method.Literals = make([]*Object, 0)
	method.Selector = selector
	method.MethodClass = class
	method.TempVarNames = make([]string, 0)

	return MethodToObject(method)
}

// GetClass returns the class of an object
func (vm *VMWithRawMemory) GetClass(obj *Object) *Object {
	// Check if it's an immediate value
	if IsImmediate(obj) {
		// Get the tag
		tag := GetTag(obj)
		switch tag {
		case TAG_INTEGER:
			return vm.IntegerClass
		case TAG_FLOAT:
			return vm.FloatClass
		case TAG_SPECIAL:
			// Check which special value it is
			if IsNilImmediate(obj) {
				return vm.NilClass
			} else if IsTrueImmediate(obj) {
				return vm.TrueClass
			} else if IsFalseImmediate(obj) {
				return vm.FalseClass
			}
		}
		// Unknown immediate value
		panic("Unknown immediate value")
	}

	// Non-immediate object
	if obj.Class == nil {
		panic("Object has no class")
	}
	return obj.Class
}

// LoadImage loads a Smalltalk image from a file
func (vm *VMWithRawMemory) LoadImage(path string) error {
	vm.Globals["Object"] = vm.ObjectClass
	return nil
}

// Execute executes the current context
func (vm *VMWithRawMemory) Execute() (*Object, error) {
	var finalResult *Object

	for vm.CurrentContext != nil {
		// Check if we need to collect garbage
		if vm.MemoryManager.ShouldCollect() {
			vm.MemoryManager.Collect((*VM)(nil)) // TODO: Convert VM to VMWithRawMemory
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
func (vm *VMWithRawMemory) ExecuteContext(context *Context) (*Object, error) {
	// This would be a copy of the ExecuteContext method from VM
	// For brevity, we'll omit the implementation here
	// In a real implementation, you would copy the method from VM
	// and update it to use the raw memory manager
	return nil, fmt.Errorf("not implemented")
}
