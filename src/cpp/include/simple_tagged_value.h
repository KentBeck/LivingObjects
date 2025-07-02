#pragma once

#include <cstdint>
#include <cassert>
#include <iostream>
#include <string>

namespace smalltalk {

/**
 * Simplified TaggedValue for Smalltalk immediate values
 * 
 * Key Design Principles:
 * 1. ONLY for immediate values (nil, true, false, integers, floats)
 * 2. NEVER for heap objects (use Object* directly for those)
 * 3. Simple discriminated union - no bit manipulation tricks
 * 4. Easy to debug, extend, and understand
 * 5. Type-safe construction and access
 * 
 * This eliminates the complex tagged pointer schemes and makes
 * the distinction clear: TaggedValue = immediate, Object* = heap object
 */
class TaggedValue {
public:
    enum class Type : uint8_t {
        NIL = 0,
        BOOLEAN,
        INTEGER, 
        FLOAT
    };
    
    // === FACTORY METHODS (preferred way to construct) ===
    
    static TaggedValue nil() {
        return TaggedValue(Type::NIL, 0);
    }
    
    static TaggedValue boolean(bool value) {
        return TaggedValue(Type::BOOLEAN, value ? 1 : 0);
    }
    
    static TaggedValue true_value() {
        return boolean(true);
    }
    
    static TaggedValue false_value() {
        return boolean(false);
    }
    
    static TaggedValue integer(int32_t value) {
        return TaggedValue(Type::INTEGER, static_cast<uint32_t>(value));
    }
    
    static TaggedValue float_value(float value) {
        TaggedValue result(Type::FLOAT, 0);
        result.float_val_ = value;
        return result;
    }
    
    // === CONSTRUCTORS ===
    
    // Default constructor creates nil
    TaggedValue() : type_(Type::NIL), int_val_(0) {}
    
    // Copy constructor
    TaggedValue(const TaggedValue& other) : type_(other.type_), int_val_(other.int_val_) {}
    
    // Assignment operator  
    TaggedValue& operator=(const TaggedValue& other) {
        if (this != &other) {
            type_ = other.type_;
            int_val_ = other.int_val_;
        }
        return *this;
    }
    
    // === TYPE CHECKING ===
    
    Type type() const { return type_; }
    
    bool is_nil() const { return type_ == Type::NIL; }
    bool is_boolean() const { return type_ == Type::BOOLEAN; }  
    bool is_integer() const { return type_ == Type::INTEGER; }
    bool is_float() const { return type_ == Type::FLOAT; }
    
    bool is_true() const { 
        return type_ == Type::BOOLEAN && int_val_ != 0; 
    }
    
    bool is_false() const { 
        return type_ == Type::BOOLEAN && int_val_ == 0; 
    }
    
    bool is_number() const { 
        return is_integer() || is_float(); 
    }
    
    // === VALUE EXTRACTION ===
    
    bool as_boolean() const {
        assert(is_boolean());
        return int_val_ != 0;
    }
    
    int32_t as_integer() const {
        assert(is_integer());
        return static_cast<int32_t>(int_val_);
    }
    
    float as_float() const {
        assert(is_float());
        return float_val_;
    }
    
    // Safe extraction with default values
    bool as_boolean_or(bool default_val) const {
        return is_boolean() ? as_boolean() : default_val;
    }
    
    int32_t as_integer_or(int32_t default_val) const {
        return is_integer() ? as_integer() : default_val;
    }
    
    float as_float_or(float default_val) const {
        return is_float() ? as_float() : default_val;
    }
    
    // === COMPARISON ===
    
    bool operator==(const TaggedValue& other) const {
        if (type_ != other.type_) return false;
        
        switch (type_) {
            case Type::NIL:
                return true;  // All nils are equal
            case Type::BOOLEAN:
            case Type::INTEGER:
                return int_val_ == other.int_val_;
            case Type::FLOAT:
                return float_val_ == other.float_val_;
            default:
                return false;
        }
    }
    
    bool operator!=(const TaggedValue& other) const {
        return !(*this == other);
    }
    
    // === ARITHMETIC SUPPORT ===
    
    // Convert to double for mixed arithmetic
    double as_number() const {
        if (is_integer()) {
            return static_cast<double>(as_integer());
        } else if (is_float()) {
            return static_cast<double>(as_float());
        } else {
            assert(false && "Not a number");
            return 0.0;
        }
    }
    
    // Promote integer to float if needed for arithmetic
    TaggedValue to_float() const {
        if (is_integer()) {
            return float_value(static_cast<float>(as_integer()));
        } else if (is_float()) {
            return *this;
        } else {
            assert(false && "Cannot convert to float");
            return TaggedValue();
        }
    }
    
    // === DEBUGGING SUPPORT ===
    
    std::string to_string() const {
        switch (type_) {
            case Type::NIL:
                return "nil";
            case Type::BOOLEAN:
                return as_boolean() ? "true" : "false";
            case Type::INTEGER:
                return std::to_string(as_integer());
            case Type::FLOAT:
                return std::to_string(as_float());
            default:
                return "<?>";
        }
    }
    
    // For debugging - get type name
    const char* type_name() const {
        switch (type_) {
            case Type::NIL: return "nil";
            case Type::BOOLEAN: return "boolean";
            case Type::INTEGER: return "integer";
            case Type::FLOAT: return "float";
            default: return "unknown";
        }
    }
    
private:
    // Private constructor used by factory methods
    TaggedValue(Type type, uint32_t int_value) : type_(type), int_val_(int_value) {}
    
    Type type_;
    
    // Union for the actual value
    union {
        uint32_t int_val_;  // Used for boolean and integer
        float float_val_;   // Used for float
    };
};

// === STREAM OUTPUT FOR DEBUGGING ===

inline std::ostream& operator<<(std::ostream& os, const TaggedValue& tv) {
    os << tv.to_string();
    return os;
}

// === UTILITY FUNCTIONS ===

/**
 * Create TaggedValue from various C++ types (convenience functions)
 */
inline TaggedValue make_tagged_value(std::nullptr_t) {
    return TaggedValue::nil();
}

inline TaggedValue make_tagged_value(bool value) {
    return TaggedValue::boolean(value);
}

inline TaggedValue make_tagged_value(int32_t value) {
    return TaggedValue::integer(value);
}

inline TaggedValue make_tagged_value(float value) {
    return TaggedValue::float_value(value);
}

inline TaggedValue make_tagged_value(double value) {
    return TaggedValue::float_value(static_cast<float>(value));
}

/**
 * Check if two TaggedValues are the same object (for Smalltalk == semantics)
 * For immediates, this is the same as value equality
 */
inline bool same_object(const TaggedValue& a, const TaggedValue& b) {
    return a == b;
}

/**
 * Smalltalk-style truthiness: nil and false are false, everything else is true
 */
inline bool is_truthy(const TaggedValue& value) {
    if (value.is_nil()) return false;
    if (value.is_boolean()) return value.as_boolean();
    return true;  // Numbers are always truthy
}

} // namespace smalltalk