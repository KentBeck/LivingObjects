package main

import (
	"testing"
)

func TestBytecodeName(t *testing.T) {
	tests := []struct {
		bytecode byte
		expected string
	}{
		{PUSH_LITERAL, "PUSH_LITERAL"},
		{PUSH_INSTANCE_VARIABLE, "PUSH_INSTANCE_VARIABLE"},
		{PUSH_TEMPORARY_VARIABLE, "PUSH_TEMPORARY_VARIABLE"},
		{PUSH_SELF, "PUSH_SELF"},
		{STORE_INSTANCE_VARIABLE, "STORE_INSTANCE_VARIABLE"},
		{STORE_TEMPORARY_VARIABLE, "STORE_TEMPORARY_VARIABLE"},
		{SEND_MESSAGE, "SEND_MESSAGE"},
		{RETURN_STACK_TOP, "RETURN_STACK_TOP"},
		{JUMP, "JUMP"},
		{JUMP_IF_TRUE, "JUMP_IF_TRUE"},
		{JUMP_IF_FALSE, "JUMP_IF_FALSE"},
		{POP, "POP"},
		{DUPLICATE, "DUPLICATE"},
		{255, "UNKNOWN"}, // Test unknown bytecode
	}

	for _, test := range tests {
		result := BytecodeName(test.bytecode)
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
		{PUSH_LITERAL, 5},
		{PUSH_INSTANCE_VARIABLE, 5},
		{PUSH_TEMPORARY_VARIABLE, 5},
		{PUSH_SELF, 1},
		{STORE_INSTANCE_VARIABLE, 5},
		{STORE_TEMPORARY_VARIABLE, 5},
		{SEND_MESSAGE, 9},
		{RETURN_STACK_TOP, 1},
		{JUMP, 5},
		{JUMP_IF_TRUE, 5},
		{JUMP_IF_FALSE, 5},
		{POP, 1},
		{DUPLICATE, 1},
		{255, 1}, // Test unknown bytecode
	}

	for _, test := range tests {
		result := InstructionSize(test.bytecode)
		if result != test.expected {
			t.Errorf("InstructionSize(%d) = %d, expected %d", test.bytecode, result, test.expected)
		}
	}
}
