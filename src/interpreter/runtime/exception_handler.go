package runtime

import (
	"smalltalklsp/interpreter/core"
)

// ExceptionHandler represents an exception handler
type ExceptionHandler struct {
	ExceptionClass *core.Object
	HandlerBlock   *core.Object
	NextHandler    *ExceptionHandler
}

// CurrentExceptionHandler is the current active exception handler
var CurrentExceptionHandler *ExceptionHandler

// IsKindOf checks if an object is an instance of a class or one of its subclasses
// This is a simplified implementation that just checks if the classes are the same
// In a real implementation, we would check the class hierarchy
func IsKindOf(obj *core.Object, class *core.Object) bool {
	return obj.Class() == class
}

// SignalException signals an exception
// If there's a handler for the exception, it will be executed
// Otherwise, it will panic with the exception
func SignalException(exception *core.Object) *core.Object {
	// If there's no handler, just panic with the exception
	if CurrentExceptionHandler == nil {
		panic(exception)
	}

	// Find a handler for this exception
	handler := CurrentExceptionHandler
	for handler != nil {
		if IsKindOf(exception, handler.ExceptionClass) {
			// Found a handler, execute it
			return ExecuteBlock(handler.HandlerBlock, []*core.Object{exception})
		}
		handler = handler.NextHandler
	}

	// No handler found, panic with the exception
	panic(exception)
}
