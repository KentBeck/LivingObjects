#include "../include/object.h"
#include "../include/interpreter.h"
#include "../include/memory_manager.h"
#include "../include/primitives.h"
#include "../include/smalltalk_class.h"
#include "../include/smalltalk_exception.h"
#include <memory>

namespace smalltalk {

namespace ObjectPrimitives {

/**
 * Primitive 70: Object new
 * Creates a new instance of the receiver class
 */
TaggedValue primitive_new(TaggedValue receiver,
                          const std::vector<TaggedValue> &args,
                          Interpreter &interpreter) {
  // Check argument count
  Primitives::checkArgumentCount(args, 0, "new");

  // Ensure receiver is a class
  Class *clazz = Primitives::ensureReceiverIsClass(receiver, "new");

  // Check if class is indexable (should use basicNew: instead)
  if (clazz->isIndexable()) {
    throw PrimitiveFailure(
        "Cannot create indexable object without size - use basicNew:");
  }

  // Allocate new instance
  Object *instance = interpreter.getMemoryManager().allocateInstance(clazz);

  return TaggedValue::fromObject(instance);
}

/**
 * Primitive 71: Object basicNew
 * Creates a new instance without calling initialize
 */
TaggedValue primitive_basic_new(TaggedValue receiver,
                                const std::vector<TaggedValue> &args,
                                Interpreter &interpreter) {
  // Check argument count
  Primitives::checkArgumentCount(args, 0, "basicNew");

  // Ensure receiver is a class
  Class *clazz = Primitives::ensureReceiverIsClass(receiver, "basicNew");

  // Check if class is indexable (should use basicNew: instead)
  if (clazz->isIndexable()) {
    throw PrimitiveFailure(
        "Cannot create indexable object without size - use basicNew:");
  }

  // Allocate new instance
  Object *instance = interpreter.getMemoryManager().allocateInstance(clazz);

  return TaggedValue::fromObject(instance);
}

/**
 * Primitive 72: Object basicNew: size
 * Creates a new indexable instance with the given size
 */
TaggedValue primitive_basic_new_size(TaggedValue receiver,
                                     const std::vector<TaggedValue> &args,
                                     Interpreter &interpreter) {
  // Check argument count
  Primitives::checkArgumentCount(args, 1, "basicNew:");

  // Ensure receiver is a class
  Class *clazz = Primitives::ensureReceiverIsClass(receiver, "basicNew:");

  // Get size argument
  TaggedValue sizeValue = args[0];
  if (!sizeValue.isSmallInteger()) {
    throw PrimitiveFailure("Size argument must be a SmallInteger");
  }

  int32_t size = sizeValue.getSmallInteger();
  if (size < 0) {
    // Throw proper ArgumentError for negative sizes
    auto exception = std::make_unique<ArgumentError>(
        "Size must be non-negative: " + std::to_string(size));
    ExceptionHandler::throwException(std::move(exception));
  }

  // Allocate new indexable instance
  Object *instance = nullptr;
  if (clazz->isByteIndexable()) {
    instance = interpreter.getMemoryManager().allocateByteIndexableInstance(
        clazz, static_cast<size_t>(size));
  } else {
    instance = interpreter.getMemoryManager().allocateIndexableInstance(
        clazz, static_cast<size_t>(size));
  }

  return TaggedValue::fromObject(instance);
}

/**
 * Primitive 75: Object identityHash
 * Returns the identity hash of the receiver
 */
TaggedValue primitive_identity_hash(TaggedValue receiver,
                                    const std::vector<TaggedValue> &args,
                                    Interpreter &interpreter) {
  (void)interpreter; // Suppress unused parameter warning
  // Check argument count
  Primitives::checkArgumentCount(args, 0, "identityHash");

  // Handle immediate values
  if (receiver.isSmallInteger()) {
    // For small integers, use the value itself as hash
    return TaggedValue::fromSmallInteger(receiver.getSmallInteger());
  }
  if (receiver.isBoolean()) {
    // For booleans, use 0 for false, 1 for true
    return TaggedValue::fromSmallInteger(receiver.getBoolean() ? 1 : 0);
  }
  if (receiver.isNil()) {
    // For nil, use a fixed hash
    return TaggedValue::fromSmallInteger(42);
  }

  // For heap objects, get hash from object header
  Object *obj = Primitives::ensureReceiverIsObject(receiver, "identityHash");
  uint32_t hash = obj->header.getHash();

  // Ensure hash fits in SmallInteger range
  if (hash > INT32_MAX) {
    hash = hash % INT32_MAX;
  }

  return TaggedValue::fromSmallInteger(static_cast<int32_t>(hash));
}

/**
 * Primitive 111: Object class
 * Returns the class of the receiver
 */
TaggedValue primitive_class(TaggedValue receiver,
                            const std::vector<TaggedValue> &args,
                            Interpreter &interpreter) {
  (void)interpreter; // Suppress unused parameter warning
  // Check argument count
  Primitives::checkArgumentCount(args, 0, "class");

  Class *receiverClass = nullptr;

  // Handle immediate values
  if (receiver.isSmallInteger()) {
    receiverClass = ClassUtils::getIntegerClass();
  } else if (receiver.isBoolean()) {
    receiverClass = receiver.getBoolean() ? ClassUtils::getTrueClass()
                                          : ClassUtils::getFalseClass();
  } else if (receiver.isNil()) {
    // UndefinedObject class for nil
    receiverClass = ClassRegistry::getInstance().getClass("UndefinedObject");
    if (receiverClass == nullptr) {
      throw PrimitiveFailure("UndefinedObject class not found");
    }
  } else if (receiver.isPointer()) {
    Object *obj = receiver.asObject();
    receiverClass = obj->getClass();
    if (receiverClass == nullptr) {
      throw PrimitiveFailure("Object has no class");
    }
  } else {
    throw PrimitiveFailure("Unknown receiver type");
  }

  return TaggedValue::fromObject(receiverClass);
}

} // namespace ObjectPrimitives

// Register object primitives
void registerObjectPrimitives() {
  PrimitiveRegistry &registry = PrimitiveRegistry::getInstance();

  registry.registerPrimitive(PrimitiveNumbers::NEW,
                             ObjectPrimitives::primitive_new);
  registry.registerPrimitive(PrimitiveNumbers::BASIC_NEW,
                             ObjectPrimitives::primitive_basic_new);
  registry.registerPrimitive(PrimitiveNumbers::BASIC_NEW_SIZE,
                             ObjectPrimitives::primitive_basic_new_size);
  registry.registerPrimitive(PrimitiveNumbers::IDENTITY_HASH,
                             ObjectPrimitives::primitive_identity_hash);
  registry.registerPrimitive(PrimitiveNumbers::CLASS,
                             ObjectPrimitives::primitive_class);
}

} // namespace smalltalk
