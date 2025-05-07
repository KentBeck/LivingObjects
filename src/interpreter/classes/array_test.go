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

func TestArrayAtAndAtPut(t *testing.T) {
	array := NewArray(3)

	// Test AtPut
	nilObj := core.MakeNilImmediate()
	trueObj := core.MakeTrueImmediate()
	falseObj := core.MakeFalseImmediate()

	array.AtPut(0, nilObj)
	array.AtPut(1, trueObj)
	array.AtPut(2, falseObj)

	// Test At
	if array.At(0) != nilObj {
		t.Errorf("array.At(0) = %v, want %v", array.At(0), nilObj)
	}
	if array.At(1) != trueObj {
		t.Errorf("array.At(1) = %v, want %v", array.At(1), trueObj)
	}
	if array.At(2) != falseObj {
		t.Errorf("array.At(2) = %v, want %v", array.At(2), falseObj)
	}

	// Test out of bounds
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("array.At(-1) did not panic")
			}
		}()
		array.At(-1)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("array.At(3) did not panic")
			}
		}()
		array.At(3)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("array.AtPut(-1, nilObj) did not panic")
			}
		}()
		array.AtPut(-1, nilObj)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("array.AtPut(3, nilObj) did not panic")
			}
		}()
		array.AtPut(3, nilObj)
	}()
}

func TestArrayCopy(t *testing.T) {
	array := NewArray(3)
	nilObj := core.MakeNilImmediate()
	trueObj := core.MakeTrueImmediate()
	falseObj := core.MakeFalseImmediate()

	array.AtPut(0, nilObj)
	array.AtPut(1, trueObj)
	array.AtPut(2, falseObj)

	copy := array.Copy()

	// Check that the copy has the same size
	if copy.Size() != array.Size() {
		t.Errorf("copy.Size() = %d, want %d", copy.Size(), array.Size())
	}

	// Check that the copy has the same elements
	for i := 0; i < array.Size(); i++ {
		if copy.At(i) != array.At(i) {
			t.Errorf("copy.At(%d) = %v, want %v", i, copy.At(i), array.At(i))
		}
	}

	// Check that modifying the copy doesn't affect the original
	intObj := core.MakeIntegerImmediate(42)
	copy.AtPut(0, intObj)
	if array.At(0) == intObj {
		t.Errorf("array.At(0) = %v after modifying copy, want %v", array.At(0), nilObj)
	}
}

func TestArrayCollect(t *testing.T) {
	array := NewArray(3)
	for i := 0; i < 3; i++ {
		array.AtPut(i, core.MakeIntegerImmediate(int64(i)))
	}

	// Collect: multiply each element by 2
	result := array.Collect(func(obj *core.Object) *core.Object {
		val := core.GetIntegerImmediate(obj)
		return core.MakeIntegerImmediate(val * 2)
	})

	// Check the result
	if result.Size() != array.Size() {
		t.Errorf("result.Size() = %d, want %d", result.Size(), array.Size())
	}

	for i := 0; i < array.Size(); i++ {
		expected := core.MakeIntegerImmediate(int64(i) * 2)
		if core.GetIntegerImmediate(result.At(i)) != core.GetIntegerImmediate(expected) {
			t.Errorf("result.At(%d) = %d, want %d", i, core.GetIntegerImmediate(result.At(i)), core.GetIntegerImmediate(expected))
		}
	}
}

func TestArraySelect(t *testing.T) {
	array := NewArray(5)
	for i := 0; i < 5; i++ {
		array.AtPut(i, core.MakeIntegerImmediate(int64(i)))
	}

	// Select: keep only even numbers
	result := array.Select(func(obj *core.Object) bool {
		val := core.GetIntegerImmediate(obj)
		return val%2 == 0
	})

	// Check the result
	if result.Size() != 3 { // 0, 2, 4
		t.Errorf("result.Size() = %d, want 3", result.Size())
	}

	expected := []int64{0, 2, 4}
	for i := 0; i < result.Size(); i++ {
		if core.GetIntegerImmediate(result.At(i)) != expected[i] {
			t.Errorf("result.At(%d) = %d, want %d", i, core.GetIntegerImmediate(result.At(i)), expected[i])
		}
	}
}

func TestArrayReject(t *testing.T) {
	array := NewArray(5)
	for i := 0; i < 5; i++ {
		array.AtPut(i, core.MakeIntegerImmediate(int64(i)))
	}

	// Reject: remove even numbers
	result := array.Reject(func(obj *core.Object) bool {
		val := core.GetIntegerImmediate(obj)
		return val%2 == 0
	})

	// Check the result
	if result.Size() != 2 { // 1, 3
		t.Errorf("result.Size() = %d, want 2", result.Size())
	}

	expected := []int64{1, 3}
	for i := 0; i < result.Size(); i++ {
		if core.GetIntegerImmediate(result.At(i)) != expected[i] {
			t.Errorf("result.At(%d) = %d, want %d", i, core.GetIntegerImmediate(result.At(i)), expected[i])
		}
	}
}

func TestArrayDo(t *testing.T) {
	array := NewArray(3)
	for i := 0; i < 3; i++ {
		array.AtPut(i, core.MakeIntegerImmediate(int64(i)))
	}

	// Do: sum all elements
	var sum int64 = 0
	array.Do(func(obj *core.Object) {
		sum += core.GetIntegerImmediate(obj)
	})

	// Check the result
	if sum != 3 { // 0 + 1 + 2
		t.Errorf("sum = %d, want 3", sum)
	}
}

func TestArrayWithIndexDo(t *testing.T) {
	array := NewArray(3)
	for i := 0; i < 3; i++ {
		array.AtPut(i, core.MakeIntegerImmediate(int64(i*10)))
	}

	// WithIndexDo: check that indices match
	array.WithIndexDo(func(i int, obj *core.Object) {
		expected := int64(i * 10)
		actual := core.GetIntegerImmediate(obj)
		if actual != expected {
			t.Errorf("array.At(%d) = %d, want %d", i, actual, expected)
		}
	})
}
