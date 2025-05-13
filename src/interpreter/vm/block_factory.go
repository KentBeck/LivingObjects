package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// NewBlock creates a block object with proper class field
func (vm *VM) NewBlock(outerContext interface{}) *core.Object {
	block := classes.NewBlockInternal(outerContext)
	blockObj := classes.BlockToObject(block)
	blockObj.SetClass(classes.ClassToObject(vm.Classes.Get(Block)))
	return blockObj
}