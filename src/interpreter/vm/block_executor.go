package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/runtime"
)

// init registers the VM as a block executor
func init() {
	// The VM will be registered as a block executor when it's created
}

// ExecuteBlock implements the runtime.BlockExecutor interface
func (vm *VM) ExecuteBlock(block *core.Object, args []*core.Object) *core.Object {
	// Check if the block is valid
	if block == nil {
		panic("ExecuteBlock: nil block")
	}

	// Convert the block to a Block
	blockObj := classes.ObjectToBlock(block)
	if blockObj == nil {
		panic("ExecuteBlock: invalid block")
	}

	// Get the outer context
	outerContext, ok := blockObj.GetOuterContext().(*Context)
	if !ok {
		panic("ExecuteBlock: invalid outer context")
	}

	// Create a method object for the block
	methodObj := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes:    blockObj.GetBytecodes(),
		Literals:     blockObj.GetLiterals(),
		TempVarNames: blockObj.GetTempVarNames(),
	}

	// Create a new context for the block execution
	blockContext := NewContext(
		classes.MethodToObject(methodObj),
		outerContext.GetReceiver(),
		args,
		outerContext,
	)

	// Save the current context
	savedContext := vm.CurrentContext

	// Set the current context to the block context
	vm.CurrentContext = blockContext

	// Execute the block
	result, err := vm.ExecuteContext(blockContext)
	if err != nil {
		panic("ExecuteBlock: " + err.Error())
	}

	// Restore the current context
	vm.CurrentContext = savedContext

	// Return the result
	return result.(*core.Object)
}

// RegisterAsBlockExecutor registers the VM as a block executor
func (vm *VM) RegisterAsBlockExecutor() {
	runtime.RegisterBlockExecutor(vm)
}
