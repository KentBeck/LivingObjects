package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/vm"
)

func TestByteArrayCreation(t *testing.T) {
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
			virtualMachine := vm.NewVM()
			byteArray := classes.ObjectToByteArray(virtualMachine.NewByteArray(tt.size))
			if byteArray.Size() != tt.want {
				t.Errorf("ByteArray.Size() = %d, want %d", byteArray.Size(), tt.want)
			}
			if byteArray.Type() != core.OBJ_BYTE_ARRAY {
				t.Errorf("ByteArray.Type() = %d, want %d", byteArray.Type(), core.OBJ_BYTE_ARRAY)
			}
			if len(byteArray.Bytes) != tt.size {
				t.Errorf("len(ByteArray.Bytes) = %d, want %d", len(byteArray.Bytes), tt.size)
			}
		})
	}
}

func TestByteArrayAtAndAtPut(t *testing.T) {
	virtualMachine := vm.NewVM()
	byteArray := classes.ObjectToByteArray(virtualMachine.NewByteArray(5))

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
