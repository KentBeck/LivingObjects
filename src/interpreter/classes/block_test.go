package classes

import (
	"testing"

	"smalltalklsp/interpreter/core"
)

func TestNewBlock(t *testing.T) {
	// Create a context
	context := "test context"
	
	// Create a block
	block := NewBlock(context)
	
	if block.Type() != core.OBJ_BLOCK {
		t.Errorf("NewBlock(context).Type() = %d, want %d", block.Type(), core.OBJ_BLOCK)
	}
	
	blockObj := ObjectToBlock(block)
	
	if blockObj.OuterContext != context {
		t.Errorf("blockObj.OuterContext = %v, want %v", blockObj.OuterContext, context)
	}
	
	if len(blockObj.Bytecodes) != 0 {
		t.Errorf("len(blockObj.Bytecodes) = %d, want 0", len(blockObj.Bytecodes))
	}
	
	if len(blockObj.Literals) != 0 {
		t.Errorf("len(blockObj.Literals) = %d, want 0", len(blockObj.Literals))
	}
	
	if len(blockObj.TempVarNames) != 0 {
		t.Errorf("len(blockObj.TempVarNames) = %d, want 0", len(blockObj.TempVarNames))
	}
}

func TestBlockToObjectAndBack(t *testing.T) {
	// Create a block
	context := "test context"
	block := ObjectToBlock(NewBlock(context))
	
	obj := BlockToObject(block)
	
	if obj.Type() != core.OBJ_BLOCK {
		t.Errorf("BlockToObject(block).Type() = %d, want %d", obj.Type(), core.OBJ_BLOCK)
	}
	
	backToBlock := ObjectToBlock(obj)
	if backToBlock.OuterContext != context {
		t.Errorf("ObjectToBlock(BlockToObject(block)).OuterContext = %v, want %v", backToBlock.OuterContext, context)
	}
}

func TestBlockString(t *testing.T) {
	// Create a block
	context := "test context"
	block := ObjectToBlock(NewBlock(context))
	
	if block.String() != "Block" {
		t.Errorf("block.String() = %q, want %q", block.String(), "Block")
	}
}

func TestBlockGetBytecodes(t *testing.T) {
	// Create a block
	context := "test context"
	block := ObjectToBlock(NewBlock(context))
	
	// Set bytecodes
	bytecodes := []byte{1, 2, 3}
	block.Bytecodes = bytecodes
	
	// Get bytecodes
	result := block.GetBytecodes()
	
	if len(result) != len(bytecodes) {
		t.Errorf("len(block.GetBytecodes()) = %d, want %d", len(result), len(bytecodes))
	}
	
	for i, b := range result {
		if b != bytecodes[i] {
			t.Errorf("block.GetBytecodes()[%d] = %d, want %d", i, b, bytecodes[i])
		}
	}
}

func TestBlockSetBytecodes(t *testing.T) {
	// Create a block
	context := "test context"
	block := ObjectToBlock(NewBlock(context))
	
	// Set bytecodes
	bytecodes := []byte{1, 2, 3}
	block.SetBytecodes(bytecodes)
	
	if len(block.Bytecodes) != len(bytecodes) {
		t.Errorf("len(block.Bytecodes) = %d, want %d", len(block.Bytecodes), len(bytecodes))
	}
	
	for i, b := range block.Bytecodes {
		if b != bytecodes[i] {
			t.Errorf("block.Bytecodes[%d] = %d, want %d", i, b, bytecodes[i])
		}
	}
}

func TestBlockGetLiterals(t *testing.T) {
	// Create a block
	context := "test context"
	block := ObjectToBlock(NewBlock(context))
	
	// Set literals
	literals := []*core.Object{
		core.MakeIntegerImmediate(1),
		core.MakeIntegerImmediate(2),
	}
	block.Literals = literals
	
	// Get literals
	result := block.GetLiterals()
	
	if len(result) != len(literals) {
		t.Errorf("len(block.GetLiterals()) = %d, want %d", len(result), len(literals))
	}
	
	for i, lit := range result {
		if lit != literals[i] {
			t.Errorf("block.GetLiterals()[%d] = %v, want %v", i, lit, literals[i])
		}
	}
}

