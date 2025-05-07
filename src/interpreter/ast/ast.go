package ast

import (
	"smalltalklsp/interpreter/core"
)

// Node is the interface for all AST nodes
type Node interface {
	// Accept accepts a visitor
	Accept(visitor Visitor) interface{}
}

// Visitor is the interface for visitors
type Visitor interface {
	// VisitMethodNode visits a method node
	VisitMethodNode(node *MethodNode) interface{}

	// VisitReturnNode visits a return node
	VisitReturnNode(node *ReturnNode) interface{}

	// VisitSelfNode visits a self node
	VisitSelfNode(node *SelfNode) interface{}

	// VisitLiteralNode visits a literal node
	VisitLiteralNode(node *LiteralNode) interface{}

	// VisitVariableNode visits a variable node
	VisitVariableNode(node *VariableNode) interface{}

	// VisitAssignmentNode visits an assignment node
	VisitAssignmentNode(node *AssignmentNode) interface{}

	// VisitMessageSendNode visits a message send node
	VisitMessageSendNode(node *MessageSendNode) interface{}

	// VisitBlockNode visits a block node
	VisitBlockNode(node *BlockNode) interface{}
}

// MethodNode represents a method definition
type MethodNode struct {
	// Selector is the method selector
	Selector string

	// Parameters are the method parameters
	Parameters []string

	// Temporaries are the method temporaries
	Temporaries []string

	// Body is the method body
	Body Node

	// Class is the method class
	Class *core.Object
}

// Accept implements the Node interface
func (n *MethodNode) Accept(visitor Visitor) interface{} {
	return visitor.VisitMethodNode(n)
}

// ReturnNode represents a return statement
type ReturnNode struct {
	// Expression is the expression to return
	Expression Node
}

// Accept implements the Node interface
func (n *ReturnNode) Accept(visitor Visitor) interface{} {
	return visitor.VisitReturnNode(n)
}

// SelfNode represents the self reference
type SelfNode struct{}

// Accept implements the Node interface
func (n *SelfNode) Accept(visitor Visitor) interface{} {
	return visitor.VisitSelfNode(n)
}

// LiteralNode represents a literal value
type LiteralNode struct {
	// Value is the literal value
	Value *core.Object
}

// Accept implements the Node interface
func (n *LiteralNode) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteralNode(n)
}

// VariableNode represents a variable reference
type VariableNode struct {
	// Name is the variable name
	Name string
}

// Accept implements the Node interface
func (n *VariableNode) Accept(visitor Visitor) interface{} {
	return visitor.VisitVariableNode(n)
}

// AssignmentNode represents an assignment
type AssignmentNode struct {
	// Variable is the variable to assign to
	Variable string

	// Expression is the expression to assign
	Expression Node
}

// Accept implements the Node interface
func (n *AssignmentNode) Accept(visitor Visitor) interface{} {
	return visitor.VisitAssignmentNode(n)
}

// MessageSendNode represents a message send
type MessageSendNode struct {
	// Receiver is the message receiver
	Receiver Node

	// Selector is the message selector
	Selector string

	// Arguments are the message arguments
	Arguments []Node
}

// Accept implements the Node interface
func (n *MessageSendNode) Accept(visitor Visitor) interface{} {
	return visitor.VisitMessageSendNode(n)
}

// BlockNode represents a block
type BlockNode struct {
	// Parameters are the block parameters
	Parameters []string

	// Temporaries are the block temporaries
	Temporaries []string

	// Body is the block body
	Body Node
}

// Accept implements the Node interface
func (n *BlockNode) Accept(visitor Visitor) interface{} {
	return visitor.VisitBlockNode(n)
}
