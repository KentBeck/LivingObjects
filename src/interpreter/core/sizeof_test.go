package core_test

import (
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/core"
)

func TestSizeOfObject(t *testing.T) {
	obj := &core.Object{}
	size := unsafe.Sizeof(*obj)

	// Assert that the size is less than or equal to 232 bytes
	const maxSize = 232
	t.Log("Object size is", size, "bytes")
	if size > maxSize {
		t.Errorf("Object size is %d bytes, which exceeds the maximum allowed size of %d bytes", size, maxSize)
	}
}
