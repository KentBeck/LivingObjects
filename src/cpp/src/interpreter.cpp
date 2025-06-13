#include "interpreter.h"
#include <stdexcept>
#include <cstring>

namespace smalltalk {

Interpreter::Interpreter(MemoryManager& memory)
    : memoryManager(memory),
      activeContext(nullptr),
      currentChunk(nullptr),
      executing(false) {
    // Initialize the stack chunk
    currentChunk = memoryManager.allocateStackChunk(1024);
}

Object* Interpreter::executeMethod(Object* method, Object* receiver, std::vector<Object*>& args) {
    // Create a new method context
    uint32_t methodObj = static_cast<uint32_t>(reinterpret_cast<uintptr_t>(method));
    MethodContext* context = memoryManager.allocateMethodContext(10 + args.size(), methodObj, receiver, nullptr);
    
    // Copy arguments to the context
    Object** slots = reinterpret_cast<Object**>(reinterpret_cast<char*>(context) + sizeof(MethodContext));
    for (size_t i = 0; i < args.size(); i++) {
        slots[i] = args[i];
    }
    
    // Set up stack pointer
    context->stackPointer = reinterpret_cast<Object*>(slots + args.size());
    
    // Execute the context
    return executeContext(context);
}

Object* Interpreter::executeContext(MethodContext* context) {
    // Save current context
    MethodContext* previousContext = activeContext;
    
    // Set new active context
    activeContext = context;
    
    // Execute bytecodes until context returns
    executeLoop();
    
    // Get the return value (top of stack)
    Object* result = top();
    
    // Restore previous context
    activeContext = previousContext;
    
    return result;
}

void Interpreter::executeLoop() {
    executing = true;
    
    while (executing && activeContext) {
        // For a basic implementation, we'll just pretend we're executing bytecodes
        // In a real implementation, we would fetch bytecodes from the method
        
        // For the minimal example, let's just execute a simple sequence
        // that pushes a value and returns
        handlePushSelf();
        handleReturnStackTop();
        
        // Stop execution after this simple sequence
        executing = false;
    }
}

void Interpreter::dispatch(Bytecode bytecode) {
    // Get current instruction pointer
    size_t ip = activeContext->instructionPointer;
    
    // Update instruction pointer for next instruction
    activeContext->instructionPointer += getInstructionSize(bytecode);
    
    // Dispatch based on bytecode
    switch (bytecode) {
        case Bytecode::PUSH_LITERAL:
            handlePushLiteral(readUInt32(ip + 1));
            break;
        case Bytecode::PUSH_INSTANCE_VARIABLE:
            handlePushInstanceVariable(readUInt32(ip + 1));
            break;
        case Bytecode::PUSH_TEMPORARY_VARIABLE:
            handlePushTemporaryVariable(readUInt32(ip + 1));
            break;
        case Bytecode::PUSH_SELF:
            handlePushSelf();
            break;
        case Bytecode::STORE_INSTANCE_VARIABLE:
            handleStoreInstanceVariable(readUInt32(ip + 1));
            break;
        case Bytecode::STORE_TEMPORARY_VARIABLE:
            handleStoreTemporaryVariable(readUInt32(ip + 1));
            break;
        case Bytecode::SEND_MESSAGE:
            handleSendMessage(readUInt32(ip + 1), readUInt32(ip + 5));
            break;
        case Bytecode::RETURN_STACK_TOP:
            handleReturnStackTop();
            break;
        case Bytecode::JUMP:
            handleJump(readUInt32(ip + 1));
            break;
        case Bytecode::JUMP_IF_TRUE:
            handleJumpIfTrue(readUInt32(ip + 1));
            break;
        case Bytecode::JUMP_IF_FALSE:
            handleJumpIfFalse(readUInt32(ip + 1));
            break;
        case Bytecode::POP:
            handlePop();
            break;
        case Bytecode::DUPLICATE:
            handleDuplicate();
            break;
        case Bytecode::CREATE_BLOCK:
            handleCreateBlock(readUInt32(ip + 1), readUInt32(ip + 5), readUInt32(ip + 9));
            break;
        case Bytecode::EXECUTE_BLOCK:
            handleExecuteBlock(readUInt32(ip + 1));
            break;
        default:
            throw std::runtime_error("Unknown bytecode");
    }
}

uint32_t Interpreter::readUInt32(size_t offset) {
    // In a real implementation, this would read from actual bytecodes
    // For this stub, just return a placeholder value
    (void)offset; // Suppress unused parameter warning
    return 0;
}

void Interpreter::push(Object* value) {
    // Push value onto the active context's stack
    Object** sp = reinterpret_cast<Object**>(activeContext->stackPointer);
    *sp = value;
    activeContext->stackPointer = reinterpret_cast<Object*>(sp + 1);
}

Object* Interpreter::pop() {
    // Pop value from the active context's stack
    Object** sp = reinterpret_cast<Object**>(activeContext->stackPointer);
    sp--;
    Object* value = *sp;
    activeContext->stackPointer = reinterpret_cast<Object*>(sp);
    return value;
}

Object* Interpreter::top() {
    // Get the top value from the active context's stack
    Object** sp = reinterpret_cast<Object**>(activeContext->stackPointer);
    return *(sp - 1);
}

// Bytecode handler implementations
void Interpreter::handlePushLiteral(uint32_t index) {
    // In a real implementation, this would get the literal from the method's literal array
    // For now, we just create a new object
    (void)index; // Suppress unused parameter warning
    Object* literal = memoryManager.allocateObject(ObjectType::OBJECT, 0);
    push(literal);
}

void Interpreter::handlePushInstanceVariable(uint32_t offset) {
    // Get the instance variable from the receiver at the given offset
    // For now, just push a new object
    (void)offset; // Suppress unused parameter warning
    Object* value = memoryManager.allocateObject(ObjectType::OBJECT, 0);
    push(value);
}

void Interpreter::handlePushTemporaryVariable(uint32_t offset) {
    // Get the temporary variable at the given offset
    Object** slots = reinterpret_cast<Object**>(reinterpret_cast<char*>(activeContext) + sizeof(MethodContext));
    push(slots[offset]);
}

void Interpreter::handlePushSelf() {
    // Push the receiver onto the stack
    push(activeContext->self);
}

void Interpreter::handleStoreInstanceVariable(uint32_t offset) {
    // Store the top of the stack into the instance variable at the given offset
    // Not implemented in this basic version
    (void)offset; // Suppress unused parameter warning
}

void Interpreter::handleStoreTemporaryVariable(uint32_t offset) {
    // Store the top of the stack into the temporary variable at the given offset
    Object* value = pop();
    Object** slots = reinterpret_cast<Object**>(reinterpret_cast<char*>(activeContext) + sizeof(MethodContext));
    slots[offset] = value;
    push(value); // Leave the value on the stack
}

void Interpreter::handleSendMessage(uint32_t selectorIndex, uint32_t argCount) {
    // Not fully implemented in this basic version
    // In a real implementation, this would look up the method and execute it
    (void)selectorIndex; // Suppress unused parameter warning
    
    // Pop arguments
    std::vector<Object*> args;
    for (uint32_t i = 0; i < argCount; i++) {
        args.push_back(pop());
    }
    
    // Pop receiver
    Object* receiver = pop();
    
    // Create a selector object
    Object* selector = memoryManager.allocateObject(ObjectType::SYMBOL, 0);
    
    // Send the message
    Object* result = sendMessage(receiver, selector, args);
    
    // Push the result
    push(result);
}

void Interpreter::handleReturnStackTop() {
    // Return from the current context
    Object* result = top();
    
    // Set the sender as the new active context
    MethodContext* sender = reinterpret_cast<MethodContext*>(activeContext->sender);
    activeContext = sender;
    
    // If there is a sender, push the result onto its stack
    if (activeContext) {
        push(result);
    }
}

void Interpreter::handleJump(uint32_t target) {
    // Jump to the target instruction
    activeContext->instructionPointer = target;
}

void Interpreter::handleJumpIfTrue(uint32_t target) {
    // Pop the condition
    Object* condition = pop();
    
    // Check if the condition is true (simplified)
    bool isTrue = (condition != nullptr);
    
    // Jump if true
    if (isTrue) {
        activeContext->instructionPointer = target;
    }
}

void Interpreter::handleJumpIfFalse(uint32_t target) {
    // Pop the condition
    Object* condition = pop();
    
    // Check if the condition is false (simplified)
    bool isFalse = (condition == nullptr);
    
    // Jump if false
    if (isFalse) {
        activeContext->instructionPointer = target;
    }
}

void Interpreter::handlePop() {
    // Pop the top value from the stack
    pop();
}

void Interpreter::handleDuplicate() {
    // Duplicate the top value on the stack
    Object* value = top();
    push(value);
}

void Interpreter::handleCreateBlock(uint32_t bytecodeSize, uint32_t literalCount, uint32_t tempVarCount) {
    // Create a block object (simplified)
    (void)bytecodeSize; // Suppress unused parameter warning
    (void)literalCount; // Suppress unused parameter warning
    (void)tempVarCount; // Suppress unused parameter warning
    
    Object* block = memoryManager.allocateObject(ObjectType::OBJECT, 0);
    
    // Push the block onto the stack
    push(block);
}

void Interpreter::handleExecuteBlock(uint32_t argCount) {
    // Not implemented in this basic version
    (void)argCount; // Suppress unused parameter warning
}

Object* Interpreter::sendMessage(Object* receiver, Object* selector, std::vector<Object*>& args) {
    // Not fully implemented in this basic version
    // In a real implementation, this would look up the method and execute it
    (void)receiver; // Suppress unused parameter warning
    (void)selector; // Suppress unused parameter warning
    (void)args; // Suppress unused parameter warning
    
    // For now, just return a new object
    return memoryManager.allocateObject(ObjectType::OBJECT, 0);
}

Object* Interpreter::lookupMethod(Object* receiver, Object* selector) {
    // Not implemented in this basic version
    (void)receiver; // Suppress unused parameter warning
    (void)selector; // Suppress unused parameter warning
    return nullptr;
}

void Interpreter::switchContext(MethodContext* newContext) {
    // Set the new active context
    activeContext = newContext;
}

} // namespace smalltalk