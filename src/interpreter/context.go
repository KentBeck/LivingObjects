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
func NewContext(method *Object, receiver ObjectInterface, arguments []*Object, sender *Context) *Context {
	if method == nil {
		panic("NewContext: nil method")
	}

	// Initialize temporary variables array with nil values
	tempVarsSize := 0
	if method.Method != nil {
		tempVarsSize = len(method.Method.TempVarNames)
	}
	tempVars := make([]*Object, tempVarsSize)
	for i := range tempVars {
		tempVars[i] = NewNil()
	}

	return &Context{
		Method:       method,
		Receiver:     receiver.(*Object),
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
		panic("stack underflow")
	}

	c.StackPointer--
	obj := c.Stack[c.StackPointer]
	return obj
}

// Top returns the top object on the stack without popping it
func (c *Context) Top() *Object {
	if c.StackPointer <= 0 {
		panic("stack underflow")
	}

	return c.Stack[c.StackPointer-1]
}

// GetTempVarByIndex gets a temporary variable by index
func (c *Context) GetTempVarByIndex(index int) *Object {
	if index < 0 || index >= len(c.TempVars) {
		panic("index out of bounds")
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
