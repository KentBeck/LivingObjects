package parser

import (
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

// TestDirectParseBlockValue directly tests parsing "[5] value" as a message send with a block receiver
func TestDirectParseBlockValue(t *testing.T) {
	// Create a class for context
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass
	classObj := (*pile.Object)(unsafe.Pointer(objectClass))

	// Create a VM for testing
	vmInstance := vm.NewVM()

	// Create a parser with the test input
	p := NewParser("[5] value", classObj, vmInstance)

	// Parse the expression
	node, err := p.ParseExpression()
	if err != nil {
		t.Fatalf("Error parsing expression: %v", err)
	}

	// Check if it's a message send node
	messageSend, ok := node.(*ast.MessageSendNode)
	if !ok {
		t.Fatalf("Expected MessageSendNode, got %T", node)
		return
	}

	// Check the selector
	if messageSend.Selector != "value" {
		t.Errorf("Expected selector 'value', got '%s'", messageSend.Selector)
	}

	// Check if receiver is a block
	blockNode, ok := messageSend.Receiver.(*ast.BlockNode)
	if !ok {
		t.Fatalf("Expected BlockNode as receiver, got %T", messageSend.Receiver)
		return
	}

	// The block body can be a literal node directly or a message send node
	// depending on how the parser is implemented
	switch body := blockNode.Body.(type) {
	case *ast.LiteralNode:
		// Verify it's an integer with value 5
		if !pile.IsIntegerImmediate(body.Value) {
			t.Fatalf("Expected integer immediate, got %v", body.Value)
		}

		value := pile.GetIntegerImmediate(body.Value)
		if value != 5 {
			t.Errorf("Expected value 5, got %d", value)
		}
	default:
		t.Fatalf("Expected LiteralNode as block body, got %T", blockNode.Body)
	}
}
