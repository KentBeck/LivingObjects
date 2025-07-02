# SmalltalkLSP Testing Guide

This document describes how to run tests for the C++ Smalltalk VM implementation.

## Quick Start

To run all tests:

```bash
make test-all
```

Or run the test script directly:

```bash
./run_all_tests.sh
```

For verbose output showing test details:

```bash
make test-verbose
```

To run expression demos:

```bash
make expressions
```

To run custom expressions:

```bash
./run_expressions.sh "3 + 4" "5 * 6"
```

## Test Categories

### 1. Basic VM Expressions (18 tests)
Tests fundamental arithmetic and comparison operations with live execution:
- Addition, subtraction, multiplication, division
- Comparisons: `<`, `>`, `=`, `~=`, `<=`, `>=`
- Boolean values: `true`, `false`, `nil`
- Nested expressions with parentheses
- Complex arithmetic with precedence

**Example**: `./build/smalltalk-vm "3 + 4"` → `Result: Integer(7)`

### 2. Block Parsing Tests (2 tests)
Tests the parsing and compilation of block syntax `[expression]`:
- Block parsing into AST nodes
- Block compilation to bytecode

**Example**: `[3 + 4]` parses to `BlockNode` and compiles to `CREATE_BLOCK` bytecode

### 3. Block Runtime Tests (2 tests)
Tests the runtime execution of block primitives:
- Block primitive execution (`Block>>value`)
- Block memory layout and context management

**Example**: Creating `BlockContext` objects and executing the `value` primitive

### 4. Block Pipeline Tests (2 tests)
Tests the complete block pipeline integration:
- Block parsing in the main VM
- Block bytecode generation (CREATE_BLOCK instruction)

### 5. Error Handling Tests (3 tests)
Tests proper error handling for invalid input:
- Invalid syntax
- Malformed blocks
- Empty input

### 6. Expression Variety Tests (4 tests)
Tests various complex expressions:
- Large numbers
- Multiple operations with precedence
- Parentheses priority
- Complex comparisons

## Test Implementation

The test suite uses a simple bash script that doesn't require external frameworks like gtest. Tests are implemented as:

1. **VM Integration Tests**: Run expressions through the main VM and show live results
2. **Unit Tests**: Build and run specific test executables for components
3. **Error Tests**: Verify that invalid input produces expected errors
4. **Expression Demos**: Interactive expression running with immediate feedback

## Test Files

- `run_all_tests.sh` - Main test runner script
- `tests/minimal_block_test.cpp` - Block primitive execution tests
- `tests/minimal_parse_test.cpp` - Block parsing tests
- `tests/simple_block_test.cpp` - Block compilation tests
- `tests/debug_block_layout.cpp` - Memory layout validation
- `run_expressions.sh` - Interactive expression runner

## Adding New Tests

To add a new test:

1. **VM Integration Test**: Add to `run_all_tests.sh`:
   ```bash
   run_test "Test Name" "./build/smalltalk-vm 'expression' | grep -q 'expected_output'"
   ```

2. **Unit Test**: Create new `.cpp` file in `tests/` directory and add build/run logic to the script.

## Current Status

✅ **All 31 tests passing**

The block implementation is fully functional:
- Parsing: `[expression]` syntax works
- Compilation: Generates correct `CREATE_BLOCK` bytecode
- Runtime: `BlockContext` objects created properly
- Primitives: `Block>>value` execution works
- Error handling: Invalid input properly rejected

## Architecture Tested

The tests validate the complete pipeline:

```
Source Code → Parser → AST → Compiler → Bytecode → VM → Result
     [3+4]  →   ✓   →  ✓  →    ✓    →    ✓    → ✓ →   7
```

Block-specific pipeline:
```
[3+4] → BlockNode → CREATE_BLOCK → BlockContext → value → New Context
```

## Running Individual Components

- **Main VM**: `./build/smalltalk-vm "expression"`
- **Expression Runner**: `./run_expressions.sh` or `make expressions`
- **Custom Expressions**: `./run_expressions.sh "your expression"`
- **REPL**: `./build/smalltalk-repl`
- **Image Tool**: `./build/image-tool <command>`
- **Block Tests**: `./build/minimal_block_test`
- **Parse Tests**: `./build/minimal_parse_test`

All tests run quickly (< 5 seconds total) and provide clear pass/fail feedback with colored output.

## Expression Examples

The system supports a wide variety of expressions:

```bash
# Arithmetic
./run_expressions.sh "3 + 4" "10 * 5" "(2 + 3) * 4"

# Comparisons  
./run_expressions.sh "5 > 3" "2 = 2" "4 ~= 7"

# Complex expressions
./run_expressions.sh "((10 + 5) * 2) / 6" "(3 + 4) > (2 * 3)"

# Blocks (parsing demo)
./build/smalltalk-vm "[3 + 4]"  # Shows parsing and compilation
```