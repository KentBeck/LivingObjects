package vm_test

import (
	"smalltalklsp/interpreter/pile"
	"testing"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/vm"
)

func TestNilClassPanic(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Create a method with the basicClass primitive using MethodBuilder
	basicClassSelector := pile.NewSymbol("basicClass")

	// Create the method using MethodBuilder
	compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"])).
		Primitive(5). // basicClass primitive
		Go("basicClass")

	// Create an object with a nil class
	objWithNilClass := &pile.Object{
		TypeField: pile.OBJ_INSTANCE,
		// ClassField is nil by default
	}

	// Create a test method that will send the basicClass message
	builder := compiler.NewMethodBuilder(pile.ObjectToClass(virtualMachine.Globals["Object"]))
	selectorIndex, builder := builder.AddLiteral(basicClassSelector)

	// Create bytecodes for the test method
	builder.PushSelf()
	builder.SendMessage(selectorIndex, 0)
	builder.ReturnStackTop()

	testMethod := builder.Go("test")

	// Create a context for the test method
	context := vm.NewContext(testMethod, objWithNilClass, []*pile.Object{}, nil)

	// Execute the test method and expect a panic
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Expected a panic when accessing basicClass on an object with nil class, but no panic occurred")
		} else {
			t.Logf("Got expected panic: %v", r)
		}
	}()

	// This should panic
	virtualMachine.ExecuteContext(context)
}