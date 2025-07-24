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
        Interpreter(MemoryManager &memory, SmalltalkImage &image);

        // Execute method
        Object *executeMethod(CompiledMethod *method, Object *receiver, std::vector<Object *> &args);

        // Execute compiled method directly
        // Suspicious that these are all so close
        TaggedValue executeCompiledMethod(const CompiledMethod &method);
        TaggedValue executeCompiledMethod(const CompiledMethod &method, MethodContext *context);
        TaggedValue executeMethodContext(MethodContext *context);
        TaggedValue executeMethodContext(MethodContext *context, CompiledMethod *method);

        // Core bytecode execution engine
        TaggedValue execute();

        // Bytecode operation helpers
        void pushLiteral();
        void pushSelf();
        void sendMessageBytecode();
        void createBlock();
        void pushTemporaryVariable();
        void storeTemporaryVariable();
        void popStack();
        void duplicate();
        TaggedValue returnStackTop();

        // Legacy bytecode handlers (deprecated)
        void handlePushTemporaryVariable(uint32_t offset);
        void handleStoreTemporaryVariable(uint32_t offset);
        void handlePop();
        void handleDuplicate();
        void handleCreateBlock(uint32_t literalIndex);

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
        SmalltalkImage &getImage() { return image; }

        // TaggedValue message sending
        TaggedValue sendMessage(TaggedValue receiver, const std::string &selector, const std::vector<TaggedValue> &args);

        // Get object class for TaggedValue
        Class *getObjectClass(TaggedValue value);

        // Exception handling
        bool findExceptionHandler(MethodContext *&handlerContext, int &handlerPC);
        void unwindToContext(MethodContext *targetContext);

    private:
        // Memory manager
        MemoryManager &memoryManager;
        SmalltalkImage &image;

        // Current context and chunk
        MethodContext *activeContext = nullptr;
        StackChunk *currentChunk = nullptr;


        // Helper for message sending
        Object *sendMessage(Object *receiver, Object *selector, std::vector<Object *> &args);

        // Simple evaluation method for testing (temporarily disabled)
        // TaggedValue evaluate(const std::string& expression);

        // Context switching
        void switchContext(MethodContext *newContext);

        // Bytecode reading helper
        uint32_t readUint32FromBytecode(const std::vector<uint8_t> &bytecodes, MethodContext *context);
    };

} // namespace smalltalk
