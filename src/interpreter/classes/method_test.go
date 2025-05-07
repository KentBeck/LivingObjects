package classes

import (
	"testing"

	"smalltalklsp/interpreter/core"
)

func TestNewMethod(t *testing.T) {
	// Create a class and selector
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")

	// Create a method
	method := NewMethod(selector, class)

	if method.Type() != core.OBJ_METHOD {
		t.Errorf("NewMethod(selector, class).Type() = %d, want %d", method.Type(), core.OBJ_METHOD)
	}

	methodObj := ObjectToMethod(method)

	if methodObj.Selector != selector {
		t.Errorf("methodObj.Selector = %v, want %v", methodObj.Selector, selector)
	}

	if methodObj.MethodClass != class {
		t.Errorf("methodObj.MethodClass = %v, want %v", methodObj.MethodClass, class)
	}

	if len(methodObj.Bytecodes) != 0 {
		t.Errorf("len(methodObj.Bytecodes) = %d, want 0", len(methodObj.Bytecodes))
	}

	if len(methodObj.Literals) != 0 {
		t.Errorf("len(methodObj.Literals) = %d, want 0", len(methodObj.Literals))
	}

	if len(methodObj.TempVarNames) != 0 {
		t.Errorf("len(methodObj.TempVarNames) = %d, want 0", len(methodObj.TempVarNames))
	}

	if methodObj.IsPrimitive {
		t.Errorf("methodObj.IsPrimitive = true, want false")
	}

	if methodObj.PrimitiveIndex != 0 {
		t.Errorf("methodObj.PrimitiveIndex = %d, want 0", methodObj.PrimitiveIndex)
	}
}

func TestMethodToObjectAndBack(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	obj := MethodToObject(method)

	if obj.Type() != core.OBJ_METHOD {
		t.Errorf("MethodToObject(method).Type() = %d, want %d", obj.Type(), core.OBJ_METHOD)
	}

	backToMethod := ObjectToMethod(obj)
	if backToMethod.Selector != selector {
		t.Errorf("ObjectToMethod(MethodToObject(method)).Selector = %v, want %v", backToMethod.Selector, selector)
	}

	// Test with wrong type
	strObj := StringToObject(NewString("test"))
	if ObjectToMethod(strObj) != nil {
		t.Errorf("ObjectToMethod(strObj) = %v, want nil", ObjectToMethod(strObj))
	}
}

func TestMethodString(t *testing.T) {
	// Create a method with a selector
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	if method.String() != "Method testMethod" {
		t.Errorf("method.String() = %q, want %q", method.String(), "Method testMethod")
	}

	// Create a method without a selector
	method.Selector = nil

	if method.String() != "a Method" {
		t.Errorf("method.String() = %q, want %q", method.String(), "a Method")
	}
}

func TestMethodGetBytecodes(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	// Set bytecodes
	bytecodes := []byte{1, 2, 3}
	method.Bytecodes = bytecodes

	// Get bytecodes
	result := method.GetBytecodes()

	if len(result) != len(bytecodes) {
		t.Errorf("len(method.GetBytecodes()) = %d, want %d", len(result), len(bytecodes))
	}

	for i, b := range result {
		if b != bytecodes[i] {
			t.Errorf("method.GetBytecodes()[%d] = %d, want %d", i, b, bytecodes[i])
		}
	}
}

func TestMethodSetBytecodes(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	// Set bytecodes
	bytecodes := []byte{1, 2, 3}
	method.SetBytecodes(bytecodes)

	if len(method.Bytecodes) != len(bytecodes) {
		t.Errorf("len(method.Bytecodes) = %d, want %d", len(method.Bytecodes), len(bytecodes))
	}

	for i, b := range method.Bytecodes {
		if b != bytecodes[i] {
			t.Errorf("method.Bytecodes[%d] = %d, want %d", i, b, bytecodes[i])
		}
	}
}

func TestMethodGetLiterals(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	// Set literals
	literals := []*core.Object{
		core.MakeIntegerImmediate(1),
		core.MakeIntegerImmediate(2),
	}
	method.Literals = literals

	// Get literals
	result := method.GetLiterals()

	if len(result) != len(literals) {
		t.Errorf("len(method.GetLiterals()) = %d, want %d", len(result), len(literals))
	}

	for i, lit := range result {
		if lit != literals[i] {
			t.Errorf("method.GetLiterals()[%d] = %v, want %v", i, lit, literals[i])
		}
	}
}

