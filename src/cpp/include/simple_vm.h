#pragma once

#include "compiled_method.h"
#include "tagged_value.h"
#include <vector>
#include <cstdint>

namespace smalltalk {

/**
 * Simple VM for executing bytecode
 */
class SimpleVM {
public:
    SimpleVM();
    
    // Execute a compiled method and return the result
    TaggedValue execute(const CompiledMethod& method);
    
private:
    // Execution state
    std::vector<TaggedValue> stack_;
    uint32_t stackPointer_;
    uint32_t instructionPointer_;
    const CompiledMethod* currentMethod_;
    
    // Stack operations
    void push(TaggedValue value);
    TaggedValue pop();
    TaggedValue top() const;
    
    // Bytecode execution
    void executeBytecode();
    
    // Instruction handlers
    void handlePushLiteral();
    void handleSendMessage();
    void handleReturn();
    
    // Primitive operations
    TaggedValue performAdd(TaggedValue left, TaggedValue right);
    TaggedValue performSubtract(TaggedValue left, TaggedValue right);
    TaggedValue performMultiply(TaggedValue left, TaggedValue right);
    TaggedValue performDivide(TaggedValue left, TaggedValue right);
    
    // Helper functions
    uint32_t readOperand();
    void checkStackUnderflow(int required) const;
};

} // namespace smalltalk
