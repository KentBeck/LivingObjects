#include "simple_vm.h"
#include "bytecode.h"
#include "symbol.h"
#include "smalltalk_object.h"
#include "smalltalk_class.h"
#include "primitive_methods.h"

#include <iostream>
#include <stdexcept>

namespace smalltalk {

SimpleVM::SimpleVM() {
    // Reserve space for the stack
    stack_.reserve(1000);
}

TaggedValue SimpleVM::execute(const CompiledMethod& method) {
    // Initialize execution state
    currentMethod_ = &method;
    instructionPointer_ = 0;
    stackPointer_ = 0;
    stack_.clear();
    
    // Execute bytecodes until we hit a return
    while (instructionPointer_ < method.getBytecodes().size()) {
        executeBytecode();
    }
    
    // Return the top of stack (should be the result)
    if (stackPointer_ == 0) {
        // No result, return nil
        return TaggedValue::nil();
    }
    
    return stack_[stackPointer_ - 1];
}

void SimpleVM::push(TaggedValue value) {
    if (stack_.size() <= stackPointer_) {
        stack_.resize(stackPointer_ + 1);
    }
    stack_[stackPointer_] = value;
    stackPointer_++;
}

TaggedValue SimpleVM::pop() {
    checkStackUnderflow(1);
    stackPointer_--;
    return stack_[stackPointer_];
}

TaggedValue SimpleVM::top() const {
    checkStackUnderflow(1);
    return stack_[stackPointer_ - 1];
}

void SimpleVM::executeBytecode() {
    if (instructionPointer_ >= currentMethod_->getBytecodes().size()) {
        throw std::runtime_error("Instruction pointer out of bounds");
    }
    
    uint8_t opcode = currentMethod_->getBytecodes()[instructionPointer_];
    instructionPointer_++;
    
    Bytecode bytecode = static_cast<Bytecode>(opcode);
    
    switch (bytecode) {
        case Bytecode::PUSH_LITERAL:
            handlePushLiteral();
            break;
            
        case Bytecode::SEND_MESSAGE:
            handleSendMessage();
            break;
            
        case Bytecode::RETURN_STACK_TOP:
            handleReturn();
            break;
            
        default:
            throw std::runtime_error("Unimplemented bytecode: " + std::string(getBytecodeString(bytecode)));
    }
}

void SimpleVM::handlePushLiteral() {
    uint32_t literalIndex = readOperand();
    TaggedValue literal = currentMethod_->getLiteral(literalIndex);
    push(literal);
}

void SimpleVM::handleSendMessage() {
    uint32_t selectorIndex = readOperand();
    uint32_t argCount = readOperand();
    
    // Pop arguments and receiver
    checkStackUnderflow(static_cast<int>(argCount + 1)); // +1 for receiver
    
    std::vector<TaggedValue> args;
    for (uint32_t i = 0; i < argCount; i++) {
        args.insert(args.begin(), pop()); // Insert at beginning to maintain order
    }
    
    TaggedValue receiver = pop();
    
    // Get the selector symbol from the literal table
    TaggedValue selectorValue = currentMethod_->getLiteral(selectorIndex);
    if (!selectorValue.isPointer()) {
        throw std::runtime_error("Selector must be a symbol");
    }
    
    Symbol* selector = selectorValue.asSymbol();
    
    // Get the receiver's class
    Class* receiverClass = receiver.getClass();
    if (receiverClass == nullptr) {
        throw std::runtime_error("Receiver has no class");
    }
    
    // Look up the method in the receiver's class
    auto method = receiverClass->lookupMethod(selector);
    if (method == nullptr) {
        throw std::runtime_error("Method not found: " + receiverClass->getName() + ">>" + selector->getName());
    }
    
    // Execute the method
    TaggedValue result = executeMethod(method, receiver, args);
    
    push(result);
}

void SimpleVM::handleReturn() {
    // Return instruction - execution will stop after this
    // The result should be on top of the stack
    instructionPointer_ = static_cast<uint32_t>(currentMethod_->getBytecodes().size()); // Force loop to exit
}

TaggedValue SimpleVM::executeMethod(std::shared_ptr<CompiledMethod> method, TaggedValue receiver, const std::vector<TaggedValue>& args) {
    // Check if this is a primitive method
    auto primitiveMethod = std::dynamic_pointer_cast<PrimitiveMethod>(method);
    if (primitiveMethod) {
        // Execute primitive directly
        return primitiveMethod->execute(receiver, args);
    }
    
    // For now, we only support primitive methods
    // In a full implementation, we would:
    // 1. Create a new execution context
    // 2. Set up the method's locals and temporaries
    // 3. Execute the method's bytecode
    // 4. Handle method return
    
    throw std::runtime_error("Non-primitive methods not yet implemented");
}



uint32_t SimpleVM::readOperand() {
    if (instructionPointer_ + 4 > currentMethod_->getBytecodes().size()) {
        throw std::runtime_error("Not enough bytes for operand");
    }
    
    const auto& bytecodes = currentMethod_->getBytecodes();
    
    // Read 4 bytes in little-endian format
    uint32_t operand = 
        static_cast<uint32_t>(bytecodes[instructionPointer_]) |
        (static_cast<uint32_t>(bytecodes[instructionPointer_ + 1]) << 8) |
        (static_cast<uint32_t>(bytecodes[instructionPointer_ + 2]) << 16) |
        (static_cast<uint32_t>(bytecodes[instructionPointer_ + 3]) << 24);
    
    instructionPointer_ += 4;
    return operand;
}

void SimpleVM::checkStackUnderflow(int required) const {
    if (static_cast<int>(stackPointer_) < required) {
        throw std::runtime_error("Stack underflow - need " + std::to_string(required) + 
                                " items but only have " + std::to_string(stackPointer_));
    }
}

} // namespace smalltalk
