package main

import (
	"fmt"
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

	// Create methods using MethodBuilder
	sizeMethod := NewMethodBuilder(objectClass).
		Selector("size").
		Go()

	atMethod := NewMethodBuilder(sequenceableCollectionClass).
		Selector("at:").
		Go()

	atPutMethod := NewMethodBuilder(sequenceableCollectionClass).
		Selector("at:put:").
		Go()

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
	// nilClassInstance := &Object{Type: OBJ_INSTANCE, Class: nil}
	// method = vm.lookupMethod(nilClassInstance, sizeSelector)
	// if method != nil {
	// 	t.Errorf("Expected nil when receiver has nil class, got %v", method)
	// }
	// To do: test that this panics

}

func TestBadLookupMethod(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The function did not panic")
		}
	}()

	// Call the function that should panic
	BadLookupMethodHelper()
	fmt.Println("This line should not be reached")
}

func BadLookupMethodHelper() {
	vm := NewVM()

	badClass := NewClass("BadClass", nil)
	badClass.InstanceVars[METHOD_DICTIONARY_IV] = nil // Set method dictionary to nil
	badClassInstance := NewInstance(badClass)
	sizeSelector := NewSymbol("size")
	vm.lookupMethod(badClassInstance, sizeSelector) // should panic
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

	// Create methods for each class using MethodBuilder
	// We don't need to store the method objects since they're already in the method dictionaries
	NewMethodBuilder(objectClass).
		Selector("size").
		Go()

	collectionSizeMethod := NewMethodBuilder(collectionClass).
		Selector("size").
		Go()

	seqCollAtMethod := NewMethodBuilder(sequenceableCollectionClass).
		Selector("at:").
		Go()

	arrayAtPutMethod := NewMethodBuilder(arrayClass).
		Selector("at:put:").
		Go()

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

	// Add a method to the dictionary using MethodBuilder
	selectorName := "test"
	method := NewMethodBuilder(class).
		Selector(selectorName).
		Go()

	// Get the method dictionary again
	methodDict2 := class.GetMethodDict()

	// Check that it's the same dictionary
	if methodDict2 != methodDict {
		t.Errorf("Expected to get the same method dictionary")
	}

	// Check that the method is in the dictionary
	if methodDict2.Entries[selectorName] != method {
		t.Errorf("Expected to find method in dictionary")
	}

}
