package main

import (
	"math"
	"unsafe"
)

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
