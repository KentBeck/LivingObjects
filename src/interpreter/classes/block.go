package classes

import (
	"unsafe"

	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/runtime"
	"smalltalklsp/interpreter/types"
)

// Block represents a Smalltalk block
type Block struct {
	core.Object
	Bytecodes    []byte
	Literals     []*core.Object
	TempVarNames []string
	OuterContext interface{} // Using interface{} to avoid circular dependency
}

// newBlock creates a new block object without setting its class field
// This is a private helper function used by vm.NewBlock
func NewBlockInternal(outerContext interface{}) *Block {
	return &Block{
		Object: core.Object{
			TypeField: core.OBJ_BLOCK,
		},
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*core.Object, 0),
		TempVarNames: make([]string, 0),
		OuterContext: outerContext,
	}
}

// NewBlock creates a new block object with proper class field
func NewBlock(outerContext interface{}) *core.Object {
	// If a factory is registered, use it to create blocks with proper class field
	factory := types.GetFactory()
	if factory != nil {
		return factory.NewBlock(outerContext)
	}
	
	// Fall back to simple block creation without class field
	// This is mainly for tests that don't need the VM
	return BlockToObject(NewBlockInternal(outerContext))
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

// OnDo implements the on:do: method for exception handling
func (b *Block) OnDo(exceptionClass *core.Object, handlerBlock *core.Object) *core.Object {
	// Convert blocks to proper types
	handlerBlockObj := ObjectToBlock(handlerBlock)

	// Store the current exception handler
	savedHandler := runtime.CurrentExceptionHandler

	// Create a new exception handler
	handler := &runtime.ExceptionHandler{
		ExceptionClass: exceptionClass,
		HandlerBlock:   handlerBlock,
		NextHandler:    savedHandler,
	}

	// Set the current exception handler
	runtime.CurrentExceptionHandler = handler

	// Execute the receiver block
	var result *core.Object

	// We need to use a defer to ensure the handler is restored
	defer func() {
		// Restore the previous exception handler
		runtime.CurrentExceptionHandler = savedHandler

		// Handle panic if it's an exception
		if r := recover(); r != nil {
			if exception, ok := r.(*core.Object); ok && exception.Type() == core.OBJ_EXCEPTION {
				// Check if the exception is of the handled class
				if runtime.IsKindOf(exception, exceptionClass) {
					// Execute the handler block with the exception as argument
					result = handlerBlockObj.ValueWithArguments([]*core.Object{exception})
				} else {
					// Re-panic for unhandled exceptions
					panic(r)
				}
			} else {
				// Re-panic for non-exception panics
				panic(r)
			}
		}
	}()

	// Execute the protected block
	result = b.Value()

	return result
}
