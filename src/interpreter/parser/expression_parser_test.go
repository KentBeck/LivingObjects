package parser

import (
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/vm"
)

// TestParseExpression tests parsing various Smalltalk expressions
func TestParseExpression(t *testing.T) {
	// Create a class for context
	objectClass := core.NewClass("Object", nil)
	classObj := (*core.Object)(unsafe.Pointer(objectClass))

	// Test cases
	tests := []struct {
		name     string
		input    string
		validate func(t *testing.T, node ast.Node)
	}{
		{
			name:  "Simple integer literal",
			input: "42",
			validate: func(t *testing.T, node ast.Node) {
				// Should be a literal node with value 42
				literalNode, ok := node.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("Expected LiteralNode, got %T", node)
				}

				if !core.IsIntegerImmediate(literalNode.Value) {
					t.Fatalf("Expected integer immediate, got %v", literalNode.Value)
				}

				value := core.GetIntegerImmediate(literalNode.Value)
				if value != 42 {
					t.Errorf("Expected value 42, got %d", value)
				}
			},
		},
		{
			name:  "Simple binary message",
			input: "2 + 3",
			validate: func(t *testing.T, node ast.Node) {
				// Should be a message send node: 2 + 3
				messageSendNode, ok := node.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected MessageSendNode, got %T", node)
				}

				// Check selector
				if messageSendNode.Selector != "+" {
					t.Errorf("Expected selector '+', got '%s'", messageSendNode.Selector)
				}

				// Check receiver (should be 2)
				receiverNode, ok := messageSendNode.Receiver.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("Expected receiver to be LiteralNode, got %T", messageSendNode.Receiver)
				}

				if !core.IsIntegerImmediate(receiverNode.Value) {
					t.Fatalf("Expected receiver to be integer immediate, got %v", receiverNode.Value)
				}

				receiverValue := core.GetIntegerImmediate(receiverNode.Value)
				if receiverValue != 2 {
					t.Errorf("Expected receiver value 2, got %d", receiverValue)
				}

				// Check argument (should be 3)
				if len(messageSendNode.Arguments) != 1 {
					t.Fatalf("Expected 1 argument, got %d", len(messageSendNode.Arguments))
				}

				argNode, ok := messageSendNode.Arguments[0].(*ast.LiteralNode)
				if !ok {
					t.Fatalf("Expected argument to be LiteralNode, got %T", messageSendNode.Arguments[0])
				}

				if !core.IsIntegerImmediate(argNode.Value) {
					t.Fatalf("Expected argument to be integer immediate, got %v", argNode.Value)
				}

				argValue := core.GetIntegerImmediate(argNode.Value)
				if argValue != 3 {
					t.Errorf("Expected argument value 3, got %d", argValue)
				}
			},
		},
		{
			name:  "Multiplication",
			input: "3 * 4",
			validate: func(t *testing.T, node ast.Node) {
				// Should be a message send node: 3 * 4
				messageSendNode, ok := node.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected MessageSendNode, got %T", node)
				}

				// Check selector
				if messageSendNode.Selector != "*" {
					t.Errorf("Expected selector '*', got '%s'", messageSendNode.Selector)
				}

				// Check receiver (should be 3)
				receiverNode, ok := messageSendNode.Receiver.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("Expected receiver to be LiteralNode, got %T", messageSendNode.Receiver)
				}

				if !core.IsIntegerImmediate(receiverNode.Value) {
					t.Fatalf("Expected receiver to be integer immediate, got %v", receiverNode.Value)
				}

				receiverValue := core.GetIntegerImmediate(receiverNode.Value)
				if receiverValue != 3 {
					t.Errorf("Expected receiver value 3, got %d", receiverValue)
				}

				// Check argument (should be 4)
				if len(messageSendNode.Arguments) != 1 {
					t.Fatalf("Expected 1 argument, got %d", len(messageSendNode.Arguments))
				}

				argNode, ok := messageSendNode.Arguments[0].(*ast.LiteralNode)
				if !ok {
					t.Fatalf("Expected argument to be LiteralNode, got %T", messageSendNode.Arguments[0])
				}

				if !core.IsIntegerImmediate(argNode.Value) {
					t.Fatalf("Expected argument to be integer immediate, got %v", argNode.Value)
				}

				argValue := core.GetIntegerImmediate(argNode.Value)
				if argValue != 4 {
					t.Errorf("Expected argument value 4, got %d", argValue)
				}
			},
		},
		{
			name:  "Multiple binary messages (left associative)",
			input: "2 + 2 * 3",
			validate: func(t *testing.T, node ast.Node) {
				// Should be a message send node: (2 + 2) * 3
				messageSendNode, ok := node.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected MessageSendNode, got %T", node)
				}

				// Check selector
				if messageSendNode.Selector != "*" {
					t.Errorf("Expected selector '*', got '%s'", messageSendNode.Selector)
				}

				// Check receiver (should be 2 + 2)
				receiverNode, ok := messageSendNode.Receiver.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected receiver to be MessageSendNode, got %T", messageSendNode.Receiver)
				}

				if receiverNode.Selector != "+" {
					t.Errorf("Expected receiver selector '+', got '%s'", receiverNode.Selector)
				}

				// Check argument (should be 3)
				if len(messageSendNode.Arguments) != 1 {
					t.Fatalf("Expected 1 argument, got %d", len(messageSendNode.Arguments))
				}

				argNode, ok := messageSendNode.Arguments[0].(*ast.LiteralNode)
				if !ok {
					t.Fatalf("Expected argument to be LiteralNode, got %T", messageSendNode.Arguments[0])
				}

				if !core.IsIntegerImmediate(argNode.Value) {
					t.Fatalf("Expected argument to be integer immediate, got %v", argNode.Value)
				}

				argValue := core.GetIntegerImmediate(argNode.Value)
				if argValue != 3 {
					t.Errorf("Expected argument value 3, got %d", argValue)
				}
			},
		},
		{
			name:  "Parenthesized expression",
			input: "(2 + 2) * 3",
			validate: func(t *testing.T, node ast.Node) {
				// Should be a message send node: (2 + 2) * 3
				messageSendNode, ok := node.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected MessageSendNode, got %T", node)
				}

				// Check selector
				if messageSendNode.Selector != "*" {
					t.Errorf("Expected selector '*', got '%s'", messageSendNode.Selector)
				}

				// Check receiver (should be 2 + 2)
				receiverNode, ok := messageSendNode.Receiver.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected receiver to be MessageSendNode, got %T", messageSendNode.Receiver)
				}

				if receiverNode.Selector != "+" {
					t.Errorf("Expected receiver selector '+', got '%s'", receiverNode.Selector)
				}

				// Check argument (should be 3)
				if len(messageSendNode.Arguments) != 1 {
					t.Fatalf("Expected 1 argument, got %d", len(messageSendNode.Arguments))
				}

				argNode, ok := messageSendNode.Arguments[0].(*ast.LiteralNode)
				if !ok {
					t.Fatalf("Expected argument to be LiteralNode, got %T", messageSendNode.Arguments[0])
				}

				if !core.IsIntegerImmediate(argNode.Value) {
					t.Fatalf("Expected argument to be integer immediate, got %v", argNode.Value)
				}

				argValue := core.GetIntegerImmediate(argNode.Value)
				if argValue != 3 {
					t.Errorf("Expected argument value 3, got %d", argValue)
				}
			},
		},
		{
			name:  "Chained binary messages",
			input: "1 + 2 + 3",
			validate: func(t *testing.T, node ast.Node) {
				// Should be a message send node: (1 + 2) + 3
				messageSendNode, ok := node.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected MessageSendNode, got %T", node)
				}

				// Check selector
				if messageSendNode.Selector != "+" {
					t.Errorf("Expected selector '+', got '%s'", messageSendNode.Selector)
				}

				// Check receiver (should be 1 + 2)
				receiverNode, ok := messageSendNode.Receiver.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected receiver to be MessageSendNode, got %T", messageSendNode.Receiver)
				}

				if receiverNode.Selector != "+" {
					t.Errorf("Expected receiver selector '+', got '%s'", receiverNode.Selector)
				}

				// Check argument (should be 3)
				if len(messageSendNode.Arguments) != 1 {
					t.Fatalf("Expected 1 argument, got %d", len(messageSendNode.Arguments))
				}

				argNode, ok := messageSendNode.Arguments[0].(*ast.LiteralNode)
				if !ok {
					t.Fatalf("Expected argument to be LiteralNode, got %T", messageSendNode.Arguments[0])
				}

				if !core.IsIntegerImmediate(argNode.Value) {
					t.Fatalf("Expected argument to be integer immediate, got %v", argNode.Value)
				}

				argValue := core.GetIntegerImmediate(argNode.Value)
				if argValue != 3 {
					t.Errorf("Expected argument value 3, got %d", argValue)
				}
			},
		},
		{
			name:  "Keyword message",
			input: "1 to: 3",
			validate: func(t *testing.T, node ast.Node) {
				// Should be a message send node: 1 to: 3
				messageSendNode, ok := node.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected MessageSendNode, got %T", node)
				}

				// Check selector
				if messageSendNode.Selector != "to:" {
					t.Errorf("Expected selector 'to:', got '%s'", messageSendNode.Selector)
				}

				// Check receiver (should be 1)
				receiverNode, ok := messageSendNode.Receiver.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("Expected receiver to be LiteralNode, got %T", messageSendNode.Receiver)
				}

				if !core.IsIntegerImmediate(receiverNode.Value) {
					t.Fatalf("Expected receiver to be integer immediate, got %v", receiverNode.Value)
				}

				receiverValue := core.GetIntegerImmediate(receiverNode.Value)
				if receiverValue != 1 {
					t.Errorf("Expected receiver value 1, got %d", receiverValue)
				}

				// Check argument (should be 3)
				if len(messageSendNode.Arguments) != 1 {
					t.Fatalf("Expected 1 argument, got %d", len(messageSendNode.Arguments))
				}

				argNode, ok := messageSendNode.Arguments[0].(*ast.LiteralNode)
				if !ok {
					t.Fatalf("Expected argument to be LiteralNode, got %T", messageSendNode.Arguments[0])
				}

				if !core.IsIntegerImmediate(argNode.Value) {
					t.Fatalf("Expected argument to be integer immediate, got %v", argNode.Value)
				}

				argValue := core.GetIntegerImmediate(argNode.Value)
				if argValue != 3 {
					t.Errorf("Expected argument value 3, got %d", argValue)
				}
			},
		},
		{
			name:  "String concatenation",
			input: "'hello' , ' world'",
			validate: func(t *testing.T, node ast.Node) {
				// Should be a message send node: 'hello' , ' world'
				messageSendNode, ok := node.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected MessageSendNode, got %T", node)
				}

				// Check selector
				if messageSendNode.Selector != "," {
					t.Errorf("Expected selector ',', got '%s'", messageSendNode.Selector)
				}

				// Check receiver (should be 'hello')
				_, isLiteral := messageSendNode.Receiver.(*ast.LiteralNode)
				if !isLiteral {
					t.Fatalf("Expected receiver to be LiteralNode, got %T", messageSendNode.Receiver)
				}

				// Check argument (should be ' world')
				if len(messageSendNode.Arguments) != 1 {
					t.Fatalf("Expected 1 argument, got %d", len(messageSendNode.Arguments))
				}

				_, isLiteral = messageSendNode.Arguments[0].(*ast.LiteralNode)
				if !isLiteral {
					t.Fatalf("Expected argument to be LiteralNode, got %T", messageSendNode.Arguments[0])
				}
			},
		},
		{
			name:  "Unary message",
			input: "'hello' size",
			validate: func(t *testing.T, node ast.Node) {
				// Should be a message send node: 'hello' size
				messageSendNode, ok := node.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected MessageSendNode, got %T", node)
				}

				// Check selector
				if messageSendNode.Selector != "size" {
					t.Errorf("Expected selector 'size', got '%s'", messageSendNode.Selector)
				}

				// Check receiver (should be 'hello')
				_, isLiteral := messageSendNode.Receiver.(*ast.LiteralNode)
				if !isLiteral {
					t.Fatalf("Expected receiver to be LiteralNode, got %T", messageSendNode.Receiver)
				}

				// Check arguments (should be empty)
				if len(messageSendNode.Arguments) != 0 {
					t.Fatalf("Expected 0 arguments, got %d", len(messageSendNode.Arguments))
				}
			},
		},
		{
			name:  "Boolean unary message",
			input: "true not",
			validate: func(t *testing.T, node ast.Node) {
				// Should be a message send node: true not
				messageSendNode, ok := node.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected MessageSendNode, got %T", node)
				}

				// Check selector
				if messageSendNode.Selector != "not" {
					t.Errorf("Expected selector 'not', got '%s'", messageSendNode.Selector)
				}

				// Check receiver (should be true)
				receiverNode, ok := messageSendNode.Receiver.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("Expected receiver to be LiteralNode, got %T", messageSendNode.Receiver)
				}

				if !core.IsTrueImmediate(receiverNode.Value) {
					t.Fatalf("Expected receiver to be true immediate, got %v", receiverNode.Value)
				}

				// Check arguments (should be empty)
				if len(messageSendNode.Arguments) != 0 {
					t.Fatalf("Expected 0 arguments, got %d", len(messageSendNode.Arguments))
				}
			},
		},
		{
			name:  "Boolean unary message with false",
			input: "false not",
			validate: func(t *testing.T, node ast.Node) {
				// Should be a message send node: false not
				messageSendNode, ok := node.(*ast.MessageSendNode)
				if !ok {
					t.Fatalf("Expected MessageSendNode, got %T", node)
				}

				// Check selector
				if messageSendNode.Selector != "not" {
					t.Errorf("Expected selector 'not', got '%s'", messageSendNode.Selector)
				}

				// Check receiver (should be false)
				receiverNode, ok := messageSendNode.Receiver.(*ast.LiteralNode)
				if !ok {
					t.Fatalf("Expected receiver to be LiteralNode, got %T", messageSendNode.Receiver)
				}

				if !core.IsFalseImmediate(receiverNode.Value) {
					t.Fatalf("Expected receiver to be false immediate, got %v", receiverNode.Value)
				}

				// Check arguments (should be empty)
				if len(messageSendNode.Arguments) != 0 {
					t.Fatalf("Expected 0 arguments, got %d", len(messageSendNode.Arguments))
				}
			},
		},
	}

	// Create a real VM for testing
	vmInstance := vm.NewVM()

	// Run the tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a parser for the expression
			p := NewParser(test.input, classObj, vmInstance)

			// Initialize tokens
			err := p.tokenize()
			if err != nil {
				t.Fatalf("Error tokenizing input: %v", err)
			}

			// Debug: print tokens
			if test.name == "Parenthesized expression" {
				t.Logf("Tokens for %s:", test.input)
				for i, token := range p.Tokens {
					t.Logf("  Token %d: Type=%d, Value=%s", i, token.Type, token.Value)
				}
			}

			// Initialize current token
			p.CurrentToken = p.Tokens[0]
			p.CurrentTokenIndex = 0

			// Parse the expression
			node, err := p.parseExpression()
			if err != nil {
				t.Fatalf("Error parsing expression: %v", err)
			}

			// Validate the result
			test.validate(t, node)
		})
	}
}

