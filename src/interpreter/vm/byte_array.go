package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/core"
)

// NewByteArrayClass creates a new ByteArray class
func (vm *VM) NewByteArrayClass() *classes.Class {
	result := classes.NewClass("ByteArray", vm.ObjectClass)

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
func (vm *VM) newByteArrayInternal(size int) *classes.ByteArray {
	return &classes.ByteArray{
		Object: core.Object{
			TypeField: core.OBJ_BYTE_ARRAY,
		},
		Bytes: make([]byte, size),
	}
}

// NewByteArray creates a new byte array object
func (vm *VM) NewByteArray(size int) *core.Object {
	byteArray := vm.newByteArrayInternal(size)
	byteArrayObj := classes.ByteArrayToObject(byteArray)
	byteArrayObj.SetClass(classes.ClassToObject(vm.ByteArrayClass))
	return byteArrayObj
}
