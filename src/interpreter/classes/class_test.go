package classes

import (
	"testing"

	"smalltalklsp/interpreter/core"
)

func TestNewClass(t *testing.T) {
	// Test with nil superclass
	class1 := NewClass("TestClass1", nil)

	if class1.Type() != core.OBJ_CLASS {
		t.Errorf("NewClass(\"TestClass1\", nil).Type() = %d, want %d", class1.Type(), core.OBJ_CLASS)
	}

	if class1.Name != "TestClass1" {
		t.Errorf("NewClass(\"TestClass1\", nil).Name = %q, want %q", class1.Name, "TestClass1")
	}

	if class1.SuperClass != nil {
		t.Errorf("NewClass(\"TestClass1\", nil).SuperClass = %v, want nil", class1.SuperClass)
	}

	// Test with a superclass
	class2 := NewClass("TestClass2", class1)

	if class2.Type() != core.OBJ_CLASS {
		t.Errorf("NewClass(\"TestClass2\", class1).Type() = %d, want %d", class2.Type(), core.OBJ_CLASS)
	}

	if class2.Name != "TestClass2" {
		t.Errorf("NewClass(\"TestClass2\", class1).Name = %q, want %q", class2.Name, "TestClass2")
	}

	if class2.SuperClass == nil {
		t.Errorf("NewClass(\"TestClass2\", class1).SuperClass is nil, want non-nil")
	}
}

func TestClassToObjectAndBack(t *testing.T) {
	class := NewClass("TestClass", nil)
	obj := ClassToObject(class)

	if obj.Type() != core.OBJ_CLASS {
		t.Errorf("ClassToObject(class).Type() = %d, want %d", obj.Type(), core.OBJ_CLASS)
	}

	backToClass := ObjectToClass(obj)
	if backToClass.Name != "TestClass" {
		t.Errorf("ObjectToClass(ClassToObject(class)).Name = %q, want %q", backToClass.Name, "TestClass")
	}
}

func TestClassString(t *testing.T) {
	class := NewClass("TestClass", nil)

	if class.String() != "Class TestClass" {
		t.Errorf("class.String() = %q, want %q", class.String(), "Class TestClass")
	}
}

func TestClassGetName(t *testing.T) {
	class := NewClass("TestClass", nil)

	if class.GetName() != "TestClass" {
		t.Errorf("class.GetName() = %q, want %q", class.GetName(), "TestClass")
	}
}

func TestClassSetName(t *testing.T) {
	class := NewClass("TestClass", nil)
	class.SetName("ModifiedClass")

	if class.Name != "ModifiedClass" {
		t.Errorf("After SetName(\"ModifiedClass\"), class.Name = %q, want %q", class.Name, "ModifiedClass")
	}
}

func TestClassGetSuperClass(t *testing.T) {
	superClass := NewClass("SuperClass", nil)
	class := NewClass("TestClass", superClass)

	if class.GetSuperClass() == nil {
		t.Errorf("class.GetSuperClass() is nil, want non-nil")
	}
}

func TestClassSetSuperClass(t *testing.T) {
	class := NewClass("TestClass", nil)
	superClass := NewClass("SuperClass", nil)

	class.SetSuperClass(ClassToObject(superClass))

	if class.SuperClass == nil {
		t.Errorf("After SetSuperClass(), class.SuperClass is nil, want non-nil")
	}
}

func TestClassGetInstanceVarNames(t *testing.T) {
	class := NewClass("TestClass", nil)

	if len(class.GetInstanceVarNames()) != 0 {
		t.Errorf("len(class.GetInstanceVarNames()) = %d, want 0", len(class.GetInstanceVarNames()))
	}
}

func TestClassAddInstanceVarName(t *testing.T) {
	class := NewClass("TestClass", nil)

	class.AddInstanceVarName("var1")
	class.AddInstanceVarName("var2")

	varNames := class.GetInstanceVarNames()

	if len(varNames) != 2 {
		t.Errorf("len(class.GetInstanceVarNames()) = %d, want 2", len(varNames))
	}

	if varNames[0] != "var1" {
		t.Errorf("varNames[0] = %q, want %q", varNames[0], "var1")
	}

	if varNames[1] != "var2" {
		t.Errorf("varNames[1] = %q, want %q", varNames[1], "var2")
	}
}

