package main

import (
	"testing"
)

func TestSimpleBlockValue(t *testing.T) {
	vm := NewVM()

	// Create a block that computes 2 + 3
	block := NewBlock(nil)
	block.SetClass(ClassToObject(vm.BlockClass))

	// Add literals for the block
	twoObj := vm.NewInteger(2)
	threeObj := vm.NewInteger(3)
	plusSymbol := NewSymbol("+")

	// Get the block struct
	blockObj := ObjectToBlock(block)
	blockObj.Literals = []*Object{twoObj, threeObj, plusSymbol}

	// Create bytecodes for the block: 2 + 3
	// 1. Push 2
	blockObj.Bytecodes = append(blockObj.Bytecodes, PUSH_LITERAL)
	blockObj.Bytecodes = append(blockObj.Bytecodes, 0, 0, 0, 0) // index 0 (2)

	// 2. Push 3
	blockObj.Bytecodes = append(blockObj.Bytecodes, PUSH_LITERAL)
	blockObj.Bytecodes = append(blockObj.Bytecodes, 0, 0, 0, 1) // index 1 (3)

	// 3. Send + message
	blockObj.Bytecodes = append(blockObj.Bytecodes, SEND_MESSAGE)
	blockObj.Bytecodes = append(blockObj.Bytecodes, 0, 0, 0, 2) // index 2 (+)
	blockObj.Bytecodes = append(blockObj.Bytecodes, 0, 0, 0, 1) // 1 argument

	// 4. Return the result
	blockObj.Bytecodes = append(blockObj.Bytecodes, RETURN_STACK_TOP)

	// Create a method that sends 'value' to the block
	builder := NewMethodBuilder(vm.ObjectClass).Selector("testBlockValue")

	// Add the block as a literal
	blockIndex, builder := builder.AddLiteral(block)

	// Add the 'value' selector as a literal
	valueSelector, builder := builder.AddLiteral(NewSymbol("value"))

	method := builder.
		PushLiteral(blockIndex).
		SendMessage(valueSelector, 0).
		ReturnStackTop().
		Go()

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
