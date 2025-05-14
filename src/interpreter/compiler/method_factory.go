package compiler

import (
	"smalltalklsp/interpreter/pile"
)

// MethodFactory defines an interface for creating Method objects
type MethodFactory interface {
	NewMethod(selector *pile.Object, class *pile.Class) *pile.Object
	NewSymbol(name string) *pile.Object
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
// If no factory is registered, it falls back to using pile.NewMethod
func CreateMethod(selector *pile.Object, class *pile.Class) *pile.Object {
	if RegisteredMethodFactory != nil {
		return RegisteredMethodFactory.NewMethod(selector, class)
	}
	
	// Fallback to direct creation (will not have class field set)
	return pile.NewMethod(selector, class)
}

// CreateSymbol creates a symbol with the registered factory
// If no factory is registered, it falls back to using pile.NewSymbol
func CreateSymbol(name string) *pile.Object {
	if RegisteredMethodFactory != nil {
		return RegisteredMethodFactory.NewSymbol(name)
	}
	
	// Fallback to direct creation (will not have class field set)
	return pile.NewSymbol(name)
}