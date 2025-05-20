package parser

import (
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

// TestBlockValueKeyword directly tests parsing "[:x | x] value: 5" as a keyword message send with a block receiver
func TestBlockValueKeyword(t *testing.T) {
	// Create a class for context
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass
	classObj := (*pile.Object)(unsafe.Pointer(objectClass))

	// Create a VM for testing
	vmInstance := vm.NewVM()

	// Create a parser with the test input
	p := NewParser("[:x | x] value: 5", classObj, vmInstance)

	// Tokenize the input manually to see what's happening
	err := p.tokenize()
	if err != nil {
		t.Fatalf("Error tokenizing input: %v", err)
	}

	// Print the tokens
	t.Logf("Tokens for [:x | x] value: 5:")
	for i, token := range p.Tokens {
		t.Logf("  Token %d: Type=%d, Value=%s", i, token.Type, token.Value)
	}

	// Reset the parser state
	p.CurrentToken = p.Tokens[0]
	p.CurrentTokenIndex = 0

	// Parse the expression
	node, err := p.ParseExpression()
	if err != nil {
		t.Fatalf("Error parsing expression: %v", err)
	}

	// Print node type for debugging
	t.Logf("Node type: %T", node)

	// Check if it's a message send node
	messageSend, ok := node.(*ast.MessageSendNode)
	if !ok {
		t.Fatalf("Expected MessageSendNode, got %T", node)
		return
	}

	// Check the selector
	if messageSend.Selector != "value:" {
		t.Errorf("Expected selector 'value:', got '%s'", messageSend.Selector)
	}

	// Check if receiver is a block
	blockNode, ok := messageSend.Receiver.(*ast.BlockNode)
	if !ok {
		t.Fatalf("Expected BlockNode as receiver, got %T", messageSend.Receiver)
		return
	}

	// Check block parameters
	if len(blockNode.Parameters) != 1 || blockNode.Parameters[0] != "x" {
		t.Errorf("Expected block parameter 'x', got %v", blockNode.Parameters)
	}

	// Check argument
	if len(messageSend.Arguments) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(messageSend.Arguments))
		return
	}

	// Check if argument is a literal node with value 5
	literalNode, ok := messageSend.Arguments[0].(*ast.LiteralNode)
	if !ok {
		t.Fatalf("Expected LiteralNode as argument, got %T", messageSend.Arguments[0])
		return
	}

	// Check the value of the literal node
	if !pile.IsIntegerImmediate(literalNode.Value) {
		t.Fatalf("Expected integer immediate, got %v", literalNode.Value)
	}

	value := pile.GetIntegerImmediate(literalNode.Value)
	if value != 5 {
		t.Errorf("Expected value 5, got %d", value)
	}
}