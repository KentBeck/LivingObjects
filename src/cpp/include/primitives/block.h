#pragma once

#include "tagged_value.h"

#include <vector>

namespace smalltalk {

// Forward declarations
class Interpreter;

/**
 * Primitive method implementations for the Block class.
 */
namespace BlockPrimitives {

    /**
     * Implements Block>>value.
     */
    TaggedValue value(TaggedValue receiver, const std::vector<TaggedValue>& args, Interpreter& interpreter);

} // namespace BlockPrimitives

} // namespace smalltalk