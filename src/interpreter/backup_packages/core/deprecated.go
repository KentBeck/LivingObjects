// Package core has been migrated to the pile package
// This package is kept for backward compatibility but is deprecated
package core

import (
	"smalltalklsp/interpreter/pile"
)

// For backward compatibility
type (
	// Object is now in the pile package
	Object = pile.Object
	
	// Class is now in the pile package
	Class = pile.Class
	
	// ObjectInterface is now in the pile package
	ObjectInterface = pile.ObjectInterface
	
	// ObjectMemory is now in the pile package
	ObjectMemory = pile.ObjectMemory
)

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
	OBJ_NIL        = pile.OBJ_NIL
	OBJ_INTEGER    = pile.OBJ_INTEGER
)

// Tag bit constants
const (
	TAG_POINTER = pile.TAG_POINTER
	TAG_SPECIAL = pile.TAG_SPECIAL
	TAG_FLOAT   = pile.TAG_FLOAT
	TAG_INTEGER = pile.TAG_INTEGER
	TAG_MASK    = pile.TAG_MASK
)

// Compatibility functions
var (
	// MakeNilImmediate is now in the pile package
	MakeNilImmediate = pile.MakeNilImmediate
	
	// MakeTrueImmediate is now in the pile package
	MakeTrueImmediate = pile.MakeTrueImmediate
	
	// MakeFalseImmediate is now in the pile package
	MakeFalseImmediate = pile.MakeFalseImmediate
	
	// MakeIntegerImmediate is now in the pile package
	MakeIntegerImmediate = pile.MakeIntegerImmediate
	
	// MakeFloatImmediate is now in the pile package
	MakeFloatImmediate = pile.MakeFloatImmediate
	
	// GetIntegerImmediate is now in the pile package
	GetIntegerImmediate = pile.GetIntegerImmediate
	
	// GetFloatImmediate is now in the pile package
	GetFloatImmediate = pile.GetFloatImmediate
	
	// IsNilImmediate is now in the pile package
	IsNilImmediate = pile.IsNilImmediate
	
	// IsTrueImmediate is now in the pile package
	IsTrueImmediate = pile.IsTrueImmediate
	
	// IsFalseImmediate is now in the pile package
	IsFalseImmediate = pile.IsFalseImmediate
	
	// IsIntegerImmediate is now in the pile package
	IsIntegerImmediate = pile.IsIntegerImmediate
	
	// IsFloatImmediate is now in the pile package
	IsFloatImmediate = pile.IsFloatImmediate
	
	// IsImmediate is now in the pile package
	IsImmediate = pile.IsImmediate
	
	// NewClass is now in the pile package
	NewClass = pile.NewClass
	
	// NewNil is now in the pile package
	NewNil = pile.NewNil
	
	// NewBoolean is now in the pile package
	NewBoolean = pile.NewBoolean
	
	// NewSymbol is now in the pile package
	NewSymbol = pile.NewSymbol
	
	// NewInstance is now in the pile package
	NewInstance = pile.NewInstance
	
	// ObjectToArray is now in the pile package
	ObjectToArray = pile.ObjectToArray
	
	// ClassToObject is now in the pile package
	ClassToObject = pile.ClassToObject
	
	// MethodToObject is now in the pile package
	MethodToObject = pile.MethodToObject
	
	// NewObjectMemory is now in the pile package
	NewObjectMemory = pile.NewObjectMemory
)