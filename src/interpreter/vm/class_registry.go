package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// ClassType enum for standard class types
type ClassType string

const (
	Object          ClassType = "Object"
	UndefinedObject ClassType = "UndefinedObject" // nil class
	True            ClassType = "True"
	False           ClassType = "False"
	Integer         ClassType = "Integer"
	Float           ClassType = "Float"
	String          ClassType = "String"
	Symbol          ClassType = "Symbol"
	Array           ClassType = "Array"
	Dictionary      ClassType = "Dictionary"
	ByteArray       ClassType = "ByteArray"
	Block           ClassType = "Block"
	Method          ClassType = "Method"
	Class           ClassType = "Class"
	ContextClass    ClassType = "Context"
	Exception       ClassType = "Exception"
)

// ClassRegistry holds references to all standard Smalltalk classes
// This centralizes class management and provides a cleaner way to
// access classes than having many separate fields in the VM struct
type ClassRegistry struct {
	// Map of registered classes by ClassType
	classesByType map[ClassType]*classes.Class
	
	// Map of registered classes by name (for non-standard classes)
	classesByName map[string]*classes.Class
}

// NewClassRegistry creates a new class registry
func NewClassRegistry() *ClassRegistry {
	return &ClassRegistry{
		classesByType: make(map[ClassType]*classes.Class),
		classesByName: make(map[string]*classes.Class),
	}
}

// Register adds a class to the registry
func (r *ClassRegistry) Register(classType ClassType, class *classes.Class) {
	if class == nil {
		return
	}
	
	r.classesByType[classType] = class
	r.classesByName[class.GetName()] = class
}

// RegisterNamed adds a class to the registry by name only (for non-standard classes)
func (r *ClassRegistry) RegisterNamed(name string, class *classes.Class) {
	if class == nil {
		return
	}
	
	r.classesByName[name] = class
}

// Get retrieves a class by type
func (r *ClassRegistry) Get(classType ClassType) *classes.Class {
	return r.classesByType[classType]
}

// GetByName retrieves a class by name
func (r *ClassRegistry) GetByName(name string) *classes.Class {
	return r.classesByName[name]
}

// All returns all registered classes
func (r *ClassRegistry) All() []*classes.Class {
	// Combine both maps (avoiding duplicates)
	uniqueClasses := make(map[*classes.Class]bool)
	
	for _, class := range r.classesByType {
		uniqueClasses[class] = true
	}
	
	for _, class := range r.classesByName {
		uniqueClasses[class] = true
	}
	
	// Convert to slice
	result := make([]*classes.Class, 0, len(uniqueClasses))
	for class := range uniqueClasses {
		result = append(result, class)
	}
	
	return result
}

// AllAsObjects returns all registered classes as core.Objects
func (r *ClassRegistry) AllAsObjects() []*core.Object {
	// Fix ambiguous reference by using a different variable name
	allClasses := r.All()
	classObjs := make([]*core.Object, len(allClasses))
	
	for i, class := range allClasses {
		classObjs[i] = classes.ClassToObject(class)
	}
	
	return classObjs
}