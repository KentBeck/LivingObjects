package vm

import (
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/runtime"
)

// init registers the VM as a block executor
func init() {
	// The VM will be registered as a block executor when it's created
}

// ExecuteBlock implements the runtime.BlockExecutor interface
func (vm *VM) ExecuteBlock(block *pile.Object, args []*pile.Object) *pile.Object {
	// Check if the block is valid
	if block == nil {
		panic("ExecuteBlock: nil block")
	}

	// Convert the block to a Block
	blockObj := pile.ObjectToBlock(block)
	if blockObj == nil {
		panic("ExecuteBlock: invalid block")
	}

	// Get the outer context
	outerContext, ok := blockObj.GetOuterContext().(*Context)
	if !ok {
		panic("ExecuteBlock: invalid outer context")
	}

	// Create a method object for the block
	methodObj := &pile.Method{
		Object: pile.Object{
			TypeField: pile.OBJ_METHOD,
		},
		Bytecodes:    blockObj.GetBytecodes(),
		Literals:     blockObj.GetLiterals(),
		TempVarNames: blockObj.GetTempVarNames(),
	}

	// Create a new context for the block execution
	blockContext := NewContext(
		pile.MethodToObject(methodObj),
		outerContext.GetReceiver(),
		args,
		outerContext,
	)

	// Manually set up the temporary variables with the arguments
	// This is necessary because the arguments passed to NewContext
	// are not automatically copied to the temporary variables
	for i, arg := range args {
		if i < len(blockObj.GetTempVarNames()) {
			blockContext.SetTempVarByIndex(i, arg)
		}
	}

	// Save the current context
	savedContext := vm.Executor.CurrentContext

	// Set the current context to the block context
	vm.Executor.CurrentContext = blockContext

	// Execute the block
	result, err := vm.ExecuteContext(blockContext)
	if err != nil {
		panic("ExecuteBlock: " + err.Error())
	}

	// Restore the current context
	vm.Executor.CurrentContext = savedContext

	// Return the result
	return result.(*pile.Object)
}

// RegisterAsBlockExecutor registers the VM as a block executor
func (vm *VM) RegisterAsBlockExecutor() {
	runtime.RegisterBlockExecutor(vm)
}
