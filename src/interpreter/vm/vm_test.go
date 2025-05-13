package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/vm"
)

func TestExecuteContextEmptyMethod(t *testing.T) {
	// Create a VM for testing
	virtualMachine := vm.NewVM()

	// Create a method with no bytecodes using MethodBuilder
	methodObj := compiler.NewMethodBuilder(virtualMachine.Classes.Get(vm.Object)).
		Selector("emptyMethod").
		Go()

	context := vm.NewContext(methodObj, virtualMachine.Classes.Get(vm.Object), []*core.Object{}, nil)

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
	builder := compiler.NewMethodBuilder(virtualMachine.Classes.Get(vm.Object)).Selector("pushMethod")
	literalIndex, builder := builder.AddLiteral(virtualMachine.NewInteger(42))
	methodObj := builder.PushLiteral(literalIndex).Go()

	context := vm.NewContext(methodObj, virtualMachine.Classes.Get(vm.Object), []*core.Object{}, nil)

	result, err := virtualMachine.ExecuteContext(context)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Method should return the value on the stack
	if core.IsIntegerImmediate(result) {
		intValue := core.GetIntegerImmediate(result)
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
	methodObj := compiler.NewMethodBuilder(virtualMachine.Classes.Get(vm.Object)).
		Selector("errorMethod").
		Go()

	// Set invalid bytecode manually
	method := classes.ObjectToMethod(methodObj)
	method.Bytecodes = []byte{255} // Invalid bytecode

	context := vm.NewContext(methodObj, virtualMachine.Classes.Get(vm.Object), []*core.Object{}, nil)

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
	if arrayObj.Type() != core.OBJ_ARRAY {
		t.Errorf("Expected array type to be OBJ_ARRAY, got %v", arrayObj.Type())
	}

	// Check that the array has the correct class
	if arrayObj.Class() == nil {
		t.Errorf("Expected array class to be set, got nil")
	}

	// Get the class of the array
	arrayClass := classes.ObjectToClass(arrayObj.Class())
	if arrayClass == nil {
		t.Errorf("Expected array class to be a valid class, got nil")
	}

	// Check that the class name is "Array"
	if arrayClass.GetName() != "Array" {
		t.Errorf("Expected array class name to be 'Array', got '%s'", arrayClass.GetName())
	}

	// Check that the class is the same as the Array class in the registry
	if arrayClass != virtualMachine.Classes.Get(vm.Array) {
		t.Errorf("Expected array class to be the Array class in the registry")
	}

	// Check that the array has the correct size
	array := classes.ObjectToArray(arrayObj)
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
