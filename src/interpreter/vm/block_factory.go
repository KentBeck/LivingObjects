package vm

import (
	"smalltalklsp/interpreter/pile"
)

// NewBlock creates a block object with proper class field
func (vm *VM) NewBlock(outerContext interface{}) *pile.Object {
	block := &pile.Block{
		Object: pile.Object{
			TypeField: pile.OBJ_BLOCK,
		},
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*pile.Object, 0),
		TempVarNames: make([]string, 0),
		OuterContext: outerContext,
	}
	blockObj := pile.BlockToObject(block)
	blockObj.SetClass(pile.ClassToObject(vm.Classes.Get(Block)))
	return blockObj
}
