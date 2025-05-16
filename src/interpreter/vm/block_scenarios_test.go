package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/bytecode"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/runtime"
	"smalltalklsp/interpreter/vm"
)

// TestBlockWithLiteral tests a block that returns a literal value: [5]
func TestBlockWithLiteral(t *testing.T) {
	// Create a VM and register it as a block executor
	virtualMachine := vm.NewVM()
	runtime.RegisterBlockExecutor(virtualMachine)

	// Create a context to serve as the outer context
	method := &pile.Method{
		Object: pile.Object{
			TypeField: pile.OBJ_METHOD,
		},
		Bytecodes:    []byte{},
		Literals:     []*pile.Object{},
		TempVarNames: []string{},
	}
	context := vm.NewContext(
		pile.MethodToObject(method),
		pile.MakeNilImmediate(),
		[]*pile.Object{},
		nil,
	)

	// Create a block with the context
	block := pile.ObjectToBlock(virtualMachine.NewBlock(context))

	// Set the block's bytecodes
	blockBytecodes := []byte{
		bytecode.PUSH_LITERAL,
		0, 0, 0, 0, // literal index 0 (the value 5)
		bytecode.RETURN_STACK_TOP,
	}
	block.SetBytecodes(blockBytecodes)

	// Set the block's literals
	block.Literals = []*pile.Object{
		pile.MakeIntegerImmediate(5), // The literal 5
	}

	// Execute the block
	result := block.Value()

	// Check that the result is 5
	if !pile.IsIntegerImmediate(result) {
		t.Errorf("Result is not an integer: %v", result)
	}

	value := pile.GetIntegerImmediate(result)
	if value != 5 {
		t.Errorf("Result = %d, want 5", value)
	}
}

// TestBlockWithExpression tests a block with an expression: [5 + 4]
func TestBlockWithExpression(t *testing.T) {
	// Create a VM and register it as a block executor
	virtualMachine := vm.NewVM()
	runtime.RegisterBlockExecutor(virtualMachine)

	// Add the + method to the Integer class
	integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])
	addMethod := &pile.Method{
		Object: pile.Object{
			TypeField: pile.OBJ_METHOD,
		},
		Bytecodes: []byte{
			// This is a primitive method, so it doesn't have bytecodes
		},
		Literals:       []*pile.Object{},
		TempVarNames:   []string{},
		IsPrimitive:    true,
		PrimitiveIndex: 1, // Primitive index for +
	}
	pile.AddClassMethod(integerClass, pile.NewSymbol("+"), pile.MethodToObject(addMethod))

	// Create a context to serve as the outer context
	method := &pile.Method{
		Object: pile.Object{
			TypeField: pile.OBJ_METHOD,
		},
		Bytecodes:    []byte{},
		Literals:     []*pile.Object{},
		TempVarNames: []string{},
	}
	context := vm.NewContext(
		pile.MethodToObject(method),
		pile.MakeNilImmediate(),
		[]*pile.Object{},
		nil,
	)

	// Create a block with the context
	block := pile.ObjectToBlock(virtualMachine.NewBlock(context))

	// Set up the block's bytecodes (normally this would be done by the compiler)
	blockBytecodes := []byte{
		bytecode.PUSH_LITERAL,
		0, 0, 0, 0, // literal index 0 (the value 5)
		bytecode.PUSH_LITERAL,
		0, 0, 0, 1, // literal index 1 (the value 4)
		bytecode.SEND_MESSAGE,
		0, 0, 0, 2, // selector index 2 (the + selector)
		0, 0, 0, 1, // arg count 1
		bytecode.RETURN_STACK_TOP,
	}
	block.SetBytecodes(blockBytecodes)

	// Set the block's literals
	block.Literals = []*pile.Object{
		pile.MakeIntegerImmediate(5), // The literal 5
		pile.MakeIntegerImmediate(4), // The literal 4
		pile.NewSymbol("+"),          // The + selector
	}

	// Execute the block
	result := block.Value()

	// Check that the result is 9
	if !pile.IsIntegerImmediate(result) {
		t.Errorf("Result is not an integer: %v", result)
	}

	value := pile.GetIntegerImmediate(result)
	if value != 9 {
		t.Errorf("Result = %d, want 9", value)
	}
}

