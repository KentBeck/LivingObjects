# Smalltalk Virtual Machine Specification

**Version**: 1.0  
**Date**: 2025-01-05  
**Purpose**: Complete specification for implementing compatible Smalltalk virtual machines

## Table of Contents

1. [Overview](#overview)
2. [Tagged Value Representation](#tagged-value-representation)
3. [Object Memory Format](#object-memory-format)
4. [Image File Format](#image-file-format)
5. [Bytecode Instruction Set](#bytecode-instruction-set)
6. [Method Context and Execution](#method-context-and-execution)
7. [Primitive Methods](#primitive-methods)
8. [Exception Handling](#exception-handling)
9. [Implementation Notes](#implementation-notes)

---

## 1. Overview

This specification defines the binary formats, instruction set, and execution semantics for a Smalltalk virtual machine. The design prioritizes:

- **Simplicity**: Easy to implement and understand
- **Correctness**: Well-defined semantics with no ambiguity
- **Compatibility**: Multiple VM implementations can share images
- **Efficiency**: Reasonable performance without excessive complexity

### 1.1 Design Principles

- **Tagged values** for immediate integers, booleans, and nil
- **Uniform object representation** with headers
- **Stack-based bytecode** execution
- **Context-based** method activation
- **Primitive methods** for performance-critical operations

### 1.2 Terminology

- **Object**: Any value in the system (immediate or heap-allocated)
- **Tagged Value**: Efficient representation using pointer tagging
- **Context**: Execution state for a method or block
- **Bytecode**: Low-level instructions executed by the VM
- **Primitive**: Native implementation of a method
- **Image**: Serialized snapshot of the object memory

---

## 2. Tagged Value Representation

Tagged values use the bottom 2 bits of a pointer-sized value to encode type information.

### 2.1 Tag Encoding

```
Bit Pattern    Type           Description
-----------    ----           -----------
xxxxxxxx00     Pointer        Heap-allocated object (4-byte aligned)
xxxxxxxx01     Special        nil, true, false
xxxxxxxx10     Float          Inline float (limited range)
xxxxxxxx11     SmallInteger   31-bit signed integer
```

### 2.2 Special Values

```
Value          Encoding       Description
-----          --------       -----------
nil            0x00000001     The nil object
true           0x00000005     Boolean true
false          0x00000009     Boolean false
```

### 2.3 SmallInteger Encoding

SmallIntegers are encoded by shifting the value left 2 bits and setting the tag to `11`:

```
Value: -2147483648 to 2147483647 (31-bit signed)
Encoding: (value << 2) | 0x03
Decoding: (encoded >> 2) as signed 32-bit
```

**Examples**:

- `0` → `0x00000003`
- `1` → `0x00000007`
- `-1` → `0xFFFFFFFB`
- `42` → `0x000000AB`

### 2.4 Pointer Encoding

Heap-allocated objects must be 4-byte aligned. The pointer is stored directly with tag `00`:

```
Encoding: pointer | 0x00
Decoding: encoded & ~0x03
```

### 2.5 Float Encoding (Simplified)

For simplicity, only a limited set of floats are encoded inline:

- `0.0` → `0x00000002`
- `1.0` → `0x00000006`
- `-1.0` → `0x0000000A`

Other floats require heap allocation (not covered in this simplified spec).

---

## 3. Object Memory Format

All heap-allocated objects share a common header structure followed by type-specific data.

### 3.1 Object Header (64 bits)

```
Bits    Field       Description
----    -----       -----------
0-23    size        Object size in slots or bytes (24 bits)
24-31   flags       Type (3 bits) + flags (5 bits)
32-63   hash        Identity hash (32 bits)
```

### 3.2 Object Types

```
Value   Type            Description
-----   ----            -----------
0       IMMEDIATE       SmallInteger, Character, Boolean (not used in headers)
1       OBJECT          General object with instance variables
2       ARRAY           Indexable collection without instance variables
3       BYTE_ARRAY      Byte-indexed collection
4       SYMBOL          Interned string
5       CONTEXT         Method or block execution context
6       CLASS           Class object
7       METHOD          Compiled method
```

### 3.3 Object Flags

```
Bit     Flag                    Description
---     ----                    -----------
0       MARKED                  Marked during garbage collection
1       REMEMBERED              In remembered set (generational GC)
2       IMMUTABLE               Cannot be modified
3       FORWARDED               Object has been moved (forwarding pointer)
4       PINNED                  Cannot be moved by GC
5       CONTAINS_POINTERS       Object contains references to other objects
6       TAGGED_VALUE_WRAPPER    Wraps an immediate value
```

### 3.4 General Object Layout

```
Offset  Size    Field
------  ----    -----
0       8       ObjectHeader
8       8       Class pointer
16      ...     Instance variables (TaggedValue array)
```

### 3.5 Array Layout

```
Offset  Size    Field
------  ----    -----
0       8       ObjectHeader (size = number of elements)
8       8       Class pointer
16      ...     Elements (TaggedValue array)
```

### 3.6 ByteArray Layout

```
Offset  Size    Field
------  ----    -----
0       8       ObjectHeader (size = number of bytes)
8       8       Class pointer
16      ...     Bytes (uint8_t array)
```

### 3.7 Symbol Layout

```
Offset  Size    Field
------  ----    -----
0       8       ObjectHeader
8       8       Class pointer
16      4       String length (uint32_t)
20      ...     String data (UTF-8 bytes, null-terminated)
```

### 3.8 Context Layout

```
Offset  Size    Field
------  ----    -----
0       8       ObjectHeader
8       8       Class pointer
16      8       self (TaggedValue)
24      8       sender (TaggedValue - pointer to sender context)
32      8       home (TaggedValue - pointer to home context for blocks)
40      8       method (pointer to CompiledMethod)
48      4       instructionPointer (uint32_t)
52      4       (padding)
56      8       stackPointer (pointer to TaggedValue)
64      ...     Stack and temporary variables (TaggedValue array)
```

### 3.9 CompiledMethod Layout

```
Offset  Size    Field
------  ----    -----
0       8       ObjectHeader
8       8       Class pointer
16      4       primitiveNumber (uint32_t, 0 = no primitive)
20      4       numArgs (uint32_t)
24      4       numTemps (uint32_t)
28      4       homeVarCount (uint32_t - for nested blocks)
32      4       bytecodeSize (uint32_t)
36      4       literalCount (uint32_t)
40      ...     Bytecodes (uint8_t array)
...     ...     Literals (TaggedValue array)
...     ...     Temporary variable names (for debugging)
```

---

## 4. Image File Format

Smalltalk images are binary files containing a serialized object graph.

### 4.1 Image Header

```
Offset  Size    Field               Description
------  ----    -----               -----------
0       4       magic               0x53544C4B ("STLK")
4       4       version             Image format version (1)
8       8       creationTime        Unix timestamp (uint64_t)
16      8       modificationTime    Unix timestamp (uint64_t)
24      4       classCount          Number of classes (uint32_t)
28      4       methodCount         Number of methods (uint32_t)
32      4       globalCount         Number of global variables (uint32_t)
36      4       metadataCount       Number of metadata entries (uint32_t)
40      8       dataOffset          Offset to object data (uint64_t)
48      ...     Reserved for future use
```

### 4.2 Class Table Section

Located immediately after header. Each entry:

```
Offset  Size    Field
------  ----    -----
0       4       nameLength (uint32_t)
4       ...     className (UTF-8 string)
...     4       superclassIndex (uint32_t, 0xFFFFFFFF = no superclass)
...     4       instanceVariableCount (uint32_t)
...     ...     instanceVariableNames (array of length-prefixed strings)
...     4       methodCount (uint32_t)
...     ...     methods (array of method references)
```

### 4.3 Global Dictionary Section

```
Offset  Size    Field
------  ----    -----
0       4       keyLength (uint32_t)
4       ...     key (UTF-8 string)
...     8       value (TaggedValue)
```

### 4.4 Metadata Section

```
Offset  Size    Field
------  ----    -----
0       4       keyLength (uint32_t)
4       ...     key (UTF-8 string)
...     4       valueLength (uint32_t)
...     ...     value (UTF-8 string)
```

### 4.5 Object Data Section

Objects are serialized in dependency order (referenced objects before referencing objects).

Each object:

```
Offset  Size    Field
------  ----    -----
0       4       objectID (uint32_t)
4       1       objectType (uint8_t)
5       3       (padding)
8       ...     object-specific data
```

### 4.6 Image Loading Process

1. Read and validate header (check magic number and version)
2. Load class table and create class objects
3. Load global dictionary
4. Load metadata
5. Load object data section, resolving references
6. Reconnect primitive methods to native implementations
7. Validate image integrity

---

## 5. Bytecode Instruction Set

All bytecodes are 1 byte opcodes followed by operands. Multi-byte operands are little-endian uint32.

### 5.1 Instruction Format

```
Opcode  Name                        Operands        Size    Description
------  ----                        --------        ----    -----------
0       PUSH_LITERAL                index(4)        5       Push literal[index] onto stack
1       PUSH_INSTANCE_VARIABLE      offset(4)       5       Push instance variable at offset
2       PUSH_TEMPORARY_VARIABLE     index(4)        5       Push temporary variable at index
3       PUSH_SELF                   -               1       Push self onto stack
4       STORE_INSTANCE_VARIABLE     offset(4)       5       Store top of stack to instance variable
5       STORE_TEMPORARY_VARIABLE    index(4)        5       Store top of stack to temporary variable
6       SEND_MESSAGE                sel(4), argc(4) 9       Send message with selector and arg count
7       RETURN_STACK_TOP            -               1       Return top of stack to sender
8       JUMP                        target(4)       5       Unconditional jump to bytecode offset
9       JUMP_IF_TRUE                target(4)       5       Jump if top of stack is true
10      JUMP_IF_FALSE               target(4)       5       Jump if top of stack is false
11      POP                         -               1       Pop and discard top of stack
12      DUPLICATE                   -               1       Duplicate top of stack
13      CREATE_BLOCK                lit(4), par(4)  9       Create block closure
14      EXECUTE_BLOCK               argc(4)         5       Execute block with arguments
```

### 5.2 Detailed Instruction Semantics

#### 5.2.1 PUSH_LITERAL (0)

**Format**: `00 [index:uint32]`

**Operation**:

1. Read 4-byte literal index from bytecode
2. Fetch `literals[index]` from current method
3. Push value onto stack

**Stack**: `... → ..., value`

**Errors**:

- Index out of bounds: throw error

**Example**:

```
Bytecode: 00 02 00 00 00
Effect: Push literals[2] onto stack
```

#### 5.2.2 PUSH_INSTANCE_VARIABLE (1)

**Format**: `01 [offset:uint32]`

**Operation**:

1. Read 4-byte offset from bytecode
2. Get receiver (self) from current context
3. Read instance variable at offset from receiver
4. Push value onto stack

**Stack**: `... → ..., value`

**Errors**:

- Offset out of bounds: throw error
- Receiver is not an object: throw error

**Notes**:

- Offset is in TaggedValue slots, not bytes
- First instance variable is at offset 0

#### 5.2.3 PUSH_TEMPORARY_VARIABLE (2)

**Format**: `02 [index:uint32]`

**Operation**:

1. Read 4-byte index from bytecode
2. If index < homeVarCount, resolve through home context chain
3. Otherwise, read from current context's temporary variables
4. Push value onto stack

**Stack**: `... → ..., value`

**Errors**:

- Index out of bounds: throw error

**Notes**:

- Temporary variables include method parameters
- Parameters come first, then local temporaries
- Nested blocks may reference outer context variables

#### 5.2.4 PUSH_SELF (3)

**Format**: `03`

**Operation**:

1. Get self from current context
2. Push self onto stack

**Stack**: `... → ..., self`

**Errors**: None

#### 5.2.5 STORE_INSTANCE_VARIABLE (4)

**Format**: `04 [offset:uint32]`

**Operation**:

1. Read 4-byte offset from bytecode
2. Pop value from stack
3. Get receiver (self) from current context
4. Store value to instance variable at offset
5. Push value back onto stack (assignment returns value)

**Stack**: `..., value → ..., value`

**Errors**:

- Offset out of bounds: throw error
- Receiver is not an object: throw error
- Object is immutable: throw error

#### 5.2.6 STORE_TEMPORARY_VARIABLE (5)

**Format**: `05 [index:uint32]`

**Operation**:

1. Read 4-byte index from bytecode
2. Pop value from stack
3. If index < homeVarCount, resolve through home context chain
4. Otherwise, store to current context's temporary variables
5. Push value back onto stack (assignment returns value)

**Stack**: `..., value → ..., value`

**Errors**:

- Index out of bounds: throw error

**Notes**:

- Nested blocks may write to outer context variables
- Assignment always returns the assigned value

#### 5.2.7 SEND_MESSAGE (6)

**Format**: `06 [selectorIndex:uint32] [argCount:uint32]`

**Operation**:

1. Read 4-byte selector index from bytecode
2. Read 4-byte argument count from bytecode
3. Pop argCount arguments from stack (last pushed = last argument)
4. Pop receiver from stack
5. Look up method for selector in receiver's class
6. If method has primitive, try primitive first
7. If primitive fails or no primitive, execute method bytecode
8. Push result onto stack

**Stack**: `..., receiver, arg1, ..., argN → ..., result`

**Errors**:

- Selector not found: throw MessageNotUnderstood
- Wrong number of arguments: throw error

**Notes**:

- Arguments are popped in reverse order (LIFO)
- Method lookup follows inheritance chain
- Primitive failure falls back to Smalltalk code

#### 5.2.8 RETURN_STACK_TOP (7)

**Format**: `07`

**Operation**:

1. Pop return value from stack
2. If no sender context, end execution with value
3. Otherwise, switch to sender context
4. Push return value onto sender's stack

**Stack**: `..., value → (sender) ..., value`

**Errors**: None

**Notes**:

- Returns from current method/block to sender
- Implicit return at end of method returns self if stack empty

#### 5.2.9 JUMP (8)

**Format**: `08 [target:uint32]`

**Operation**:

1. Read 4-byte target offset from bytecode
2. Set instruction pointer to target

**Stack**: `... → ...` (unchanged)

**Errors**:

- Target out of bounds: throw error

**Notes**:

- Target is absolute bytecode offset in current method
- Used for loops and control flow

#### 5.2.10 JUMP_IF_TRUE (9)

**Format**: `09 [target:uint32]`

**Operation**:

1. Read 4-byte target offset from bytecode
2. Pop value from stack
3. If value is true, set instruction pointer to target
4. Otherwise, continue to next instruction

**Stack**: `..., value → ...`

**Errors**:

- Value is not a boolean: throw error
- Target out of bounds: throw error

**Notes**:

- Only true (not truthy) causes jump
- Value is consumed (popped)

#### 5.2.11 JUMP_IF_FALSE (10)

**Format**: `0A [target:uint32]`

**Operation**:

1. Read 4-byte target offset from bytecode
2. Pop value from stack
3. If value is false, set instruction pointer to target
4. Otherwise, continue to next instruction

**Stack**: `..., value → ...`

**Errors**:

- Value is not a boolean: throw error
- Target out of bounds: throw error

**Notes**:

- Only false (not falsy) causes jump
- Value is consumed (popped)

#### 5.2.12 POP (11)

**Format**: `0B`

**Operation**:

1. Pop and discard top value from stack

**Stack**: `..., value → ...`

**Errors**:

- Stack underflow: throw error

**Notes**:

- Used to discard expression results
- Statement separator in Smalltalk (`.`) generates POP

#### 5.2.13 DUPLICATE (12)

**Format**: `0C`

**Operation**:

1. Read top value from stack without popping
2. Push copy of value onto stack

**Stack**: `..., value → ..., value, value`

**Errors**:

- Stack underflow: throw error

**Notes**:

- Used for cascaded messages
- Value is not deep-copied, just reference duplicated

#### 5.2.14 CREATE_BLOCK (13)

**Format**: `0D [literalIndex:uint32] [paramCount:uint32]`

**Operation**:

1. Read 4-byte literal index (points to block's CompiledMethod)
2. Read 4-byte parameter count
3. Get block's CompiledMethod from literals[literalIndex]
4. Create BlockContext with:
   - self = current context's self
   - home = current context
   - method = block's CompiledMethod
5. Push BlockContext onto stack

**Stack**: `... → ..., blockContext`

**Errors**:

- Literal index out of bounds: throw error
- Literal is not a CompiledMethod: throw error

**Notes**:

- Block captures home context for closure variables
- Block is not executed until `value` message sent
- Nested blocks create chain of home contexts

#### 5.2.15 EXECUTE_BLOCK (14)

**Format**: `0E [argCount:uint32]`

**Operation**:

1. Read 4-byte argument count
2. Pop argCount arguments from stack
3. Pop block context from stack
4. Create new execution context for block:
   - Copy arguments to block's temporary variables
   - Set sender to current context
   - Set instruction pointer to 0
5. Switch to block context and begin execution
6. When block returns, push result onto sender's stack

**Stack**: `..., block, arg1, ..., argN → ..., result`

**Errors**:

- Block is not a BlockContext: throw error
- Wrong number of arguments: throw error

**Notes**:

- Block execution can access home context variables
- Block return returns to sender, not home
- Used by Block>>value, Block>>value:, etc.

---

## 6. Method Context and Execution

### 6.1 Context Creation

When a method is invoked:

1. **Allocate MethodContext** with space for:

   - Fixed context fields (64 bytes)
   - Temporary variables (numTemps slots)
   - Stack space (implementation-defined, minimum 16 slots)

2. **Initialize context fields**:

   - `self` = receiver
   - `sender` = current active context (or nil)
   - `home` = nil (for methods) or home context (for blocks)
   - `method` = CompiledMethod being executed
   - `instructionPointer` = 0
   - `stackPointer` = start of stack (after temp vars)

3. **Copy arguments** to first N temporary variables

4. **Initialize remaining temporaries** to nil

5. **Switch to new context** and begin execution

### 6.2 Stack Management

The stack grows upward in memory:

```
Context Layout:
+------------------+  ← Context base
| Fixed fields     |  (64 bytes)
+------------------+  ← Temporary variables start
| Temp var 0       |  (method parameters first)
| Temp var 1       |
| ...              |
| Temp var N-1     |
+------------------+  ← Stack start (stackPointer initially points here)
| Stack slot 0     |
| Stack slot 1     |
| ...              |
+------------------+  ← Stack end
```

**Stack Operations**:

- **Push**: Store value at stackPointer, increment stackPointer
- **Pop**: Decrement stackPointer, read value
- **Top**: Read value at stackPointer-1 without modifying

**Stack Bounds**:

- Underflow: stackPointer < stack start
- Overflow: stackPointer >= stack end

### 6.3 Method Lookup

1. Start with receiver's class
2. Look in class's method dictionary for selector
3. If not found, try superclass
4. Repeat until method found or no superclass
5. If not found, send `doesNotUnderstand:` message
6. If that fails, throw MessageNotUnderstood exception

### 6.4 Block Closures

Blocks capture their home context for accessing outer variables:

```
Method context:
  | x |
  x := 42.
  [x + 1] value  ← Block accesses x from home context
```

**Block Creation**:

1. CREATE_BLOCK creates BlockContext
2. BlockContext.home points to creating context
3. Block's CompiledMethod stored in literals

**Block Execution**:

1. EXECUTE_BLOCK or Block>>value message
2. Create new context with:
   - sender = current context
   - home = block's home context
   - self = block's home context's self
3. Copy arguments to block's temporaries
4. Execute block's bytecode
5. Return to sender (not home)

**Variable Resolution**:

- If variable index < homeVarCount: search home chain
- Otherwise: use current context's temporaries

### 6.5 Return Semantics

**Method Return** (`^expression`):

- Returns to sender context
- Pushes result onto sender's stack

**Implicit Return**:

- At end of method bytecode
- Returns self if stack is empty
- Otherwise returns top of stack

**Block Return** (`^expression` in block):

- Returns from home method, not just block
- Unwinds contexts until home's sender reached
- Non-local return

**Block Value Return**:

- Normal block completion
- Returns to sender (block caller)
- Does not unwind past sender

---

## 7. Primitive Methods

Primitive methods are native implementations for performance-critical operations.

### 7.1 Primitive Numbers

```
Number  Method                  Description
------  ------                  -----------
1       SmallInteger>>+         Integer addition
2       SmallInteger>>-         Integer subtraction
3       SmallInteger>><         Integer less than
4       SmallInteger>>>         Integer greater than
5       SmallInteger>><=        Integer less than or equal
6       SmallInteger>>>=        Integer greater than or equal
7       SmallInteger>>=         Integer equality
8       SmallInteger>>~=        Integer inequality
9       SmallInteger>>*         Integer multiplication
10      SmallInteger>>/         Integer division
11      SmallInteger>>//        Integer modulo

60      Array>>at:              Array element access
61      Array>>at:put:          Array element store
62      Array>>size             Array size

63      String>>at:             String character access
64      String>>at:put:         String character store
65      String>>,               String concatenation
66      String>>size            String size

70      Object>>new             Create new instance
71      Object>>basicNew        Create new instance (no initialization)
72      Object>>basicNew:       Create new instance with size
75      Object>>identityHash    Object identity hash
111     Object>>class           Get object's class

154     True>>ifTrue:           Conditional execution (true case)
155     True>>ifFalse:          Conditional execution (false case)
156     True>>ifTrue:ifFalse:   Conditional execution (both cases)
157     False>>ifTrue:          Conditional execution (true case)
158     False>>ifFalse:         Conditional execution (false case)
159     False>>ifTrue:ifFalse:  Conditional execution (both cases)

201     Block>>value            Execute block with no arguments
202     Block>>value:           Execute block with one argument

1000    Exception>>mark         Exception handler marker (always fails)
1001    Exception>>signal       Signal (throw) exception
```

### 7.2 Primitive Semantics

#### 7.2.1 Integer Arithmetic Primitives (1, 2, 9, 10, 11)

**Arguments**: receiver (SmallInteger), argument (SmallInteger)

**Operation**:

1. Extract integer values from tagged values
2. Perform operation
3. Check for overflow
4. Return result as SmallInteger

**Failure Conditions**:

- Argument is not a SmallInteger
- Result overflows 31-bit range
- Division by zero (for `/` and `//`)

**Fallback**: Execute Smalltalk implementation if available

#### 7.2.2 Integer Comparison Primitives (3-8)

**Arguments**: receiver (SmallInteger), argument (SmallInteger)

**Operation**:

1. Extract integer values
2. Perform comparison
3. Return true or false

**Failure Conditions**:

- Argument is not a SmallInteger

**Fallback**: Execute Smalltalk implementation

#### 7.2.3 Array Access Primitives (60-62)

**Array>>at: (60)**:

- Arguments: receiver (Array), index (SmallInteger)
- Returns: element at index (1-based)
- Fails if: index out of bounds, receiver not Array

**Array>>at:put: (61)**:

- Arguments: receiver (Array), index (SmallInteger), value (any)
- Returns: value
- Fails if: index out of bounds, receiver not Array, array immutable

**Array>>size (62)**:

- Arguments: receiver (Array)
- Returns: number of elements as SmallInteger
- Fails if: receiver not Array

#### 7.2.4 String Primitives (63-66)

**String>>at: (63)**:

- Arguments: receiver (String), index (SmallInteger)
- Returns: character at index as SmallInteger (Unicode code point)
- Fails if: index out of bounds, receiver not String

**String>>, (65)**:

- Arguments: receiver (String), argument (String)
- Returns: new String with concatenated content
- Fails if: argument not String, memory allocation fails

**String>>size (66)**:

- Arguments: receiver (String)
- Returns: number of characters as SmallInteger
- Fails if: receiver not String

#### 7.2.5 Object Creation Primitives (70-72)

**Object>>new (70)**:

- Arguments: receiver (Class)
- Returns: new instance of class
- Operation:
  1. Allocate object with class's instance variable count
  2. Initialize all instance variables to nil
  3. Send `initialize` message to new instance
  4. Return instance

**Object>>basicNew (71)**:

- Same as `new` but skips `initialize` message

**Object>>basicNew: (72)**:

- Arguments: receiver (Class), size (SmallInteger)
- Returns: new instance with specified size
- Used for: Array, String, ByteArray creation

#### 7.2.6 Boolean Conditional Primitives (154-159)

**True>>ifTrue: (154)**:

- Arguments: receiver (true), block (Block)
- Returns: result of evaluating block
- Operation: Execute block and return result

**True>>ifFalse: (155)**:

- Arguments: receiver (true), block (Block)
- Returns: nil
- Operation: Do not execute block, return nil

**True>>ifTrue:ifFalse: (156)**:

- Arguments: receiver (true), trueBlock (Block), falseBlock (Block)
- Returns: result of evaluating trueBlock
- Operation: Execute trueBlock and return result

**False>>ifTrue: (157)**:

- Arguments: receiver (false), block (Block)
- Returns: nil
- Operation: Do not execute block, return nil

**False>>ifFalse: (158)**:

- Arguments: receiver (false), block (Block)
- Returns: result of evaluating block
- Operation: Execute block and return result

**False>>ifTrue:ifFalse: (159)**:

- Arguments: receiver (false), trueBlock (Block), falseBlock (Block)
- Returns: result of evaluating falseBlock
- Operation: Execute falseBlock and return result

**Notes**:

- These primitives enable efficient conditional execution
- Blocks are only evaluated when needed (short-circuit evaluation)
- Receiver must be exactly true or false (not truthy/falsy)

#### 7.2.7 Block Execution Primitives (201-202)

**Block>>value (201)**:

- Arguments: receiver (BlockContext)
- Returns: result of block execution
- Operation:
  1. Create execution context for block
  2. Set sender to current context
  3. Execute block bytecode
  4. Return result

**Block>>value: (202)**:

- Arguments: receiver (BlockContext), argument (any)
- Returns: result of block execution
- Operation:
  1. Create execution context for block
  2. Copy argument to first temporary variable
  3. Set sender to current context
  4. Execute block bytecode
  5. Return result

**Failure Conditions**:

- Receiver is not a BlockContext
- Wrong number of arguments for block

---

## 8. Exception Handling

### 8.1 Exception Model

Exceptions use context unwinding and handler blocks:

```smalltalk
[
    "protected code"
    10 / 0
] on: ZeroDivisionError do: [:ex |
    "handler code"
    'caught'
]
```

### 8.2 Exception Handling Mechanism

**Exception Signaling**:

1. Create exception object
2. Search context chain for handler
3. If found, unwind to handler context
4. Execute handler block with exception as argument
5. If not found, terminate with unhandled exception

**Handler Search**:

1. Start with current context
2. Check if context has exception handler marker (primitive 1000)
3. If yes, check if exception class matches handler class
4. If match, use this handler
5. Otherwise, continue to sender context
6. Repeat until handler found or no more contexts

**Context Unwinding**:

1. Execute `ensure:` blocks while unwinding
2. Clean up resources
3. Stop at handler context
4. Execute handler block

### 8.3 Exception Classes

```
Exception
  ├─ Error
  │   ├─ MessageNotUnderstood
  │   ├─ ZeroDivisionError
  │   ├─ IndexError
  │   ├─ ArgumentError
  │   └─ NameError
  ├─ Warning
  └─ Notification
```

### 8.4 Exception Handling Bytecode Pattern

```
Method with exception handler:
  PUSH_LITERAL <exception_class>
  PUSH_LITERAL <handler_block>
  SEND_MESSAGE #on:do: 2

The on:do: method:
  - Marks context with primitive 1000
  - Executes protected block
  - If exception occurs, executes handler block
```

---

## 9. Implementation Notes

### 9.1 Memory Alignment

- All objects must be 4-byte aligned (for tagged pointers)
- TaggedValue arrays must be 8-byte aligned
- Bytecode arrays can be byte-aligned

### 9.2 Garbage Collection

This specification does not mandate a specific GC algorithm, but implementations should:

- Support object marking (MARKED flag)
- Support generational GC (REMEMBERED flag)
- Support object pinning (PINNED flag)
- Handle forwarding pointers (FORWARDED flag)

### 9.3 Symbol Interning

Symbols must be interned (unique per string):

- Maintain global symbol table
- Symbol equality is pointer equality
- Symbol hash is based on string content

### 9.4 Method Lookup Caching

Implementations may cache method lookups for performance:

- Inline caching at call sites
- Global method cache
- Class-based method dictionaries

### 9.5 Primitive Failure Handling

When a primitive fails:

1. Check if method has Smalltalk bytecode
2. If yes, execute bytecode as fallback
3. If no, throw error

### 9.6 Stack Overflow Protection

Implementations should:

- Detect stack overflow before it occurs
- Throw StackOverflow exception
- Provide reasonable default stack size (minimum 1024 slots)

### 9.7 Instruction Pointer Management

- IP points to next instruction to execute
- After reading opcode, IP points to first operand byte
- After reading operands, IP points to next instruction
- JUMP instructions set IP directly

### 9.8 Context Switching

When switching contexts:

1. Save current context state (IP, SP)
2. Set activeContext to new context
3. Restore new context state
4. Continue execution

### 9.9 Bytecode Verification

Implementations may verify bytecode before execution:

- Check instruction bounds
- Verify operand validity
- Validate jump targets
- Check stack depth

### 9.10 Debugging Support

Implementations should support:

- Instruction-level stepping
- Breakpoints
- Stack traces
- Variable inspection
- Method source mapping

---

## 10. Corner Cases and Edge Conditions

### 10.1 Empty Stack Return

**Scenario**: Method ends with empty stack

**Behavior**: Return self

**Example**:

```smalltalk
method
    "no explicit return"

→ Returns self
```

### 10.2 Block Non-Local Return

**Scenario**: Block contains `^expression`

**Behavior**: Return from home method, not just block

**Example**:

```smalltalk
method
    [^42] value.
    99

→ Returns 42, not 99
```

### 10.3 Nested Block Variable Access

**Scenario**: Inner block accesses outer block's variable

**Behavior**: Search home chain for variable

**Example**:

```smalltalk
method
    | x |
    x := 1.
    [| y |
        y := 2.
        [x + y] value
    ] value

→ Inner block accesses both x and y through home chain
```

### 10.4 Block Outliving Home Context

**Scenario**: Block returned from method, then executed

**Behavior**: Block still accesses home context variables

**Example**:

```smalltalk
makeCounter
    | count |
    count := 0.
    ^[count := count + 1]

counter := self makeCounter.
counter value.  → 1
counter value.  → 2
```

**Implementation**: Home context must not be garbage collected while block exists

### 10.5 Primitive Failure with No Fallback

**Scenario**: Primitive fails and method has no bytecode

**Behavior**: Throw error

**Example**:

```smalltalk
SmallInteger>>+ arg
    <primitive: 1>
    "no fallback code"

1 + 'string'  → Error: primitive failed, no fallback
```

### 10.6 Message Send to nil

**Scenario**: Send message to nil

**Behavior**: Normal message send (nil is an object)

**Example**:

```smalltalk
nil isNil  → true
nil class  → UndefinedObject
```

### 10.7 Division by Zero

**Scenario**: Integer division by zero

**Behavior**: Primitive fails, fallback code signals exception

**Example**:

```smalltalk
10 / 0  → ZeroDivisionError
```

### 10.8 Array Index Out of Bounds

**Scenario**: Access array element beyond size

**Behavior**: Primitive fails, fallback code signals exception

**Example**:

```smalltalk
#(1 2 3) at: 5  → IndexError
```

### 10.9 Method Not Found

**Scenario**: Send message with no matching method

**Behavior**: Send `doesNotUnderstand:` message

**Example**:

```smalltalk
42 unknownMessage  → MessageNotUnderstood
```

### 10.10 Stack Overflow

**Scenario**: Infinite recursion

**Behavior**: Detect stack overflow, throw exception

**Example**:

```smalltalk
infinite
    self infinite

self infinite  → StackOverflow
```

### 10.11 Context Unwinding with ensure:

**Scenario**: Exception thrown with ensure: blocks

**Behavior**: Execute all ensure: blocks while unwinding

**Example**:

```smalltalk
[
    [
        error signal
    ] ensure: [
        'cleanup 1' print
    ]
] ensure: [
    'cleanup 2' print
]

→ Prints 'cleanup 1' then 'cleanup 2' before exception propagates
```

### 10.12 Boolean Conditional with Non-Boolean

**Scenario**: ifTrue: sent to non-boolean

**Behavior**: Method not found (only True and False implement ifTrue:)

**Example**:

```smalltalk
42 ifTrue: [99]  → MessageNotUnderstood
```

### 10.13 Block Argument Count Mismatch

**Scenario**: Block called with wrong number of arguments

**Behavior**: Error

**Example**:

```smalltalk
[:x | x + 1] value: 1 value: 2  → ArgumentError
```

### 10.14 Immutable Object Modification

**Scenario**: Attempt to modify immutable object

**Behavior**: Primitive fails, throw error

**Example**:

```smalltalk
#(1 2 3) at: 1 put: 99  → Error: array is immutable
```

### 10.15 Symbol Interning

**Scenario**: Create multiple symbols with same string

**Behavior**: Return same symbol object

**Example**:

```smalltalk
#abc == #abc  → true (pointer equality)
'abc' asSymbol == 'abc' asSymbol  → true
```

---

## 11. Compliance and Testing

### 11.1 Compliance Requirements

A compliant VM implementation must:

1. **Support all bytecode instructions** as specified
2. **Implement tagged value encoding** correctly
3. **Handle all primitive methods** or provide fallbacks
4. **Support exception handling** mechanism
5. **Implement proper context management**
6. **Handle all corner cases** as specified

### 11.2 Test Suite

Implementations should pass:

- **Expression tests**: 64 test cases covering all language features
- **Bytecode tests**: Verify each instruction works correctly
- **Primitive tests**: Test all primitive methods
- **Exception tests**: Test exception handling and unwinding
- **Block tests**: Test closures and non-local returns
- **Image tests**: Test image save/load functionality

### 11.3 Interoperability

Compliant VMs should:

- **Load images** created by other compliant VMs
- **Execute bytecode** identically to other VMs
- **Produce identical results** for same inputs

---

## 12. Version History

**Version 1.0** (2025-01-05):

- Initial specification
- Complete bytecode instruction set
- Tagged value representation
- Object memory format
- Image file format
- Primitive method definitions
- Exception handling mechanism
- Corner case documentation

---

## 13. References

- **Smalltalk-80: The Language and its Implementation** (Blue Book)
- **Smalltalk-80: The Interactive Programming Environment** (Orange Book)
- **Efficient Implementation of the Smalltalk-80 System** (Deutsch & Schiffman)
- **Current implementation**: `src/cpp/` directory

---

## Appendix A: Quick Reference

### Bytecode Opcodes

```
0=PUSH_LITERAL  1=PUSH_IVAR     2=PUSH_TEMP     3=PUSH_SELF
4=STORE_IVAR    5=STORE_TEMP    6=SEND_MESSAGE  7=RETURN
8=JUMP          9=JUMP_IF_TRUE  10=JUMP_IF_FALSE 11=POP
12=DUPLICATE    13=CREATE_BLOCK 14=EXECUTE_BLOCK
```

### Tagged Value Tags

```
00=Pointer  01=Special  10=Float  11=Integer
```

### Object Types

```
0=IMMEDIATE  1=OBJECT  2=ARRAY  3=BYTE_ARRAY
4=SYMBOL     5=CONTEXT 6=CLASS  7=METHOD
```

### Common Primitive Numbers

```
1-11=Integer ops  60-62=Array  63-66=String
70-72=Object      154-159=Boolean  201-202=Block
```

---

**End of Specification**
