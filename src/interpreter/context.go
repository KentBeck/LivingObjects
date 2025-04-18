package main

// Context represents a method activation context
type Context struct {
	Method       *Object
	Receiver     *Object
	Arguments    []*Object
	TempVars     []*Object // Temporary variables stored by index
	Sender       *Context
	PC           int
	Stack        []*Object
	StackPointer int
}

// NewContext creates a new method activation context
func NewContext(method *Object, receiver *Object, arguments []*Object, sender *Context) *Context {
	// Initialize temporary variables array with nil values
	tempVarsSize := 0
	if method != nil && method.Method != nil {
		tempVarsSize = len(method.Method.TempVarNames)
	}
	tempVars := make([]*Object, tempVarsSize)
	for i := range tempVars {
		tempVars[i] = NewNil()
	}

	return &Context{
		Method:       method,
		Receiver:     receiver,
		Arguments:    arguments,
		TempVars:     tempVars,
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

// GetTempVar gets a temporary variable by name (for backward compatibility)
func (c *Context) GetTempVar(name string) *Object {
	// Find the index of the name in the method's temporary variable names
	if c.Method == nil || c.Method.Method == nil {
		return NewNil()
	}

	for i, tempName := range c.Method.Method.TempVarNames {
		if tempName == name {
			return c.GetTempVarByIndex(i)
		}
	}

	return NewNil()
}

// SetTempVar sets a temporary variable by name (for backward compatibility)
func (c *Context) SetTempVar(name string, value *Object) {
	// Find the index of the name in the method's temporary variable names
	if c.Method == nil || c.Method.Method == nil {
		return
	}

	for i, tempName := range c.Method.Method.TempVarNames {
		if tempName == name {
			c.SetTempVarByIndex(i, value)
			return
		}
	}
}

// GetTempVarByIndex gets a temporary variable by index
func (c *Context) GetTempVarByIndex(index int) *Object {
	if index < 0 || index >= len(c.TempVars) {
		return NewNil()
	}

	return c.TempVars[index]
}

// SetTempVarByIndex sets a temporary variable by index
func (c *Context) SetTempVarByIndex(index int, value *Object) {
	if index < 0 || index >= len(c.TempVars) {
		return
	}

	c.TempVars[index] = value
}
