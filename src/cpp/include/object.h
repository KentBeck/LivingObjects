#pragma once

#include <cstdint>
#include <cstddef>

namespace smalltalk {

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
    
    // Constructor
    Object(ObjectType objectType, size_t objectSize, uint32_t objectHash = 0)
        : header(objectType, objectSize, objectHash) {}
};

} // namespace smalltalk
