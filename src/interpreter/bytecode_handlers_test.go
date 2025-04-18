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

func TestExecutePushInstanceVariable(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a class with an instance variable
	class := NewClass("TestClass", nil)
	class.InstanceVarNames = append(class.InstanceVarNames, "testVar")

	// Create an instance with a value for the instance variable
	instance := NewInstance(class)
	instance.SetInstanceVarByIndex(0, NewInteger(42))

	// Create a method
	methodObj := NewMethod(NewSymbol("test"), class)
	methodObj.Method.Bytecodes = []byte{
		PUSH_INSTANCE_VARIABLE, 0, 0, 0, 0, // Push instance variable at index 0
	}

	// Create a context
	context := NewContext(methodObj, instance, []*Object{}, nil)

	// Execute the bytecode
	err := vm.ExecutePushInstanceVariable(context)
	if err != nil {
		t.Errorf("ExecutePushInstanceVariable returned an error: %v", err)
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

func TestExecutePushTemporaryVariable(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method with a temporary variable
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)
	methodObj.Method.TempVarNames = append(methodObj.Method.TempVarNames, "temp")
	methodObj.Method.Bytecodes = []byte{
		PUSH_TEMPORARY_VARIABLE, 0, 0, 0, 0, // Push temporary variable at index 0
	}

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Set the temporary variable
	context.SetTempVar("temp", NewInteger(42))

	// Execute the bytecode
	err := vm.ExecutePushTemporaryVariable(context)
	if err != nil {
		t.Errorf("ExecutePushTemporaryVariable returned an error: %v", err)
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

func TestExecuteStoreInstanceVariable(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a class with an instance variable
	class := NewClass("TestClass", nil)
	class.InstanceVarNames = append(class.InstanceVarNames, "testVar")

	// Create an instance
	instance := NewInstance(class)

	// Create a method
	methodObj := NewMethod(NewSymbol("test"), class)
	methodObj.Method.Bytecodes = []byte{
		STORE_INSTANCE_VARIABLE, 0, 0, 0, 0, // Store into instance variable at index 0
	}

	// Create a context
	context := NewContext(methodObj, instance, []*Object{}, nil)

	// Push a value onto the stack
	context.Push(NewInteger(42))

	// Execute the bytecode
	err := vm.ExecuteStoreInstanceVariable(context)
	if err != nil {
		t.Errorf("ExecuteStoreInstanceVariable returned an error: %v", err)
	}

	// Check the instance variable
	value := instance.GetInstanceVarByIndex(0)
	if value.Type != OBJ_INTEGER || value.IntegerValue != 42 {
		t.Errorf("Expected instance variable to be 42, got %v", value)
	}

	// Check the stack (the value should be pushed back onto the stack)
	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	// Check the value on the stack
	stackValue := context.Pop()
	if stackValue.Type != OBJ_INTEGER || stackValue.IntegerValue != 42 {
		t.Errorf("Expected 42 on the stack, got %v", stackValue)
	}
}

func TestExecuteStoreTemporaryVariable(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method with a temporary variable
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)
	methodObj.Method.TempVarNames = append(methodObj.Method.TempVarNames, "temp")
	methodObj.Method.Bytecodes = []byte{
		STORE_TEMPORARY_VARIABLE, 0, 0, 0, 0, // Store into temporary variable at index 0
	}

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Push a value onto the stack
	context.Push(NewInteger(42))

	// Execute the bytecode
	err := vm.ExecuteStoreTemporaryVariable(context)
	if err != nil {
		t.Errorf("ExecuteStoreTemporaryVariable returned an error: %v", err)
	}

	// Check the temporary variable
	value := context.GetTempVar("temp")
	if value.Type != OBJ_INTEGER || value.IntegerValue != 42 {
		t.Errorf("Expected temporary variable to be 42, got %v", value)
	}

	// Check the stack (the value should be pushed back onto the stack)
	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	// Check the value on the stack
	stackValue := context.Pop()
	if stackValue.Type != OBJ_INTEGER || stackValue.IntegerValue != 42 {
		t.Errorf("Expected 42 on the stack, got %v", stackValue)
	}
}

func TestExecuteReturnStackTop(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Push a value onto the stack
	context.Push(NewInteger(42))

	// Execute the bytecode
	result, err := vm.ExecuteReturnStackTop(context)
	if err != nil {
		t.Errorf("ExecuteReturnStackTop returned an error: %v", err)
	}

	// Check the result
	if result.Type != OBJ_INTEGER || result.IntegerValue != 42 {
		t.Errorf("Expected result to be 42, got %v", result)
	}

	// Check the stack (the value should be popped)
	if context.StackPointer != 0 {
		t.Errorf("Expected stack pointer to be 0, got %d", context.StackPointer)
	}
}

func TestExecuteJump(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)
	methodObj.Method.Bytecodes = []byte{
		JUMP, 0, 0, 0, 10, // Jump to offset 10
	}

	// Create a context
	context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	// Execute the bytecode
	skipIncrement, err := vm.ExecuteJump(context)
	if err != nil {
		t.Errorf("ExecuteJump returned an error: %v", err)
	}

	// Check the PC
	if context.PC != 10 {
		t.Errorf("Expected PC to be 10, got %d", context.PC)
	}

	// Check the skipIncrement flag
	if !skipIncrement {
		t.Errorf("Expected skipIncrement to be true")
	}
}

func TestExecuteJumpIfTrue(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)
	methodObj.Method.Bytecodes = []byte{
		JUMP_IF_TRUE, 0, 0, 0, 10, // Jump to offset 10 if true
	}

	// Test with true condition
	{
		// Create a context
		context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

		// Push a true value onto the stack
		context.Push(NewBoolean(true))

		// Execute the bytecode
		skipIncrement, err := vm.ExecuteJumpIfTrue(context)
		if err != nil {
			t.Errorf("ExecuteJumpIfTrue returned an error: %v", err)
		}

		// Check the PC
		if context.PC != 10 {
			t.Errorf("Expected PC to be 10, got %d", context.PC)
		}

		// Check the skipIncrement flag
		if !skipIncrement {
			t.Errorf("Expected skipIncrement to be true")
		}
	}

	// Test with false condition
	{
		// Create a context
		context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

		// Push a false value onto the stack
		context.Push(NewBoolean(false))

		// Execute the bytecode
		skipIncrement, err := vm.ExecuteJumpIfTrue(context)
		if err != nil {
			t.Errorf("ExecuteJumpIfTrue returned an error: %v", err)
		}

		// Check the PC (should not change)
		if context.PC != 0 {
			t.Errorf("Expected PC to be 0, got %d", context.PC)
		}

		// Check the skipIncrement flag
		if skipIncrement {
			t.Errorf("Expected skipIncrement to be false")
		}
	}
}

