#pragma once

#include "bytecode.h"
#include "context.h"
#include "memory_manager.h"
#include "object.h"
#include "tagged_value.h"
#include "compiled_method.h"

#include <cstdint>
#include <vector>

namespace smalltalk
{

    class SmalltalkImage; // Forward declaration

    class Interpreter
    {
    public:
        // Constructor
        Interpreter(MemoryManager &memory, SmalltalkImage& image);

        // Execute method
        Object *executeMethod(CompiledMethod *method, Object *receiver, std::vector<Object *> &args);

        // Execute compiled method directly
        TaggedValue executeCompiledMethod(const CompiledMethod &method);
        
        // Execute compiled method with specific context
        TaggedValue executeCompiledMethod(const CompiledMethod &method, MethodContext *context);

        // Unified method context execution (replaces executeLoop + dispatch)
        TaggedValue executeMethodContext(MethodContext *context);
        
        // Direct method execution without hash lookup (fixes architectural issue)
        TaggedValue executeMethodContext(MethodContext *context, CompiledMethod *method);

        // Bytecode handlers
        void handlePushTemporaryVariable(uint32_t offset);
        void handleStoreTemporaryVariable(uint32_t offset);
        void handlePop();
        void handleDuplicate();
        void handleCreateBlock(uint32_t literalIndex, uint32_t parameterCount);

        // Stack operations
        void push(TaggedValue value);
        TaggedValue pop();
        TaggedValue top();

        // Context access
        MethodContext *getCurrentContext() const { return activeContext; }
        void setCurrentContext(MethodContext *context) { activeContext = context; }
        
        // Current method access (for block execution)
        CompiledMethod *getCurrentMethod() const { return currentMethod; }

        // Memory manager access
        MemoryManager &getMemoryManager() { return memoryManager; }
        
        // Image access
        SmalltalkImage& getImage() { return image; }

        // TaggedValue message sending
        TaggedValue sendMessage(TaggedValue receiver, const std::string &selector, const std::vector<TaggedValue> &args);

        // Get object class for TaggedValue
        Class *getObjectClass(TaggedValue value);
        
        // Exception handling
        bool findExceptionHandler(const std::string& exceptionClass, MethodContext*& handlerContext, int& handlerPC);
        void unwindToContext(MethodContext* targetContext);

    private:
        // Memory manager
        MemoryManager &memoryManager;
        SmalltalkImage& image;

        // Current context and chunk
        MethodContext *activeContext = nullptr;
        StackChunk *currentChunk = nullptr;
        
        // Current method being executed (eliminates hash lookup)
        CompiledMethod *currentMethod = nullptr;

        // Internal state
        bool executing = false;

        // Helper for message sending
        Object *sendMessage(Object *receiver, Object *selector, std::vector<Object *> &args);

        // Simple evaluation method for testing (temporarily disabled)
        // TaggedValue evaluate(const std::string& expression);

        // Context switching
        void switchContext(MethodContext *newContext);

        // Stack bounds checking helpers
        Object **getStackStart(MethodContext *context);
        Object **getStackEnd(MethodContext *context);
        Object **getCurrentStackPointer(MethodContext *context);
        void validateStackBounds(MethodContext *context, Object **stackPointer);
    };

} // namespace smalltalk
