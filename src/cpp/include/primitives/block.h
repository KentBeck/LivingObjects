#pragma once

#include "tagged_value.h"

#include <vector>

namespace smalltalk {

// Forward declarations
class Interpreter;

/**
 * Primitive numbers for Block operations.
 */
namespace PrimitiveNumbers {
    constexpr int BLOCK_VALUE = 81;
} // namespace PrimitiveNumbers

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