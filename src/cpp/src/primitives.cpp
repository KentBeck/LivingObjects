#include "../include/primitives.h"
#include "../include/interpreter.h"
#include "../include/smalltalk_class.h"
#include "../include/object.h"

namespace smalltalk {

// PrimitiveRegistry implementation
PrimitiveRegistry& PrimitiveRegistry::getInstance() {
    static PrimitiveRegistry instance;
    return instance;
}

void PrimitiveRegistry::registerPrimitive(int primitiveNumber, PrimitiveFunction function) {
    primitives_[primitiveNumber] = function;
}

TaggedValue PrimitiveRegistry::callPrimitive(int primitiveNumber, TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter) {
    auto it = primitives_.find(primitiveNumber);
    if (it == primitives_.end()) {
        throw PrimitiveFailure("Primitive " + std::to_string(primitiveNumber) + " not found");
    }
    
    return it->second(receiver, args, interpreter);
}

bool PrimitiveRegistry::hasPrimitive(int primitiveNumber) const {
    return primitives_.find(primitiveNumber) != primitives_.end();
}

std::vector<int> PrimitiveRegistry::getAllPrimitiveNumbers() const {
    std::vector<int> numbers;
    for (const auto& pair : primitives_) {
        numbers.push_back(pair.first);
    }
    return numbers;
}

void PrimitiveRegistry::clear() {
    primitives_.clear();
}

void PrimitiveRegistry::initializeCorePrimitives() {
    // Forward declarations for primitive registration functions
    extern void registerObjectPrimitives();
    extern void registerArrayPrimitives();
    extern void registerIntegerPrimitives();
    extern void registerBlockPrimitives();
    
    // Register all core primitive groups
    registerObjectPrimitives();
    registerArrayPrimitives();
    // registerIntegerPrimitives(); // Will implement later
    // registerBlockPrimitives(); // Will implement later
}

// Primitives namespace implementation
namespace Primitives {

void initialize() {
    PrimitiveRegistry::getInstance().initializeCorePrimitives();
}

TaggedValue callPrimitive(int primitiveNumber, TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter) {
    return PrimitiveRegistry::getInstance().callPrimitive(primitiveNumber, receiver, args, interpreter);
}

void checkArgumentCount(const std::vector<TaggedValue>& args, size_t expected, const std::string& primitiveName) {
    if (args.size() != expected) {
        throw PrimitiveFailure("Primitive " + primitiveName + " expects " + 
                              std::to_string(expected) + " arguments, got " + 
                              std::to_string(args.size()));
    }
}

void checkReceiverType(TaggedValue receiver, ObjectType expectedType, const std::string& primitiveName) {
    if (!receiver.isPointer()) {
        throw PrimitiveFailure("Primitive " + primitiveName + " expects object receiver");
    }
    
    Object* obj = receiver.asObject();
    if (obj->header.getType() != expectedType) {
        throw PrimitiveFailure("Primitive " + primitiveName + " expects receiver of type " + 
                              std::to_string(static_cast<int>(expectedType)));
    }
}

Class* ensureReceiverIsClass(TaggedValue receiver, const std::string& primitiveName) {
    if (!receiver.isPointer()) {
        throw PrimitiveFailure("Primitive " + primitiveName + " expects class receiver");
    }
    
    Object* obj = receiver.asObject();
    if (obj->header.getType() != ObjectType::CLASS) {
        throw PrimitiveFailure("Primitive " + primitiveName + " expects class receiver");
    }
    
    return static_cast<Class*>(obj);
}

Object* ensureReceiverIsObject(TaggedValue receiver, const std::string& primitiveName) {
    if (!receiver.isPointer()) {
        throw PrimitiveFailure("Primitive " + primitiveName + " expects object receiver");
    }
    
    return receiver.asObject();
}

} // namespace Primitives

} // namespace smalltalk