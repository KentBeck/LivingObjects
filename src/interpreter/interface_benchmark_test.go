package main

import (
	"testing"
)

// BenchmarkGetClassDirect benchmarks the direct GetClass method
func BenchmarkGetClassDirect(b *testing.B) {
	vm := NewVM()
	obj := vm.NewInteger(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vm.GetClass(obj)
	}
}

// BenchmarkGetClassInterface benchmarks the interface-based GetClass method with an interface parameter
func BenchmarkGetClassInterface(b *testing.B) {
	vm := NewVM()
	obj := vm.NewInteger(42)
	var objInterface ObjectInterface = obj

	// We'll use the Type method from the interface to simulate using an interface
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// We need to cast back to *Object to use GetClass
		objAsPtr := objInterface.(*Object)
		_ = vm.GetClass(objAsPtr)
	}
}

// BenchmarkLookupMethodDirect benchmarks the direct lookupMethod method
func BenchmarkLookupMethodDirect(b *testing.B) {
	vm := NewVM()
	obj := vm.NewInteger(42)
	selector := NewSymbol("+")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vm.lookupMethod(obj, selector)
	}
}

// BenchmarkLookupMethodInterface benchmarks the lookupMethod method with interface usage
func BenchmarkLookupMethodInterface(b *testing.B) {
	vm := NewVM()
	obj := vm.NewInteger(42)
	selector := NewSymbol("+")
	var objInterface ObjectInterface = obj

	// We'll use the Type method from the interface to simulate using an interface
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// We need to cast back to *Object to use lookupMethod
		objAsPtr := objInterface.(*Object)
		_ = vm.lookupMethod(objAsPtr, selector)
	}
}

// BenchmarkObjectTypeDirectAccess benchmarks direct access to the type field
func BenchmarkObjectTypeDirectAccess(b *testing.B) {
	obj := &Object{
		type1: OBJ_INTEGER,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = obj.type1
	}
}

// BenchmarkObjectTypeInterface benchmarks access to the type field through the interface
func BenchmarkObjectTypeInterface(b *testing.B) {
	obj := &Object{
		type1: OBJ_INTEGER,
	}
	var objInterface ObjectInterface = obj

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = objInterface.Type()
	}
}

// BenchmarkObjectClassDirectAccess benchmarks direct access to the class field
func BenchmarkObjectClassDirectAccess(b *testing.B) {
	vm := NewVM()
	obj := &Object{
		class: vm.IntegerClass,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = obj.class
	}
}

// BenchmarkObjectClassInterface benchmarks access to the class field through the interface
func BenchmarkObjectClassInterface(b *testing.B) {
	vm := NewVM()
	obj := &Object{
		class: vm.IntegerClass,
	}
	var objInterface ObjectInterface = obj

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = objInterface.Class()
	}
}

// BenchmarkComplexOperationDirect benchmarks a complex operation using direct *Object
func BenchmarkComplexOperationDirect(b *testing.B) {
	vm := NewVM()
	obj := vm.NewInteger(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Get the class
		class := vm.GetClass(obj)

		// Check if it's an immediate value
		isImm := IsImmediate(obj)

		// Get the type - for immediate values, we need to check the tag bits
		var typ ObjectType
		if isImm {
			// For immediate values, the type is encoded in the tag bits
			if IsIntegerImmediate(obj) {
				typ = OBJ_INTEGER
			} else if IsTrueImmediate(obj) || IsFalseImmediate(obj) {
				typ = OBJ_BOOLEAN
			} else if IsNilImmediate(obj) {
				typ = OBJ_NIL
			}
		} else {
			// For regular objects, we can use the Type method
			typ = obj.Type()
		}

		// Combine the results to prevent the compiler from optimizing away the operations
		_ = class
		_ = isImm
		_ = typ
	}
}

// BenchmarkComplexOperationInterface benchmarks a complex operation using ObjectInterface
func BenchmarkComplexOperationInterface(b *testing.B) {
	vm := NewVM()
	obj := vm.NewInteger(42)
	var objInterface ObjectInterface = obj

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Cast to *Object for operations that require it
		objAsPtr := objInterface.(*Object)

		// Get the class
		class := vm.GetClass(objAsPtr)

		// Check if it's an immediate value
		isImm := IsImmediate(objAsPtr)

		// Get the type - for immediate values, we need to check the tag bits
		var typ ObjectType
		if isImm {
			// For immediate values, the type is encoded in the tag bits
			if IsIntegerImmediate(objAsPtr) {
				typ = OBJ_INTEGER
			} else if IsTrueImmediate(objAsPtr) || IsFalseImmediate(objAsPtr) {
				typ = OBJ_BOOLEAN
			} else if IsNilImmediate(objAsPtr) {
				typ = OBJ_NIL
			}
		} else {
			// For regular objects, we can use the Type method from the interface
			typ = objInterface.Type()
		}

		// Combine the results to prevent the compiler from optimizing away the operations
		_ = class
		_ = isImm
		_ = typ
	}
}
