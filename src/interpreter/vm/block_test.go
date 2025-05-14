package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

func TestBasicBlock(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Create a block with proper class field
	block := virtualMachine.NewBlock(nil)

	// Check that the block is of the correct class
	blockClass := virtualMachine.GetClass(block)
	if blockClass != virtualMachine.Classes.Get(vm.Block) {
		t.Errorf("Expected block class to be Block class, got %v", blockClass)
	}

	// Check that the block has the correct type
	if block.Type() != pile.OBJ_BLOCK {
		t.Errorf("Expected block type to be OBJ_BLOCK, got %v", block.Type())
	}

	// Check that the block has the correct string representation
	if block.String() != "Block" {
		t.Errorf("Expected block string to be 'Block', got %s", block.String())
	}
}
