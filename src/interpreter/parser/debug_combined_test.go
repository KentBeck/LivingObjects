package parser

import (
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

// TestDebugCombined is a simplified version of TestCombined for debugging purposes
func TestDebugCombined(t *testing.T) {
	// Create a class for context
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass
	classObj := (*pile.Object)(unsafe.Pointer(objectClass))

	// Create a VM instance
	vmInstance := vm.NewVM()

	// Test parsing a boolean value
	t.Run("TestBooleanTrue", func(t *testing.T) {
		// Create a parser
		parser := NewParser("true", classObj, vmInstance)

		// Parse the expression
		node, err := parser.ParseExpression()
		if err != nil {
			t.Fatalf("Error parsing expression: %v", err)
		}

		// Check that the node is a literal node
		literalNode, ok := node.(*ast.LiteralNode)
		if !ok {
			t.Fatalf("Expected literal node, got %T", node)
		}

		// Check that the literal is true
		if !pile.IsTrueImmediate(literalNode.Value) {
			t.Fatalf("Expected true immediate, got %v", literalNode.Value)
		}
	})

	// Test parsing a boolean value
	t.Run("TestBooleanFalse", func(t *testing.T) {
		// Create a parser
		parser := NewParser("false", classObj, vmInstance)

		// Parse the expression
		node, err := parser.ParseExpression()
		if err != nil {
			t.Fatalf("Error parsing expression: %v", err)
		}

		// Check that the node is a literal node
		literalNode, ok := node.(*ast.LiteralNode)
		if !ok {
			t.Fatalf("Expected literal node, got %T", node)
		}

		// Check that the literal is false
		if !pile.IsFalseImmediate(literalNode.Value) {
			t.Fatalf("Expected false immediate, got %v", literalNode.Value)
		}
	})
}