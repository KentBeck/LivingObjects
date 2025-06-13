# Smalltalk Context Stack Memory Layout

## Core Context Object Structure

Each context is a 64-bit aligned object with the following structure:

- **Header Word** (64 bits)
  - **Discriminator** (3 bits): Indicates context type
    - 000: Method context
    - 001: Block context
    - 010: Special context (primitives)
    - 011: Reserved
    - 100: Reserved
    - 101: Reserved
    - 110: Reserved
    - 111: Stack chunk boundary marker
  - **Flags** (5 bits): Various state flags
    - Bit 0: Has been materialized
    - Bit 1: Has been scanned by GC
    - Bit 2: Contains pointers
    - Bit 3: Reserved
    - Bit 4: Reserved
  - **Size** (24 bits): Size of context in slots
  - **Method** (32 bits): Reference to method object or block

- **Stack Pointer** (64 bits): Points to current top of stack for this context
- **Sender** (64 bits): Reference to sender context
- **Home** (64 bits): For block contexts, points to method context
- **Self** (64 bits): Reference to receiver object
- **Instruction Pointer** (64 bits): Current execution position
- **Temporaries and Arguments** (variable): Method/block temporaries and arguments
- **Stack** (variable): Evaluation stack growing upward

## Stack Chunk Architecture

Contexts are organized in stack chunks:

- **Stack Chunk** (fixed size, e.g., 64KB)
  - **Header** (64 bits): Chunk metadata
    - Previous chunk pointer
    - Next chunk pointer
    - Allocation pointer
  - **Contexts** (variable): Contiguous contexts
  - **Free Space** (variable): Available for allocation

## Memory Management Strategies

### Allocation

1. Fast allocation from current stack chunk
2. When chunk full, allocate new chunk
3. Link chunks in doubly-linked list

### Growth

1. Contexts grow upward within their chunk
2. Stack overflow triggers materialization into heap

### Context Slot Management

1. Arguments and temporaries stored at fixed offsets
2. Stack slots dynamically allocated/freed

## Stack Operations

### Context Creation

1. Allocate space in current chunk
2. Initialize header and pointers
3. Copy arguments from caller

### Context Switching

1. Save current IP and SP
2. Load new context's IP and SP
3. Resume execution

### Return Handling

1. Restore sender context
2. Push return value on sender's stack
3. Resume execution at sender's IP

## Optimization Strategies

### Context Caching

1. Keep pool of recycled contexts
2. Avoid full allocation for common patterns

### Lazy Materialization

1. Contexts start as stack-only
2. Only materialize to heap when:
   - Returned as value
   - Survived long enough
   - Referenced by heap object

### GC Scanning

1. Stack chunks scanned in-place
2. Materialized contexts scanned like normal objects