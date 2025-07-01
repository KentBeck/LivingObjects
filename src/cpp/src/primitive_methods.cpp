#include "primitive_methods.h"
#include "smalltalk_object.h"
#include "smalltalk_class.h"
#include <stdexcept>
#include <sstream>

namespace smalltalk {

// PrimitiveMethod implementation
PrimitiveMethod::PrimitiveMethod(int primitiveNumber, PrimitiveFunction function)
    : primitiveNumber_(primitiveNumber), function_(function) {
}

TaggedValue PrimitiveMethod::execute(TaggedValue receiver, const std::vector<TaggedValue>& args) const {
    if (function_) {
        return function_(receiver, args);
    }
    throw std::runtime_error("Primitive method has no implementation");
}

std::string PrimitiveMethod::toString() const {
    std::ostringstream oss;
    oss << "PrimitiveMethod { primitive: " << primitiveNumber_ << " }";
    return oss.str();
}

// PrimitiveRegistry implementation
PrimitiveRegistry& PrimitiveRegistry::getInstance() {
    static PrimitiveRegistry instance;
    return instance;
}

void PrimitiveRegistry::registerPrimitive(int primitiveNumber, PrimitiveMethod::PrimitiveFunction function) {
    primitives_[primitiveNumber] = function;
}

std::shared_ptr<PrimitiveMethod> PrimitiveRegistry::createPrimitiveMethod(int primitiveNumber) const {
    auto it = primitives_.find(primitiveNumber);
    if (it != primitives_.end()) {
        return std::make_shared<PrimitiveMethod>(primitiveNumber, it->second);
    }
    return nullptr;
}

bool PrimitiveRegistry::hasPrimitive(int primitiveNumber) const {
    return primitives_.find(primitiveNumber) != primitives_.end();
}

void PrimitiveRegistry::initializeCorePrimitives() {
    // Register Integer arithmetic primitives
    registerPrimitive(PrimitiveNumbers::INTEGER_ADD, IntegerPrimitives::add);
    registerPrimitive(PrimitiveNumbers::INTEGER_SUBTRACT, IntegerPrimitives::subtract);
    registerPrimitive(PrimitiveNumbers::INTEGER_MULTIPLY, IntegerPrimitives::multiply);
    registerPrimitive(PrimitiveNumbers::INTEGER_DIVIDE, IntegerPrimitives::divide);
    
    // Register Integer comparison primitives
    registerPrimitive(PrimitiveNumbers::INTEGER_LESS_THAN, IntegerPrimitives::lessThan);
    registerPrimitive(PrimitiveNumbers::INTEGER_GREATER_THAN, IntegerPrimitives::greaterThan);
    registerPrimitive(PrimitiveNumbers::INTEGER_EQUAL, IntegerPrimitives::equal);
    registerPrimitive(PrimitiveNumbers::INTEGER_NOT_EQUAL, IntegerPrimitives::notEqual);
    registerPrimitive(PrimitiveNumbers::INTEGER_LESS_THAN_OR_EQUAL, IntegerPrimitives::lessThanOrEqual);
    registerPrimitive(PrimitiveNumbers::INTEGER_GREATER_THAN_OR_EQUAL, IntegerPrimitives::greaterThanOrEqual);
}

// Integer primitive implementations
namespace IntegerPrimitives {
    
    void checkArgumentCount(const std::vector<TaggedValue>& args, size_t expected) {
        if (args.size() != expected) {
            throw std::runtime_error("Wrong number of arguments: expected " + 
                                   std::to_string(expected) + ", got " + std::to_string(args.size()));
        }
    }
    
    void checkIntegerReceiver(TaggedValue receiver) {
        if (!receiver.isInteger()) {
            throw std::runtime_error("Receiver must be an integer");
        }
    }
    
    void checkIntegerArgument(TaggedValue arg, size_t index) {
        if (!arg.isInteger()) {
            throw std::runtime_error("Argument " + std::to_string(index) + " must be an integer");
        }
    }
    
    TaggedValue add(TaggedValue receiver, const std::vector<TaggedValue>& args) {
        checkArgumentCount(args, 1);
        checkIntegerReceiver(receiver);
        checkIntegerArgument(args[0], 0);
        
        int32_t result = receiver.asInteger() + args[0].asInteger();
        return TaggedValue(result);
    }
    
    TaggedValue subtract(TaggedValue receiver, const std::vector<TaggedValue>& args) {
        checkArgumentCount(args, 1);
        checkIntegerReceiver(receiver);
        checkIntegerArgument(args[0], 0);
        
        int32_t result = receiver.asInteger() - args[0].asInteger();
        return TaggedValue(result);
    }
    