// TestParseArrayLiteral tests parsing array literals
func TestParseArrayLiteral(t *testing.T) {
	// Create a class for context
	objectClass := core.NewClass("Object", nil)
	classObj := (*core.Object)(unsafe.Pointer(objectClass))

	// Create a real VM for testing
	vmInstance := vm.NewVM()

	// Create a parser for the expression
	p := NewParser("#(1 2 3)", classObj, vmInstance)

	// Initialize tokens
	err := p.tokenize()
	if err != nil {
		t.Fatalf("Error tokenizing input: %v", err)
	}

	// Debug: print tokens
	t.Logf("Tokens for #(1 2 3):")
	for i, token := range p.Tokens {
		t.Logf("  Token %d: Type=%d, Value=%s", i, token.Type, token.Value)
	}

	// Initialize current token
	p.CurrentToken = p.Tokens[0]
	p.CurrentTokenIndex = 0

	// Parse the expression
	node, err := p.parseExpression()
	if err != nil {
		t.Fatalf("Error parsing expression: %v", err)
	}

	// Should be a literal node with an Array object
	literalNode, ok := node.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("Expected LiteralNode, got %T", node)
	}

	// Check that the value is an Array object
	if literalNode.Value.Type() != core.OBJ_ARRAY {
		t.Fatalf("Expected Array object, got %v", literalNode.Value.Type())
	}

	// Convert to Array and check its properties
	array := classes.ObjectToArray(literalNode.Value)

	// Check array size
	if array.Size() != 3 {
		t.Fatalf("Expected array size 3, got %d", array.Size())
	}

	// Check first element (should be 1)
	elem1 := array.At(0)
	if !core.IsIntegerImmediate(elem1) {
		t.Fatalf("Expected element 1 to be integer immediate, got %v", elem1)
	}

	elem1Value := core.GetIntegerImmediate(elem1)
	if elem1Value != 1 {
		t.Errorf("Expected element 1 value 1, got %d", elem1Value)
	}

	// Check second element (should be 2)
	elem2 := array.At(1)
	if !core.IsIntegerImmediate(elem2) {
		t.Fatalf("Expected element 2 to be integer immediate, got %v", elem2)
	}

	elem2Value := core.GetIntegerImmediate(elem2)
	if elem2Value != 2 {
		t.Errorf("Expected element 2 value 2, got %d", elem2Value)
	}

	// Check third element (should be 3)
	elem3 := array.At(2)
	if !core.IsIntegerImmediate(elem3) {
		t.Fatalf("Expected element 3 to be integer immediate, got %v", elem3)
	}

	elem3Value := core.GetIntegerImmediate(elem3)
	if elem3Value != 3 {
		t.Errorf("Expected element 3 value 3, got %d", elem3Value)
	}
}

