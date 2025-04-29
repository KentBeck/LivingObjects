package main

import (
	"unsafe"
)

// MethodToObject converts a Method to an Object
func MethodToObject(m *Method) *Object {
	return (*Object)(unsafe.Pointer(m))
}

// ObjectToMethod converts an Object to a Method
func ObjectToMethod(o *Object) *Method {
	return (*Method)(unsafe.Pointer(o))
}

// ClassToObject converts a Class to an Object
func ClassToObject(c *Class) *Object {
	return (*Object)(unsafe.Pointer(c))
}

// ObjectToClass converts an Object to a Class
func ObjectToClass(o *Object) *Class {
	return (*Class)(unsafe.Pointer(o))
}

// GetStringValue gets the string value of a string
func GetStringValue(obj *Object) string {
	if obj.Type() != OBJ_STRING {
		return ""
	}
	return ObjectToString(obj).Value
}

// GetClassName gets the name of a class
func GetClassName(obj *Object) string {
	if obj.Type() != OBJ_CLASS {
		return ""
	}
	return ObjectToClass(obj).Name
}
