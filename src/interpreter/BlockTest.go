package main

import (
	"testing"
)

func TestBlockCreation(t *testing.T) {
	vm := NewVM()

	// Create a block class
	blockClass := NewClass("Block", vm.ObjectClass)

	// Create a method that creates and returns a block
	builder := NewMethodBuilder(vm.ObjectClass).Selector("createBlock")
	
	// Add a literal for the block class
	blockClassIndex, builder := builder.AddLiteral(blockClass)
	
	// Create bytecodes for the method:
	// 1. Create a new block instance
	builder.PushLiteral(blockClassIndex)
	
	// 2. Return the block
	builder.ReturnStackTop()
	
	// Finalize the method
	createBlockMethod := builder.Go()

	// Create a context for the method
	context := NewContext(createBlockMethod, vm.ObjectClass, []*Object{}, nil)

	// Execute the context
	result, err := vm.ExecuteContext(context)
	if err != nil {
		t.Errorf("Error executing context: %v", err)
	}

	// Check that the result is a block
	if result.Class != blockClass {
		t.Errorf("Expected result to be a block, got %v", result)
	}
}

func TestBlockExecution(t *testing.T) {
	vm := NewVM()

	// Create a block class
	blockClass := NewClass("Block", vm.ObjectClass)

	// Add a 'value' method to the block class that returns 42
	valueBuilder := NewMethodBuilder(blockClass).Selector("value")
	
	// Add a literal for the return value
	fortyTwoIndex, valueBuilder := valueBuilder.AddLiteral(vm.NewInteger(42))
	
	// Create bytecodes for the method:
	// 1. Push the return value
	valueBuilder.PushLiteral(fortyTwoIndex)
	
	// 2. Return the value
	valueBuilder.ReturnStackTop()
	
	// Finalize the method
	valueBuilder.Go()

	// Create a method that creates a block and sends it the 'value' message
	builder := NewMethodBuilder(vm.ObjectClass).Selector("executeBlock")
	
	// Add literals for the block class and the 'value' selector
	blockClassIndex, builder := builder.AddLiteral(blockClass)
	valueSelector, builder := builder.AddLiteral(NewSymbol("value"))
	
	// Create bytecodes for the method:
	// 1. Create a new block instance
	builder.PushLiteral(blockClassIndex)
	
	// 2. Send the 'value' message to the block
	builder.SendMessage(valueSelector, 0)
	
	// 3. Return the result
	builder.ReturnStackTop()
	
	// Finalize the method
	executeBlockMethod := builder.Go()

	// Create a context for the method
	context := NewContext(executeBlockMethod, vm.ObjectClass, []*Object{}, nil)

	// Execute the context
	result, err := vm.ExecuteContext(context)
	if err != nil {
		t.Errorf("Error executing context: %v", err)
	}

	// Check that the result is 42
	if IsIntegerImmediate(result) {
		intValue := GetIntegerImmediate(result)
		if intValue != 42 {
			t.Errorf("Expected result to be 42, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}
}

func TestBlockWithArguments(t *testing.T) {
	vm := NewVM()

	// Create a block class
	blockClass := NewClass("Block", vm.ObjectClass)

	// Add a 'value:' method to the block class that returns its argument
	valueWithArgBuilder := NewMethodBuilder(blockClass).Selector("value:")
	
	// Create bytecodes for the method:
	// 1. Push the first argument
	valueWithArgBuilder.PushTemporaryVariable(0)
	
	// 2. Return the value
	valueWithArgBuilder.ReturnStackTop()
	
	// Finalize the method
	valueWithArgBuilder.Go()

	// Create a method that creates a block and sends it the 'value:' message with 42 as the argument
	builder := NewMethodBuilder(vm.ObjectClass).Selector("executeBlockWithArg")
	
	// Add literals for the block class, the 'value:' selector, and the argument
	blockClassIndex, builder := builder.AddLiteral(blockClass)
	valueWithArgSelector, builder := builder.AddLiteral(NewSymbol("value:"))
	fortyTwoIndex, builder := builder.AddLiteral(vm.NewInteger(42))
	
	// Create bytecodes for the method:
	// 1. Create a new block instance
	builder.PushLiteral(blockClassIndex)
	
	// 2. Push the argument
	builder.PushLiteral(fortyTwoIndex)
	
	// 3. Send the 'value:' message to the block with the argument
	builder.SendMessage(valueWithArgSelector, 1)
	
	// 4. Return the result
	builder.ReturnStackTop()
	
	// Finalize the method
	executeBlockWithArgMethod := builder.Go()

	// Create a context for the method
	context := NewContext(executeBlockWithArgMethod, vm.ObjectClass, []*Object{}, nil)

	// Execute the context
	result, err := vm.ExecuteContext(context)
	if err != nil {
		t.Errorf("Error executing context: %v", err)
	}

	// Check that the result is 42
	if IsIntegerImmediate(result) {
		intValue := GetIntegerImmediate(result)
		if intValue != 42 {
			t.Errorf("Expected result to be 42, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}
}

func TestBlockWithClosure(t *testing.T) {
	vm := NewVM()

	// Create a block class
	blockClass := NewClass("Block", vm.ObjectClass)

	// Add instance variables to the block class for the closure
	blockClass.InstanceVarNames = append(blockClass.InstanceVarNames, "outerContext")

	// Add a 'value' method to the block class that accesses a variable from the outer context
	valueBuilder := NewMethodBuilder(blockClass).Selector("value")
	
	// Add a literal for the 'outerContext' instance variable name
	outerContextIndex, valueBuilder := valueBuilder.AddLiteral(NewSymbol("outerContext"))
	
	// Create bytecodes for the method:
	// 1. Push self (the block)
	valueBuilder.PushSelf()
	
	// 2. Push the outerContext instance variable
	valueBuilder.PushInstanceVariable(0)
	
	// 3. Return the value
	valueBuilder.ReturnStackTop()
	
	// Finalize the method
	valueBuilder.Go()

	// Create a method that creates a block with a closure and sends it the 'value' message
	builder := NewMethodBuilder(vm.ObjectClass).Selector("executeBlockWithClosure")
	
	// Add literals for the block class, the 'value' selector, and the closure value
	blockClassIndex, builder := builder.AddLiteral(blockClass)
	valueSelector, builder := builder.AddLiteral(NewSymbol("value"))
	closureValueIndex, builder := builder.AddLiteral(vm.NewInteger(42))
	
	// Create bytecodes for the method:
	// 1. Create a new block instance
	builder.PushLiteral(blockClassIndex)
	
	// 2. Store the closure value in the block's outerContext instance variable
	builder.Duplicate() // Duplicate the block for later use
	builder.PushLiteral(closureValueIndex)
	builder.StoreInstanceVariable(0) // Store in the outerContext instance variable
	
	// 3. Send the 'value' message to the block
	builder.SendMessage(valueSelector, 0)
	
	// 4. Return the result
	builder.ReturnStackTop()
	
	// Finalize the method
	executeBlockWithClosureMethod := builder.Go()

	// Create a context for the method
	context := NewContext(executeBlockWithClosureMethod, vm.ObjectClass, []*Object{}, nil)

	// Execute the context
	result, err := vm.ExecuteContext(context)
	if err != nil {
		t.Errorf("Error executing context: %v", err)
	}

	// Check that the result is 42
	if IsIntegerImmediate(result) {
		intValue := GetIntegerImmediate(result)
		if intValue != 42 {
			t.Errorf("Expected result to be 42, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}
}
