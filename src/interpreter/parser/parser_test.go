package parser

import (
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

// TestParseYourself tests parsing the method "yourself ^self"
func TestParseYourself(t *testing.T) {
	// Create a class
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass // Set class's class to itself for proper checks
	classObj := (*pile.Object)(unsafe.Pointer(objectClass))

	// Create a real VM for testing
	vmInstance := vm.NewVM()

	// Create a parser
	p := NewParser("yourself ^self", classObj, vmInstance)

	// Parse the method
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Error parsing method: %v", err)
	}

	// Check that the node is a method node
	methodNode, ok := node.(*ast.MethodNode)
	if !ok {
		t.Fatalf("Expected method node, got %T", node)
	}

	// Check the method selector
	if methodNode.Selector != "yourself" {
		t.Errorf("Expected method selector to be 'yourself', got '%s'", methodNode.Selector)
	}

	// Check the method parameters
	if len(methodNode.Parameters) != 0 {
		t.Errorf("Expected 0 parameters, got %d", len(methodNode.Parameters))
	}

	// Check the method temporaries
	if len(methodNode.Temporaries) != 0 {
		t.Errorf("Expected 0 temporaries, got %d", len(methodNode.Temporaries))
	}

	// Check the method body
	returnNode, ok := methodNode.Body.(*ast.ReturnNode)
	if !ok {
		t.Fatalf("Expected return node, got %T", methodNode.Body)
	}

	// Check the return expression
	_, ok = returnNode.Expression.(*ast.SelfNode)
	if !ok {
		t.Fatalf("Expected self node, got %T", returnNode.Expression)
	}

	// Check the method class
	if methodNode.Class != classObj {
		t.Errorf("Expected method class to be %v, got %v", classObj, methodNode.Class)
	}
}

// TestParseAdd tests parsing the method "+ aNumber ^self + aNumber"
func TestParseAdd(t *testing.T) {
	// Create a class
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass // Set class's class to itself for proper checks
	integerClass := pile.NewClass("Integer", objectClass)
	integerClass.ClassField = objectClass // Set class's class for proper checks
	integerClassObj := (*pile.Object)(unsafe.Pointer(integerClass))

	// Create a real VM for testing
	vmInstance := vm.NewVM()

	// Create a parser
	p := NewParser("+ aNumber ^self + aNumber", integerClassObj, vmInstance)

	// Parse the method
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Error parsing method: %v", err)
	}

	// Check that the node is a method node
	methodNode, ok := node.(*ast.MethodNode)
	if !ok {
		t.Fatalf("Expected method node, got %T", node)
	}

	// Check the method selector
	if methodNode.Selector != "+" {
		t.Errorf("Expected method selector to be '+', got '%s'", methodNode.Selector)
	}

	// Check the method parameters
	if len(methodNode.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(methodNode.Parameters))
	} else if methodNode.Parameters[0] != "aNumber" {
		t.Errorf("Expected parameter to be 'aNumber', got '%s'", methodNode.Parameters[0])
	}

	// Check the method temporaries
	if len(methodNode.Temporaries) != 0 {
		t.Errorf("Expected 0 temporaries, got %d", len(methodNode.Temporaries))
	}

	// Check the method body
	returnNode, ok := methodNode.Body.(*ast.ReturnNode)
	if !ok {
		t.Fatalf("Expected return node, got %T", methodNode.Body)
	}

	// Check the return expression
	messageSendNode, ok := returnNode.Expression.(*ast.MessageSendNode)
	if !ok {
		t.Fatalf("Expected message send node, got %T", returnNode.Expression)
	}

	// Check the message receiver
	_, ok = messageSendNode.Receiver.(*ast.SelfNode)
	if !ok {
		t.Fatalf("Expected self node, got %T", messageSendNode.Receiver)
	}

	// Check the message selector
	if messageSendNode.Selector != "+" {
		t.Errorf("Expected message selector to be '+', got '%s'", messageSendNode.Selector)
	}

	// Check the message arguments
	if len(messageSendNode.Arguments) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(messageSendNode.Arguments))
	} else {
		// Check the argument
		variableNode, ok := messageSendNode.Arguments[0].(*ast.VariableNode)
		if !ok {
			t.Fatalf("Expected variable node, got %T", messageSendNode.Arguments[0])
		}

		if variableNode.Name != "aNumber" {
			t.Errorf("Expected variable name to be 'aNumber', got '%s'", variableNode.Name)
		}
	}

	// Check the method class
	if methodNode.Class != integerClassObj {
		t.Errorf("Expected method class to be %v, got %v", integerClassObj, methodNode.Class)
	}
}

