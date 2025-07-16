# Design Document

## Overview

The C++ Smalltalk VM has a comprehensive expression test suite, but several categories of tests are failing due to parsing limitations and missing implementations. This design addresses the systematic fixes needed to make all "should pass" tests actually pass.

Based on the analysis of the current codebase and test failures, the main issues are:

1. **Parser limitations** - Missing support for certain syntax constructs
2. **Missing primitive implementations** - Some basic operations aren't implemented
3. **Incomplete method dispatch** - Some message sends aren't properly handled
4. **Block execution issues** - Block creation and execution has problems

## Architecture

The fix will involve modifications to several key components:

### Parser Enhancements (SimpleParser)

- Fix array creation syntax parsing (`Array new: 3`)
- Fix string operation parsing (`'hello' , ' world'`)
- Ensure variable declaration parsing works correctly
- Fix block parameter parsing (`:x |`)
- Add support for conditional message parsing (`ifTrue:ifFalse:`)

### Primitive Method Implementations

- String concatenation (`,` operator) - already implemented
- String size method - already implemented
- Array creation and access methods
- Conditional execution methods (`ifTrue:`, `ifFalse:`, `ifTrue:ifFalse:`)

### Block Execution System

- Fix block value execution
- Fix parameterized block execution
- Ensure proper context creation for blocks

### Collection Support

- Array literal parsing (`#(1 2 3)`)
- Array access methods (`at:`, `size`)
- Symbol literal parsing (`#abc`)

## Components and Interfaces

### 1. SimpleParser Enhancements

**Current Issues:**

- Parse error on `Array new: 3` - "Unexpected characters at end of input"
- Parse error on `'hello' , ' world'` - "Unexpected characters at end of input"
- Parse error on `| x | x := 42. x` - "Unexpected character: |"
- Parse error on `[:x | x + 1] value: 5` - "Unexpected character: :"

**Required Changes:**

- Fix keyword message parsing to handle `new:` properly
- Ensure binary operator parsing handles `,` correctly
- Fix temporary variable parsing (already implemented but has issues)
- Fix block parameter parsing

### 2. Method Dispatch System

**Current Issues:**

- `'hello' size` fails with "Method not found: size"
- Block execution fails with "Stack pointer below stack start"

**Required Changes:**

- Ensure string primitive methods are properly registered
- Fix block execution context management
- Add missing conditional methods to Boolean classes

### 3. Collection Implementation

**Current Issues:**

- `#(1 2 3) at: 2` fails with parse error
- `#abc` symbol parsing needs verification

**Required Changes:**

- Implement array access primitives
- Ensure symbol creation works correctly
- Add array size method

### 4. Boolean and Conditional System

**Current Issues:**

- `true ifTrue: [42]` fails with parse error
- Conditional expressions not implemented

**Required Changes:**

- Add `ifTrue:`, `ifFalse:`, `ifTrue:ifFalse:` methods to Boolean classes
- Implement proper conditional execution

## Data Models

### Expression Test Structure

```cpp
struct ExpressionTest {
    std::string expression;
    std::string expectedResult;
    bool shouldPass;
    std::string category;
};
```

### Parser State

The parser maintains position and input state, with methods for:

- Tokenization and lookahead
- Error reporting with position information
- Recursive descent parsing

### Primitive Registry

Maps primitive numbers to implementation functions:

```cpp
PrimitiveFunction = TaggedValue(*)(TaggedValue receiver,
                                  const std::vector<TaggedValue>& args,
                                  Interpreter& interpreter);
```

## Error Handling

### Parse Errors

- Provide clear error messages with position information
- Handle unexpected characters gracefully
- Support recovery where possible

### Runtime Errors

- Proper exception handling for method not found
- Index out of bounds for collections
- Type checking for primitive operations

### Test Framework Integration

- Clear pass/fail reporting
- Categorized test results
- Detailed failure information

## Testing Strategy

### Test Categories to Fix

1. **object_creation** - Fix `Array new: 3` parsing
2. **string_operations** - Fix `'hello' , ' world'` and ensure `'hello' size` works
3. **variables** - Fix `| x | x := 42. x` parsing
4. **blocks** - Fix `[3 + 4] value` and `[:x | x + 1] value: 5`
5. **conditionals** - Implement `ifTrue:ifFalse:` methods
6. **collections** - Fix `#(1 2 3)` parsing and array methods

### Verification Approach

1. Run expression tests after each fix
2. Verify that previously passing tests still pass
3. Ensure error messages are clear for genuinely invalid syntax
4. Test edge cases and boundary conditions

### Integration Testing

- Test combinations of fixed features
- Verify complex expressions work correctly
- Ensure performance is not degraded

## Implementation Priority

### Phase 1: Parser Fixes

1. Fix keyword message parsing for `Array new: 3`
2. Fix binary operator parsing for string concatenation
3. Fix temporary variable declaration parsing
4. Fix block parameter parsing

### Phase 2: Method Implementation

1. Verify string primitive methods are working
2. Implement array creation and access methods
3. Add boolean conditional methods
4. Fix block execution context issues

### Phase 3: Collection Support

1. Ensure array literal parsing works
2. Implement array primitive methods
3. Verify symbol creation and handling

### Phase 4: Integration and Testing

1. Run full test suite
2. Fix any remaining issues
3. Verify all "should pass" tests actually pass
4. Document any limitations or known issues
