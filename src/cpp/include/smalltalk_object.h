#pragma once

#include <cstdint>
#include <cstddef>
#include <string>

namespace smalltalk {

// Forward declarations
class Class;
class Symbol;

/**
 * Object is the root of the Smalltalk class hierarchy.
 * Every Smalltalk object inherits from Object.
 */
class Object {
public:
    Object(Class* objectClass);
    virtual ~Object() = default;
    
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
    
    // Object header for GC (size, type, etc.)
    struct Header {
        uint32_t size;        // Size in bytes
        uint32_t type;        // Object type identifier
        Class* class_;        // Class pointer
        uint32_t padding;     // Alignment padding
    };
    
    // Get the object header
    Header* getHeader() {
        return reinterpret_cast<Header*>(reinterpret_cast<char*>(this) - sizeof(Header));
    }
    
    const Header* getHeader() const {
        return reinterpret_cast<const Header*>(reinterpret_cast<const char*>(this) - sizeof(Header));
    }
    
protected:
    Class* class_;  // Every object knows its class
};

/**
 * SmallInteger represents immediate integer values.
 * These are not heap-allocated objects but are represented
 * via tagged pointers in TaggedValue.
 */
class SmallInteger : public Object {
public:
    SmallInteger(int32_t value, Class* integerClass);
    
    int32_t getValue() const { return value_; }
    void setValue(int32_t value) { value_ = value; }
    
    std::string toString() const override;
    
private:
    int32_t value_;
};

/**
 * Boolean represents true and false values.
 * Like SmallInteger, these are typically immediate values.
 */
class Boolean : public Object {
public:
    Boolean(bool value, Class* booleanClass);
    
    bool getValue() const { return value_; }
    void setValue(bool value) { value_ = value; }
    
    std::string toString() const override;
    
private:
    bool value_;
};

} // namespace smalltalk