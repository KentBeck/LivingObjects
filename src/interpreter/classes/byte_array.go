package classes

import (
	"fmt"
	"unsafe"

	"smalltalklsp/interpreter/core"
)

// ByteArray represents a Smalltalk byte array object
type ByteArray struct {
	core.Object
	Bytes []byte
}

// ByteArrayToObject converts a ByteArray to an Object
func ByteArrayToObject(ba *ByteArray) *core.Object {
	return (*core.Object)(unsafe.Pointer(ba))
}

// ObjectToByteArray converts an Object to a ByteArray
func ObjectToByteArray(o *core.Object) *ByteArray {
	return (*ByteArray)(unsafe.Pointer(o))
}

// String returns a string representation of the byte array object
func (ba *ByteArray) String() string {
	return fmt.Sprintf("ByteArray(%d)", len(ba.Bytes))
}

// Size returns the size of the byte array
func (ba *ByteArray) Size() int {
	return len(ba.Bytes)
}

// At returns the byte at the given index
func (ba *ByteArray) At(index int) byte {
	if index < 0 || index >= len(ba.Bytes) {
		panic("index out of bounds")
	}
	return ba.Bytes[index]
}

// AtPut sets the byte at the given index
func (ba *ByteArray) AtPut(index int, value byte) {
	if index < 0 || index >= len(ba.Bytes) {
		panic("index out of bounds")
	}
	ba.Bytes[index] = value
}
