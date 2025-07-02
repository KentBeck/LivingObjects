#include "tagged_value.h"
#include "object.h"
#include "smalltalk_class.h"
#include "smalltalk_string.h"
#include "symbol.h"

namespace smalltalk {

bool TaggedValue::isObjectOfClass(Class* clazz) const {
    if (!isPointer()) {
        return false;
    }
    
    try {
        Object* obj = asObject();
        return obj->getClass() == clazz;
    } catch (...) {
        return false;
    }
}

Class* TaggedValue::getClass() const {
    if (isInteger()) {
        return ClassUtils::getIntegerClass();
    } else if (isBoolean()) {
        return ClassUtils::getBooleanClass();
    } else if (isNil()) {
        // nil is a special case - it's the sole instance of UndefinedObject
        // For now, we'll return Object class
        return ClassUtils::getObjectClass();
    } else if (isPointer()) {
        try {
            // Try as Symbol first
            asSymbol();  // Just check if it's a symbol
            return ClassUtils::getSymbolClass();
        } catch (...) {
            try {
                // Try as String
                Object* obj = asObject();
                if (StringUtils::isString(*this)) {
                    return ClassUtils::getStringClass();
                }
                // Try as generic Object
                return obj->getClass();
            } catch (...) {
                return nullptr;
            }
        }
    }
    return nullptr;
}

} // namespace smalltalk