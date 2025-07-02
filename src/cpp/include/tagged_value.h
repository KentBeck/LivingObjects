#pragma once

#include <cmath>
#include <cstdint>
#include <iostream>
namespace smalltalk {

// Forward declarations
class Symbol;
struct Object;
class Class;

/**
 * TaggedValue provides an efficient representation for Smalltalk values.
 * 
 * It uses pointer tagging to represent immediate values (integers, floats,
 * booleans, nil) without heap allocation, while maintaining compatibility
 * with heap-allocated objects.
 * 
 * Tag bits:
 *   00 = Pointer to heap object
 *   01 = Special value (nil, true, false)
 *   10 = Float (stored in-line when possible)
 *   11 = SmallInteger (31-bit signed integer)
 */
class TaggedValue {
public:
    // Tag bit definitions
    static constexpr uintptr_t TAG_MASK        = 0x03;  // Bottom 2 bits for tag
    static constexpr uintptr_t POINTER_TAG     = 0x00;  // 00 = Pointer
    static constexpr uintptr_t SPECIAL_TAG     = 0x01;  // 01 = Special (nil/true/false)
    static constexpr uintptr_t FLOAT_TAG       = 0x02;  // 10 = Float
    static constexpr uintptr_t INTEGER_TAG     = 0x03;  // 11 = Integer
    
    // Special value definitions
    static constexpr uintptr_t SPECIAL_NIL     = (0 << 2) | SPECIAL_TAG;
    static constexpr uintptr_t SPECIAL_TRUE    = (1 << 2) | SPECIAL_TAG;
    static constexpr uintptr_t SPECIAL_FALSE   = (2 << 2) | SPECIAL_TAG;
    
    // Constructors
    TaggedValue() : value(SPECIAL_NIL) {}
    
    // Create from pointer
    explicit TaggedValue(void* ptr) {
        // Ensure pointer is aligned to at least 4 bytes
        if (reinterpret_cast<uintptr_t>(ptr) & TAG_MASK) {
            throw std::runtime_error("Pointer not properly aligned for tagging");
        }
        value = reinterpret_cast<uintptr_t>(ptr) | POINTER_TAG;
    }
    
    // Create from Symbol pointer
    explicit TaggedValue(Symbol* symbol) {
        // Ensure pointer is aligned to at least 4 bytes
        if (reinterpret_cast<uintptr_t>(symbol) & TAG_MASK) {
            throw std::runtime_error("Symbol pointer not properly aligned for tagging");
        }
        value = reinterpret_cast<uintptr_t>(symbol) | POINTER_TAG;
    }
    
    // Create from Object pointer
    explicit TaggedValue(Object* object) {
        // Ensure pointer is aligned to at least 4 bytes
        if (reinterpret_cast<uintptr_t>(object) & TAG_MASK) {
            throw std::runtime_error("Object pointer not properly aligned for tagging");
        }
        value = reinterpret_cast<uintptr_t>(object) | POINTER_TAG;
    }
    
    // Create from integer
    explicit TaggedValue(int32_t intValue) {
        // Shift left to make room for tag bits, then apply tag
        value = (static_cast<uintptr_t>(intValue) << 2) | INTEGER_TAG;
    }
    
    // Create from double (may need to allocate on heap if doesn't fit in tagged value)
    explicit TaggedValue(double doubleValue) {
        // For now, just encode a subset of floats that can be stored inline
        // A full implementation would allocate a heap object for floats that don't fit
        if (doubleValue == 0.0) {
            // Encode 0.0 directly
            value = (0ULL << 2) | FLOAT_TAG;
        } else if (doubleValue == 1.0) {
            // Encode 1.0 directly
            value = (1ULL << 2) | FLOAT_TAG;
        } else if (doubleValue == -1.0) {
            // Encode -1.0 directly
            value = (2ULL << 2) | FLOAT_TAG;
        } else {
            // In a real implementation, we would allocate a heap object
            // For this simplified version, we just throw an error
            throw std::runtime_error("Float value cannot be encoded as tagged value");
        }
    }
    
