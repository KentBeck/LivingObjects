package parser

import (
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

// TestImmediateHandling tests the handling of immediate values in minimal context
func TestImmediateHandling(t *testing.T) {
	// Create a class for context
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass // Set class's class to itself
	classObj := (*pile.Object)(unsafe.Pointer(objectClass))
	
	// Create a VM for testing
	vmInstance := vm.NewVM()
	
	// Test parsing true
	testTrueImmediate(t, classObj, vmInstance)
	
	// Test parsing false
	testFalseImmediate(t, classObj, vmInstance)
}

func testTrueImmediate(t *testing.T, classObj *pile.Object, vmInstance *vm.VM) {
	// Create a parser for the expression
	p := NewParser("true", classObj, vmInstance)
	
	// Initialize tokens
	err := p.tokenize()
	if err != nil {
		t.Fatalf("Error tokenizing input: %v", err)
	}
	
	// Initialize current token
	p.CurrentToken = p.Tokens[0]
	p.CurrentTokenIndex = 0
	
	// Parse the expression
	node, err := p.ParseExpression()
	if err != nil {
		t.Fatalf("Error parsing expression: %v", err)
	}
	
	// Should be a literal node with true
	literalNode, ok := node.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("Expected LiteralNode, got %T", node)
	}
	
	// Check for true immediate
	if !pile.IsTrueImmediate(literalNode.Value) {
		t.Fatalf("Expected true immediate value, got %v", literalNode.Value)
	}
}

func testFalseImmediate(t *testing.T, classObj *pile.Object, vmInstance *vm.VM) {
	// Create a parser for the expression
	p := NewParser("false", classObj, vmInstance)
	
	// Initialize tokens
	err := p.tokenize()
	if err != nil {
		t.Fatalf("Error tokenizing input: %v", err)
	}
	
	// Initialize current token
	p.CurrentToken = p.Tokens[0]
	p.CurrentTokenIndex = 0
	
	// Parse the expression
	node, err := p.ParseExpression()
	if err != nil {
		t.Fatalf("Error parsing expression: %v", err)
	}
	
	// Should be a literal node with false
	literalNode, ok := node.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("Expected LiteralNode, got %T", node)
	}
	
	// Check for false immediate
	if !pile.IsFalseImmediate(literalNode.Value) {
		t.Fatalf("Expected false immediate value, got %v", literalNode.Value)
	}
}