package vm_test

import (
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/bytecode"
	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

func TestExecutePushLiteral(t *testing.T) {
	virtualMachine := vm.NewVM()

	builder := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"]))
	literalIndex, builder := builder.AddLiteral(virtualMachine.NewInteger(42))
	methodObj := builder.PushLiteral(literalIndex).Go("test")

	context := vm.NewContext(methodObj, pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])), []*pile.Object{}, nil)

	err := virtualMachine.ExecutePushLiteral(context)
	if err != nil {
		t.Errorf("ExecutePushLiteral returned an error: %v", err)
	}

	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	value := context.Pop()
	// Check for immediate integer
	if pile.IsIntegerImmediate(value) {
		intValue := pile.GetIntegerImmediate(value)
		if intValue != 42 {
			t.Errorf("Expected 42 on the stack, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", value)
	}
}

func TestExecutePushSelf(t *testing.T) {
	virtualMachine := vm.NewVM()

	methodObj := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"])).
		Go("test")

	// Convert to pile.Class
	pileClass := (*pile.Class)(unsafe.Pointer(pile.ObjectToClass(virtualMachine.Globals["Object"])))
	receiver := pile.NewInstance(pileClass)

	context := vm.NewContext(methodObj, receiver, []*pile.Object{}, nil)

	err := virtualMachine.ExecutePushSelf(context)
	if err != nil {
		t.Errorf("ExecutePushSelf returned an error: %v", err)
	}

	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	value := context.Pop()
	if value != receiver {
		t.Errorf("Expected receiver on the stack, got %v", value)
	}
}

func TestExecutePop(t *testing.T) {
	virtualMachine := vm.NewVM()

	methodObj := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"])).
		Go("test")

	context := vm.NewContext(methodObj, pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])), []*pile.Object{}, nil)

	context.Push(virtualMachine.NewInteger(42))

	err := virtualMachine.ExecutePop(context)
	if err != nil {
		t.Errorf("ExecutePop returned an error: %v", err)
	}

	if context.StackPointer != 0 {
		t.Errorf("Expected stack pointer to be 0, got %d", context.StackPointer)
	}
}

