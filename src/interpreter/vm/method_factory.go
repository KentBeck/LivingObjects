package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// NewMethod creates a method object with proper class field
func (vm *VM) NewMethod(selector *core.Object, class *classes.Class) *core.Object {
	method := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*core.Object, 0),
		Selector:     selector,
		MethodClass:  class,
		TempVarNames: make([]string, 0),
	}

	methodObj := classes.MethodToObject(method)
	methodObj.SetClass(classes.ClassToObject(vm.Classes.Get(Object))) // Methods are instances of the Object class for now
	return methodObj
}