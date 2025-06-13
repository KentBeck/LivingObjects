#pragma once

#include <cstdint>
#include <cstddef>

namespace smalltalk {

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
    PINNED          = 4  // Cannot be moved by GC
};

// Object header structure (64 bits)
struct ObjectHeader {
    uint64_t size : 24;      // Size in slots or bytes
    uint64_t flags : 5;      // Various flags
    uint64_t type : 3;       // Object type
    uint64_t hash : 32;      // Identity hash
    
    // Constructor
    ObjectHeader(ObjectType type, size_t size, uint32_t hash = 0)
        : size(size), flags(0), type(static_cast<uint64_t>(type)), hash(hash) {}
    
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
    
    // Constructor
    Object(ObjectType type, size_t size, uint32_t hash = 0)
        : header(type, size, hash) {}
};

} // namespace smalltalk