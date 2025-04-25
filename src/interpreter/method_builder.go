package main

type MethodBuilder struct {
	bytecodes      []byte
	literals       []*Object
	tempVarNames   []string
	isPrimitive    bool
	primitiveIndex int
}

func NewMethodBuilder() *MethodBuilder {
	return &MethodBuilder{
		bytecodes:      make([]byte, 0),
		literals:       make([]*Object, 0),
		tempVarNames:   make([]string, 0),
		isPrimitive:    false,
		primitiveIndex: 0,
	}
}

func (mb *MethodBuilder) AddPrimitive(primitive int) {
	mb.primitiveIndex = primitive
	mb.isPrimitive = true
}

func (mb *MethodBuilder) Install(class *Object, selector *Symbol) {
	methodDict := class.GetMethodDict()
	method := mb.method(class, selector)
	methodDict.Entries[GetSymbolValue(SymbolToObject(selector))] = method
}

func (mb *MethodBuilder) method(class *Object, selector *Symbol) *Object {
	// Create the method object
	method := NewMethod(SymbolToObject(selector), class)
	method.Method.Bytecodes = mb.bytecodes
	method.Method.Literals = mb.literals
	method.Method.TempVarNames = mb.tempVarNames
	method.Method.IsPrimitive = mb.isPrimitive
	method.Method.PrimitiveIndex = mb.primitiveIndex
	return method
}
