package pile_test

import (
	"testing"

	"smalltalklsp/interpreter/pile"
)

func TestNewDictionary(t *testing.T) {
	dict := pile.NewDictionary()

	if dict.Type() != pile.OBJ_DICTIONARY {
		t.Errorf("NewDictionary().Type() = %d, want %d", dict.Type(), pile.OBJ_DICTIONARY)
	}

	dictObj := pile.ObjectToDictionary(dict)
	if dictObj.GetEntryCount() != 0 {
		t.Errorf("ObjectToDictionary(NewDictionary()).GetEntryCount() = %d, want 0", dictObj.GetEntryCount())
	}
}

func TestDictionaryToObjectAndBack(t *testing.T) {
	dictObj := &pile.Dictionary{
		Object: pile.Object{
			TypeField: pile.OBJ_DICTIONARY,
		},
		Entries: make(map[string]*pile.Object),
	}

	obj := pile.DictionaryToObject(dictObj)

	if obj.Type() != pile.OBJ_DICTIONARY {
		t.Errorf("DictionaryToObject(dictObj).Type() = %d, want %d", obj.Type(), pile.OBJ_DICTIONARY)
	}

	backToDict := pile.ObjectToDictionary(obj)
	if backToDict.GetEntryCount() != 0 {
		t.Errorf("ObjectToDictionary(DictionaryToObject(dictObj)).GetEntryCount() = %d, want 0", backToDict.GetEntryCount())
	}
}

func TestDictionaryString(t *testing.T) {
	tests := []struct {
		name     string
		entries  int
		expected string
	}{
		{"Empty dictionary", 0, "Dictionary(0)"},
		{"Dictionary with entries", 3, "Dictionary(3)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dict := pile.ObjectToDictionary(pile.NewDictionary())

			// Add entries if needed
			for i := 0; i < tt.entries; i++ {
				dict.SetEntry(string(rune('a'+i)), pile.MakeIntegerImmediate(int64(i)))
			}

			if dict.String() != tt.expected {
				t.Errorf("dict.String() = %q, want %q", dict.String(), tt.expected)
			}
		})
	}
}

func TestDictionaryGetEntries(t *testing.T) {
	dict := pile.ObjectToDictionary(pile.NewDictionary())

	// Add some entries
	dict.SetEntry("a", pile.MakeIntegerImmediate(1))
	dict.SetEntry("b", pile.MakeIntegerImmediate(2))

	entries := dict.GetEntries()

	if len(entries) != 2 {
		t.Errorf("len(dict.GetEntries()) = %d, want 2", len(entries))
	}

	if pile.GetIntegerImmediate(entries["a"]) != 1 {
		t.Errorf("entries[\"a\"] = %d, want 1", pile.GetIntegerImmediate(entries["a"]))
	}

	if pile.GetIntegerImmediate(entries["b"]) != 2 {
		t.Errorf("entries[\"b\"] = %d, want 2", pile.GetIntegerImmediate(entries["b"]))
	}
}

func TestDictionaryGetEntry(t *testing.T) {
	dict := pile.ObjectToDictionary(pile.NewDictionary())

	// Add some entries
	dict.SetEntry("a", pile.MakeIntegerImmediate(1))
	dict.SetEntry("b", pile.MakeIntegerImmediate(2))

	// Test existing entries
	if pile.GetIntegerImmediate(dict.GetEntry("a")) != 1 {
		t.Errorf("dict.GetEntry(\"a\") = %d, want 1", pile.GetIntegerImmediate(dict.GetEntry("a")))
	}

	if pile.GetIntegerImmediate(dict.GetEntry("b")) != 2 {
		t.Errorf("dict.GetEntry(\"b\") = %d, want 2", pile.GetIntegerImmediate(dict.GetEntry("b")))
	}

	// Test non-existent entry
	if dict.GetEntry("c") != nil {
		t.Errorf("dict.GetEntry(\"c\") = %v, want nil", dict.GetEntry("c"))
	}
}

func TestDictionarySetEntry(t *testing.T) {
	dict := pile.ObjectToDictionary(pile.NewDictionary())

	// Add a new entry
	dict.SetEntry("a", pile.MakeIntegerImmediate(1))

	if pile.GetIntegerImmediate(dict.GetEntry("a")) != 1 {
		t.Errorf("After SetEntry(\"a\", 1), dict.GetEntry(\"a\") = %d, want 1", pile.GetIntegerImmediate(dict.GetEntry("a")))
	}

	// Update an existing entry
	dict.SetEntry("a", pile.MakeIntegerImmediate(2))

	if pile.GetIntegerImmediate(dict.GetEntry("a")) != 2 {
		t.Errorf("After SetEntry(\"a\", 2), dict.GetEntry(\"a\") = %d, want 2", pile.GetIntegerImmediate(dict.GetEntry("a")))
	}
}