func TestExecuteDuplicate(t *testing.T) {
	virtualMachine := vm.NewVM()

	methodObj := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"])).
		Go("test")

	context := vm.NewContext(methodObj, pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])), []*pile.Object{}, nil)

	context.Push(virtualMachine.NewInteger(42))

	err := virtualMachine.ExecuteDuplicate(context)
	if err != nil {
		t.Errorf("ExecuteDuplicate returned an error: %v", err)
	}

	if context.StackPointer != 2 {
		t.Errorf("Expected stack pointer to be 2, got %d", context.StackPointer)
	}

	value1 := context.Pop()
	value2 := context.Pop()

	// Check for immediate integer
	if pile.IsIntegerImmediate(value1) {
		intValue := pile.GetIntegerImmediate(value1)
		if intValue != 42 {
			t.Errorf("Expected 42 on the stack, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", value1)
	}

	// Check for immediate integer
	if pile.IsIntegerImmediate(value2) {
		intValue := pile.GetIntegerImmediate(value2)
		if intValue != 42 {
			t.Errorf("Expected 42 on the stack, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", value2)
	}
}

func TestExecuteSendMessage(t *testing.T) {
	virtualMachine := vm.NewVM()

	// The Integer addition primitive is already defined by the VM

	builder := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"]))
	twoIndex, builder := builder.AddLiteral(virtualMachine.NewInteger(2))
	threeIndex, builder := builder.AddLiteral(virtualMachine.NewInteger(3))
	plusIndex, builder := builder.AddLiteral(pile.NewSymbol("+"))

	methodObj := builder.
		PushLiteral(twoIndex).
		PushLiteral(threeIndex).
		SendMessage(plusIndex, 1).
		Go("test")

	context := vm.NewContext(methodObj, pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])), []*pile.Object{}, nil)

	context.PC = 10 // After the two PUSH_LITERAL instructions

	context.Push(virtualMachine.NewInteger(2)) // Receiver
	context.Push(virtualMachine.NewInteger(3)) // Argument

	result, err := virtualMachine.ExecuteSendMessage(context)
	if err != nil {
		t.Errorf("ExecuteSendMessage returned an error: %v", err)
	}

	if result == nil {
		t.Errorf("Expected a result, got nil")
	} else if pile.IsIntegerImmediate(result) {
		intValue := pile.GetIntegerImmediate(result)
		if intValue != 5 {
			t.Errorf("Expected result to be 5, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}
}

func TestExecutePushInstanceVariable(t *testing.T) {
	virtualMachine := vm.NewVM()

	class := pile.NewClass("TestClass", nil)
	pile.AddClassInstanceVarName(class, "testVar")

	instance := pile.NewInstance(class)
	instance.SetInstanceVarByIndex(0, virtualMachine.NewInteger(42))

	methodObj := compiler.NewMethodBuilder(class).
		PushInstanceVariable(0).
		Go("test")

	context := vm.NewContext(methodObj, instance, []*pile.Object{}, nil)

	err := virtualMachine.ExecutePushInstanceVariable(context)
	if err != nil {
		t.Errorf("ExecutePushInstanceVariable returned an error: %v", err)
	}

	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	value := context.Pop()
	// Check for immediate integer
	if pile.IsIntegerImmediate(value) {
		intValue := pile.GetIntegerImmediate(value)
		if intValue != 42 {
			t.Errorf("Expected 42 on the stack, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", value)
	}
}

func TestExecutePushTemporaryVariable(t *testing.T) {
	virtualMachine := vm.NewVM()

	methodObj := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"])).
		TempVars([]string{"temp"}).
		PushTemporaryVariable(0).
		Go("test")

	context := vm.NewContext(methodObj, pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])), []*pile.Object{}, nil)

	context.SetTempVarByIndex(0, virtualMachine.NewInteger(42))

	err := virtualMachine.ExecutePushTemporaryVariable(context)
	if err != nil {
		t.Errorf("ExecutePushTemporaryVariable returned an error: %v", err)
	}

	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	value := context.Pop()
	// Check for immediate integer
	if pile.IsIntegerImmediate(value) {
		intValue := pile.GetIntegerImmediate(value)
		if intValue != 42 {
			t.Errorf("Expected 42 on the stack, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", value)
	}
}

func TestExecuteStoreInstanceVariable(t *testing.T) {
	virtualMachine := vm.NewVM()

	class := pile.NewClass("TestClass", nil)
	pile.AddClassInstanceVarName(class, "testVar")

	instance := pile.NewInstance(class)

	methodObj := compiler.NewMethodBuilder(class).
		StoreInstanceVariable(0).
		Go("test")

	context := vm.NewContext(methodObj, instance, []*pile.Object{}, nil)

	context.Push(virtualMachine.NewInteger(42))

	err := virtualMachine.ExecuteStoreInstanceVariable(context)
	if err != nil {
		t.Errorf("ExecuteStoreInstanceVariable returned an error: %v", err)
	}

	value := instance.GetInstanceVarByIndex(0)
	// Check for immediate integer
	if pile.IsIntegerImmediate(value) {
		intValue := pile.GetIntegerImmediate(value)
		if intValue != 42 {
			t.Errorf("Expected instance variable to be 42, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", value)
	}

	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	stackValue := context.Pop()
	// Check for immediate integer
	if pile.IsIntegerImmediate(stackValue) {
		intValue := pile.GetIntegerImmediate(stackValue)
		if intValue != 42 {
			t.Errorf("Expected 42 on the stack, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", stackValue)
	}
}

func TestExecuteStoreTemporaryVariable(t *testing.T) {
	virtualMachine := vm.NewVM()

	methodObj := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"])).
		TempVars([]string{"temp"}).
		StoreTemporaryVariable(0).
		Go("test")

	context := vm.NewContext(methodObj, pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])), []*pile.Object{}, nil)

	context.Push(virtualMachine.NewInteger(42))

	err := virtualMachine.ExecuteStoreTemporaryVariable(context)
	if err != nil {
		t.Errorf("ExecuteStoreTemporaryVariable returned an error: %v", err)
	}

	value := context.GetTempVarByIndex(0)
	// Check for immediate integer
	if pile.IsIntegerImmediate(value) {
		intValue := pile.GetIntegerImmediate(value)
		if intValue != 42 {
			t.Errorf("Expected temporary variable to be 42, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", value)
	}

	if context.StackPointer != 1 {
		t.Errorf("Expected stack pointer to be 1, got %d", context.StackPointer)
	}

	stackValue := context.Pop()
	// Check for immediate integer
	if pile.IsIntegerImmediate(stackValue) {
		intValue := pile.GetIntegerImmediate(stackValue)
		if intValue != 42 {
			t.Errorf("Expected 42 on the stack, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", stackValue)
	}
}

func TestExecuteReturnStackTop(t *testing.T) {
	virtualMachine := vm.NewVM()

	methodObj := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"])).
		Go("test")

	context := vm.NewContext(methodObj, pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])), []*pile.Object{}, nil)

	context.Push(virtualMachine.NewInteger(42))

	result, err := virtualMachine.ExecuteReturnStackTop(context)
	if err != nil {
		t.Errorf("ExecuteReturnStackTop returned an error: %v", err)
	}

	// Check for immediate integer
	if pile.IsIntegerImmediate(result) {
		intValue := pile.GetIntegerImmediate(result)
		if intValue != 42 {
			t.Errorf("Expected result to be 42, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}

	if context.StackPointer != 0 {
		t.Errorf("Expected stack pointer to be 0, got %d", context.StackPointer)
	}
}

func TestExecuteJump(t *testing.T) {
	virtualMachine := vm.NewVM()

	methodObj := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"])).
		Jump(10).
		// Add some dummy bytecodes to make the jump valid
		PushSelf().PushSelf().PushSelf().PushSelf().PushSelf().
		PushSelf().PushSelf().PushSelf().PushSelf().PushSelf().
		PushSelf().PushSelf().PushSelf().PushSelf().PushSelf().
		Go("test")

	context := vm.NewContext(methodObj, pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])), []*pile.Object{}, nil)

	skipIncrement, err := virtualMachine.ExecuteJump(context)
	if err != nil {
		t.Errorf("ExecuteJump returned an error: %v", err)
	}

	expectedPC := 0 + bytecode.InstructionSize(bytecode.JUMP) + 10
	if context.PC != expectedPC {
		t.Errorf("Expected PC to be %d, got %d", expectedPC, context.PC)
	}

	if !skipIncrement {
		t.Errorf("Expected skipIncrement to be true")
	}
}

