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
	Type           ObjectType
	Class          *Object
	IntegerValue   int64
	BooleanValue   bool
	StringValue    string
	SymbolValue    string
	InstanceVars   map[string]*Object
	Elements       []*Object
	Entries        map[string]*Object
	Method         *Method
	Bytecodes      []byte
	Literals       []*Object
	Selector       *Object
	SuperClass     *Object
	InstanceVarNames []string
	Moved          bool  // Used for garbage collection
	ForwardingPtr  *Object // Used for garbage collection
}

// Method represents a Smalltalk method
type Method struct {
	Bytecodes      []byte
	Literals       []*Object
	Selector       *Object
	Class          *Object
	TempVarNames   []string
}

// NewInteger creates a new integer object
func NewInteger(value int64) *Object {
	return &Object{
		Type:         OBJ_INTEGER,
		IntegerValue: value,
	}
}

// NewBoolean creates a new boolean object
func NewBoolean(value bool) *Object {
	return &Object{
		Type:         OBJ_BOOLEAN,
		BooleanValue: value,
	}
}

// NewNil creates a new nil object
func NewNil() *Object {
	return &Object{
		Type: OBJ_NIL,
	}
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
	return &Object{
		Type:         OBJ_INSTANCE,
		Class:        class,
		InstanceVars: make(map[string]*Object),
	}
}

// NewClass creates a new class object
func NewClass(name string, superClass *Object) *Object {
	return &Object{
		Type:            OBJ_CLASS,
		SymbolValue:     name,
		SuperClass:      superClass,
		InstanceVarNames: make([]string, 0),
		InstanceVars:    make(map[string]*Object),
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
	if o.Type == OBJ_BOOLEAN {
		return o.BooleanValue
	}
	return o.Type != OBJ_NIL
}

// String returns a string representation of the object
func (o *Object) String() string {
	switch o.Type {
	case OBJ_INTEGER:
		return fmt.Sprintf("%d", o.IntegerValue)
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
