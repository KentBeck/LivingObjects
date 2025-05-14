package pile

// ExceptionHandler represents an exception handler
type ExceptionHandler struct {
	ExceptionClass *Object
	HandlerBlock   *Object
	NextHandler    *ExceptionHandler
}

// CurrentExceptionHandler is the current active exception handler
var CurrentExceptionHandler *ExceptionHandler

// IsKindOf checks if an object is an instance of a class or one of its subclasses
// This is a simplified implementation that just checks if the classes are the same
// In a real implementation, we would check the class hierarchy
func IsKindOf(obj *Object, class *Object) bool {
	return obj.Class() == class
}

// BlockExecutor is an interface for executing blocks
type BlockExecutor interface {
	// ExecuteBlock executes a block with the given arguments and returns the result
	ExecuteBlock(block *Object, args []*Object) *Object
}

// RegisterBlockExecutor registers a block executor
// This is called by the VM to register itself as the block executor
var currentBlockExecutor BlockExecutor

// RegisterBlockExecutor registers a block executor
func RegisterBlockExecutor(executor BlockExecutor) {
	currentBlockExecutor = executor
}

// GetCurrentBlockExecutor returns the current block executor
func GetCurrentBlockExecutor() BlockExecutor {
	return currentBlockExecutor
}

// ExecuteBlock executes a block with the given arguments and returns the result
// This is set as a variable to allow for testing
var ExecuteBlock = func(block *Object, args []*Object) *Object {
	if currentBlockExecutor == nil {
		// If no block executor is registered, return nil
		// This is a temporary solution for testing
		return MakeNilImmediate()
	}
	return currentBlockExecutor.ExecuteBlock(block, args)
}

// SetBlockExecutor sets the block executor function
// This allows external packages to provide a block executor without circular imports
func SetBlockExecutor(executor func(block *Object, args []*Object) *Object) {
	ExecuteBlock = executor
}

// SignalException signals an exception
// If there's a handler for the exception, it will be executed
// Otherwise, it will panic with the exception
func SignalException(exception *Object) *Object {
	// If there's no handler, just panic with the exception
	if CurrentExceptionHandler == nil {
		panic(exception)
	}

	// Find a handler for this exception
	handler := CurrentExceptionHandler
	for handler != nil {
		if IsKindOf(exception, handler.ExceptionClass) {
			// Found a handler, execute it
			return ExecuteBlock(handler.HandlerBlock, []*Object{exception})
		}
		handler = handler.NextHandler
	}

	// No handler found, panic with the exception
	panic(exception)
}