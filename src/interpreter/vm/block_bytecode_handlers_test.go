package vm

import (
	"testing"

	"smalltalklsp/interpreter/bytecode"
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

func TestExecuteCreateBlock(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method with a CREATE_BLOCK bytecode
	method := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes: []byte{
			bytecode.CREATE_BLOCK,
			0, 0, 0, 10, // bytecode size
			0, 0, 0, 2, // literal count
			0, 0, 0, 3, // temp var count
		},
		Literals:     []*core.Object{},
		TempVarNames: []string{},
	}

	// Create a context
	context := NewContext(
		classes.MethodToObject(method),
		core.MakeNilImmediate(),
		[]*core.Object{},
		nil,
	)

	// Execute the CREATE_BLOCK bytecode
	err := vm.ExecuteCreateBlock(context)
	if err != nil {
		t.Errorf("ExecuteCreateBlock returned an error: %v", err)
	}

	// Check that a block was pushed onto the stack
	if context.StackPointer != 1 {
		t.Errorf("Stack pointer = %d, want 1", context.StackPointer)
	}

	// Check that the block is valid
	block := context.Pop()
	if block == nil {
		t.Errorf("Block is nil")
	}

	if block.Type() != core.OBJ_BLOCK {
		t.Errorf("Block type = %d, want %d", block.Type(), core.OBJ_BLOCK)
	}

	// Convert to a Block
	blockObj := classes.ObjectToBlock(block)
	if blockObj == nil {
		t.Errorf("Failed to convert to Block")
	}

	// Check the block's bytecodes
	if len(blockObj.GetBytecodes()) != 10 {
		t.Errorf("Block bytecode size = %d, want 10", len(blockObj.GetBytecodes()))
	}

	// Check the block's literals
	if len(blockObj.GetLiterals()) != 2 {
		t.Errorf("Block literal count = %d, want 2", len(blockObj.GetLiterals()))
	}

	// Check the block's temp var names
	if len(blockObj.GetTempVarNames()) != 3 {
		t.Errorf("Block temp var count = %d, want 3", len(blockObj.GetTempVarNames()))
	}
}

func TestExecuteExecuteBlock(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create a method with an EXECUTE_BLOCK bytecode
	method := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes: []byte{
			bytecode.EXECUTE_BLOCK,
			0, 0, 0, 2, // arg count
		},
		Literals:     []*core.Object{},
		TempVarNames: []string{},
	}

	// Create a context
	context := NewContext(
		classes.MethodToObject(method),
		core.MakeNilImmediate(),
		[]*core.Object{},
		nil,
	)

	// Create a block with proper class field
	block := classes.ObjectToBlock(vm.NewBlock(context))

	// Push the block onto the stack
	context.Push(classes.BlockToObject(block))

	// Push some arguments onto the stack
	context.Push(core.MakeIntegerImmediate(1))
	context.Push(core.MakeIntegerImmediate(2))

	// Execute the EXECUTE_BLOCK bytecode
	result, err := vm.ExecuteExecuteBlock(context)
	if err != nil {
		t.Errorf("ExecuteExecuteBlock returned an error: %v", err)
	}

	// Check that the result is nil (since the block doesn't do anything)
	if !core.IsNilImmediate(result) {
		t.Errorf("Result = %v, want nil", result)
	}

	// Check that the result was pushed onto the stack
	if context.StackPointer != 1 {
		t.Errorf("Stack pointer = %d, want 1", context.StackPointer)
	}

	// Check that the result on the stack is nil
	stackResult := context.Pop()
	if !core.IsNilImmediate(stackResult) {
		t.Errorf("Stack result = %v, want nil", stackResult)
	}
}
