package main

import (
	"fmt"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/vm"
)

func main() {
	// Create a VM
	virtualMachine := vm.NewVM()

	// Test string literals
	fmt.Println("Testing string literals...")
	str1 := classes.NewString("hello")
	fmt.Printf("String 1: %s\n", str1.GetValue())

	// Test string concatenation
	fmt.Println("\nTesting string concatenation...")
	str2 := classes.NewString(" world")
	result := str1.Concat(str2)
	fmt.Printf("Concatenated: %s\n", result.GetValue())

	// Test string concatenation primitive
	fmt.Println("\nTesting string concatenation primitive...")

	// Create a method for the string class
	stringClass := virtualMachine.StringClass
	commaMethod := &classes.Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes:      []byte{},
		Literals:       []*core.Object{},
		TempVarNames:   []string{},
		IsPrimitive:    true,
		PrimitiveIndex: 30, // Primitive index for string concatenation
	}
	stringClass.AddMethod(classes.NewSymbol(","), classes.MethodToObject(commaMethod))

	// Convert strings to objects
	str1Obj := classes.StringToObject(str1)
	str2Obj := classes.StringToObject(str2)

	// Execute the primitive
	selector := classes.NewSymbol(",")
	method := stringClass.LookupMethod(selector)

	if method == nil {
		fmt.Println("Error: Method not found")
		return
	}

	// selector is already an *core.Object
	resultObj := virtualMachine.ExecutePrimitive(str1Obj, selector, []*core.Object{str2Obj}, method)

	if resultObj == nil {
		fmt.Println("Error: Primitive returned nil")
		return
	}

	resultStr := classes.ObjectToString(resultObj)
	fmt.Printf("Primitive result: %s\n", resultStr.GetValue())

	// Test the string tests
	fmt.Println("\nTesting string tests...")
	fmt.Println("'hello' should return 'hello'")
	fmt.Println("'hello', ' world' should return 'hello world'")
}