    TaggedValue multiply(TaggedValue receiver, const std::vector<TaggedValue>& args) {
        checkArgumentCount(args, 1);
        checkIntegerReceiver(receiver);
        checkIntegerArgument(args[0], 0);
        
        int32_t result = receiver.asInteger() * args[0].asInteger();
        return TaggedValue(result);
    }
    
    TaggedValue divide(TaggedValue receiver, const std::vector<TaggedValue>& args) {
        checkArgumentCount(args, 1);
        checkIntegerReceiver(receiver);
        checkIntegerArgument(args[0], 0);
        
        int32_t divisor = args[0].asInteger();
        if (divisor == 0) {
            throw std::runtime_error("Division by zero");
        }
        
        int32_t result = receiver.asInteger() / divisor;
        return TaggedValue(result);
    }
    
    TaggedValue lessThan(TaggedValue receiver, const std::vector<TaggedValue>& args) {
        checkArgumentCount(args, 1);
        checkIntegerReceiver(receiver);
        checkIntegerArgument(args[0], 0);
        
        bool result = receiver.asInteger() < args[0].asInteger();
        return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
    }
    
    TaggedValue greaterThan(TaggedValue receiver, const std::vector<TaggedValue>& args) {
        checkArgumentCount(args, 1);
        checkIntegerReceiver(receiver);
        checkIntegerArgument(args[0], 0);
        
        bool result = receiver.asInteger() > args[0].asInteger();
        return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
    }
    
    TaggedValue equal(TaggedValue receiver, const std::vector<TaggedValue>& args) {
        checkArgumentCount(args, 1);
        checkIntegerReceiver(receiver);
        checkIntegerArgument(args[0], 0);
        
        bool result = receiver.asInteger() == args[0].asInteger();
        return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
    }
    
    TaggedValue notEqual(TaggedValue receiver, const std::vector<TaggedValue>& args) {
        checkArgumentCount(args, 1);
        checkIntegerReceiver(receiver);
        checkIntegerArgument(args[0], 0);
        
        bool result = receiver.asInteger() != args[0].asInteger();
        return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
    }
    
    TaggedValue lessThanOrEqual(TaggedValue receiver, const std::vector<TaggedValue>& args) {
        checkArgumentCount(args, 1);
        checkIntegerReceiver(receiver);
        checkIntegerArgument(args[0], 0);
        
        bool result = receiver.asInteger() <= args[0].asInteger();
        return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
    }
    
    TaggedValue greaterThanOrEqual(TaggedValue receiver, const std::vector<TaggedValue>& args) {
        checkArgumentCount(args, 1);
        checkIntegerReceiver(receiver);
        checkIntegerArgument(args[0], 0);
        
        bool result = receiver.asInteger() >= args[0].asInteger();
        return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
    }
}

// Integer class setup utilities
namespace IntegerClassSetup {
    
    void addPrimitiveMethod(Class* clazz, const std::string& selector, int primitiveNumber) {
        Symbol* selectorSymbol = Symbol::intern(selector);
        auto& registry = PrimitiveRegistry::getInstance();
        
        auto primitiveMethod = registry.createPrimitiveMethod(primitiveNumber);
        if (primitiveMethod) {
            clazz->addMethod(selectorSymbol, primitiveMethod);
        } else {
            throw std::runtime_error("Unknown primitive number: " + std::to_string(primitiveNumber));
        }
    }
    
    void addPrimitiveMethods(Class* integerClass) {
        // Arithmetic methods
        addPrimitiveMethod(integerClass, "+", PrimitiveNumbers::INTEGER_ADD);
        addPrimitiveMethod(integerClass, "-", PrimitiveNumbers::INTEGER_SUBTRACT);
        addPrimitiveMethod(integerClass, "*", PrimitiveNumbers::INTEGER_MULTIPLY);
        addPrimitiveMethod(integerClass, "/", PrimitiveNumbers::INTEGER_DIVIDE);
        
        // Comparison methods
        addPrimitiveMethod(integerClass, "<", PrimitiveNumbers::INTEGER_LESS_THAN);
        addPrimitiveMethod(integerClass, ">", PrimitiveNumbers::INTEGER_GREATER_THAN);
        addPrimitiveMethod(integerClass, "=", PrimitiveNumbers::INTEGER_EQUAL);
        addPrimitiveMethod(integerClass, "~=", PrimitiveNumbers::INTEGER_NOT_EQUAL);
        addPrimitiveMethod(integerClass, "<=", PrimitiveNumbers::INTEGER_LESS_THAN_OR_EQUAL);
        addPrimitiveMethod(integerClass, ">=", PrimitiveNumbers::INTEGER_GREATER_THAN_OR_EQUAL);
    }
}

} // namespace smalltalk