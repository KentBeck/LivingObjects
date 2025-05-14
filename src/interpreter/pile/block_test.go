package pile_test

import (
	"testing"

	"smalltalklsp/interpreter/pile"
)

func TestNewBlock(t *testing.T) {
	// Create a context
	context := "test context"
	
	// Create a block
	block := pile.NewBlock(context)
	
	if block.Type() != pile.OBJ_BLOCK {
		t.Errorf("NewBlock(context).Type() = %d, want %d", block.Type(), pile.OBJ_BLOCK)
	}
	
	blockObj := pile.ObjectToBlock(block)
	
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
	block := pile.ObjectToBlock(pile.NewBlock(context))
	
	obj := pile.BlockToObject(block)
	
	if obj.Type() != pile.OBJ_BLOCK {
		t.Errorf("BlockToObject(block).Type() = %d, want %d", obj.Type(), pile.OBJ_BLOCK)
	}
	
	backToBlock := pile.ObjectToBlock(obj)
	if backToBlock.OuterContext != context {
		t.Errorf("ObjectToBlock(BlockToObject(block)).OuterContext = %v, want %v", backToBlock.OuterContext, context)
	}
}

func TestBlockString(t *testing.T) {
	// Create a block
	context := "test context"
	block := pile.ObjectToBlock(pile.NewBlock(context))
	
	if block.String() != "Block" {
		t.Errorf("block.String() = %q, want %q", block.String(), "Block")
	}
}

func TestBlockGetBytecodes(t *testing.T) {
	// Create a block
	context := "test context"
	block := pile.ObjectToBlock(pile.NewBlock(context))
	
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
	block := pile.ObjectToBlock(pile.NewBlock(context))
	
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
	block := pile.ObjectToBlock(pile.NewBlock(context))
	
	// Set literals
	literals := []*pile.Object{
		pile.MakeIntegerImmediate(1),
		pile.MakeIntegerImmediate(2),
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
	block := pile.ObjectToBlock(pile.NewBlock(context))
	
	// Add literals
	lit1 := pile.MakeIntegerImmediate(1)
	lit2 := pile.MakeIntegerImmediate(2)
	
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
	block := pile.ObjectToBlock(pile.NewBlock(context))
	
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
	block := pile.ObjectToBlock(pile.NewBlock(context))
	
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
	block := pile.ObjectToBlock(pile.NewBlock(context))
	
	// Get outer context
	result := block.GetOuterContext()
	
	if result != context {
		t.Errorf("block.GetOuterContext() = %v, want %v", result, context)
	}
}

func TestBlockSetOuterContext(t *testing.T) {
	// Create a block
	context1 := "test context 1"
	block := pile.ObjectToBlock(pile.NewBlock(context1))
	
	// Set a new outer context
	context2 := "test context 2"
	block.SetOuterContext(context2)
	
	if block.OuterContext != context2 {
		t.Errorf("block.OuterContext = %v, want %v", block.OuterContext, context2)
	}
}

// MockBlockExecutor is a mock implementation of the BlockExecutor interface
type MockBlockExecutor struct {
	ReturnValue *pile.Object
}

// ExecuteBlock implements the BlockExecutor interface
func (m *MockBlockExecutor) ExecuteBlock(block *pile.Object, args []*pile.Object) *pile.Object {
	return m.ReturnValue
}

// MockExceptionExecutor is a mock implementation of the BlockExecutor interface that panics with an exception
type MockExceptionExecutor struct {
	ExceptionClass *pile.Object
}

// ExecuteBlock implements the BlockExecutor interface
func (m *MockExceptionExecutor) ExecuteBlock(block *pile.Object, args []*pile.Object) *pile.Object {
	// Create an exception with the specified class
	exception := pile.NewException(m.ExceptionClass)
	
	// Set message text for the exception
	exceptionObj := pile.ObjectToException(exception)
	messageText := pile.StringToObject(pile.NewString("Test exception"))
	exceptionObj.SetMessageText(messageText)
	
	// Panic with the exception to trigger the exception handling mechanism
	panic(exception)
}

func TestBlockValue(t *testing.T) {
	// Create a block
	context := "test context"
	block := pile.ObjectToBlock(pile.NewBlock(context))
	
	// Register a mock block executor
	oldExecutor := pile.GetCurrentBlockExecutor()
	defer func() {
		pile.RegisterBlockExecutor(oldExecutor)
	}()
	
	// Create a mock block executor that returns the value 42
	mockExecutor := &MockBlockExecutor{
		ReturnValue: pile.MakeIntegerImmediate(42),
	}
	pile.RegisterBlockExecutor(mockExecutor)
	
	// Execute the block
	result := block.Value()
	
	// Check that the result is 42
	if !pile.IsIntegerImmediate(result) {
		t.Fatalf("Expected an integer, got %v", result)
	}
	
	value := pile.GetIntegerImmediate(result)
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}
}

func TestBlockValueWithArguments(t *testing.T) {
	// Create a block
	context := "test context"
	block := pile.ObjectToBlock(pile.NewBlock(context))
	
	// Register a mock block executor
	oldExecutor := pile.GetCurrentBlockExecutor()
	defer func() {
		pile.RegisterBlockExecutor(oldExecutor)
	}()
	
	// Create a mock block executor that returns the value 42
	mockExecutor := &MockBlockExecutor{
		ReturnValue: pile.MakeIntegerImmediate(42),
	}
	pile.RegisterBlockExecutor(mockExecutor)
	
	// Execute the block with arguments
	args := []*pile.Object{
		pile.MakeIntegerImmediate(1),
		pile.MakeIntegerImmediate(2),
	}
	result := block.ValueWithArguments(args)
	
	// Check that the result is 42
	if !pile.IsIntegerImmediate(result) {
		t.Fatalf("Expected an integer, got %v", result)
	}
	
	value := pile.GetIntegerImmediate(result)
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}
}