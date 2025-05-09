package bytecode

// Bytecode constants
const (
	// Bytecodes
	PUSH_LITERAL             byte = 0  // Push a literal from the literals array (followed by 4-byte index)
	PUSH_INSTANCE_VARIABLE   byte = 1  // Push an instance variable value (followed by 4-byte offset)
	PUSH_TEMPORARY_VARIABLE  byte = 2  // Push a temporary variable value (followed by 4-byte offset)
	PUSH_SELF                byte = 3  // Push self onto the stack
	STORE_INSTANCE_VARIABLE  byte = 4  // Store a value into an instance variable (followed by 4-byte offset)
	STORE_TEMPORARY_VARIABLE byte = 5  // Store a value into a temporary variable (followed by 4-byte offset)
	SEND_MESSAGE             byte = 6  // Send a message (followed by 4-byte selector index and 4-byte arg count)
	RETURN_STACK_TOP         byte = 7  // Return the value on top of the stack
	JUMP                     byte = 8  // Jump to a different bytecode (followed by 4-byte target)
	JUMP_IF_TRUE             byte = 9  // Jump if top of stack is true (followed by 4-byte target)
	JUMP_IF_FALSE            byte = 10 // Jump if top of stack is false (followed by 4-byte target)
	POP                      byte = 11 // Pop the top value from the stack
	DUPLICATE                byte = 12 // Duplicate the top value on the stack
	CREATE_BLOCK             byte = 13 // Create a block (followed by 4-byte bytecode size, 4-byte literal count, 4-byte temp var count)
	EXECUTE_BLOCK            byte = 14 // Execute a block (followed by 4-byte arg count)
)

// InstructionSize returns the size of the instruction in bytes (including the opcode)
func InstructionSize(bytecode byte) int {
	switch bytecode {
	case PUSH_LITERAL, PUSH_INSTANCE_VARIABLE, PUSH_TEMPORARY_VARIABLE,
		STORE_INSTANCE_VARIABLE, STORE_TEMPORARY_VARIABLE,
		JUMP, JUMP_IF_TRUE, JUMP_IF_FALSE:
		return 5 // 1 byte opcode + 4 byte operand
	case SEND_MESSAGE:
		return 9 // 1 byte opcode + 4 byte selector index + 4 byte arg count
	case CREATE_BLOCK:
		return 13 // 1 byte opcode + 4 byte bytecode size + 4 byte literal count + 4 byte temp var count
	case EXECUTE_BLOCK:
		return 5 // 1 byte opcode + 4 byte arg count
	case PUSH_SELF, RETURN_STACK_TOP, POP, DUPLICATE:
		return 1 // 1 byte opcode
	default:
		return 1 // Default to 1 byte for unknown bytecodes
	}
}

// BytecodeName returns the name of the bytecode
func BytecodeName(bytecode byte) string {
	switch bytecode {
	case PUSH_LITERAL:
		return "PUSH_LITERAL"
	case PUSH_INSTANCE_VARIABLE:
		return "PUSH_INSTANCE_VARIABLE"
	case PUSH_TEMPORARY_VARIABLE:
		return "PUSH_TEMPORARY_VARIABLE"
	case PUSH_SELF:
		return "PUSH_SELF"
	case STORE_INSTANCE_VARIABLE:
		return "STORE_INSTANCE_VARIABLE"
	case STORE_TEMPORARY_VARIABLE:
		return "STORE_TEMPORARY_VARIABLE"
	case SEND_MESSAGE:
		return "SEND_MESSAGE"
	case RETURN_STACK_TOP:
		return "RETURN_STACK_TOP"
	case JUMP:
		return "JUMP"
	case JUMP_IF_TRUE:
		return "JUMP_IF_TRUE"
	case JUMP_IF_FALSE:
		return "JUMP_IF_FALSE"
	case POP:
		return "POP"
	case DUPLICATE:
		return "DUPLICATE"
	case CREATE_BLOCK:
		return "CREATE_BLOCK"
	case EXECUTE_BLOCK:
		return "EXECUTE_BLOCK"
	default:
		return "UNKNOWN"
	}
}

// ReadUint32 reads a 4-byte unsigned integer from the bytecode array
func ReadUint32(bytecode []byte, pc int) uint32 {
	return uint32(bytecode[pc]) |
		uint32(bytecode[pc+1])<<8 |
		uint32(bytecode[pc+2])<<16 |
		uint32(bytecode[pc+3])<<24
}

// WriteUint32 writes a 4-byte unsigned integer to the bytecode array
func WriteUint32(bytecode []byte, pc int, value uint32) {
	bytecode[pc] = byte(value)
	bytecode[pc+1] = byte(value >> 8)
	bytecode[pc+2] = byte(value >> 16)
	bytecode[pc+3] = byte(value >> 24)
}
