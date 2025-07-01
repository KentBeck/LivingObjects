#include "smalltalk_object.h"
#include "smalltalk_class.h"

#include <sstream>

namespace smalltalk {

Object::Object(Class* objectClass) : class_(objectClass) {
    // Every object must have a class
    if (objectClass == nullptr) {
        // This can happen during bootstrap when Object class doesn't exist yet
        // We'll set it later during initialization
    }
}

std::string Object::toString() const {
    std::ostringstream oss;
    if (class_ != nullptr) {
        oss << "a " << class_->getName();
    } else {
        oss << "an Object";
    }
    oss << "@" << std::hex << reinterpret_cast<uintptr_t>(this);
    return oss.str();
}

// SmallInteger implementation
SmallInteger::SmallInteger(int32_t value, Class* integerClass) 
    : Object(integerClass), value_(value) {
}

std::string SmallInteger::toString() const {
    return std::to_string(value_);
}

// Boolean implementation
Boolean::Boolean(bool value, Class* booleanClass) 
    : Object(booleanClass), value_(value) {
}

std::string Boolean::toString() const {
    return value_ ? "true" : "false";
}

} // namespace smalltalk