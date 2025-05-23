package pile

import (
	"fmt"
	"unsafe"
)

// Method represents a Smalltalk method
type Method struct {
	Object
	Bytecodes      []byte
	Literals       []*Object
	Selector       *Object
	TempVarNames   []string
	MethodClass    *Class
	IsPrimitive    bool
	PrimitiveIndex int
}

// newMethod creates a new method object without setting its class field
// This is a private helper function used by vm.NewMethod
func NewMethodInternal(selector *Object, class *Class) *Method {
	return &Method{
		Object: Object{
			TypeField: OBJ_METHOD,
		},
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*Object, 0),
		Selector:     selector,
		MethodClass:  class,
		TempVarNames: make([]string, 0),
	}
}

// NewMethod creates a new method object (deprecated - use vm.NewMethod instead)
func NewMethod(selector *Object, class *Class) *Object {
	return MethodToObject(NewMethodInternal(selector, class))
}

// MethodToObject converts a Method to an Object
func MethodToObject(m *Method) *Object {
	return (*Object)(unsafe.Pointer(m))
}

// ObjectToMethod converts an Object to a Method
func ObjectToMethod(o *Object) *Method {
	if o == nil || o.Type() != OBJ_METHOD {
		return nil
	}
	return (*Method)(unsafe.Pointer(o))
}

// String returns a string representation of the method object
func (m *Method) String() string {
	if m.Selector != nil {
		return fmt.Sprintf("Method %s", GetSymbolValue(m.Selector))
	}
	return "a Method"
}

// GetBytecodes returns the bytecodes of the method
func (m *Method) GetBytecodes() []byte {
	return m.Bytecodes
}

// SetBytecodes sets the bytecodes of the method
func (m *Method) SetBytecodes(bytecodes []byte) {
	m.Bytecodes = bytecodes
}

// GetLiterals returns the literals of the method
func (m *Method) GetLiterals() []*Object {
	return m.Literals
}

// AddLiteral adds a literal to the method
func (m *Method) AddLiteral(literal *Object) {
	m.Literals = append(m.Literals, literal)
}

// GetSelector returns the selector of the method
func (m *Method) GetSelector() *Object {
	return m.Selector
}

// SetSelector sets the selector of the method
func (m *Method) SetSelector(selector *Object) {
	m.Selector = selector
}

// GetTempVarNames returns the temporary variable names of the method
func (m *Method) GetTempVarNames() []string {
	return m.TempVarNames
}

// AddTempVarName adds a temporary variable name to the method
func (m *Method) AddTempVarName(name string) {
	m.TempVarNames = append(m.TempVarNames, name)
}

// GetMethodClass returns the class of the method
func (m *Method) GetMethodClass() *Class {
	return m.MethodClass
}

// SetMethodClass sets the class of the method
func (m *Method) SetMethodClass(class *Class) {
	m.MethodClass = class
}

// IsPrimitiveMethod returns true if the method is a primitive
func (m *Method) IsPrimitiveMethod() bool {
	return m.IsPrimitive
}

// SetPrimitive sets whether the method is a primitive
func (m *Method) SetPrimitive(isPrimitive bool) {
	m.IsPrimitive = isPrimitive
}

// GetPrimitiveIndex returns the primitive index of the method
func (m *Method) GetPrimitiveIndex() int {
	return m.PrimitiveIndex
}

// SetPrimitiveIndex sets the primitive index of the method
func (m *Method) SetPrimitiveIndex(index int) {
	m.PrimitiveIndex = index
}