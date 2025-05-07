package runtime

import (
	"smalltalklsp/interpreter/core"
)

// BlockExecutor is an interface for executing blocks
type BlockExecutor interface {
	// ExecuteBlock executes a block with the given arguments and returns the result
	ExecuteBlock(block *core.Object, args []*core.Object) *core.Object
}

// RegisterBlockExecutor registers a block executor
// This is called by the VM to register itself as the block executor
var currentBlockExecutor BlockExecutor

// RegisterBlockExecutor registers a block executor
func RegisterBlockExecutor(executor BlockExecutor) {
	currentBlockExecutor = executor
}

// ExecuteBlock executes a block with the given arguments and returns the result
func ExecuteBlock(block *core.Object, args []*core.Object) *core.Object {
	if currentBlockExecutor == nil {
		// If no block executor is registered, return nil
		// This is a temporary solution for testing
		return core.MakeNilImmediate()
	}
	return currentBlockExecutor.ExecuteBlock(block, args)
}
