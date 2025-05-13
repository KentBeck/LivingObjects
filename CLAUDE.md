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

## Known Issues

- There are circular dependencies between VM and compiler packages that need resolution
- The ByteArray class is being implemented to avoid import cycles