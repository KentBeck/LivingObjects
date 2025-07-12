#include "../include/primitives.h"
#include "../include/interpreter.h"
#include "../include/smalltalk_exception.h"
#include "../include/smalltalk_class.h"
#include "../include/memory_manager.h"
#include <memory>

namespace smalltalk {

namespace ExceptionPrimitives {

/**
 * Primitive 1000: Exception marker
 * This primitive always fails but serves as a marker that this method
 * handles exceptions. During stack unwinding, we check for this primitive.
 */
TaggedValue primitive_exception_mark(TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter) {
    (void)receiver;
    (void)args;
    (void)interpreter;
    
    // This primitive always fails to trigger fallback to Smalltalk code
    throw PrimitiveFailure("Exception handler marker - always fails");
}

/**
 * Primitive 1001: Signal exception
 * Throws the receiver as an exception
 */
TaggedValue primitive_signal(TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter) {
    (void)args;
    (void)interpreter;
    
    // For now, we'll throw a C++ exception with the receiver's class name
    // In a full implementation, this would create a proper Smalltalk exception object
    
    if (!receiver.isPointer()) {
        throw PrimitiveFailure("Can only signal object exceptions");
    }
    
    Object* obj = receiver.asObject();
    if (!obj) {
        throw PrimitiveFailure("Cannot signal nil");
    }
    
    // Get the exception class name
    Class* exceptionClass = obj->getClass();
    std::string className = exceptionClass ? exceptionClass->getName() : "Exception";
    
    // Create and throw the appropriate exception type
    if (className == "ZeroDivisionError") {
        ExceptionHandler::throwException(std::make_unique<ZeroDivisionError>());
    } else if (className == "NameError") {
        ExceptionHandler::throwException(std::make_unique<NameError>("unknown"));
    } else if (className == "IndexError") {
        ExceptionHandler::throwException(std::make_unique<IndexError>());
    } else if (className == "ArgumentError") {
        ExceptionHandler::throwException(std::make_unique<ArgumentError>());
    } else if (className == "MessageNotUnderstood") {
        ExceptionHandler::throwException(std::make_unique<MessageNotUnderstood>("Object", "unknown"));
    } else {
        ExceptionHandler::throwException(std::make_unique<RuntimeError>(className));
    }
    
    // Never reached
    return TaggedValue::nil();
}

} // namespace ExceptionPrimitives

// Register exception primitives
void registerExceptionPrimitives() {
    PrimitiveRegistry& registry = PrimitiveRegistry::getInstance();
    
    registry.registerPrimitive(PrimitiveNumbers::EXCEPTION_MARK, ExceptionPrimitives::primitive_exception_mark);
    registry.registerPrimitive(PrimitiveNumbers::EXCEPTION_SIGNAL, ExceptionPrimitives::primitive_signal);
}

} // namespace smalltalk