package parser

import (
	"fmt"
	"strings"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/core"
)

// Parser parses Smalltalk code into an AST
type Parser struct {
	// Input is the input string to parse
	Input string

	// Class is the class the method belongs to
	Class *core.Object

	// Position is the current position in the input
	Position int

	// CurrentChar is the current character being processed
	CurrentChar byte

	// Tokens are the tokens extracted from the input
	Tokens []Token

	// CurrentToken is the current token being processed
	CurrentToken Token

	// CurrentTokenIndex is the index of the current token
	CurrentTokenIndex int
}

// TokenType represents the type of a token
type TokenType int

const (
	// Token types
	TOKEN_IDENTIFIER TokenType = iota
	TOKEN_NUMBER
	TOKEN_STRING
	TOKEN_SYMBOL
	TOKEN_KEYWORD
	TOKEN_SPECIAL
	TOKEN_EOF
)

// Token represents a token in the input
type Token struct {
	// Type is the type of the token
	Type TokenType

	// Value is the value of the token
	Value string
}

// NewParser creates a new parser
func NewParser(input string, class *core.Object) *Parser {
	p := &Parser{
		Input:             input,
		Class:             class,
		Position:          0,
		CurrentTokenIndex: 0,
		Tokens:            []Token{},
	}

	if len(input) > 0 {
		p.CurrentChar = input[0]
	}

	return p
}

// Parse parses the input and returns an AST
func (p *Parser) Parse() (ast.Node, error) {
	// Tokenize the input
	err := p.tokenize()
	if err != nil {
		return nil, err
	}

	// Parse the method
	return p.parseMethod()
}

// tokenize tokenizes the input
func (p *Parser) tokenize() error {
	for p.Position < len(p.Input) {
		// Skip whitespace
		if p.isWhitespace(p.CurrentChar) {
			p.advance()
			continue
		}

		// Parse identifiers
		if p.isAlpha(p.CurrentChar) {
			p.Tokens = append(p.Tokens, p.parseIdentifier())
			continue
		}

		// Parse numbers
		if p.isDigit(p.CurrentChar) {
			p.Tokens = append(p.Tokens, p.parseNumber())
			continue
		}

		// Parse special characters
		if p.isSpecial(p.CurrentChar) {
			p.Tokens = append(p.Tokens, p.parseSpecial())
			continue
		}

		// Parse strings
		if p.CurrentChar == '\'' {
			token, err := p.parseString()
			if err != nil {
				return err
			}
			p.Tokens = append(p.Tokens, token)
			continue
		}

		// Parse symbols
		if p.CurrentChar == '#' {
			token, err := p.parseSymbol()
			if err != nil {
				return err
			}
			p.Tokens = append(p.Tokens, token)
			continue
		}

		// Skip comments
		if p.CurrentChar == '"' {
			err := p.skipComment()
			if err != nil {
				return err
			}
			continue
		}

		// Unknown character
		return fmt.Errorf("unknown character: %c", p.CurrentChar)
	}

	// Add EOF token
	p.Tokens = append(p.Tokens, Token{Type: TOKEN_EOF, Value: ""})

	return nil
}

// parseMethod parses a method
func (p *Parser) parseMethod() (ast.Node, error) {
	// Initialize the current token
	p.CurrentToken = p.Tokens[0]

	// Parse the method selector
	selector, parameters, err := p.parseMethodSelector()
	if err != nil {
		return nil, err
	}

	// Parse temporary variables
	temporaries, err := p.parseTemporaries()
	if err != nil {
		return nil, err
	}

	// Parse the method body
	body, err := p.parseStatements()
	if err != nil {
		return nil, err
	}

	// Create the method node
	methodNode := &ast.MethodNode{
		Selector:    selector,
		Parameters:  parameters,
		Temporaries: temporaries,
		Body:        body,
		Class:       p.Class,
	}

	return methodNode, nil
}