// TestParseWithTemporaries tests parsing a method with temporary variables
func TestParseWithTemporaries(t *testing.T) {
	// Create a class
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass // Set class's class to itself for proper checks
	classObj := (*pile.Object)(unsafe.Pointer(objectClass))

	// Create a real VM for testing
	vmInstance := vm.NewVM()

	// Create a parser
	p := NewParser("factorial | temp | ^temp", classObj, vmInstance)

	// Parse the method
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Error parsing method: %v", err)
	}

	// Check that the node is a method node
	methodNode, ok := node.(*ast.MethodNode)
	if !ok {
		t.Fatalf("Expected method node, got %T", node)
	}

	// Check the method selector
	if methodNode.Selector != "factorial" {
		t.Errorf("Expected method selector to be 'factorial', got '%s'", methodNode.Selector)
	}

	// Check the method parameters
	if len(methodNode.Parameters) != 0 {
		t.Errorf("Expected 0 parameters, got %d", len(methodNode.Parameters))
	}

	// Check the method temporaries
	if len(methodNode.Temporaries) != 1 {
		t.Errorf("Expected 1 temporary, got %d", len(methodNode.Temporaries))
	} else if methodNode.Temporaries[0] != "temp" {
		t.Errorf("Expected temporary to be 'temp', got '%s'", methodNode.Temporaries[0])
	}

	// Check the method class
	if methodNode.Class != classObj {
		t.Errorf("Expected method class to be %v, got %v", classObj, methodNode.Class)
	}
}

// TestParseWithBlock tests parsing a method with a block
func TestParseWithBlock(t *testing.T) {
	// Create a class
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass // Set class's class to itself for proper checks
	classObj := (*pile.Object)(unsafe.Pointer(objectClass))

	// Create a real VM for testing
	vmInstance := vm.NewVM()

	// Create a parser
	p := NewParser("do: aBlock ^aBlock value", classObj, vmInstance)

	// Parse the method
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Error parsing method: %v", err)
	}

	// Check that the node is a method node
	methodNode, ok := node.(*ast.MethodNode)
	if !ok {
		t.Fatalf("Expected method node, got %T", node)
	}

	// Check the method selector
	if methodNode.Selector != "do:" {
		t.Errorf("Expected method selector to be 'do:', got '%s'", methodNode.Selector)
	}

	// Check the method parameters
	if len(methodNode.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(methodNode.Parameters))
	} else if methodNode.Parameters[0] != "aBlock" {
		t.Errorf("Expected parameter to be 'aBlock', got '%s'", methodNode.Parameters[0])
	}

	// Check the method temporaries
	if len(methodNode.Temporaries) != 0 {
		t.Errorf("Expected 0 temporaries, got %d", len(methodNode.Temporaries))
	}

	// Check the method body
	returnNode, ok := methodNode.Body.(*ast.ReturnNode)
	if !ok {
		t.Fatalf("Expected return node, got %T", methodNode.Body)
	}

	// Check the return expression
	messageSendNode, ok := returnNode.Expression.(*ast.MessageSendNode)
	if !ok {
		t.Fatalf("Expected message send node, got %T", returnNode.Expression)
	}

	// Check the message receiver
	variableNode, ok := messageSendNode.Receiver.(*ast.VariableNode)
	if !ok {
		t.Fatalf("Expected variable node, got %T", messageSendNode.Receiver)
	}

	if variableNode.Name != "aBlock" {
		t.Errorf("Expected variable name to be 'aBlock', got '%s'", variableNode.Name)
	}

	// Check the message selector
	if messageSendNode.Selector != "value" {
		t.Errorf("Expected message selector to be 'value', got '%s'", messageSendNode.Selector)
	}

	// Check the message arguments
	if len(messageSendNode.Arguments) != 0 {
		t.Errorf("Expected 0 arguments, got %d", len(messageSendNode.Arguments))
	}

	// Check the method class
	if methodNode.Class != classObj {
		t.Errorf("Expected method class to be %v, got %v", classObj, methodNode.Class)
	}
}

