# C++ Bytecode Interpreter Implementation Plan

## Overview
Port the Go Smalltalk bytecode interpreter to C++ while maintaining full test compatibility. The C++ version will use tagged pointers for performance and prepare for eventual JIT compilation.

## Test Analysis Summary
**130+ test functions** across these categories:
- **VM Core**: Bytecode handlers, context management, message sending (47 tests)
- **Object System**: Classes, methods, blocks, arrays, strings, symbols (62 tests)  
- **Compiler**: Parsing, bytecode generation, method building (21 tests)
- **Primitives**: Integer/float math, boolean operations, memory access
- **Integration**: Factorial computation, expression evaluation

## Core Components to Port

### 1. Object System (`pile` package → `objects/`)
- **Object**: Base object with tagged pointer support
- **Class**: Class definitions with method dictionary
- **Immediate Values**: Tagged integers, floats, booleans, nil
- **Collections**: Array, Dictionary, String, Symbol, ByteArray
- **Methods & Blocks**: Executable code objects

### 2. Virtual Machine (`vm` package → `vm/`)
- **Context**: Execution context with stack management
- **Bytecode Handlers**: 13 core bytecode operations
- **Message Dispatch**: Method lookup and invocation
- **Primitives**: Built-in operations (+, -, *, <, etc.)

### 3. Memory Management (`pile/memory.go` → `memory/`)
- **Stop & Copy GC**: Two semi-space garbage collector
- **Tagged Pointers**: Handle immediate values correctly
- **Object Allocation**: Proper alignment and initialization

### 4. Compiler (`compiler` package → `compiler/`)
- **Parser**: Smalltalk source to AST
- **Bytecode Generator**: AST to bytecode compilation
- **Method Builder**: Method object construction

## Implementation Phases

### Phase 1: Core Object System (Week 1-2)
**Goals**: Basic object model with tagged pointers working

**Tasks**:
1. **Tagged Pointer System**
   - Implement `TaggedValue` union type
   - Support for immediate integers, floats, booleans, nil
   - Pointer tagging with 2-bit tags (00=pointer, 01=special, 10=float, 11=int)

2. **Base Object Types**
   - `Object` base class with class pointer
   - `Class` with method dictionary and instance variables
   - Memory layout compatible with GC

3. **Test Infrastructure**
   - C++ test framework (Google Test)
   - Port immediate value tests (14 tests)
   - Port basic object tests (12 tests)

**Deliverables**:
- `objects/tagged_value.h/cpp`
- `objects/object.h/cpp` 
- `objects/class.h/cpp`
- Working immediate value tests

### Phase 2: Collections & Primitives (Week 3)
**Goals**: Core collection types and arithmetic working

**Tasks**:
1. **Collection Classes**
   - Array with bounds checking
   - Dictionary with hash table
   - String with UTF-8 support
   - Symbol table for interned strings

2. **Primitive Operations**
   - Integer arithmetic (+, -, *, /, <, >, =)
   - Float operations with precision handling
   - Boolean logic (and, or, not)
   - Memory access primitives

3. **Port Collection Tests** (35 tests)
   - Array operations and bounds checking
   - Dictionary get/set/remove operations
   - String manipulation and comparison
   - Symbol interning and equality

**Deliverables**:
- `objects/array.h/cpp`, `objects/dictionary.h/cpp`
- `objects/string.h/cpp`, `objects/symbol.h/cpp`
- `vm/primitives.h/cpp`
- Passing collection and primitive tests

### Phase 3: Virtual Machine Core (Week 4-5)
**Goals**: Bytecode execution and context management

**Tasks**:
1. **Execution Context**
   - Stack management with overflow protection  
   - Temporary variable storage
   - Method context chaining

2. **Bytecode Handlers**
   - 13 core bytecode operations
   - Efficient dispatch mechanism (computed goto or jump table)
   - Stack manipulation (push, pop, duplicate)

3. **Message Dispatch**
   - Method lookup with inheritance
   - Argument passing and stack management
   - Return value handling

4. **Port VM Tests** (47 tests)
   - Context operations
   - Bytecode handler execution
   - Message sending scenarios
   - Stack management edge cases

**Deliverables**:
- `vm/context.h/cpp`
- `vm/bytecode_handlers.h/cpp`
- `vm/message_dispatch.h/cpp`
- `vm/vm.h/cpp`
- Passing VM core tests