// TestArrayLiteralWithKeywordMessage tests parsing array literals with keyword messages
func TestArrayLiteralWithKeywordMessage(t *testing.T) {
	// Create a class for context
	objectClass := core.NewClass("Object", nil)
	classObj := (*core.Object)(unsafe.Pointer(objectClass))

	// Create a real VM for testing
	vmInstance := vm.NewVM()

	// Create a parser for the expression
	p := NewParser("#(1 2 3) at: 2", classObj, vmInstance)

	// Initialize tokens
	err := p.tokenize()
	if err != nil {
		t.Fatalf("Error tokenizing input: %v", err)
	}

	// Initialize current token
	p.CurrentToken = p.Tokens[0]
	p.CurrentTokenIndex = 0

	// Parse the expression
	node, err := p.parseExpression()
	if err != nil {
		t.Fatalf("Error parsing expression: %v", err)
	}

	// Should be a message send node: #(1 2 3) at: 2
	messageSendNode, ok := node.(*ast.MessageSendNode)
	if !ok {
		t.Fatalf("Expected MessageSendNode, got %T", node)
	}

	// Check selector
	if messageSendNode.Selector != "at:" {
		t.Errorf("Expected selector 'at:', got '%s'", messageSendNode.Selector)
	}

	// Check receiver (should be a LiteralNode with an Array object)
	receiverNode, ok := messageSendNode.Receiver.(*ast.LiteralNode)
	if !ok {
		t.Fatalf("Expected receiver to be LiteralNode, got %T", messageSendNode.Receiver)
	}

	// Check that the value is an Array object
	if receiverNode.Value.Type() != core.OBJ_ARRAY {
		t.Fatalf("Expected Array object, got %v", receiverNode.Value.Type())
	}

	// Convert to Array and check its properties
	array := classes.ObjectToArray(receiverNode.Value)

	// Check array size
	if array.Size() != 3 {
		t.Fatalf("Expected array size 3, got %d", array.Size())
	}

	// Check argument (should be 2)
	if len(messageSendNode.Arguments) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(messageSendNode.Arguments))
	}

	argNode, ok := messageSendNode.Arguments[0].(*ast.LiteralNode)
	if !ok {
		t.Fatalf("Expected argument to be LiteralNode, got %T", messageSendNode.Arguments[0])
	}

	if !core.IsIntegerImmediate(argNode.Value) {
		t.Fatalf("Expected argument to be integer immediate, got %v", argNode.Value)
	}

	argValue := core.GetIntegerImmediate(argNode.Value)
	if argValue != 2 {
		t.Errorf("Expected argument value 2, got %d", argValue)
	}
}
