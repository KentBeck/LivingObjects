package main

import (
	"testing"
)

// TestLookupMethod tests the lookupMethod function
func TestLookupMethod(t *testing.T) {
	vm := NewVM()

	// Create a class hierarchy
	objectClass := NewClass("Object", nil)
	collectionClass := NewClass("Collection", objectClass)
	sequenceableCollectionClass := NewClass("SequenceableCollection", collectionClass)
	arrayClass := NewClass("Array", sequenceableCollectionClass)

	// Create method selectors
	sizeSelector := NewSymbol("size")
	atSelector := NewSymbol("at:")
	atPutSelector := NewSymbol("at:put:")

	// Create methods
	sizeMethod := NewMethod(sizeSelector, objectClass)
	atMethod := NewMethod(atSelector, sequenceableCollectionClass)
	atPutMethod := NewMethod(atPutSelector, sequenceableCollectionClass)

	// Add methods to classes
	objectMethodDict := objectClass.GetMethodDict()
	objectMethodDict.Entries[sizeSelector.SymbolValue] = sizeMethod

	seqCollMethodDict := sequenceableCollectionClass.GetMethodDict()
	seqCollMethodDict.Entries[atSelector.SymbolValue] = atMethod
	seqCollMethodDict.Entries[atPutSelector.SymbolValue] = atPutMethod

	// Create an instance of Array
	arrayInstance := NewInstance(arrayClass)

	// Test cases

	// 1. Look up a method defined in a superclass (Object)
	method := vm.lookupMethod(arrayInstance, sizeSelector)
	if method != sizeMethod {
		t.Errorf("Expected to find size method from Object class, got %v", method)
	}

	// 2. Look up a method defined in an ancestor class (SequenceableCollection)
	method = vm.lookupMethod(arrayInstance, atSelector)
	if method != atMethod {
		t.Errorf("Expected to find at: method from SequenceableCollection class, got %v", method)
	}

	// 3. Look up a method defined in an ancestor class (SequenceableCollection)
	method = vm.lookupMethod(arrayInstance, atPutSelector)
	if method != atPutMethod {
		t.Errorf("Expected to find at:put: method from SequenceableCollection class, got %v", method)
	}

	// 4. Look up a method that doesn't exist
	notFoundSelector := NewSymbol("notFound")
	method = vm.lookupMethod(arrayInstance, notFoundSelector)
	if method != nil {
		t.Errorf("Expected nil for non-existent method, got %v", method)
	}

	// 5. Look up a method on a class object directly
	method = vm.lookupMethod(arrayClass, sizeSelector)
	if method != sizeMethod {
		t.Errorf("Expected to find size method from Object class when looking up on class, got %v", method)
	}

	// 6. Test with nil class
	nilClassInstance := &Object{Type: OBJ_INSTANCE, Class: nil}
	method = vm.lookupMethod(nilClassInstance, sizeSelector)
	if method != nil {
		t.Errorf("Expected nil when receiver has nil class, got %v", method)
	}

	// 7. Test with nil method dictionary
	badClass := NewClass("BadClass", nil)
	badClass.InstanceVars[METHOD_DICTIONARY_IV] = nil // Set method dictionary to nil
	badClassInstance := NewInstance(badClass)
	method = vm.lookupMethod(badClassInstance, sizeSelector)
	if method != nil {
		t.Errorf("Expected nil when class has nil method dictionary, got %v", method)
	}
}

