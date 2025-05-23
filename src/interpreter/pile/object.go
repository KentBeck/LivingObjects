package pile

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
	OBJ_BYTE_ARRAY
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
	case OBJ_BYTE_ARRAY:
		return "ByteArray"
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
	if o.Type() != OBJ_CLASS {
		panic("object is not a class")
	}
	class := (*Class)(unsafe.Pointer(o))
	return class.MethodDictionary
}

// Class represents a Smalltalk class object
type Class struct {
	Object
	Name             string
	SuperClass       *Object
	InstanceVarNames []string
	MethodDictionary *Object  // Direct reference to the method dictionary
}

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