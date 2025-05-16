package vm

import (
	"smalltalklsp/interpreter/pile"
)

// NewClass creates a new class object with proper class field
func (vm *VM) NewClass(name string, superClass *pile.Class) *pile.Class {
	// For classes, we need a special instance variable for the method dictionary
	// We'll store it at index 0
	instVars := make([]*pile.Object, 1)
	instVars[0] = vm.NewDictionary() // methodDict at index 0
	
	// Create a pile.Class
	class := pile.NewClass(name, superClass)
	
	// Store the method dictionary
	class.InstanceVarsField = instVars
	
	// A class's class should be a metaclass, but for now we'll use ObjectClass
	classObj := pile.ClassToObject(class)
	classObj.SetClass(vm.Globals["Object"])
	
	return class
}