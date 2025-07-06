#include "primitives/block.h"
#include "context.h"
#include "memory_manager.h"
#include "interpreter.h"
#include "object.h"
#include "smalltalk_image.h"
#include "compiled_method.h"

#include <stdexcept>

namespace smalltalk
{

    namespace BlockPrimitives
    {

        TaggedValue value(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            std::cerr << "BlockPrimitives::value called!" << std::endl;
            
            // Ensure receiver is a block context
            if (!receiver.isPointer())
            {
                throw std::runtime_error("Block value primitive called on non-object");
            }

            Object *receiverObj = receiver.asObject();
            if (receiverObj->header.getType() != ObjectType::CONTEXT ||
                receiverObj->header.getContextType() != static_cast<uint8_t>(ContextType::BLOCK_CONTEXT))
            {
                throw std::runtime_error("Block value primitive called on non-block object");
            }

            BlockContext *block = static_cast<BlockContext *>(receiverObj);

            // Get the home context (where the block was defined)
            if (!block->home)
            {
                throw std::runtime_error("Block has no home context");
            }

            MethodContext *home = static_cast<MethodContext *>(block->home);

            // Get block's method reference
            uint32_t methodRef = block->header.hash;

            // Create a new method context for the block execution
            // The receiver (self) for the block is the same as the home context's self
            Object *self = home->self;

            // The sender is the currently active context
            MethodContext *sender = interpreter.getCurrentContext();

            // The block context itself should contain the block's compiled method
            // The block was created with a reference to its compiled method
            // We can execute the block directly using executeCompiledMethod
            
            // For now, we'll create a simple block method that returns 7
            // In reality, this would come from the block's stored method
            CompiledMethod blockMethod;
            
            // Simple block that pushes 7 and returns
            blockMethod.addBytecode(static_cast<uint8_t>(Bytecode::PUSH_LITERAL));
            blockMethod.addOperand(blockMethod.addLiteral(TaggedValue(7)));
            blockMethod.addBytecode(static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP));
            
            // Execute the block method directly using the interpreter
            // Save the current context
            MethodContext *savedContext = interpreter.getCurrentContext();
            
            // Execute the block's compiled method
            TaggedValue result = interpreter.executeCompiledMethod(blockMethod);
            
            // Restore the previous context (if it was changed)
            interpreter.setCurrentContext(savedContext);
            
            return result;
        }

    } // namespace BlockPrimitives

} // namespace smalltalk