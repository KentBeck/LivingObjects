#pragma once

#include "bytecode.h"
#include "object.h"
#include "context.h"
#include "memory_manager.h"
#include <cstdint>
#include <vector>

namespace smalltalk {

class Interpreter {
public:
    // Constructor
    Interpreter(MemoryManager& memory);
    
    // Execute method
    Object* executeMethod(Object* method, Object* receiver, std::vector<Object*>& args);
    
    // Execute context
    Object* executeContext(MethodContext* context);
    
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
    void push(Object* value);
    Object* pop();
    Object* top();
    
private:
    // Memory manager
    MemoryManager& memoryManager;
    
    // Current context and chunk
    MethodContext* activeContext;
    StackChunk* currentChunk;
    
    // Internal state
    bool executing;
    
    // Read 32-bit value from bytecode
    uint32_t readUInt32(size_t offset);
    
    // Helper for message sending
    Object* sendMessage(Object* receiver, Object* selector, std::vector<Object*>& args);
    
    // Method lookup
    Object* lookupMethod(Object* receiver, Object* selector);
    
    // Context switching
    void switchContext(MethodContext* newContext);
};

} // namespace smalltalk