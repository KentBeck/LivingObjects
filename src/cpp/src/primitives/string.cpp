#include "../../include/primitives.h"
#include "../../include/smalltalk_string.h"
#include "../../include/smalltalk_class.h"
#include <stdexcept>

namespace smalltalk
{

    // String primitive implementations
    namespace StringPrimitives
    {

        void checkArgumentCount(const std::vector<TaggedValue> &args, size_t expected)
        {
            if (args.size() != expected)
            {
                throw std::runtime_error("Wrong number of arguments: expected " +
                                         std::to_string(expected) + ", got " + std::to_string(args.size()));
            }
        }

        void checkStringReceiver(TaggedValue receiver)
        {
            // Simplified check: just verify it's a pointer to an object
            // We'll trust that the method dispatch system only calls string primitives on strings
            if (!receiver.isPointer())
            {
                throw std::runtime_error("Receiver must be a string (got non-pointer)");
            }

            try
            {
                Object *obj = receiver.asObject();
                if (!obj)
                {
                    throw std::runtime_error("Receiver must be a string (got null object)");
                }
                // For now, we'll trust that if it's a pointer to an object, it's a string
                // since the method dispatch should only call string primitives on strings
            }
            catch (...)
            {
                throw std::runtime_error("Receiver must be a string (exception during object access)");
            }
        }

        void checkStringArgument(TaggedValue arg, size_t index)
        {
            // Simplified check: just verify it's a pointer to an object
            if (!arg.isPointer())
            {
                throw std::runtime_error("Argument " + std::to_string(index) + " must be a string (got non-pointer)");
            }

            try
            {
                Object *obj = arg.asObject();
                if (!obj)
                {
                    throw std::runtime_error("Argument " + std::to_string(index) + " must be a string (got null object)");
                }
            }
            catch (...)
            {
                throw std::runtime_error("Argument " + std::to_string(index) + " must be a string (exception during object access)");
            }
        }

        TaggedValue concatenate(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 1);
            checkStringReceiver(receiver);
            checkStringArgument(args[0], 0);

            String *receiverStr = static_cast<String *>(receiver.asObject());
            String *argStr = static_cast<String *>(args[0].asObject());

            String *result = receiverStr->concatenate(argStr);
            return StringUtils::createTaggedString(result);
        }

        TaggedValue size(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 0);
            checkStringReceiver(receiver);

            String *receiverStr = static_cast<String *>(receiver.asObject());
            int32_t size = static_cast<int32_t>(receiverStr->size());
            return TaggedValue(size);
        }
    }

    // Register string primitives
    void registerStringPrimitives()
    {
        auto &registry = PrimitiveRegistry::getInstance();

        // Register String primitives using standard Smalltalk primitive numbers
        registry.registerPrimitive(PrimitiveNumbers::STRING_CONCAT, StringPrimitives::concatenate);
        registry.registerPrimitive(PrimitiveNumbers::STRING_SIZE, StringPrimitives::size);
    }

} // namespace smalltalk
