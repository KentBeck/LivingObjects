package main

import (
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
func IsImmediate(obj *Object) bool {
	// Convert the pointer to an integer
	ptr := uintptr(unsafe.Pointer(obj))

	// Check if the bottom two bits are not 00
	return (ptr & TAG_MASK) != TAG_POINTER
}

// GetTag returns the tag bits of a value
func GetTag(obj *Object) int {
	// Convert the pointer to an integer
	ptr := uintptr(unsafe.Pointer(obj))

	// Return the bottom two bits
	return int(ptr & TAG_MASK)
}

// IsNilImmediate returns true if the value is the immediate nil value
func IsNilImmediate(obj *Object) bool {
	// Convert the pointer to an integer
	ptr := uintptr(unsafe.Pointer(obj))

	// Check if it's the nil immediate value
	return ptr == SPECIAL_NIL
}

// IsTrueImmediate returns true if the value is the immediate true value
func IsTrueImmediate(obj *Object) bool {
	// Convert the pointer to an integer
	ptr := uintptr(unsafe.Pointer(obj))

	// Check if it's the true immediate value
	return ptr == SPECIAL_TRUE
}

// IsFalseImmediate returns true if the value is the immediate false value
func IsFalseImmediate(obj *Object) bool {
	// Convert the pointer to an integer
	ptr := uintptr(unsafe.Pointer(obj))

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
