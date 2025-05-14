package types

import (
	"smalltalklsp/interpreter/core"
)

// ObjectFactory defines an interface for creating Smalltalk objects
type ObjectFactory interface {
	// NewBlock creates a block object with proper class field
	NewBlock(outerContext interface{}) *core.Object
}

// DefaultFactory is a singleton instance of the ObjectFactory
var DefaultFactory ObjectFactory

// RegisterFactory registers the default object factory
func RegisterFactory(factory ObjectFactory) {
	DefaultFactory = factory
}

// GetFactory returns the default object factory
func GetFactory() ObjectFactory {
	return DefaultFactory
}