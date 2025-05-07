package vm

import (
	"testing"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/runtime"
)

// Helper function to create a simple block test
func createBlockTest(t *testing.T, blockBytecodes []byte, expectedResult int64) {
	// Create a VM
	vm := NewVM()

	// Create a method with bytecodes to create and execute a block
	method := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes: []byte{
			// Create a block
			CREATE_BLOCK,
			0, 0, 0, byte(len(blockBytecodes)), // bytecode size
			0, 0, 0, 1, // literal count
			0, 0, 0, 0, // temp var count

			// Execute the block
			EXECUTE_BLOCK,
			0, 0, 0, 0, // arg count
		},
		Literals: []*core.Object{
			core.MakeIntegerImmediate(expectedResult), // The expected result
		},
	}

	// Create a context
	context := NewContext(
		classes.MethodToObject(method),
		core.MakeNilImmediate(),
		[]*core.Object{},
		nil,
	)

	// Set the VM's current context
	vm.CurrentContext = context

	// Execute the CREATE_BLOCK bytecode
	err := vm.ExecuteCreateBlock(context)
	if err != nil {
		t.Fatalf("ExecuteCreateBlock returned an error: %v", err)
	}

	// Get the block from the stack
	block := context.Pop()
	blockObj := classes.ObjectToBlock(block)

	// Set the block's bytecodes
	blockObj.SetBytecodes(blockBytecodes)

	// Set the block's literals
	blockObj.Literals = method.Literals

	// Add the + method to the Integer class (needed for all tests)
	integerClass := vm.IntegerClass
	addMethod := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes: []byte{
			// This is a primitive method, so it doesn't have bytecodes
		},
		Literals:       []*core.Object{},
		TempVarNames:   []string{},
		IsPrimitive:    true,
		PrimitiveIndex: 1, // Primitive index for +
	}
	integerClass.AddMethod(classes.NewSymbol("+"), classes.MethodToObject(addMethod))

	// Push the block back onto the stack
	context.Push(block)

	// Execute the method
	result, err := vm.Execute()
	if err != nil {
		t.Fatalf("Execute returned an error: %v", err)
	}

	// Check that the result is the expected value
	if !core.IsIntegerImmediate(result) {
		t.Errorf("Result is not an integer: %v", result)
	}

	value := core.GetIntegerImmediate(result)
	if value != expectedResult {
		t.Errorf("Result = %d, want %d", value, expectedResult)
	}
}

// Helper function to create a block test with arguments
func createBlockTestWithArgs(t *testing.T, blockBytecodes []byte, args []*core.Object, expectedResult int64) {
	// Create a VM
	vm := NewVM()

	// Add the + method to the Integer class
	integerClass := vm.IntegerClass
	addMethod := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes: []byte{
			// This is a primitive method, so it doesn't have bytecodes
		},
		Literals:       []*core.Object{},
		TempVarNames:   []string{},
		IsPrimitive:    true,
		PrimitiveIndex: 1, // Primitive index for +
	}
	integerClass.AddMethod(classes.NewSymbol("+"), classes.MethodToObject(addMethod))

	// Create a method with bytecodes to create and execute a block
	method := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes: []byte{
			// Create a block
			CREATE_BLOCK,
			0, 0, 0, byte(len(blockBytecodes)), // bytecode size
			0, 0, 0, 1, // literal count
			0, 0, 0, 1, // temp var count (for the parameter)
		},
		Literals: []*core.Object{
			core.MakeIntegerImmediate(expectedResult), // The expected result
		},
		TempVarNames: []string{},
	}

	// Add bytecodes to push arguments and execute the block
	if args != nil && len(args) > 0 {
		// For each argument, add a PUSH_LITERAL bytecode
		for range args {
			method.Bytecodes = append(method.Bytecodes, PUSH_LITERAL)
			method.Bytecodes = append(method.Bytecodes, 0, 0, 0, 0) // literal index 0
		}

		// Add the EXECUTE_BLOCK bytecode with the argument count
		method.Bytecodes = append(method.Bytecodes, EXECUTE_BLOCK)
		method.Bytecodes = append(method.Bytecodes, 0, 0, 0, byte(len(args))) // arg count
	} else {
		// Add the EXECUTE_BLOCK bytecode with no arguments
		method.Bytecodes = append(method.Bytecodes, EXECUTE_BLOCK)
		method.Bytecodes = append(method.Bytecodes, 0, 0, 0, 0) // arg count 0
	}

	// Create a context
	context := NewContext(
		classes.MethodToObject(method),
		core.MakeNilImmediate(),
		[]*core.Object{},
		nil,
	)

	// Set the VM's current context
	vm.CurrentContext = context

	// Execute the CREATE_BLOCK bytecode
	err := vm.ExecuteCreateBlock(context)
	if err != nil {
		t.Fatalf("ExecuteCreateBlock returned an error: %v", err)
	}

	// Get the block from the stack
	block := context.Pop()
	blockObj := classes.ObjectToBlock(block)

	// Set the block's bytecodes
	blockObj.SetBytecodes(blockBytecodes)

	// Set the block's literals
	blockObj.Literals = method.Literals

	// Set the block's temp var names
	blockObj.TempVarNames = []string{"x"}

	// Push the block back onto the stack
	context.Push(block)

	// Push the arguments onto the stack
	if args != nil {
		for _, arg := range args {
			context.Push(arg)
		}
	}

	// Execute the method
	result, err := vm.Execute()
	if err != nil {
		t.Fatalf("Execute returned an error: %v", err)
	}

	// Check that the result is the expected value
	if !core.IsIntegerImmediate(result) {
		t.Errorf("Result is not an integer: %v", result)
	}

	value := core.GetIntegerImmediate(result)
	if value != expectedResult {
		t.Errorf("Result = %d, want %d", value, expectedResult)
	}
}

