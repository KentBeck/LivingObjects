// Package classes has been migrated to the pile package
// This package is kept for backward compatibility but is deprecated
package classes

import (
	"smalltalklsp/interpreter/pile"
)

// For backward compatibility
type (
	// Object is now in the pile package
	Object = pile.Object
	
	// Class is now in the pile package
	Class = pile.Class
	
	// Array is now in the pile package
	Array = pile.Array
	
	// Dictionary is now in the pile package
	Dictionary = pile.Dictionary
	
	// String is now in the pile package
	String = pile.String
	
	// Symbol is now in the pile package
	Symbol = pile.Symbol
	
	// Method is now in the pile package
	Method = pile.Method
	
	// Block is now in the pile package
	Block = pile.Block
	
	// ByteArray is now in the pile package
	ByteArray = pile.ByteArray
)

// ObjectType is now in the pile package
type ObjectType = pile.ObjectType

// Object type constants
const (
	OBJ_INSTANCE   = pile.OBJ_INSTANCE
	OBJ_CLASS      = pile.OBJ_CLASS
	OBJ_ARRAY      = pile.OBJ_ARRAY
	OBJ_DICTIONARY = pile.OBJ_DICTIONARY
	OBJ_STRING     = pile.OBJ_STRING
	OBJ_SYMBOL     = pile.OBJ_SYMBOL
	OBJ_METHOD     = pile.OBJ_METHOD
	OBJ_BLOCK      = pile.OBJ_BLOCK
	OBJ_BYTE_ARRAY = pile.OBJ_BYTE_ARRAY
)

// Compatibility functions
var (
	// NewMethod is now in the pile package
	NewMethod = pile.NewMethod
	
	// NewSymbol is now in the pile package
	NewSymbol = pile.NewSymbol
	
	// NewSymbolInternal is now in the pile package
	NewSymbolInternal = pile.NewSymbolInternal
	
	// NewStringInternal is now in the pile package
	NewStringInternal = pile.NewStringInternal
	
	// NewClass is now in the pile package
	NewClass = pile.NewClass
	
	// NewDictionaryInternal is now in the pile package
	NewDictionaryInternal = pile.NewDictionaryInternal
	
	// ObjectToClass is now in the pile package
	ObjectToClass = pile.ObjectToClass
	
	// ObjectToMethod is now in the pile package
	ObjectToMethod = pile.ObjectToMethod
	
	// ObjectToBlock is now in the pile package
	ObjectToBlock = pile.ObjectToBlock
	
	// ObjectToSymbol is now in the pile package
	ObjectToSymbol = pile.ObjectToSymbol
	
	// ObjectToString is now in the pile package
	ObjectToString = pile.ObjectToString
	
	// ObjectToArray is now in the pile package
	ObjectToArray = pile.ObjectToArray
	
	// ObjectToByteArray is now in the pile package
	ObjectToByteArray = pile.ObjectToByteArray
	
	// BlockToObject is now in the pile package
	BlockToObject = pile.BlockToObject
	
	// ClassToObject is now in the pile package
	ClassToObject = pile.ClassToObject
	
	// MethodToObject is now in the pile package
	MethodToObject = pile.MethodToObject
	
	// ByteArrayToObject is now in the pile package
	ByteArrayToObject = pile.ByteArrayToObject
	
	// SymbolToObject is now in the pile package
	SymbolToObject = pile.SymbolToObject
	
	// StringToObject is now in the pile package
	StringToObject = pile.StringToObject
	
	// ArrayToObject is now in the pile package
	ArrayToObject = pile.ArrayToObject
	
	// DictionaryToObject is now in the pile package
	DictionaryToObject = pile.DictionaryToObject
	
	// GetSymbolValue is now in the pile package
	GetSymbolValue = pile.GetSymbolValue
	
	// GetClassMethodDictionary is now in the pile package
	GetClassMethodDictionary = pile.GetClassMethodDictionary
	
	// GetClassName is now in the pile package
	GetClassName = pile.GetClassName
	
	// GetClassInstanceVarNames is now in the pile package
	GetClassInstanceVarNames = pile.GetClassInstanceVarNames
	
	// GetClassSuperClass is now in the pile package
	GetClassSuperClass = pile.GetClassSuperClass
)