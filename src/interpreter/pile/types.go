package pile

// Factory is an interface for creating Smalltalk objects
type Factory interface {
	// NewBlock creates a new block object
	NewBlock(outerContext interface{}) *Object
}

// Current factory instance
var currentFactory Factory

// RegisterFactory registers a factory
func RegisterFactory(factory Factory) {
	currentFactory = factory
}

// GetFactory returns the current factory
func GetFactory() Factory {
	return currentFactory
}