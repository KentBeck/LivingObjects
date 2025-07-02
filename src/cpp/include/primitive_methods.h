#pragma once

#include "compiled_method.h"
#include "tagged_value.h"
#include "symbol.h"

#include <functional>
#include <memory>

namespace smalltalk {

// Forward declarations
struct Object;
class Class;

/**
 * PrimitiveMethod represents a method implemented in C++ rather than Smalltalk bytecode.
 * These are used for basic operations like arithmetic, comparison, etc.
 */
class PrimitiveMethod : public CompiledMethod {
public:
    // Function signature for primitive implementations
    using PrimitiveFunction = std::function<TaggedValue(TaggedValue receiver, const std::vector<TaggedValue>& args, class Interpreter& interpreter)>;
    
    PrimitiveMethod(int primitiveNumber, PrimitiveFunction function);
    
    // Execute the primitive
    TaggedValue execute(TaggedValue receiver, const std::vector<TaggedValue>& args, class Interpreter& interpreter) const;
    
    // Get primitive number
    int getPrimitiveNumber() const { return primitiveNumber_; }
    
    // Override CompiledMethod methods
    std::string toString() const override;
    
private:
    int primitiveNumber_;
    PrimitiveFunction function_;
};

/**
 * PrimitiveRegistry manages all primitive methods in the system.
 */
class PrimitiveRegistry {
public:
    // Singleton access
    static PrimitiveRegistry& getInstance();
    
    // Register a primitive
    void registerPrimitive(int primitiveNumber, PrimitiveMethod::PrimitiveFunction function);
    
    // Create a primitive method
    std::shared_ptr<PrimitiveMethod> createPrimitiveMethod(int primitiveNumber) const;
    
    // Check if a primitive exists
    bool hasPrimitive(int primitiveNumber) const;
    
    // Initialize all core primitives
    void initializeCorePrimitives();
    
private:
    PrimitiveRegistry() = default;
    std::unordered_map<int, PrimitiveMethod::PrimitiveFunction> primitives_;
};

/**
 * Primitive method implementations for Integer class
 */
namespace IntegerPrimitives {
    // Arithmetic primitives
    TaggedValue add(TaggedValue receiver, const std::vector<TaggedValue>& args, class Interpreter& interpreter);
    TaggedValue subtract(TaggedValue receiver, const std::vector<TaggedValue>& args, class Interpreter& interpreter);
    TaggedValue multiply(TaggedValue receiver, const std::vector<TaggedValue>& args, class Interpreter& interpreter);
    TaggedValue divide(TaggedValue receiver, const std::vector<TaggedValue>& args, class Interpreter& interpreter);
    
    // Comparison primitives
    TaggedValue lessThan(TaggedValue receiver, const std::vector<TaggedValue>& args, class Interpreter& interpreter);
    TaggedValue greaterThan(TaggedValue receiver, const std::vector<TaggedValue>& args, class Interpreter& interpreter);
    TaggedValue equal(TaggedValue receiver, const std::vector<TaggedValue>& args, class Interpreter& interpreter);
    TaggedValue notEqual(TaggedValue receiver, const std::vector<TaggedValue>& args, class Interpreter& interpreter);
    TaggedValue lessThanOrEqual(TaggedValue receiver, const std::vector<TaggedValue>& args, class Interpreter& interpreter);
    TaggedValue greaterThanOrEqual(TaggedValue receiver, const std::vector<TaggedValue>& args, class Interpreter& interpreter);
    
    // Utility functions
    void checkArgumentCount(const std::vector<TaggedValue>& args, size_t expected);
    void checkIntegerReceiver(TaggedValue receiver);
    void checkIntegerArgument(TaggedValue arg, size_t index);
}

/**
 * Primitive numbers (following Smalltalk convention)
 */
namespace PrimitiveNumbers {
    constexpr int INTEGER_ADD = 1;
    constexpr int INTEGER_SUBTRACT = 2;
    constexpr int INTEGER_MULTIPLY = 3;
    constexpr int INTEGER_DIVIDE = 4;
    constexpr int INTEGER_LESS_THAN = 5;
    constexpr int INTEGER_GREATER_THAN = 6;
    constexpr int INTEGER_EQUAL = 7;
    constexpr int INTEGER_NOT_EQUAL = 8;
    constexpr int INTEGER_LESS_THAN_OR_EQUAL = 9;
    constexpr int INTEGER_GREATER_THAN_OR_EQUAL = 10;
}

/**
 * Utility functions for setting up Integer class with primitive methods
 */
namespace IntegerClassSetup {
    // Add all primitive methods to Integer class
    void addPrimitiveMethods(Class* integerClass);
    
    // Create a primitive method and add it to a class
    void addPrimitiveMethod(Class* clazz, const std::string& selector, int primitiveNumber);
}

} // namespace smalltalk