#pragma once

#include <cstdint>
#include <cstddef>
#include <string>

namespace smalltalk
{

    // Forward declarations
    class Class;
    enum class ContextType : uint8_t;

    // Object header bitfield sizes
    const uint64_t OBJECT_HEADER_SIZE_BITS = 24;
    const uint64_t OBJECT_HEADER_FLAGS_BITS = 8; // Now includes ObjectType (3 bits) + flags (3 bits) + 2 spare
    const uint64_t OBJECT_HEADER_HASH_BITS = 32;

    // ObjectType is packed into the flags field (bits 3-5)
    const uint64_t OBJECT_TYPE_SHIFT = 3;
    const uint64_t OBJECT_TYPE_MASK = 0x7;  // 3 bits mask
    const uint64_t OBJECT_FLAGS_MASK = 0x7; // 3 bits mask for actual flags

    // ContextType is packed into the top 3 bits of hash field for CONTEXT objects
    const uint64_t CONTEXT_TYPE_SHIFT = 29;
    const uint64_t CONTEXT_TYPE_MASK = 0x7;        // 3 bits mask
    const uint64_t CONTEXT_HASH_MASK = 0x1FFFFFFF; // 29 bits mask for actual hash

    // Object type discriminators
    enum class ObjectType : uint8_t
    {
        IMMEDIATE = 0,  // SmallInteger, Character, Boolean
        OBJECT = 1,     // General object
        ARRAY = 2,      // Indexable objects without instance variables
        BYTE_ARRAY = 3, // Byte arrays
        SYMBOL = 4,     // Interned strings
        CONTEXT = 5,    // Method/block context
        CLASS = 6,      // Class object
        METHOD = 7      // Compiled method
    };

    // Object flags
    enum class ObjectFlag : uint8_t
    {
        MARKED = 0,              // Marked during GC
        REMEMBERED = 1,          // In the remembered set
        IMMUTABLE = 2,           // Cannot be modified
        FORWARDED = 3,           // Object has been forwarded
        PINNED = 4,              // Cannot be moved by GC
        CONTAINS_POINTERS = 5,   // Object contains pointers to other objects
        TAGGED_VALUE_WRAPPER = 6 // Object wraps a TaggedValue (for immediate values on stack)
    };

    // Object header structure (64 bits)
    struct ObjectHeader
    {
        uint64_t size : OBJECT_HEADER_SIZE_BITS;   // Size in slots or bytes
        uint64_t flags : OBJECT_HEADER_FLAGS_BITS; // Flags (3 bits) + ObjectType (3 bits) + 2 spare
        uint64_t hash : OBJECT_HEADER_HASH_BITS;   // Identity hash

        // Constructor
        ObjectHeader(ObjectType objectType, size_t objectSize, uint32_t objectHash = 0)
            : size(objectSize), flags(static_cast<uint64_t>(objectType) << OBJECT_TYPE_SHIFT), hash(objectHash) {}

        // ObjectType operations
        ObjectType getType() const
        {
            return static_cast<ObjectType>((flags >> OBJECT_TYPE_SHIFT) & OBJECT_TYPE_MASK);
        }

        void setType(ObjectType objectType)
        {
            flags = (flags & ~(OBJECT_TYPE_MASK << OBJECT_TYPE_SHIFT)) |
                    (static_cast<uint64_t>(objectType) << OBJECT_TYPE_SHIFT);
        }

        // ContextType operations (stored in top 3 bits of hash field for CONTEXT objects)
        uint8_t getContextType() const
        {
            return static_cast<uint8_t>((hash >> CONTEXT_TYPE_SHIFT) & CONTEXT_TYPE_MASK);
        }

        void setContextType(uint8_t contextType)
        {
            hash = (hash & CONTEXT_HASH_MASK) | (static_cast<uint64_t>(contextType) << CONTEXT_TYPE_SHIFT);
        }

        // Get the actual hash value (29 bits for CONTEXT objects, 32 bits for others)
        uint32_t getHash() const
        {
            if (getType() == ObjectType::CONTEXT)
            {
                return static_cast<uint32_t>(hash & CONTEXT_HASH_MASK);
            }
            return static_cast<uint32_t>(hash);
        }

        // Set the actual hash value (preserves ContextType for CONTEXT objects)
        void setHash(uint32_t hashValue)
        {
            if (getType() == ObjectType::CONTEXT)
            {
                hash = (hash & ~CONTEXT_HASH_MASK) | (hashValue & CONTEXT_HASH_MASK);
            }
            else
            {
                hash = hashValue;
            }
        }

        // Flag operations
        void setFlag(ObjectFlag flag)
        {
            flags |= (1 << static_cast<uint8_t>(flag));
        }

        bool hasFlag(ObjectFlag flag) const
        {
            return (flags & (1 << static_cast<uint8_t>(flag))) != 0;
        }

        void clearFlag(ObjectFlag flag)
        {
            flags &= ~(1 << static_cast<uint8_t>(flag));
        }
    };

    // Base object structure
    struct Object
    {
        ObjectHeader header;
        Class *class_ = nullptr; // Every object knows its class

        // Constructor
        Object(ObjectType objectType, size_t objectSize, uint32_t objectHash = 0)
            : header(objectType, objectSize, objectHash) {}

        // Constructor with class
        Object(ObjectType objectType, size_t objectSize, Class *objectClass, uint32_t objectHash = 0)
            : header(objectType, objectSize, objectHash), class_(objectClass) {}

        // Get the class of this object
        Class *getClass() const { return class_; }

        // Set the class (used during object creation)
        void setClass(Class *objectClass) { class_ = objectClass; }

        // Object identity
        bool operator==(const Object &other) const
        {
            return this == &other;
        }

        bool operator!=(const Object &other) const
        {
            return this != &other;
        }

        // Hash code for object identity
        size_t hash() const
        {
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
    class SmallInteger : public Object
    {
    public:
        SmallInteger(int32_t value, Class *integerClass)
            : Object(ObjectType::IMMEDIATE, sizeof(SmallInteger), integerClass), value_(value) {}

        int32_t getValue() const { return value_; }
        void setValue(int32_t value) { value_ = value; }

        std::string toString() const override
        {
            return std::to_string(value_);
        }

    private:
        int32_t value_;
    };

    /**
     * Boolean represents true and false values.
     * Like SmallInteger, these are typically immediate values.
     */
    class Boolean : public Object
    {
    public:
        Boolean(bool value, Class *booleanClass)
            : Object(ObjectType::IMMEDIATE, sizeof(Boolean), booleanClass), value_(value) {}

        bool getValue() const { return value_; }
        void setValue(bool value) { value_ = value; }

        std::string toString() const override
        {
            return value_ ? "true" : "false";
        }

    private:
        bool value_;
    };

} // namespace smalltalk
