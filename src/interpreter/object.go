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
	type1            ObjectType
	class            *Object
	moved            bool      // Used for garbage collection
	forwardingPtr    *Object   // Used for garbage collection
	instanceVars     []*Object // Instance variables stored by index
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

type ObjectInterface interface {
	Type() ObjectType
	SetType(t ObjectType)
	Class() *Object
	SetClass(class *Object)
	Moved() bool
	SetMoved(moved bool)
	ForwardingPtr() *Object
	SetForwardingPtr(ptr *Object)
	InstanceVars() []*Object
	SetInstanceVars(vars []*Object)
	IsTrue() bool // for now
	String() string
}

func (o *Object) Type() ObjectType {
	return o.type1
}

func (o *Object) SetType(t ObjectType) {
	o.type1 = t
}

// Class returns the class of the object
func (o *Object) Class() *Object {
	return o.class
}

// SetClass sets the class of the object
func (o *Object) SetClass(class *Object) {
	o.class = class
}

// Moved returns whether the object has been moved during garbage collection
func (o *Object) Moved() bool {
	return o.moved
}

// SetMoved sets whether the object has been moved during garbage collection
func (o *Object) SetMoved(moved bool) {
	o.moved = moved
}

// ForwardingPtr returns the forwarding pointer of the object
func (o *Object) ForwardingPtr() *Object {
	return o.forwardingPtr
}

// SetForwardingPtr sets the forwarding pointer of the object
func (o *Object) SetForwardingPtr(ptr *Object) {
	o.forwardingPtr = ptr
}

// InstanceVars returns the instance variables of the object
func (o *Object) InstanceVars() []*Object {
	return o.instanceVars
}

// SetInstanceVars sets the instance variables of the object
func (o *Object) SetInstanceVars(vars []*Object) {
	o.instanceVars = vars
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
func NewNil() ObjectInterface {
	return MakeNilImmediate()
}

// NewString creates a new string object
func NewString(value string) *String {
	str := &String{
		Value: value,
	}
	str.type1 = OBJ_STRING
	return str
}

// NewSymbol creates a new symbol object
func NewSymbol(value string) *Object {
	sym := &Symbol{
		Value: value,
	}
	sym.type1 = OBJ_SYMBOL
	return SymbolToObject(sym)
}

// NewArray creates a new array object
func NewArray(size int) *Object {
	obj := &Object{
		type1:    OBJ_ARRAY,
		Elements: make([]*Object, size),
	}
	return obj
}

// NewDictionary creates a new dictionary object
func NewDictionary() *Object {
	obj := &Object{
		type1:   OBJ_DICTIONARY,
		Entries: make(map[string]*Object),
	}
	return obj
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
		instVars[i] = NewNil().(*Object)
	}

	obj := &Object{
		type1:        OBJ_INSTANCE,
		class:        class,
		instanceVars: instVars,
	}
	return obj
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
	sym.type1 = OBJ_CLASS
	sym.SuperClass = superClass
	sym.InstanceVarNames = make([]string, 0)
	sym.SetInstanceVars(instVars)

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

	obj := &Object{
		type1:  OBJ_METHOD,
		Method: method,
	}
	return obj
}

// NewBlock creates a new block object
func NewBlock(outerContext *Context) *Object {
	block := &Block{
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*Object, 0),
		TempVarNames: make([]string, 0),
		OuterContext: outerContext,
	}

	obj := &Object{
		type1: OBJ_BLOCK,
		Block: block,
	}
	return obj
}

// IsTrue returns true if the object is considered true in Smalltalk
func (o *Object) IsTrue() bool {
	// maybe should be an error if it's not a boolean?
	return IsTrueImmediate(o)
}

// String returns a string representation of the object
func (o *Object) String() string {
	// Check if it's an immediate value first
	if IsImmediate(o) {
		if IsNilImmediate(o) {
			return "nil"
		}
		if IsTrueImmediate(o) {
			return "true"
		}
		if IsFalseImmediate(o) {
			return "false"
		}
		if IsIntegerImmediate(o) {
			return fmt.Sprintf("%d", GetIntegerImmediate(o))
		}
		if IsFloatImmediate(o) {
			return fmt.Sprintf("%g", GetFloatImmediate(o))
		}
		return "Immediate value"
	}

	// Handle regular objects
	switch o.Type() {
	case OBJ_STRING:
		str := ObjectToString(o)
		return fmt.Sprintf("'%s'", str.Value)
	case OBJ_SYMBOL:
		sym := ObjectToSymbol(o)
		return fmt.Sprintf("#%s", sym.Value)
	case OBJ_CLASS:
		if o.Class() != nil && o.Class().Type() == OBJ_CLASS {
			class := ObjectToClass(o)
			if class.Name != "" {
				return fmt.Sprintf("Class %s", class.Name)
			}
		}
		return "Class Object"
	case OBJ_INSTANCE:
		if o.Class() == nil {
			return "an Object"
		}
		if o.Type() != OBJ_CLASS || len(o.InstanceVars()) == 0 {
			class := ObjectToClass(o.Class())
			return fmt.Sprintf("a %s", class.Name)
		}
		class := ObjectToClass(o)
		return fmt.Sprintf("Class %s", class.Name)
	case OBJ_ARRAY:
		return fmt.Sprintf("Array(%d)", len(o.Elements))
	case OBJ_DICTIONARY:
		return fmt.Sprintf("Dictionary(%d)", len(o.Entries))
	case OBJ_BLOCK:
		return "Block"
	case OBJ_METHOD:
		if o.Method != nil && o.Method.Selector != nil {
			return fmt.Sprintf("Method %s", GetSymbolValue(o.Method.Selector))
		}
		return "a Method"
	default:
		return "Unknown object"
	}
}

// GetInstanceVarByIndex gets an instance variable by index
func (o *Object) GetInstanceVarByIndex(index int) *Object {
	if index < 0 || index >= len(o.InstanceVars()) {
		panic("index out of bounds")
	}

	return o.InstanceVars()[index]
}

// SetInstanceVarByIndex sets an instance variable by index
func (o *Object) SetInstanceVarByIndex(index int, value *Object) {
	if index < 0 || index >= len(o.InstanceVars()) {
		panic("index out of bounds")
	}

	vars := o.InstanceVars()
	vars[index] = value
}

// GetMethodDict gets the method dictionary for a class
func (o *Object) GetMethodDict() *Object {
	if o.Type() != OBJ_CLASS || len(o.InstanceVars()) == 0 {
		panic("object is not a class or has no instance variables")
	}

	// Method dictionary is stored at index 0 for classes
	return o.InstanceVars()[METHOD_DICTIONARY_IV]
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
	if o.Type() == OBJ_SYMBOL {
		return ObjectToSymbol(o).Value
	}
	panic("GetSymbolValue: not a symbol")
}
