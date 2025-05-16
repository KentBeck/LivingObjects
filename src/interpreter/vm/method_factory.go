package vm

import (
	"smalltalklsp/interpreter/pile"
)

// NewMethod creates a method object with proper class field
func (vm *VM) NewMethod(selector *pile.Object, class *pile.Class) *pile.Object {
	method := pile.NewMethod(selector, class)
	methodObj := method // Already an Object
	methodObj.SetClass(vm.Globals["Method"]) // Methods are instances of the Method class
	return methodObj
}