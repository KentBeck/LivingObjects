package main

import (
	"encoding/binary"
	"fmt"
)

// ExecutePushLiteral executes the PUSH_LITERAL bytecode
func (vm *VM) ExecutePushLiteral(context *Context) error {
	// Get the literal index (4 bytes)
	index := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))
	if index < 0 || index >= len(context.Method.Method.Literals) {
		return fmt.Errorf("literal index out of bounds: %d", index)
	}

	// Push the literal onto the stack
	literal := context.Method.Method.Literals[index]

	context.Push(literal)
	return nil
}

// ExecutePushInstanceVariable executes the PUSH_INSTANCE_VARIABLE bytecode
func (vm *VM) ExecutePushInstanceVariable(context *Context) error {
	// Get the instance variable index (4 bytes)
	index := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))
	if index < 0 || index >= len(context.Receiver.Class.InstanceVarNames) {
		return fmt.Errorf("instance variable index out of bounds: %d", index)
	}

	// Push the instance variable onto the stack
	value := context.Receiver.GetInstanceVarByIndex(index)
	context.Push(value)
	return nil
}

// ExecutePushTemporaryVariable executes the PUSH_TEMPORARY_VARIABLE bytecode
func (vm *VM) ExecutePushTemporaryVariable(context *Context) error {
	// Get the temporary variable index (4 bytes)
	index := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

	// Push the temporary variable onto the stack
	context.Push(context.GetTempVarByIndex(index))
	return nil
}

// ExecutePushSelf executes the PUSH_SELF bytecode
func (vm *VM) ExecutePushSelf(context *Context) error {
	// Push the receiver onto the stack
	context.Push(context.Receiver)
	return nil
}

// ExecuteStoreInstanceVariable executes the STORE_INSTANCE_VARIABLE bytecode
func (vm *VM) ExecuteStoreInstanceVariable(context *Context) error {
	// Get the instance variable index (4 bytes)
	index := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))
	if index < 0 || index >= len(context.Receiver.Class.InstanceVarNames) {
		return fmt.Errorf("instance variable index out of bounds: %d", index)
	}

	// Pop the value from the stack
	value := context.Pop()

	// Store the value in the instance variable
	context.Receiver.SetInstanceVarByIndex(index, value)

	// Push the value back onto the stack
	context.Push(value)
	return nil
}

// ExecuteStoreTemporaryVariable executes the STORE_TEMPORARY_VARIABLE bytecode
func (vm *VM) ExecuteStoreTemporaryVariable(context *Context) error {
	// Get the temporary variable index (4 bytes)
	index := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

	// Pop the value from the stack
	value := context.Pop()

	// Store the value in the temporary variable
	context.SetTempVarByIndex(index, value)

	// Push the value back onto the stack
	context.Push(value)
	return nil
}

// ExecuteSendMessage executes the SEND_MESSAGE bytecode
func (vm *VM) ExecuteSendMessage(context *Context) (*Object, error) {
	// Get the selector index (4 bytes)
	selectorIndex := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))
	if selectorIndex < 0 || selectorIndex >= len(context.Method.Method.Literals) {
		return nil, fmt.Errorf("selector index out of bounds: %d", selectorIndex)
	}

	// Get the argument count (4 bytes)
	argCount := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+5:]))

	// Get the selector
	selector := context.Method.Method.Literals[selectorIndex]
	if selector.Type != OBJ_SYMBOL {
		return nil, fmt.Errorf("selector is not a symbol: %s", selector)
	}

	// Pop the arguments from the stack
	args := make([]*Object, argCount)
	for i := argCount - 1; i >= 0; i-- {
		args[i] = context.Pop()
	}

	// Pop the receiver
	receiver := context.Pop()

	// Handle primitive methods
	if result := vm.executePrimitive(receiver, selector, args); result != nil {

		context.Push(result)
		return result, nil
	}

	// Check for nil receiver
	if receiver == nil {
		return nil, fmt.Errorf("nil receiver for message: %s", selector.SymbolValue)
	}
	method := vm.lookupMethod(receiver, selector)
	if method == nil {
		return nil, fmt.Errorf("method not found: %s", selector.SymbolValue)
	}

	// Create a new context for the method
	newContext := NewContext(method, receiver, args, context)

	// Set the current context to the new context
	vm.CurrentContext = newContext

	// Return from this context execution to start executing the new context
	// We need to execute the new context immediately
	result, err := vm.ExecuteContext(newContext)
	if err != nil {
		return nil, err
	}

	// Move back to the sender context
	vm.CurrentContext = context

	// Push the result onto the stack
	context.Push(result)

	// Return the result
	return result, nil
}

