# Core Classes & Methods Analysis for Self-Hosting Smalltalk

This document analyzes the core classes and methods needed for implementing the three fundamental components of a self-hosting Smalltalk system: the testing framework, parser, and compiler.

## SUnit Testing Framework

### Core Classes Needed:

#### TestCase

```smalltalk
TestCase
  - setUp, tearDown
  - assert:, deny:, assert:equals:
  - should:raise:, shouldnt:raise:
  - fail, fail:
  - run, runCase
```

#### TestSuite

```smalltalk
TestSuite
  - addTest:, addTestsFor:
  - run, tests
  - name, name:
```

#### TestResult

```smalltalk
TestResult
  - initialize, runCount, errorCount, failureCount
  - addError:, addFailure:, addPass:
  - errors, failures, passed
  - printOn:
```

#### TestRunner

```smalltalk
TestRunner
  - run:, runSuite:
  - report, printResults
```

### Supporting Classes Needed:

#### Exception (basic hierarchy)

```smalltalk
Exception
  - signal, signal:
  - messageText, messageText:

Error (subclass of Exception)
  - for test failures and errors
```

#### Collection protocol

```smalltalk
Collection protocol
  - OrderedCollection (already have)
  - do:, collect:, select: (for managing test lists)
```

#### String

```smalltalk
String
  - =, printString, asString (for test names and messages)
  - , (concatenation for building messages)
```

#### Stream

```smalltalk
Stream
  - WriteStream for building output
  - nextPutAll:, contents
```

## Parser

### Core Classes Needed:

#### Parser

```smalltalk
Parser
  - parseMethod, parseClass
  - parseExpression, parseStatement
  - parseBlock, parseArray
  - parseTemporaries, parseArguments
```

#### Token

```smalltalk
Token
  - type, value, position
  - isIdentifier, isKeyword, isLiteral, etc.
```

#### Scanner/Tokenizer

```smalltalk
Scanner/Tokenizer
  - nextToken, peekToken
  - position, atEnd
  - scanIdentifier, scanNumber, scanString, scanSymbol
```

#### ParseNode (AST nodes)

```smalltalk
ParseNode (AST nodes)
  - MethodNode, ClassNode, BlockNode
  - MessageNode, VariableNode, LiteralNode
  - AssignmentNode, ReturnNode
```

#### SyntaxError

```smalltalk
SyntaxError
  - position, message
  - signal:at:
```

### Supporting Classes Needed:

#### String

```smalltalk
String
  - at:, size, copyFrom:to:
  - first, last, isEmpty
  - isDigit, isLetter, isAlphaNumeric
  - asSymbol, asNumber
```

#### Character

```smalltalk
Character
  - isDigit, isLetter, isAlphaNumeric
  - isSeparator, isUppercase, isLowercase
  - =, asString
```

#### OrderedCollection

```smalltalk
OrderedCollection
  - add:, at:, size, do:
  - first, last, isEmpty
```

#### Dictionary

```smalltalk
Dictionary
  - at:put:, at:ifAbsent:
  - includesKey: (for keyword tables)
```

#### Symbol

```smalltalk
Symbol
  - = (for comparing selectors)
  - asString, printString
```

## Compiler

### Core Classes Needed:

#### Compiler

```smalltalk
Compiler
  - compile:, compileMethod:in:
  - compileBlock:, compileExpression:
  - generateBytecode, optimize
```

#### CompiledMethod

```smalltalk
CompiledMethod
  - selector, bytecodes, literals
  - numArgs, numTemps
  - primitive, header
```

#### BytecodeGenerator

```smalltalk
BytecodeGenerator
  - pushLiteral:, pushTemp:, pushInstVar:
  - storeTemp:, storeInstVar:
  - send:numArgs:, return
  - jump:, jumpIf:, jumpIfNot:
```

#### MethodContext

```smalltalk
MethodContext
  - receiver, arguments, temporaries
  - method, sender, pc
  - push:, pop, top
```

#### Bytecode constants

```smalltalk
Bytecode constants
  - PUSH_LITERAL, PUSH_TEMP, etc.
  - SEND_MESSAGE, RETURN, etc.
```

### Supporting Classes Needed:

#### Class

```smalltalk
Class
  - name, superclass, methods
  - addMethod:, lookupMethod:
  - instanceVariableNames
  - new (for creating instances)
```

#### Method

```smalltalk
Method
  - selector, bytecodes, literals
  - numArgs, numTemps
  - valueWithReceiver:arguments:
```

#### Symbol

```smalltalk
Symbol
  - = (for selector comparison)
  - hash (for method lookup)
```

#### Array

```smalltalk
Array
  - at:, at:put:, size
  - new: (for literals array)
```

#### ByteArray

```smalltalk
ByteArray
  - at:, at:put:, size
  - new: (for bytecode storage)
```

#### Integer

```smalltalk
Integer
  - +, -, *, / (for bytecode generation)
  - =, <, > (for comparisons)
  - bitAnd:, bitOr:, bitShift: (for encoding)
```

#### Block

```smalltalk
Block
  - value, value:, value:value:
  - numArgs, method
```

## Priority Order for Implementation

### Phase 1: Testing Foundation

1. **Exception** (basic signal/handle)
2. **String** (=, printString, ,)
3. **TestCase** (assert:, deny:, setUp, tearDown)
4. **TestResult** (basic counting)
5. **TestRunner** (basic execution)

### Phase 2: Parser Foundation

1. **Character** (classification methods)
2. **String** (at:, size, copyFrom:to:, character methods)
3. **Symbol** (=, asString, creation)
4. **Scanner** (basic tokenization)
5. **ParseNode** hierarchy (basic AST)
6. **Parser** (expression parsing first)

### Phase 3: Compiler Foundation

1. **Array** (literals storage)
2. **ByteArray** (bytecode storage)
3. **Integer** (bytecode encoding)
4. **CompiledMethod** (basic structure)
5. **BytecodeGenerator** (basic opcodes)
6. **Compiler** (expression compilation first)

## Critical Dependencies

### String is Fundamental

- Needed by testing (test names, error messages)
- Needed by parser (source code manipulation)
- Needed by compiler (selector handling)
- **Should be implemented first after basic testing**

### Symbol System

- Critical for method selectors
- Needed for variable names
- Required for proper identity semantics
- **Must work correctly for compiler**

### Collection Protocol

- OrderedCollection (already have)
- Array (for literals and arguments)
- Dictionary (for method lookup, variable binding)
- **Array is most critical for compiler**

### Block Closures

- Needed for testing (should:raise:, etc.)
- Needed for parser (error handling blocks)
- Needed for compiler (code generation blocks)
- **Complex but essential for all three systems**

## Recommended Implementation Order

Based on the dependency analysis, the implementation order should be:

1. **Basic Exception + String** (enables testing)
2. **TestCase + TestRunner** (enables TDD)
3. **Character + Symbol** (enables parsing)
4. **Array + Integer** (enables compilation)
5. **Block** (enables advanced features in all systems)

This order ensures that each component has its dependencies available when needed, while maintaining the ability to test each component as it's developed.

## Notes

- This analysis focuses on the minimal set of classes and methods needed for self-hosting
- Additional convenience methods can be added later once the core functionality is working
- The emphasis is on simplicity and correctness rather than completeness or performance
- Each class should be implemented with comprehensive test coverage before moving to the next
