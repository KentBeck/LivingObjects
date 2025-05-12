package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/vm"
)

// TestBooleanNotPrimitives tests the True not and False not primitives
func TestBooleanNotPrimitives(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Test True not primitive
	t.Run("True not", func(t *testing.T) {
		// Add primitive methods to the True class
		trueClass := virtualMachine.TrueClass
		notSelector := core.NewSymbol("not")
		notMethod := compiler.NewMethodBuilder(trueClass).
			Selector("not").
			Primitive(50). // True not primitive
			Go()

		// Get the true object
		trueObj := core.MakeTrueImmediate()

		// Execute the primitive
		result := virtualMachine.ExecutePrimitive(trueObj, notSelector, []*core.Object{}, notMethod)

		// Check that the result is not nil
		if result == nil {
			t.Errorf("True not primitive returned nil")
			return
		}

		// Check that the result is false
		if !core.IsFalseImmediate(result) {
			t.Errorf("Expected false result, got %v", result.Type())
			return
		}
	})

	// Test False not primitive
	t.Run("False not", func(t *testing.T) {
		// Add primitive methods to the False class
		falseClass := virtualMachine.FalseClass
		notSelector := core.NewSymbol("not")
		notMethod := compiler.NewMethodBuilder(falseClass).
			Selector("not").
			Primitive(51). // False not primitive
			Go()

		// Get the false object
		falseObj := core.MakeFalseImmediate()

		// Execute the primitive
		result := virtualMachine.ExecutePrimitive(falseObj, notSelector, []*core.Object{}, notMethod)

		// Check that the result is not nil
		if result == nil {
			t.Errorf("False not primitive returned nil")
			return
		}

		// Check that the result is true
		if !core.IsTrueImmediate(result) {
			t.Errorf("Expected true result, got %v", result.Type())
			return
		}
	})
}
