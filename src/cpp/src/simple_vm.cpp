#include "simple_vm.h"
#include "bytecode.h"
#include "symbol.h"

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
    std::string selectorName = selector->getName();
    
    // Handle primitive operations based on selector name
    TaggedValue result;
    
    if (selectorName == "+") {
        if (argCount != 1) {
            throw std::runtime_error("Add expects exactly 1 argument");
        }
        result = performAdd(receiver, args[0]);
    } else if (selectorName == "-") {
        if (argCount != 1) {
            throw std::runtime_error("Subtract expects exactly 1 argument");
        }
        result = performSubtract(receiver, args[0]);
    } else if (selectorName == "*") {
        if (argCount != 1) {
            throw std::runtime_error("Multiply expects exactly 1 argument");
        }
        result = performMultiply(receiver, args[0]);
    } else if (selectorName == "/") {
        if (argCount != 1) {
            throw std::runtime_error("Divide expects exactly 1 argument");
        }
        result = performDivide(receiver, args[0]);
    } else if (selectorName == "<") {
        if (argCount != 1) {
            throw std::runtime_error("Less than expects exactly 1 argument");
        }
        result = performLessThan(receiver, args[0]);
    } else if (selectorName == ">") {
        if (argCount != 1) {
            throw std::runtime_error("Greater than expects exactly 1 argument");
        }
        result = performGreaterThan(receiver, args[0]);
    } else if (selectorName == "=") {
        if (argCount != 1) {
            throw std::runtime_error("Equal expects exactly 1 argument");
        }
        result = performEqual(receiver, args[0]);
    } else if (selectorName == "~=") {
        if (argCount != 1) {
            throw std::runtime_error("Not equal expects exactly 1 argument");
        }
        result = performNotEqual(receiver, args[0]);
    } else if (selectorName == "<=") {
        if (argCount != 1) {
            throw std::runtime_error("Less than or equal expects exactly 1 argument");
        }
        result = performLessThanOrEqual(receiver, args[0]);
    } else if (selectorName == ">=") {
        if (argCount != 1) {
            throw std::runtime_error("Greater than or equal expects exactly 1 argument");
        }
        result = performGreaterThanOrEqual(receiver, args[0]);
    } else {
        throw std::runtime_error("Unknown message selector: " + selectorName);
    }
    
    push(result);
}

void SimpleVM::handleReturn() {
    // Return instruction - execution will stop after this
    // The result should be on top of the stack
    instructionPointer_ = static_cast<uint32_t>(currentMethod_->getBytecodes().size()); // Force loop to exit
}

TaggedValue SimpleVM::performAdd(TaggedValue left, TaggedValue right) {
    if (left.isInteger() && right.isInteger()) {
        int32_t result = left.asInteger() + right.asInteger();
        return TaggedValue(result);
    }
    throw std::runtime_error("Add operation only supports integers in this simple VM");
}

TaggedValue SimpleVM::performSubtract(TaggedValue left, TaggedValue right) {
    if (left.isInteger() && right.isInteger()) {
        int32_t result = left.asInteger() - right.asInteger();
        return TaggedValue(result);
    }
    throw std::runtime_error("Subtract operation only supports integers in this simple VM");
}

TaggedValue SimpleVM::performMultiply(TaggedValue left, TaggedValue right) {
    if (left.isInteger() && right.isInteger()) {
        int32_t result = left.asInteger() * right.asInteger();
        return TaggedValue(result);
    }
    throw std::runtime_error("Multiply operation only supports integers in this simple VM");
}

TaggedValue SimpleVM::performDivide(TaggedValue left, TaggedValue right) {
    if (left.isInteger() && right.isInteger()) {
        if (right.asInteger() == 0) {
            throw std::runtime_error("Division by zero");
        }
        int32_t result = left.asInteger() / right.asInteger();
        return TaggedValue(result);
    }
    throw std::runtime_error("Divide operation only supports integers in this simple VM");
}

TaggedValue SimpleVM::performLessThan(TaggedValue left, TaggedValue right) {
    if (left.isInteger() && right.isInteger()) {
        bool result = left.asInteger() < right.asInteger();
        return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
    }
    throw std::runtime_error("Less than operation only supports integers in this simple VM");
}

TaggedValue SimpleVM::performGreaterThan(TaggedValue left, TaggedValue right) {
    if (left.isInteger() && right.isInteger()) {
        bool result = left.asInteger() > right.asInteger();
        return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
    }
    throw std::runtime_error("Greater than operation only supports integers in this simple VM");
}

TaggedValue SimpleVM::performEqual(TaggedValue left, TaggedValue right) {
    if (left.isInteger() && right.isInteger()) {
        bool result = left.asInteger() == right.asInteger();
        return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
    }
    throw std::runtime_error("Equal operation only supports integers in this simple VM");
}

TaggedValue SimpleVM::performNotEqual(TaggedValue left, TaggedValue right) {
    if (left.isInteger() && right.isInteger()) {
        bool result = left.asInteger() != right.asInteger();
        return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
    }
    throw std::runtime_error("Not equal operation only supports integers in this simple VM");
}

TaggedValue SimpleVM::performLessThanOrEqual(TaggedValue left, TaggedValue right) {
    if (left.isInteger() && right.isInteger()) {
        bool result = left.asInteger() <= right.asInteger();
        return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
    }
    throw std::runtime_error("Less than or equal operation only supports integers in this simple VM");
}

TaggedValue SimpleVM::performGreaterThanOrEqual(TaggedValue left, TaggedValue right) {
    if (left.isInteger() && right.isInteger()) {
        bool result = left.asInteger() >= right.asInteger();
        return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
    }
    throw std::runtime_error("Greater than or equal operation only supports integers in this simple VM");
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
