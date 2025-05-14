package compiler

import (
	"encoding/binary"
	"fmt"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/bytecode"
	"smalltalklsp/interpreter/pile"
)

// BytecodeCompiler compiles AST to bytecode
type BytecodeCompiler struct {
	// Method is the method being compiled
	Method *pile.Method

	// Literals are the literals used in the method
	Literals []*pile.Object

	// Bytecodes are the bytecodes generated
	Bytecodes []byte

	// TempVarNames are the temporary variable names
	TempVarNames []string

	// Class is the class the method belongs to
	Class *pile.Object
}

// NewBytecodeCompiler creates a new bytecode compiler
func NewBytecodeCompiler(class *pile.Object) *BytecodeCompiler {
	return &BytecodeCompiler{
		Method:       nil,
		Literals:     []*pile.Object{},
		Bytecodes:    []byte{},
		TempVarNames: []string{},
		Class:        class,
	}
}

// Compile compiles an AST node to bytecode
func (c *BytecodeCompiler) Compile(node ast.Node) *pile.Method {
	// Create a new method
	c.Method = &pile.Method{
		Object: pile.Object{
			TypeField: pile.OBJ_METHOD,
		},
		Bytecodes:    []byte{},
		Literals:     []*pile.Object{},
		TempVarNames: []string{},
	}

	// Visit the node
	node.Accept(c)

	// Set the method bytecodes and literals
	c.Method.Bytecodes = c.Bytecodes
	c.Method.Literals = c.Literals
	c.Method.TempVarNames = c.TempVarNames

	// Set the method class
	c.Method.SetMethodClass(pile.ObjectToClass(c.Class))

	return c.Method
}

// VisitMethodNode visits a method node
func (c *BytecodeCompiler) VisitMethodNode(node *ast.MethodNode) interface{} {
	// Set the method selector
	c.Method.SetSelector(pile.NewSymbol(node.Selector))

	// Set the temporary variable names
	c.TempVarNames = append(c.TempVarNames, node.Parameters...)
	c.TempVarNames = append(c.TempVarNames, node.Temporaries...)
	c.Method.TempVarNames = c.TempVarNames

	// Compile the method body
	node.Body.Accept(c)

	return nil
}

// VisitReturnNode visits a return node
func (c *BytecodeCompiler) VisitReturnNode(node *ast.ReturnNode) interface{} {
	// Compile the expression
	node.Expression.Accept(c)

	// Add the return bytecode
	c.Bytecodes = append(c.Bytecodes, bytecode.RETURN_STACK_TOP)

	return nil
}

// VisitSelfNode visits a self node
func (c *BytecodeCompiler) VisitSelfNode(node *ast.SelfNode) interface{} {
	// Add the push self bytecode
	c.Bytecodes = append(c.Bytecodes, bytecode.PUSH_SELF)

	return nil
}

// VisitLiteralNode visits a literal node
func (c *BytecodeCompiler) VisitLiteralNode(node *ast.LiteralNode) interface{} {
	// Add the literal to the literals array
	literalIndex := c.addLiteral(node.Value)

	// Add the push literal bytecode
	c.Bytecodes = append(c.Bytecodes, bytecode.PUSH_LITERAL)

	// Add the literal index (4 bytes)
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, uint32(literalIndex))
	c.Bytecodes = append(c.Bytecodes, indexBytes...)

	return nil
}

// VisitVariableNode visits a variable node
func (c *BytecodeCompiler) VisitVariableNode(node *ast.VariableNode) interface{} {
	// Check if the variable is a temporary variable
	for i, name := range c.TempVarNames {
		if name == node.Name {
			// Add the push temporary variable bytecode
			c.Bytecodes = append(c.Bytecodes, bytecode.PUSH_TEMPORARY_VARIABLE)

			// Add the temporary variable index (4 bytes)
			indexBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(indexBytes, uint32(i))
			c.Bytecodes = append(c.Bytecodes, indexBytes...)

			return nil
		}
	}

	// Check if the variable is an instance variable
	// TODO: Implement instance variable lookup

	// If we get here, the variable is not found
	panic(fmt.Sprintf("Variable not found: %s", node.Name))
}

