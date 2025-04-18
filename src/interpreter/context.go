package main

// Context represents a method activation context
type Context struct {
	Method       *Object
	Receiver     *Object
	Arguments    []*Object
	TempVars     map[string]*Object
	Sender       *Context
	PC           int
	Stack        []*Object
	StackPointer int
}

// NewContext creates a new method activation context
func NewContext(method *Object, receiver *Object, arguments []*Object, sender *Context) *Context {
	return &Context{
		Method:       method,
		Receiver:     receiver,
		Arguments:    arguments,
		TempVars:     make(map[string]*Object),
		Sender:       sender,
		PC:           0,
		Stack:        make([]*Object, 100), // Initial stack size
		StackPointer: 0,
	}
}

// Push pushes an object onto the stack
func (c *Context) Push(obj *Object) {
	// Grow stack if needed
	if c.StackPointer >= len(c.Stack) {
		newStack := make([]*Object, len(c.Stack)*2)
		copy(newStack, c.Stack)
		c.Stack = newStack
	}
	
	c.Stack[c.StackPointer] = obj
	c.StackPointer++
}

// Pop pops an object from the stack
func (c *Context) Pop() *Object {
	if c.StackPointer <= 0 {
		return NewNil() // Stack underflow
	}
	
	c.StackPointer--
	return c.Stack[c.StackPointer]
}

// Top returns the top object on the stack without popping it
func (c *Context) Top() *Object {
	if c.StackPointer <= 0 {
		return NewNil() // Stack underflow
	}
	
	return c.Stack[c.StackPointer-1]
}

// GetTempVar gets a temporary variable by name
func (c *Context) GetTempVar(name string) *Object {
	if obj, ok := c.TempVars[name]; ok {
		return obj
	}
	return NewNil()
}

// SetTempVar sets a temporary variable by name
func (c *Context) SetTempVar(name string, value *Object) {
	c.TempVars[name] = value
}

// GetTempVarByIndex gets a temporary variable by index
func (c *Context) GetTempVarByIndex(index int) *Object {
	if index < 0 || index >= len(c.Method.Method.TempVarNames) {
		return NewNil()
	}
	
	name := c.Method.Method.TempVarNames[index]
	return c.GetTempVar(name)
}

// SetTempVarByIndex sets a temporary variable by index
func (c *Context) SetTempVarByIndex(index int, value *Object) {
	if index < 0 || index >= len(c.Method.Method.TempVarNames) {
		return
	}
	
	name := c.Method.Method.TempVarNames[index]
	c.SetTempVar(name, value)
}
