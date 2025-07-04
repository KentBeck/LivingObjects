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

    class Interpreter
    {
    public:
        // Constructor
        Interpreter(MemoryManager &memory);

        // Execute method
        Object *executeMethod(Object *method, Object *receiver, std::vector<Object *> &args);

        // Execute compiled method directly
        TaggedValue executeCompiledMethod(const CompiledMethod &method);

        // Execute context
        Object *executeContext(MethodContext *context);

        // Main bytecode execution loop
        void executeLoop();

        // Bytecode dispatch
        void dispatch(Bytecode bytecode);

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
        void push(Object *value);
        Object *pop();
        Object *top();

        // Context access
        MethodContext *getCurrentContext() const { return activeContext; }
        void setCurrentContext(MethodContext *context) { activeContext = context; }

        // Memory manager access
        MemoryManager &getMemoryManager() { return memoryManager; }

        // TaggedValue message sending
        TaggedValue sendMessage(TaggedValue receiver, const std::string &selector, const std::vector<TaggedValue> &args);

        // Get object class for TaggedValue
        Class *getObjectClass(TaggedValue value);

    private:
        // Memory manager
        MemoryManager &memoryManager;

        // Current context and chunk
        MethodContext *activeContext = nullptr;
        StackChunk *currentChunk = nullptr;

        // Internal state
        bool executing = false;

        // Read 32-bit value from bytecode
        uint32_t readUInt32(size_t offset);

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
