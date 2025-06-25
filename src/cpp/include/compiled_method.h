#pragma once

#include "tagged_value.h"
#include <vector>
#include <cstdint>
#include <string>

namespace smalltalk {

/**
 * A compiled method containing bytecode and literal values
 */
class CompiledMethod {
public:
    CompiledMethod() = default;
    
    // Add bytecode
    void addBytecode(uint8_t bytecode) {
        bytecodes_.push_back(bytecode);
    }
    
    // Add a 4-byte operand (little-endian)
    void addOperand(uint32_t operand) {
        bytecodes_.push_back(static_cast<uint8_t>(operand & 0xFF));
        bytecodes_.push_back(static_cast<uint8_t>((operand >> 8) & 0xFF));
        bytecodes_.push_back(static_cast<uint8_t>((operand >> 16) & 0xFF));
        bytecodes_.push_back(static_cast<uint8_t>((operand >> 24) & 0xFF));
    }
    
    // Add a literal value and return its index
    uint32_t addLiteral(TaggedValue value) {
        uint32_t index = static_cast<uint32_t>(literals_.size());
        literals_.push_back(value);
        return index;
    }
    
    // Getters
    const std::vector<uint8_t>& getBytecodes() const { return bytecodes_; }
    const std::vector<TaggedValue>& getLiterals() const { return literals_; }
    
    // Get literal by index
    TaggedValue getLiteral(uint32_t index) const {
        if (index >= literals_.size()) {
            throw std::runtime_error("Literal index out of bounds");
        }
        return literals_[index];
    }
    
    // Debug output
    std::string toString() const {
        std::string result = "CompiledMethod {\n";
        result += "  Bytecodes: [";
        for (size_t i = 0; i < bytecodes_.size(); ++i) {
            if (i > 0) result += ", ";
            result += std::to_string(static_cast<int>(bytecodes_[i]));
        }
        result += "]\n";
        
        result += "  Literals: [";
        for (size_t i = 0; i < literals_.size(); ++i) {
            if (i > 0) result += ", ";
            if (literals_[i].isInteger()) {
                result += std::to_string(literals_[i].asInteger());
            } else {
                result += "?";
            }
        }
        result += "]\n";
        result += "}";
        
        return result;
    }
    
private:
    std::vector<uint8_t> bytecodes_;
    std::vector<TaggedValue> literals_;
};

} // namespace smalltalk