func TestDictionaryGetEntryCount(t *testing.T) {
	dict := pile.ObjectToDictionary(pile.NewDictionary())

	if dict.GetEntryCount() != 0 {
		t.Errorf("dict.GetEntryCount() = %d, want 0", dict.GetEntryCount())
	}

	// Add some entries
	dict.SetEntry("a", pile.MakeIntegerImmediate(1))
	dict.SetEntry("b", pile.MakeIntegerImmediate(2))

	if dict.GetEntryCount() != 2 {
		t.Errorf("dict.GetEntryCount() = %d, want 2", dict.GetEntryCount())
	}
}

func TestDictionaryRemoveEntry(t *testing.T) {
	dict := pile.ObjectToDictionary(pile.NewDictionary())

	// Add some entries
	dict.SetEntry("a", pile.MakeIntegerImmediate(1))
	dict.SetEntry("b", pile.MakeIntegerImmediate(2))

	// Remove an entry
	dict.RemoveEntry("a")

	if dict.GetEntryCount() != 1 {
		t.Errorf("After RemoveEntry(\"a\"), dict.GetEntryCount() = %d, want 1", dict.GetEntryCount())
	}

	if dict.GetEntry("a") != nil {
		t.Errorf("After RemoveEntry(\"a\"), dict.GetEntry(\"a\") = %v, want nil", dict.GetEntry("a"))
	}

	// Remove a non-existent entry (should not cause an error)
	dict.RemoveEntry("c")

	if dict.GetEntryCount() != 1 {
		t.Errorf("After RemoveEntry(\"c\"), dict.GetEntryCount() = %d, want 1", dict.GetEntryCount())
	}
}

func TestDictionaryHasKey(t *testing.T) {
	dict := pile.ObjectToDictionary(pile.NewDictionary())

	// Add some entries
	dict.SetEntry("a", pile.MakeIntegerImmediate(1))
	dict.SetEntry("b", pile.MakeIntegerImmediate(2))

	// Test existing keys
	if !dict.HasKey("a") {
		t.Errorf("dict.HasKey(\"a\") = false, want true")
	}

	if !dict.HasKey("b") {
		t.Errorf("dict.HasKey(\"b\") = false, want true")
	}

	// Test non-existent key
	if dict.HasKey("c") {
		t.Errorf("dict.HasKey(\"c\") = true, want false")
	}
}

func TestDictionaryKeys(t *testing.T) {
	dict := pile.ObjectToDictionary(pile.NewDictionary())

	// Add some entries
	dict.SetEntry("a", pile.MakeIntegerImmediate(1))
	dict.SetEntry("b", pile.MakeIntegerImmediate(2))

	keys := dict.Keys()

	if keys.Size() != 2 {
		t.Errorf("dict.Keys().Size() = %d, want 2", keys.Size())
	}

	// Check that all keys are present
	foundA := false
	foundB := false

	for i := 0; i < keys.Size(); i++ {
		key := pile.ObjectToString(keys.At(i)).Value
		if key == "a" {
			foundA = true
		} else if key == "b" {
			foundB = true
		}
	}

	if !foundA {
		t.Errorf("Key \"a\" not found in dict.Keys()")
	}

	if !foundB {
		t.Errorf("Key \"b\" not found in dict.Keys()")
	}
}

func TestDictionaryValues(t *testing.T) {
	dict := pile.ObjectToDictionary(pile.NewDictionary())

	// Add some entries
	dict.SetEntry("a", pile.MakeIntegerImmediate(1))
	dict.SetEntry("b", pile.MakeIntegerImmediate(2))

	values := dict.Values()

	if values.Size() != 2 {
		t.Errorf("dict.Values().Size() = %d, want 2", values.Size())
	}

	// Check that all values are present
	found1 := false
	found2 := false

	for i := 0; i < values.Size(); i++ {
		value := pile.GetIntegerImmediate(values.At(i))
		if value == 1 {
			found1 = true
		} else if value == 2 {
			found2 = true
		}
	}

	if !found1 {
		t.Errorf("Value 1 not found in dict.Values()")
	}

	if !found2 {
		t.Errorf("Value 2 not found in dict.Values()")
	}
}

func TestDictionaryDo(t *testing.T) {
	dict := pile.ObjectToDictionary(pile.NewDictionary())

	// Add some entries
	dict.SetEntry("a", pile.MakeIntegerImmediate(int64(1)))
	dict.SetEntry("b", pile.MakeIntegerImmediate(int64(2)))

	// Use Do to sum all values
	var sum int64 = 0
	dict.Do(func(key string, value *pile.Object) {
		sum += pile.GetIntegerImmediate(value)
	})

	if sum != 3 { // 1 + 2
		t.Errorf("sum = %d, want 3", sum)
	}
}