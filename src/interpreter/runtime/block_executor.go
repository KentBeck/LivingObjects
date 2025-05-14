package runtime

import (
	"smalltalklsp/interpreter/pile"
)

// BlockExecutor is an interface for executing blocks
type BlockExecutor interface {
	// ExecuteBlock executes a block with the given arguments and returns the result
	ExecuteBlock(block *pile.Object, args []*pile.Object) *pile.Object
}

// RegisterBlockExecutor registers a block executor
// This is called by the VM to register itself as the block executor
var currentBlockExecutor BlockExecutor

// RegisterBlockExecutor registers a block executor
func RegisterBlockExecutor(executor BlockExecutor) {
	currentBlockExecutor = executor
	
	// Also set the pile.ExecuteBlock function to use our executor
	pile.SetBlockExecutor(func(block *pile.Object, args []*pile.Object) *pile.Object {
		return executor.ExecuteBlock(block, args)
	})
}

// GetCurrentBlockExecutor returns the current block executor
func GetCurrentBlockExecutor() BlockExecutor {
	return currentBlockExecutor
}

// ExecuteBlock executes a block with the given arguments and returns the result
func ExecuteBlock(block *pile.Object, args []*pile.Object) *pile.Object {
	if currentBlockExecutor == nil {
		// If no block executor is registered, return nil
		// This is a temporary solution for testing
		return pile.MakeNilImmediate()
	}
	return currentBlockExecutor.ExecuteBlock(block, args)
}