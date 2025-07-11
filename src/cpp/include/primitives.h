#pragma once

#include "tagged_value.h"
#include "object.h"
#include <vector>
#include <functional>
#include <unordered_map>

namespace smalltalk {

// Forward declarations
class Interpreter;
class Class;

/**
 * Primitive function type
 * Takes receiver, arguments, and interpreter context
 * Returns a TaggedValue result
 */
using PrimitiveFunction = std::function<TaggedValue(TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter)>;

/**
 * Exception thrown when a primitive fails and should fall back to Smalltalk code
 */
class PrimitiveFailure : public std::exception {
public:
    explicit PrimitiveFailure(std::string message) : message_(std::move(message)) {}
    const char* what() const noexcept override { return message_.c_str(); }
    
private:
    std::string message_;
};

/**
 * Primitive numbers for standard Smalltalk primitives
 */
namespace PrimitiveNumbers {
    // Object primitives
    const int NEW = 70;               // Object new
    const int BASIC_NEW = 71;         // Object basicNew
    const int BASIC_NEW_SIZE = 72;    // Object basicNew: size
    const int IDENTITY_HASH = 75;     // Object identityHash
    const int CLASS = 111;            // Object class
    
    // Integer arithmetic
    const int SMALL_INT_ADD = 1;      // SmallInteger +
    const int SMALL_INT_SUB = 2;      // SmallInteger -
    const int SMALL_INT_MUL = 9;      // SmallInteger *
    const int SMALL_INT_DIV = 10;     // SmallInteger /
    const int SMALL_INT_MOD = 11;     // SmallInteger //
    
    // Integer comparison
    const int SMALL_INT_LT = 3;       // SmallInteger <
    const int SMALL_INT_GT = 4;       // SmallInteger >
    const int SMALL_INT_LE = 5;       // SmallInteger <=
    const int SMALL_INT_GE = 6;       // SmallInteger >=
    const int SMALL_INT_EQ = 7;       // SmallInteger =
    const int SMALL_INT_NE = 8;       // SmallInteger ~=
    
    // Block evaluation
    const int BLOCK_VALUE = 201;      // Block value
    const int BLOCK_VALUE_ARG = 202;  // Block value:
    
    // Array primitives
    const int ARRAY_AT = 60;          // Array at:
    const int ARRAY_AT_PUT = 61;      // Array at:put:
    const int ARRAY_SIZE = 62;        // Array size
    
    // String primitives
    const int STRING_AT = 63;         // String at:
    const int STRING_AT_PUT = 64;     // String at:put:
    const int STRING_SIZE = 66;       // String size
    const int STRING_CONCAT = 65;     // String ,
}

/**
 * PrimitiveRegistry manages all primitive functions in the system
 */
class PrimitiveRegistry {
public:
    // Singleton access
    static PrimitiveRegistry& getInstance();
    
    // Register a primitive function
    void registerPrimitive(int primitiveNumber, PrimitiveFunction function);
    
    // Call a primitive
    TaggedValue callPrimitive(int primitiveNumber, TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter);
    
    // Check if a primitive exists
    bool hasPrimitive(int primitiveNumber) const;
    
    // Get all registered primitive numbers
    std::vector<int> getAllPrimitiveNumbers() const;
    
    // Clear all primitives
    void clear();
    
    // Initialize core primitives
    void initializeCorePrimitives();
    
private:
    PrimitiveRegistry() = default;
    std::unordered_map<int, PrimitiveFunction> primitives_;
};

// Convenience functions for primitive operations
namespace Primitives {
    // Initialize all core primitives
    void initialize();
    
    // Call a primitive with error handling
    TaggedValue callPrimitive(int primitiveNumber, TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter);
    
    // Helper to check argument count
    void checkArgumentCount(const std::vector<TaggedValue>& args, size_t expected, const std::string& primitiveName);
    
    // Helper to check receiver type
    void checkReceiverType(TaggedValue receiver, ObjectType expectedType, const std::string& primitiveName);
    
    // Helper to ensure receiver is a class
    Class* ensureReceiverIsClass(TaggedValue receiver, const std::string& primitiveName);
    
    // Helper to ensure receiver is an object
    Object* ensureReceiverIsObject(TaggedValue receiver, const std::string& primitiveName);
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