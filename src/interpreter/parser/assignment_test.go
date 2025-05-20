package parser

import (
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

// TestAssignmentExpression tests parsing a simple assignment expression
func TestAssignmentExpression(t *testing.T) {
	// Create a class for context
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass
	classObj := (*pile.Object)(unsafe.Pointer(objectClass))

	// Create a VM for testing
	vmInstance := vm.NewVM()

	// Create a parser with the test input
	p := NewParser("x := 5", classObj, vmInstance)

	// Tokenize the input manually to see what's happening
	err := p.tokenize()
	if err != nil {
		t.Fatalf("Error tokenizing input: %v", err)
	}

	// Print the tokens
	t.Logf("Tokens for x := 5:")
	for i, token := range p.Tokens {
		t.Logf("  Token %d: Type=%d, Value=%s", i, token.Type, token.Value)
		// Debug information for assignment
		if token.Value == ":=" {
			t.Logf("  Assignment token found: Type=%d, ASSIGNMENT=%d", token.Type, TOKEN_ASSIGNMENT)
		}
	}

	// Reset the parser state
	p.CurrentToken = p.Tokens[0]
	p.CurrentTokenIndex = 0

	// Add a debug print before parsing to see the current token
	t.Logf("Before parsing - Current token: Type=%d, Value=%s", p.CurrentToken.Type, p.CurrentToken.Value)
	
	// Debug the isAssignment check
	t.Logf("Is assignment check: %v", p.isAssignment())
	
	// Debug the next token
	if p.CurrentTokenIndex+1 < len(p.Tokens) {
		t.Logf("Next token: Type=%d, Value=%s", p.Tokens[p.CurrentTokenIndex+1].Type, p.Tokens[p.CurrentTokenIndex+1].Value)
	}
	
	// Parse the expression
	node, err := p.ParseExpression()
	if err != nil {
		t.Fatalf("Error parsing expression: %v", err)
	}

	// Print node type for debugging
	t.Logf("Node type: %T", node)

	// Check if it's an assignment node
	assignmentNode, ok := node.(*ast.AssignmentNode)
	if !ok {
		t.Fatalf("Expected AssignmentNode, got %T", node)
		return
	}

	// Check the variable name
	if assignmentNode.Variable != "x" {
		t.Errorf("Expected variable name 'x', got '%s'", assignmentNode.Variable)
	}

	// Check the expression
	literalNode, ok := assignmentNode.Expression.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("Expected LiteralNode as expression, got %T", assignmentNode.Expression)
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