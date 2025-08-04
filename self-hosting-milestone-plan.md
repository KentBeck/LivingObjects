# Self-Hosting Smalltalk Milestone Plan (Simple & Test-Driven)

## Goal

Build a Smalltalk system that can read Smalltalk source code, compile it to bytecode, and execute it - essentially creating a self-hosting Smalltalk implementation that can rebuild itself from source.

**Key Principles**:

1. **Test-driven development** - Every change is validated by tests
2. **Simplicity first** - Choose the simplest solution that works
3. **Correctness over speed** - Get it working right, ignore performance

## Current State Analysis

### What We Have ✅

- **C++ VM Infrastructure**: Working bytecode interpreter, memory management, primitive operations
- **Basic Smalltalk Classes**: SmalltalkCompiler, SmalltalkClass, SmalltalkMethod, SmalltalkObject
- **Core Collections**: OrderedCollection, Dictionary, ByteArray with just-enough implementations
- **Expression Evaluation**: 93.75% of expressions working (60/64 tests passing)
- **Basic Parser**: Simple tokenizer and recursive descent parser in Smalltalk
- **Bytecode Generation**: Method compilation to bytecode format

### What We Need 🎯

- **SUnit Prerequisites**: ZEROTH PRIORITY - test everything in C++ necessary for writing a testing framework
- **SUnit Testing Framework**: FIRST PRIORITY - safety net for all changes
- **Complete Parser & Compiler**: Full Smalltalk syntax support in Smalltalk
- **Core Class Library**: String, Array, Integer, Boolean, Block classes with all essential methods
- **Image Building**: Ability to serialize/deserialize the object graph
- **Bootstrap Process**: Mechanism to build the system from source files

---

## Phase 1: SUnit Testing Framework (Week 1)

### 1.1 Minimal SUnit Implementation

**Goal**: Get basic testing infrastructure working immediately

**Tasks**:

- Implement minimal TestCase class with core assertion methods
- Add basic `assert:`, `deny:`, `assert:equals:` methods
- Create simple TestSuite class for organizing tests
- Implement TestResult class for tracking pass/fail
- Add basic TestRunner that can execute tests and report results
- Create `setUp` and `tearDown` methods

**Acceptance Criteria**:

- Can write simple test cases that pass/fail correctly
- Test runner shows clear pass/fail results
- Basic assertions work reliably
- Can organize tests into suites

**Simplicity Focus**: Start with the absolute minimum needed to write and run tests. Don't worry about fancy features or reporting.

### 1.2 Test Infrastructure for Current Code

**Goal**: Create comprehensive test coverage for existing functionality

**Tasks**:

- Write tests for OrderedCollection (all 25+ methods)
- Write tests for Dictionary (all 20+ methods)
- Write tests for ByteArray (all 15+ methods)
- Write tests for current SmalltalkCompiler functionality
- Write tests for current expression evaluation (document the 60/64 passing tests)
- Create test categories matching our development phases

**Acceptance Criteria**:

- All existing functionality has test coverage
- Tests pass consistently (establish baseline)
- Test failures clearly indicate what broke
- Tests serve as documentation of expected behavior

**Simplicity Focus**: Write straightforward tests that verify basic functionality. Don't over-engineer the test structure.

---

## Phase 2: Core Infrastructure with TDD (Weeks 2-3)

### 2.1 Fix Current Failing Tests (Test-Driven)

**Goal**: Get to 100% expression test pass rate using TDD

**Tasks**:

- **Write tests first** for the 4 currently failing expression cases
- Fix symbol literal parsing (`#abc` should return Symbol, not Object)
- Fix nested block variable scoping issue
- Implement Boolean conditional methods (`ifTrue:`, `ifFalse:`, `ifTrue:ifFalse:`)
- **Run full test suite after each fix**

**Acceptance Criteria**:

- All 64 expression tests pass
- No regressions in existing functionality
- New functionality has comprehensive test coverage

**Simplicity Focus**: Fix one issue at a time. Choose the simplest implementation that makes the tests pass.

### 2.2 Complete String Class Implementation (Test-Driven)

**Goal**: Full String class with comprehensive test coverage

**Tasks**:

