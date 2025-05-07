package classes

import (
	"testing"

	"smalltalklsp/interpreter/core"
)

func TestNewString(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"Empty string", ""},
		{"Simple string", "hello"},
		{"String with spaces", "hello world"},
		{"String with special chars", "hello\nworld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := NewString(tt.value)
			if str.Value != tt.value {
				t.Errorf("NewString(%q).Value = %q, want %q", tt.value, str.Value, tt.value)
			}
			if str.Type() != core.OBJ_STRING {
				t.Errorf("NewString(%q).Type() = %d, want %d", tt.value, str.Type(), core.OBJ_STRING)
			}
		})
	}
}

func TestStringToObjectAndBack(t *testing.T) {
	str := NewString("hello")
	obj := StringToObject(str)

	if obj.Type() != core.OBJ_STRING {
		t.Errorf("StringToObject(str).Type() = %d, want %d", obj.Type(), core.OBJ_STRING)
	}

	backToString := ObjectToString(obj)
	if backToString.Value != "hello" {
		t.Errorf("ObjectToString(StringToObject(str)).Value = %q, want %q", backToString.Value, "hello")
	}
}

func TestStringString(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"Empty string", "", "''"},
		{"Simple string", "hello", "'hello'"},
		{"String with spaces", "hello world", "'hello world'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := NewString(tt.value)
			if str.String() != tt.want {
				t.Errorf("NewString(%q).String() = %q, want %q", tt.value, str.String(), tt.want)
			}
		})
	}
}

func TestStringGetValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"Empty string", ""},
		{"Simple string", "hello"},
		{"String with spaces", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := NewString(tt.value)
			if str.GetValue() != tt.value {
				t.Errorf("NewString(%q).GetValue() = %q, want %q", tt.value, str.GetValue(), tt.value)
			}
		})
	}
}

func TestStringSetValue(t *testing.T) {
	str := NewString("original")
	str.SetValue("modified")

	if str.Value != "modified" {
		t.Errorf("After SetValue(%q), str.Value = %q, want %q", "modified", str.Value, "modified")
	}
}

func TestStringLength(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  int
	}{
		{"Empty string", "", 0},
		{"Simple string", "hello", 5},
		{"String with spaces", "hello world", 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := NewString(tt.value)
			if str.Length() != tt.want {
				t.Errorf("NewString(%q).Length() = %d, want %d", tt.value, str.Length(), tt.want)
			}
		})
	}
}

func TestStringCharAt(t *testing.T) {
	str := NewString("hello")

	tests := []struct {
		name  string
		index int
		want  byte
	}{
		{"First char", 0, 'h'},
		{"Middle char", 2, 'l'},
		{"Last char", 4, 'o'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if str.CharAt(tt.index) != tt.want {
				t.Errorf("str.CharAt(%d) = %c, want %c", tt.index, str.CharAt(tt.index), tt.want)
			}
		})
	}

	// Test out of bounds
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("str.CharAt(-1) did not panic")
			}
		}()
		str.CharAt(-1)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("str.CharAt(5) did not panic")
			}
		}()
		str.CharAt(5)
	}()
}

func TestStringSubstring(t *testing.T) {
	str := NewString("hello world")

	tests := []struct {
		name  string
		start int
		end   int
		want  string
	}{
		{"Full string", 0, 11, "hello world"},
		{"First word", 0, 5, "hello"},
		{"Second word", 6, 11, "world"},
		{"Middle chars", 2, 9, "llo wor"},
		{"Empty substring", 5, 5, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := str.Substring(tt.start, tt.end)
			if result.Value != tt.want {
				t.Errorf("str.Substring(%d, %d) = %q, want %q", tt.start, tt.end, result.Value, tt.want)
			}
		})
	}

	// Test invalid ranges
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("str.Substring(-1, 5) did not panic")
			}
		}()
		str.Substring(-1, 5)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("str.Substring(5, 12) did not panic")
			}
		}()
		str.Substring(5, 12)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("str.Substring(6, 3) did not panic")
			}
		}()
		str.Substring(6, 3)
	}()
}

func TestStringConcat(t *testing.T) {
	str1 := NewString("hello")
	str2 := NewString(" world")

	result := str1.Concat(str2)

	if result.Value != "hello world" {
		t.Errorf("str1.Concat(str2) = %q, want %q", result.Value, "hello world")
	}

	// Check that the original strings are unchanged
	if str1.Value != "hello" {
		t.Errorf("After concat, str1.Value = %q, want %q", str1.Value, "hello")
	}
	if str2.Value != " world" {
		t.Errorf("After concat, str2.Value = %q, want %q", str2.Value, " world")
	}
}

func TestStringEqual(t *testing.T) {
	str1 := NewString("hello")
	str2 := NewString("hello")
	str3 := NewString("world")

	if !str1.Equal(str2) {
		t.Errorf("str1.Equal(str2) = false, want true")
	}

	if str1.Equal(str3) {
		t.Errorf("str1.Equal(str3) = true, want false")
	}
}

func TestGetStringValue(t *testing.T) {
	// Test with a string object
	strObj := StringToObject(NewString("hello"))
	if GetStringValue(strObj) != "hello" {
		t.Errorf("GetStringValue(strObj) = %q, want %q", GetStringValue(strObj), "hello")
	}

	// Test that GetStringValue panics with an immediate value
	t.Run("panic with immediate value", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("GetStringValue with immediate value did not panic")
			}
		}()

		intObj := core.MakeIntegerImmediate(int64(42))
		GetStringValue(intObj) // This should panic
	})

	// Test that GetStringValue panics with a non-string object
	t.Run("panic with non-string object", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("GetStringValue with non-string object did not panic")
			}
		}()

		// Create a simple non-string object
		nonStringObj := &core.Object{
			TypeField: core.OBJ_INSTANCE, // Not a string
		}

		GetStringValue(nonStringObj) // This should panic
	})
}
