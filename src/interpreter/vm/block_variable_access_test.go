package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/bytecode"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/runtime"
	"smalltalklsp/interpreter/vm"
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
	virtualMachine := vm.NewVM()
	runtime.RegisterBlockExecutor(virtualMachine)

	// The block bytecodes ([a := 2])
	blockBytecodes := []byte{
		// Push the literal 2
		bytecode.PUSH_LITERAL,
		0, 0, 0, 0, // literal index 0 (the value 2)

		// Store it in the outer context's temporary variable 'a'
		bytecode.STORE_TEMPORARY_VARIABLE,
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
		bytecode.PUSH_LITERAL,
		0, 0, 0, 0, // literal index 0 (the value 1)
		bytecode.STORE_TEMPORARY_VARIABLE,
		0, 0, 0, 0, // temp var index 0 (a)

		// Create a block that assigns 2 to 'a'
		bytecode.CREATE_BLOCK,
		0, 0, 0, byte(len(blockBytecodes)), // bytecode size
		0, 0, 0, 1, // literal count (1 literal)
		0, 0, 0, 0, // temp var count (0 temp vars)

		// Execute the block
		bytecode.EXECUTE_BLOCK,
		0, 0, 0, 0, // arg count 0

		// Return the value of 'a'
		bytecode.PUSH_TEMPORARY_VARIABLE,
		0, 0, 0, 0, // temp var index 0 (a)
		bytecode.RETURN_STACK_TOP,
	}

	// Create a method with the bytecodes
	method := &pile.Method{
		Object: pile.Object{
			TypeField: pile.OBJ_METHOD,
		},
		Bytecodes: methodBytecodes,
		Literals: []*pile.Object{
			pile.MakeIntegerImmediate(1), // The literal 1
			pile.MakeIntegerImmediate(2), // The literal 2
		},
		TempVarNames: []string{"a"}, // One temporary variable 'a'
	}

	// Create a context
	context := vm.NewContext(
		pile.MethodToObject(method),
		pile.MakeNilImmediate(),
		[]*pile.Object{},
		nil,
	)

	// Set the VM's current context
	virtualMachine.Executor.CurrentContext = context

	// Execute the PUSH_LITERAL and STORE_TEMPORARY_VARIABLE bytecodes
	// This sets a := 1
	err := virtualMachine.ExecutePushLiteral(context)
	if err != nil {
		t.Fatalf("ExecutePushLiteral returned an error: %v", err)
	}
	context.PC += bytecode.InstructionSize(bytecode.PUSH_LITERAL)

	err = virtualMachine.ExecuteStoreTemporaryVariable(context)
	if err != nil {
		t.Fatalf("ExecuteStoreTemporaryVariable returned an error: %v", err)
	}
	context.PC += bytecode.InstructionSize(bytecode.STORE_TEMPORARY_VARIABLE)

	// Execute the CREATE_BLOCK bytecode
	err = virtualMachine.ExecuteCreateBlock(context)
	if err != nil {
		t.Fatalf("ExecuteCreateBlock returned an error: %v", err)
	}

	// Get the block from the stack
	block := context.Pop()
	blockObj := pile.ObjectToBlock(block)

	// Set the block's bytecodes
	blockObj.SetBytecodes(blockBytecodes)

	// Set the block's literals
	blockObj.Literals = []*pile.Object{
		pile.MakeIntegerImmediate(2), // The literal 2
	}

	// Push the block back onto the stack
	context.Push(block)

	// Advance the PC to the EXECUTE_BLOCK bytecode
	context.PC += bytecode.InstructionSize(bytecode.CREATE_BLOCK)

	// Execute the EXECUTE_BLOCK bytecode
	_, err = virtualMachine.ExecuteExecuteBlock(context)
	if err != nil {
		t.Fatalf("ExecuteExecuteBlock returned an error: %v", err)
	}

	// Advance the PC to the PUSH_TEMPORARY_VARIABLE bytecode
	context.PC += bytecode.InstructionSize(bytecode.EXECUTE_BLOCK)

	// Execute the PUSH_TEMPORARY_VARIABLE bytecode
	err = virtualMachine.ExecutePushTemporaryVariable(context)
	if err != nil {
		t.Fatalf("ExecutePushTemporaryVariable returned an error: %v", err)
	}

	// Advance the PC to the RETURN_STACK_TOP bytecode
	context.PC += bytecode.InstructionSize(bytecode.PUSH_TEMPORARY_VARIABLE)

	// Execute the RETURN_STACK_TOP bytecode
	result, err := virtualMachine.ExecuteReturnStackTop(context)
	if err != nil {
		t.Fatalf("ExecuteReturnStackTop returned an error: %v", err)
	}

	// Check that the result is 2
	if !pile.IsIntegerImmediate(result) {
		t.Errorf("Result is not an integer: %v", result)
	}

	value := pile.GetIntegerImmediate(result)
	if value != 2 {
		t.Errorf("Result = %d, want 2", value)
	}
}
