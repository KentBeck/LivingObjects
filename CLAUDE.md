# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SmalltalkLSP is a Smalltalk implementation with Language Server Protocol (LSP) support. The project consists of:

1. A Go-based bytecode interpreter (in `src/interpreter/`)
2. Smalltalk class implementations (in `src/*.st`)
3. A JavaScript interpreter (in `src/js-interpreter/`)
4. LSP implementations for editor integration

The project is gradually migrating from a Go-based interpreter to a self-hosted Smalltalk system, as outlined in `smalltalk-implementation-plan.md`.

## Building and Running

### Go Interpreter

```bash
# Navigate to the interpreter directory
cd src/interpreter

# Build the interpreter
make build
# OR
go build -o smalltalk-vm

# Run the interpreter
make run
# OR
./smalltalk-vm

# Run the factorial demo
./smalltalk-vm demo

# Clean build artifacts
make clean
```

### JavaScript Interpreter

```bash
# Navigate to the JS interpreter directory
cd src/js-interpreter

# Install dependencies
npm install

# Run tests
npm test
```

## Running Tests

### Go Tests

```bash
# Run all Go tests
cd src/interpreter
go test ./...

# Run tests for specific packages
go test ./core ./vm ./classes
```

### JavaScript Tests

```bash
# Run all JS tests
cd src/js-interpreter
npm test
```

## Project Architecture

### Go Interpreter Structure

The Go interpreter is organized into packages:

- **core**: Base Object type, memory management, immediate values
- **classes**: Smalltalk class implementations (Array, Dictionary, String, Symbol, Block, etc.)
- **vm**: VM implementation with context, bytecode handlers, and primitives
- **compiler**: Method builder and parsing subsystem
- **bytecode**: Bytecode constants and definitions

### Smalltalk Implementation

The Smalltalk implementation consists of `.st` files in the `src/` directory:

- Core classes (Object, Collection, Dictionary)
- LSP implementation (SmalltalkLSP, LSPMessageHandler)
- Runtime components (SmalltalkContext, SmalltalkMethod)
- Testing framework (TestRunner)

### Implementation Plan

The project follows a phased migration approach:

1. Move runtime data structures to Smalltalk objects
2. Implement directory loading and SUnit tests
3. Move the compiler to Smalltalk
4. Move the parser to Smalltalk
5. Implement Smalltalk source code management in Smalltalk

## Key Files

- **src/interpreter/main.go**: Entry point for the Go interpreter
- **src/Main.st**: Main Smalltalk entry point
- **src/SmalltalkLSP.st**: LSP implementation
- **src/SmalltalkMethod.st**: Method representations
- **src/SmalltalkClass.st**: Class representations
- **src/OrderedCollection.st**: Collection implementation
- **src/BytecodeInterpreter.st**: Bytecode interpreter

## C++ VM Implementation

The project also includes a C++ VM implementation in `src/cpp/` with the following structure:

- **src/**: Core VM implementation (interpreter.cpp, memory_manager.cpp, etc.)
- **include/**: Header files for VM components
- **tests/**: Test suite including comprehensive expression tests
- **plan.md**: Detailed 5-phase refactoring plan for VM improvements

### C++ VM Building and Testing

```bash
# Navigate to the C++ implementation
cd src/cpp

# Build the VM
make

# Run all tests including expression tests
make test

# Run specific expression tests
./tests/all_expressions_test

# Clean build artifacts
make clean
```

## Expert Role Instructions

When working on this codebase, Claude should act as an expert in:

### C++ and Virtual Machine Implementation
- Apply deep knowledge of C++ best practices, memory management, and performance optimization
- Understand VM internals including bytecode interpretation, garbage collection, and object models
- Recognize common VM patterns like tagged values, context switching, and method dispatch
- Identify performance bottlenecks and architectural issues in VM design
- Suggest improvements based on modern VM implementation techniques

### Smalltalk Language Semantics
- Understand Smalltalk object model, message passing, and block closures
- Apply knowledge of Smalltalk-80 semantics and proper language behavior
- Recognize fake implementations and suggest real functionality
- Ensure implementations match expected Smalltalk behavior and conventions

### Progress Tracking

When working on VM improvements, always:

1. **Run Expression Tests**: Use `src/cpp/tests/all_expressions_test.cpp` as the primary progress indicator
   - Run tests before and after changes to measure improvement
   - Report which test categories are passing/failing
   - Use test results to identify next implementation priorities

2. **Track Implementation Status**: 
   - Monitor the 7 identified fake functions and their implementation status
   - Update todo lists with specific technical details
   - Reference the 5-phase refactoring plan in `src/cpp/plan.md` for long-term goals

3. **Performance Analysis**:
   - Consider performance implications of architectural decisions
   - Suggest optimizations based on VM best practices
   - Identify opportunities for the high-performance raw memory stack approach

4. **Code Quality**:
   - Ensure proper error handling and bounds checking
   - Maintain consistent object model throughout the VM
   - Follow C++ best practices for memory safety and performance

## Expression Test Categories

The test suite in `all_expressions_test.cpp` covers:
- **arithmetic**: Basic mathematical operations
- **comparison**: Relational operators and equality
- **object_creation**: Object instantiation
- **strings**: String literals and operations
- **literals**: Boolean and nil values
- **variables**: Variable assignment and access
- **blocks**: Block creation and execution
- **conditionals**: Control flow (not yet implemented)
- **collections**: Array and collection operations (not yet implemented)
- **dictionaries**: Dictionary operations (not yet implemented)
- **class_creation**: Dynamic class creation (not yet implemented)

Use these categories to measure implementation progress and prioritize work.

## Known Issues

- There are circular dependencies between VM and compiler packages that need resolution
- The ByteArray class is being implemented to avoid import cycles
- Mixed Object* and TaggedValue representations cause architectural inconsistencies
- Instance variable access has Object*/TaggedValue conversion issues
- Block execution infrastructure needs proper bytecode handling