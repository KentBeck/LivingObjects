#pragma once

#include "tagged_value.h"
#include <string>
#include <stdexcept>

namespace smalltalk {

/**
 * SimpleInterpreter provides basic evaluation of Smalltalk expressions.
 * 
 * This is the minimal interpreter needed to get the expression "3" working.
 * It focuses on immediate values and basic object identity.
 */
class SimpleInterpreter {
public:
    SimpleInterpreter() = default;
    
    /**
     * Evaluate a simple Smalltalk expression.
     * Currently supports:
     * - Integer literals (e.g., "3", "42", "-17")
     * - Special values ("nil", "true", "false")
     */
    TaggedValue evaluate(const std::string& expression);
    
private:
    /**
     * Parse an integer from a string.
     * Returns true if successful, false otherwise.
     */
    bool tryParseInteger(const std::string& str, int32_t& result);
    
    /**
     * Trim whitespace from a string.
     */
    std::string trim(const std::string& str);
};

} // namespace smalltalk
