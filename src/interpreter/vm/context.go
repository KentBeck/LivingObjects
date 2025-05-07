package vm

import (
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// Context represents a method activation context
type Context struct {
	Method       *core.Object
	Receiver     core.ObjectInterface
	Arguments    []*core.Object
	TempVars     []core.ObjectInterface // Temporary variables stored by index
	Sender       *Context
	PC           int
	Stack        []*core.Object
	StackPointer int
}

// NewContext creates a new method activation context
func NewContext(method *core.Object, receiver core.ObjectInterface, arguments []*core.Object, sender *Context) *Context {
	if method == nil {
		panic("NewContext: nil method")
	}
	methodObj := classes.ObjectToMethod(method)
	if methodObj == nil { // temporary
		panic("NewContext: nil method")
	}

	// Initialize temporary variables array with nil values
	tempVarsSize := len(methodObj.GetTempVarNames())
	tempVars := make([]core.ObjectInterface, tempVarsSize)
	for i := range tempVars {
		tempVars[i] = core.NewNil()
	}

	return &Context{
		Method:       method,
		Receiver:     receiver,
		Arguments:    arguments,
		TempVars:     tempVars,
		Sender:       sender,
		PC:           0,
		Stack:        make([]*core.Object, 100), // Initial stack size
		StackPointer: 0,
	}
}

// Push pushes an object onto the stack
func (c *Context) Push(obj core.ObjectInterface) {
	// Grow stack if needed
	if c.StackPointer >= len(c.Stack) {
		newStack := make([]*core.Object, len(c.Stack)*2)
		copy(newStack, c.Stack)
		c.Stack = newStack
	}

	// Handle nil values
	if obj == nil {
		c.Stack[c.StackPointer] = nil
	} else {
		c.Stack[c.StackPointer] = obj.(*core.Object)
	}
	c.StackPointer++
}

// Pop pops an object from the stack
func (c *Context) Pop() *core.Object {
	if c.StackPointer <= 0 {
		panic("stack underflow")
	}

	c.StackPointer--
	obj := c.Stack[c.StackPointer]
	return obj
}

// Top returns the top object on the stack without popping it
func (c *Context) Top() *core.Object {
	if c.StackPointer <= 0 {
		panic("stack underflow")
	}

	return c.Stack[c.StackPointer-1]
}

// GetTempVarByIndex gets a temporary variable by index
func (c *Context) GetTempVarByIndex(index int) core.ObjectInterface {
	if index < 0 || index >= len(c.TempVars) {
		panic("index out of bounds")
	}

	return c.TempVars[index]
}

// SetTempVarByIndex sets a temporary variable by index
func (c *Context) SetTempVarByIndex(index int, value core.ObjectInterface) {
	if index < 0 || index >= len(c.TempVars) {
		return
	}

	if value == nil {
		c.TempVars[index] = nil
	} else {
		c.TempVars[index] = value.(*core.Object)
	}
}

// GetMethod returns the method of the context
func (c *Context) GetMethod() *core.Object {
	return c.Method
}

// GetReceiver returns the receiver of the context
func (c *Context) GetReceiver() *core.Object {
	return c.Receiver.(*core.Object)
}

// GetArguments returns the arguments of the context
func (c *Context) GetArguments() []*core.Object {
	return c.Arguments
}

// GetTempVars returns the temporary variables of the context
func (c *Context) GetTempVars() []*core.Object {
	result := make([]*core.Object, len(c.TempVars))
	for i, tempVar := range c.TempVars {
		if tempVar != nil {
			result[i] = tempVar.(*core.Object)
		}
	}
	return result
}

// GetStack returns the stack of the context
func (c *Context) GetStack() []*core.Object {
	return c.Stack
}

// GetStackPointer returns the stack pointer of the context
func (c *Context) GetStackPointer() int {
	return c.StackPointer
}

// GetSender returns the sender context
func (c *Context) GetSender() interface{} {
	return c.Sender
}

// SetSender sets the sender context
func (c *Context) SetSender(sender *Context) {
	c.Sender = sender
}

// GetPC returns the program counter
func (c *Context) GetPC() int {
	return c.PC
}

// SetPC sets the program counter
func (c *Context) SetPC(pc int) {
	c.PC = pc
}
