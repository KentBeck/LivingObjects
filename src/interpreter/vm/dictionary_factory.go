package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// NewDictionary creates a dictionary object with proper class field
func (vm *VM) NewDictionary() *core.Object {
	dict := &classes.Dictionary{
		Object: core.Object{
			TypeField: core.OBJ_DICTIONARY,
		},
		Entries: make(map[string]*core.Object),
	}
	
	dictObj := classes.DictionaryToObject(dict)
	dictObj.SetClass(classes.ClassToObject(vm.ObjectClass)) // Dictionary is an instance of Object class for now
	return dictObj
}