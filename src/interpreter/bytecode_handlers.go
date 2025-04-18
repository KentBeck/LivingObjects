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
	fmt.Printf("PUSH_LITERAL: %s\n", literal)
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

	// Get the instance variable name
	name := context.Receiver.Class.InstanceVarNames[index]

	// Push the instance variable onto the stack
	if value, ok := context.Receiver.InstanceVars[name]; ok {
		context.Push(value)
	} else {
		context.Push(vm.NilObject)
	}
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

	// Get the instance variable name
	name := context.Receiver.Class.InstanceVarNames[index]

	// Pop the value from the stack
	value := context.Pop()

	// Store the value in the instance variable
	context.Receiver.InstanceVars[name] = value

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

	fmt.Printf("SEND_MESSAGE: %s %s with %d args\n", receiver, selector.SymbolValue, argCount)

	// Handle primitive methods
	if result := vm.executePrimitive(receiver, selector, args); result != nil {
		fmt.Printf("PRIMITIVE RESULT: %s\n", result)
		context.Push(result)
		return result, nil
	}

	// Look up the method
	method := vm.lookupMethod(receiver, selector)
	if method == nil {
		return nil, fmt.Errorf("method not found: %s", selector.SymbolValue)
	}

	// Create a new context for the method
	newContext := NewContext(method, receiver, args, context)

	// Set the current context to the new context
	vm.CurrentContext = newContext

	// Return from this context execution to start executing the new context
	return nil, nil
}

// ExecuteReturnStackTop executes the RETURN_STACK_TOP bytecode
func (vm *VM) ExecuteReturnStackTop(context *Context) (*Object, error) {
	// Pop the return value from the stack
	returnValue := context.Pop()

	fmt.Printf("RETURN_STACK_TOP: %s\n", returnValue)

	// Return the value
	return returnValue, nil
}

// ExecuteJump executes the JUMP bytecode
func (vm *VM) ExecuteJump(context *Context) (bool, error) {
	// Get the jump target (4 bytes)
	target := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

	// Set the PC to the target
	context.PC = target

	// Skip the normal PC increment
	return true, nil
}

// ExecuteJumpIfTrue executes the JUMP_IF_TRUE bytecode
func (vm *VM) ExecuteJumpIfTrue(context *Context) (bool, error) {
	// Get the jump target (4 bytes)
	target := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

	// Pop the condition from the stack
	condition := context.Pop()

	// If the condition is true, jump to the target
	if condition.IsTrue() {
		context.PC = target
		return true, nil
	}

	return false, nil
}

// ExecuteJumpIfFalse executes the JUMP_IF_FALSE bytecode
func (vm *VM) ExecuteJumpIfFalse(context *Context) (bool, error) {
	// Get the jump target (4 bytes)
	target := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

	// Pop the condition from the stack
	condition := context.Pop()

	// If the condition is false, jump to the target
	if !condition.IsTrue() {
		context.PC = target
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