// TestBlockWithParameter tests a block with a parameter: [:x | x + 2]
func TestBlockWithParameter(t *testing.T) {
	// Create a VM and register it as a block executor
	virtualMachine := vm.NewVM()
	runtime.RegisterBlockExecutor(virtualMachine)

	// Add the + method to the Integer class
	integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])
	addMethod := &pile.Method{
		Object: pile.Object{
			TypeField: pile.OBJ_METHOD,
		},
		Bytecodes: []byte{
			// This is a primitive method, so it doesn't have bytecodes
		},
		Literals:       []*pile.Object{},
		TempVarNames:   []string{},
		IsPrimitive:    true,
		PrimitiveIndex: 1, // Primitive index for +
	}
	pile.AddClassMethod(integerClass, pile.NewSymbol("+"), pile.MethodToObject(addMethod))

	// Create a context to serve as the outer context
	method := &pile.Method{
		Object: pile.Object{
			TypeField: pile.OBJ_METHOD,
		},
		Bytecodes:    []byte{},
		Literals:     []*pile.Object{},
		TempVarNames: []string{},
	}
	context := vm.NewContext(
		pile.MethodToObject(method),
		pile.MakeNilImmediate(),
		[]*pile.Object{},
		nil,
	)

	// Create a block with the context
	block := pile.ObjectToBlock(virtualMachine.NewBlock(context))

	// Set up the block's bytecodes (normally this would be done by the compiler)
	// This implements [:x | x + 2]
	blockBytecodes := []byte{
		// Push the temporary variable 'x' (parameter)
		bytecode.PUSH_TEMPORARY_VARIABLE,
		0, 0, 0, 0, // temp var index 0

		// Push the literal 2
		bytecode.PUSH_LITERAL,
		0, 0, 0, 0, // literal index 0 (the value 2)

		// Send the + message
		bytecode.SEND_MESSAGE,
		0, 0, 0, 1, // selector index 1 (the + selector)
		0, 0, 0, 1, // arg count 1

		// Return the result
		bytecode.RETURN_STACK_TOP,
	}
	block.SetBytecodes(blockBytecodes)

	// Set the block's literals
	block.Literals = []*pile.Object{
		pile.MakeIntegerImmediate(2), // The literal 2
		pile.NewSymbol("+"),          // The + selector
	}

	// Set the block's temp var names
	block.TempVarNames = []string{"x"}

	// Execute the block with an argument
	result := block.ValueWithArguments([]*pile.Object{
		pile.MakeIntegerImmediate(5), // The argument 5
	})

	// Check that the result is 7 (5 + 2)
	if !pile.IsIntegerImmediate(result) {
		t.Errorf("Result is not an integer: %v", result)
	}

	value := pile.GetIntegerImmediate(result)
	if value != 7 {
		t.Errorf("Result = %d, want 7", value)
	}
}

// TestBlockWithNonLocalReturn tests a block with a non-local return: [^7]
func TestBlockWithNonLocalReturn(t *testing.T) {
	// Create a VM and register it as a block executor
	virtualMachine := vm.NewVM()
	runtime.RegisterBlockExecutor(virtualMachine)

	// Create a method that executes a block with a non-local return
	outerMethod := &pile.Method{
		Object: pile.Object{
			TypeField: pile.OBJ_METHOD,
		},
		Bytecodes: []byte{},
		Literals: []*pile.Object{
			pile.MakeIntegerImmediate(7),  // The return value 7
			pile.MakeIntegerImmediate(42), // The value that should not be returned
		},
		TempVarNames: []string{},
	}

	// Create a context for the outer method
	outerContext := vm.NewContext(
		pile.MethodToObject(outerMethod),
		pile.MakeNilImmediate(),
		[]*pile.Object{},
		nil,
	)

	// Create a block with the outer context
	block := pile.ObjectToBlock(virtualMachine.NewBlock(outerContext))

	// Set up the block's bytecodes (normally this would be done by the compiler)
	blockBytecodes := []byte{
		bytecode.PUSH_LITERAL,
		0, 0, 0, 0, // literal index 0 (the value 7)
		bytecode.RETURN_STACK_TOP, // This should return from the outer method, not just the block
	}
	block.SetBytecodes(blockBytecodes)

	// Set the block's literals
	block.Literals = outerMethod.Literals

	// Execute the block
	result := block.Value()

	// Check that the result is 7 (from the block's non-local return)
	if !pile.IsIntegerImmediate(result) {
		t.Errorf("Result is not an integer: %v", result)
	}

	value := pile.GetIntegerImmediate(result)
	if value != 7 {
		t.Errorf("Result = %d, want 7", value)
	}
}
