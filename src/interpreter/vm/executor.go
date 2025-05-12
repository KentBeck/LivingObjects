package vm

import (
	"fmt"

	"smalltalklsp/interpreter/bytecode"
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// Executor handles the execution of bytecode in contexts
type Executor struct {
	// VM is a reference to the virtual machine
	VM *VM

	// CurrentContext is the context currently being executed
	CurrentContext *Context
}

// NewExecutor creates a new executor
func NewExecutor(vm *VM) *Executor {
	return &Executor{
		VM: vm,
	}
}

// Execute executes the current context
func (e *Executor) Execute() (core.ObjectInterface, error) {
	var finalResult core.ObjectInterface

	for e.CurrentContext != nil {
		// Execute the current context
		result, err := e.ExecuteContext(e.CurrentContext)
		if err != nil {
			return nil, err
		}

		// Save the result if this is the top-level context
		if e.CurrentContext.Sender == nil {
			finalResult = result
		}

		// Move to the sender context
		e.CurrentContext = e.CurrentContext.Sender

		// If we have a sender, push the result onto its stack
		if e.CurrentContext != nil {
			e.CurrentContext.Push(result)
		}
	}

	return finalResult, nil
}

// ExecuteContext executes a single context until it returns
func (e *Executor) ExecuteContext(context *Context) (core.ObjectInterface, error) {
	// Execute the context
	for {
		// Get the method
		method := classes.ObjectToMethod(context.Method)

		// Check if we've reached the end of the method
		if context.PC >= len(method.GetBytecodes()) {
			// Reached end of bytecode array

			// If we've reached the end of the method, return the top of the stack
			// This handles the case where we jump to the end of the bytecode array
			if context.StackPointer > 0 {
				returnValue := context.Pop()
				return returnValue, nil
			}
			return e.VM.NilObject, nil
		}

		// Get the current bytecode
		opcode := method.GetBytecodes()[context.PC]

		// Get the instruction size
		size := bytecode.InstructionSize(opcode)

		// Execute the bytecode
		var err error
		var skipIncrement bool

		switch opcode {
		case bytecode.PUSH_LITERAL:
			err = e.VM.ExecutePushLiteral(context)

		case bytecode.PUSH_INSTANCE_VARIABLE:
			err = e.VM.ExecutePushInstanceVariable(context)

		case bytecode.PUSH_TEMPORARY_VARIABLE:
			err = e.VM.ExecutePushTemporaryVariable(context)

		case bytecode.PUSH_SELF:
			err = e.VM.ExecutePushSelf(context)

		case bytecode.STORE_INSTANCE_VARIABLE:
			err = e.VM.ExecuteStoreInstanceVariable(context)

		case bytecode.STORE_TEMPORARY_VARIABLE:
			err = e.VM.ExecuteStoreTemporaryVariable(context)

		case bytecode.SEND_MESSAGE:
			returnValue, err := e.VM.ExecuteSendMessage(context)
			if err == nil {
				if returnValue != nil {
					// We got a result from a primitive method
					// Continue execution in the current context
					context.PC += size
					continue
				} else {
					// A nil return value with no error means we've started a new context
					return e.VM.NilObject, nil
				}
			}

		case bytecode.RETURN_STACK_TOP:
			returnValue, err := e.VM.ExecuteReturnStackTop(context)
			if err == nil {
				return returnValue, nil
			}

		case bytecode.JUMP:
			skipIncrement, err = e.VM.ExecuteJump(context)
			if err == nil && skipIncrement {
				continue
			}

		case bytecode.JUMP_IF_TRUE:
			skipIncrement, err = e.VM.ExecuteJumpIfTrue(context)
			if err == nil && skipIncrement {
				continue
			}

		case bytecode.JUMP_IF_FALSE:
			skipIncrement, err = e.VM.ExecuteJumpIfFalse(context)
			if err == nil && skipIncrement {
				continue
			}

		case bytecode.POP:
			err = e.VM.ExecutePop(context)

		case bytecode.DUPLICATE:
			err = e.VM.ExecuteDuplicate(context)

		case bytecode.CREATE_BLOCK:
			err = e.VM.ExecuteCreateBlock(context)

		case bytecode.EXECUTE_BLOCK:
			returnValue, err := e.VM.ExecuteExecuteBlock(context)
			if err == nil {
				if returnValue != nil {
					// We got a result from executing the block
					// Continue execution in the current context
					context.PC += size
					continue
				} else {
					// A nil return value with no error means we've started a new context
					return e.VM.NilObject, nil
				}
			}

		default:
			return nil, fmt.Errorf("unknown bytecode: %d", opcode)
		}

		// Check for errors
		if err != nil {
			return nil, err
		}

		// Increment the PC
		context.PC += size
	}
}