// TestLookupMethodWithInheritance tests method lookup with multiple levels of inheritance
func TestLookupMethodWithInheritance(t *testing.T) {
	vm := NewVM()

	// Create a deeper class hierarchy
	objectClass := NewClass("Object", nil)
	collectionClass := NewClass("Collection", objectClass)
	sequenceableCollectionClass := NewClass("SequenceableCollection", collectionClass)
	arrayClass := NewClass("Array", sequenceableCollectionClass)

	// Create method selectors
	sizeSelector := NewSymbol("size")
	atSelector := NewSymbol("at:")
	atPutSelector := NewSymbol("at:put:")

	// Create methods for each class
	objectSizeMethod := NewMethod(sizeSelector, objectClass)
	collectionSizeMethod := NewMethod(sizeSelector, collectionClass)
	seqCollAtMethod := NewMethod(atSelector, sequenceableCollectionClass)
	arrayAtPutMethod := NewMethod(atPutSelector, arrayClass)

	// Add methods to classes
	objectMethodDict := objectClass.GetMethodDict()
	objectMethodDict.Entries[sizeSelector.SymbolValue] = objectSizeMethod

	collectionMethodDict := collectionClass.GetMethodDict()
	collectionMethodDict.Entries[sizeSelector.SymbolValue] = collectionSizeMethod

	seqCollMethodDict := sequenceableCollectionClass.GetMethodDict()
	seqCollMethodDict.Entries[atSelector.SymbolValue] = seqCollAtMethod

	arrayMethodDict := arrayClass.GetMethodDict()
	arrayMethodDict.Entries[atPutSelector.SymbolValue] = arrayAtPutMethod

	// Create an instance of Array
	arrayInstance := NewInstance(arrayClass)

	// Test cases

	// 1. Method should be found in the receiver's class first
	method := vm.lookupMethod(arrayInstance, atPutSelector)
	if method != arrayAtPutMethod {
		t.Errorf("Expected to find at:put: method in Array class, got %v", method)
	}

	// 2. Method should be found in parent class if not in receiver's class
	method = vm.lookupMethod(arrayInstance, atSelector)
	if method != seqCollAtMethod {
		t.Errorf("Expected to find at: method in SequenceableCollection class, got %v", method)
	}

	// 3. Method should be found in closest ancestor that defines it
	method = vm.lookupMethod(arrayInstance, sizeSelector)
	if method != collectionSizeMethod {
		t.Errorf("Expected to find size method in Collection class, got %v", method)
	}

	// 4. Method lookup should work with class objects too
	method = vm.lookupMethod(arrayClass, sizeSelector)
	if method != collectionSizeMethod {
		t.Errorf("Expected to find size method in Collection class when looking up on class, got %v", method)
	}
}

// TestGetMethodDict tests the GetMethodDict method
func TestGetMethodDict(t *testing.T) {
	// Create a class
	class := NewClass("TestClass", nil)

	// Get the method dictionary
	methodDict := class.GetMethodDict()

	// Check that it's a dictionary
	if methodDict.Type != OBJ_DICTIONARY {
		t.Errorf("Expected method dictionary to be a dictionary, got %v", methodDict.Type)
	}

	// Check that it's empty
	if len(methodDict.Entries) != 0 {
		t.Errorf("Expected method dictionary to be empty, got %d entries", len(methodDict.Entries))
	}

	// Add a method to the dictionary
	selector := NewSymbol("test")
	method := NewMethod(selector, class)
	methodDict.Entries[selector.SymbolValue] = method

	// Get the method dictionary again
	methodDict2 := class.GetMethodDict()

	// Check that it's the same dictionary
	if methodDict2 != methodDict {
		t.Errorf("Expected to get the same method dictionary")
	}

	// Check that the method is in the dictionary
	if methodDict2.Entries[selector.SymbolValue] != method {
		t.Errorf("Expected to find method in dictionary")
	}

	// Test with non-class object
	instance := NewInstance(class)
	methodDict = instance.GetMethodDict()
	if methodDict.Type != OBJ_NIL {
		t.Errorf("Expected nil for non-class object, got %v", methodDict.Type)
	}

	// Test with nil instance variables
	badClass := &Object{Type: OBJ_CLASS, InstanceVars: nil}
	methodDict = badClass.GetMethodDict()
	if methodDict.Type != OBJ_NIL {
		t.Errorf("Expected nil for class with nil instance variables, got %v", methodDict.Type)
	}
}
