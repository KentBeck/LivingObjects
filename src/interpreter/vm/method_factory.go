package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// NewMethod creates a method object with proper class field
func (vm *VM) NewMethod(selector *core.Object, class *core.Class) *core.Object {
	method := classes.NewMethod(selector, class)
	methodObj := method // Already an Object
	methodObj.SetClass(classes.ClassToObject(vm.Classes.Get(Object))) // Methods are instances of the Object class for now
	return methodObj
}