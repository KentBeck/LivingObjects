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
	methodObj := compiler.NewMethodBuilder(virtualMachine.ObjectClass).
		Selector("emptyMethod").
		Go()

	context := vm.NewContext(methodObj, virtualMachine.ObjectClass, []*core.Object{}, nil)

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
	builder := compiler.NewMethodBuilder(virtualMachine.ObjectClass).Selector("pushMethod")
	literalIndex, builder := builder.AddLiteral(virtualMachine.NewInteger(42))
	methodObj := builder.PushLiteral(literalIndex).Go()

	context := vm.NewContext(methodObj, virtualMachine.ObjectClass, []*core.Object{}, nil)

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
	methodObj := compiler.NewMethodBuilder(virtualMachine.ObjectClass).
		Selector("errorMethod").
		Go()

	// Set invalid bytecode manually
	method := classes.ObjectToMethod(methodObj)
	method.Bytecodes = []byte{255} // Invalid bytecode

	context := vm.NewContext(methodObj, virtualMachine.ObjectClass, []*core.Object{}, nil)

	_, err := virtualMachine.ExecuteContext(context)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}
