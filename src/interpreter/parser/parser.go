package parser

import (
	"fmt"
	"strings"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/pile"
)

// Parser parses Smalltalk code into an AST
type Parser struct {
	// Input is the input string to parse
	Input string

	// Class is the class the method belongs to
	Class *pile.Object

	// VM is the virtual machine used for creating literals and accessing globals
	VM interface {
		NewInteger(value int64) *pile.Object
		NewString(value string) *pile.Object
		NewArray(size int) *pile.Object
		GetGlobal(name string) *pile.Object
	}

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
func NewParser(input string, class *pile.Object, vm interface {
	NewInteger(value int64) *pile.Object
	NewString(value string) *pile.Object
	NewArray(size int) *pile.Object
	GetGlobal(name string) *pile.Object
}) *Parser {
	p := &Parser{
		Input:             input,
		Class:             class,
		VM:                vm,
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

// ParseExpression parses the input and returns an AST
func (p *Parser) ParseExpression() (ast.Node, error) {
	// Tokenize the input
	err := p.tokenize()
	if err != nil {
		return nil, err
	}

	// Initialize the current token
	p.CurrentToken = p.Tokens[0]
	p.CurrentTokenIndex = 0

	// Check if the input starts with a return statement
	if p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == "^" {
		// Skip the return token
		p.advanceToken()

		// Parse the expression
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		// Create a return node
		return &ast.ReturnNode{
			Expression: expr,
		}, nil
	}

	// Parse the expression
	return p.parseKeywordMessage()
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

		// Initialize the current token index if needed
		if p.CurrentTokenIndex >= len(p.Tokens) {
			return nil, fmt.Errorf("unexpected end of input after return token")
		}

		// Parse the expression
		expression, err := p.parseKeywordMessage()
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
	// In Smalltalk, expressions are evaluated with the following precedence:
	// 1. Parenthesized expressions
	// 2. Unary messages (e.g., obj size)
	// 3. Binary messages (e.g., a + b)
	// 4. Keyword messages (e.g., dict at: key put: value)

	// Start with a primary expression
	return p.parseKeywordMessage()
}

// parseKeywordMessage parses a keyword message (lowest precedence)
func (p *Parser) parseKeywordMessage() (ast.Node, error) {
	// First parse a binary expression
	receiver, err := p.parseBinaryMessage()
	if err != nil {
		return nil, err
	}

	// Check if there's a keyword message
	if p.CurrentToken.Type == TOKEN_IDENTIFIER && strings.HasSuffix(p.CurrentToken.Value, ":") {
		// Collect all keyword parts and arguments
		var keywordParts []string
		var arguments []ast.Node

		for p.CurrentToken.Type == TOKEN_IDENTIFIER && strings.HasSuffix(p.CurrentToken.Value, ":") {
			// Add the keyword part
			keywordParts = append(keywordParts, p.CurrentToken.Value)
			p.advanceToken()

			// Parse the argument (which can be any expression except a keyword message)
			arg, err := p.parseBinaryMessage()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
		}

		// Combine the keyword parts to form the selector
		selector := strings.Join(keywordParts, "")

		return &ast.MessageSendNode{
			Receiver:  receiver,
			Selector:  selector,
			Arguments: arguments,
		}, nil
	}

	return receiver, nil
}

// parseBinaryMessage parses a binary message (medium precedence)
func (p *Parser) parseBinaryMessage() (ast.Node, error) {
	// First parse a unary message
	left, err := p.parseUnaryMessage()
	if err != nil {
		return nil, err
	}

	// Parse a chain of binary messages
	// Binary operators are special characters like +, -, *, /, <, >, etc.
	// But NOT ), ], or other non-binary operators
	for p.CurrentToken.Type == TOKEN_SPECIAL && 
		p.CurrentToken.Value != ")" && 
		p.CurrentToken.Value != "]" && 
		p.CurrentToken.Value != "." {
		
		// Get the binary selector
		selector := p.CurrentToken.Value
		p.advanceToken()

		// Parse the right operand (which can be any expression except a binary or keyword message)
		right, err := p.parseUnaryMessage()
		if err != nil {
			return nil, err
		}

		// Create a message send node
		left = &ast.MessageSendNode{
			Receiver:  left,
			Selector:  selector,
			Arguments: []ast.Node{right},
		}
	}

	return left, nil
}

// parseUnaryMessage parses a unary message (highest precedence)
func (p *Parser) parseUnaryMessage() (ast.Node, error) {
	// First parse a primary expression
	receiver, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	// Parse a chain of unary messages
	for p.CurrentToken.Type == TOKEN_IDENTIFIER && !strings.HasSuffix(p.CurrentToken.Value, ":") {
		// Get the unary selector
		selector := p.CurrentToken.Value
		p.advanceToken()

		// Create a message send node
		receiver = &ast.MessageSendNode{
			Receiver:  receiver,
			Selector:  selector,
			Arguments: []ast.Node{},
		}
	}

	return receiver, nil
}

// parsePrimary parses a primary expression
func (p *Parser) parsePrimary() (ast.Node, error) {
	// Handle self
	if p.CurrentToken.Type == TOKEN_IDENTIFIER && p.CurrentToken.Value == "self" {
		p.advanceToken()
		return &ast.SelfNode{}, nil
	}

	// Handle true and false
	if p.CurrentToken.Type == TOKEN_IDENTIFIER && p.CurrentToken.Value == "true" {
		p.advanceToken()

		// For parser tests, we'll use immediate values since they're safer for test comparison
		// and we're not keeping them across GC cycles in tests
		return &ast.LiteralNode{
			Value: pile.MakeTrueImmediate(),
		}, nil
	}

	if p.CurrentToken.Type == TOKEN_IDENTIFIER && p.CurrentToken.Value == "false" {
		p.advanceToken()

		// For parser tests, we'll use immediate values since they're safer for test comparison
		// and we're not keeping them across GC cycles in tests
		return &ast.LiteralNode{
			Value: pile.MakeFalseImmediate(),
		}, nil
	}

	// Handle string literals
	if p.CurrentToken.Type == TOKEN_STRING {
		// Create a string literal node using the VM
		literalNode := &ast.LiteralNode{
			Value: p.VM.NewString(p.CurrentToken.Value),
		}
		p.advanceToken()
		return literalNode, nil
	}

	// Handle number literals
	if p.CurrentToken.Type == TOKEN_NUMBER {
		// Parse the number
		var value int64
		fmt.Sscanf(p.CurrentToken.Value, "%d", &value)

		// Create a number literal node using the VM
		literalNode := &ast.LiteralNode{
			Value: p.VM.NewInteger(value),
		}
		p.advanceToken()
		return literalNode, nil
	}

	// Handle block expressions
	if p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == "[" {
		return p.parseBlock()
	}

	// Handle array literals - must check this before parenthesized expressions
	if p.CurrentToken.Type == TOKEN_SYMBOL && p.CurrentToken.Value == "(" {
		return p.parseArrayLiteral()
	}

	// Handle parenthesized expressions
	if p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == "(" {
		p.advanceToken() // Skip the opening parenthesis

		// Parse the expression inside the parentheses
		expr, err := p.parseKeywordMessage()
		if err != nil {
			return nil, err
		}

		// Expect a closing parenthesis
		if p.CurrentToken.Type != TOKEN_SPECIAL || p.CurrentToken.Value != ")" {
			return nil, fmt.Errorf("expected closing parenthesis, got %v", p.CurrentToken)
		}
		p.advanceToken() // Skip the closing parenthesis

		return expr, nil
	}

	// Handle variables and globals
	if p.CurrentToken.Type == TOKEN_IDENTIFIER {
		name := p.CurrentToken.Value
		p.advanceToken()

		// Check if this is a global variable (starting with uppercase)
		if len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z' {
			// Look up the global in the VM
			globalObj := p.VM.GetGlobal(name)

			// If it's a class or other global, return it as a literal node
			// TODO look it up at runtime since the value may have changed
			return &ast.LiteralNode{Value: globalObj}, nil
		}

		// Otherwise, treat it as a regular variable
		return &ast.VariableNode{Name: name}, nil
	}

	return nil, fmt.Errorf("expected primary expression, got %v", p.CurrentToken)
}

// parseArrayLiteral parses an array literal like #(1 2 3)
func (p *Parser) parseArrayLiteral() (ast.Node, error) {
	// Skip the opening symbol token (the # has already been handled by the tokenizer)
	p.advanceToken()

	// Check if the next token is the opening parenthesis
	if p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == "(" {
		// Skip the opening parenthesis
		p.advanceToken()
	} else {
		return nil, fmt.Errorf("expected opening parenthesis for array literal, got %v", p.CurrentToken)
	}

	// Parse the array elements
	var elements []ast.Node

	// Continue parsing elements until we reach the closing parenthesis
	for p.CurrentToken.Type != TOKEN_SPECIAL || p.CurrentToken.Value != ")" {
		// Parse the element (only literals are allowed in array literals)
		if p.CurrentToken.Type == TOKEN_NUMBER {
			// Parse number literal
			var value int64
			fmt.Sscanf(p.CurrentToken.Value, "%d", &value)

			element := &ast.LiteralNode{
				Value: p.VM.NewInteger(value),
			}
			elements = append(elements, element)
			p.advanceToken()
		} else if p.CurrentToken.Type == TOKEN_STRING {
			// Parse string literal
			element := &ast.LiteralNode{
				Value: p.VM.NewString(p.CurrentToken.Value),
			}
			elements = append(elements, element)
			p.advanceToken()
		} else if p.CurrentToken.Type == TOKEN_IDENTIFIER &&
			(p.CurrentToken.Value == "true" || p.CurrentToken.Value == "false") {
			// Parse boolean literal
			var value *pile.Object
			if p.CurrentToken.Value == "true" {
				value = pile.MakeTrueImmediate()
			} else {
				value = pile.MakeFalseImmediate()
			}
			element := &ast.LiteralNode{
				Value: value,
			}
			elements = append(elements, element)
			p.advanceToken()
		} else {
			return nil, fmt.Errorf("unexpected token in array literal: %v", p.CurrentToken)
		}

		// If we've reached the end of the array, break
		if p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == ")" {
			break
		}
	}

	// Expect a closing parenthesis
	if p.CurrentToken.Type != TOKEN_SPECIAL || p.CurrentToken.Value != ")" {
		return nil, fmt.Errorf("expected closing parenthesis for array literal, got %v", p.CurrentToken)
	}
	p.advanceToken() // Skip the closing parenthesis

	// Create an actual Array object using the VM
	arrayObj := p.VM.NewArray(len(elements))
	array := pile.ObjectToArray(arrayObj)

	// Fill the array with the parsed elements
	for i, element := range elements {
		// We need to extract the actual value from each element node
		if literalNode, ok := element.(*ast.LiteralNode); ok {
			array.AtPut(i, literalNode.Value)
		} else {
			return nil, fmt.Errorf("expected literal node for array element, got %T", element)
		}
	}

	// Create a literal node with the array object
	return &ast.LiteralNode{
		Value: arrayObj,
	}, nil
}

// parseBlock parses a block expression
func (p *Parser) parseBlock() (ast.Node, error) {
	// Skip the opening bracket
	p.advanceToken()

	// Parse block parameters (if any)
	var parameters []string

	// Check for the block parameter pattern, which starts with a colon
	if p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == ":" {
		// Skip the colon
		p.advanceToken()

		// Expect an identifier (the parameter name)
		if p.CurrentToken.Type != TOKEN_IDENTIFIER {
			return nil, fmt.Errorf("expected identifier after : in block parameter, got %v", p.CurrentToken)
		}

		// Add the parameter name
		parameters = append(parameters, p.CurrentToken.Value)
		p.advanceToken() // Skip the parameter name

		// Check for additional parameters (would be another : token)
		for p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == ":" {
			// Skip the colon
			p.advanceToken()

			// Expect an identifier
			if p.CurrentToken.Type != TOKEN_IDENTIFIER {
				return nil, fmt.Errorf("expected identifier after : in block parameter, got %v", p.CurrentToken)
			}

			// Add the parameter name
			parameters = append(parameters, p.CurrentToken.Value)
			p.advanceToken() // Skip the parameter name
		}

		// After parameters, expect a | token
		if p.CurrentToken.Type != TOKEN_SPECIAL || p.CurrentToken.Value != "|" {
			return nil, fmt.Errorf("expected | after block parameters, got %v", p.CurrentToken)
		}

		// Skip the | token
		p.advanceToken()
	} else {
		// No parameters, check for temporaries
		if p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == "|" {
			// Has temporaries
			p.advanceToken() // Skip the opening |
		}
	}

	// Parse temporary variables if they exist
	var temporaries []string
	
	// Parse temporaries if we didn't just parse parameters (in which case we've already consumed the |)
	if len(parameters) == 0 && p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == "|" {
		// Skip the opening bar for temporaries
		p.advanceToken()

		// Parse temporary variable names
		for p.CurrentToken.Type == TOKEN_IDENTIFIER {
			temporaries = append(temporaries, p.CurrentToken.Value)
			p.advanceToken()
		}

		// Expect a closing bar
		if p.CurrentToken.Type != TOKEN_SPECIAL || p.CurrentToken.Value != "|" {
			return nil, fmt.Errorf("expected closing | for block temporaries, got %v", p.CurrentToken)
		}

		// Skip the closing bar
		p.advanceToken()
	} else if len(parameters) == 0 {
		// No parameters and no temporaries, do nothing
	}

	// Parse the block body
	// We'll handle a sequence of expressions separated by periods
	var bodyExpressions []ast.Node

	// Parse expressions until we reach the closing bracket or EOF
	for p.CurrentToken.Type != TOKEN_EOF && (p.CurrentToken.Type != TOKEN_SPECIAL || p.CurrentToken.Value != "]") {
		// Parse an expression
		expr, err := p.parseKeywordMessage()
		if err != nil {
			return nil, err
		}

		bodyExpressions = append(bodyExpressions, expr)

		// If the next token is a period, skip it and continue parsing
		if p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == "." {
			p.advanceToken()
		} else if p.CurrentToken.Type == TOKEN_EOF || (p.CurrentToken.Type != TOKEN_SPECIAL || p.CurrentToken.Value != "]") {
			// If it's not a period or closing bracket (and not EOF), it's a syntax error
			if p.CurrentToken.Type != TOKEN_EOF {
				return nil, fmt.Errorf("expected period or closing bracket in block, got %v", p.CurrentToken)
			}
			// If it's EOF, we'll just end the block
			break
		}
	}

	// Skip the closing bracket
	if p.CurrentToken.Type == TOKEN_SPECIAL && p.CurrentToken.Value == "]" {
		p.advanceToken()
	}

	// If the block has no expressions, return a nil block
	if len(bodyExpressions) == 0 {
		return &ast.BlockNode{
			Parameters:  parameters,
			Temporaries: temporaries,
			Body:        &ast.LiteralNode{Value: pile.MakeNilImmediate()},
		}, nil
	}

	// Use the last expression as the block's body (for now, ignoring statements before the last one)
	body := bodyExpressions[len(bodyExpressions)-1]

	// Create the block node
	blockNode := &ast.BlockNode{
		Parameters:  parameters,
		Temporaries: temporaries,
		Body:        body,
	}

	return blockNode, nil
}

// We don't need parseMessageSend anymore as it's been replaced by the more specific
// parseUnaryMessage, parseBinaryMessage, and parseKeywordMessage methods

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
	return strings.ContainsRune("+-*/=<>[](){}^.|:,~", rune(c))
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

	// If the next character is an opening parenthesis, it's an array literal
	if p.CurrentChar == '(' {
		// Return a special token for array literals
		return Token{Type: TOKEN_SYMBOL, Value: "("}, nil
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