// VisitAssignmentNode visits an assignment node
func (c *BytecodeCompiler) VisitAssignmentNode(node *ast.AssignmentNode) interface{} {
	// Compile the expression
	node.Expression.Accept(c)

	// Check if the variable is a temporary variable
	for i, name := range c.TempVarNames {
		if name == node.Variable {
			// Add the store temporary variable bytecode
			c.Bytecodes = append(c.Bytecodes, bytecode.STORE_TEMPORARY_VARIABLE)

			// Add the temporary variable index (4 bytes)
			indexBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(indexBytes, uint32(i))
			c.Bytecodes = append(c.Bytecodes, indexBytes...)

			return nil
		}
	}

	// Check if the variable is an instance variable
	// TODO: Implement instance variable lookup

	// If we get here, the variable is not found
	panic(fmt.Sprintf("Variable not found: %s", node.Variable))
}

// VisitMessageSendNode visits a message send node
func (c *BytecodeCompiler) VisitMessageSendNode(node *ast.MessageSendNode) interface{} {
	// Compile the receiver
	node.Receiver.Accept(c)

	// Compile the arguments
	for _, arg := range node.Arguments {
		arg.Accept(c)
	}

	// Create a symbol and add it to the literals array
	symbol := pile.NewSymbol(node.Selector)
	selectorIndex := c.addLiteral(symbol)

	// Add the send message bytecode
	c.Bytecodes = append(c.Bytecodes, bytecode.SEND_MESSAGE)

	// Add the selector index (4 bytes)
	selectorIndexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(selectorIndexBytes, uint32(selectorIndex))
	c.Bytecodes = append(c.Bytecodes, selectorIndexBytes...)

	// Add the argument count (4 bytes)
	argCountBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(argCountBytes, uint32(len(node.Arguments)))
	c.Bytecodes = append(c.Bytecodes, argCountBytes...)

	return nil
}

// VisitBlockNode visits a block node
func (c *BytecodeCompiler) VisitBlockNode(node *ast.BlockNode) interface{} {
	// Create a new bytecode compiler for the block
	blockCompiler := NewBytecodeCompiler(c.Class)

	// Set the temporary variable names
	blockCompiler.TempVarNames = append(blockCompiler.TempVarNames, node.Parameters...)
	blockCompiler.TempVarNames = append(blockCompiler.TempVarNames, node.Temporaries...)

	// Compile the block body
	node.Body.Accept(blockCompiler)

	// Add the create block bytecode
	c.Bytecodes = append(c.Bytecodes, bytecode.CREATE_BLOCK)

	// Add the bytecode size (4 bytes)
	bytecodeSize := len(blockCompiler.Bytecodes)
	bytecodeSizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytecodeSizeBytes, uint32(bytecodeSize))
	c.Bytecodes = append(c.Bytecodes, bytecodeSizeBytes...)

	// Add the literal count (4 bytes)
	literalCount := len(blockCompiler.Literals)
	literalCountBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(literalCountBytes, uint32(literalCount))
	c.Bytecodes = append(c.Bytecodes, literalCountBytes...)

	// Add the temporary variable count (4 bytes)
	tempVarCount := len(blockCompiler.TempVarNames)
	tempVarCountBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tempVarCountBytes, uint32(tempVarCount))
	c.Bytecodes = append(c.Bytecodes, tempVarCountBytes...)

	// Add the block bytecodes
	c.Bytecodes = append(c.Bytecodes, blockCompiler.Bytecodes...)

	// Add the block literals to the method literals
	for _, literal := range blockCompiler.Literals {
		c.addLiteral(literal)
	}

	return nil
}

// addLiteral adds a literal to the literals array and returns its index
func (c *BytecodeCompiler) addLiteral(literal *pile.Object) int {
	// Check if the literal already exists
	for i, l := range c.Literals {
		if l == literal {
			return i
		}
	}

	// Add the literal
	c.Literals = append(c.Literals, literal)
	return len(c.Literals) - 1
}
