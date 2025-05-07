package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/vm"
)

func TestBytecodeName(t *testing.T) {
	tests := []struct {
		bytecode byte
		expected string
	}{
		{vm.PUSH_LITERAL, "PUSH_LITERAL"},
		{vm.PUSH_INSTANCE_VARIABLE, "PUSH_INSTANCE_VARIABLE"},
		{vm.PUSH_TEMPORARY_VARIABLE, "PUSH_TEMPORARY_VARIABLE"},
		{vm.PUSH_SELF, "PUSH_SELF"},
		{vm.STORE_INSTANCE_VARIABLE, "STORE_INSTANCE_VARIABLE"},
		{vm.STORE_TEMPORARY_VARIABLE, "STORE_TEMPORARY_VARIABLE"},
		{vm.SEND_MESSAGE, "SEND_MESSAGE"},
		{vm.RETURN_STACK_TOP, "RETURN_STACK_TOP"},
		{vm.JUMP, "JUMP"},
		{vm.JUMP_IF_TRUE, "JUMP_IF_TRUE"},
		{vm.JUMP_IF_FALSE, "JUMP_IF_FALSE"},
		{vm.POP, "POP"},
		{vm.DUPLICATE, "DUPLICATE"},
		{255, "UNKNOWN"}, // Test unknown bytecode
	}

	for _, test := range tests {
		result := vm.BytecodeName(test.bytecode)
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
		{vm.PUSH_LITERAL, 5},
		{vm.PUSH_INSTANCE_VARIABLE, 5},
		{vm.PUSH_TEMPORARY_VARIABLE, 5},
		{vm.PUSH_SELF, 1},
		{vm.STORE_INSTANCE_VARIABLE, 5},
		{vm.STORE_TEMPORARY_VARIABLE, 5},
		{vm.SEND_MESSAGE, 9},
		{vm.RETURN_STACK_TOP, 1},
		{vm.JUMP, 5},
		{vm.JUMP_IF_TRUE, 5},
		{vm.JUMP_IF_FALSE, 5},
		{vm.POP, 1},
		{vm.DUPLICATE, 1},
		{255, 1}, // Test unknown bytecode
	}

	for _, test := range tests {
		result := vm.InstructionSize(test.bytecode)
		if result != test.expected {
			t.Errorf("InstructionSize(%d) = %d, expected %d", test.bytecode, result, test.expected)
		}
	}
}
