package main

import (
	"testing"
)

func TestExecutePushLiteral(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method with literals
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)
	methodObj.Method.Literals = append(methodObj.Method.Literals, NewInteger(42))
	methodObj.Method.Bytecodes = []byte{
		PUSH_LITERAL,
		0, 0, 0, 0, // Index 0
	}

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Execute the bytecode
	err := vm.ExecutePushLiteral(context)
	if err != nil {
		t.Errorf("ExecutePushLiteral returned an error: %v", err)
	}

	// Check the stack
	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	// Check the value on the stack
	value := context.Pop()
	if value.Type != OBJ_INTEGER || value.IntegerValue != 42 {
		t.Errorf("Expected 42 on the stack, got %v", value)
	}
}

func TestExecutePushSelf(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)

	// Create a receiver
	receiver := NewInstance(vm.ObjectClass)

	// Create a context
	context := NewContext(methodObj, receiver, []*Object{}, nil)

	// Execute the bytecode
	err := vm.ExecutePushSelf(context)
	if err != nil {
		t.Errorf("ExecutePushSelf returned an error: %v", err)
	}

	// Check the stack
	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	// Check the value on the stack
	value := context.Pop()
	if value != receiver {
		t.Errorf("Expected receiver on the stack, got %v", value)
	}
}

func TestExecutePop(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Push a value onto the stack
	context.Push(NewInteger(42))

	// Execute the bytecode
	err := vm.ExecutePop(context)
	if err != nil {
		t.Errorf("ExecutePop returned an error: %v", err)
	}

	// Check the stack
	if context.StackPointer != 0 {
		t.Errorf("Expected stack pointer to be 0, got %d", context.StackPointer)
	}
}

func TestExecuteDuplicate(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Push a value onto the stack
	context.Push(NewInteger(42))

	// Execute the bytecode
	err := vm.ExecuteDuplicate(context)
	if err != nil {
		t.Errorf("ExecuteDuplicate returned an error: %v", err)
	}

	// Check the stack
	if context.StackPointer != 2 {
		t.Errorf("Expected stack pointer to be 2, got %d", context.StackPointer)
	}

	// Check the values on the stack
	value1 := context.Pop()
	value2 := context.Pop()
	if value1.Type != OBJ_INTEGER || value1.IntegerValue != 42 {
		t.Errorf("Expected 42 on the stack, got %v", value1)
	}
	if value2.Type != OBJ_INTEGER || value2.IntegerValue != 42 {
		t.Errorf("Expected 42 on the stack, got %v", value2)
	}
}

func TestExecuteSendMessage(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method with literals
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)
	methodObj.Method.Literals = append(methodObj.Method.Literals, NewInteger(2))  // Literal 0
	methodObj.Method.Literals = append(methodObj.Method.Literals, NewInteger(3))  // Literal 1
	methodObj.Method.Literals = append(methodObj.Method.Literals, NewSymbol("+")) // Literal 2
	methodObj.Method.Bytecodes = []byte{
		PUSH_LITERAL, 0, 0, 0, 0, // Push 2
		PUSH_LITERAL, 0, 0, 0, 1, // Push 3
		SEND_MESSAGE, 0, 0, 0, 2, 0, 0, 0, 1, // Send + with 1 arg
	}

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Set up the context's PC to point to the SEND_MESSAGE bytecode
	context.PC = 10 // After the two PUSH_LITERAL instructions

	// Push the receiver and argument onto the stack
	context.Push(NewInteger(2)) // Receiver
	context.Push(NewInteger(3)) // Argument

	// Execute the bytecode
	result, err := vm.ExecuteSendMessage(context)
	if err != nil {
		t.Errorf("ExecuteSendMessage returned an error: %v", err)
	}

	// Check the result
	if result == nil {
		t.Errorf("Expected a result, got nil")
	} else if result.Type != OBJ_INTEGER || result.IntegerValue != 5 {
		t.Errorf("Expected result to be 5, got %v", result)
	}
}
