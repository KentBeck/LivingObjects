package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

func TestExecuteContextEmptyMethod(t *testing.T) {
	// Create a VM for testing
	virtualMachine := vm.NewVM()

	// Create a method with no bytecodes using MethodBuilder
	methodObj := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"])).
		Go("emptyMethod")

	context := vm.NewContext(methodObj, pile.ObjectToClass(virtualMachine.Globals["Object"]), []*pile.Object{}, nil)

	result, err := virtualMachine.ExecuteContext(context)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Empty method should return nil
	if result != virtualMachine.NilObject {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestExecuteContextWithStackValue(t *testing.T) {
	// Create a VM for testing
	virtualMachine := vm.NewVM()

	// Create a method that pushes a value onto the stack using MethodBuilder
	builder := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"]))
	literalIndex, builder := builder.AddLiteral(virtualMachine.NewInteger(42))
	methodObj := builder.PushLiteral(literalIndex).Go("pushMethod")

	context := vm.NewContext(methodObj, pile.ObjectToClass(virtualMachine.Globals["Object"]), []*pile.Object{}, nil)

	result, err := virtualMachine.ExecuteContext(context)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Method should return the value on the stack
	if pile.IsIntegerImmediate(result) {
		intValue := pile.GetIntegerImmediate(result)
		if intValue != 42 {
			t.Errorf("Expected 42, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}
}

func TestExecuteContextWithError(t *testing.T) {
	// Create a VM for testing
	virtualMachine := vm.NewVM()

	// Create a method with an invalid bytecode
	// Since we can't use the fluent API for invalid bytecodes, we'll create the method
	// and then manually set the bytecodes
	methodObj := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"])).
		Go("errorMethod")

	// Set invalid bytecode manually
	method := pile.ObjectToMethod(methodObj)
	method.Bytecodes = []byte{255} // Invalid bytecode

	context := vm.NewContext(methodObj, pile.ObjectToClass(virtualMachine.Globals["Object"]), []*pile.Object{}, nil)

	_, err := virtualMachine.ExecuteContext(context)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}

// TestNewArray tests that when we create an Array through the VM, the class is set correctly
func TestNewArray(t *testing.T) {
	// Create a new VM
	virtualMachine := vm.NewVM()

	// Create a new array
	arrayObj := virtualMachine.NewArray(3)

	// Check that the array is of the correct type
	if arrayObj.Type() != pile.OBJ_ARRAY {
		t.Errorf("Expected array type to be OBJ_ARRAY, got %v", arrayObj.Type())
	}

	// Check that the array has the correct class
	if arrayObj.Class() == nil {
		t.Errorf("Expected array class to be set, got nil")
	}

	// Get the class of the array
	arrayClass := pile.ObjectToClass(arrayObj.Class())
	if arrayClass == nil {
		t.Errorf("Expected array class to be a valid class, got nil")
	}

	// Check that the class name is "Array"
	if pile.GetClassName(arrayClass) != "Array" {
		t.Errorf("Expected array class name to be 'Array', got '%s'", pile.GetClassName(arrayClass))
	}

	// Check that the class is the same as the Array class in the registry
	if arrayClass != pile.ObjectToClass(virtualMachine.Globals["Array"]) {
		t.Errorf("Expected array class to be the Array class in the registry")
	}

	// Check that the array has the correct size
	array := pile.ObjectToArray(arrayObj)
	if array.Size() != 3 {
		t.Errorf("Expected array size to be 3, got %d", array.Size())
	}

	// Check that the array elements are initialized to nil
	for i := 0; i < array.Size(); i++ {
		elem := array.At(i)
		if elem != nil {
			t.Errorf("Expected array element %d to be nil, got %v", i, elem)
		}
	}
}
