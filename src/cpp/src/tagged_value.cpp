#include "tagged_value.h"
#include "memory_manager.h"
#include "object.h"
#include "smalltalk_class.h"
#include "smalltalk_string.h"
#include "symbol.h"

namespace smalltalk {

bool TaggedValue::isObjectOfClass(Class *clazz) const {
  if (!isPointer()) {
    return false;
  }

  try {
    Object *obj = asObject();
    return obj->getClass() == clazz;
  } catch (...) {
    return false;
  }
}

Class *TaggedValue::getClass() const {
  if (isInteger()) {
    return ClassUtils::getIntegerClass();
  } else if (isBoolean()) {
    return getBoolean() ? ClassUtils::getTrueClass()
                        : ClassUtils::getFalseClass();
  } else if (isNil()) {
    // nil's class is UndefinedObject
    return ClassUtils::getUndefinedObjectClass();
  } else if (isPointer()) {
    try {
      Object *obj = asObject();
      return obj ? obj->getClass() : nullptr;
    } catch (...) {
      return nullptr;
    }
  }
  return nullptr;
}

TaggedValue TaggedValue::fromObject(Object *object) {
  if (object == nullptr) {
    return TaggedValue::nil();
  }

  // Check if this is a tagged value wrapper
  if (object->header.hasFlag(ObjectFlag::TAGGED_VALUE_WRAPPER)) {
    // Check if it's a boxed integer
    if (object->getClass() == ClassUtils::getIntegerClass()) {
      TaggedValue *valueSlot = reinterpret_cast<TaggedValue *>(
          reinterpret_cast<char *>(object) + sizeof(Object));
      return *valueSlot;
    }
    // Check if it's a boxed boolean
    if (object->getClass() == ClassUtils::getBooleanClass() ||
        object->getClass() == ClassUtils::getTrueClass() ||
        object->getClass() == ClassUtils::getFalseClass()) {
      TaggedValue *valueSlot = reinterpret_cast<TaggedValue *>(
          reinterpret_cast<char *>(object) + sizeof(Object));
      return *valueSlot;
    }
    // If it's a generic tagged value wrapper, extract the original TaggedValue
    TaggedValue *storedValue = reinterpret_cast<TaggedValue *>(
        reinterpret_cast<char *>(object) + sizeof(Object));
    return *storedValue;
  }

  // Regular object pointer
  return TaggedValue(object);
}

Object *TaggedValue::toObject(MemoryManager &memoryManager) const {
  // Convert TaggedValue to Object* for legacy compatibility
  if (isPointer()) {
    return asObject();
  } else if (isInteger()) {
    return memoryManager.allocateInteger(asInteger());
  } else if (isBoolean()) {
    return memoryManager.allocateBoolean(asBoolean());
  } else if (isNil()) {
    // Nil is a special object, not an immediate value that needs boxing
    return ClassRegistry::getInstance().getClass("UndefinedObject");
  } else {
    // Should not happen for now, but handle other immediate types if they arise
    throw std::runtime_error("Unhandled immediate value type in toObject");
  }
}

} // namespace smalltalk
