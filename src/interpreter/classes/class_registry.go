package classes

import (
	"unsafe"
	"smalltalklsp/interpreter/core"
)

// ClassRegistry is a global registry of class references
// that helps break circular dependencies in class initialization
var ClassRegistry = &Registry{
	initialized: false,
}

// Registry holds references to standard classes
type Registry struct {
	ObjectClass    *core.Class
	StringClass    *core.Class
	ArrayClass     *core.Class
	SymbolClass    *core.Class
	DictionaryClass *core.Class
	BlockClass     *core.Class
	MethodClass    *core.Class
	ByteArrayClass *core.Class
	
	// Track initialization state
	initialized bool
}

// Initialize sets up the registry with class references
// It should be called once by the VM after all classes are created
func (r *Registry) Initialize(
	objectClass, stringClass, arrayClass, symbolClass, 
	dictionaryClass, blockClass, methodClass, byteArrayClass *core.Class) {
	
	r.ObjectClass = objectClass
	r.StringClass = stringClass
	r.ArrayClass = arrayClass
	r.SymbolClass = symbolClass
	r.DictionaryClass = dictionaryClass
	r.BlockClass = blockClass
	r.MethodClass = methodClass
	r.ByteArrayClass = byteArrayClass
	
	r.initialized = true
}

// IsInitialized returns true if the registry has been initialized
func (r *Registry) IsInitialized() bool {
	return r.initialized
}

// SetClassField sets the class field of an object
// If the registry is not initialized, it does nothing
func (r *Registry) SetClassField(obj *core.Object, class *core.Class) {
	if r.initialized {
		obj.SetClass((*core.Object)(unsafe.Pointer(class)))
	}
}