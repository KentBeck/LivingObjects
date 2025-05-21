package vm

import (
	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/pile"
)

// NewByteArrayClass creates a new ByteArray class
func (vm *VM) NewByteArrayClass() *pile.Class {
	objectClass := pile.ObjectToClass(vm.Globals["Object"])
	result := pile.NewClass("ByteArray", objectClass)

	// Add primitive methods to the ByteArray class - create a new builder for each method
	
	// at: method (returns the byte at the given index)
	compiler.NewMethodBuilder(result).Primitive(50).Go("at:")

	// at:put: method (sets the byte at the given index)
	compiler.NewMethodBuilder(result).Primitive(51).Go("at:put:")

	return result
}

// newByteArrayInternal creates a new byte array object without setting its class
// This is a private helper function used by NewByteArray
func (vm *VM) newByteArrayInternal(size int) *pile.ByteArray {
	return &pile.ByteArray{
		Object: pile.Object{
			TypeField: pile.OBJ_BYTE_ARRAY,
		},
		Bytes: make([]byte, size),
	}
}

// NewByteArray creates a new byte array object
func (vm *VM) NewByteArray(size int) *pile.Object {
	byteArray := vm.newByteArrayInternal(size)
	byteArrayObj := pile.ByteArrayToObject(byteArray)
	byteArrayObj.SetClass(vm.Globals["ByteArray"])
	return byteArrayObj
}