- **Write comprehensive String tests first** (focus on essential methods)
- Implement String class in Smalltalk with simple instance variables
- Add core methods: `size`, `at:`, `at:put:`, `copyFrom:to:`, `=`, `hash`
- Add string operations: `,` (concatenation), `asUppercase`, `asLowercase`
- Add conversion methods: `asSymbol`, `asNumber`, `printString`
- Add searching: `indexOf:`, `indexOf:ifAbsent:`, `includes:`
- **Run tests after each method implementation**

**Acceptance Criteria**:

- All String tests pass
- String methods work correctly from compiled Smalltalk code
- No regressions in existing functionality

**Simplicity Focus**: Implement the most straightforward version of each method. Don't optimize for edge cases initially.

### 2.3 Complete Array Class Implementation (Test-Driven)

**Goal**: Full Array class with test-driven development

**Tasks**:

- **Write comprehensive Array tests first** (focus on core functionality)
- Implement Array class with simple instance variables
- Add core methods: `size`, `at:`, `at:put:`, `copyFrom:to:`
- Add enumeration: `do:`, `collect:`, `select:`, `reject:`, `detect:`
- Add utility methods: `includes:`, `indexOf:`
- Implement array literals parsing (`#(1 2 3)`)
- **Verify all tests pass after each addition**

**Acceptance Criteria**:

- All Array tests pass
- Array functionality works from Smalltalk source
- Full test suite continues to pass

**Simplicity Focus**: Start with basic array operations. Add complexity only when needed and tested.

### 2.4 Complete Integer and Number Classes (Test-Driven)

**Goal**: Full numeric class hierarchy with comprehensive testing

**Tasks**:

- **Write tests for all numeric operations first**
- Implement Integer class with arithmetic operations
- Add comparison methods: `<`, `>`, `<=`, `>=`, `=`, `~=`
- Add utility methods: `abs`, `negated`, `min:`, `max:`
- Add iteration: `to:do:`, `to:by:do:`, `timesRepeat:`
- Add conversion: `asString`, `printString`
- **Run full test suite after each change**

**Acceptance Criteria**:

- All numeric tests pass
- Arithmetic expressions work from Smalltalk source
- No regressions in existing functionality

**Simplicity Focus**: Implement straightforward numeric operations. Don't worry about floating point or complex number types initially.

---

## Phase 3: Enhanced Parser & Compiler with Test Coverage (Weeks 4-5)

### 3.1 Complete Smalltalk Parser (Test-Driven)

**Goal**: Parse all essential Smalltalk syntax constructs with comprehensive tests

**Tasks**:

- **Write parser tests first** for all syntax constructs
- Extend parser to handle method definitions with parameters
- Add support for class definitions (`Object subclass: #MyClass`)
- Implement proper temporary variable scoping
- Add support for instance variable declarations
- Handle method categories and class-side methods
- Parse complex expressions with proper precedence
- **Run parser tests after each enhancement**

**Acceptance Criteria**:

- All parser tests pass
- Can parse complete Smalltalk class definitions
- Method definitions with multiple parameters work
- Complex expressions parse with correct precedence
- All existing tests continue to pass

**Simplicity Focus**: Implement a straightforward recursive descent parser. Don't worry about error recovery or advanced parsing techniques.

### 3.2 Enhanced Bytecode Compiler (Test-Driven)

**Goal**: Generate correct bytecode for all constructs with test validation

**Tasks**:

- **Write bytecode generation tests first**
- Implement proper variable scoping (instance, temporary, class)
- Add support for block compilation with closures
- Generate correct bytecode for method definitions
- Handle message sends with multiple arguments
- Implement proper return handling
- Generate bytecode for class creation
- **Verify bytecode correctness with execution tests**

**Acceptance Criteria**:

- All bytecode generation tests pass
- Compiled methods execute correctly
- Block closures work with proper variable capture
- Method calls with multiple arguments work
- Class creation bytecode executes properly
- Full test suite continues to pass

**Simplicity Focus**: Generate straightforward bytecode. Don't optimize for size or execution speed.

---

## Phase 4: Core Class Library Completion with TDD (Weeks 6-7)

### 4.1 Block Class Implementation (Test-Driven)

**Goal**: Complete Block class with proper closure semantics and full test coverage

**Tasks**:

