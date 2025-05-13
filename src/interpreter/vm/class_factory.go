package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// NewClass creates a new class object with proper class field
func (vm *VM) NewClass(name string, superClass *classes.Class) *classes.Class {
	// For classes, we need a special instance variable for the method dictionary
	// We'll store it at index 0
	instVars := make([]*core.Object, 1)
	instVars[0] = vm.NewDictionary() // methodDict at index 0
	
	class := &classes.Class{
		Object: core.Object{
			TypeField:         core.OBJ_CLASS,
			InstanceVarsField: instVars,
		},
		SuperClass:       classes.ClassToObject(superClass),
		InstanceVarNames: make([]string, 0),
		Name:             name,
	}
	
	// A class's class should be a metaclass, but for now we'll use ObjectClass
	classObj := classes.ClassToObject(class)
	classObj.SetClass(classes.ClassToObject(vm.Classes.Get(Object)))
	
	return class
}