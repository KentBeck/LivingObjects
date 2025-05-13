package vm

import (
	"encoding/binary"
	"fmt"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// ExecuteCreateBlock executes the CREATE_BLOCK bytecode
func (vm *VM) ExecuteCreateBlock(context *Context) error {
	// Get the method
	method := classes.ObjectToMethod(context.Method)

	// Get the bytecode size (4 bytes)
	bytecodeSize := int(binary.BigEndian.Uint32(method.GetBytecodes()[context.PC+1:]))

	// Get the literal count (4 bytes)
	literalCount := int(binary.BigEndian.Uint32(method.GetBytecodes()[context.PC+5:]))

	// Get the temp var count (4 bytes)
	tempVarCount := int(binary.BigEndian.Uint32(method.GetBytecodes()[context.PC+9:]))

	// Create a new block
	block := classes.ObjectToBlock(vm.NewBlock(context))

	// Set the bytecodes
	// In a real implementation, we would extract the bytecodes from the method
	// For now, we'll just create an empty bytecode array
	block.SetBytecodes(make([]byte, bytecodeSize))

	// Set the literals
	// In a real implementation, we would extract the literals from the method
	// For now, we'll just create an empty literal array
	for i := 0; i < literalCount; i++ {
		block.AddLiteral(core.MakeNilImmediate())
	}

	// Set the temporary variable names
	// In a real implementation, we would extract the temp var names from the method
	// For now, we'll just create empty temp var names
	for i := 0; i < tempVarCount; i++ {
		block.AddTempVarName(fmt.Sprintf("temp%d", i))
	}

	// Push the block onto the stack
	context.Push(classes.BlockToObject(block))

	return nil
}

// ExecuteExecuteBlock executes the EXECUTE_BLOCK bytecode
func (vm *VM) ExecuteExecuteBlock(context *Context) (*core.Object, error) {
	// Get the method
	method := classes.ObjectToMethod(context.Method)

	// Get the argument count (4 bytes)
	argCount := int(binary.BigEndian.Uint32(method.GetBytecodes()[context.PC+1:]))

	// Pop the arguments from the stack
	args := make([]*core.Object, argCount)
	for i := argCount - 1; i >= 0; i-- {
		args[i] = context.Pop()
	}

	// Pop the block from the stack
	blockObj := context.Pop()
	if blockObj == nil {
		return nil, fmt.Errorf("nil block")
	}

	// Check if it's a block
	if blockObj.Type() != core.OBJ_BLOCK {
		return nil, fmt.Errorf("not a block: %v", blockObj)
	}

	// Convert to a Block
	block := classes.ObjectToBlock(blockObj)

	// Execute the block
	result := block.ValueWithArguments(args)

	// Push the result onto the stack
	context.Push(result)

	return result, nil
}
