package pile_test

import (
	"testing"

	"smalltalklsp/interpreter/pile"
)

func TestNewByteArrayInternal(t *testing.T) {
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
			ba := pile.NewByteArrayInternal(tt.size)
			if ba.Size() != tt.want {
				t.Errorf("NewByteArrayInternal(%d).Size() = %d, want %d", tt.size, ba.Size(), tt.want)
			}
			if ba.Type() != pile.OBJ_BYTE_ARRAY {
				t.Errorf("NewByteArrayInternal(%d).Type() = %d, want %d", tt.size, ba.Type(), pile.OBJ_BYTE_ARRAY)
			}
			if len(ba.Bytes) != tt.size {
				t.Errorf("len(NewByteArrayInternal(%d).Bytes) = %d, want %d", tt.size, len(ba.Bytes), tt.size)
			}
		})
	}
}

func TestByteArrayToObjectAndBack(t *testing.T) {
	ba := pile.NewByteArrayInternal(5)
	obj := pile.ByteArrayToObject(ba)

	if obj.Type() != pile.OBJ_BYTE_ARRAY {
		t.Errorf("ByteArrayToObject(ba).Type() = %d, want %d", obj.Type(), pile.OBJ_BYTE_ARRAY)
	}

	backToBA := pile.ObjectToByteArray(obj)
	if backToBA.Size() != 5 {
		t.Errorf("ObjectToByteArray(ByteArrayToObject(ba)).Size() = %d, want 5", backToBA.Size())
	}
}

func TestByteArrayString(t *testing.T) {
	tests := []struct {
		name string
		size int
		want string
	}{
		{"Empty byte array", 0, "ByteArray(0)"},
		{"Small byte array", 5, "ByteArray(5)"},
		{"Large byte array", 100, "ByteArray(100)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ba := pile.NewByteArrayInternal(tt.size)
			if ba.String() != tt.want {
				t.Errorf("NewByteArrayInternal(%d).String() = %s, want %s", tt.size, ba.String(), tt.want)
			}
		})
	}
}

func TestByteArrayAtAndAtPut(t *testing.T) {
	ba := pile.NewByteArrayInternal(5)
	
	// Test AtPut and At
	ba.AtPut(0, 42)
	if ba.At(0) != 42 {
		t.Errorf("After ba.AtPut(0, 42), ba.At(0) = %d, want 42", ba.At(0))
	}
	
	// Test out of bounds
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected ba.At(-1) to panic")
			}
		}()
		ba.At(-1)
	}()
	
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected ba.At(5) to panic")
			}
		}()
		ba.At(5)
	}()
	
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected ba.AtPut(-1, 42) to panic")
			}
		}()
		ba.AtPut(-1, 42)
	}()
	
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected ba.AtPut(5, 42) to panic")
			}
		}()
		ba.AtPut(5, 42)
	}()
}