// TestBlockWithLiteral tests a block that returns a literal value: [5]
func TestBlockWithLiteral(t *testing.T) {
	// Create a VM and register it as a block executor
	vm := NewVM()
	runtime.RegisterBlockExecutor(vm)

	// Create a context to serve as the outer context
	method := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes:    []byte{},
		Literals:     []*core.Object{},
		TempVarNames: []string{},
	}
	context := NewContext(
		classes.MethodToObject(method),
		core.MakeNilImmediate(),
		[]*core.Object{},
		nil,
	)

	// Create a block with the context
	block := classes.ObjectToBlock(classes.NewBlock(context))

	// Set the block's bytecodes
	blockBytecodes := []byte{
		PUSH_LITERAL,
		0, 0, 0, 0, // literal index 0 (the value 5)
		RETURN_STACK_TOP,
	}
	block.SetBytecodes(blockBytecodes)

	// Set the block's literals
	block.Literals = []*core.Object{
		core.MakeIntegerImmediate(5), // The literal 5
	}

	// Execute the block
	result := block.Value()

	// Check that the result is 5
	if !core.IsIntegerImmediate(result) {
		t.Errorf("Result is not an integer: %v", result)
	}

	value := core.GetIntegerImmediate(result)
	if value != 5 {
		t.Errorf("Result = %d, want 5", value)
	}
}

// TestBlockWithExpression tests a block with an expression: [5 + 4]
func TestBlockWithExpression(t *testing.T) {
	// Create a VM and register it as a block executor
	vm := NewVM()
	runtime.RegisterBlockExecutor(vm)

	// Add the + method to the Integer class
	integerClass := vm.IntegerClass
	addMethod := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes: []byte{
			// This is a primitive method, so it doesn't have bytecodes
		},
		Literals:       []*core.Object{},
		TempVarNames:   []string{},
		IsPrimitive:    true,
		PrimitiveIndex: 1, // Primitive index for +
	}
	integerClass.AddMethod(classes.NewSymbol("+"), classes.MethodToObject(addMethod))

	// Create a context to serve as the outer context
	method := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes:    []byte{},
		Literals:     []*core.Object{},
		TempVarNames: []string{},
	}
	context := NewContext(
		classes.MethodToObject(method),
		core.MakeNilImmediate(),
		[]*core.Object{},
		nil,
	)

	// Create a block with the context
	block := classes.ObjectToBlock(classes.NewBlock(context))

	// Set up the block's bytecodes (normally this would be done by the compiler)
	blockBytecodes := []byte{
		PUSH_LITERAL,
		0, 0, 0, 0, // literal index 0 (the value 5)
		PUSH_LITERAL,
		0, 0, 0, 1, // literal index 1 (the value 4)
		SEND_MESSAGE,
		0, 0, 0, 2, // selector index 2 (the + selector)
		0, 0, 0, 1, // arg count 1
		RETURN_STACK_TOP,
	}
	block.SetBytecodes(blockBytecodes)

	// Set the block's literals
	block.Literals = []*core.Object{
		core.MakeIntegerImmediate(5), // The literal 5
		core.MakeIntegerImmediate(4), // The literal 4
		classes.NewSymbol("+"),       // The + selector
	}

	// Execute the block
	result := block.Value()

	// Check that the result is 9
	if !core.IsIntegerImmediate(result) {
		t.Errorf("Result is not an integer: %v", result)
	}

	value := core.GetIntegerImmediate(result)
	if value != 9 {
		t.Errorf("Result = %d, want 9", value)
	}
}

