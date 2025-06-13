# C++ Bytecode Interpreter Implementation Plan

This document outlines the implementation plan for a C++ bytecode interpreter for the Living Objects Smalltalk system, focusing on achieving basic execution capabilities with the necessary components for a functioning Smalltalk environment.

## Phase 1: Core Memory System

### Memory Management

1. **Object Memory**
   - Define base object structure with header word
   - Implement object allocation routines
   - Create memory page management system

2. **Stop & Copy Garbage Collector**
   - Implement two-space memory model
   - Add forwarding pointer mechanism
   - Create object copying routines
   - Implement root set identification
   - Add mark and sweep phases

3. **Object Format**
   - Define object header format (type, size, flags)
   - Implement object reference handling
   - Add support for immediate values (SmallInteger, Character, Boolean)

### Timeline: 2-3 weeks

## Phase 2: Object System

1. **Class Structure**
   - Implement metaclass hierarchy
   - Create class object structure
   - Add method dictionary support
   - Implement instance variable layout

2. **Instance Creation**
   - Add instance allocation
   - Implement instance initialization
   - Create factory methods

3. **Basic Object Protocol**
   - Implement identity comparison
   - Add class access
   - Create primitive method support

### Timeline: 2-3 weeks

## Phase 3: Context & Execution

1. **Context Implementation**
   - Implement context stack memory layout as specified
   - Create context allocation/deallocation
   - Add support for method and block contexts
   - Implement context switching

2. **Bytecode Interpreter**
   - Create bytecode dispatch table
   - Implement bytecodes from the Go implementation
   - Add stack manipulation instructions
   - Implement control flow instructions
   - Add object instantiation instructions

3. **Message Sending**
   - Implement message lookup
   - Add method activation
   - Create return handling
   - Implement primitive invocation
   - Add super sends

### Timeline: 3-4 weeks

## Phase 4: Base Classes & Primitives

1. **Core Classes**
   - Implement Object class
   - Add Collection hierarchy
   - Create String and Symbol
   - Implement Array and ByteArray
   - Add Dictionary support

2. **Primitive Methods**
   - Implement arithmetic primitives
   - Add object creation primitives
   - Create comparison primitives
   - Implement collection access primitives
   - Add I/O primitives

3. **Block Support**
   - Implement block creation
   - Add block evaluation
   - Create non-local return handling
   - Implement block closure support

### Timeline: 3-4 weeks

## Phase 5: Image System

1. **Object Serialization**
   - Implement object graph traversal
   - Create reference mapping
   - Add object serialization format

2. **Image File Format**
   - Define image file header
   - Implement version information
   - Create object table

3. **Image Operations**
   - Implement image saving
   - Add image loading
   - Create snapshot capability
   - Implement startup sequence

### Timeline: 2-3 weeks

## Phase 6: Integration & Testing

1. **Test Harness**
   - Create unit test framework
   - Implement test runners
   - Add automated test suites

2. **Benchmarking**
   - Create performance benchmarks
   - Implement memory usage tracking
   - Add execution profiling

3. **Debugging Support**
   - Implement execution tracing
   - Add breakpoint mechanism
   - Create object inspectors
   - Implement error handling

### Timeline: 2-3 weeks

## Implementation Details

### Memory Management

```cpp
// Object header structure
struct ObjectHeader {
    uint64_t size : 24;      // Size in slots
    uint64_t flags : 5;      // Various flags
    uint64_t type : 3;       // Object type
    uint64_t hash : 32;      // Identity hash
};

// Base object structure
struct Object {
    ObjectHeader header;
    // Variable-sized data follows
};

// Memory manager
class MemoryManager {
public:
    // Allocation
    Object* allocateObject(size_t size);
    
    // Garbage collection
    void collectGarbage();
    
    // Memory spaces
    void* fromSpace;
    void* toSpace;
    size_t spaceSize;
};
```

### Context Structure

```cpp
// Context header
struct ContextHeader {
    uint64_t size : 24;      // Size in slots
    uint64_t flags : 5;      // Various flags
    uint64_t type : 3;       // Context type
    uint64_t method : 32;    // Method reference
};

// Method context
struct MethodContext {
    ContextHeader header;
    Object* stackPointer;    // Current stack top
    Object* sender;          // Sender context
    Object* self;            // Receiver
    uint64_t instructionPointer; // Current IP
    // Variable-sized temporaries and stack follow
};

// Stack chunk
struct StackChunk {
    Object* previousChunk;
    Object* nextChunk;
    void* allocationPointer;
    // Contexts follow
};
```

### Bytecode Definitions