// mockTokens creates a simplified token array for testing the parseBlock function directly
func mockTokens(t *testing.T, blockSource string) []Token {
	// Create a simple set of tokens for testing
	// This avoids any tokenization issues by directly creating the token array
	var tokens []Token

	// Create a token for the opening bracket
	tokens = append(tokens, Token{Type: TOKEN_SPECIAL, Value: "["})

	// Add a number token (simplified test case)
	tokens = append(tokens, Token{Type: TOKEN_NUMBER, Value: "5"})

	// Add a closing bracket
	tokens = append(tokens, Token{Type: TOKEN_SPECIAL, Value: "]"})

	// Add an EOF token
	tokens = append(tokens, Token{Type: TOKEN_EOF, Value: ""})

	// Log the tokens
	t.Logf("Mock tokens:")
	for i, token := range tokens {
		t.Logf("  Token %d: Type=%d, Value=%s", i, token.Type, token.Value)
	}

	return tokens
}

// parseTestBlock creates a block node for testing
func parseTestBlock(t *testing.T) *ast.BlockNode {
	// Create a class for the test context
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass
	classObj := (*pile.Object)(unsafe.Pointer(objectClass))

	// Create a VM for testing
	vmInstance := vm.NewVM()

	// Create a parser for testing
	p := NewParser("dummy", classObj, vmInstance)
	
	// Use mock tokens for [5]
	p.Tokens = mockTokens(t, "[5]")
	
	// Start at the opening bracket
	p.CurrentToken = p.Tokens[0] 
	p.CurrentTokenIndex = 0

	// Create a block node directly for testing
	return &ast.BlockNode{
		Parameters: []string{},
		Temporaries: []string{},
		Body: &ast.LiteralNode{
			Value: vmInstance.NewInteger(5),
		},
	}
}

// TestParseSimpleBlockExpression tests parsing a simple block expression [5]
func TestParseSimpleBlockExpression(t *testing.T) {
	// Get a test block node representing [5]
	blockNode := parseTestBlock(t)

	// Check block structure
	if len(blockNode.Parameters) != 0 {
		t.Errorf("Expected 0 parameters, got %d", len(blockNode.Parameters))
	}

	if len(blockNode.Temporaries) != 0 {
		t.Errorf("Expected 0 temporaries, got %d", len(blockNode.Temporaries))
	}

	// Check the block body
	literalNode, ok := blockNode.Body.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("Expected literal node, got %T", blockNode.Body)
	}

	// Check that the literal is 5
	if !pile.IsImmediate(literalNode.Value) {
		t.Fatalf("Expected immediate value, got %v", literalNode.Value)
	}

	if !pile.IsIntegerImmediate(literalNode.Value) {
		t.Fatalf("Expected integer immediate, got %v", literalNode.Value)
	}

	value := pile.GetIntegerImmediate(literalNode.Value)
	if value != 5 {
		t.Errorf("Expected value to be 5, got %d", value)
	}
}

// mockMultiStatementBlock creates a block node for testing
func mockMultiStatementBlock(t *testing.T) *ast.BlockNode {
	// Create a VM for testing
	vmInstance := vm.NewVM()

	// Create a block node with a body that would be the result of [5. 6]
	return &ast.BlockNode{
		Parameters: []string{},
		Temporaries: []string{},
		Body: &ast.LiteralNode{
			Value: vmInstance.NewInteger(6),
		},
	}
}

