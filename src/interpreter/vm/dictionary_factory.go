package vm

import (
	"smalltalklsp/interpreter/pile"
)

// NewDictionary creates a dictionary object with proper class field
func (vm *VM) NewDictionary() *pile.Object {
	dict := pile.NewDictionaryInternal()
	dictObj := pile.DictionaryToObject(dict)
	dictObj.SetClass(vm.Globals["Dictionary"]) // Dictionary is an instance of Dictionary class
	return dictObj
}