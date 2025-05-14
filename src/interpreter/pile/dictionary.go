package pile

import (
	"fmt"
	"unsafe"
)

// Dictionary represents a Smalltalk dictionary object
type Dictionary struct {
	Object
	Entries map[string]*Object // later Object->Object
}

// newDictionary creates a new dictionary object without setting its class field
// This is a private helper function used by vm.NewDictionary
func NewDictionaryInternal() *Dictionary {
	return &Dictionary{
		Object: Object{
			TypeField: OBJ_DICTIONARY,
		},
		Entries: make(map[string]*Object),
	}
}

// NewDictionary creates a new dictionary object (deprecated - use vm.NewDictionary instead)
func NewDictionary() *Object {
	return DictionaryToObject(NewDictionaryInternal())
}

// DictionaryToObject converts a Dictionary to an Object
func DictionaryToObject(d *Dictionary) *Object {
	return (*Object)(unsafe.Pointer(d))
}

// ObjectToDictionary converts an Object to a Dictionary
func ObjectToDictionary(o ObjectInterface) *Dictionary {
	return (*Dictionary)(unsafe.Pointer(o.(*Object)))
}

// String returns a string representation of the dictionary object
func (d *Dictionary) String() string {
	return fmt.Sprintf("Dictionary(%d)", d.GetEntryCount())
}

// GetEntries returns the entries of the dictionary
func (d *Dictionary) GetEntries() map[string]*Object {
	return d.Entries
}

// GetEntry gets an entry from the dictionary
func (d *Dictionary) GetEntry(key string) *Object {
	return d.Entries[key]
}

// SetEntry sets an entry in the dictionary
func (d *Dictionary) SetEntry(key string, value *Object) {
	d.Entries[key] = value
}

// GetEntryCount returns the number of entries in the dictionary
func (d *Dictionary) GetEntryCount() int {
	return len(d.Entries)
}

// RemoveEntry removes an entry from the dictionary
func (d *Dictionary) RemoveEntry(key string) {
	delete(d.Entries, key)
}

// HasKey returns true if the dictionary has the given key
func (d *Dictionary) HasKey(key string) bool {
	_, ok := d.Entries[key]
	return ok
}

// Keys returns an array of all keys in the dictionary
func (d *Dictionary) Keys() *Array {
	keys := NewArray(len(d.Entries))
	i := 0
	for key := range d.Entries {
		keys.Elements[i] = StringToObject(NewString(key))
		i++
	}
	return keys
}

// Values returns an array of all values in the dictionary
func (d *Dictionary) Values() *Array {
	values := NewArray(len(d.Entries))
	i := 0
	for _, value := range d.Entries {
		values.Elements[i] = value
		i++
	}
	return values
}

// Do applies a function to each key-value pair in the dictionary
func (d *Dictionary) Do(fn func(string, *Object)) {
	for key, value := range d.Entries {
		fn(key, value)
	}
}

// Copy returns a copy of the dictionary
func (d *Dictionary) Copy() *Dictionary {
	newDict := &Dictionary{
		Object: Object{
			TypeField: OBJ_DICTIONARY,
		},
		Entries: make(map[string]*Object, len(d.Entries)),
	}
	for key, value := range d.Entries {
		newDict.Entries[key] = value
	}
	return newDict
}

// Merge merges another dictionary into this one
func (d *Dictionary) Merge(other *Dictionary) {
	for key, value := range other.Entries {
		d.Entries[key] = value
	}
}