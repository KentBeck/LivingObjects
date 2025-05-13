package classes

import (
	"testing"

	"smalltalklsp/interpreter/core"
)

func TestNewByteArray(t *testing.T) {
	tests := []struct {
		name string
		size int
		want int
	}{
		{"Empty byte array", 0, 0},
		{"Small byte array", 5, 5},
		{"Large byte array", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			byteArray := NewByteArray(tt.size)
			if byteArray.Size() != tt.want {
				t.Errorf("NewByteArray(%d).Size() = %d, want %d", tt.size, byteArray.Size(), tt.want)
			}
			if byteArray.Type() != core.OBJ_BYTE_ARRAY {
				t.Errorf("NewByteArray(%d).Type() = %d, want %d", tt.size, byteArray.Type(), core.OBJ_BYTE_ARRAY)
			}
			if len(byteArray.Bytes) != tt.size {
				t.Errorf("len(NewByteArray(%d).Bytes) = %d, want %d", tt.size, len(byteArray.Bytes), tt.size)
			}
		})
	}
}

func TestByteArrayAtAndAtPut(t *testing.T) {
	byteArray := NewByteArray(5)

	// Test initial values (should be 0)
	for i := 0; i < byteArray.Size(); i++ {
		if byteArray.At(i) != 0 {
			t.Errorf("Initial value at index %d = %d, want 0", i, byteArray.At(i))
		}
	}

	// Test setting and getting values
	testValues := []byte{42, 255, 0, 127, 128}
	for i, value := range testValues {
		byteArray.AtPut(i, value)
	}

	for i, want := range testValues {
		got := byteArray.At(i)
		if got != want {
			t.Errorf("byteArray.At(%d) = %d, want %d", i, got, want)
		}
	}

	// Test bounds checking
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for index out of bounds, but got none")
			}
		}()
		byteArray.At(-1)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for index out of bounds, but got none")
			}
		}()
		byteArray.At(byteArray.Size())
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for index out of bounds, but got none")
			}
		}()
		byteArray.AtPut(-1, 0)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for index out of bounds, but got none")
			}
		}()
		byteArray.AtPut(byteArray.Size(), 0)
	}()
}

func TestByteArrayCopy(t *testing.T) {
	byteArray := NewByteArray(5)
	for i := 0; i < byteArray.Size(); i++ {
		byteArray.AtPut(i, byte(i+1))
	}

	// Test full copy
	copy := byteArray.Copy()
	if copy.Size() != byteArray.Size() {
		t.Errorf("copy.Size() = %d, want %d", copy.Size(), byteArray.Size())
	}
	for i := 0; i < byteArray.Size(); i++ {
		if copy.At(i) != byteArray.At(i) {
			t.Errorf("copy.At(%d) = %d, want %d", i, copy.At(i), byteArray.At(i))
		}
	}

	// Modify the copy and check that the original is unchanged
	copy.AtPut(0, 99)
	if byteArray.At(0) == 99 {
		t.Errorf("Original byte array was modified when copy was modified")
	}
}

func TestByteArrayCopyFrom(t *testing.T) {
	byteArray := NewByteArray(5)
	for i := 0; i < byteArray.Size(); i++ {
		byteArray.AtPut(i, byte(i+1))
	}

	// Test partial copy
	partial := byteArray.CopyFrom(1, 3)
	if partial.Size() != 3 {
		t.Errorf("partial.Size() = %d, want 3", partial.Size())
	}
	for i := 0; i < partial.Size(); i++ {
		if partial.At(i) != byteArray.At(i+1) {
			t.Errorf("partial.At(%d) = %d, want %d", i, partial.At(i), byteArray.At(i+1))
		}
	}

	// Test bounds checking
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for start index out of bounds, but got none")
			}
		}()
		byteArray.CopyFrom(-1, 3)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for end index out of bounds, but got none")
			}
		}()
		byteArray.CopyFrom(1, byteArray.Size())
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for end < start, but got none")
			}
		}()
		byteArray.CopyFrom(3, 1)
	}()
}

func TestByteArrayUint32AtAndUint32AtPut(t *testing.T) {
	byteArray := NewByteArray(8)

	// Test setting and getting uint32 values
	testValues := []uint32{42, 0xFFFFFFFF}
	byteArray.Uint32AtPut(0, testValues[0])
	byteArray.Uint32AtPut(4, testValues[1])

	for i, want := range testValues {
		got := byteArray.Uint32At(i * 4)
		if got != want {
			t.Errorf("byteArray.Uint32At(%d) = %d, want %d", i*4, got, want)
		}
	}

	// Test bounds checking
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for index out of bounds, but got none")
			}
		}()
		byteArray.Uint32At(-1)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for index out of bounds, but got none")
			}
		}()
		byteArray.Uint32At(byteArray.Size() - 3)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for index out of bounds, but got none")
			}
		}()
		byteArray.Uint32AtPut(-1, 0)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for index out of bounds, but got none")
			}
		}()
		byteArray.Uint32AtPut(byteArray.Size() - 3, 0)
	}()
}
