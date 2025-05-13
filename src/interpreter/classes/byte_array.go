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

// ByteArray is created by the VM using NewByteArray

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

// Copy returns a copy of the byte array
func (ba *ByteArray) Copy() *ByteArray {
	newBA := &ByteArray{
		Object: core.Object{
			TypeField: core.OBJ_BYTE_ARRAY,
		},
		Bytes: make([]byte, len(ba.Bytes)),
	}
	copy(newBA.Bytes, ba.Bytes)
	return newBA
}

// CopyFrom returns a new byte array containing the bytes from startIndex to endIndex
func (ba *ByteArray) CopyFrom(startIndex, endIndex int) *ByteArray {
	if startIndex < 0 || startIndex >= len(ba.Bytes) {
		panic("start index out of bounds")
	}
	if endIndex < startIndex || endIndex >= len(ba.Bytes) {
		panic("end index out of bounds")
	}

	newSize := endIndex - startIndex + 1
	newBA := &ByteArray{
		Object: core.Object{
			TypeField: core.OBJ_BYTE_ARRAY,
		},
		Bytes: make([]byte, newSize),
	}
	copy(newBA.Bytes, ba.Bytes[startIndex:endIndex+1])
	return newBA
}

// Uint32At reads a 32-bit unsigned integer from the byte array at the given index
func (ba *ByteArray) Uint32At(index int) uint32 {
	if index < 0 || index+3 >= len(ba.Bytes) {
		panic("index out of bounds for uint32")
	}

	return uint32(ba.Bytes[index]) |
		uint32(ba.Bytes[index+1])<<8 |
		uint32(ba.Bytes[index+2])<<16 |
		uint32(ba.Bytes[index+3])<<24
}

// Uint32AtPut writes a 32-bit unsigned integer to the byte array at the given index
func (ba *ByteArray) Uint32AtPut(index int, value uint32) {
	if index < 0 || index+3 >= len(ba.Bytes) {
		panic("index out of bounds for uint32")
	}

	ba.Bytes[index] = byte(value)
	ba.Bytes[index+1] = byte(value >> 8)
	ba.Bytes[index+2] = byte(value >> 16)
	ba.Bytes[index+3] = byte(value >> 24)
}
