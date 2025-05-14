package pile_test

import (
	"testing"

	"smalltalklsp/interpreter/pile"
)

// TestBlockOnDoNoException tests the Block>>on:do: method when no exception is raised
func TestBlockOnDoNoException(t *testing.T) {
	// Create a class for the exception
	objectClass := pile.NewClass("Object", nil)
	exceptionClass := pile.NewClass("Exception", objectClass)

	// Register a mock block executor
	oldExecutor := pile.GetCurrentBlockExecutor()
	defer func() {
		pile.RegisterBlockExecutor(oldExecutor)
	}()

	// Create a mock block executor that returns the value 42
	mockExecutor := &MockBlockExecutor{
		ReturnValue: pile.MakeIntegerImmediate(42),
	}
	pile.RegisterBlockExecutor(mockExecutor)

	// Create a protected block
	protectedBlock := pile.ObjectToBlock(pile.NewBlock(nil))

	// Create a handler block
	handlerBlock := pile.ObjectToBlock(pile.NewBlock(nil))

	// Execute the on:do: method
	result := protectedBlock.OnDo(pile.ClassToObject(exceptionClass), pile.BlockToObject(handlerBlock))

	// Check that the result is 42 (from the protected block)
	if !pile.IsIntegerImmediate(result) {
		t.Fatalf("Expected an integer, got %v", result)
	}

	value := pile.GetIntegerImmediate(result)
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}
}