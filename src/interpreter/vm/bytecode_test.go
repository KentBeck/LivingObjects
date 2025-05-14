package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/bytecode"
)

func TestBytecodeName(t *testing.T) {
	tests := []struct {
		bytecode byte
		expected string
	}{
		{bytecode.PUSH_LITERAL, "PUSH_LITERAL"},
		{bytecode.PUSH_INSTANCE_VARIABLE, "PUSH_INSTANCE_VARIABLE"},
		{bytecode.PUSH_TEMPORARY_VARIABLE, "PUSH_TEMPORARY_VARIABLE"},
		{bytecode.PUSH_SELF, "PUSH_SELF"},
		{bytecode.STORE_INSTANCE_VARIABLE, "STORE_INSTANCE_VARIABLE"},
		{bytecode.STORE_TEMPORARY_VARIABLE, "STORE_TEMPORARY_VARIABLE"},
		{bytecode.SEND_MESSAGE, "SEND_MESSAGE"},
		{bytecode.RETURN_STACK_TOP, "RETURN_STACK_TOP"},
		{bytecode.JUMP, "JUMP"},
		{bytecode.JUMP_IF_TRUE, "JUMP_IF_TRUE"},
		{bytecode.JUMP_IF_FALSE, "JUMP_IF_FALSE"},
		{bytecode.POP, "POP"},
		{bytecode.DUPLICATE, "DUPLICATE"},
		{255, "UNKNOWN"}, // Test unknown bytecode
	}

	for _, test := range tests {
		result := bytecode.BytecodeName(test.bytecode)
		if result != test.expected {
			t.Errorf("BytecodeName(%d) = %s, expected %s", test.bytecode, result, test.expected)
		}
	}
}

func TestInstructionSize(t *testing.T) {
	tests := []struct {
		bytecode byte
		expected int
	}{
		{bytecode.PUSH_LITERAL, 5},
		{bytecode.PUSH_INSTANCE_VARIABLE, 5},
		{bytecode.PUSH_TEMPORARY_VARIABLE, 5},
		{bytecode.PUSH_SELF, 1},
		{bytecode.STORE_INSTANCE_VARIABLE, 5},
		{bytecode.STORE_TEMPORARY_VARIABLE, 5},
		{bytecode.SEND_MESSAGE, 9},
		{bytecode.RETURN_STACK_TOP, 1},
		{bytecode.JUMP, 5},
		{bytecode.JUMP_IF_TRUE, 5},
		{bytecode.JUMP_IF_FALSE, 5},
		{bytecode.POP, 1},
		{bytecode.DUPLICATE, 1},
		{255, 1}, // Test unknown bytecode
	}

	for _, test := range tests {
		result := bytecode.InstructionSize(test.bytecode)
		if result != test.expected {
			t.Errorf("InstructionSize(%d) = %d, expected %d", test.bytecode, result, test.expected)
		}
	}
}