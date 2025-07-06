#include "../../include/primitives.h"
#include "../../include/smalltalk_class.h"
#include "../../include/symbol.h"
#include "../../include/compiled_method.h"
#include "../../include/smalltalk_exception.h"
#include <stdexcept>

namespace smalltalk
{

    // Integer primitive implementations
    namespace IntegerPrimitives
    {

        void checkArgumentCount(const std::vector<TaggedValue> &args, size_t expected)
        {
            if (args.size() != expected)
            {
                throw std::runtime_error("Wrong number of arguments: expected " +
                                         std::to_string(expected) + ", got " + std::to_string(args.size()));
            }
        }

        void checkIntegerReceiver(TaggedValue receiver)
        {
            if (!receiver.isInteger())
            {
                throw std::runtime_error("Receiver must be an integer");
            }
        }

        void checkIntegerArgument(TaggedValue arg, size_t index)
        {
            if (!arg.isInteger())
            {
                throw std::runtime_error("Argument " + std::to_string(index) + " must be an integer");
            }
        }

        TaggedValue add(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 1);
            checkIntegerReceiver(receiver);
            checkIntegerArgument(args[0], 0);

            int32_t result = receiver.asInteger() + args[0].asInteger();
            return TaggedValue(result);
        }

        TaggedValue subtract(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 1);
            checkIntegerReceiver(receiver);
            checkIntegerArgument(args[0], 0);

            int32_t result = receiver.asInteger() - args[0].asInteger();
            return TaggedValue(result);
        }

        TaggedValue multiply(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 1);
            checkIntegerReceiver(receiver);
            checkIntegerArgument(args[0], 0);

            int32_t result = receiver.asInteger() * args[0].asInteger();
            return TaggedValue(result);
        }

        TaggedValue divide(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            checkArgumentCount(args, 1);
            checkIntegerReceiver(receiver);
            checkIntegerArgument(args[0], 0);

            int32_t divisor = args[0].asInteger();
            if (divisor == 0)
            {
                // Throw proper Smalltalk exception directly
                throw std::runtime_error("ZeroDivisionError");
            }

            int32_t result = receiver.asInteger() / divisor;
            return TaggedValue(result);
        }

        TaggedValue lessThan(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 1);
            checkIntegerReceiver(receiver);
            checkIntegerArgument(args[0], 0);

            bool result = receiver.asInteger() < args[0].asInteger();
            return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
        }

        TaggedValue greaterThan(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 1);
            checkIntegerReceiver(receiver);
            checkIntegerArgument(args[0], 0);

            bool result = receiver.asInteger() > args[0].asInteger();
            return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
        }

        TaggedValue equal(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 1);
            checkIntegerReceiver(receiver);
            checkIntegerArgument(args[0], 0);

            bool result = receiver.asInteger() == args[0].asInteger();
            return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
        }

        TaggedValue notEqual(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 1);
            checkIntegerReceiver(receiver);
            checkIntegerArgument(args[0], 0);

            bool result = receiver.asInteger() != args[0].asInteger();
            return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
        }

        TaggedValue lessThanOrEqual(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 1);
            checkIntegerReceiver(receiver);
            checkIntegerArgument(args[0], 0);

            bool result = receiver.asInteger() <= args[0].asInteger();
            return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
        }

        TaggedValue greaterThanOrEqual(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 1);
            checkIntegerReceiver(receiver);
            checkIntegerArgument(args[0], 0);

            bool result = receiver.asInteger() >= args[0].asInteger();
            return result ? TaggedValue::trueValue() : TaggedValue::falseValue();
        }
    }

    // IntegerClassSetup implementation
    namespace IntegerClassSetup
    {
        void addPrimitiveMethod(Class* clazz, const std::string& selector, int primitiveNumber)
        {
            // Create a symbol for the selector
            Symbol* selectorSymbol = Symbol::intern(selector);
            
            // Create a CompiledMethod with the primitive number
            auto primitiveMethod = std::make_shared<CompiledMethod>();
            primitiveMethod->primitiveNumber = primitiveNumber;
            
            // Add the method to the class
            clazz->addMethod(selectorSymbol, primitiveMethod);
        }
        
        void addPrimitiveMethods(Class* integerClass)
        {
            if (!integerClass) {
                throw std::runtime_error("Integer class is null");
            }
            
            // Add arithmetic primitive methods using correct primitive numbers
            addPrimitiveMethod(integerClass, "+", PrimitiveNumbers::SMALL_INT_ADD);
            addPrimitiveMethod(integerClass, "-", PrimitiveNumbers::SMALL_INT_SUB);
            addPrimitiveMethod(integerClass, "*", PrimitiveNumbers::SMALL_INT_MUL);
            addPrimitiveMethod(integerClass, "/", PrimitiveNumbers::SMALL_INT_DIV);
            
            // Add comparison primitive methods
            addPrimitiveMethod(integerClass, "<", PrimitiveNumbers::SMALL_INT_LT);
            addPrimitiveMethod(integerClass, ">", PrimitiveNumbers::SMALL_INT_GT);
            addPrimitiveMethod(integerClass, "=", PrimitiveNumbers::SMALL_INT_EQ);
            addPrimitiveMethod(integerClass, "~=", PrimitiveNumbers::SMALL_INT_NE);
            addPrimitiveMethod(integerClass, "<=", PrimitiveNumbers::SMALL_INT_LE);
            addPrimitiveMethod(integerClass, ">=", PrimitiveNumbers::SMALL_INT_GE);
        }
    }


} // namespace smalltalk