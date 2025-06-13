#pragma once

#include <cstdint>

namespace smalltalk {

// Bytecode constants based on Go implementation
enum class Bytecode : uint8_t {
    PUSH_LITERAL             = 0,  // Push a literal from the literals array (followed by 4-byte index)
    PUSH_INSTANCE_VARIABLE   = 1,  // Push an instance variable value (followed by 4-byte offset)
    PUSH_TEMPORARY_VARIABLE  = 2,  // Push a temporary variable value (followed by 4-byte offset)
    PUSH_SELF                = 3,  // Push self onto the stack
    STORE_INSTANCE_VARIABLE  = 4,  // Store a value into an instance variable (followed by 4-byte offset)
    STORE_TEMPORARY_VARIABLE = 5,  // Store a value into a temporary variable (followed by 4-byte offset)
    SEND_MESSAGE             = 6,  // Send a message (followed by 4-byte selector index and 4-byte arg count)
    RETURN_STACK_TOP         = 7,  // Return the value on top of the stack
    JUMP                     = 8,  // Jump to a different bytecode (followed by 4-byte target)
    JUMP_IF_TRUE             = 9,  // Jump if top of stack is true (followed by 4-byte target)
    JUMP_IF_FALSE            = 10, // Jump if top of stack is false (followed by 4-byte target)
    POP                      = 11, // Pop the top value from the stack
    DUPLICATE                = 12, // Duplicate the top value on the stack
    CREATE_BLOCK             = 13, // Create a block (followed by 4-byte bytecode size, 4-byte literal count, 4-byte temp var count)
    EXECUTE_BLOCK            = 14  // Execute a block (followed by 4-byte arg count)
};

// Instruction size (in bytes, including opcode)
inline int getInstructionSize(Bytecode bytecode) {
    switch (bytecode) {
        case Bytecode::PUSH_LITERAL:
        case Bytecode::PUSH_INSTANCE_VARIABLE:
        case Bytecode::PUSH_TEMPORARY_VARIABLE:
        case Bytecode::STORE_INSTANCE_VARIABLE:
        case Bytecode::STORE_TEMPORARY_VARIABLE:
        case Bytecode::JUMP:
        case Bytecode::JUMP_IF_TRUE:
        case Bytecode::JUMP_IF_FALSE:
        case Bytecode::EXECUTE_BLOCK:
            return 5; // 1 byte opcode + 4 byte operand
        case Bytecode::SEND_MESSAGE:
            return 9; // 1 byte opcode + 4 byte selector index + 4 byte arg count
        case Bytecode::CREATE_BLOCK:
            return 13; // 1 byte opcode + 4 byte bytecode size + 4 byte literal count + 4 byte temp var count
        case Bytecode::PUSH_SELF:
        case Bytecode::RETURN_STACK_TOP:
        case Bytecode::POP:
        case Bytecode::DUPLICATE:
            return 1; // 1 byte opcode
        default:
            return 1; // Default to 1 byte for unknown bytecodes
    }
}

// Get bytecode name
inline const char* getBytecodeString(Bytecode bytecode) {
    switch (bytecode) {
        case Bytecode::PUSH_LITERAL:            return "PUSH_LITERAL";
        case Bytecode::PUSH_INSTANCE_VARIABLE:  return "PUSH_INSTANCE_VARIABLE";
        case Bytecode::PUSH_TEMPORARY_VARIABLE: return "PUSH_TEMPORARY_VARIABLE";
        case Bytecode::PUSH_SELF:               return "PUSH_SELF";
        case Bytecode::STORE_INSTANCE_VARIABLE: return "STORE_INSTANCE_VARIABLE";
        case Bytecode::STORE_TEMPORARY_VARIABLE:return "STORE_TEMPORARY_VARIABLE";
        case Bytecode::SEND_MESSAGE:            return "SEND_MESSAGE";
        case Bytecode::RETURN_STACK_TOP:        return "RETURN_STACK_TOP";
        case Bytecode::JUMP:                    return "JUMP";
        case Bytecode::JUMP_IF_TRUE:            return "JUMP_IF_TRUE";
        case Bytecode::JUMP_IF_FALSE:           return "JUMP_IF_FALSE";
        case Bytecode::POP:                     return "POP";
        case Bytecode::DUPLICATE:               return "DUPLICATE";
        case Bytecode::CREATE_BLOCK:            return "CREATE_BLOCK";
        case Bytecode::EXECUTE_BLOCK:           return "EXECUTE_BLOCK";
        default:                                return "UNKNOWN";
    }
}

} // namespace smalltalk