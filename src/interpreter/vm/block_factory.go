package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// NewBlock creates a block object with proper class field
func (vm *VM) NewBlock(outerContext interface{}) *core.Object {
	block := &classes.Block{
		Object: core.Object{
			TypeField: core.OBJ_BLOCK,
		},
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*core.Object, 0),
		TempVarNames: make([]string, 0),
		OuterContext: outerContext,
	}
	
	blockObj := classes.BlockToObject(block)
	blockObj.SetClass(classes.ClassToObject(vm.Classes.Get(Block)))
	return blockObj
}