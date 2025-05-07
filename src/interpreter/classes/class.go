package classes

import (
	"fmt"
	"unsafe"

	"smalltalklsp/interpreter/core"
)

// Class represents a Smalltalk class object
type Class struct {
	core.Object
	Name             string
	SuperClass       *core.Object
	InstanceVarNames []string
}

const METHOD_DICTIONARY_IV = 0

// NewClass creates a new class object
func NewClass(name string, superClass *Class) *Class {
	// For classes, we need a special instance variable for the method dictionary
	// We'll store it at index 0
	instVars := make([]*core.Object, 1)
	instVars[0] = NewDictionary() // methodDict at index 0

	result := &Class{
		Object: core.Object{
			TypeField:         core.OBJ_CLASS,
			InstanceVarsField: instVars,
		},
		SuperClass:       ClassToObject(superClass),
		InstanceVarNames: make([]string, 0),
		Name:             name,
	}

	return result
}

// ClassToObject converts a Class to an Object
func ClassToObject(c *Class) *core.Object {
	return (*core.Object)(unsafe.Pointer(c))
}

// ObjectToClass converts an Object to a Class
func ObjectToClass(o *core.Object) *Class {
	return (*Class)(unsafe.Pointer(o))
}

// String returns a string representation of the class object
func (c *Class) String() string {
	return fmt.Sprintf("Class %s", c.Name)
}

// GetName returns the name of the class
func (c *Class) GetName() string {
	return c.Name
}

// SetName sets the name of the class
func (c *Class) SetName(name string) {
	c.Name = name
}

// GetSuperClass returns the superclass of the class
func (c *Class) GetSuperClass() *core.Object {
	return c.SuperClass
}

// SetSuperClass sets the superclass of the class
func (c *Class) SetSuperClass(superClass *core.Object) {
	c.SuperClass = superClass
}

// GetInstanceVarNames returns the instance variable names of the class
func (c *Class) GetInstanceVarNames() []string {
	return c.InstanceVarNames
}

// AddInstanceVarName adds an instance variable name to the class
func (c *Class) AddInstanceVarName(name string) {
	c.InstanceVarNames = append(c.InstanceVarNames, name)
}

// GetMethodDictionary returns the method dictionary of the class
func (c *Class) GetMethodDictionary() *Dictionary {
	methodDict := c.InstanceVars()[METHOD_DICTIONARY_IV]
	return ObjectToDictionary(methodDict)
}

// AddMethod adds a method to the class
func (c *Class) AddMethod(selector *core.Object, method *core.Object) {
	methodDict := c.GetMethodDictionary()
	selectorSymbol := ObjectToSymbol(selector)
	methodDict.SetEntry(selectorSymbol.Value, method)
}

// LookupMethod looks up a method in the class hierarchy
func (c *Class) LookupMethod(selector *core.Object) *core.Object {
	// Look for the method in this class
	methodDict := c.GetMethodDictionary()
	selectorSymbol := ObjectToSymbol(selector)
	method := methodDict.GetEntry(selectorSymbol.Value)
	if method != nil {
		return method
	}

	// If not found, look in the superclass
	if c.SuperClass != nil {
		superClass := ObjectToClass(c.SuperClass)
		return superClass.LookupMethod(selector)
	}

	// Method not found
	return nil
}

// NewInstance creates a new instance of the class
func (c *Class) NewInstance() *core.Object {
	// Initialize instance variables array with nil values
	instVarsSize := 0
	if c != nil && len(c.InstanceVarNames) > 0 {
		instVarsSize = len(c.InstanceVarNames)
	}
	instVars := make([]*core.Object, instVarsSize)
	for i := range instVars {
		instVars[i] = core.MakeNilImmediate()
	}

	obj := &core.Object{
		TypeField:         core.OBJ_INSTANCE,
		ClassField:        (*core.Class)(unsafe.Pointer(c)),
		InstanceVarsField: instVars,
	}
	return obj
}

// GetClassName gets the name of a class
func GetClassName(obj core.ObjectInterface) string {
	if obj.Type() != core.OBJ_CLASS {
		return ""
	}
	class := ObjectToClass(obj.(*core.Object))
	return class.GetName()
}
