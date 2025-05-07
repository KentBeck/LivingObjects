package testing

import (
	"fmt"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/parser"
)

// ParseExpression parses a Smalltalk expression and wraps it in a method
func ParseExpression(input string, class *core.Object) (ast.Node, error) {
	// Wrap the expression in a method
	methodSource := fmt.Sprintf("evaluate\n^%s", input)

	// Create a parser
	p := parser.NewParser(methodSource, class)

	// Parse the method
	return p.Parse()
}

// WrapExpressionInMethod wraps an expression in a method for execution
func WrapExpressionInMethod(expression string) string {
	return fmt.Sprintf("evaluate\n^%s", expression)
}
