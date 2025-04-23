package main

import (
	"fmt"
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
	BooleanValue     bool
	StringValue      string
	SymbolValue      string
	InstanceVars     []*Object // Instance variables stored by index
	Elements         []*Object
	Entries          map[string]*Object
	Method           *Method
	Bytecodes        []byte
	Literals         []*Object
	Selector         *Object
	SuperClass       *Object
	InstanceVarNames []string
	Moved            bool    // Used for garbage collection
	ForwardingPtr    *Object // Used for garbage collection
}

const METHOD_DICTIONARY_IV = 0

// Method represents a Smalltalk method
type Method struct {
	Bytecodes      []byte
	Literals       []*Object
	Selector       *Object
	Class          *Object
	TempVarNames   []string
	IsPrimitive    bool
	PrimitiveIndex int
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
func NewString(value string) *Object {
	return &Object{
		Type:        OBJ_STRING,
		StringValue: value,
	}
}

// NewSymbol creates a new symbol object
func NewSymbol(value string) *Object {
	return &Object{
		Type:        OBJ_SYMBOL,
		SymbolValue: value,
	}
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

	return &Object{
		Type:             OBJ_CLASS,
		SymbolValue:      name,
		SuperClass:       superClass,
		InstanceVarNames: make([]string, 0),
		InstanceVars:     instVars,
	}
}

// NewMethod creates a new method object
func NewMethod(selector *Object, class *Object) *Object {
	method := &Method{
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*Object, 0),
		Selector:     selector,
		Class:        class,
		TempVarNames: make([]string, 0),
	}

	return &Object{
		Type:   OBJ_METHOD,
		Method: method,
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
		if IsIntegerImmediate(o) {
			return true
		}
		// Immediate float is true
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
		if o.BooleanValue {
			return "true"
		}
		return "false"
	case OBJ_NIL:
		return "nil"
	case OBJ_STRING:
		return fmt.Sprintf("'%s'", o.StringValue)
	case OBJ_SYMBOL:
		return fmt.Sprintf("#%s", o.SymbolValue)
	case OBJ_ARRAY:
		return fmt.Sprintf("Array(%d)", len(o.Elements))
	case OBJ_DICTIONARY:
		return fmt.Sprintf("Dictionary(%d)", len(o.Entries))
	case OBJ_INSTANCE:
		if o.Class != nil {
			return fmt.Sprintf("a %s", o.Class.SymbolValue)
		}
		return "an Object"
	case OBJ_CLASS:
		return fmt.Sprintf("Class %s", o.SymbolValue)
	case OBJ_METHOD:
		if o.Method.Selector != nil {
			return fmt.Sprintf("Method %s", o.Method.Selector.SymbolValue)
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