### Phase 4: Memory Management (Week 6)
**Goals**: Stop & copy garbage collection working with tagged pointers

**Tasks**:
1. **Two-Space Collector**
   - Semi-space allocation and copying
   - Root set scanning (globals, stacks)
   - Object traversal with proper tagged pointer handling

2. **Tagged Pointer GC Integration**
   - Skip scanning immediate values
   - Handle forwarding pointers correctly
   - Preserve tagged values during collection

3. **Memory Layout**
   - Proper object alignment
   - Header word optimization
   - Forwarding pointer mechanics

4. **Port Memory Tests** (8 tests)
   - Allocation and collection cycles
   - Reference updating during GC
   - Memory pressure scenarios

**Deliverables**:
- `memory/gc.h/cpp`
- `memory/allocator.h/cpp`
- Passing GC and memory tests

### Phase 5: Blocks & Control Flow (Week 7)
**Goals**: Block objects and non-local returns

**Tasks**:
1. **Block Objects**
   - Closures with captured variables
   - Block compilation and execution
   - Outer context references

2. **Control Flow**
   - Jump instructions (conditional/unconditional)
   - Non-local returns from blocks
   - Exception handling framework

3. **Port Block Tests** (18 tests)
   - Block creation and execution
   - Variable capture scenarios
   - Nested block handling
   - Non-local return edge cases

**Deliverables**:
- `objects/block.h/cpp`
- `vm/control_flow.h/cpp`
- Passing block and control flow tests

### Phase 6: Compiler Integration (Week 8)
**Goals**: Parse Smalltalk source and generate bytecode

**Tasks**:
1. **Parser**
   - Smalltalk syntax parsing
   - AST construction
   - Error handling and recovery

2. **Code Generation**
   - AST to bytecode translation
   - Literal table construction
   - Method object assembly

3. **Port Compiler Tests** (21 tests)
   - Expression parsing
   - Method compilation
   - Bytecode generation accuracy

**Deliverables**:
- `compiler/parser.h/cpp`
- `compiler/codegen.h/cpp`
- Passing compiler tests

### Phase 7: Integration & Optimization (Week 9-10)
**Goals**: Full system integration and performance

**Tasks**:
1. **Integration Testing**
   - Port remaining integration tests (8 tests)
   - Factorial computation verification
   - End-to-end expression evaluation

2. **Performance Optimization**
   - Bytecode dispatch optimization (threaded code)
   - Inline caching for method dispatch
   - Memory allocation tuning

3. **Build System**
   - CMake configuration
   - Cross-platform compatibility
   - Test automation

**Deliverables**:
- Complete working C++ interpreter
- All 130+ tests passing
- Performance benchmarks vs Go version

## Build & Test Strategy

### Build System
```cmake
# CMakeLists.txt structure
project(SmalltalkCppVM)
add_subdirectory(objects)
add_subdirectory(vm) 
add_subdirectory(memory)
add_subdirectory(compiler)
add_subdirectory(tests)
```

### Test Porting
- **1:1 test mapping**: Each Go test gets C++ equivalent
- **Test data preservation**: Same test cases, expected results
- **Continuous validation**: Tests pass after each phase

### Performance Targets
- **Startup**: ≤ 10ms (vs Go's ~50ms)
- **Bytecode dispatch**: ≥ 100M ops/sec
- **GC pause**: ≤ 1ms for small heaps
- **Memory overhead**: ≤ 2x object size

## Risk Mitigation

### Technical Risks
1. **Tagged pointer GC bugs**: Extensive testing, reference Go implementation
2. **Memory corruption**: Valgrind, AddressSanitizer integration  
3. **Performance regression**: Continuous benchmarking vs Go

### Schedule Risks
1. **Complex GC debugging**: Allocate extra time for Phase 4
2. **Block semantics**: Non-local returns are subtle, plan for iteration

## Success Criteria
- ✅ All 130+ Go tests pass in C++
- ✅ Factorial computation produces identical results
- ✅ Memory management stable under stress
- ✅ Performance equal or better than Go version
- ✅ Clean valgrind/sanitizer runs
- ✅ Ready for JIT compiler integration

## Next Steps After Completion
1. **JIT Compiler**: Template-based code generation
2. **Inline Caching**: Polymorphic method dispatch optimization  
3. **Generational GC**: Reduce collection overhead
4. **LSP Integration**: Connect to language server
