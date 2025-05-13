package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// NewDictionary creates a dictionary object with proper class field
func (vm *VM) NewDictionary() *core.Object {
	dict := classes.NewDictionaryInternal()
	dictObj := classes.DictionaryToObject(dict)
	dictObj.SetClass(classes.ClassToObject(vm.Classes.Get(Object))) // Dictionary is an instance of Object class for now
	return dictObj
}