package utils

import (
	"unsafe"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// MethodToObject converts a Method to an Object
func MethodToObject(m *classes.Method) *core.Object {
	return (*core.Object)(unsafe.Pointer(m))
}

// ObjectToMethod converts an Object to a Method
func ObjectToMethod(o *core.Object) *classes.Method {
	if o == nil || o.Type() != core.OBJ_METHOD {
		return nil
	}
	return (*classes.Method)(unsafe.Pointer(o))
}

// ClassToObject converts a Class to an Object
func ClassToObject(c *classes.Class) *core.Object {
	return (*core.Object)(unsafe.Pointer(c))
}

// ObjectToClass converts an Object to a Class
func ObjectToClass(o *core.Object) *classes.Class {
	return (*classes.Class)(unsafe.Pointer(o))
}

// StringToObject converts a String to an Object
func StringToObject(s *classes.String) *core.Object {
	return (*core.Object)(unsafe.Pointer(s))
}

// ObjectToString converts an Object to a String
func ObjectToString(o core.ObjectInterface) *classes.String {
	return (*classes.String)(unsafe.Pointer(o.(*core.Object)))
}

// SymbolToObject converts a Symbol to an Object
func SymbolToObject(s *classes.Symbol) *core.Object {
	return (*core.Object)(unsafe.Pointer(s))
}

// ObjectToSymbol converts an Object to a Symbol
func ObjectToSymbol(o core.ObjectInterface) *classes.Symbol {
	return (*classes.Symbol)(unsafe.Pointer(o.(*core.Object)))
}

// ArrayToObject converts an Array to an Object
func ArrayToObject(a *classes.Array) *core.Object {
	return (*core.Object)(unsafe.Pointer(a))
}

// ObjectToArray converts an Object to an Array
func ObjectToArray(o *core.Object) *classes.Array {
	return (*classes.Array)(unsafe.Pointer(o))
}

// BlockToObject converts a Block to an Object
func BlockToObject(b *classes.Block) *core.Object {
	return (*core.Object)(unsafe.Pointer(b))
}

// ObjectToBlock converts an Object to a Block
func ObjectToBlock(o *core.Object) *classes.Block {
	return (*classes.Block)(unsafe.Pointer(o))
}

// DictionaryToObject converts a Dictionary to an Object
func DictionaryToObject(d *classes.Dictionary) *core.Object {
	return (*core.Object)(unsafe.Pointer(d))
}

// ObjectToDictionary converts an Object to a Dictionary
func ObjectToDictionary(o core.ObjectInterface) *classes.Dictionary {
	return (*classes.Dictionary)(unsafe.Pointer(o.(*core.Object)))
}

// GetStringValue gets the string value of a string
func GetStringValue(obj *core.Object) string {
	if obj.Type() != core.OBJ_STRING {
		return ""
	}
	return ObjectToString(obj).GetValue()
}

// GetSymbolValue gets the value from a Symbol object
func GetSymbolValue(o core.ObjectInterface) string {
	if o.Type() == core.OBJ_SYMBOL {
		return ObjectToSymbol(o.(*core.Object)).GetValue()
	}
	panic("GetSymbolValue: not a symbol")
}
