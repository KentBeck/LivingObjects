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

        // Execute context
        Object *executeContext(MethodContext *context);

        // Unified method context execution (replaces executeLoop + dispatch)
        TaggedValue executeMethodContext(MethodContext *context);

        // Bytecode handlers
        void handlePushLiteral(uint32_t index);
        void handlePushInstanceVariable(uint32_t offset);
        void handlePushTemporaryVariable(uint32_t offset);
        void handlePushSelf();
        void handleStoreInstanceVariable(uint32_t offset);
        void handleStoreTemporaryVariable(uint32_t offset);
        void handleSendMessage(uint32_t selectorIndex, uint32_t argCount);
        void handleReturnStackTop();
        void handleJump(uint32_t target);
        void handleJumpIfTrue(uint32_t target);
        void handleJumpIfFalse(uint32_t target);
        void handlePop();
        void handleDuplicate();
        void handleCreateBlock(uint32_t bytecodeSize, uint32_t literalCount, uint32_t tempVarCount);
        void handleExecuteBlock(uint32_t argCount);

        // Stack operations
        void push(TaggedValue value);
        TaggedValue pop();
        TaggedValue top();

        // Context access
        MethodContext *getCurrentContext() const { return activeContext; }
        void setCurrentContext(MethodContext *context) { activeContext = context; }

        // Memory manager access
        MemoryManager &getMemoryManager() { return memoryManager; }
        
        // Image access
        SmalltalkImage& getImage() { return image; }

        // TaggedValue message sending
        TaggedValue sendMessage(TaggedValue receiver, const std::string &selector, const std::vector<TaggedValue> &args);

        // Get object class for TaggedValue
        Class *getObjectClass(TaggedValue value);

    private:
        // Memory manager
        MemoryManager &memoryManager;
        SmalltalkImage& image;

        // Current context and chunk
        MethodContext *activeContext = nullptr;
        StackChunk *currentChunk = nullptr;

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
