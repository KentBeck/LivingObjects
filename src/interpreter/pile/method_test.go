package pile_test

import (
	"testing"

	"smalltalklsp/interpreter/pile"
)

func TestNewMethod(t *testing.T) {
	// Create a simple class for testing
	class := pile.NewClass("TestClass", nil)
	
	// Create a selector symbol
	selector := pile.NewSymbol("testMethod")
	
	// Create a method
	method := pile.NewMethod(selector, class)
	
	// Check the method type
	if method.Type() != pile.OBJ_METHOD {
		t.Errorf("NewMethod().Type() = %d, want %d", method.Type(), pile.OBJ_METHOD)
	}
	
	// Convert to Method and check properties
	methodObj := pile.ObjectToMethod(method)
	if methodObj == nil {
		t.Fatal("ObjectToMethod(method) returned nil")
	}
	
	// Check method class
	if methodObj.GetMethodClass() != class {
		t.Errorf("methodObj.GetMethodClass() != class")
	}
	
	// Check method selector
	if methodObj.GetSelector() != selector {
		t.Errorf("methodObj.GetSelector() != selector")
	}
	
	// Check bytecodes initialization
	if len(methodObj.GetBytecodes()) != 0 {
		t.Errorf("len(methodObj.GetBytecodes()) = %d, want 0", len(methodObj.GetBytecodes()))
	}
	
	// Check literals initialization
	if len(methodObj.GetLiterals()) != 0 {
		t.Errorf("len(methodObj.GetLiterals()) = %d, want 0", len(methodObj.GetLiterals()))
	}
	
	// Check temp vars initialization
	if len(methodObj.GetTempVarNames()) != 0 {
		t.Errorf("len(methodObj.GetTempVarNames()) = %d, want 0", len(methodObj.GetTempVarNames()))
	}
}

func TestMethodToObjectAndBack(t *testing.T) {
	// Create a method directly
	methodObj := &pile.Method{
		Object: pile.Object{
			TypeField: pile.OBJ_METHOD,
		},
		Selector:       pile.NewSymbol("testMethod"),
		MethodClass:    pile.NewClass("TestClass", nil),
		Bytecodes:      make([]byte, 0),
		Literals:       make([]*pile.Object, 0),
		TempVarNames:   make([]string, 0),
		IsPrimitive:    false,
		PrimitiveIndex: 0,
	}
	
	// Convert to Object and back
	obj := pile.MethodToObject(methodObj)
	
	if obj.Type() != pile.OBJ_METHOD {
		t.Errorf("MethodToObject(methodObj).Type() = %d, want %d", obj.Type(), pile.OBJ_METHOD)
	}
	
	backToMethod := pile.ObjectToMethod(obj)
	if backToMethod == nil {
		t.Fatal("ObjectToMethod(MethodToObject(methodObj)) returned nil")
	}
}

func TestMethodString(t *testing.T) {
	// Create a method with a selector
	class := pile.NewClass("TestClass", nil)
	selector := pile.NewSymbol("testMethod")
	method := pile.ObjectToMethod(pile.NewMethod(selector, class))
	
	expected := "Method testMethod"
	if method.String() != expected {
		t.Errorf("method.String() = %q, want %q", method.String(), expected)
	}
}

func TestMethodBytecodes(t *testing.T) {
	// Create a method
	class := pile.NewClass("TestClass", nil)
	selector := pile.NewSymbol("testMethod")
	method := pile.ObjectToMethod(pile.NewMethod(selector, class))
	
	// Set bytecodes
	bytecodes := []byte{1, 2, 3, 4, 5}
	method.SetBytecodes(bytecodes)
	
	// Check bytecodes
	if len(method.GetBytecodes()) != len(bytecodes) {
		t.Errorf("len(method.GetBytecodes()) = %d, want %d", len(method.GetBytecodes()), len(bytecodes))
	}
	
	for i, b := range method.GetBytecodes() {
		if b != bytecodes[i] {
			t.Errorf("method.GetBytecodes()[%d] = %d, want %d", i, b, bytecodes[i])
		}
	}
}

func TestMethodLiterals(t *testing.T) {
	// Create a method
	class := pile.NewClass("TestClass", nil)
	selector := pile.NewSymbol("testMethod")
	method := pile.ObjectToMethod(pile.NewMethod(selector, class))
	
	// Add literals
	literal1 := pile.MakeIntegerImmediate(42)
	literal2 := pile.StringToObject(pile.NewString("hello"))
	
	method.AddLiteral(literal1)
	method.AddLiteral(literal2)
	
	// Check literals
	literals := method.GetLiterals()
	if len(literals) != 2 {
		t.Errorf("len(method.GetLiterals()) = %d, want 2", len(literals))
	}
	
	if literals[0] != literal1 {
		t.Errorf("method.GetLiterals()[0] != literal1")
	}
	
	if literals[1] != literal2 {
		t.Errorf("method.GetLiterals()[1] != literal2")
	}
}

func TestMethodTempVarNames(t *testing.T) {
	// Create a method
	class := pile.NewClass("TestClass", nil)
	selector := pile.NewSymbol("testMethod")
	method := pile.ObjectToMethod(pile.NewMethod(selector, class))
	
	// Add temp var names
	method.AddTempVarName("temp1")
	method.AddTempVarName("temp2")
	
	// Check temp var names
	tempVarNames := method.GetTempVarNames()
	if len(tempVarNames) != 2 {
		t.Errorf("len(method.GetTempVarNames()) = %d, want 2", len(tempVarNames))
	}
	
	if tempVarNames[0] != "temp1" {
		t.Errorf("method.GetTempVarNames()[0] = %q, want %q", tempVarNames[0], "temp1")
	}
	
	if tempVarNames[1] != "temp2" {
		t.Errorf("method.GetTempVarNames()[1] = %q, want %q", tempVarNames[1], "temp2")
	}
}

func TestMethodIsPrimitive(t *testing.T) {
	// Create a method
	class := pile.NewClass("TestClass", nil)
	selector := pile.NewSymbol("testMethod")
	method := pile.ObjectToMethod(pile.NewMethod(selector, class))
	
	// Check default value
	if method.IsPrimitiveMethod() {
		t.Errorf("method.IsPrimitiveMethod() = true, want false")
	}
	
	// Set as primitive
	method.SetPrimitive(true)
	if !method.IsPrimitiveMethod() {
		t.Errorf("After SetPrimitive(true), method.IsPrimitiveMethod() = false, want true")
	}
	
	// Set primitive index
	method.SetPrimitiveIndex(42)
	if method.GetPrimitiveIndex() != 42 {
		t.Errorf("After SetPrimitiveIndex(42), method.GetPrimitiveIndex() = %d, want 42", method.GetPrimitiveIndex())
	}
}