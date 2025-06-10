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

	// Parse the expression
	node, err := p.ParseExpression()
	if err != nil {
		t.Fatalf("Error parsing expression: %v", err)
	}

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
