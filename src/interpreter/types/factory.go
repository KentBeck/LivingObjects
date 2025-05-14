package types

import (
	"smalltalklsp/interpreter/pile"
)

// ObjectFactory defines an interface for creating Smalltalk objects
type ObjectFactory interface {
	// NewBlock creates a block object with proper class field
	NewBlock(outerContext interface{}) *pile.Object
}

// DefaultFactory is a singleton instance of the ObjectFactory
var DefaultFactory ObjectFactory

// RegisterFactory registers the default object factory
func RegisterFactory(factory ObjectFactory) {
	DefaultFactory = factory
	
	// Set up the hook for pile.NewBlock to use DefaultFactory
	pile.SetFactoryRegisterHook(func(block *pile.Object, outerContext interface{}) *pile.Object {
		if DefaultFactory != nil {
			return DefaultFactory.NewBlock(outerContext)
		}
		return block
	})
}

// GetFactory returns the default object factory
func GetFactory() ObjectFactory {
	return DefaultFactory
}