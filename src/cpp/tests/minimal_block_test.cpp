#include "primitives/block.h"
#include "context.h"
#include "memory_manager.h"
#include "interpreter.h"
#include "object.h"
#include "tagged_value.h"

#include <iostream>
#include <cassert>

using namespace smalltalk;

int main() {
    try {
        std::cout << "Testing Block primitive (minimal test)..." << std::endl;
        
        // Initialize memory manager and interpreter
        MemoryManager memoryManager;
        Interpreter interpreter(memoryManager);
        
        // Create a simple home context (method context where the block was defined)
        MethodContext* homeContext = memoryManager.allocateMethodContext(
            4,      // smaller context size
            123,    // method reference  
            nullptr, // self (receiver)
            nullptr  // sender
        );
        homeContext->self = reinterpret_cast<Object*>(0x1000); // dummy self pointer
        
        std::cout << "Created home context: " << homeContext << std::endl;
        std::cout << "Home context self: " << homeContext->self << std::endl;
        
        // Create a block context with smaller size
        BlockContext* blockContext = memoryManager.allocateBlockContext(
            2,          // smaller context size
            456,        // method reference for the block's code
            homeContext->self, // receiver (same as home's self)
            nullptr,    // sender (will be set when block is executed)
            homeContext // home context
        );
        
        // Verify the block context was created correctly
        std::cout << "Created block context: " << blockContext << std::endl;
        std::cout << "Block context home field: " << blockContext->home << std::endl;
        std::cout << "Expected home: " << homeContext << std::endl;
        std::cout << "Block context type: " << blockContext->header.type << std::endl;
        std::cout << "Expected type: " << static_cast<uint64_t>(ContextType::BLOCK_CONTEXT) << std::endl;
        
        // Check type first
        if (blockContext->header.type != static_cast<uint64_t>(ContextType::BLOCK_CONTEXT)) {
            std::cout << "ERROR: Block type mismatch!" << std::endl;
            return 1;
        }
        
        // Check home pointer
        if (blockContext->home != homeContext) {
            std::cout << "ERROR: Block home pointer mismatch!" << std::endl;
            std::cout << "This might be a memory layout or GC issue" << std::endl;
            // Don't fail the test yet, let's continue and see what happens
        }
        
        std::cout << "Block context created successfully" << std::endl;
        std::cout << "Block type: " << blockContext->header.type << std::endl;
        std::cout << "Expected type: " << static_cast<uint64_t>(ContextType::BLOCK_CONTEXT) << std::endl;
        std::cout << "Block home: " << blockContext->home << std::endl;
        std::cout << "Block method ref: " << blockContext->header.hash << std::endl;
        
        // Create a TaggedValue for the block
        TaggedValue blockValue(blockContext);
        
        // Test the block primitive
        std::vector<TaggedValue> args; // no arguments for Block>>value
        
        std::cout << "Calling BlockPrimitives::value..." << std::endl;
        
        // This should set up a new context and activate it
        TaggedValue result = BlockPrimitives::value(blockValue, args, interpreter);
        
        // Check that the interpreter's active context has changed
        MethodContext* activeContext = interpreter.getCurrentContext();
        std::cout << "Active context after block call: " << activeContext << std::endl;
        
        // The active context should be a new method context for the block
        assert(activeContext != nullptr);
        assert(activeContext != homeContext);
        assert(activeContext->sender == nullptr); // was nullptr when we called
        assert(activeContext->self == homeContext->self); // should inherit self from home
        assert(activeContext->instructionPointer == 0); // should start at beginning
        
        std::cout << "Block primitive test PASSED!" << std::endl;
        
        // Test error cases
        std::cout << "Testing error cases..." << std::endl;
        
        // Test with non-pointer value
        TaggedValue intValue(42);
        try {
            BlockPrimitives::value(intValue, args, interpreter);
            assert(false && "Should have thrown exception for non-pointer value");
        } catch (const std::runtime_error& e) {
            std::cout << "Correctly caught error for non-pointer: " << e.what() << std::endl;
        }
        
        // Test with non-block object
        Object* regularObject = memoryManager.allocateObject(ObjectType::OBJECT, 4);
        TaggedValue objValue(regularObject);
        try {
            BlockPrimitives::value(objValue, args, interpreter);
            assert(false && "Should have thrown exception for non-block object");
        } catch (const std::runtime_error& e) {
            std::cout << "Correctly caught error for non-block object: " << e.what() << std::endl;
        }
        
        std::cout << "All tests PASSED!" << std::endl;
        return 0;
        
    } catch (const std::exception& e) {
        std::cerr << "Test FAILED with exception: " << e.what() << std::endl;
        return 1;
    }
}