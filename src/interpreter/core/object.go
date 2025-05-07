package core

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
	OBJ_EXCEPTION
)

// Object represents a Smalltalk object
type Object struct {
	TypeField          ObjectType
	ClassField         *Class
	MovedField         bool      // Used for garbage collection
	ForwardingPtrField *Object   // Used for garbage collection
	InstanceVarsField  []*Object // Instance variables stored by index
}

// ObjectInterface defines the interface for all Smalltalk objects
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
	GetInstanceVarByIndex(index int) *Object
	SetInstanceVarByIndex(index int, value *Object)
	IsTrue() bool // for now
	String() string
}

func (o *Object) Type() ObjectType {
	return o.TypeField
}

func (o *Object) SetType(t ObjectType) {
	o.TypeField = t
}

// Class returns the class of the object
func (o *Object) Class() *Object {
	return (*Object)(unsafe.Pointer(o.ClassField))
}

// SetClass sets the class of the object
func (o *Object) SetClass(class *Object) {
	o.ClassField = (*Class)(unsafe.Pointer(class))
}

// Moved returns whether the object has been moved during garbage collection
func (o *Object) Moved() bool {
	return o.MovedField
}

// SetMoved sets whether the object has been moved during garbage collection
func (o *Object) SetMoved(moved bool) {
	o.MovedField = moved
}

// ForwardingPtr returns the forwarding pointer of the object
func (o *Object) ForwardingPtr() *Object {
	return o.ForwardingPtrField
}

// SetForwardingPtr sets the forwarding pointer of the object
func (o *Object) SetForwardingPtr(ptr *Object) {
	o.ForwardingPtrField = ptr
}

// InstanceVars returns the instance variables of the object
func (o *Object) InstanceVars() []*Object {
	return o.InstanceVarsField
}

func (o *Object) GetInstanceVarByIndex(index int) *Object {
	if index < 0 || index >= len(o.InstanceVars()) {
		panic("index out of bounds")
	}

	return o.InstanceVars()[index]
}

func (o *Object) SetInstanceVarByIndex(index int, value *Object) {
	if index < 0 || index >= len(o.InstanceVars()) {
		panic("index out of bounds")
	}

	vars := o.InstanceVars()
	vars[index] = value
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

	// Handle regular objects based on type
	switch o.Type() {
	case OBJ_STRING:
		str := (*String)(unsafe.Pointer(o))
		return fmt.Sprintf("'%s'", str.Value)
	case OBJ_SYMBOL:
		sym := (*Symbol)(unsafe.Pointer(o))
		return fmt.Sprintf("#%s", sym.Value)
	case OBJ_CLASS:
		if o.Class() != nil && o.Class().Type() == OBJ_CLASS {
			class := (*Class)(unsafe.Pointer(o))
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
			class := (*Class)(unsafe.Pointer(o.Class()))
			return fmt.Sprintf("a %s", class.Name)
		}
		class := (*Class)(unsafe.Pointer(o))
		return fmt.Sprintf("Class %s", class.Name)
	case OBJ_ARRAY:
		array := (*Array)(unsafe.Pointer(o))
		return fmt.Sprintf("Array(%d)", len(array.Elements))
	case OBJ_DICTIONARY:
		dict := (*Dictionary)(unsafe.Pointer(o))
		return fmt.Sprintf("Dictionary(%d)", dict.GetEntryCount())
	case OBJ_BLOCK:
		return "Block"
	case OBJ_METHOD:
		method := (*Method)(unsafe.Pointer(o))
		if method != nil && method.Selector != nil {
			sym := (*Symbol)(unsafe.Pointer(method.Selector))
			return fmt.Sprintf("Method %s", sym.Value)
		}
		return "a Method"
	default:
		return "Unknown object"
	}
}

// GetMethodDict gets the method dictionary for a class
func (o *Object) GetMethodDict() *Object {
	if o.Type() != OBJ_CLASS || len(o.InstanceVars()) == 0 {
		panic("object is not a class or has no instance variables")
	}
	return o.InstanceVars()[METHOD_DICTIONARY_IV]
}

