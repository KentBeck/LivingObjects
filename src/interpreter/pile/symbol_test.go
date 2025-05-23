package pile_test

import (
	"testing"

	"smalltalklsp/interpreter/pile"
)

func TestNewSymbol(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"Empty symbol", ""},
		{"Simple symbol", "hello"},
		{"Symbol with spaces", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sym := pile.NewSymbol(tt.value)
			if sym.Type() != pile.OBJ_SYMBOL {
				t.Errorf("NewSymbol(%q).Type() = %d, want %d", tt.value, sym.Type(), pile.OBJ_SYMBOL)
			}
			
			// Convert to Symbol and check value
			symObj := pile.ObjectToSymbol(sym)
			if symObj.Value != tt.value {
				t.Errorf("ObjectToSymbol(NewSymbol(%q)).Value = %q, want %q", tt.value, symObj.Value, tt.value)
			}
		})
	}
}

func TestSymbolToObjectAndBack(t *testing.T) {
	symObj := &pile.Symbol{
		Object: pile.Object{
			TypeField: pile.OBJ_SYMBOL,
		},
		Value: "hello",
	}
	
	obj := pile.SymbolToObject(symObj)
	
	if obj.Type() != pile.OBJ_SYMBOL {
		t.Errorf("SymbolToObject(symObj).Type() = %d, want %d", obj.Type(), pile.OBJ_SYMBOL)
	}
	
	backToSymbol := pile.ObjectToSymbol(obj)
	if backToSymbol.Value != "hello" {
		t.Errorf("ObjectToSymbol(SymbolToObject(symObj)).Value = %q, want %q", backToSymbol.Value, "hello")
	}
}

func TestSymbolString(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"Empty symbol", "", "#"},
		{"Simple symbol", "hello", "#hello"},
		{"Symbol with spaces", "hello world", "#hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sym := pile.ObjectToSymbol(pile.NewSymbol(tt.value))
			if sym.String() != tt.want {
				t.Errorf("ObjectToSymbol(NewSymbol(%q)).String() = %q, want %q", tt.value, sym.String(), tt.want)
			}
		})
	}
}

func TestSymbolGetValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"Empty symbol", ""},
		{"Simple symbol", "hello"},
		{"Symbol with spaces", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sym := pile.ObjectToSymbol(pile.NewSymbol(tt.value))
			if sym.GetValue() != tt.value {
				t.Errorf("ObjectToSymbol(NewSymbol(%q)).GetValue() = %q, want %q", tt.value, sym.GetValue(), tt.value)
			}
		})
	}
}

func TestSymbolSetValue(t *testing.T) {
	sym := pile.ObjectToSymbol(pile.NewSymbol("original"))
	sym.SetValue("modified")
	
	if sym.Value != "modified" {
		t.Errorf("After SetValue(%q), sym.Value = %q, want %q", "modified", sym.Value, "modified")
	}
}

func TestSymbolLength(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  int
	}{
		{"Empty symbol", "", 0},
		{"Simple symbol", "hello", 5},
		{"Symbol with spaces", "hello world", 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sym := pile.ObjectToSymbol(pile.NewSymbol(tt.value))
			if sym.Length() != tt.want {
				t.Errorf("ObjectToSymbol(NewSymbol(%q)).Length() = %d, want %d", tt.value, sym.Length(), tt.want)
			}
		})
	}
}

func TestSymbolEqual(t *testing.T) {
	sym1 := pile.ObjectToSymbol(pile.NewSymbol("hello"))
	sym2 := pile.ObjectToSymbol(pile.NewSymbol("hello"))
	sym3 := pile.ObjectToSymbol(pile.NewSymbol("world"))
	
	if !sym1.Equal(sym2) {
		t.Errorf("sym1.Equal(sym2) = false, want true")
	}
	
	if sym1.Equal(sym3) {
		t.Errorf("sym1.Equal(sym3) = true, want false")
	}
}

func TestGetSymbolValue(t *testing.T) {
	// Test with a symbol object
	symObj := pile.NewSymbol("hello")
	if pile.GetSymbolValue(symObj) != "hello" {
		t.Errorf("GetSymbolValue(symObj) = %q, want %q", pile.GetSymbolValue(symObj), "hello")
	}
	
	// Test with a non-symbol object
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("GetSymbolValue(nilObj) did not panic")
			}
		}()
		nilObj := pile.MakeNilImmediate()
		pile.GetSymbolValue(nilObj)
	}()
}