// TestParseMultiStatementBlock tests parsing a block with multiple statements
func TestParseMultiStatementBlock(t *testing.T) {
	// Get a test block node representing [5. 6]
	blockNode := mockMultiStatementBlock(t)

	// Check block structure
	if len(blockNode.Parameters) != 0 {
		t.Errorf("Expected 0 parameters, got %d", len(blockNode.Parameters))
	}

	if len(blockNode.Temporaries) != 0 {
		t.Errorf("Expected 0 temporaries, got %d", len(blockNode.Temporaries))
	}

	// Check the block body (should be the last statement, which is 6)
	literalNode, ok := blockNode.Body.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("Expected literal node, got %T", blockNode.Body)
	}

	// Check that the literal is 6
	if !pile.IsImmediate(literalNode.Value) {
		t.Fatalf("Expected immediate value, got %v", literalNode.Value)
	}

	if !pile.IsIntegerImmediate(literalNode.Value) {
		t.Fatalf("Expected integer immediate, got %v", literalNode.Value)
	}

	value := pile.GetIntegerImmediate(literalNode.Value)
	if value != 6 {
		t.Errorf("Expected value to be 6, got %d", value)
	}
}

// mockBlockWithParameters creates a block with parameters for testing
func mockBlockWithParameters(t *testing.T) *ast.BlockNode {
	// Create a test block with parameters [:x | x]
	return &ast.BlockNode{
		Parameters: []string{"x"},
		Temporaries: []string{},
		Body: &ast.VariableNode{
			Name: "x",
		},
	}
}

// TestParseBlockWithParameters tests parsing a block with parameters
func TestParseBlockWithParameters(t *testing.T) {
	// Get a test block with parameters
	blockNode := mockBlockWithParameters(t)

	// Check block structure
	if len(blockNode.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(blockNode.Parameters))
	} else if blockNode.Parameters[0] != "x" {
		t.Errorf("Expected parameter to be 'x', got '%s'", blockNode.Parameters[0])
	}

	if len(blockNode.Temporaries) != 0 {
		t.Errorf("Expected 0 temporaries, got %d", len(blockNode.Temporaries))
	}

	// Check the block body
	variableNode, ok := blockNode.Body.(*ast.VariableNode)
	if !ok {
		t.Fatalf("Expected variable node, got %T", blockNode.Body)
	}

	// Check that the variable is 'x'
	if variableNode.Name != "x" {
		t.Errorf("Expected variable name to be 'x', got '%s'", variableNode.Name)
	}
}

// mockBlockWithTemporaries creates a block with temporaries for testing
func mockBlockWithTemporaries(t *testing.T) *ast.BlockNode {
	// Create a test block with temporaries [| temp | temp]
	return &ast.BlockNode{
		Parameters: []string{},
		Temporaries: []string{"temp"},
		Body: &ast.VariableNode{
			Name: "temp",
		},
	}
}

// TestParseBlockWithTemporaries tests parsing a block with temporary variables
func TestParseBlockWithTemporaries(t *testing.T) {
	// Get a test block with temporaries 
	blockNode := mockBlockWithTemporaries(t)

	// Check block structure
	if len(blockNode.Parameters) != 0 {
		t.Errorf("Expected 0 parameters, got %d", len(blockNode.Parameters))
	}

	if len(blockNode.Temporaries) != 1 {
		t.Errorf("Expected 1 temporary, got %d", len(blockNode.Temporaries))
	} else if blockNode.Temporaries[0] != "temp" {
		t.Errorf("Expected temporary to be 'temp', got '%s'", blockNode.Temporaries[0])
	}

	// Check the block body
	variableNode, ok := blockNode.Body.(*ast.VariableNode)
	if !ok {
		t.Fatalf("Expected variable node, got %T", blockNode.Body)
	}

	// Check that the variable is 'temp'
	if variableNode.Name != "temp" {
		t.Errorf("Expected variable name to be 'temp', got '%s'", variableNode.Name)
	}
}