// Class represents a Smalltalk class object
type Class struct {
	Object
	Name             string
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

// Method represents a Smalltalk method
type Method struct {
	Object
	Bytecodes      []byte
	Literals       []*Object
	Selector       *Object
	TempVarNames   []string
	MethodClass    *Class
	IsPrimitive    bool
	PrimitiveIndex int
}

// Block represents a Smalltalk block
type Block struct {
	Object
	Bytecodes    []byte
	Literals     []*Object
	TempVarNames []string
	OuterContext interface{} // Changed to interface{} to avoid circular dependency
}

// Array represents a Smalltalk array object
type Array struct {
	Object
	Elements []*Object
}

// Dictionary represents a Smalltalk dictionary object
type Dictionary struct {
	Object
	Entries map[string]*Object // later Object->Object
}

// GetEntries returns the entries of the dictionary
func (d *Dictionary) GetEntries() map[string]*Object {
	return d.Entries
}

// GetEntry gets an entry from the dictionary
func (d *Dictionary) GetEntry(key string) *Object {
	return d.Entries[key]
}

// SetEntry sets an entry in the dictionary
func (d *Dictionary) SetEntry(key string, value *Object) {
	d.Entries[key] = value
}

// GetEntryCount returns the number of entries in the dictionary
func (d *Dictionary) GetEntryCount() int {
	return len(d.Entries)
}

const METHOD_DICTIONARY_IV = 0

// NewInstance creates a new instance of a class
func NewInstance(class *Class) *Object {
	// Initialize instance variables array with nil values
	instVarsSize := 0
	if class != nil && len(class.InstanceVarNames) > 0 {
		instVarsSize = len(class.InstanceVarNames)
	}
	instVars := make([]*Object, instVarsSize)
	for i := range instVars {
		instVars[i] = MakeNilImmediate()
	}

	obj := &Object{
		TypeField:         OBJ_INSTANCE,
		ClassField:        class,
		InstanceVarsField: instVars,
	}
	return obj
}

// NewString creates a new string object
func NewString(value string) *String {
	str := &String{
		Object: Object{
			TypeField: OBJ_STRING,
		},
		Value: value,
	}
	return str
}

// NewSymbol creates a new symbol object
func NewSymbol(value string) *Object {
	sym := &Symbol{
		Object: Object{
			TypeField: OBJ_SYMBOL,
		},
		Value: value,
	}
	return (*Object)(unsafe.Pointer(sym))
}

// NewArray creates a new array object
func NewArray(size int) *Array {
	obj := &Array{
		Object: Object{
			TypeField: OBJ_ARRAY,
		},
		Elements: make([]*Object, size),
	}
	return obj
}

// NewDictionary creates a new dictionary object
func NewDictionary() *Object {
	entries := make(map[string]*Object)
	obj := &Dictionary{
		Object: Object{
			TypeField: OBJ_DICTIONARY,
		},
		Entries: entries,
	}
	return (*Object)(unsafe.Pointer(obj))
}

// NewClass creates a new class object
func NewClass(name string, superClass *Class) *Class {
	// For classes, we need a special instance variable for the method dictionary
	// We'll store it at index 0
	instVars := make([]*Object, 1)
	instVars[0] = NewDictionary() // methodDict at index 0

	result := &Class{
		Object: Object{
			TypeField:         OBJ_CLASS,
			InstanceVarsField: instVars,
		},
		SuperClass:       (*Object)(unsafe.Pointer(superClass)),
		InstanceVarNames: make([]string, 0),
		Name:             name,
	}

	return result
}

// NewMethod creates a new method object
func NewMethod(selector *Object, class *Class) *Object {
	method := &Method{
		Object: Object{
			TypeField: OBJ_METHOD,
		},
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*Object, 0),
		Selector:     selector,
		MethodClass:  class,
		TempVarNames: make([]string, 0),
	}

	return (*Object)(unsafe.Pointer(method))
}

// NewBlock creates a new block object
func NewBlock(outerContext interface{}) *Object {
	block := &Block{
		Object: Object{
			TypeField: OBJ_BLOCK,
		},
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*Object, 0),
		TempVarNames: make([]string, 0),
		OuterContext: outerContext,
	}

	return (*Object)(unsafe.Pointer(block))
}
