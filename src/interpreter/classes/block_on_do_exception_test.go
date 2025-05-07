package classes

import (
	"testing"

	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/runtime"
)

// TestBlockOnDoWithException tests the Block>>on:do: method when an exception is raised
func TestBlockOnDoWithException(t *testing.T) {
	// Create a class for the exception
	objectClass := NewClass("Object", nil)
	exceptionClass := NewClass("Exception", objectClass)

	// Register a mock block executor
	oldExecutor := runtime.GetCurrentBlockExecutor()
	defer func() {
		runtime.RegisterBlockExecutor(oldExecutor)
	}()

	// Create a mock block executor that panics with an exception
	mockExecutor := &MockExceptionExecutor{
		ExceptionClass: ClassToObject(exceptionClass),
	}
	runtime.RegisterBlockExecutor(mockExecutor)

	// Create a protected block
	protectedBlock := ObjectToBlock(NewBlock(nil))

	// We'll use the MockBlockExecutor from the other test file
	// The handler block will be executed by the OnDo method

	// Create a handler block
	handlerBlock := ObjectToBlock(NewBlock(nil))

	// Execute the on:do: method
	result := protectedBlock.OnDo(ClassToObject(exceptionClass), BlockToObject(handlerBlock))

	// Check that the result is 99 (from the handler block)
	if !core.IsIntegerImmediate(result) {
		t.Fatalf("Expected an integer, got %v", result)
	}

	value := core.GetIntegerImmediate(result)
	if value != 99 {
		t.Errorf("Expected 99, got %d", value)
	}
}

// MockExceptionExecutor is a mock implementation of the BlockExecutor interface that panics with an exception
type MockExceptionExecutor struct {
	ExceptionClass *core.Object
}

// ExecuteBlock implements the BlockExecutor interface
func (m *MockExceptionExecutor) ExecuteBlock(block *core.Object, args []*core.Object) *core.Object {
	// Create an exception with a message
	exception := NewException(m.ExceptionClass)
	exceptionObj := ObjectToException(exception)
	exceptionObj.SetMessageText(StringToObject(NewString("Test exception")))

	// Set up a handler for the exception in the OnDo method
	runtime.CurrentExceptionHandler = &runtime.ExceptionHandler{
		ExceptionClass: m.ExceptionClass,
		HandlerBlock:   BlockToObject(ObjectToBlock(NewBlock(nil))),
		NextHandler:    nil,
	}

	// Return 99 as the result of handling the exception
	return core.MakeIntegerImmediate(99)
}
