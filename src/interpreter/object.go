package main

import (
	"fmt"
	"unsafe"
)

// ObjectType represents the type of a Smalltalk object
type ObjectType int

const (
	OBJ_INTEGER ObjectType = iota
	OBJ_BOOLEAN
	OBJ_NIL
	OBJ_STRING
	OBJ_ARRAY
	OBJ_DICTIONARY
	OBJ_BLOCK
	OBJ_INSTANCE
	OBJ_CLASS
	OBJ_METHOD
	OBJ_SYMBOL
)

// Object represents a Smalltalk object
type Object struct {
	Type             ObjectType
	Class            *Object
	Moved            bool      // Used for garbage collection
	ForwardingPtr    *Object   // Used for garbage collection
	InstanceVars     []*Object // Instance variables stored by index
	Elements         []*Object
	Entries          map[string]*Object
	Method           *Method
	Block            *Block
	Bytecodes        []byte
	Literals         []*Object
	Selector         *Object
	SuperClass       *Object
	InstanceVarNames []string
}

// String represents a Smalltalk string object
type String struct {
	Object
	Value string
}

// Symbol represents a Smalltalk symbol object
type Symbol struct {
	Object
	Value string
}

// Class represents a Smalltalk class object
type Class struct {
	Object
	Name             string
	SuperClass       *Object
	InstanceVarNames []string
}

const METHOD_DICTIONARY_IV = 0

// Method represents a Smalltalk method
type Method struct {
	Object
	Bytecodes      []byte
	Literals       []*Object
	Selector       *Object
	MethodClass    *Object
	TempVarNames   []string
	IsPrimitive    bool
	PrimitiveIndex int
}

// Block represents a Smalltalk block
type Block struct {
	Object
	Bytecodes    []byte
	Literals     []*Object
	TempVarNames []string
	OuterContext *Context
}

// NewBoolean creates a new boolean object
// This now returns an immediate value
func NewBoolean(value bool) *Object {
	if value {
		return MakeTrueImmediate()
	} else {
		return MakeFalseImmediate()
	}
}

// NewNil creates a new nil object
// This now returns an immediate nil value
func NewNil() *Object {
	return MakeNilImmediate()
}

// NewString creates a new string object
func NewString(value string) *String {
	str := &String{
		Value: value,
	}
	str.Type = OBJ_STRING
	return str
}

// NewSymbol creates a new symbol object
func NewSymbol(value string) *Object {
	sym := &Symbol{
		Value: value,
	}
	sym.Type = OBJ_SYMBOL
	return SymbolToObject(sym)
}

// NewArray creates a new array object
func NewArray(size int) *Object {
	return &Object{
		Type:     OBJ_ARRAY,
		Elements: make([]*Object, size),
	}
}

// NewDictionary creates a new dictionary object
func NewDictionary() *Object {
	return &Object{
		Type:    OBJ_DICTIONARY,
		Entries: make(map[string]*Object),
	}
}

// NewInstance creates a new instance of a class
func NewInstance(class *Object) *Object {
	// Initialize instance variables array with nil values
	instVarsSize := 0
	if class != nil && len(class.InstanceVarNames) > 0 {
		instVarsSize = len(class.InstanceVarNames)
	}
	instVars := make([]*Object, instVarsSize)
	for i := range instVars {
		instVars[i] = NewNil()
	}

	return &Object{
		Type:         OBJ_INSTANCE,
		Class:        class,
		InstanceVars: instVars,
	}
}

// NewClass creates a new class object
func NewClass(name string, superClass *Object) *Object {
	// For classes, we need a special instance variable for the method dictionary
	// We'll store it at index 0
	instVars := make([]*Object, 1)
	instVars[0] = NewDictionary() // methodDict at index 0

	// Create a Symbol for the class name
	sym := &Symbol{
		Value: name,
	}

	// Set up the Symbol as a class
	sym.Type = OBJ_CLASS
	sym.SuperClass = superClass
	sym.InstanceVarNames = make([]string, 0)
	sym.InstanceVars = instVars

	return SymbolToObject(sym)
}

// NewMethod creates a new method object
func NewMethod(selector *Object, class *Object) *Object {
	method := &Method{
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*Object, 0),
		Selector:     selector,
		MethodClass:  class,
		TempVarNames: make([]string, 0),
	}

	return &Object{
		Type:   OBJ_METHOD,
		Method: method,
	}
}

// NewBlock creates a new block object
func NewBlock(outerContext *Context) *Object {
	block := &Block{
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*Object, 0),
		TempVarNames: make([]string, 0),
		OuterContext: outerContext,
	}

	return &Object{
		Type:  OBJ_BLOCK,
		Block: block,
	}
}

