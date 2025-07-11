#include "primitives/block.h"
#include "context.h"
#include "memory_manager.h"
#include "interpreter.h"
#include "object.h"
#include "smalltalk_image.h"
#include "compiled_method.h"
#include "primitives.h"

#include <stdexcept>

namespace smalltalk
{

    namespace BlockPrimitives
    {

        TaggedValue value(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            // REAL BLOCK EXECUTION: The receiver IS the block we want to execute
            
            // Validate that receiver is a block context
            if (!receiver.isPointer()) {
                throw std::runtime_error("Block value: receiver must be pointer");
            }
            
            Object* receiverObj = receiver.asObject();
            if (receiverObj->header.getType() != ObjectType::CONTEXT ||
                receiverObj->header.getContextType() != static_cast<uint8_t>(ContextType::BLOCK_CONTEXT)) {
                throw std::runtime_error("Block value: receiver must be block context");
            }
            
            BlockContext* blockContext = static_cast<BlockContext*>(receiverObj);
            
            // Get home context
            if (blockContext->home.isNil() || !blockContext->home.isPointer()) {
                throw std::runtime_error("Block value: invalid home context");
            }
            
            MethodContext* homeContext = static_cast<MethodContext*>(blockContext->home.asObject());
            
            // The block's compiled method is stored directly in the block context's first slot
            char* contextEnd = reinterpret_cast<char*>(blockContext) + sizeof(BlockContext);
            TaggedValue* slots = reinterpret_cast<TaggedValue*>(contextEnd);
            TaggedValue blockMethodValue = slots[0];
            
            if (!blockMethodValue.isPointer()) {
                throw std::runtime_error("Block value: block method is not a pointer");
            }
            
            // Cast to CompiledMethod
            CompiledMethod* blockMethod = reinterpret_cast<CompiledMethod*>(blockMethodValue.asObject());
            
            // Set up execution context for the block
            size_t tempVarCount = blockMethod->getTempVars().size();
            size_t argCount = args.size();
            size_t contextSize = std::max(tempVarCount, argCount) + 20; // temp vars + args + stack space
            
            TaggedValue selfValue = homeContext->self; // Block executes with home context's self
            TaggedValue senderValue = TaggedValue::fromObject(interpreter.getCurrentContext());
            
            MethodContext* blockMethodContext = interpreter.getMemoryManager().allocateMethodContext(
                contextSize,
                0, // No hash needed for direct execution
                selfValue,
                senderValue
            );
            
            // Set up temporary variable and argument slots
            char* methodContextEnd = reinterpret_cast<char*>(blockMethodContext) + sizeof(MethodContext);
            TaggedValue* methodSlots = reinterpret_cast<TaggedValue*>(methodContextEnd);
            
            // LEXICAL SCOPING: Get home context's temporary variables for copying
            char* homeContextEnd = reinterpret_cast<char*>(homeContext) + sizeof(MethodContext);
            TaggedValue* homeSlots = reinterpret_cast<TaggedValue*>(homeContextEnd);
            
            // Get the home context's method to know how many home variables there are
            // We need to find the home method's CompiledMethod to get its temp var count
            // For now, let's use a simple approach: copy all home context variables
            // and let the compiler indexing handle the rest
            
            // LEXICAL SCOPING: The block's temp vars are structured as: [home vars][block params][block temps]
            // We need to copy the home context's variables to the first slots
            
            // Get the number of home variables from the compiled method
            size_t homeVarCount = blockMethod->homeVarCount;
            
            // Copy home context variables to the first slots
            for (size_t i = 0; i < homeVarCount && i < tempVarCount; i++) {
                methodSlots[i] = homeSlots[i];
            }
            
            // Copy block arguments to temporary variable slots after home variables
            // In Smalltalk, block parameters become temporary variables after outer scope vars
            for (size_t i = 0; i < argCount && (homeVarCount + i) < tempVarCount; i++) {
                methodSlots[homeVarCount + i] = args[i];
            }
            
            // Initialize remaining temporary variables to nil
            for (size_t i = homeVarCount + argCount; i < tempVarCount; i++) {
                methodSlots[i] = TaggedValue::nil();
            }
            
            // Set stack pointer to start after temporary variables
            blockMethodContext->stackPointer = methodSlots + tempVarCount;
            
            // Execute the block's compiled method directly
            TaggedValue result = interpreter.executeMethodContext(blockMethodContext, blockMethod);
            
            return result;
        }

    } // namespace BlockPrimitives

} // namespace smalltalk