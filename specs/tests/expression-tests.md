# Smalltalk Expression Test Suite

This document specifies the canonical test suite for Smalltalk expression evaluation. These tests verify that the VM correctly executes Smalltalk code across all major language features.

## Test Format

Each test consists of:

- **Expression**: The Smalltalk code to execute
- **Expected Result**: The expected return value (as string)
- **Should Pass**: Whether the test should succeed (true) or fail with exception (false)
- **Category**: The feature category being tested

## Test Categories

### 1. Arithmetic Operations

Basic integer arithmetic with primitives 1, 2, 9, 10, 11.

| Expression     | Expected | Status  |
| -------------- | -------- | ------- |
| `3 + 4`        | `7`      | ✅ Pass |
| `5 - 2`        | `3`      | ✅ Pass |
| `2 * 3`        | `6`      | ✅ Pass |
| `10 / 2`       | `5`      | ✅ Pass |
| `(3 + 2) * 4`  | `20`     | ✅ Pass |
| `10 - 2 * 3`   | `24`     | ✅ Pass |
| `(10 - 2) / 4` | `2`      | ✅ Pass |

### 2. Comparison Operations

Integer comparisons with primitives 3-8.

| Expression           | Expected | Status  |
| -------------------- | -------- | ------- |
| `3 < 5`              | `true`   | ✅ Pass |
| `7 > 2`              | `true`   | ✅ Pass |
| `3 = 3`              | `true`   | ✅ Pass |
| `4 ~= 5`             | `true`   | ✅ Pass |
| `4 <= 4`             | `true`   | ✅ Pass |
| `5 >= 3`             | `true`   | ✅ Pass |
| `5 < 3`              | `false`  | ✅ Pass |
| `2 > 7`              | `false`  | ✅ Pass |
| `3 = 4`              | `false`  | ✅ Pass |
| `(3 + 2) < (4 * 2)`  | `true`   | ✅ Pass |
| `(10 - 3) > (2 * 3)` | `true`   | ✅ Pass |
| `(6 / 2) = (1 + 2)`  | `true`   | ✅ Pass |

### 3. Object Creation

Object instantiation with primitives 70, 71, 72.

| Expression     | Expected          | Status  |
| -------------- | ----------------- | ------- |
| `Object new`   | `Object`          | ✅ Pass |
| `Array new: 3` | `<Array size: 3>` | ✅ Pass |

### 4. String Operations

String literals and operations with primitives 63-67.

| Expression           | Expected      | Status  |
| -------------------- | ------------- | ------- |
| `'hello'`            | `hello`       | ✅ Pass |
| `'world'`            | `world`       | ✅ Pass |
| `'hello' , ' world'` | `hello world` | ✅ Pass |
| `'hello' size`       | `5`           | ✅ Pass |

### 5. Literals

Special literal values and their classes.

| Expression    | Expected          | Status  |
| ------------- | ----------------- | ------- |
| `true`        | `true`            | ✅ Pass |
| `false`       | `false`           | ✅ Pass |
| `nil`         | `nil`             | ✅ Pass |
| `#abc`        | `Symbol(abc)`     | ✅ Pass |
| `true class`  | `True`            | ✅ Pass |
| `false class` | `False`           | ✅ Pass |
| `nil class`   | `UndefinedObject` | ✅ Pass |

### 6. Variables

Temporary variable assignment and access.

| Expression             | Expected | Status  |
| ---------------------- | -------- | ------- |
| `\| x \| x := 42. x`   | `42`     | ✅ Pass |
| `\| x \| (x := 5) + 1` | `6`      | ✅ Pass |

### 7. Block Closures

Block creation, execution, and lexical scoping.

#### Basic Blocks

| Expression                              | Expected | Status  |
| --------------------------------------- | -------- | ------- |
| `[] value`                              | `nil`    | ✅ Pass |
| `[3 + 4] value`                         | `7`      | ✅ Pass |
| `[:x \| x + 1] value: 5`                | `6`      | ✅ Pass |
| `[\| x \| x := 5. x + 1] value`         | `6`      | ✅ Pass |
| `[:y \|\| x \| x := 5. x + 7] value: 3` | `12`     | ✅ Pass |

#### Lexical Scoping