```cpp
// Bytecode constants based on Go implementation
enum Bytecode : uint8_t {
    PUSH_LITERAL             = 0,  // Push a literal from the literals array (followed by 4-byte index)
    PUSH_INSTANCE_VARIABLE   = 1,  // Push an instance variable value (followed by 4-byte offset)
    PUSH_TEMPORARY_VARIABLE  = 2,  // Push a temporary variable value (followed by 4-byte offset)
    PUSH_SELF                = 3,  // Push self onto the stack
    STORE_INSTANCE_VARIABLE  = 4,  // Store a value into an instance variable (followed by 4-byte offset)
    STORE_TEMPORARY_VARIABLE = 5,  // Store a value into a temporary variable (followed by 4-byte offset)
    SEND_MESSAGE             = 6,  // Send a message (followed by 4-byte selector index and 4-byte arg count)
    RETURN_STACK_TOP         = 7,  // Return the value on top of the stack
    JUMP                     = 8,  // Jump to a different bytecode (followed by 4-byte target)
    JUMP_IF_TRUE             = 9,  // Jump if top of stack is true (followed by 4-byte target)
    JUMP_IF_FALSE            = 10, // Jump if top of stack is false (followed by 4-byte target)
    POP                      = 11, // Pop the top value from the stack
    DUPLICATE                = 12, // Duplicate the top value on the stack
    CREATE_BLOCK             = 13, // Create a block (followed by 4-byte bytecode size, 4-byte literal count, 4-byte temp var count)
    EXECUTE_BLOCK            = 14  // Execute a block (followed by 4-byte arg count)
};

// Instruction size (in bytes, including opcode)
int getInstructionSize(Bytecode bytecode) {
    switch (bytecode) {
        case PUSH_LITERAL:
        case PUSH_INSTANCE_VARIABLE:
        case PUSH_TEMPORARY_VARIABLE:
        case STORE_INSTANCE_VARIABLE:
        case STORE_TEMPORARY_VARIABLE:
        case JUMP:
        case JUMP_IF_TRUE:
        case JUMP_IF_FALSE:
        case EXECUTE_BLOCK:
            return 5; // 1 byte opcode + 4 byte operand
        case SEND_MESSAGE:
            return 9; // 1 byte opcode + 4 byte selector index + 4 byte arg count
        case CREATE_BLOCK:
            return 13; // 1 byte opcode + 4 byte bytecode size + 4 byte literal count + 4 byte temp var count
        case PUSH_SELF:
        case RETURN_STACK_TOP:
        case POP:
        case DUPLICATE:
            return 1; // 1 byte opcode
        default:
            return 1; // Default to 1 byte for unknown bytecodes
    }
}
```

### Bytecode Interpreter

```cpp
// Interpreter
class Interpreter {
public:
    // Execution
    void executeMethod(Object* method, Object* receiver, Object** args, int argCount);
    void executeContext(MethodContext* context);
    
    // Bytecode dispatch
    void dispatch(Bytecode bytecode);
    
    // Bytecode handlers
    void handlePushLiteral(uint32_t index);
    void handlePushInstanceVariable(uint32_t offset);
    void handlePushTemporaryVariable(uint32_t offset);
    void handlePushSelf();
    void handleStoreInstanceVariable(uint32_t offset);
    void handleStoreTemporaryVariable(uint32_t offset);
    void handleSendMessage(uint32_t selectorIndex, uint32_t argCount);
    void handleReturnStackTop();
    void handleJump(uint32_t target);
    void handleJumpIfTrue(uint32_t target);
    void handleJumpIfFalse(uint32_t target);
    void handlePop();
    void handleDuplicate();
    void handleCreateBlock(uint32_t bytecodeSize, uint32_t literalCount, uint32_t tempVarCount);
    void handleExecuteBlock(uint32_t argCount);
    
    // Current state
    MethodContext* activeContext;
    StackChunk* currentChunk;
};
```

### Image System

```cpp
// Image operations
class ImageManager {
public:
    // Save current state to file
    bool saveImage(const char* filename);
    
    // Load state from file
    bool loadImage(const char* filename);
    
    // Create snapshot
    Object* createSnapshot();
    
    // Startup sequence
    void startup();
};
```

## Getting to Basic Execution

To achieve the minimum viable product with basic execution capabilities:

1. Implement core memory system with object allocation
2. Create basic class structure for Object, Class, Method
3. Implement context stack memory layout 
4. Add bytecode interpreter with the 15 bytecodes from the Go implementation
5. Implement message sending with method lookup
6. Add primitive method support for essential operations
7. Create minimal image save/load capability

This focused approach will allow running simple Smalltalk expressions like:

```smalltalk
3 + 4.
'Hello' , ' World'.
OrderedCollection new add: 1; add: 2; yourself.
```

With these core components in place, the system can be incrementally enhanced to support more complex Smalltalk functionality.