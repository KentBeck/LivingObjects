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

// Context flags
enum class ContextFlag : uint8_t {
    MATERIALIZED         = 0, // Context has been materialized to heap
    GC_SCANNED           = 1, // Context has been scanned by GC
    CONTAINS_POINTERS    = 2, // Context contains pointers
    RESERVED_3           = 3, // Reserved
    RESERVED_4           = 4  // Reserved
};

// Context header structure (64 bits)
struct ContextHeader {
    uint64_t size : 24;      // Size in slots
    uint64_t flags : 5;      // Various flags
    uint64_t type : 3;       // Context type
    uint64_t method : 32;    // Reference to method object or block
    
    // Constructor
    ContextHeader(ContextType contextType, size_t contextSize, uint32_t methodRef = 0)
        : size(contextSize), flags(0), type(static_cast<uint64_t>(contextType)), method(methodRef) {}
    
    // Flag operations
    void setFlag(ContextFlag flag) {
        flags |= (1 << static_cast<uint8_t>(flag));
    }
    
    bool hasFlag(ContextFlag flag) const {
        return (flags & (1 << static_cast<uint8_t>(flag))) != 0;
    }
    
    void clearFlag(ContextFlag flag) {
        flags &= ~(1 << static_cast<uint8_t>(flag));
    }
};

// Method context structure
struct MethodContext {
    ContextHeader header;
    Object* stackPointer;    // Current stack top
    Object* sender;          // Sender context
    Object* self;            // Receiver
    uint64_t instructionPointer; // Current IP
    // Variable-sized temporaries and stack follow
    
    // Constructor
    MethodContext(size_t contextSize, uint32_t methodRef, Object* receiver, Object* senderContext)
        : header(ContextType::METHOD_CONTEXT, contextSize, methodRef),
          stackPointer(nullptr),
          sender(senderContext),
          self(receiver),
          instructionPointer(0) {
        header.setFlag(ContextFlag::CONTAINS_POINTERS);
    }
};

// Block context structure
struct BlockContext : public MethodContext {
    Object* home;            // Home context (method context)
    
    // Constructor
    BlockContext(size_t contextSize, uint32_t methodRef, Object* receiver, Object* senderContext, Object* homeContext)
        : MethodContext(contextSize, methodRef, receiver, senderContext),
          home(homeContext) {
        header.type = static_cast<uint64_t>(ContextType::BLOCK_CONTEXT);
    }
};

// Stack chunk structure
struct StackChunk {
    ContextHeader header;    // Chunk header with STACK_CHUNK_BOUNDARY type
    StackChunk* previousChunk; // Previous chunk in the chain
    StackChunk* nextChunk;   // Next chunk in the chain
    void* allocationPointer; // Current allocation position
    // Contexts follow
    
    // Constructor
    StackChunk(size_t chunkSize)
        : header(ContextType::STACK_CHUNK_BOUNDARY, chunkSize),
          previousChunk(nullptr),
          nextChunk(nullptr),
          allocationPointer(nullptr) {}
};

} // namespace smalltalk
