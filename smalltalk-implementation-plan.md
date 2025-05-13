# Smalltalk Implementation Plan

## Goal

Move from a Go-based Smalltalk interpreter to a self-hosted Smalltalk system by progressively migrating components from Go to Smalltalk.

## Phase 0: Move Runtime Data Structures to Smalltalk Objects

### Step 1: Core Object System Completion

- **Implement ByteArray class**

  - Create ByteArray class in Smalltalk
  - Implement primitives for byte access
  - Ensure ByteArray can store compiled method bytecodes

- **Complete Collection classes**

  - Implement OrderedCollection
  - Implement Dictionary with proper hashing
  - Implement Set
  - Ensure all collections use Smalltalk objects internally

- **Implement Stream classes**
  - ReadStream, WriteStream, ReadWriteStream
  - FileStream for file I/O
  - These will be essential for file loading

### Step 2: Runtime Data Structure Migration

- **Method objects**

  - Store bytecodes in ByteArray instead of Go slice
  - Store literals in Smalltalk Array
  - Store temporary variable names in Smalltalk Array

- **Context objects**

  - Store stack in Smalltalk Array
  - Store temporary variables in Smalltalk Array
  - Ensure context chain uses Smalltalk objects

- **Class objects**
  - Store method dictionary as Smalltalk Dictionary
  - Store instance variables as Smalltalk Array
  - Store superclass as Smalltalk reference

### Step 3: VM State Migration

- **Replace VM globals with Smalltalk Dictionary**

  - Move from Go map to Smalltalk Dictionary
  - Ensure all global lookups use Smalltalk objects

- **Replace immediate values with proper Smalltalk objects**

  - Ensure nil, true, false are Smalltalk objects
  - Ensure integers and floats are properly represented

- **Migrate VM execution state**
  - Move executor state to Smalltalk objects
  - Ensure all VM operations use Smalltalk objects

## Phase 1: Directory Loading and SUnit Tests

### Step 1: File System Access

- **Implement file system primitives**

  - Directory listing primitive
  - File reading primitive
  - File writing primitive

- **Create FileSystem class in Smalltalk**
  - Methods for directory operations
  - Methods for file operations
  - Error handling for file operations

### Step 2: Smalltalk Code Loader

- **Create SmalltalkLoader class**

  - Load all .st files from a directory
  - Parse and compile each file
  - Register classes and methods

- **Implement dependency resolution**
  - Handle class dependencies
  - Process files in correct order
  - Track and report errors

### Step 3: SUnit Test Framework

- **Implement TestCase class**

  - setUp and tearDown methods
  - assert methods
  - test running infrastructure

- **Implement TestSuite class**

  - Collect related tests
  - Run multiple tests
  - Report results

- **Create TestRunner**
  - Discover tests
  - Run tests
  - Report results

## Phase 2: Move Compiler to Smalltalk

### Step 1: AST Classes

- **Define base AST node classes**

  - ASTNode as abstract base class
  - Node types for all language constructs
  - Visitor pattern for traversal

- **Implement AST construction**
  - Methods to build AST nodes
  - Helper methods for common patterns
  - Validation methods

### Step 2: Bytecode Generation

- **Create BytecodeGenerator class**

  - Visit AST nodes
  - Generate appropriate bytecodes
  - Manage literals and variables

- **Implement method compilation**

  - Compile method from AST
  - Set up method header
  - Link to class

- **Add optimization passes**
  - Basic optimizations
  - Inline simple methods
  - Constant folding

## Phase 3: Move Parser to Smalltalk

### Step 1: Scanner Implementation

- **Create Scanner class**

  - Token recognition
  - Handle comments and whitespace
  - Position tracking

- **Implement token classification**
  - Identifiers, keywords, operators
  - Numbers, strings, symbols
  - Special characters

### Step 2: Parser Implementation

- **Create Parser class**

  - Use Scanner for tokens
  - Implement recursive descent parsing
  - Build AST

- **Implement parsing rules**

  - Method definitions
  - Expressions and statements
  - Message sends with precedence

- **Add error handling**
  - Detect syntax errors
  - Recover from errors
  - Report helpful messages

## Phase 4: Smalltalk Source in Smalltalk

### Step 1: Source Code Management

- **Create SourceCode class**

  - Store source for methods
  - Link source to compiled methods
  - Support source retrieval

- **Implement change tracking**
  - Track changes to classes and methods
  - Support undoing changes
  - Group related changes

### Step 2: Persistence

- **Implement image save/load**

  - Save object memory to disk
  - Load object memory from disk
  - Initialize system after loading

- **Add change logging**
  - Record all changes
  - Support replaying changes
  - Enable recovery from crashes

## Implementation Strategy

1. **Incremental migration**

   - Keep Go implementation working
   - Replace one component at a time
   - Test thoroughly at each step

2. **Bootstrap approach**

   - Use Go implementation to compile Smalltalk code
   - Use Smalltalk code to replace Go implementation
   - Eventually eliminate Go dependency

3. **Test-driven development**
   - Write tests for each component
   - Run tests after each change
   - Ensure compatibility with existing code
