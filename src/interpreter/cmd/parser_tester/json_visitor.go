package main

import (
	"fmt"
	"strings"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/pile"
)

// JSONVisitor converts an AST to a JSON string
type JSONVisitor struct{}

// VisitMethodNode visits a method node
func (v *JSONVisitor) VisitMethodNode(node *ast.MethodNode) interface{} {
	bodyJSON := ""
	if node.Body != nil {
		bodyJSON = node.Body.Accept(v).(string)
	}

	// Convert array strings to JSON array format
	paramsJSON := formatStringArray(node.Parameters)
	tempsJSON := formatStringArray(node.Temporaries)

	return fmt.Sprintf(`{
  "type": "MethodNode",
  "selector": "%s",
  "parameters": %s,
  "temporaries": %s,
  "body": %s
}`, node.Selector, paramsJSON, tempsJSON, bodyJSON)
}

// VisitReturnNode visits a return node
func (v *JSONVisitor) VisitReturnNode(node *ast.ReturnNode) interface{} {
	exprJSON := "null"
	if node.Expression != nil {
		exprJSON = node.Expression.Accept(v).(string)
	}

	return fmt.Sprintf(`{
  "type": "ReturnNode",
  "expression": %s
}`, exprJSON)
}

// VisitSelfNode visits a self node
func (v *JSONVisitor) VisitSelfNode(node *ast.SelfNode) interface{} {
	return `{
  "type": "SelfNode"
}`
}

// VisitLiteralNode visits a literal node
func (v *JSONVisitor) VisitLiteralNode(node *ast.LiteralNode) interface{} {
	literalJSON := "null"
	if node.Value != nil {
		// Try to convert the literal value based on its type
		if pile.IsIntegerImmediate(node.Value) {
			literalJSON = fmt.Sprintf(`{"type": "Integer", "value": %d}`, 
				pile.GetIntegerImmediate(node.Value))
		} else if pile.IsTrueImmediate(node.Value) {
			literalJSON = `{"type": "Boolean", "value": true}`
		} else if pile.IsFalseImmediate(node.Value) {
			literalJSON = `{"type": "Boolean", "value": false}`
		} else if pile.IsNilImmediate(node.Value) {
			literalJSON = `{"type": "Nil"}`
		} else if pile.IsFloatImmediate(node.Value) {
			literalJSON = fmt.Sprintf(`{"type": "Float", "value": %f}`, 
				pile.GetFloatImmediate(node.Value))
		} else if node.Value.Type() == pile.OBJ_STRING {
			str := pile.ObjectToString(node.Value)
			literalJSON = fmt.Sprintf(`{"type": "String", "value": "%s"}`, escapeString(str.GetValue()))
		} else if node.Value.Type() == pile.OBJ_SYMBOL {
			sym := pile.ObjectToSymbol(node.Value)
			literalJSON = fmt.Sprintf(`{"type": "Symbol", "value": "%s"}`, escapeString(sym.GetValue()))
		} else {
			// For unknown types, just use a generic description
			literalJSON = fmt.Sprintf(`{"type": "Object", "objectType": %d}`, node.Value.Type())
		}
	}

	return fmt.Sprintf(`{
  "type": "LiteralNode",
  "value": %s
}`, literalJSON)
}

// VisitVariableNode visits a variable node
func (v *JSONVisitor) VisitVariableNode(node *ast.VariableNode) interface{} {
	return fmt.Sprintf(`{
  "type": "VariableNode",
  "name": "%s"
}`, node.Name)
}

// VisitAssignmentNode visits an assignment node
func (v *JSONVisitor) VisitAssignmentNode(node *ast.AssignmentNode) interface{} {
	exprJSON := "null"
	if node.Expression != nil {
		exprJSON = node.Expression.Accept(v).(string)
	}

	return fmt.Sprintf(`{
  "type": "AssignmentNode",
  "variable": "%s",
  "expression": %s
}`, node.Variable, exprJSON)
}

// VisitMessageSendNode visits a message send node
func (v *JSONVisitor) VisitMessageSendNode(node *ast.MessageSendNode) interface{} {
	receiverJSON := "null"
	if node.Receiver != nil {
		receiverJSON = node.Receiver.Accept(v).(string)
	}

	argsJSON := "[]"
	if len(node.Arguments) > 0 {
		args := make([]string, len(node.Arguments))
		for i, arg := range node.Arguments {
			if arg != nil {
				args[i] = arg.Accept(v).(string)
			} else {
				args[i] = "null"
			}
		}
		argsJSON = fmt.Sprintf("[\n    %s\n  ]", strings.Join(args, ",\n    "))
	}

	return fmt.Sprintf(`{
  "type": "MessageSendNode",
  "receiver": %s,
  "selector": "%s",
  "arguments": %s
}`, receiverJSON, node.Selector, argsJSON)
}

// VisitBlockNode visits a block node
func (v *JSONVisitor) VisitBlockNode(node *ast.BlockNode) interface{} {
	bodyJSON := "null"
	if node.Body != nil {
		bodyJSON = node.Body.Accept(v).(string)
	}

	// Convert array strings to JSON array format
	paramsJSON := formatStringArray(node.Parameters)
	tempsJSON := formatStringArray(node.Temporaries)

	return fmt.Sprintf(`{
  "type": "BlockNode",
  "parameters": %s,
  "temporaries": %s,
  "body": %s
}`, paramsJSON, tempsJSON, bodyJSON)
}

// Helper functions

// formatStringArray formats a string array as a JSON array
func formatStringArray(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}

	quoted := make([]string, len(arr))
	for i, s := range arr {
		quoted[i] = fmt.Sprintf(`"%s"`, s)
	}

	return fmt.Sprintf("[%s]", strings.Join(quoted, ", "))
}

// escapeString escapes special characters in a string for JSON
func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}