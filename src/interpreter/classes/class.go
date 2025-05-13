package classes

import (
	"fmt"
	"unsafe"

	"smalltalklsp/interpreter/core"
)

// Don't define a Class type in the classes package anymore
// Use core.Class directly

const METHOD_DICTIONARY_IV = 0

// NewClass creates a new class object
func NewClass(name string, superClass *core.Class) *core.Class {
	// For classes, we need a special instance variable for the method dictionary
	// We'll store it at index 0
	instVars := make([]*core.Object, 1)
	instVars[0] = NewDictionary() // methodDict at index 0

	// Create a core.Class
	result := core.NewClass(name, superClass)
	
	// Set the dictionary
	result.InstanceVarsField = instVars

	return result
}

// ClassToObject converts a Class to an Object
func ClassToObject(c *core.Class) *core.Object {
	return (*core.Object)(unsafe.Pointer(c))
}

// ObjectToClass converts an Object to a Class
func ObjectToClass(o *core.Object) *core.Class {
	return (*core.Class)(unsafe.Pointer(o))
}

// GetClassString returns a string representation of the class object
func GetClassString(c *core.Class) string {
	return fmt.Sprintf("Class %s", c.Name)
}

// GetClassName returns the name of the class
func GetClassName(c *core.Class) string {
	return c.Name
}

// SetClassName sets the name of the class
func SetClassName(c *core.Class, name string) {
	c.Name = name
}

// GetClassSuperClass returns the superclass of the class
func GetClassSuperClass(c *core.Class) *core.Object {
	return c.SuperClass
}

// SetClassSuperClass sets the superclass of the class
func SetClassSuperClass(c *core.Class, superClass *core.Object) {
	c.SuperClass = superClass
}

// GetClassInstanceVarNames returns the instance variable names of the class
func GetClassInstanceVarNames(c *core.Class) []string {
	return c.InstanceVarNames
}

// AddClassInstanceVarName adds an instance variable name to the class
func AddClassInstanceVarName(c *core.Class, name string) {
	c.InstanceVarNames = append(c.InstanceVarNames, name)
}

// GetClassMethodDictionary returns the method dictionary of the class
func GetClassMethodDictionary(c *core.Class) *Dictionary {
	methodDict := c.InstanceVars()[METHOD_DICTIONARY_IV]
	return ObjectToDictionary(methodDict)
}

// AddClassMethod adds a method to the class
func AddClassMethod(c *core.Class, selector *core.Object, method *core.Object) {
	methodDict := GetClassMethodDictionary(c)
	selectorSymbol := ObjectToSymbol(selector)
	methodDict.SetEntry(selectorSymbol.Value, method)
}

// LookupClassMethod looks up a method in the class hierarchy
func LookupClassMethod(c *core.Class, selector *core.Object) *core.Object {
	// Look for the method in this class
	methodDict := GetClassMethodDictionary(c)
	selectorSymbol := ObjectToSymbol(selector)
	method := methodDict.GetEntry(selectorSymbol.Value)
	if method != nil {
		return method
	}

	// If not found, look in the superclass
	if c.SuperClass != nil {
		superClass := ObjectToClass(c.SuperClass)
		return LookupClassMethod(superClass, selector)
	}

	// Method not found
	return nil
}

// NewClassInstance creates a new instance of the class
func NewClassInstance(c *core.Class) *core.Object {
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
		ClassField:        c,
		InstanceVarsField: instVars,
	}
	return obj
}

// GetClassNameFromObject gets the name of a class
func GetClassNameFromObject(obj core.ObjectInterface) string {
	if obj.Type() != core.OBJ_CLASS {
		return ""
	}
	class := ObjectToClass(obj.(*core.Object))
	return GetClassName(class)
}