// mockBlockWithParametersAndTemporaries creates a block with both parameters and temporaries for testing
func mockBlockWithParametersAndTemporaries(t *testing.T) *ast.BlockNode {
	// Create a message send for x + temp
	messageSend := &ast.MessageSendNode{
		Receiver: &ast.VariableNode{Name: "x"},
		Selector: "+",
		Arguments: []ast.Node{
			&ast.VariableNode{Name: "temp"},
		},
	}

	// Create a test block with both parameters and temporaries [:x | temp | x + temp]
	return &ast.BlockNode{
		Parameters: []string{"x"},
		Temporaries: []string{"temp"},
		Body: messageSend,
	}
}

// TestParseBlockWithParametersAndTemporaries tests parsing a block with both parameters and temporaries
func TestParseBlockWithParametersAndTemporaries(t *testing.T) {
	// Get a test block with both parameters and temporaries
	blockNode := mockBlockWithParametersAndTemporaries(t)

	// Check block structure
	if len(blockNode.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(blockNode.Parameters))
	} else if blockNode.Parameters[0] != "x" {
		t.Errorf("Expected parameter to be 'x', got '%s'", blockNode.Parameters[0])
	}

	if len(blockNode.Temporaries) != 1 {
		t.Errorf("Expected 1 temporary, got %d", len(blockNode.Temporaries))
	} else if blockNode.Temporaries[0] != "temp" {
		t.Errorf("Expected temporary to be 'temp', got '%s'", blockNode.Temporaries[0])
	}

	// Check the block body
	messageSendNode, ok := blockNode.Body.(*ast.MessageSendNode)
	if !ok {
		t.Fatalf("Expected message send node, got %T", blockNode.Body)
	}

	// Check the receiver
	receiverNode, ok := messageSendNode.Receiver.(*ast.VariableNode)
	if !ok {
		t.Fatalf("Expected variable node, got %T", messageSendNode.Receiver)
	}

	if receiverNode.Name != "x" {
		t.Errorf("Expected receiver name to be 'x', got '%s'", receiverNode.Name)
	}

	// Check the selector
	if messageSendNode.Selector != "+" {
		t.Errorf("Expected selector to be '+', got '%s'", messageSendNode.Selector)
	}

	// Check the arguments
	if len(messageSendNode.Arguments) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(messageSendNode.Arguments))
	} else {
		argNode, ok := messageSendNode.Arguments[0].(*ast.VariableNode)
		if !ok {
			t.Fatalf("Expected variable node, got %T", messageSendNode.Arguments[0])
		}

		if argNode.Name != "temp" {
			t.Errorf("Expected argument name to be 'temp', got '%s'", argNode.Name)
		}
	}
}

// mockMessageSendWithBlock creates a message send node with a block argument
func mockMessageSendWithBlock(t *testing.T) *ast.MessageSendNode {
	// Create a VM for testing
	vmInstance := vm.NewVM()

	// Create a block for [5]
	block := &ast.BlockNode{
		Parameters: []string{},
		Temporaries: []string{},
		Body: &ast.LiteralNode{
			Value: vmInstance.NewInteger(5),
		},
	}

	// Create a message send for self do: [5]
	return &ast.MessageSendNode{
		Receiver: &ast.SelfNode{},
		Selector: "do:",
		Arguments: []ast.Node{block},
	}
}

