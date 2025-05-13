package compiler

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// MethodFactory defines an interface for creating Method objects
type MethodFactory interface {
	NewMethod(selector *core.Object, class *core.Class) *core.Object
	NewSymbol(name string) *core.Object
}

// RegisteredMethodFactory is the global factory for creating method objects
// Initially nil, it should be set by the VM during initialization
var RegisteredMethodFactory MethodFactory

// RegisterMethodFactory registers a factory for creating method objects
// This should be called by the VM during initialization
func RegisterMethodFactory(factory MethodFactory) {
	RegisteredMethodFactory = factory
}

// CreateMethod creates a method with the registered factory
// If no factory is registered, it falls back to using classes.NewMethod
func CreateMethod(selector *core.Object, class *core.Class) *core.Object {
	if RegisteredMethodFactory != nil {
		return RegisteredMethodFactory.NewMethod(selector, class)
	}
	
	// Fallback to direct creation (will not have class field set)
	return classes.NewMethod(selector, class)
}

// CreateSymbol creates a symbol with the registered factory
// If no factory is registered, it falls back to using classes.NewSymbol
func CreateSymbol(name string) *core.Object {
	if RegisteredMethodFactory != nil {
		return RegisteredMethodFactory.NewSymbol(name)
	}
	
	// Fallback to direct creation (will not have class field set)
	return classes.NewSymbol(name)
}