| Expression                                      | Expected | Status  |
| ----------------------------------------------- | -------- | ------- |
| `\| y \| y := 3. [\| x \| x := 5. x + y] value` | `8`      | ✅ Pass |
| `\| z y \| y := 3. z := 2. [z + y] value`       | `5`      | ✅ Pass |
| `[self] value`                                  | `Object` | ✅ Pass |

#### Nested Blocks

| Expression                                                                           | Expected | Status  |
| ------------------------------------------------------------------------------------ | -------- | ------- |
| `[\| x \| [\| y \| y := 5. x := y] value. x] value`                                  | `5`      | ✅ Pass |
| `[\| x \| [\| y \| [\| z \| x := 'x'. y := 'y'. z:= 'z'] value. x , y] value] value` | `xy`     | ✅ Pass |

#### Variable Shadowing

| Expression                                                                          | Expected  | Status  |
| ----------------------------------------------------------------------------------- | --------- | ------- |
| `\| x \| x := 'outer'. [:x \| x := 'inner'. x] value: 'param'. x`                   | `outer`   | ✅ Pass |
| `\| x \| x := 'outer'. [\| x \| x := 'inner'. x] value. x`                          | `outer`   | ✅ Pass |
| `\| x \| x := 'outer'. [:x \| [\| z \| x] value] value: 'param'`                    | `param`   | ✅ Pass |
| `\| x \| x := 'outer'. [:x \| [\| z \| x := 'deep'. x] value. x] value: 'param'. x` | `outer`   | ✅ Pass |
| `\| x \| x := 'outer'. [\| y \| [\| z \| x := 'changed'] value. y] value. x`        | `changed` | ✅ Pass |

#### Sequential Mutations

| Expression                                                                                  | Expected | Status  |
| ------------------------------------------------------------------------------------------- | -------- | ------- |
| `\| x \| x := 'O'. [x := x , '1'] value. [x := x , '2'] value. x`                           | `O12`    | ✅ Pass |
| `\| a \| a := 'A'. [\| y \| [\| z \| a := a , '1'. y := 'Y'. z := 'Z'] value. a , y] value` | `A1Y`    | ✅ Pass |

#### Shadow Chains

| Expression                                                                    | Expected | Status  |
| ----------------------------------------------------------------------------- | -------- | ------- |
| `\| x \| x := 'O'. [\| x \| x := 'M'. [\| x \| x := 'I'. x] value] value`     | `I`      | ✅ Pass |
| `\| x \| x := 'o'. [:y \| [\| z \| x , y] value] value: 'mid'`                | `omid`   | ✅ Pass |
| `\| a \| a := 'A'. [\| a \| a := 'M'. [\| b \| a := a , 'i'. a] value] value` | `Mi`     | ✅ Pass |

#### Block Methods

Blocks can receive messages and call their own methods.

| Expression                         | Expected | Status  |
| ---------------------------------- | -------- | ------- |
| `[42] identity`                    | `Object` | ✅ Pass |
| `[42] test`                        | `999`    | ✅ Pass |
| `[42] callTest`                    | `999`    | ✅ Pass |
| `[42] callValue`                   | `42`     | ✅ Pass |
| `[100] ensureSimple: [200]`        | `100`    | ✅ Pass |
| `[100] testTemp: [200]`            | `100`    | ✅ Pass |
| `[100] testAssign: [200]`          | `100`    | ✅ Pass |
| `[100] testSelfValueAssign: [200]` | `100`    | ✅ Pass |
| `[100] ensure: [200]`              | `100`    | ✅ Pass |

### 8. Conditionals

Boolean conditional execution (implemented in Smalltalk, not primitives).

| Expression                           | Expected | Status     |
| ------------------------------------ | -------- | ---------- |
| `(3 < 4) ifTrue: [10] ifFalse: [20]` | `10`     | ✅ Pass    |
| `true ifTrue: [42]`                  | `42`     | ✅ Pass    |
| `false ifTrue: [1] ifFalse: [7]`     | `7`      | ✅ Pass    |
| `42 isNil`                           | `false`  | ⏳ Pending |
| `42 ifNotNil: [7]`                   | `7`      | ⏳ Pending |
| `42 ifNil: [8]`                      | `nil`    | ⏳ Pending |

### 9. Collections

Array and collection operations (primitives 60-62).

| Expression       | Expected | Status     |
| ---------------- | -------- | ---------- |
| `#(1 2 3) at: 2` | `2`      | ⏳ Pending |
| `#(1 2 3) size`  | `3`      | ⏳ Pending |

### 10. Dictionaries