func TestBlockAddLiteral(t *testing.T) {
	// Create a block
	context := "test context"
	block := ObjectToBlock(NewBlock(context))
	
	// Add literals
	lit1 := core.MakeIntegerImmediate(1)
	lit2 := core.MakeIntegerImmediate(2)
	
	block.AddLiteral(lit1)
	block.AddLiteral(lit2)
	
	if len(block.Literals) != 2 {
		t.Errorf("len(block.Literals) = %d, want 2", len(block.Literals))
	}
	
	if block.Literals[0] != lit1 {
		t.Errorf("block.Literals[0] = %v, want %v", block.Literals[0], lit1)
	}
	
	if block.Literals[1] != lit2 {
		t.Errorf("block.Literals[1] = %v, want %v", block.Literals[1], lit2)
	}
}

func TestBlockGetTempVarNames(t *testing.T) {
	// Create a block
	context := "test context"
	block := ObjectToBlock(NewBlock(context))
	
	// Set temp var names
	tempVars := []string{"temp1", "temp2"}
	block.TempVarNames = tempVars
	
	// Get temp var names
	result := block.GetTempVarNames()
	
	if len(result) != len(tempVars) {
		t.Errorf("len(block.GetTempVarNames()) = %d, want %d", len(result), len(tempVars))
	}
	
	for i, name := range result {
		if name != tempVars[i] {
			t.Errorf("block.GetTempVarNames()[%d] = %q, want %q", i, name, tempVars[i])
		}
	}
}

func TestBlockAddTempVarName(t *testing.T) {
	// Create a block
	context := "test context"
	block := ObjectToBlock(NewBlock(context))
	
	// Add temp var names
	block.AddTempVarName("temp1")
	block.AddTempVarName("temp2")
	
	if len(block.TempVarNames) != 2 {
		t.Errorf("len(block.TempVarNames) = %d, want 2", len(block.TempVarNames))
	}
	
	if block.TempVarNames[0] != "temp1" {
		t.Errorf("block.TempVarNames[0] = %q, want %q", block.TempVarNames[0], "temp1")
	}
	
	if block.TempVarNames[1] != "temp2" {
		t.Errorf("block.TempVarNames[1] = %q, want %q", block.TempVarNames[1], "temp2")
	}
}

func TestBlockGetOuterContext(t *testing.T) {
	// Create a block
	context := "test context"
	block := ObjectToBlock(NewBlock(context))
	
	// Get outer context
	result := block.GetOuterContext()
	
	if result != context {
		t.Errorf("block.GetOuterContext() = %v, want %v", result, context)
	}
}

func TestBlockSetOuterContext(t *testing.T) {
	// Create a block
	context1 := "test context 1"
	block := ObjectToBlock(NewBlock(context1))
	
	// Set a new outer context
	context2 := "test context 2"
	block.SetOuterContext(context2)
	
	if block.OuterContext != context2 {
		t.Errorf("block.OuterContext = %v, want %v", block.OuterContext, context2)
	}
}

func TestBlockValue(t *testing.T) {
	// Create a block
	context := "test context"
	block := ObjectToBlock(NewBlock(context))
	
	// Call Value
	result := block.Value()
	
	// Currently, Value just returns nil
	if !core.IsNilImmediate(result) {
		t.Errorf("block.Value() = %v, want nil", result)
	}
}

func TestBlockValueWithArguments(t *testing.T) {
	// Create a block
	context := "test context"
	block := ObjectToBlock(NewBlock(context))
	
	// Call ValueWithArguments
	args := []*core.Object{
		core.MakeIntegerImmediate(1),
		core.MakeIntegerImmediate(2),
	}
	result := block.ValueWithArguments(args)
	
	// Currently, ValueWithArguments just returns nil
	if !core.IsNilImmediate(result) {
		t.Errorf("block.ValueWithArguments(args) = %v, want nil", result)
	}
}