func TestClassGetMethodDictionary(t *testing.T) {
	class := NewClass("TestClass", nil)

	methodDict := class.GetMethodDictionary()

	if methodDict == nil {
		t.Errorf("class.GetMethodDictionary() is nil, want non-nil")
	}

	if methodDict.GetEntryCount() != 0 {
		t.Errorf("methodDict.GetEntryCount() = %d, want 0", methodDict.GetEntryCount())
	}
}

func TestClassAddMethod(t *testing.T) {
	class := NewClass("TestClass", nil)

	// Create a selector and method
	selector := NewSymbol("testMethod")
	method := NewMethod(selector, class)

	// Add the method to the class
	class.AddMethod(selector, method)

	// Check that the method was added
	methodDict := class.GetMethodDictionary()

	if methodDict.GetEntryCount() != 1 {
		t.Errorf("methodDict.GetEntryCount() = %d, want 1", methodDict.GetEntryCount())
	}

	if methodDict.GetEntry("testMethod") != method {
		t.Errorf("methodDict.GetEntry(\"testMethod\") = %v, want %v", methodDict.GetEntry("testMethod"), method)
	}
}

func TestClassLookupMethod(t *testing.T) {
	// Create a class hierarchy
	superClass := NewClass("SuperClass", nil)
	class := NewClass("TestClass", superClass)

	// Add a method to the superclass
	superSelector := NewSymbol("superMethod")
	superMethod := NewMethod(superSelector, superClass)
	superClass.AddMethod(superSelector, superMethod)

	// Add a method to the class
	classSelector := NewSymbol("classMethod")
	classMethod := NewMethod(classSelector, class)
	class.AddMethod(classSelector, classMethod)

	// Look up methods
	foundSuperMethod := class.LookupMethod(superSelector)
	foundClassMethod := class.LookupMethod(classSelector)
	notFoundMethod := class.LookupMethod(NewSymbol("nonExistentMethod"))

	// Check results
	if foundSuperMethod != superMethod {
		t.Errorf("class.LookupMethod(superSelector) = %v, want %v", foundSuperMethod, superMethod)
	}

	if foundClassMethod != classMethod {
		t.Errorf("class.LookupMethod(classSelector) = %v, want %v", foundClassMethod, classMethod)
	}

	if notFoundMethod != nil {
		t.Errorf("class.LookupMethod(\"nonExistentMethod\") = %v, want nil", notFoundMethod)
	}
}

func TestClassNewInstance(t *testing.T) {
	// Create a class with instance variables
	class := NewClass("TestClass", nil)
	class.AddInstanceVarName("var1")
	class.AddInstanceVarName("var2")

	// Create an instance
	instance := class.NewInstance()

	if instance.Type() != core.OBJ_INSTANCE {
		t.Errorf("instance.Type() = %d, want %d", instance.Type(), core.OBJ_INSTANCE)
	}

	if instance.Class() == nil {
		t.Errorf("instance.Class() is nil, want non-nil")
	}

	// Check instance variables
	instVars := instance.InstanceVars()

	if len(instVars) != 2 {
		t.Errorf("len(instance.InstanceVars()) = %d, want 2", len(instVars))
	}

	// Check that instance variables are initialized to nil
	for i, v := range instVars {
		if !core.IsNilImmediate(v) {
			t.Errorf("instVars[%d] = %v, want nil", i, v)
		}
	}
}

func TestGetClassName(t *testing.T) {
	// Test with a class object
	class := NewClass("TestClass", nil)
	classObj := ClassToObject(class)

	if GetClassName(classObj) != "TestClass" {
		t.Errorf("GetClassName(classObj) = %q, want %q", GetClassName(classObj), "TestClass")
	}

	// Test with a non-class object
	// Create a string object instead of using an immediate value
	strObj := StringToObject(NewString("test"))
	if GetClassName(strObj) != "" {
		t.Errorf("GetClassName(strObj) = %q, want %q", GetClassName(strObj), "")
	}
}