// parseMethodSelector parses a method selector
func (p *Parser) parseMethodSelector() (string, []string, error) {
	// Handle binary selectors
	if p.CurrentToken.Type == TOKEN_SPECIAL {
		selector := p.CurrentToken.Value
		p.advanceToken()

		// Parse the parameter
		if p.CurrentToken.Type != TOKEN_IDENTIFIER {
			return "", nil, fmt.Errorf("expected identifier, got %v", p.CurrentToken)
		}

		parameter := p.CurrentToken.Value
		p.advanceToken()

		return selector, []string{parameter}, nil
	}

	// Handle keyword selectors
	if p.CurrentToken.Type == TOKEN_IDENTIFIER && strings.HasSuffix(p.CurrentToken.Value, ":") {
		selector := p.CurrentToken.Value
		p.advanceToken()

		// Parse the parameter
		if p.CurrentToken.Type != TOKEN_IDENTIFIER {
			return "", nil, fmt.Errorf("expected identifier, got %v", p.CurrentToken)
		}

		parameter := p.CurrentToken.Value
		p.advanceToken()

		return selector, []string{parameter}, nil
	}

	// Handle unary selectors
	if p.CurrentToken.Type == TOKEN_IDENTIFIER {
		selector := p.CurrentToken.Value
		p.advanceToken()

		// No parameters for unary selectors
		return selector, []string{}, nil
	}

	return "", nil, fmt.Errorf("expected identifier or special, got %v", p.CurrentToken)
}

// parseTemporaries parses temporary variables
func (p *Parser) parseTemporaries() ([]string, error) {
	// Check if there are temporary variables
	if p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == "|" {
		p.advanceToken()

		// Parse the temporary variable names
		temporaries := []string{}

		// Parse each temporary variable
		for p.CurrentToken.Type == TOKEN_IDENTIFIER {
			temporaries = append(temporaries, p.CurrentToken.Value)
			p.advanceToken()
		}

		// Check for the closing |
		if p.CurrentToken.Type != TOKEN_SPECIAL || p.CurrentToken.Value != "|" {
			return nil, fmt.Errorf("expected |, got %v", p.CurrentToken)
		}

		p.advanceToken()

		return temporaries, nil
	}

	// No temporary variables
	return []string{}, nil
}

// parseStatements parses statements
func (p *Parser) parseStatements() (ast.Node, error) {
	// For now, we only handle a single return statement
	// Skip any statements before the return
	for p.CurrentToken.Type != TOKEN_EOF && !(p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == "^") {
		// Just advance to the next token
		p.advanceToken()

		// If we've reached the end of the tokens, break
		if p.CurrentTokenIndex >= len(p.Tokens) {
			break
		}
	}

	// Parse the return statement
	if p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == "^" {
		p.advanceToken()

		// Parse the expression
		expression, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		// Create the return node
		returnNode := &ast.ReturnNode{
			Expression: expression,
		}

		return returnNode, nil
	}

	return nil, fmt.Errorf("expected return statement, got %v", p.CurrentToken)
}

// parseExpression parses an expression
func (p *Parser) parseExpression() (ast.Node, error) {
	// Parse the primary expression
	primary, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	// Check if there's a message send
	if p.CurrentToken.Type == TOKEN_SPECIAL || p.CurrentToken.Type == TOKEN_IDENTIFIER || p.CurrentToken.Type == TOKEN_KEYWORD {
		return p.parseMessageSend(primary)
	}

	return primary, nil
}

// parsePrimary parses a primary expression
func (p *Parser) parsePrimary() (ast.Node, error) {
	// Handle self
	if p.CurrentToken.Type == TOKEN_IDENTIFIER && p.CurrentToken.Value == "self" {
		p.advanceToken()
		return &ast.SelfNode{}, nil
	}

	// Handle variables
	if p.CurrentToken.Type == TOKEN_IDENTIFIER {
		name := p.CurrentToken.Value
		p.advanceToken()
		return &ast.VariableNode{Name: name}, nil
	}

	return nil, fmt.Errorf("expected primary expression, got %v", p.CurrentToken)
}

// parseMessageSend parses a message send
func (p *Parser) parseMessageSend(receiver ast.Node) (ast.Node, error) {
	// Handle binary messages
	if p.CurrentToken.Type == TOKEN_SPECIAL {
		selector := p.CurrentToken.Value
		p.advanceToken()

		// Parse the argument
		argument, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}

		return &ast.MessageSendNode{
			Receiver:  receiver,
			Selector:  selector,
			Arguments: []ast.Node{argument},
		}, nil
	}

	// Handle unary messages
	if p.CurrentToken.Type == TOKEN_IDENTIFIER {
		selector := p.CurrentToken.Value
		p.advanceToken()

		return &ast.MessageSendNode{
			Receiver:  receiver,
			Selector:  selector,
			Arguments: []ast.Node{},
		}, nil
	}

	// Handle keyword messages (not implemented yet)

	return nil, fmt.Errorf("expected message selector, got %v", p.CurrentToken)
}

