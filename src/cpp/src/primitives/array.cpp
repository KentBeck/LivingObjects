#include "../include/primitives.h"
#include "../include/interpreter.h"
#include "../include/smalltalk_class.h"
#include "../include/memory_manager.h"
#include "../include/smalltalk_exception.h"
#include "../include/object.h"
#include <memory>

namespace smalltalk {

namespace ArrayPrimitives {

/**
 * Primitive 72: Array new: size
 * Creates a new array with the given size, all elements initialized to nil
 */
TaggedValue primitive_new_size(TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter) {
    // Check argument count
    Primitives::checkArgumentCount(args, 1, "new:");
    
    // Ensure receiver is a class
    Class* clazz = Primitives::ensureReceiverIsClass(receiver, "new:");
    
    // Get size argument
    TaggedValue sizeValue = args[0];
    if (!sizeValue.isSmallInteger()) {
        throw PrimitiveFailure("Size argument must be a SmallInteger");
    }
    
    int32_t size = sizeValue.getSmallInteger();
    if (size < 0) {
        // Throw proper ArgumentError for negative sizes
        auto exception = std::make_unique<ArgumentError>("Size must be non-negative: " + std::to_string(size));
        ExceptionHandler::throwException(std::move(exception));
    }
    
    // Allocate new indexable instance (Array)
    Object* instance = interpreter.getMemoryManager().allocateIndexableInstance(clazz, static_cast<size_t>(size));
    
    // Array elements are automatically initialized to nullptr by allocateIndexableInstance
    // which will be interpreted as nil when accessed
    
    return TaggedValue::fromObject(instance);
}

/**
 * Primitive 60: Array at: index
 * Returns the element at the given 1-based index
 */
TaggedValue primitive_at(TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter) {
    (void)interpreter; // Suppress unused parameter warning
    
    // Check argument count
    Primitives::checkArgumentCount(args, 1, "at:");
    
    // Ensure receiver is an object
    Object* obj = Primitives::ensureReceiverIsObject(receiver, "at:");
    
    // Check that it's an array type
    if (obj->header.getType() != ObjectType::ARRAY) {
        throw PrimitiveFailure("at: can only be sent to arrays");
    }
    
    // Get index argument
    TaggedValue indexValue = args[0];
    if (!indexValue.isSmallInteger()) {
        throw PrimitiveFailure("Index must be a SmallInteger");
    }
    
    int32_t index = indexValue.getSmallInteger();
    if (index < 1) {
        throw PrimitiveFailure("Array index must be >= 1");
    }
    
    // Convert to 0-based index
    size_t arrayIndex = static_cast<size_t>(index - 1);
    
    // Get array size from header
    size_t arraySize = obj->header.size;
    
    if (arrayIndex >= arraySize) {
        throw PrimitiveFailure("Array index out of bounds");
    }
    
    // Get the element pointer
    Object** elements = reinterpret_cast<Object**>(reinterpret_cast<char*>(obj) + sizeof(Object));
    Object* element = elements[arrayIndex];
    
    // If element is nullptr (uninitialized), return nil
    if (element == nullptr) {
        return TaggedValue::nil();
    }
    
    return TaggedValue::fromObject(element);
}

/**
 * Primitive 61: Array at: index put: value
 * Sets the element at the given 1-based index to the value
 */
TaggedValue primitive_at_put(TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter) {
    (void)interpreter; // Suppress unused parameter warning
    
    // Check argument count
    Primitives::checkArgumentCount(args, 2, "at:put:");
    
    // Ensure receiver is an object
    Object* obj = Primitives::ensureReceiverIsObject(receiver, "at:put:");
    
    // Check that it's an array type
    if (obj->header.getType() != ObjectType::ARRAY) {
        throw PrimitiveFailure("at:put: can only be sent to arrays");
    }
    
    // Get index argument
    TaggedValue indexValue = args[0];
    if (!indexValue.isSmallInteger()) {
        throw PrimitiveFailure("Index must be a SmallInteger");
    }
    
    int32_t index = indexValue.getSmallInteger();
    if (index < 1) {
        throw PrimitiveFailure("Array index must be >= 1");
    }
    
    // Convert to 0-based index
    size_t arrayIndex = static_cast<size_t>(index - 1);
    
    // Get array size from header
    size_t arraySize = obj->header.size;
    
    if (arrayIndex >= arraySize) {
        throw PrimitiveFailure("Array index out of bounds");
    }
    
    // Get value argument
    TaggedValue value = args[1];
    Object* valueObject = nullptr;
    
    // Convert TaggedValue to Object*
    if (value.isNil()) {
        valueObject = nullptr; // nil is stored as nullptr
    } else if (value.isPointer()) {
        valueObject = value.asObject();
    } else if (value.isSmallInteger()) {
        // Box integer immediate value
        valueObject = interpreter.getMemoryManager().allocateInteger(value.getSmallInteger());
    } else if (value.isBoolean()) {
        // Box boolean immediate value
        valueObject = interpreter.getMemoryManager().allocateBoolean(value.getBoolean());
    } else {
        // Handle other immediate value types (e.g., float)
        throw PrimitiveFailure("Unsupported immediate value type for array storage");
    }
    
    // Set the element
    Object** elements = reinterpret_cast<Object**>(reinterpret_cast<char*>(obj) + sizeof(Object));
    elements[arrayIndex] = valueObject;
    
    // Return the value that was stored
    return value;
}

/**
 * Primitive 62: Array size
 * Returns the size of the array
 */
TaggedValue primitive_size(TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter) {
    (void)interpreter; // Suppress unused parameter warning
    
    // Check argument count
    Primitives::checkArgumentCount(args, 0, "size");
    
    // Ensure receiver is an object
    Object* obj = Primitives::ensureReceiverIsObject(receiver, "size");
    
    // Check that it's an array type
    if (obj->header.getType() != ObjectType::ARRAY) {
        throw PrimitiveFailure("size can only be sent to arrays");
    }
    
    // Get array size from header
    size_t arraySize = obj->header.size;
    
    // Return as SmallInteger
    if (arraySize > INT32_MAX) {
        throw PrimitiveFailure("Array size too large for SmallInteger");
    }
    
    return TaggedValue::fromSmallInteger(static_cast<int32_t>(arraySize));
}

} // namespace ArrayPrimitives

// Register array primitives
void registerArrayPrimitives() {
    PrimitiveRegistry& registry = PrimitiveRegistry::getInstance();
    
    registry.registerPrimitive(PrimitiveNumbers::BASIC_NEW_SIZE, ArrayPrimitives::primitive_new_size);
    registry.registerPrimitive(PrimitiveNumbers::ARRAY_AT, ArrayPrimitives::primitive_at);
    registry.registerPrimitive(PrimitiveNumbers::ARRAY_AT_PUT, ArrayPrimitives::primitive_at_put);
    registry.registerPrimitive(PrimitiveNumbers::ARRAY_SIZE, ArrayPrimitives::primitive_size);
}

} // namespace smalltalk