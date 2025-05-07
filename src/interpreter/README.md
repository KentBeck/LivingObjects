# SmalltalkLSP Bytecode Interpreter

This is a bytecode interpreter for the SmalltalkLSP project, implemented in Go.

## Features

- Minimal bytecode set for Smalltalk semantics
- Stop & copy garbage collection
- Image-based persistence
- Support for 2^32 temporary and instance variables
- Modular architecture with clear separation of concerns

## Directory Structure

The interpreter is organized into the following packages:

- **core**: Core functionality including base Object type, memory management, and immediate values
- **classes**: Smalltalk class implementations (Array, Dictionary, String, Symbol, Block, etc.)
- **vm**: Virtual Machine implementation including context, bytecode handlers, and primitives
- **compiler**: Compiler subsystem including method builder
- **utils**: Utility functions including type conversions and image loading/saving

## Bytecode Set

The interpreter uses a minimal bytecode set:

1. **PUSH_LITERAL** (0): Push a literal from the literals array (followed by 4-byte index)
2. **PUSH_INSTANCE_VARIABLE** (1): Push an instance variable value (followed by 4-byte offset)
3. **PUSH_TEMPORARY_VARIABLE** (2): Push a temporary variable value (followed by 4-byte offset)
4. **PUSH_SELF** (3): Push self onto the stack
5. **STORE_INSTANCE_VARIABLE** (4): Store a value into an instance variable (followed by 4-byte offset)
6. **STORE_TEMPORARY_VARIABLE** (5): Store a value into a temporary variable (followed by 4-byte offset)
7. **SEND_MESSAGE** (6): Send a message (followed by 4-byte selector index and 4-byte arg count)
8. **RETURN_STACK_TOP** (7): Return the value on top of the stack
9. **JUMP** (8): Jump to a different bytecode (followed by 4-byte target)
10. **JUMP_IF_TRUE** (9): Jump if top of stack is true (followed by 4-byte target)
11. **JUMP_IF_FALSE** (10): Jump if top of stack is false (followed by 4-byte target)
12. **POP** (11): Pop the top value from the stack
13. **DUPLICATE** (12): Duplicate the top value on the stack

## Building and Running

```bash
# Build the interpreter
go build -o smalltalk-vm

# Run with an image file
./smalltalk-vm image.st

# Run the factorial demo
./smalltalk-vm demo
```

## Object Memory and Garbage Collection

The interpreter uses a stop & copy garbage collection algorithm:

1. Memory is divided into two equal-sized semi-spaces: from-space and to-space
2. Objects are allocated in from-space during normal execution
3. When memory is low, live objects are copied from from-space to to-space
4. References are updated to point to the new locations
5. Spaces are swapped after collection completes

## Image Format

The image format consists of:

1. A header with magic number, version, and object counts
2. A serialized object graph
3. References to global variables and the root object

## Next Steps

- Implement a parser for Smalltalk source code
- Add support for closures and blocks
- Implement a compiler from Smalltalk to bytecode
- Add debugging support
- Integrate with the LSP server
