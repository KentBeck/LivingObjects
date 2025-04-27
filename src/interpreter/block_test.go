package main

import (
	"testing"
)

func TestSimpleBlockValue(t *testing.T) {
	vm := NewVM()

	// Create a block that computes 2 + 3
	block := NewBlock(nil)
	block.Class = vm.BlockClass

	// Add literals for the block
	twoObj := vm.NewInteger(2)
	threeObj := vm.NewInteger(3)
	plusSymbol := NewSymbol("+")

	block.Block.Literals = []*Object{twoObj, threeObj, plusSymbol}

	// Create bytecodes for the block: 2 + 3
	// 1. Push 2
	block.Block.Bytecodes = append(block.Block.Bytecodes, PUSH_LITERAL)
	block.Block.Bytecodes = append(block.Block.Bytecodes, 0, 0, 0, 0) // index 0 (2)

	// 2. Push 3
	block.Block.Bytecodes = append(block.Block.Bytecodes, PUSH_LITERAL)
	block.Block.Bytecodes = append(block.Block.Bytecodes, 0, 0, 0, 1) // index 1 (3)

	// 3. Send + message
	block.Block.Bytecodes = append(block.Block.Bytecodes, SEND_MESSAGE)
	block.Block.Bytecodes = append(block.Block.Bytecodes, 0, 0, 0, 2) // index 2 (+)
	block.Block.Bytecodes = append(block.Block.Bytecodes, 0, 0, 0, 1) // 1 argument

	// 4. Return the result
	block.Block.Bytecodes = append(block.Block.Bytecodes, RETURN_STACK_TOP)

	// Create a method that sends 'value' to the block
	builder := NewMethodBuilder(vm.ObjectClass).Selector("testBlockValue")

	// Add the block as a literal
	blockIndex, builder := builder.AddLiteral(block)

	// Add the 'value' selector as a literal
	valueSelector, builder := builder.AddLiteral(NewSymbol("value"))

	// Create bytecodes for the method:
	// 1. Push the block
	builder.PushLiteral(blockIndex)

	// 2. Send 'value' message to the block
	builder.SendMessage(valueSelector, 0)

	// 3. Return the result
	builder.ReturnStackTop()

	// Finalize the method
	method := builder.Go()

	// Create a context for the method
	context := NewContext(method, vm.ObjectClass, []*Object{}, nil)

	// Execute the context
	result, err := vm.ExecuteContext(context)
	if err != nil {
		t.Errorf("Error executing context: %v", err)
	}

	// Check that the result is 5
	if IsIntegerImmediate(result) {
		intValue := GetIntegerImmediate(result)
		if intValue != 5 {
			t.Errorf("Expected result to be 5, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}
}
