package main

// MethodBuilder provides a fluent interface for creating methods
type MethodBuilder struct {
	class          *Object
	selectorName   string
	selectorObj    *Object
	bytecodes      []byte
	literals       []*Object
	tempVarNames   []string
	isPrimitive    bool
	primitiveIndex int
}

// NewMethodBuilder creates a new MethodBuilder for the given class
func NewMethodBuilder(class *Object) *MethodBuilder {
	return &MethodBuilder{
		class:          class,
		bytecodes:      make([]byte, 0),
		literals:       make([]*Object, 0),
		tempVarNames:   make([]string, 0),
		isPrimitive:    false,
		primitiveIndex: 0,
	}
}

// Selector sets the selector for the method
func (mb *MethodBuilder) Selector(name string) *MethodBuilder {
	mb.selectorName = name
	mb.selectorObj = NewSymbol(name)
	return mb
}

// Primitive marks the method as a primitive with the given index
func (mb *MethodBuilder) Primitive(index int) *MethodBuilder {
	mb.primitiveIndex = index
	mb.isPrimitive = true
	return mb
}

// Bytecodes adds bytecodes to the method
func (mb *MethodBuilder) Bytecodes(bytecodes []byte) *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecodes...)
	return mb
}

// Literals adds multiple literals to the method
func (mb *MethodBuilder) Literals(literals []*Object) *MethodBuilder {
	mb.literals = append(mb.literals, literals...)
	return mb
}

// AddLiteral adds a single literal to the method and returns its index
func (mb *MethodBuilder) AddLiteral(literal *Object) (int, *MethodBuilder) {
	index := len(mb.literals)
	mb.literals = append(mb.literals, literal)
	return index, mb
}

// TempVars adds temporary variable names to the method
func (mb *MethodBuilder) TempVars(names []string) *MethodBuilder {
	mb.tempVarNames = append(mb.tempVarNames, names...)
	return mb
}

// Go finalizes the method creation and adds it to the class's method dictionary
func (mb *MethodBuilder) Go() *Object {
	if mb.selectorObj == nil {
		panic("Selector not set. Call Selector() first.")
	}

	// Create the method object
	method := NewMethod(mb.selectorObj, mb.class)

	// Set the method properties
	method.Method.Bytecodes = mb.bytecodes
	method.Method.Literals = mb.literals
	method.Method.TempVarNames = mb.tempVarNames
	method.Method.IsPrimitive = mb.isPrimitive
	method.Method.PrimitiveIndex = mb.primitiveIndex

	// Add the method to the method dictionary
	symbolValue := GetSymbolValue(mb.selectorObj)
	methodDict := mb.class.GetMethodDict()
	methodDict.Entries[symbolValue] = method

	// Reset the builder state for reuse
	mb.bytecodes = make([]byte, 0)
	mb.literals = make([]*Object, 0)
	mb.tempVarNames = make([]string, 0)
	mb.isPrimitive = false
	mb.primitiveIndex = 0
	// Note: We don't reset the class or selector as they might be reused

	return method
}
