package classes

import (
	"testing"

	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/runtime"
)

// TestBlockOnDoNoException tests the Block>>on:do: method when no exception is raised
func TestBlockOnDoNoException(t *testing.T) {
	// Create a class for the exception
	objectClass := NewClass("Object", nil)
	exceptionClass := NewClass("Exception", objectClass)

	// Register a mock block executor
	oldExecutor := runtime.GetCurrentBlockExecutor()
	defer func() {
		runtime.RegisterBlockExecutor(oldExecutor)
	}()

	// Create a mock block executor that returns the value 42
	mockExecutor := &MockBlockExecutor{
		ReturnValue: core.MakeIntegerImmediate(42),
	}
	runtime.RegisterBlockExecutor(mockExecutor)

	// Create a protected block
	protectedBlock := ObjectToBlock(NewBlock(nil))

	// Create a handler block
	handlerBlock := ObjectToBlock(NewBlock(nil))

	// Execute the on:do: method
	result := protectedBlock.OnDo(ClassToObject(exceptionClass), BlockToObject(handlerBlock))

	// Check that the result is 42 (from the protected block)
	if !core.IsIntegerImmediate(result) {
		t.Fatalf("Expected an integer, got %v", result)
	}

	value := core.GetIntegerImmediate(result)
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}
}

// MockBlockExecutor is a mock implementation of the BlockExecutor interface
type MockBlockExecutor struct {
	ReturnValue *core.Object
}

// ExecuteBlock implements the BlockExecutor interface
func (m *MockBlockExecutor) ExecuteBlock(block *core.Object, args []*core.Object) *core.Object {
	return m.ReturnValue
}
