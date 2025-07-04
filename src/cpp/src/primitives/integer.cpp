#include "../../include/primitive_methods.h"
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
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 1);
            checkIntegerReceiver(receiver);
            checkIntegerArgument(args[0], 0);

            int32_t divisor = args[0].asInteger();
            if (divisor == 0)
            {
                throw std::runtime_error("Division by zero");
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

} // namespace smalltalk