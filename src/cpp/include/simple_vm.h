#pragma once

#include "compiled_method.h"
#include "tagged_value.h"

#include <cstdint>
#include <vector>
#include <memory>

namespace smalltalk {

// Forward declarations
class Class;

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
    uint32_t stackPointer_ = 0;
    uint32_t instructionPointer_ = 0;
    const CompiledMethod* currentMethod_ = nullptr;
    
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
    
    // Method execution
    TaggedValue executeMethod(std::shared_ptr<CompiledMethod> method, TaggedValue receiver, const std::vector<TaggedValue>& args);
    

    
    // Helper functions
    uint32_t readOperand();
    void checkStackUnderflow(int required) const;
};

} // namespace smalltalk
