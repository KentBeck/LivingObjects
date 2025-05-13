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

func TestGetClassString(t *testing.T) {
	class := NewClass("TestClass", nil)

	if GetClassString(class) != "Class TestClass" {
		t.Errorf("GetClassString(class) = %q, want %q", GetClassString(class), "Class TestClass")
	}
}

func TestGetClassName(t *testing.T) {
	class := NewClass("TestClass", nil)

	if GetClassName(class) != "TestClass" {
		t.Errorf("GetClassName(class) = %q, want %q", GetClassName(class), "TestClass")
	}
}

func TestSetClassName(t *testing.T) {
	class := NewClass("TestClass", nil)
	SetClassName(class, "ModifiedClass")

	if class.Name != "ModifiedClass" {
		t.Errorf("After SetClassName(class, \"ModifiedClass\"), class.Name = %q, want %q", class.Name, "ModifiedClass")
	}
}

func TestGetClassSuperClass(t *testing.T) {
	superClass := NewClass("SuperClass", nil)
	class := NewClass("TestClass", superClass)

	if GetClassSuperClass(class) == nil {
		t.Errorf("GetClassSuperClass(class) is nil, want non-nil")
	}
}

func TestSetClassSuperClass(t *testing.T) {
	class := NewClass("TestClass", nil)
	superClass := NewClass("SuperClass", nil)

	SetClassSuperClass(class, ClassToObject(superClass))

	if class.SuperClass == nil {
		t.Errorf("After SetClassSuperClass(), class.SuperClass is nil, want non-nil")
	}
}

func TestGetClassInstanceVarNames(t *testing.T) {
	class := NewClass("TestClass", nil)

	if len(GetClassInstanceVarNames(class)) != 0 {
		t.Errorf("len(GetClassInstanceVarNames(class)) = %d, want 0", len(GetClassInstanceVarNames(class)))
	}
}

func TestAddClassInstanceVarName(t *testing.T) {
	class := NewClass("TestClass", nil)

	AddClassInstanceVarName(class, "var1")
	AddClassInstanceVarName(class, "var2")

	varNames := GetClassInstanceVarNames(class)

	if len(varNames) != 2 {
		t.Errorf("len(GetClassInstanceVarNames(class)) = %d, want 2", len(varNames))
	}

	if varNames[0] != "var1" {
		t.Errorf("varNames[0] = %q, want %q", varNames[0], "var1")
	}

	if varNames[1] != "var2" {
		t.Errorf("varNames[1] = %q, want %q", varNames[1], "var2")
	}
}

func TestGetClassMethodDictionary(t *testing.T) {
	class := NewClass("TestClass", nil)

	methodDict := GetClassMethodDictionary(class)

	if methodDict == nil {
		t.Errorf("GetClassMethodDictionary(class) is nil, want non-nil")
	}

	if methodDict.GetEntryCount() != 0 {
		t.Errorf("methodDict.GetEntryCount() = %d, want 0", methodDict.GetEntryCount())
	}
}

func TestAddClassMethod(t *testing.T) {
	class := NewClass("TestClass", nil)

	// Create a selector and method
	selector := NewSymbol("testMethod")
	method := NewMethod(selector, class)

	// Add the method to the class
	AddClassMethod(class, selector, method)

	// Check that the method was added
	methodDict := GetClassMethodDictionary(class)

	if methodDict.GetEntryCount() != 1 {
		t.Errorf("methodDict.GetEntryCount() = %d, want 1", methodDict.GetEntryCount())
	}

	if methodDict.GetEntry("testMethod") != method {
		t.Errorf("methodDict.GetEntry(\"testMethod\") = %v, want %v", methodDict.GetEntry("testMethod"), method)
	}
}

func TestLookupClassMethod(t *testing.T) {
	// Create a class hierarchy
	superClass := NewClass("SuperClass", nil)
	class := NewClass("TestClass", superClass)

	// Add a method to the superclass
	superSelector := NewSymbol("superMethod")
	superMethod := NewMethod(superSelector, superClass)
	AddClassMethod(superClass, superSelector, superMethod)

	// Add a method to the class
	classSelector := NewSymbol("classMethod")
	classMethod := NewMethod(classSelector, class)
	AddClassMethod(class, classSelector, classMethod)

	// Look up methods
	foundSuperMethod := LookupClassMethod(class, superSelector)
	foundClassMethod := LookupClassMethod(class, classSelector)
	notFoundMethod := LookupClassMethod(class, NewSymbol("nonExistentMethod"))

	// Check results
	if foundSuperMethod != superMethod {
		t.Errorf("LookupClassMethod(class, superSelector) = %v, want %v", foundSuperMethod, superMethod)
	}

	if foundClassMethod != classMethod {
		t.Errorf("LookupClassMethod(class, classSelector) = %v, want %v", foundClassMethod, classMethod)
	}

	if notFoundMethod != nil {
		t.Errorf("LookupClassMethod(class, \"nonExistentMethod\") = %v, want nil", notFoundMethod)
	}
}

func TestNewClassInstance(t *testing.T) {
	// Create a class with instance variables
	class := NewClass("TestClass", nil)
	AddClassInstanceVarName(class, "var1")
	AddClassInstanceVarName(class, "var2")

	// Create an instance
	instance := NewClassInstance(class)

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

func TestGetClassNameFromObject(t *testing.T) {
	// Test with a class object
	class := NewClass("TestClass", nil)
	classObj := ClassToObject(class)

	if GetClassNameFromObject(classObj) != "TestClass" {
		t.Errorf("GetClassNameFromObject(classObj) = %q, want %q", GetClassNameFromObject(classObj), "TestClass")
	}

	// Test with a non-class object
	// Create a string object instead of using an immediate value
	strObj := StringToObject(NewString("test"))
	if GetClassNameFromObject(strObj) != "" {
		t.Errorf("GetClassNameFromObject(strObj) = %q, want %q", GetClassNameFromObject(strObj), "")
	}
}