Dictionary operations (primitives 700-703).

| Expression       | Expected       | Status     |
| ---------------- | -------------- | ---------- |
| `Dictionary new` | `<Dictionary>` | ⏳ Pending |

### 11. Exception Handling

Exception signaling and handling.

#### Expected Exceptions

These tests should fail with specific exceptions:

| Expression                 | Expected Exception     | Status         |
| -------------------------- | ---------------------- | -------------- |
| `10 / 0`                   | `ZeroDivisionError`    | ❌ Should Fail |
| `undefined_variable`       | `NameError`            | ❌ Should Fail |
| `'hello' at: 10`           | `IndexError`           | ❌ Should Fail |
| `Object new unknownMethod` | `MessageNotUnderstood` | ❌ Should Fail |
| `Array new: -1`            | `ArgumentError`        | ❌ Should Fail |

#### Exception Handling Constructs

| Expression                                             | Expected            | Status     |
| ------------------------------------------------------ | ------------------- | ---------- |
| `[10 / 0] ensure: [42]`                                | `42`                | ⏳ Pending |
| `[10 / 0] on: ZeroDivisionError do: [:ex \| 'caught']` | `caught`            | ⏳ Pending |
| `[1 + 2] ensure: [3 + 4]`                              | `3`                 | ✅ Pass    |
| `ZeroDivisionError signal: 'test error'`               | `ZeroDivisionError` | ⏳ Pending |

### 12. Class Creation

Dynamic class creation.

| Expression                | Expected         | Status     |
| ------------------------- | ---------------- | ---------- |
| `Object subclass: #Point` | `<Class: Point>` | ⏳ Pending |

### 13. Method Execution

Testing the executeMethod API.

| Expression | Expected | Status  |
| ---------- | -------- | ------- |
| `^ 42`     | `42`     | ✅ Pass |

## Test Status Legend

- ✅ **Pass**: Test currently passes
- ⏳ **Pending**: Feature not yet implemented
- ❌ **Should Fail**: Test should raise an exception

## Implementation Notes

### Primitives Required

- **1-11**: Integer arithmetic and comparison
- **60-62**: Array operations
- **63-67**: String operations
- **70-72**: Object creation
- **75**: Identity hash
- **111**: Class introspection
- **201-202**: Block execution
- **700-703**: Dictionary operations

### Bytecode Required

- **0**: PUSH_LITERAL
- **1-2**: PUSH_INSTANCE_VARIABLE, PUSH_TEMPORARY_VARIABLE
- **3**: PUSH_SELF
- **4-5**: STORE_INSTANCE_VARIABLE, STORE_TEMPORARY_VARIABLE
- **6**: SEND_MESSAGE
- **7**: RETURN_STACK_TOP
- **8-10**: JUMP, JUMP_IF_FALSE, JUMP_IF_TRUE
- **11-12**: POP, DUPLICATE
- **13-14**: CREATE_BLOCK, EXECUTE_BLOCK

### Smalltalk Methods Required

Boolean conditionals should be implemented as methods in True/False classes:

- `True>>ifTrue:`, `True>>ifFalse:`, `True>>ifTrue:ifFalse:`
- `False>>ifTrue:`, `False>>ifFalse:`, `False>>ifTrue:ifFalse:`

## Running Tests

The test suite is implemented in `src/cpp/tests/all_expressions_test.cpp` and can be run with:

```bash
make test
```

Or directly:

```bash
./build/bin/vm_tests
```

## Test Coverage Summary

Current implementation status by category:

- ✅ **Arithmetic**: 7/7 tests passing
- ✅ **Comparison**: 12/12 tests passing
- ✅ **Object Creation**: 2/2 tests passing
- ✅ **Strings**: 4/4 tests passing
- ✅ **Literals**: 7/7 tests passing
- ✅ **Variables**: 2/2 tests passing
- ✅ **Blocks**: 30/30 tests passing
- ✅ **Block Methods**: 9/9 tests passing
- ⏳ **Conditionals**: 3/6 tests passing (need isNil, ifNil:, ifNotNil:)
- ⏳ **Collections**: 0/2 tests passing
- ⏳ **Dictionaries**: 0/1 tests passing
- ⏳ **Exception Handling**: 1/9 tests passing
- ⏳ **Class Creation**: 0/1 tests passing
- ✅ **Method Execution**: 1/1 tests passing

**Total**: 77/92 tests passing (84%)
