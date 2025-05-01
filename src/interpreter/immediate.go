package main

import (
	"math"
	"unsafe"
)

// Tag bits for immediate values
const (
	TAG_POINTER = 0x0 // 00 - Pointer to an object
	TAG_SPECIAL = 0x1 // 01 - Special value (nil, true, false)
	TAG_FLOAT   = 0x2 // 10 - 62-bit float
	TAG_INTEGER = 0x3 // 11 - 62-bit integer
	TAG_MASK    = 0x3 // Mask for the tag bits
)

// Special values
const (
	SPECIAL_NIL   = 0x1 // 01 (TAG_SPECIAL | 0 << 2)
	SPECIAL_TRUE  = 0x5 // 101 (TAG_SPECIAL | 1 << 2)
	SPECIAL_FALSE = 0x9 // 1001 (TAG_SPECIAL | 2 << 2)
)

// IsImmediate returns true if the value is an immediate value
func IsImmediate(obj ObjectInterface) bool {
	converted := obj.(*Object)
	ptr := uintptr(unsafe.Pointer(converted))
	return (ptr & TAG_MASK) != TAG_POINTER
}

// GetTag returns the tag bits of a value
func GetTag(obj ObjectInterface) int {
	// Convert the pointer to an integer
	converted := obj.(*Object)
	ptr := uintptr(unsafe.Pointer(converted))

	// Return the bottom two bits
	return int(ptr & TAG_MASK)
}

// IsNilImmediate returns true if the value is the immediate nil value
func IsNilImmediate(obj ObjectInterface) bool {
	// Convert the pointer to an integer
	converted := obj.(*Object)
	ptr := uintptr(unsafe.Pointer(converted))

	// Check if it's the nil immediate value
	return ptr == SPECIAL_NIL
}

// IsTrueImmediate returns true if the value is the immediate true value
func IsTrueImmediate(obj ObjectInterface) bool {
	// Convert the pointer to an integer
	converted := obj.(*Object)
	ptr := uintptr(unsafe.Pointer(converted))

	// Check if it's the true immediate value
	return ptr == SPECIAL_TRUE
}

// IsFalseImmediate returns true if the value is the immediate false value
func IsFalseImmediate(obj ObjectInterface) bool {
	// Convert the pointer to an integer
	converted := obj.(*Object)
	ptr := uintptr(unsafe.Pointer(converted))

	// Check if it's the false immediate value
	return ptr == SPECIAL_FALSE
}

// MakeNilImmediate returns the immediate nil value
func MakeNilImmediate() *Object {
	// Convert the immediate nil value to a pointer
	return (*Object)(unsafe.Pointer(uintptr(SPECIAL_NIL)))
}

// MakeTrueImmediate returns the immediate true value
func MakeTrueImmediate() *Object {
	// Convert the immediate true value to a pointer
	return (*Object)(unsafe.Pointer(uintptr(SPECIAL_TRUE)))
}

// MakeFalseImmediate returns the immediate false value
func MakeFalseImmediate() *Object {
	// Convert the immediate false value to a pointer
	return (*Object)(unsafe.Pointer(uintptr(SPECIAL_FALSE)))
}

// IsIntegerImmediate returns true if the value is an immediate integer
func IsIntegerImmediate(obj ObjectInterface) bool {
	converted := obj.(*Object)
	ptr := uintptr(unsafe.Pointer(converted))
	return (ptr & TAG_MASK) == TAG_INTEGER
}

// MakeIntegerImmediate returns an immediate integer value
func MakeIntegerImmediate(value int64) *Object {
	// Ensure the value fits in 62 bits (signed)
	if value > 0x1FFFFFFFFFFFFFFF || value < -0x2000000000000000 {
		panic("Integer value too large for immediate representation")
	}

	// Shift the value left by 2 bits and set the tag bits
	imm := (uintptr(value) << 2) | TAG_INTEGER

	// Convert to a pointer
	return (*Object)(unsafe.Pointer(imm))
}

// GetIntegerImmediate extracts the integer value from an immediate integer
func GetIntegerImmediate(obj ObjectInterface) int64 {
	converted := obj.(*Object)
	ptr := uintptr(unsafe.Pointer(converted))
	unsigned := ptr >> 2

	// Handle sign extension for negative numbers
	if (unsigned & (1 << 61)) != 0 {
		// It's a negative number, sign extend
		unsigned |= 0xC000000000000000
	}

	return int64(unsigned)
}

// IsFloatImmediate returns true if the value is an immediate float
func IsFloatImmediate(obj *Object) bool {
	// Convert the pointer to an integer
	ptr := uintptr(unsafe.Pointer(obj))

	// Check if the tag is TAG_FLOAT
	return (ptr & TAG_MASK) == TAG_FLOAT
}

// MakeFloatImmediate returns an immediate float value
func MakeFloatImmediate(value float64) *Object {
	// Convert the float to bits
	bits := math.Float64bits(value)
	// We lose some precision, but that's acceptable for most use cases
	// The bottom 2 bits are used for the tag
	imm := (bits >> 2 << 2) | TAG_FLOAT

	// Convert to a pointer
	return (*Object)(unsafe.Pointer(uintptr(imm)))
}

// GetFloatImmediate extracts the float value from an immediate float
func GetFloatImmediate(obj *Object) float64 {
	// Convert the pointer to an integer
	ptr := uintptr(unsafe.Pointer(obj))

	// Remove the tag bits
	bits := ptr & ^uintptr(TAG_MASK)

	// Convert to float64
	return math.Float64frombits(uint64(bits))
}
