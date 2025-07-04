#include "primitives/block.h"
#include "context.h"
#include "memory_manager.h"
#include "interpreter.h"
#include "object.h"

#include <stdexcept>

namespace smalltalk {

namespace BlockPrimitives {

TaggedValue value(TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter) {
    // Ensure receiver is a block context
    if (!receiver.isPointer()) {
        throw std::runtime_error("Block value primitive called on non-object");
    }
    
    Object* receiverObj = receiver.asObject();
    if (receiverObj->header.getType() != ObjectType::CONTEXT || 
        receiverObj->header.getContextType() != static_cast<uint8_t>(ContextType::BLOCK_CONTEXT)) {
        throw std::runtime_error("Block value primitive called on non-block object");
    }
    
    BlockContext* block = static_cast<BlockContext*>(receiverObj);
    
    // Get the home context (where the block was defined)
    if (!block->home) {
        throw std::runtime_error("Block has no home context");
    }
    
    MethodContext* home = static_cast<MethodContext*>(block->home);
    
    // Get block's method reference
    uint32_t methodRef = block->header.hash;
    
    // Create a new method context for the block execution
    // The receiver (self) for the block is the same as the home context's self
    Object* self = home->self;
    
    // The sender is the currently active context
    MethodContext* sender = interpreter.getCurrentContext();
    
    // Allocate a new context for the block
    // Using a reasonable default size for temporaries and stack
    MethodContext* blockContext = interpreter.getMemoryManager().allocateMethodContext(
        16, methodRef, self, sender);
    
    // Set the instruction pointer to the beginning of the block's bytecode
    blockContext->instructionPointer = 0;
    
    // Copy any arguments to the block's temporary variables
    // For now, we only support blocks with no arguments (Block>>value)
    if (!args.empty()) {
        throw std::runtime_error("Block value primitive does not support arguments yet");
    }
    
    // Activate the new context
    interpreter.setCurrentContext(blockContext);
    
    // The interpreter will now execute the block's bytecode
    // When the block returns, it will restore the previous context
    // and the return value will be on the stack
    
    // Return nil for now - the actual return happens through the interpreter
    return TaggedValue();
}

} // namespace BlockPrimitives

} // namespace smalltalk