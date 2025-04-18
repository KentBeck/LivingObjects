package main

import (
	"encoding/binary"
	"fmt"
)

// VM represents the Smalltalk virtual machine
type VM struct {
	Globals        map[string]*Object
	CurrentContext *Context
	ObjectMemory   *ObjectMemory

	// Special objects
	NilObject   *Object
	TrueObject  *Object
	FalseObject *Object
	ObjectClass *Object
}

// NewVM creates a new virtual machine
func NewVM() *VM {
	vm := &VM{
		Globals:      make(map[string]*Object),
		ObjectMemory: NewObjectMemory(),
	}

	// Initialize special objects
	vm.NilObject = NewNil()
	vm.TrueObject = NewBoolean(true)
	vm.FalseObject = NewBoolean(false)

	return vm
}

// LoadImage loads a Smalltalk image from a file
func (vm *VM) LoadImage(path string) error {
	// For now, we'll just create a simple test image in memory
	// In a real implementation, this would load from a file

	// Create basic classes
	vm.ObjectClass = NewClass("Object", nil)
	vm.Globals["Object"] = vm.ObjectClass

	// Create a simple test method: 2 + 3
	twoObj := NewInteger(2)
	threeObj := NewInteger(3)
	plusSymbol := NewSymbol("+")

	// Create a method that adds 2 and 3
	methodObj := NewMethod(plusSymbol, vm.ObjectClass)

	// Add literals
	methodObj.Method.Literals = append(methodObj.Method.Literals, twoObj)
	methodObj.Method.Literals = append(methodObj.Method.Literals, threeObj)
	methodObj.Method.Literals = append(methodObj.Method.Literals, plusSymbol)

	// Create bytecodes for: 2 + 3
	// PUSH_LITERAL 0 (2)
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, PUSH_LITERAL)
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, 0, 0, 0, 0) // Index 0

	// PUSH_LITERAL 1 (3)
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, PUSH_LITERAL)
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, 0, 0, 0, 1) // Index 1

	// SEND_MESSAGE 2 ("+") with 1 argument
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, SEND_MESSAGE)
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, 0, 0, 0, 2) // Selector index 2
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, 0, 0, 0, 1) // 1 argument

	// RETURN_STACK_TOP
	methodObj.Method.Bytecodes = append(methodObj.Method.Bytecodes, RETURN_STACK_TOP)

	// Create a context for this method
	vm.CurrentContext = NewContext(methodObj, vm.ObjectClass, []*Object{}, nil)

	return nil
}

// Execute executes the current context
func (vm *VM) Execute() error {
	for vm.CurrentContext != nil {
		// Check if we need to collect garbage
		if vm.ObjectMemory.ShouldCollect() {
			vm.ObjectMemory.Collect(vm)
		}

		// Execute the current context
		result, err := vm.ExecuteContext(vm.CurrentContext)
		if err != nil {
			return err
		}

		// Move to the sender context
		vm.CurrentContext = vm.CurrentContext.Sender

		// If we have a sender, push the result onto its stack
		if vm.CurrentContext != nil {
			vm.CurrentContext.Push(result)
		} else {
			// No sender, print the result
			fmt.Printf("Result: %s\n", result)
		}
	}

	return nil
}

