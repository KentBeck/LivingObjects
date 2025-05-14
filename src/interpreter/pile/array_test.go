package pile_test

import (
	"testing"

	"smalltalklsp/interpreter/pile"
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
			array := pile.NewArray(tt.size)
			if array.Size() != tt.want {
				t.Errorf("NewArray(%d).Size() = %d, want %d", tt.size, array.Size(), tt.want)
			}
			if array.Type() != pile.OBJ_ARRAY {
				t.Errorf("NewArray(%d).Type() = %d, want %d", tt.size, array.Type(), pile.OBJ_ARRAY)
			}
			if len(array.Elements) != tt.size {
				t.Errorf("len(NewArray(%d).Elements) = %d, want %d", tt.size, len(array.Elements), tt.size)
			}
		})
	}
}

func TestArrayToObjectAndBack(t *testing.T) {
	array := pile.NewArray(5)
	obj := pile.ArrayToObject(array)

	if obj.Type() != pile.OBJ_ARRAY {
		t.Errorf("ArrayToObject(array).Type() = %d, want %d", obj.Type(), pile.OBJ_ARRAY)
	}

	backToArray := pile.ObjectToArray(obj)
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
			array := pile.NewArray(tt.size)
			if array.String() != tt.want {
				t.Errorf("NewArray(%d).String() = %s, want %s", tt.size, array.String(), tt.want)
			}
		})
	}
}

func TestArrayAtAndAtPut(t *testing.T) {
	array := pile.NewArray(5)
	intObj := pile.MakeIntegerImmediate(42)
	
	// Test AtPut and At
	array.AtPut(0, intObj)
	if array.At(0) != intObj {
		t.Errorf("After array.AtPut(0, intObj), array.At(0) != intObj")
	}
	
	// Test out of bounds
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected array.At(-1) to panic")
			}
		}()
		array.At(-1)
	}()
	
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected array.At(5) to panic")
			}
		}()
		array.At(5)
	}()
	
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected array.AtPut(-1, intObj) to panic")
			}
		}()
		array.AtPut(-1, intObj)
	}()
	
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected array.AtPut(5, intObj) to panic")
			}
		}()
		array.AtPut(5, intObj)
	}()
}