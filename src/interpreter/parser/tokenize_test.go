package parser

import (
	"testing"
)

// TestTokenizeBlockValue tests the tokenization of the "[5] value" expression
func TestTokenizeBlockValue(t *testing.T) {
	// Create a parser with the test input
	p := NewParser("[5] value", nil, nil)

	// Tokenize the input
	err := p.tokenize()
	if err != nil {
		t.Fatalf("Error tokenizing input: %v", err)
	}

	// Verify token count
	expectedTokenCount := 5 // [, 5, ], value, EOF
	if len(p.Tokens) != expectedTokenCount {
		t.Errorf("Expected %d tokens, got %d", expectedTokenCount, len(p.Tokens))
	}

	// Verify token types and values
	if len(p.Tokens) >= expectedTokenCount {
		// Token 0: [
		if p.Tokens[0].Type != TOKEN_SPECIAL || p.Tokens[0].Value != "[" {
			t.Errorf("Expected Token 0 to be [ (special), got %v", p.Tokens[0])
		}

		// Token 1: 5
		if p.Tokens[1].Type != TOKEN_NUMBER || p.Tokens[1].Value != "5" {
			t.Errorf("Expected Token 1 to be 5 (number), got %v", p.Tokens[1])
		}

		// Token 2: ]
		if p.Tokens[2].Type != TOKEN_SPECIAL || p.Tokens[2].Value != "]" {
			t.Errorf("Expected Token 2 to be ] (special), got %v", p.Tokens[2])
		}

		// Token 3: value
		if p.Tokens[3].Type != TOKEN_IDENTIFIER || p.Tokens[3].Value != "value" {
			t.Errorf("Expected Token 3 to be value (identifier), got %v", p.Tokens[3])
		}

		// Token 4: EOF
		if p.Tokens[4].Type != TOKEN_EOF {
			t.Errorf("Expected Token 4 to be EOF, got %v", p.Tokens[4])
		}
	}
}
