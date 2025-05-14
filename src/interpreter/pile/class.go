package pile

import (
	"fmt"
	"unsafe"
)


// NewClass creates a new class object
func NewClass(name string, superClass *Class) *Class {
	// Create a Class
	result := &Class{
		Object: Object{
			TypeField:         OBJ_CLASS,
			InstanceVarsField: make([]*Object, 0),
		},
		SuperClass:       (*Object)(unsafe.Pointer(superClass)),
		InstanceVarNames: make([]string, 0),
		Name:             name,
		MethodDictionary: NewDictionary(),
	}
	
	return result
}

// ClassToObject converts a Class to an Object
func ClassToObject(c *Class) *Object {
	return (*Object)(unsafe.Pointer(c))
}

// ObjectToClass converts an Object to a Class
func ObjectToClass(o *Object) *Class {
	return (*Class)(unsafe.Pointer(o))
}

// GetClassString returns a string representation of the class object
func GetClassString(c *Class) string {
	return fmt.Sprintf("Class %s", c.Name)
}

// GetClassName returns the name of the class
func GetClassName(c *Class) string {
	return c.Name
}

// SetClassName sets the name of the class
func SetClassName(c *Class, name string) {
	c.Name = name
}

// GetClassSuperClass returns the superclass of the class
func GetClassSuperClass(c *Class) *Object {
	return c.SuperClass
}

// SetClassSuperClass sets the superclass of the class
func SetClassSuperClass(c *Class, superClass *Object) {
	c.SuperClass = superClass
}

// GetClassInstanceVarNames returns the instance variable names of the class
func GetClassInstanceVarNames(c *Class) []string {
	return c.InstanceVarNames
}

// AddClassInstanceVarName adds an instance variable name to the class
func AddClassInstanceVarName(c *Class, name string) {
	c.InstanceVarNames = append(c.InstanceVarNames, name)
}

// GetClassMethodDictionary returns the method dictionary of the class
func GetClassMethodDictionary(c *Class) *Dictionary {
	return ObjectToDictionary(c.MethodDictionary)
}

// AddClassMethod adds a method to the class
func AddClassMethod(c *Class, selector *Object, method *Object) {
	methodDict := GetClassMethodDictionary(c)
	selectorSymbol := ObjectToSymbol(selector)
	methodDict.SetEntry(selectorSymbol.Value, method)
}

// LookupClassMethod looks up a method in the class hierarchy
func LookupClassMethod(c *Class, selector *Object) *Object {
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
func NewClassInstance(c *Class) *Object {
	// Initialize instance variables array with nil values
	instVarsSize := 0
	if c != nil && len(c.InstanceVarNames) > 0 {
		instVarsSize = len(c.InstanceVarNames)
	}
	instVars := make([]*Object, instVarsSize)
	for i := range instVars {
		instVars[i] = MakeNilImmediate()
	}

	obj := &Object{
		TypeField:         OBJ_INSTANCE,
		ClassField:        c,
		InstanceVarsField: instVars,
	}
	return obj
}

// GetClassNameFromObject gets the name of a class
func GetClassNameFromObject(obj ObjectInterface) string {
	if obj.Type() != OBJ_CLASS {
		return ""
	}
	class := ObjectToClass(obj.(*Object))
	return GetClassName(class)
}
