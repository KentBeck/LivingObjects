package vm_test

import (
	"smalltalklsp/interpreter/pile"
	"testing"

	"smalltalklsp/interpreter/vm"
)

// TestBooleanNotMethods tests the True not and False not methods
func TestBooleanNotMethods(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Test True not method
	t.Run("True not", func(t *testing.T) {
		// Get the true object
		trueObj := pile.MakeTrueImmediate()

		// Get the not selector
		notSelector := pile.NewSymbol("not")

		// Look up the method
		method := virtualMachine.LookupMethod(trueObj, notSelector)
		if method == nil {
			t.Fatalf("Could not find not method for True class")
		}

		// Create a context for executing the method
		context := vm.NewContext(method, trueObj, []*pile.Object{}, nil)

		// Execute the method
		result, err := virtualMachine.ExecuteContext(context)
		if err != nil {
			t.Fatalf("Error executing True not method: %v", err)
		}

		// Check that the result is not nil
		if result == nil {
			t.Errorf("True not method returned nil")
			return
		}

		// Check that the result is false
		if !pile.IsFalseImmediate(result.(*pile.Object)) {
			t.Errorf("Expected false result, got %v", result.(*pile.Object).Type())
			return
		}
	})

	// Test False not method
	t.Run("False not", func(t *testing.T) {
		// Get the false object
		falseObj := pile.MakeFalseImmediate()

		// Get the not selector
		notSelector := pile.NewSymbol("not")

		// Look up the method
		method := virtualMachine.LookupMethod(falseObj, notSelector)
		if method == nil {
			t.Fatalf("Could not find not method for False class")
		}

		// Create a context for executing the method
		context := vm.NewContext(method, falseObj, []*pile.Object{}, nil)

		// Execute the method
		result, err := virtualMachine.ExecuteContext(context)
		if err != nil {
			t.Fatalf("Error executing False not method: %v", err)
		}

		// Check that the result is not nil
		if result == nil {
			t.Errorf("False not method returned nil")
			return
		}

		// Check that the result is true
		if !pile.IsTrueImmediate(result.(*pile.Object)) {
			t.Errorf("Expected true result, got %v", result.(*pile.Object).Type())
			return
		}
	})
}