func TestExecuteJumpIfTrue(t *testing.T) {
	virtualMachine := vm.NewVM()

	methodObj := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"])).
		JumpIfTrue(10).
		// Add some dummy bytecodes to make the jump valid
		PushSelf().PushSelf().PushSelf().PushSelf().PushSelf().
		PushSelf().PushSelf().PushSelf().PushSelf().PushSelf().
		PushSelf().PushSelf().PushSelf().PushSelf().PushSelf().
		Go("test")

	// Test with true condition
	{
		context := vm.NewContext(methodObj, pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])), []*pile.Object{}, nil)

		context.Push(virtualMachine.TrueObject)

		skipIncrement, err := virtualMachine.ExecuteJumpIfTrue(context)
		if err != nil {
			t.Errorf("ExecuteJumpIfTrue returned an error: %v", err)
		}

		expectedPC := 0 + bytecode.InstructionSize(bytecode.JUMP_IF_TRUE) + 10
		if context.PC != expectedPC {
			t.Errorf("Expected PC to be %d, got %d", expectedPC, context.PC)
		}

		if !skipIncrement {
			t.Errorf("Expected skipIncrement to be true")
		}
	}

	// Test with false condition
	{
		context := vm.NewContext(methodObj, pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])), []*pile.Object{}, nil)

		context.Push(pile.NewBoolean(false))

		skipIncrement, err := virtualMachine.ExecuteJumpIfTrue(context)
		if err != nil {
			t.Errorf("ExecuteJumpIfTrue returned an error: %v", err)
		}

		if context.PC != 0 {
			t.Errorf("Expected PC to be 0, got %d", context.PC)
		}

		if skipIncrement {
			t.Errorf("Expected skipIncrement to be false")
		}
	}
}

func TestExecuteJumpIfFalse(t *testing.T) {
	virtualMachine := vm.NewVM()

	methodObj := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"])).
		JumpIfFalse(10).
		// Add some dummy bytecodes to make the jump valid
		PushSelf().PushSelf().PushSelf().PushSelf().PushSelf().
		PushSelf().PushSelf().PushSelf().PushSelf().PushSelf().
		PushSelf().PushSelf().PushSelf().PushSelf().PushSelf().
		Go("test")

	// Test with false condition
	{
		context := vm.NewContext(methodObj, pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])), []*pile.Object{}, nil)

		context.Push(pile.NewBoolean(false))

		skipIncrement, err := virtualMachine.ExecuteJumpIfFalse(context)
		if err != nil {
			t.Errorf("ExecuteJumpIfFalse returned an error: %v", err)
		}

		expectedPC := 0 + bytecode.InstructionSize(bytecode.JUMP_IF_FALSE) + 10
		if context.PC != expectedPC {
			t.Errorf("Expected PC to be %d, got %d", expectedPC, context.PC)
		}

		if !skipIncrement {
			t.Errorf("Expected skipIncrement to be true")
		}
	}

	// Test with true condition
	{
		context := vm.NewContext(methodObj, pile.ClassToObject(pile.ObjectToClass(virtualMachine.Globals["Object"])), []*pile.Object{}, nil)

		context.Push(pile.NewBoolean(true))

		skipIncrement, err := virtualMachine.ExecuteJumpIfFalse(context)
		if err != nil {
			t.Errorf("ExecuteJumpIfFalse returned an error: %v", err)
		}

		if context.PC != 0 {
			t.Errorf("Expected PC to be 0, got %d", context.PC)
		}

		if skipIncrement {
			t.Errorf("Expected skipIncrement to be false")
		}
	}
}
