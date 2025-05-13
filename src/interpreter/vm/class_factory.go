package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// NewClass creates a new class object with proper class field
func (vm *VM) NewClass(name string, superClass *core.Class) *core.Class {
	// For classes, we need a special instance variable for the method dictionary
	// We'll store it at index 0
	instVars := make([]*core.Object, 1)
	instVars[0] = vm.NewDictionary() // methodDict at index 0
	
	// Create a core.Class
	class := core.NewClass(name, superClass)
	
	// Store the method dictionary
	class.InstanceVarsField = instVars
	
	// A class's class should be a metaclass, but for now we'll use ObjectClass
	classObj := classes.ClassToObject(class)
	classObj.SetClass(classes.ClassToObject(vm.Classes.Get(Object)))
	
	return class
}