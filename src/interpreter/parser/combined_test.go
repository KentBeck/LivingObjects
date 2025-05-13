package parser

import (
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/vm"
)

// TestCombined runs multiple parser tests in a controlled sequence
func TestCombined(t *testing.T) {
	// Create a class for context
	objectClass := core.NewClass("Object", nil)
	objectClass.ClassField = objectClass // Set class's class to itself
	classObj := (*core.Object)(unsafe.Pointer(objectClass))
	
	// Create a VM for testing
	vmInstance := vm.NewVM()
	
	t.Run("TestBoolean", func(t *testing.T) {
		runBooleanTest(t, "true", true, classObj, vmInstance)
		runBooleanTest(t, "false", false, classObj, vmInstance)
	})
	
	t.Run("TestInteger", func(t *testing.T) {
		runIntegerTest(t, "42", 42, classObj, vmInstance)
	})
	
	t.Run("TestArray", func(t *testing.T) {
		// Create a parser for the array expression
		p := NewParser("#(1 2 3)", classObj, vmInstance)
		
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
		
		// Basic check for array type
		literalNode, ok := node.(*ast.LiteralNode)
		if !ok {
			t.Fatalf("Expected LiteralNode, got %T", node)
		}
		
		if literalNode.Value.Type() != core.OBJ_ARRAY {
			t.Fatalf("Expected array type, got %v", literalNode.Value.Type())
		}
	})
}

func runBooleanTest(t *testing.T, input string, expectedValue bool, classObj *core.Object, vmInstance *vm.VM) {
	// Create a parser for the expression
	p := NewParser(input, classObj, vmInstance)
	
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
	
	// Should be a literal node with expected boolean value
	literalNode, ok := node.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("Expected LiteralNode, got %T", node)
	}
	
	// Check for expected immediate value
	if expectedValue {
		if !core.IsTrueImmediate(literalNode.Value) {
			t.Fatalf("Expected true immediate value, got %v", literalNode.Value)
		}
	} else {
		if !core.IsFalseImmediate(literalNode.Value) {
			t.Fatalf("Expected false immediate value, got %v", literalNode.Value)
		}
	}
}

func runIntegerTest(t *testing.T, input string, expectedValue int, classObj *core.Object, vmInstance *vm.VM) {
	// Create a parser for the expression
	p := NewParser(input, classObj, vmInstance)
	
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
	
	// Should be a literal node with the expected integer value
	literalNode, ok := node.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("Expected LiteralNode, got %T", node)
	}
	
	if !core.IsIntegerImmediate(literalNode.Value) {
		t.Fatalf("Expected integer immediate, got %v", literalNode.Value)
	}
	
	value := core.GetIntegerImmediate(literalNode.Value)
	if int(value) != expectedValue {
		t.Fatalf("Expected value %d, got %d", expectedValue, value)
	}
}