// ExecuteContext executes a single context until it returns
func (vm *VM) ExecuteContext(context *Context) (*Object, error) {
	for {
		// Check if we've reached the end of the method
		if context.PC >= len(context.Method.Method.Bytecodes) {
			return vm.NilObject, nil
		}

		// Get the current bytecode
		bytecode := context.Method.Method.Bytecodes[context.PC]

		// Get the instruction size
		size := InstructionSize(bytecode)

		// Execute the bytecode
		switch bytecode {
		case PUSH_LITERAL:
			// Get the literal index (4 bytes)
			index := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))
			if index < 0 || index >= len(context.Method.Method.Literals) {
				return nil, fmt.Errorf("literal index out of bounds: %d", index)
			}

			// Push the literal onto the stack
			context.Push(context.Method.Method.Literals[index])

		case PUSH_INSTANCE_VARIABLE:
			// Get the instance variable index (4 bytes)
			index := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))
			if index < 0 || index >= len(context.Receiver.Class.InstanceVarNames) {
				return nil, fmt.Errorf("instance variable index out of bounds: %d", index)
			}

			// Get the instance variable name
			name := context.Receiver.Class.InstanceVarNames[index]

			// Push the instance variable onto the stack
			if value, ok := context.Receiver.InstanceVars[name]; ok {
				context.Push(value)
			} else {
				context.Push(vm.NilObject)
			}

		case PUSH_TEMPORARY_VARIABLE:
			// Get the temporary variable index (4 bytes)
			index := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

			// Push the temporary variable onto the stack
			context.Push(context.GetTempVarByIndex(index))

		case PUSH_SELF:
			// Push the receiver onto the stack
			context.Push(context.Receiver)

		case STORE_INSTANCE_VARIABLE:
			// Get the instance variable index (4 bytes)
			index := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))
			if index < 0 || index >= len(context.Receiver.Class.InstanceVarNames) {
				return nil, fmt.Errorf("instance variable index out of bounds: %d", index)
			}

			// Get the instance variable name
			name := context.Receiver.Class.InstanceVarNames[index]

			// Pop the value from the stack
			value := context.Pop()

			// Store the value in the instance variable
			context.Receiver.InstanceVars[name] = value

			// Push the value back onto the stack
			context.Push(value)

		case STORE_TEMPORARY_VARIABLE:
			// Get the temporary variable index (4 bytes)
			index := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

			// Pop the value from the stack
			value := context.Pop()

			// Store the value in the temporary variable
			context.SetTempVarByIndex(index, value)

			// Push the value back onto the stack
			context.Push(value)

		case SEND_MESSAGE:
			// Get the selector index (4 bytes)
			selectorIndex := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))
			if selectorIndex < 0 || selectorIndex >= len(context.Method.Method.Literals) {
				return nil, fmt.Errorf("selector index out of bounds: %d", selectorIndex)
			}

			// Get the argument count (4 bytes)
			argCount := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+5:]))

			// Get the selector
			selector := context.Method.Method.Literals[selectorIndex]
			if selector.Type != OBJ_SYMBOL {
				return nil, fmt.Errorf("selector is not a symbol: %s", selector)
			}

			// Pop the arguments from the stack
			args := make([]*Object, argCount)
			for i := argCount - 1; i >= 0; i-- {
				args[i] = context.Pop()
			}

			// Pop the receiver
			receiver := context.Pop()

			// Handle primitive methods
			if result := vm.executePrimitive(receiver, selector, args); result != nil {
				context.Push(result)
				break
			}

			// Look up the method
			method := vm.lookupMethod(receiver, selector)
			if method == nil {
				return nil, fmt.Errorf("method not found: %s", selector.SymbolValue)
			}

			// Create a new context for the method
			newContext := NewContext(method, receiver, args, context)

			// Set the current context to the new context
			vm.CurrentContext = newContext

			// Return from this context execution to start executing the new context
			return nil, nil

		case RETURN_STACK_TOP:
			// Pop the return value from the stack
			returnValue := context.Pop()

			// Return the value
			return returnValue, nil

		case JUMP:
			// Get the jump target (4 bytes)
			target := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

			// Set the PC to the target
			context.PC = target

			// Skip the normal PC increment
			continue

		case JUMP_IF_TRUE:
			// Get the jump target (4 bytes)
			target := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

			// Pop the condition from the stack
			condition := context.Pop()

			// If the condition is true, jump to the target
			if condition.IsTrue() {
				context.PC = target
				continue
			}

		case JUMP_IF_FALSE:
			// Get the jump target (4 bytes)
			target := int(binary.BigEndian.Uint32(context.Method.Method.Bytecodes[context.PC+1:]))

			// Pop the condition from the stack
			condition := context.Pop()

			// If the condition is false, jump to the target
			if !condition.IsTrue() {
				context.PC = target
				continue
			}

		case POP:
			// Pop the top value from the stack
			context.Pop()

		case DUPLICATE:
			// Duplicate the top value on the stack
			value := context.Top()
			context.Push(value)

		default:
			return nil, fmt.Errorf("unknown bytecode: %d", bytecode)
		}

		// Increment the PC
		context.PC += size
	}
}

// executePrimitive executes a primitive method
func (vm *VM) executePrimitive(receiver *Object, selector *Object, args []*Object) *Object {
	// Handle primitive methods like + - * / for integers
	if receiver.Type == OBJ_INTEGER {
		switch selector.SymbolValue {
		case "+":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER {
				return NewInteger(receiver.IntegerValue + args[0].IntegerValue)
			}
		case "-":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER {
				return NewInteger(receiver.IntegerValue - args[0].IntegerValue)
			}
		case "*":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER {
				return NewInteger(receiver.IntegerValue * args[0].IntegerValue)
			}
		case "/":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER && args[0].IntegerValue != 0 {
				return NewInteger(receiver.IntegerValue / args[0].IntegerValue)
			}
		case "=":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER {
				return NewBoolean(receiver.IntegerValue == args[0].IntegerValue)
			}
		case "<":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER {
				return NewBoolean(receiver.IntegerValue < args[0].IntegerValue)
			}
		case ">":
			if len(args) == 1 && args[0].Type == OBJ_INTEGER {
				return NewBoolean(receiver.IntegerValue > args[0].IntegerValue)
			}
		}
	}

	// No primitive method found
	return nil
}

// lookupMethod looks up a method in a class hierarchy
func (vm *VM) lookupMethod(receiver *Object, selector *Object) *Object {
	// Get the class of the receiver
	class := receiver
	if receiver.Type != OBJ_CLASS {
		class = receiver.Class
	}

	// Look up the method in the class hierarchy
	for class != nil {
		// Check if the class has a method dictionary
		if methodDict, ok := class.InstanceVars["methodDict"]; ok && methodDict.Type == OBJ_DICTIONARY {
			// Check if the method dictionary has the selector
			if method, ok := methodDict.Entries[selector.SymbolValue]; ok {
				return method
			}
		}

		// Move up the class hierarchy
		class = class.SuperClass
	}

	// Method not found
	return nil
}