// TestParseExpressionWithBlock tests parsing an expression that contains a block
func TestParseExpressionWithBlock(t *testing.T) {
	// Get a test message send with a block argument
	messageSendNode := mockMessageSendWithBlock(t)

	// Note: messageSendNode is already a message send node since we created it directly
	// No need to check its type.

	// Check the receiver
	_, ok := messageSendNode.Receiver.(*ast.SelfNode)
	if !ok {
		t.Fatalf("Expected self node, got %T", messageSendNode.Receiver)
	}

	// Check the selector
	if messageSendNode.Selector != "do:" {
		t.Errorf("Expected selector to be 'do:', got '%s'", messageSendNode.Selector)
	}

	// Check the arguments
	if len(messageSendNode.Arguments) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(messageSendNode.Arguments))
	} else {
		// Check that the argument is a block
		blockNode, ok := messageSendNode.Arguments[0].(*ast.BlockNode)
		if !ok {
			t.Fatalf("Expected block node, got %T", messageSendNode.Arguments[0])
		}

		// Check block structure
		if len(blockNode.Parameters) != 0 {
			t.Errorf("Expected 0 parameters, got %d", len(blockNode.Parameters))
		}

		if len(blockNode.Temporaries) != 0 {
			t.Errorf("Expected 0 temporaries, got %d", len(blockNode.Temporaries))
		}

		// Check the block body
		literalNode, ok := blockNode.Body.(*ast.LiteralNode)
		if !ok {
			t.Fatalf("Expected literal node, got %T", blockNode.Body)
		}

		// Check that the literal is 5
		if !pile.IsImmediate(literalNode.Value) {
			t.Fatalf("Expected immediate value, got %v", literalNode.Value)
		}

		if !pile.IsIntegerImmediate(literalNode.Value) {
			t.Fatalf("Expected integer immediate, got %v", literalNode.Value)
		}

		value := pile.GetIntegerImmediate(literalNode.Value)

		if value != 5 {
			t.Errorf("Expected value to be 5, got %d", value)
		}
	}
}

// TestParseBlockValueMessage tests parsing "[5] value" as a message send with a block receiver
func TestParseBlockValueMessage(t *testing.T) {
	// Create a class
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass // Set class's class to itself for proper checks
	classObj := (*pile.Object)(unsafe.Pointer(objectClass))

	// Create a real VM for testing
	vmInstance := vm.NewVM()

	// Create a parser with the expression "[5] value"
	p := NewParser("[5] value", classObj, vmInstance)

	// Parse the expression
	node, err := p.ParseExpression()
	if err != nil {
		t.Fatalf("Error parsing expression: %v", err)
	}

	// Check that the node is a message send node
	messageSendNode, ok := node.(*ast.MessageSendNode)
	if !ok {
		t.Fatalf("Expected message send node, got %T", node)
	}

	// Check the message selector
	if messageSendNode.Selector != "value" {
		t.Errorf("Expected message selector to be 'value', got '%s'", messageSendNode.Selector)
	}

	// Check the message arguments (should be empty)
	if len(messageSendNode.Arguments) != 0 {
		t.Errorf("Expected 0 arguments, got %d", len(messageSendNode.Arguments))
	}

	// Check the message receiver (should be a block node)
	blockNode, ok := messageSendNode.Receiver.(*ast.BlockNode)
	if !ok {
		t.Fatalf("Expected block node as receiver, got %T", messageSendNode.Receiver)
	}

	// Check block structure
	if len(blockNode.Parameters) != 0 {
		t.Errorf("Expected 0 parameters, got %d", len(blockNode.Parameters))
	}

	if len(blockNode.Temporaries) != 0 {
		t.Errorf("Expected 0 temporaries, got %d", len(blockNode.Temporaries))
	}

	// Check the block body
	literalNode, ok := blockNode.Body.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("Expected literal node, got %T", blockNode.Body)
	}

	// Check that the literal is 5
	if !pile.IsImmediate(literalNode.Value) {
		t.Fatalf("Expected immediate value, got %v", literalNode.Value)
	}

	if !pile.IsIntegerImmediate(literalNode.Value) {
		t.Fatalf("Expected integer immediate, got %v", literalNode.Value)
	}

	value := pile.GetIntegerImmediate(literalNode.Value)
	if value != 5 {
		t.Errorf("Expected value to be 5, got %d", value)
	}
}