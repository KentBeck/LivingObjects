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

// NewByteArray creates a new byte array object
func (vm *VM) NewByteArray(size int) *core.Object {
	byteArray := classes.NewByteArray(size)
	byteArrayObj := classes.ByteArrayToObject(byteArray)
	byteArrayObj.SetClass(classes.ClassToObject(vm.ByteArrayClass))
	return byteArrayObj
}