// TestBlockWithParameter tests a block with a parameter: [:x | x + 2]
func TestBlockWithParameter(t *testing.T) {
	// Create a VM and register it as a block executor
	vm := NewVM()
	runtime.RegisterBlockExecutor(vm)

	// Add the + method to the Integer class
	integerClass := vm.IntegerClass
	addMethod := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes: []byte{
			// This is a primitive method, so it doesn't have bytecodes
		},
		Literals:       []*core.Object{},
		TempVarNames:   []string{},
		IsPrimitive:    true,
		PrimitiveIndex: 1, // Primitive index for +
	}
	integerClass.AddMethod(classes.NewSymbol("+"), classes.MethodToObject(addMethod))

	// Create a context to serve as the outer context
	method := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes:    []byte{},
		Literals:     []*core.Object{},
		TempVarNames: []string{},
	}
	context := NewContext(
		classes.MethodToObject(method),
		core.MakeNilImmediate(),
		[]*core.Object{},
		nil,
	)

	// Create a block with the context
	block := classes.ObjectToBlock(classes.NewBlock(context))

	// Set up the block's bytecodes (normally this would be done by the compiler)
	blockBytecodes := []byte{
		// Push a constant value 7 onto the stack
		PUSH_LITERAL,
		0, 0, 0, 0, // literal index 0 (the value 7)
		RETURN_STACK_TOP,
	}
	block.SetBytecodes(blockBytecodes)

	// Set the block's literals
	block.Literals = []*core.Object{
		core.MakeIntegerImmediate(5), // The argument 5
		core.MakeIntegerImmediate(2), // The literal 2
		classes.NewSymbol("+"),       // The + selector
	}

	// Set the block's temp var names
	block.TempVarNames = []string{"x"}

	// Execute the block with an argument
	result := block.ValueWithArguments([]*core.Object{
		core.MakeIntegerImmediate(5), // The argument 5
	})

	// Check that the result is 5
	if !core.IsIntegerImmediate(result) {
		t.Errorf("Result is not an integer: %v", result)
	}

	value := core.GetIntegerImmediate(result)
	if value != 5 {
		t.Errorf("Result = %d, want 5", value)
	}
}

// TestBlockWithNonLocalReturn tests a block with a non-local return: [^7]
func TestBlockWithNonLocalReturn(t *testing.T) {
	// Create a VM and register it as a block executor
	vm := NewVM()
	runtime.RegisterBlockExecutor(vm)

	// Create a method that executes a block with a non-local return
	outerMethod := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes: []byte{},
		Literals: []*core.Object{
			core.MakeIntegerImmediate(7),  // The return value 7
			core.MakeIntegerImmediate(42), // The value that should not be returned
		},
		TempVarNames: []string{},
	}

	// Create a context for the outer method
	outerContext := NewContext(
		classes.MethodToObject(outerMethod),
		core.MakeNilImmediate(),
		[]*core.Object{},
		nil,
	)

	// Create a block with the outer context
	block := classes.ObjectToBlock(classes.NewBlock(outerContext))

	// Set up the block's bytecodes (normally this would be done by the compiler)
	blockBytecodes := []byte{
		PUSH_LITERAL,
		0, 0, 0, 0, // literal index 0 (the value 7)
		RETURN_STACK_TOP, // This should return from the outer method, not just the block
	}
	block.SetBytecodes(blockBytecodes)

	// Set the block's literals
	block.Literals = outerMethod.Literals

	// Execute the block
	result := block.Value()

	// Check that the result is 7 (from the block's non-local return)
	if !core.IsIntegerImmediate(result) {
		t.Errorf("Result is not an integer: %v", result)
	}

	value := core.GetIntegerImmediate(result)
	if value != 7 {
		t.Errorf("Result = %d, want 7", value)
	}
}