- **Write comprehensive Block tests first**
- Implement Block class with simple instance variables
- Add core methods: `value`, `value:`, `value:value:`, `valueWithArguments:`
- Implement basic closure variable capture
- Add utility methods: `whileTrue:`, `whileFalse:`
- Handle block parameters and temporary variables
- **Test closure behavior thoroughly**

**Acceptance Criteria**:

- All Block tests pass
- Block evaluation works correctly
- Closure variables are properly captured
- Block iteration methods work
- No regressions in existing functionality

**Simplicity Focus**: Implement the simplest closure mechanism that works correctly. Don't optimize for complex scoping scenarios initially.

### 4.2 Essential Collection Classes (Test-Driven)

**Goal**: Add only the most essential collection types

**Tasks**:

- **Write tests for essential collection operations first**
- Implement Collection abstract class with basic protocol
- Add Set class with simple uniqueness constraints
- Add basic stream classes (ReadStream, WriteStream)
- **Test collection protocol consistency**

**Acceptance Criteria**:

- All collection tests pass
- Essential collection types work correctly
- Collection protocols are consistent
- Stream operations work properly
- Full test suite continues to pass

**Simplicity Focus**: Implement only the collections actually needed for self-hosting. Skip advanced collection types.

### 4.3 Basic Exception Handling System (Test-Driven)

**Goal**: Minimal exception handling infrastructure

**Tasks**:

- **Write exception handling tests first**
- Implement basic Exception class
- Add Error class for common errors
- Implement simple `on:do:` exception handling
- Add basic `ensure:` cleanup blocks
- Add common exception types needed for compiler (ParseError, etc.)
- **Test exception propagation**

**Acceptance Criteria**:

- All exception tests pass
- Basic exception handling works correctly
- `ensure:` blocks execute properly
- Exception propagation works for compiler needs
- No regressions in existing functionality

**Simplicity Focus**: Implement only the exception handling needed for the compiler and basic system operation.

---

## Phase 5: Image Building & Serialization with TDD (Week 8)

### 5.1 Simple Object Serialization (Test-Driven)

**Goal**: Basic serialize/deserialize object graphs

**Tasks**:

- **Write serialization tests first**
- Implement simple object graph traversal
- Add basic serialization for all core classes
- Handle circular references with simple approach
- Create basic image file format
- **Test round-trip serialization**

**Acceptance Criteria**:

- All serialization tests pass
- Can serialize essential object graphs
- Deserialization recreates objects correctly
- Circular references are handled
- Image files work correctly
- Full test suite continues to pass

**Simplicity Focus**: Use the simplest serialization approach that works. Don't worry about file size or format efficiency.

### 5.2 Basic Image Building Process (Test-Driven)

**Goal**: Build images from source code with simple approach

**Tasks**:

- **Write image building tests first**
- Create simple source file loading mechanism
- Implement basic class installation process
- Add method compilation and installation
- Handle basic dependencies between classes
- Create simple bootstrap image building
- **Test image building with basic source configurations**

**Acceptance Criteria**:

- All image building tests pass
- Can build basic image from source files
- Class dependencies are resolved correctly
- Bootstrap process works reliably
- Full test suite continues to pass

**Simplicity Focus**: Implement the simplest image building process that works. Don't worry about incremental updates or complex dependency resolution.

---

## Phase 6: Self-Hosting Bootstrap with TDD (Week 9)

### 6.1 Bootstrap Compiler (Test-Driven)

**Goal**: Smalltalk compiler that can compile itself

**Tasks**:

- **Write self-compilation tests first**
- Ensure compiler can parse its own source
- Verify compiler can compile its own bytecode
- Test basic recursive compilation process
- Handle essential compiler bootstrapping cases
- Add basic compiler error handling
- **Test self-compilation thoroughly**

**Acceptance Criteria**:

- All self-compilation tests pass
- Compiler successfully compiles itself
- Recursive compilation produces correct results
- Bootstrap process works reliably
- Full test suite continues to pass

**Simplicity Focus**: Implement the simplest self-compilation approach. Don't worry about bootstrapping optimizations.

### 6.2 Self-Hosting Verification (Test-Driven)

**Goal**: Verify complete self-hosting capability

**Tasks**:

- **Write self-hosting verification tests**
- Build image entirely from Smalltalk source
- Verify image can rebuild itself
- Test all core functionality in self-hosted image
- Run complete test suite in self-hosted environment
- Document any limitations
- **Test self-hosting under basic conditions**

**Acceptance Criteria**:

- All self-hosting tests pass
- Self-hosted image passes all tests
- Image can successfully rebuild itself
- All core Smalltalk features work correctly
- Test suite runs successfully in self-hosted environment

**Simplicity Focus**: Verify that basic self-hosting works correctly. Don't worry about edge cases or complex scenarios initially.

---

## Phase 7: Final Testing & Validation (Week 10)

### 7.1 Comprehensive Testing

**Goal**: Thorough validation of self-hosting system

**Tasks**:

- Run complete test suite in multiple configurations
- Test essential edge cases and error conditions
- Validate basic memory management
- Verify compatibility with essential Smalltalk code
- Test image building with various class hierarchies
- **Achieve comprehensive test coverage**

**Acceptance Criteria**:

- All tests pass in self-hosted environment
- System handles essential edge cases gracefully
- System is stable and reliable for basic use
- No critical bugs or crashes

**Simplicity Focus**: Focus on testing that the system works correctly for its intended purpose. Don't over-test edge cases.

### 7.2 Documentation & Examples

**Goal**: Complete documentation with working examples

**Tasks**:

- Document the bootstrap process with simple examples
- Create basic examples of self-hosting usage
- Write essential developer documentation
- Add basic troubleshooting guides
- Document known limitations
- **Ensure all documentation examples work**

**Acceptance Criteria**:

- Documentation is complete and accurate
- All examples work correctly
- Developers can understand and use the system
- Known limitations are clearly documented

**Simplicity Focus**: Write clear, simple documentation that gets developers started quickly.

---

## Success Metrics (Simple & Tested)

### Technical Metrics

- **Test Coverage**: Comprehensive test coverage of all functionality
- **Self-Compilation**: System can compile itself from source (verified by tests)
- **Reliability**: Bootstrap process works consistently (verified by tests)
- **Simplicity**: Code is easy to understand and modify

### Functional Metrics

- **Complete Smalltalk**: All essential Smalltalk features implemented (tested)
- **SUnit Integration**: Full testing framework available and working
- **Image Building**: Can build images from source reliably (tested)
- **Developer Experience**: Easy to understand and use (documented with tests)

## Simplicity-First Benefits

### Development Speed

- **Faster Implementation**: Simple solutions are quicker to implement
- **Easier Debugging**: Simple code is easier to understand and fix
- **Reduced Complexity**: Fewer moving parts means fewer things can break
- **Clear Intent**: Simple code clearly expresses what it's supposed to do

### Maintainability

- **Easy to Modify**: Simple code is easier to change and extend
- **Clear Architecture**: Simple design is easier to understand
- **Fewer Bugs**: Simple code has fewer places for bugs to hide
- **Better Testing**: Simple code is easier to test thoroughly

## Risk Mitigation (Simplicity-Focused)

### Technical Risks

- **Over-Engineering**: Avoided by choosing simplest solutions first
- **Circular Dependencies**: Handled with simple, straightforward approaches
- **Memory Issues**: Addressed with basic, well-tested memory management
- **Compatibility Issues**: Tested with essential Smalltalk code samples

### Quality Risks

- **Regression Bugs**: Prevented by comprehensive test suite
- **Integration Issues**: Caught early by continuous testing
- **Complexity Creep**: Avoided by maintaining simplicity focus
- **Feature Bloat**: Prevented by implementing only essential features

## Timeline Summary

- **Week 1**: SUnit framework + test existing code
- **Weeks 2-3**: Core classes (String, Array, Integer, Boolean) with TDD
- **Weeks 4-5**: Enhanced parser & compiler with TDD
- **Weeks 6-7**: Essential class library with TDD
- **Week 8**: Basic image building & serialization with TDD
- **Week 9**: Self-hosting bootstrap with TDD
- **Week 10**: Final validation & documentation

**Key Principles**:

1. Every change is test-driven
2. Choose the simplest solution that works
3. Implement only what's needed for self-hosting
4. Correctness over cleverness