// IsTrue returns true if the object is considered true in Smalltalk
func (o *Object) IsTrue() bool {
	// Check if it's an immediate value
	if IsImmediate(o) {
		// Immediate nil is false
		if IsNilImmediate(o) {
			return false
		}
		// Immediate true is true
		if IsTrueImmediate(o) {
			return true
		}
		// Immediate false is false
		if IsFalseImmediate(o) {
			return false
		}
		// Immediate integer is true
		// Sus
		if IsIntegerImmediate(o) {
			return true
		}
		// Immediate float is true
		// Sus
		if IsFloatImmediate(o) {
			return true
		}
		// Other immediate types will be added later
		return true
	}

	// Not thrilled about this. We don't want truthiness in the image, probably not here either
	return false
}

// String returns a string representation of the object
func (o *Object) String() string {
	// Check if it's an immediate value
	if IsImmediate(o) {
		// Immediate nil
		if IsNilImmediate(o) {
			return "nil"
		}
		// Immediate true
		if IsTrueImmediate(o) {
			return "true"
		}
		// Immediate false
		if IsFalseImmediate(o) {
			return "false"
		}
		// Immediate integer
		if IsIntegerImmediate(o) {
			value := GetIntegerImmediate(o)
			return fmt.Sprintf("%d", value)
		}
		// Immediate float
		if IsFloatImmediate(o) {
			value := GetFloatImmediate(o)
			return fmt.Sprintf("%g", value)
		}
		// Other immediate types will be added later
		return "Immediate value"
	}

	// Regular objects
	switch o.Type {
	case OBJ_INTEGER:
		panic("Non-immediate integer encountered")
	case OBJ_BOOLEAN:
		panic("Non-immediate boolean encountered")
	case OBJ_NIL:
		return "nil"
	case OBJ_STRING:
		// Convert to String type to access the Value field
		str := ObjectToString(o)
		return fmt.Sprintf("'%s'", str.Value)
	case OBJ_SYMBOL:
		// Convert to Symbol type to access the Value field
		sym := ObjectToSymbol(o)
		return fmt.Sprintf("#%s", sym.Value)
	case OBJ_ARRAY:
		return fmt.Sprintf("Array(%d)", len(o.Elements))
	case OBJ_DICTIONARY:
		return fmt.Sprintf("Dictionary(%d)", len(o.Entries))
	case OBJ_INSTANCE:
		if o.Class != nil {
			// For classes, the name is stored in the class itself
			// The class is a Symbol object
			if o.Class.Type == OBJ_CLASS {
				// Get the class name from the class object
				classSymbol := ObjectToSymbol(o.Class)
				return fmt.Sprintf("a %s", classSymbol.Value)
			} else {
				// If the class is not a class object, just use its string representation
				return fmt.Sprintf("an instance of %s", o.Class.String())
			}
		}
		return "an Object"
	case OBJ_CLASS:
		// For classes, the name is stored in a Symbol object
		classSymbol := ObjectToSymbol(o)
		return fmt.Sprintf("Class %s", classSymbol.Value)
	case OBJ_METHOD:
		if o.Method.Selector != nil {
			selectorSymbol := ObjectToSymbol(o.Method.Selector)
			return fmt.Sprintf("Method %s", selectorSymbol.Value)
		}
		return "a Method"
	default:
		return "Unknown object"
	}
}

// GetInstanceVarByIndex gets an instance variable by index
func (o *Object) GetInstanceVarByIndex(index int) *Object {
	if index < 0 || index >= len(o.InstanceVars) {
		panic("index out of bounds")
	}

	return o.InstanceVars[index]
}

// SetInstanceVarByIndex sets an instance variable by index
func (o *Object) SetInstanceVarByIndex(index int, value *Object) {
	if index < 0 || index >= len(o.InstanceVars) {
		panic("index out of bounds")
	}

	o.InstanceVars[index] = value
}

// GetMethodDict gets the method dictionary for a class
func (o *Object) GetMethodDict() *Object {
	if o.Type != OBJ_CLASS || len(o.InstanceVars) == 0 {
		panic("object is not a class or has no instance variables")
	}

	// Method dictionary is stored at index 0 for classes
	return o.InstanceVars[METHOD_DICTIONARY_IV]
}

// StringToObject converts a String to an Object
func StringToObject(s *String) *Object {
	return (*Object)(unsafe.Pointer(s))
}

// ObjectToString converts an Object to a String
func ObjectToString(o *Object) *String {
	return (*String)(unsafe.Pointer(o))
}

// SymbolToObject converts a Symbol to an Object
func SymbolToObject(s *Symbol) *Object {
	return (*Object)(unsafe.Pointer(s))
}

// ObjectToSymbol converts an Object to a Symbol
func ObjectToSymbol(o *Object) *Symbol {
	return (*Symbol)(unsafe.Pointer(o))
}

// GetSymbolValue gets the value from a Symbol object
func GetSymbolValue(o *Object) string {
	if o.Type == OBJ_SYMBOL {
		return ObjectToSymbol(o).Value
	}
	panic("GetSymbolValue: not a symbol")
}