func TestExecuteJumpIfFalse(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method
	methodObj := NewMethod(NewSymbol("test"), vm.ObjectClass)
	methodObj.Method.Bytecodes = []byte{
		JUMP_IF_FALSE, 0, 0, 0, 10, // Jump to offset 10 if false
	}

	// Test with false condition
	{
		// Create a context
		context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

		// Push a false value onto the stack
		context.Push(NewBoolean(false))

		// Execute the bytecode
		skipIncrement, err := vm.ExecuteJumpIfFalse(context)
		if err != nil {
			t.Errorf("ExecuteJumpIfFalse returned an error: %v", err)
		}

		// Check the PC
		if context.PC != 10 {
			t.Errorf("Expected PC to be 10, got %d", context.PC)
		}

		// Check the skipIncrement flag
		if !skipIncrement {
			t.Errorf("Expected skipIncrement to be true")
		}
	}

	// Test with true condition
	{
		// Create a context
		context := NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

		// Push a true value onto the stack
		context.Push(NewBoolean(true))

		// Execute the bytecode
		skipIncrement, err := vm.ExecuteJumpIfFalse(context)
		if err != nil {
			t.Errorf("ExecuteJumpIfFalse returned an error: %v", err)
		}

		// Check the PC (should not change)
		if context.PC != 0 {
			t.Errorf("Expected PC to be 0, got %d", context.PC)
		}

		// Check the skipIncrement flag
		if skipIncrement {
			t.Errorf("Expected skipIncrement to be false")
		}
	}
}
