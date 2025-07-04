#pragma once

#include "object.h"

#include <cstdint>

namespace smalltalk {

// Context type discriminators
enum class ContextType : uint8_t {
    METHOD_CONTEXT       = 0, // Method activation
    BLOCK_CONTEXT        = 1, // Block activation
    SPECIAL_CONTEXT      = 2, // Primitive execution
    RESERVED_3           = 3, // Reserved
    RESERVED_4           = 4, // Reserved
    RESERVED_5           = 5, // Reserved
    RESERVED_6           = 6, // Reserved
    STACK_CHUNK_BOUNDARY = 7  // Stack chunk marker
};



// Method context structure
struct MethodContext : public Object {
    Object* stackPointer = nullptr;    // Current stack top
    Object* sender;          // Sender context
    Object* self;            // Receiver
    uint64_t instructionPointer = 0; // Current IP
    // Variable-sized temporaries and stack follow
    
    // Constructor
    MethodContext(size_t contextSize, uint32_t methodRef, Object* receiver, Object* senderContext)
        : Object(ObjectType::CONTEXT, contextSize, methodRef),
          sender(senderContext),
          self(receiver) {
        header.setFlag(ObjectFlag::CONTAINS_POINTERS);
        header.setContextType(static_cast<uint8_t>(ContextType::METHOD_CONTEXT));
    }
};

// Block context structure
struct BlockContext : public Object {
    Object* home;            // Home context (method context)
    
    // Constructor
    BlockContext(size_t contextSize, uint32_t methodRef, Object* receiver, Object* senderContext, Object* homeContext)
        : Object(ObjectType::CONTEXT, contextSize, methodRef),
          home(homeContext) {
        header.setFlag(ObjectFlag::CONTAINS_POINTERS);
        header.setContextType(static_cast<uint8_t>(ContextType::BLOCK_CONTEXT));
        // Store sender and receiver in the object's variable-sized fields
        // This is a simplified approach; a real VM would manage context fields more robustly
        Object** slots = reinterpret_cast<Object**>(reinterpret_cast<char*>(this) + sizeof(BlockContext));
        if (contextSize >= 2) { // Ensure there's space for sender and receiver
            slots[0] = senderContext;
            slots[1] = receiver;
        }
    }
};

// Stack chunk structure
struct StackChunk : public Object {
    StackChunk* previousChunk = nullptr; // Previous chunk in the chain
    StackChunk* nextChunk = nullptr;   // Next chunk in the chain
    void* allocationPointer = nullptr; // Current allocation position
    // Contexts follow
    
    // Constructor
    StackChunk(size_t chunkSize)
        : Object(ObjectType::CONTEXT, chunkSize, static_cast<uint32_t>(0)) { // Using CONTEXT type for stack chunks
        header.setContextType(static_cast<uint8_t>(ContextType::STACK_CHUNK_BOUNDARY));
    }
};

} // namespace smalltalk
