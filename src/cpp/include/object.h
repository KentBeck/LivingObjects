#pragma once

#include <cstdint>
#include <cstddef>
#include <string>

namespace smalltalk {

// Forward declarations
class Class;

// Object header bitfield sizes
const uint64_t OBJECT_HEADER_SIZE_BITS = 24;
const uint64_t OBJECT_HEADER_FLAGS_BITS = 5;
const uint64_t OBJECT_HEADER_TYPE_BITS = 3;
const uint64_t OBJECT_HEADER_HASH_BITS = 32;

// Object type discriminators
enum class ObjectType : uint8_t {
    IMMEDIATE       = 0, // SmallInteger, Character, Boolean
    OBJECT          = 1, // General object
    ARRAY           = 2, // Indexable objects without instance variables
    BYTE_ARRAY      = 3, // Byte arrays
    SYMBOL          = 4, // Interned strings
    CONTEXT         = 5, // Method/block context
    CLASS           = 6, // Class object
    METHOD          = 7  // Compiled method
};

// Object flags
enum class ObjectFlag : uint8_t {
    MARKED          = 0, // Marked during GC
    REMEMBERED      = 1, // In the remembered set
    IMMUTABLE       = 2, // Cannot be modified
    FORWARDED       = 3, // Object has been forwarded
    PINNED          = 4,  // Cannot be moved by GC
    CONTAINS_POINTERS = 5 // Object contains pointers to other objects
};

// Object header structure (64 bits)
struct ObjectHeader {
    uint64_t size : OBJECT_HEADER_SIZE_BITS;      // Size in slots or bytes
    uint64_t flags : OBJECT_HEADER_FLAGS_BITS;      // Various flags
    uint64_t type : OBJECT_HEADER_TYPE_BITS;       // Object type
    uint64_t hash : OBJECT_HEADER_HASH_BITS;      // Identity hash
    
    // Constructor
    ObjectHeader(ObjectType objectType, size_t objectSize, uint32_t objectHash = 0)
        : size(objectSize), flags(0), type(static_cast<uint64_t>(objectType)), hash(objectHash) {}
    
    // Flag operations
    void setFlag(ObjectFlag flag) {
        flags |= (1 << static_cast<uint8_t>(flag));
    }
    
    bool hasFlag(ObjectFlag flag) const {
        return (flags & (1 << static_cast<uint8_t>(flag))) != 0;
    }
    
    void clearFlag(ObjectFlag flag) {
        flags &= ~(1 << static_cast<uint8_t>(flag));
    }
};

// Base object structure
struct Object {
    ObjectHeader header;
    Class* class_ = nullptr;  // Every object knows its class
    
    // Constructor
    Object(ObjectType objectType, size_t objectSize, uint32_t objectHash = 0)
        : header(objectType, objectSize, objectHash) {}
    
    // Constructor with class
    Object(ObjectType objectType, size_t objectSize, Class* objectClass, uint32_t objectHash = 0)
        : header(objectType, objectSize, objectHash), class_(objectClass) {}
    
    // Get the class of this object
    Class* getClass() const { return class_; }
    
    // Set the class (used during object creation)
    void setClass(Class* objectClass) { class_ = objectClass; }
    
    // Object identity
    bool operator==(const Object& other) const {
        return this == &other;
    }
    
    bool operator!=(const Object& other) const {
        return this != &other;
    }
    
    // Hash code for object identity
    size_t hash() const {
        return reinterpret_cast<size_t>(this);
    }
    
    // String representation for debugging
    virtual std::string toString() const;
    
    // Virtual destructor for proper cleanup
    virtual ~Object() = default;
};

/**
 * SmallInteger represents immediate integer values.
 * These are not heap-allocated objects but are represented
 * via tagged pointers in TaggedValue.
 */
class SmallInteger : public Object {
public:
    SmallInteger(int32_t value, Class* integerClass)
        : Object(ObjectType::IMMEDIATE, sizeof(SmallInteger), integerClass), value_(value) {}
    
    int32_t getValue() const { return value_; }
    void setValue(int32_t value) { value_ = value; }
    
    std::string toString() const override {
        return std::to_string(value_);
    }
    
private:
    int32_t value_;
};

/**
 * Boolean represents true and false values.
 * Like SmallInteger, these are typically immediate values.
 */
class Boolean : public Object {
public:
    Boolean(bool value, Class* booleanClass)
        : Object(ObjectType::IMMEDIATE, sizeof(Boolean), booleanClass), value_(value) {}
    
    bool getValue() const { return value_; }
    void setValue(bool value) { value_ = value; }
    
    std::string toString() const override {
        return value_ ? "true" : "false";
    }
    
private:
    bool value_;
};

} // namespace smalltalk
