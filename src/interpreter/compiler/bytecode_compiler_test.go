package compiler

import (
	"testing"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/bytecode"
	"smalltalklsp/interpreter/pile"
)

// TestCompileYourself tests compiling the method Object>>yourself ^self
func TestCompileYourself(t *testing.T) {
	// Create a class
	objectClass := pile.NewClass("Object", nil)

	// Create the AST for Object>>yourself ^self
	methodNode := &ast.MethodNode{
		Selector:    "yourself",
		Parameters:  []string{},
		Temporaries: []string{},
		Body: &ast.ReturnNode{
			Expression: &ast.SelfNode{},
		},
		Class: pile.ClassToObject(objectClass),
	}

	// Create a bytecode compiler
	compiler := NewBytecodeCompiler(pile.ClassToObject(objectClass))

	// Compile the method
	method := compiler.Compile(methodNode)

	// Check the method selector
	if method.GetSelector() == nil {
		t.Errorf("Method selector is nil")
	} else {
		selectorSymbol := pile.ObjectToSymbol(method.GetSelector())
		if selectorSymbol.GetValue() != "yourself" {
			t.Errorf("Expected method selector to be 'yourself', got '%s'", selectorSymbol.GetValue())
		}
	}

	// Check the method bytecodes
	expectedBytecodes := []byte{
		bytecode.PUSH_SELF,        // Push self
		bytecode.RETURN_STACK_TOP, // Return the value on top of the stack
	}

	if len(method.Bytecodes) != len(expectedBytecodes) {
		t.Errorf("Expected bytecode length to be %d, got %d", len(expectedBytecodes), len(method.Bytecodes))
	} else {
		for i, b := range expectedBytecodes {
			if method.Bytecodes[i] != b {
				t.Errorf("Expected bytecode at index %d to be %d, got %d", i, b, method.Bytecodes[i])
			}
		}
	}

	// Check the method literals
	if len(method.Literals) != 0 {
		t.Errorf("Expected 0 literals, got %d", len(method.Literals))
	}

	// Check the method temporary variables
	if len(method.TempVarNames) != 0 {
		t.Errorf("Expected 0 temporary variables, got %d", len(method.TempVarNames))
	}

	// Check the method class
	if method.GetMethodClass() != objectClass {
		t.Errorf("Expected method class to be %v, got %v", objectClass, method.GetMethodClass())
	}
}

// TestCompileAdd tests compiling the method Integer>>+ aNumber ^self + aNumber
func TestCompileAdd(t *testing.T) {
	// Create a class
	objectClass := pile.NewClass("Object", nil)
	integerClass := pile.NewClass("Integer", objectClass)

	// Create the AST for Integer>>+ aNumber ^self + aNumber
	methodNode := &ast.MethodNode{
		Selector:    "+",
		Parameters:  []string{"aNumber"},
		Temporaries: []string{},
		Body: &ast.ReturnNode{
			Expression: &ast.MessageSendNode{
				Receiver: &ast.SelfNode{},
				Selector: "+",
				Arguments: []ast.Node{
					&ast.VariableNode{
						Name: "aNumber",
					},
				},
			},
		},
		Class: pile.ClassToObject(integerClass),
	}

	// Create a bytecode compiler
	compiler := NewBytecodeCompiler(pile.ClassToObject(integerClass))

	// Compile the method
	method := compiler.Compile(methodNode)

	// Check the method selector
	if method.GetSelector() == nil {
		t.Errorf("Method selector is nil")
	} else {
		selectorSymbol := pile.ObjectToSymbol(method.GetSelector())
		if selectorSymbol.GetValue() != "+" {
			t.Errorf("Expected method selector to be '+', got '%s'", selectorSymbol.GetValue())
		}
	}

	// Check the method bytecodes
	expectedBytecodes := []byte{
		bytecode.PUSH_SELF,               // Push self
		bytecode.PUSH_TEMPORARY_VARIABLE, // Push aNumber
		0, 0, 0, 0,                       // Temporary variable index 0
		bytecode.SEND_MESSAGE, // Send message +
		0, 0, 0, 0,            // Selector index 0
		0, 0, 0, 1, // Argument count 1
		bytecode.RETURN_STACK_TOP, // Return the value on top of the stack
	}

	if len(method.Bytecodes) != len(expectedBytecodes) {
		t.Errorf("Expected bytecode length to be %d, got %d", len(expectedBytecodes), len(method.Bytecodes))
	} else {
		for i, b := range expectedBytecodes {
			if method.Bytecodes[i] != b {
				t.Errorf("Expected bytecode at index %d to be %d, got %d", i, b, method.Bytecodes[i])
			}
		}
	}

	// Check the method literals
	if len(method.Literals) != 1 {
		t.Errorf("Expected 1 literal, got %d", len(method.Literals))
	} else {
		// Check that the literal is the + symbol
		literal := method.Literals[0]
		if literal.Type() != pile.OBJ_SYMBOL {
			t.Errorf("Expected literal to be a symbol, got %v", literal.Type())
		} else {
			symbol := pile.ObjectToSymbol(literal)
			if symbol.GetValue() != "+" {
				t.Errorf("Expected literal to be '+', got '%s'", symbol.GetValue())
			}
		}
	}

	// Check the method temporary variables
	if len(method.TempVarNames) != 1 {
		t.Errorf("Expected 1 temporary variable, got %d", len(method.TempVarNames))
	} else {
		if method.TempVarNames[0] != "aNumber" {
			t.Errorf("Expected temporary variable to be 'aNumber', got '%s'", method.TempVarNames[0])
		}
	}

	// Check the method class
	if method.GetMethodClass() != integerClass {
		t.Errorf("Expected method class to be %v, got %v", integerClass, method.GetMethodClass())
	}
}