// advance advances to the next character
func (p *Parser) advance() {
	p.Position++
	if p.Position < len(p.Input) {
		p.CurrentChar = p.Input[p.Position]
	}
}

// advanceToken advances to the next token
func (p *Parser) advanceToken() {
	p.CurrentTokenIndex++
	if p.CurrentTokenIndex < len(p.Tokens) {
		p.CurrentToken = p.Tokens[p.CurrentTokenIndex]
	}
}

// isWhitespace returns true if the character is whitespace
func (p *Parser) isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// isAlpha returns true if the character is alphabetic
func (p *Parser) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// isDigit returns true if the character is a digit
func (p *Parser) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isSpecial returns true if the character is a special character
func (p *Parser) isSpecial(c byte) bool {
	return strings.ContainsRune("+-*/=<>[](){}^.|:", rune(c))
}

// parseIdentifier parses an identifier
func (p *Parser) parseIdentifier() Token {
	var value strings.Builder

	for p.Position < len(p.Input) && (p.isAlpha(p.CurrentChar) || p.isDigit(p.CurrentChar)) {
		value.WriteByte(p.CurrentChar)
		p.advance()
	}

	// Check if it's a keyword
	if p.Position < len(p.Input) && p.CurrentChar == ':' {
		value.WriteByte(':')
		p.advance()
		return Token{Type: TOKEN_IDENTIFIER, Value: value.String()}
	}

	return Token{Type: TOKEN_IDENTIFIER, Value: value.String()}
}

// parseNumber parses a number
func (p *Parser) parseNumber() Token {
	var value strings.Builder

	for p.Position < len(p.Input) && p.isDigit(p.CurrentChar) {
		value.WriteByte(p.CurrentChar)
		p.advance()
	}

	// Handle decimal point
	if p.Position < len(p.Input) && p.CurrentChar == '.' {
		// Make sure the next character is a digit
		if p.Position+1 < len(p.Input) && p.isDigit(p.Input[p.Position+1]) {
			value.WriteByte('.')
			p.advance()

			for p.Position < len(p.Input) && p.isDigit(p.CurrentChar) {
				value.WriteByte(p.CurrentChar)
				p.advance()
			}
		}
	}

	return Token{Type: TOKEN_NUMBER, Value: value.String()}
}

// parseSpecial parses a special character
func (p *Parser) parseSpecial() Token {
	value := string(p.CurrentChar)
	p.advance()
	return Token{Type: TOKEN_SPECIAL, Value: value}
}

// parseString parses a string
func (p *Parser) parseString() (Token, error) {
	var value strings.Builder

	// Skip the opening quote
	p.advance()

	for p.Position < len(p.Input) && p.CurrentChar != '\'' {
		// Handle escaped quotes
		if p.CurrentChar == '\'' && p.Position+1 < len(p.Input) && p.Input[p.Position+1] == '\'' {
			value.WriteByte('\'')
			p.advance() // Skip the first quote
			p.advance() // Skip the second quote
			continue
		}

		value.WriteByte(p.CurrentChar)
		p.advance()
	}

	// Skip the closing quote
	if p.Position < len(p.Input) && p.CurrentChar == '\'' {
		p.advance()
	} else {
		return Token{}, fmt.Errorf("unterminated string")
	}

	return Token{Type: TOKEN_STRING, Value: value.String()}, nil
}

// parseSymbol parses a symbol
func (p *Parser) parseSymbol() (Token, error) {
	// Skip the # character
	p.advance()

	// If the next character is a quote, parse a string symbol
	if p.CurrentChar == '\'' {
		token, err := p.parseString()
		if err != nil {
			return Token{}, err
		}
		return Token{Type: TOKEN_SYMBOL, Value: token.Value}, nil
	}

	// Otherwise, parse an identifier symbol
	if p.isAlpha(p.CurrentChar) {
		token := p.parseIdentifier()
		return Token{Type: TOKEN_SYMBOL, Value: token.Value}, nil
	}

	return Token{}, fmt.Errorf("invalid symbol")
}

// skipComment skips a comment
func (p *Parser) skipComment() error {
	// Skip the opening quote
	p.advance()

	for p.Position < len(p.Input) && p.CurrentChar != '"' {
		p.advance()
	}

	// Skip the closing quote
	if p.Position < len(p.Input) && p.CurrentChar == '"' {
		p.advance()
	} else {
		return fmt.Errorf("unterminated comment")
	}

	return nil
}
