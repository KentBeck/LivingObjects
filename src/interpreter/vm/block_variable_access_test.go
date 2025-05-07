package vm

import (
	"testing"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/runtime"
)

// TestBlockModifiesLocalVariable tests a method that:
// 1. Declares a local variable 'a'
// 2. Assigns 1 to 'a'
// 3. Creates a block that assigns 2 to 'a'
// 4. Executes the block
// 5. Returns 'a'
// The expected result is 2, showing that blocks can modify variables in their outer context.
func TestBlockModifiesLocalVariable(t *testing.T) {
	// Create a VM and register it as a block executor
	vm := NewVM()
	runtime.RegisterBlockExecutor(vm)

	// The block bytecodes ([a := 2])
	blockBytecodes := []byte{
		// Push the literal 2
		PUSH_LITERAL,
		0, 0, 0, 0, // literal index 0 (the value 2)

		// Store it in the outer context's temporary variable 'a'
		STORE_TEMPORARY_VARIABLE,
		0, 0, 0, 0, // temp var index 0 (a)

		// Return nil (implicit)
	}

	// Create a method with the following Smalltalk code:
	// method
	//   | a |
	//   a := 1.
	//   [a := 2] value.
	//   ^a
	methodBytecodes := []byte{
		// Store 1 in temporary variable 'a'
		PUSH_LITERAL,
		0, 0, 0, 0, // literal index 0 (the value 1)
		STORE_TEMPORARY_VARIABLE,
		0, 0, 0, 0, // temp var index 0 (a)

		// Create a block that assigns 2 to 'a'
		CREATE_BLOCK,
		0, 0, 0, byte(len(blockBytecodes)), // bytecode size
		0, 0, 0, 1, // literal count (1 literal)
		0, 0, 0, 0, // temp var count (0 temp vars)

		// Execute the block
		EXECUTE_BLOCK,
		0, 0, 0, 0, // arg count 0

		// Return the value of 'a'
		PUSH_TEMPORARY_VARIABLE,
		0, 0, 0, 0, // temp var index 0 (a)
		RETURN_STACK_TOP,
	}

	// Create a method with the bytecodes
	method := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes: methodBytecodes,
		Literals: []*core.Object{
			core.MakeIntegerImmediate(1), // The literal 1
			core.MakeIntegerImmediate(2), // The literal 2
		},
		TempVarNames: []string{"a"}, // One temporary variable 'a'
	}

	// Create a context
	context := NewContext(
		classes.MethodToObject(method),
		core.MakeNilImmediate(),
		[]*core.Object{},
		nil,
	)

	// Set the VM's current context
	vm.CurrentContext = context

	// Execute the PUSH_LITERAL and STORE_TEMPORARY_VARIABLE bytecodes
	// This sets a := 1
	err := vm.ExecutePushLiteral(context)
	if err != nil {
		t.Fatalf("ExecutePushLiteral returned an error: %v", err)
	}
	context.PC += InstructionSize(PUSH_LITERAL)

	err = vm.ExecuteStoreTemporaryVariable(context)
	if err != nil {
		t.Fatalf("ExecuteStoreTemporaryVariable returned an error: %v", err)
	}
	context.PC += InstructionSize(STORE_TEMPORARY_VARIABLE)

	// Execute the CREATE_BLOCK bytecode
	err = vm.ExecuteCreateBlock(context)
	if err != nil {
		t.Fatalf("ExecuteCreateBlock returned an error: %v", err)
	}

	// Get the block from the stack
	block := context.Pop()
	blockObj := classes.ObjectToBlock(block)

	// Set the block's bytecodes
	blockObj.SetBytecodes(blockBytecodes)

	// Set the block's literals
	blockObj.Literals = []*core.Object{
		core.MakeIntegerImmediate(2), // The literal 2
	}

	// Push the block back onto the stack
	context.Push(block)

	// Advance the PC to the EXECUTE_BLOCK bytecode
	context.PC += InstructionSize(CREATE_BLOCK)

	// Execute the EXECUTE_BLOCK bytecode
	_, err = vm.ExecuteExecuteBlock(context)
	if err != nil {
		t.Fatalf("ExecuteExecuteBlock returned an error: %v", err)
	}

	// Advance the PC to the PUSH_TEMPORARY_VARIABLE bytecode
	context.PC += InstructionSize(EXECUTE_BLOCK)

	// Execute the PUSH_TEMPORARY_VARIABLE bytecode
	err = vm.ExecutePushTemporaryVariable(context)
	if err != nil {
		t.Fatalf("ExecutePushTemporaryVariable returned an error: %v", err)
	}

	// Advance the PC to the RETURN_STACK_TOP bytecode
	context.PC += InstructionSize(PUSH_TEMPORARY_VARIABLE)

	// Execute the RETURN_STACK_TOP bytecode
	result, err := vm.ExecuteReturnStackTop(context)
	if err != nil {
		t.Fatalf("ExecuteReturnStackTop returned an error: %v", err)
	}

	// Check that the result is 2
	if !core.IsIntegerImmediate(result) {
		t.Errorf("Result is not an integer: %v", result)
	}

	value := core.GetIntegerImmediate(result)
	if value != 2 {
		t.Errorf("Result = %d, want 2", value)
	}
}
