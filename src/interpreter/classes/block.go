package classes

import (
	"unsafe"

	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/runtime"
)

// Block represents a Smalltalk block
type Block struct {
	core.Object
	Bytecodes    []byte
	Literals     []*core.Object
	TempVarNames []string
	OuterContext interface{} // Using interface{} to avoid circular dependency
}

// NewBlock creates a new block object
func NewBlock(outerContext interface{}) *core.Object {
	block := &Block{
		Object: core.Object{
			TypeField: core.OBJ_BLOCK,
		},
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*core.Object, 0),
		TempVarNames: make([]string, 0),
		OuterContext: outerContext,
	}

	return BlockToObject(block)
}

// BlockToObject converts a Block to an Object
func BlockToObject(b *Block) *core.Object {
	return (*core.Object)(unsafe.Pointer(b))
}

// ObjectToBlock converts an Object to a Block
func ObjectToBlock(o *core.Object) *Block {
	return (*Block)(unsafe.Pointer(o))
}

// String returns a string representation of the block object
func (b *Block) String() string {
	return "Block"
}

// GetBytecodes returns the bytecodes of the block
func (b *Block) GetBytecodes() []byte {
	return b.Bytecodes
}

// SetBytecodes sets the bytecodes of the block
func (b *Block) SetBytecodes(bytecodes []byte) {
	b.Bytecodes = bytecodes
}

// GetLiterals returns the literals of the block
func (b *Block) GetLiterals() []*core.Object {
	return b.Literals
}

// AddLiteral adds a literal to the block
func (b *Block) AddLiteral(literal *core.Object) {
	b.Literals = append(b.Literals, literal)
}

// GetTempVarNames returns the temporary variable names of the block
func (b *Block) GetTempVarNames() []string {
	return b.TempVarNames
}

// AddTempVarName adds a temporary variable name to the block
func (b *Block) AddTempVarName(name string) {
	b.TempVarNames = append(b.TempVarNames, name)
}

// GetOuterContext returns the outer context of the block
func (b *Block) GetOuterContext() interface{} {
	return b.OuterContext
}

// SetOuterContext sets the outer context of the block
func (b *Block) SetOuterContext(outerContext interface{}) {
	b.OuterContext = outerContext
}

// Value evaluates the block with the given arguments
func (b *Block) Value(args ...*core.Object) *core.Object {
	return b.ValueWithArguments(args)
}

// ValueWithArguments evaluates the block with the given arguments
func (b *Block) ValueWithArguments(args []*core.Object) *core.Object {
	// Convert the block to an Object
	blockObj := BlockToObject(b)

	// Use the runtime package to execute the block
	return runtime.ExecuteBlock(blockObj, args)
}
