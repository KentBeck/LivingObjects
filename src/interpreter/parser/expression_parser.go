package parser

import (
	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/core"
)

// ParseExpression parses a standalone expression
func ParseExpression(input string, class *core.Object) (ast.Node, error) {
	// Wrap the expression in a method with a return statement
	wrappedInput := "^ " + input

	// Create a new parser
	p := NewParser(wrappedInput, class)

	// Tokenize the input
	err := p.tokenize()
	if err != nil {
		return nil, err
	}

	// Initialize the current token
	p.CurrentToken = p.Tokens[0]

	// Parse the method (which now contains our expression as a return statement)
	methodNode, err := p.parseMethod()
	if err != nil {
		return nil, err
	}

	// Set a default selector
	if methodNode, ok := methodNode.(*ast.MethodNode); ok {
		methodNode.Selector = "doIt"
	}

	return methodNode, nil
}
