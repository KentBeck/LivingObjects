package vm

import (
	"smalltalklsp/interpreter/bytecode"
)

// Re-export bytecode constants for backward compatibility
const (
	PUSH_LITERAL             = bytecode.PUSH_LITERAL
	PUSH_INSTANCE_VARIABLE   = bytecode.PUSH_INSTANCE_VARIABLE
	PUSH_TEMPORARY_VARIABLE  = bytecode.PUSH_TEMPORARY_VARIABLE
	PUSH_SELF                = bytecode.PUSH_SELF
	STORE_INSTANCE_VARIABLE  = bytecode.STORE_INSTANCE_VARIABLE
	STORE_TEMPORARY_VARIABLE = bytecode.STORE_TEMPORARY_VARIABLE
	SEND_MESSAGE             = bytecode.SEND_MESSAGE
	RETURN_STACK_TOP         = bytecode.RETURN_STACK_TOP
	JUMP                     = bytecode.JUMP
	JUMP_IF_TRUE             = bytecode.JUMP_IF_TRUE
	JUMP_IF_FALSE            = bytecode.JUMP_IF_FALSE
	POP                      = bytecode.POP
	DUPLICATE                = bytecode.DUPLICATE
	CREATE_BLOCK             = bytecode.CREATE_BLOCK
	EXECUTE_BLOCK            = bytecode.EXECUTE_BLOCK
)

// InstructionSize returns the size of the instruction in bytes (including the opcode)
func InstructionSize(code byte) int {
	return bytecode.InstructionSize(code)
}

// BytecodeName returns the name of the bytecode
func BytecodeName(code byte) string {
	return bytecode.BytecodeName(code)
}

// ReadUint32 reads a 4-byte unsigned integer from the bytecode array
func ReadUint32(bytecodeArray []byte, pc int) uint32 {
	return bytecode.ReadUint32(bytecodeArray, pc)
}

// WriteUint32 writes a 4-byte unsigned integer to the bytecode array
func WriteUint32(bytecodeArray []byte, pc int, value uint32) {
	bytecode.WriteUint32(bytecodeArray, pc, value)
}
