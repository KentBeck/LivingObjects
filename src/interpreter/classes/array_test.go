package classes

import (
	"testing"

	"smalltalklsp/interpreter/core"
)

func TestNewArray(t *testing.T) {
	tests := []struct {
		name string
		size int
		want int
	}{
		{"Empty array", 0, 0},
		{"Small array", 5, 5},
		{"Large array", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			array := NewArray(tt.size)
			if array.Size() != tt.want {
				t.Errorf("NewArray(%d).Size() = %d, want %d", tt.size, array.Size(), tt.want)
			}
			if array.Type() != core.OBJ_ARRAY {
				t.Errorf("NewArray(%d).Type() = %d, want %d", tt.size, array.Type(), core.OBJ_ARRAY)
			}
			if len(array.Elements) != tt.size {
				t.Errorf("len(NewArray(%d).Elements) = %d, want %d", tt.size, len(array.Elements), tt.size)
			}
		})
	}
}

func TestArrayToObjectAndBack(t *testing.T) {
	array := NewArray(5)
	obj := ArrayToObject(array)

	if obj.Type() != core.OBJ_ARRAY {
		t.Errorf("ArrayToObject(array).Type() = %d, want %d", obj.Type(), core.OBJ_ARRAY)
	}

	backToArray := ObjectToArray(obj)
	if backToArray.Size() != 5 {
		t.Errorf("ObjectToArray(ArrayToObject(array)).Size() = %d, want 5", backToArray.Size())
	}
}

func TestArrayString(t *testing.T) {
	tests := []struct {
		name string
		size int
		want string
	}{
		{"Empty array", 0, "Array(0)"},
		{"Small array", 5, "Array(5)"},
		{"Large array", 100, "Array(100)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			array := NewArray(tt.size)
			if array.String() != tt.want {
				t.Errorf("NewArray(%d).String() = %s, want %s", tt.size, array.String(), tt.want)
			}
		})
	}
}
