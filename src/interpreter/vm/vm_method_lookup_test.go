package vm_test

import (
	"smalltalklsp/interpreter/pile"
	"fmt"
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/vm"
)

// TestLookupMethod tests the lookupMethod function
func TestLookupMethod(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Create a class hierarchy
	objectClass := pile.NewClass("Object", nil)
	collectionClass := pile.NewClass("Collection", objectClass)
	sequenceableCollectionClass := pile.NewClass("SequenceableCollection", collectionClass)
	arrayClass := pile.NewClass("Array", sequenceableCollectionClass)

	// Create method selectors
	sizeSelector := pile.NewSymbol("size")
	atSelector := pile.NewSymbol("at:")
	atPutSelector := pile.NewSymbol("at:put:")

	// Create methods using MethodBuilder
	sizeMethod := compiler.NewMethodBuilder(objectClass).
		Selector("size").
		Go()

	atMethod := compiler.NewMethodBuilder(sequenceableCollectionClass).
		Selector("at:").
		Go()

	atPutMethod := compiler.NewMethodBuilder(sequenceableCollectionClass).
		Selector("at:put:").
		Go()

	// Create an instance of Array
	arrayInstance := pile.NewInstance((*pile.Class)(unsafe.Pointer(arrayClass)))

	// Test cases

	// 1. Look up a method defined in a superclass (Object)
	method := virtualMachine.LookupMethod(arrayInstance, sizeSelector)
	if method != sizeMethod {
		t.Errorf("Expected to find size method from Object class, got %v", method)
	}

	// 2. Look up a method defined in an ancestor class (SequenceableCollection)
	method = virtualMachine.LookupMethod(arrayInstance, atSelector)
	if method != atMethod {
		t.Errorf("Expected to find at: method from SequenceableCollection class, got %v", method)
	}

	// 3. Look up a method defined in an ancestor class (SequenceableCollection)
	method = virtualMachine.LookupMethod(arrayInstance, atPutSelector)
	if method != atPutMethod {
		t.Errorf("Expected to find at:put: method from SequenceableCollection class, got %v", method)
	}

	// 4. Look up a method that doesn't exist
	notFoundSelector := pile.NewSymbol("notFound")
	method = virtualMachine.LookupMethod(arrayInstance, notFoundSelector)
	if method != nil {
		t.Errorf("Expected nil for non-existent method, got %v", method)
	}

	// 5. Look up a method on a class object directly
	method = virtualMachine.LookupMethod(pile.ClassToObject(arrayClass), sizeSelector)
	if method != sizeMethod {
		t.Errorf("Expected to find size method from Object class when looking up on class, got %v", method)
	}

	// 6. Test with nil class
	// nilClassInstance := &core.Object{Type: OBJ_INSTANCE, Class: nil}
	// method = virtualMachine.LookupMethod(nilClassInstance, sizeSelector)
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
	virtualMachine := vm.NewVM()

	badClass := pile.NewClass("BadClass", nil)
	badClass.MethodDictionary = nil // Set method dictionary to nil
	badClassInstance := pile.NewInstance((*pile.Class)(unsafe.Pointer(badClass)))
	sizeSelector := pile.NewSymbol("size")

	// Add a panic handler to make sure we actually panic
	defer func() {
		if r := recover(); r == nil {
			panic("Expected panic but none occurred")
		}
	}()

	virtualMachine.LookupMethod(badClassInstance, sizeSelector) // should panic
}

// TestLookupMethodWithInheritance tests method lookup with multiple levels of inheritance
func TestLookupMethodWithInheritance(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Create a deeper class hierarchy
	objectClass := pile.NewClass("Object", nil)
	collectionClass := pile.NewClass("Collection", objectClass)
	sequenceableCollectionClass := pile.NewClass("SequenceableCollection", collectionClass)
	arrayClass := pile.NewClass("Array", sequenceableCollectionClass)

	// Create method selectors
	sizeSelector := pile.NewSymbol("size")
	atSelector := pile.NewSymbol("at:")
	atPutSelector := pile.NewSymbol("at:put:")

	// Create methods for each class using MethodBuilder
	// We don't need to store the method objects since they're already in the method dictionaries
	compiler.NewMethodBuilder(objectClass).
		Selector("size").
		Go()

	collectionSizeMethod := compiler.NewMethodBuilder(collectionClass).
		Selector("size").
		Go()

	seqCollAtMethod := compiler.NewMethodBuilder(sequenceableCollectionClass).
		Selector("at:").
		Go()

	arrayAtPutMethod := compiler.NewMethodBuilder(arrayClass).
		Selector("at:put:").
		Go()

	// Create an instance of Array
	arrayInstance := pile.NewInstance((*pile.Class)(unsafe.Pointer(arrayClass)))

	// Test cases

	// 1. Method should be found in the receiver's class first
	method := virtualMachine.LookupMethod(arrayInstance, atPutSelector)
	if method != arrayAtPutMethod {
		t.Errorf("Expected to find at:put: method in Array class, got %v", method)
	}

	// 2. Method should be found in parent class if not in receiver's class
	method = virtualMachine.LookupMethod(arrayInstance, atSelector)
	if method != seqCollAtMethod {
		t.Errorf("Expected to find at: method in SequenceableCollection class, got %v", method)
	}

	// 3. Method should be found in closest ancestor that defines it
	method = virtualMachine.LookupMethod(arrayInstance, sizeSelector)
	if method != collectionSizeMethod {
		t.Errorf("Expected to find size method in Collection class, got %v", method)
	}

	// 4. Method lookup should work with class objects too
	method = virtualMachine.LookupMethod(pile.ClassToObject(arrayClass), sizeSelector)
	if method != collectionSizeMethod {
		t.Errorf("Expected to find size method in Collection class when looking up on class, got %v", method)
	}
}

// TestGetMethodDict tests the GetMethodDict method
func TestGetMethodDict(t *testing.T) {
	// Create a class
	class := pile.NewClass("TestClass", nil)

	// Get the method dictionary
	methodDict := class.GetMethodDict()

	// Check that it's a dictionary
	if methodDict.Type() != pile.OBJ_DICTIONARY {
		t.Errorf("Expected method dictionary to be a dictionary, got %v", methodDict.Type())
	}

	// Check that it's empty
	// Convert to Dictionary to access entries
	dict := pile.ObjectToDictionary(methodDict)
	if len(dict.Entries) != 0 {
		t.Errorf("Expected method dictionary to be empty, got %d entries", len(dict.Entries))
	}

	// Add a method to the dictionary using MethodBuilder
	selectorName := "test"
	method := compiler.NewMethodBuilder(class).
		Selector(selectorName).
		Go()

	// Get the method dictionary again
	methodDict2 := class.GetMethodDict()

	// Check that it's the same dictionary
	if methodDict2 != methodDict {
		t.Errorf("Expected to get the same method dictionary")
	}

	// Check that the method is in the dictionary
	// Convert to Dictionary to access entries
	dict2 := pile.ObjectToDictionary(methodDict2)
	if dict2.Entries[selectorName] != method {
		t.Errorf("Expected to find method in dictionary")
	}
}