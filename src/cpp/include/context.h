#pragma once

#include "object.h"
#include "tagged_value.h"

#include <cassert>
#include <cstdint>

namespace smalltalk
{

  // Context type discriminators
  enum class ContextType : uint8_t
  {
    METHOD_CONTEXT = 0,       // Method activation
    BLOCK_CONTEXT = 1,        // Block activation
    STACK_CHUNK_BOUNDARY = 2, // Stack chunk marker
    RESERVED_3 = 3,           // Reserved
    RESERVED_4 = 4,           // Reserved
    RESERVED_5 = 5,           // Reserved
    RESERVED_6 = 6,           // Reserved
    RESERVED_7 = 7            // Reserved
  };

  // Forward declaration
  class CompiledMethod;

  // Method context structure
  struct MethodContext : public Object
  {
    TaggedValue *stackPointer = nullptr; // Current stack top
    TaggedValue sender;                  // Sender context (TaggedValue for consistency)
    TaggedValue self;                    // Receiver (TaggedValue for consistency)
    TaggedValue home;                    // Home context for blocks (nil for regular methods)
    uint64_t instructionPointer = 0;     // Current IP
    CompiledMethod *method;              // Direct pointer to the compiled method (required)
    // Variable-sized temporaries and stack follow

    // Constructor
    MethodContext(size_t contextSize, TaggedValue receiver,
                  TaggedValue senderContext, TaggedValue homeContext,
                  CompiledMethod *compiledMethod = nullptr)
        : Object(ObjectType::CONTEXT, contextSize, static_cast<uint32_t>(0)),
          sender(senderContext), self(receiver), home(homeContext),
          method(compiledMethod)
    {
      // Note: compiledMethod can be nullptr during bootstrap
      header.setFlag(ObjectFlag::CONTAINS_POINTERS);
      header.setContextType(static_cast<uint8_t>(ContextType::METHOD_CONTEXT));
    }
  };

  // Block context structure
  struct BlockContext : public Object
  {
    TaggedValue
        home; // Home context (method context) - TaggedValue for consistency

    // Constructor
    BlockContext(size_t contextSize, TaggedValue receiver,
                 TaggedValue senderContext, TaggedValue homeContext)
        : Object(ObjectType::CONTEXT, contextSize, static_cast<uint32_t>(0)),
          home(homeContext)
    {
      header.setFlag(ObjectFlag::CONTAINS_POINTERS);
      header.setContextType(static_cast<uint8_t>(ContextType::BLOCK_CONTEXT));
      // Store sender and receiver in the object's variable-sized fields as
      // TaggedValue
      TaggedValue *slots = reinterpret_cast<TaggedValue *>(
          reinterpret_cast<char *>(this) + sizeof(BlockContext));
      if (contextSize >= 2)
      { // Ensure there's space for sender and receiver
        slots[0] = senderContext;
        slots[1] = receiver;
      }
    }
  };

  // Stack chunk structure
  struct StackChunk : public Object
  {
    StackChunk *previousChunk = nullptr; // Previous chunk in the chain
    StackChunk *nextChunk = nullptr;     // Next chunk in the chain
    void *allocationPointer = nullptr;   // Current allocation position
    // Contexts follow

    // Constructor
    StackChunk(size_t chunkSize)
        : Object(
              ObjectType::CONTEXT, chunkSize,
              static_cast<uint32_t>(0))
    { // Using CONTEXT type for stack chunks
      header.setContextType(
          static_cast<uint8_t>(ContextType::STACK_CHUNK_BOUNDARY));
    }
  };

} // namespace smalltalk