func TestMethodAddLiteral(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	// Add literals
	lit1 := core.MakeIntegerImmediate(1)
	lit2 := core.MakeIntegerImmediate(2)

	method.AddLiteral(lit1)
	method.AddLiteral(lit2)

	if len(method.Literals) != 2 {
		t.Errorf("len(method.Literals) = %d, want 2", len(method.Literals))
	}

	if method.Literals[0] != lit1 {
		t.Errorf("method.Literals[0] = %v, want %v", method.Literals[0], lit1)
	}

	if method.Literals[1] != lit2 {
		t.Errorf("method.Literals[1] = %v, want %v", method.Literals[1], lit2)
	}
}

func TestMethodGetSelector(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	// Get selector
	result := method.GetSelector()

	if result != selector {
		t.Errorf("method.GetSelector() = %v, want %v", result, selector)
	}
}

func TestMethodSetSelector(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector1 := NewSymbol("testMethod1")
	method := ObjectToMethod(NewMethod(selector1, class))

	// Set a new selector
	selector2 := NewSymbol("testMethod2")
	method.SetSelector(selector2)

	if method.Selector != selector2 {
		t.Errorf("method.Selector = %v, want %v", method.Selector, selector2)
	}
}

func TestMethodGetTempVarNames(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	// Set temp var names
	tempVars := []string{"temp1", "temp2"}
	method.TempVarNames = tempVars

	// Get temp var names
	result := method.GetTempVarNames()

	if len(result) != len(tempVars) {
		t.Errorf("len(method.GetTempVarNames()) = %d, want %d", len(result), len(tempVars))
	}

	for i, name := range result {
		if name != tempVars[i] {
			t.Errorf("method.GetTempVarNames()[%d] = %q, want %q", i, name, tempVars[i])
		}
	}
}

func TestMethodAddTempVarName(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	// Add temp var names
	method.AddTempVarName("temp1")
	method.AddTempVarName("temp2")

	if len(method.TempVarNames) != 2 {
		t.Errorf("len(method.TempVarNames) = %d, want 2", len(method.TempVarNames))
	}

	if method.TempVarNames[0] != "temp1" {
		t.Errorf("method.TempVarNames[0] = %q, want %q", method.TempVarNames[0], "temp1")
	}

	if method.TempVarNames[1] != "temp2" {
		t.Errorf("method.TempVarNames[1] = %q, want %q", method.TempVarNames[1], "temp2")
	}
}

func TestMethodGetMethodClass(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	// Get method class
	result := method.GetMethodClass()

	if result != class {
		t.Errorf("method.GetMethodClass() = %v, want %v", result, class)
	}
}

func TestMethodSetMethodClass(t *testing.T) {
	// Create a method
	class1 := NewClass("TestClass1", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class1))

	// Set a new method class
	class2 := NewClass("TestClass2", nil)
	method.SetMethodClass(class2)

	if method.MethodClass != class2 {
		t.Errorf("method.MethodClass = %v, want %v", method.MethodClass, class2)
	}
}

func TestMethodIsPrimitiveMethod(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	// Initially not a primitive
	if method.IsPrimitiveMethod() {
		t.Errorf("method.IsPrimitiveMethod() = true, want false")
	}

	// Set as primitive
	method.IsPrimitive = true

	if !method.IsPrimitiveMethod() {
		t.Errorf("method.IsPrimitiveMethod() = false, want true")
	}
}

func TestMethodSetPrimitive(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	// Set as primitive
	method.SetPrimitive(true)

	if !method.IsPrimitive {
		t.Errorf("method.IsPrimitive = false, want true")
	}

	// Set as non-primitive
	method.SetPrimitive(false)

	if method.IsPrimitive {
		t.Errorf("method.IsPrimitive = true, want false")
	}
}

func TestMethodGetPrimitiveIndex(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	// Set primitive index
	method.PrimitiveIndex = 42

	// Get primitive index
	result := method.GetPrimitiveIndex()

	if result != 42 {
		t.Errorf("method.GetPrimitiveIndex() = %d, want 42", result)
	}
}

func TestMethodSetPrimitiveIndex(t *testing.T) {
	// Create a method
	class := NewClass("TestClass", nil)
	selector := NewSymbol("testMethod")
	method := ObjectToMethod(NewMethod(selector, class))

	// Set primitive index
	method.SetPrimitiveIndex(42)

	if method.PrimitiveIndex != 42 {
		t.Errorf("method.PrimitiveIndex = %d, want 42", method.PrimitiveIndex)
	}
}
