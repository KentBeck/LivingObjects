#include "../../include/primitives.h"
#include "../../include/smalltalk_string.h"
#include "../../include/symbol.h"
#include "../../include/smalltalk_class.h"
#include "../../include/smalltalk_exception.h"
#include <stdexcept>
#include <memory>

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

        /**
         * String at: primitive - returns character at given index (1-based)
         */
        TaggedValue at(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            (void)interpreter; // Suppress unused parameter warning
            checkArgumentCount(args, 1);
            
            // Check that receiver is a string
            if (!StringUtils::isString(receiver)) {
                throw std::runtime_error("Receiver must be a string");
            }
            
            // Check that argument is an integer
            if (!args[0].isInteger()) {
                throw std::runtime_error("Index must be an integer");
            }
            
            int32_t index = args[0].asInteger();
            
            // Get string content
            String* str = static_cast<String*>(receiver.asObject());
            std::string content = str->getContent();
            
            // Smalltalk uses 1-based indexing
            if (index < 1 || index > static_cast<int32_t>(content.length())) {
                // Throw proper IndexError
                auto exception = std::make_unique<IndexError>("Index " + std::to_string(index) + " out of bounds for string of size " + std::to_string(content.length()));
                ExceptionHandler::throwException(std::move(exception));
            }
            
            // Return character at index (convert to Character object)
            char character = content[index - 1]; // Convert to 0-based
            
            // For now, return as integer (character code)
            // In full Smalltalk, this would be a Character object
            return TaggedValue(static_cast<int32_t>(character));
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

        /**
         * String asSymbol primitive - returns interned Symbol for receiver content
         */
        TaggedValue asSymbol(TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter)
        {
            (void)interpreter;
            checkArgumentCount(args, 0);

            if (!StringUtils::isString(receiver)) {
                throw std::runtime_error("Receiver must be a string");
            }
            String* str = static_cast<String*>(receiver.asObject());
            Symbol* sym = Symbol::intern(str->getContent());
            return TaggedValue::fromObject(sym);
        }
    }

    // Register string primitives
    void registerStringPrimitives()
    {
        auto &registry = PrimitiveRegistry::getInstance();

        // Register String primitives using standard Smalltalk primitive numbers
        registry.registerPrimitive(PrimitiveNumbers::STRING_AT, StringPrimitives::at);
        registry.registerPrimitive(PrimitiveNumbers::STRING_CONCAT, StringPrimitives::concatenate);
        registry.registerPrimitive(PrimitiveNumbers::STRING_SIZE, StringPrimitives::size);
        registry.registerPrimitive(PrimitiveNumbers::STRING_AS_SYMBOL, StringPrimitives::asSymbol);
    }

} // namespace smalltalk