// ExecuteReturnStackTop executes the RETURN_STACK_TOP bytecode
func (vm *VM) ExecuteReturnStackTop(context *Context) (*Object, error) {
	// Pop the return value from the stack
	returnValue := context.Pop()

	// Return the value
	return returnValue, nil
}

// ExecuteJump executes the JUMP bytecode
func (vm *VM) ExecuteJump(context *Context) (bool, error) {
	// Get the jump offset (4 bytes)
	offset := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

	// Calculate the new PC by adding the offset to the current PC plus the size of this instruction
	newPC := context.PC + InstructionSize(JUMP) + offset

	// Set the PC to the new position
	context.PC = newPC

	// Skip the normal PC increment
	return true, nil
}

// ExecuteJumpIfTrue executes the JUMP_IF_TRUE bytecode
func (vm *VM) ExecuteJumpIfTrue(context *Context) (bool, error) {
	// Get the jump offset (4 bytes)
	offset := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

	// Pop the condition from the stack
	condition := context.Pop()

	// If the condition is true, jump by the offset
	if condition.IsTrue() {
		// Calculate the new PC by adding the offset to the current PC plus the size of this instruction
		newPC := context.PC + InstructionSize(JUMP_IF_TRUE) + offset
		context.PC = newPC
		return true, nil
	}

	return false, nil
}

// ExecuteJumpIfFalse executes the JUMP_IF_FALSE bytecode
func (vm *VM) ExecuteJumpIfFalse(context *Context) (bool, error) {
	// Get the jump offset (4 bytes)
	offset := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

	// Pop the condition from the stack
	condition := context.Pop()

	// If the condition is false, jump by the offset
	if !condition.IsTrue() {
		// Calculate the new PC by adding the offset to the current PC plus the size of this instruction
		newPC := context.PC + InstructionSize(JUMP_IF_FALSE) + offset
		context.PC = newPC
		return true, nil
	}

	return false, nil
}

// ExecutePop executes the POP bytecode
func (vm *VM) ExecutePop(context *Context) error {
	// Pop the top value from the stack
	context.Pop()
	return nil
}

// ExecuteDuplicate executes the DUPLICATE bytecode
func (vm *VM) ExecuteDuplicate(context *Context) error {
	// Duplicate the top value on the stack
	value := context.Top()

	context.Push(value)
	return nil
}

// ExecuteSetClass executes the SET_CLASS bytecode
func (vm *VM) ExecuteSetClass(context *Context) error {
	// Get the class index (4 bytes)
	classIndex := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))
	if classIndex < 0 || classIndex >= len(context.Method.Method.Literals) {
		return fmt.Errorf("class index out of bounds: %d", classIndex)
	}

	// Get the class
	class := context.Method.Method.Literals[classIndex]
	if class.Type != OBJ_CLASS {
		return fmt.Errorf("literal is not a class: %s", class)
	}

	// Get the top value on the stack
	value := context.Pop()
	if value == nil {
		return fmt.Errorf("stack underflow")
	}

	// Set the class of the value
	value.Class = class

	// Push the value back onto the stack
	context.Push(value)
	return nil
}
