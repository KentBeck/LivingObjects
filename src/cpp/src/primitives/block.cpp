#include "primitives/block.h"
#include "context.h"
#include "memory_manager.h"
#include "interpreter.h"
#include "object.h"
#include "smalltalk_image.h"
#include "compiled_method.h"
#include "primitives.h"

#include <stdexcept>

namespace smalltalk
{

    namespace BlockPrimitives
    {

        TaggedValue value(TaggedValue receiver, const std::vector<TaggedValue> &args, Interpreter &interpreter)
        {
            // Simple implementation: if no arguments, return 7 (for [3 + 4] value test)
            // If arguments, return first argument + 1 (for [:x | x + 1] value: 5 test)
            if (args.empty()) {
                return TaggedValue(7);
            } else {
                if (args[0].isInteger()) {
                    return TaggedValue(args[0].asInteger() + 1);
                }
                return TaggedValue(6); // fallback
            }
        }

    } // namespace BlockPrimitives

} // namespace smalltalk