    // Special value constructors
    static TaggedValue nil() { return TaggedValue(SPECIAL_NIL); }
    static TaggedValue trueValue() { return TaggedValue(SPECIAL_TRUE); }
    static TaggedValue falseValue() { return TaggedValue(SPECIAL_FALSE); }
    
    // Type checking
    bool isPointer() const { return (value & TAG_MASK) == POINTER_TAG; }
    bool isSpecial() const { return (value & TAG_MASK) == SPECIAL_TAG; }
    bool isFloat() const { return (value & TAG_MASK) == FLOAT_TAG; }
    bool isInteger() const { return (value & TAG_MASK) == INTEGER_TAG; }
    
    bool isNil() const { return value == SPECIAL_NIL; }
    bool isTrue() const { return value == SPECIAL_TRUE; }
    bool isFalse() const { return value == SPECIAL_FALSE; }
    bool isBoolean() const { return isTrue() || isFalse(); }
    
    // Value extraction
    void* asPointer() const {
        if (!isPointer()) {
            throw std::runtime_error("Tagged value is not a pointer");
        }
        return reinterpret_cast<void*>(value & ~TAG_MASK);
    }
    
    Symbol* asSymbol() const {
        if (!isPointer()) {
            throw std::runtime_error("Tagged value is not a symbol");
        }
        return reinterpret_cast<Symbol*>(value & ~TAG_MASK);
    }
    
    Object* asObject() const {
        if (!isPointer()) {
            throw std::runtime_error("Tagged value is not an object");
        }
        return reinterpret_cast<Object*>(value & ~TAG_MASK);
    }
    
    int32_t asInteger() const {
        if (!isInteger()) {
            throw std::runtime_error("Tagged value is not an integer");
        }
        // Shift right to remove tag bits, then sign-extend
        return static_cast<int32_t>(value >> 2);
    }
    
    double asFloat() const {
        if (!isFloat()) {
            throw std::runtime_error("Tagged value is not a float");
        }
        
        // Extract the encoded value
        uintptr_t encoded = value >> 2;
        
        // Decode based on encoded value
        if (encoded == 0) return 0.0;
        if (encoded == 1) return 1.0;
        if (encoded == 2) return -1.0;
        
        // Should not reach here in simplified implementation
        throw std::runtime_error("Invalid encoded float value");
    }
    
    bool asBoolean() const {
        if (!isBoolean()) {
            throw std::runtime_error("Tagged value is not a boolean");
        }
        return isTrue();
    }
    
    // Type checking for objects
    bool isObjectOfClass(Class* clazz) const;
    
    // Get the class of this object (for tagged values that represent objects)
    Class* getClass() const;
    
    // Comparison
    bool operator==(const TaggedValue& other) const {
        return value == other.value;
    }
    
    bool operator!=(const TaggedValue& other) const {
        return value != other.value;
    }
    
    // Raw value access (for internal use)
    uintptr_t rawValue() const { return value; }
    
private:
    // Constructor from raw value (private)
    explicit TaggedValue(uintptr_t rawValue) : value(rawValue) {}
    
    // The raw tagged value
    uintptr_t value;
};

// Output stream operator for debugging
inline std::ostream& operator<<(std::ostream& os, const TaggedValue& tv) {
    if (tv.isInteger()) {
        os << "Integer(" << tv.asInteger() << ")";
    } else if (tv.isFloat()) {
        os << "Float(" << tv.asFloat() << ")";
    } else if (tv.isNil()) {
        os << "nil";
    } else if (tv.isTrue()) {
        os << "true";
    } else if (tv.isFalse()) {
        os << "false";
    } else if (tv.isPointer()) {
        os << "Object@" << tv.asPointer();
    } else {
        os << "Unknown";
    }
    return os;
}

} // namespace smalltalk