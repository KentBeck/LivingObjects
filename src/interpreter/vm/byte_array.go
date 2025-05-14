package vm

import (
	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/pile"
)

// NewByteArrayClass creates a new ByteArray class
func (vm *VM) NewByteArrayClass() *pile.Class {
	result := pile.NewClass("ByteArray", vm.Classes.Get(Object))

	// Add primitive methods to the ByteArray class
	builder := compiler.NewMethodBuilder(result)

	// at: method (returns the byte at the given index)
	builder.Selector("at:").Primitive(50).Go()

	// at:put: method (sets the byte at the given index)
	builder.Selector("at:put:").Primitive(51).Go()

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
	byteArrayObj.SetClass(pile.ClassToObject(vm.Classes.Get(ByteArray)))
	return byteArrayObj
}
