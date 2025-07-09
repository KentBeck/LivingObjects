# Implementation Plan: Self-Hosted Smalltalk Compiler

## Goal: Parse and Compile Smalltalk Code Inside the VM

Enable Living Objects to compile Smalltalk source code from within Smalltalk itself, eliminating the C++ parser and compiler dependencies.

## Current Status

### ✅ Already Working
- Basic VM with bytecode execution
- Method lookup and message passing
- Block creation and execution (9/9 block tests passing)
- Exception hierarchy defined
- Integer, String, and Array primitives
- Method compilation from Smalltalk source (via C++ MethodCompiler)

### ❌ Not Yet Working
- Exception handling (ensure:, on:do:)
- Class creation from Smalltalk
- Collection classes beyond basic Array
- Testing framework (SUnit)
- Parser in Smalltalk
- Compiler in Smalltalk

## Phase 1: Exception Handling (Week 1)
**Goal: Complete exception handling to enable proper error recovery**

### Tasks
1. **Implement ensure: primitive**
   - Add primitive for ensure: that handles unwind protection
   - Ensure cleanup blocks run even during exceptions
   - Test: `[10/0] ensure: [Transcript show: 'cleanup']`

2. **Implement on:do: primitive**
   - Add primitive for exception catching
   - Support exception class matching
   - Test: `[10/0] on: ZeroDivisionError do: [:ex | 'caught']`

3. **Stack unwinding**
   - Implement proper stack unwinding during exceptions
   - Preserve exception information during unwind
   - Support non-local returns through ensure: blocks

### Verification
```smalltalk
"All these should work:"
[1/0] ensure: [cleanupDone := true].
[1/0] on: ZeroDivisionError do: [:ex | ex description].
[[1/0] ensure: [x := 1]] on: Error do: [:ex | x]. "x = 1"
```

## Phase 2: Core Classes (Week 2)
**Goal: Implement essential classes needed for compilation**

### OrderedCollection
```smalltalk
Object subclass: #OrderedCollection
    instanceVariableNames: 'array firstIndex lastIndex'
    classVariableNames: ''
```
- Basic add:, at:, size, do:
- Growing/shrinking behavior
- Used by parser for AST nodes

### Dictionary  
```smalltalk
Object subclass: #Dictionary
    instanceVariableNames: 'keys values tally'
    classVariableNames: ''
```
- at:put:, at:, at:ifAbsent:
- keys, values, do:
- Used for symbol tables, method dictionaries

### Stream Classes
```smalltalk
Object subclass: #Stream
Object subclass: #ReadStream
Object subclass: #WriteStream
```
- Basic stream protocol for parsing
- next, peek, atEnd
- nextPut:, nextPutAll:

### Verification
```smalltalk
| dict |
dict := Dictionary new.
dict at: #name put: 'Living Objects'.
dict at: #name "Returns 'Living Objects'"
```

## Phase 3: Testing Framework (Week 3)
**Goal: Minimal SUnit implementation for testing**

### TestCase
```smalltalk
Object subclass: #TestCase
    instanceVariableNames: 'testSelector'
    classVariableNames: ''

TestCase>>run
    self setUp.
    [self perform: testSelector]
        ensure: [self tearDown]
```

### TestResult
```smalltalk
Object subclass: #TestResult
    instanceVariableNames: 'passed failed errors'
    classVariableNames: ''
```

### Basic Assertions
- assert:
- deny:  
- assert:equals:
- should:raise:

### Verification
```smalltalk
TestCase subclass: #ExampleTest
    instanceVariableNames: ''
    classVariableNames: ''

ExampleTest>>testAddition
    self assert: 2 + 2 equals: 4

ExampleTest new run "Returns TestResult"
```

## Phase 4: Parser in Smalltalk (Week 4-5)
**Goal: Recursive descent parser written in Smalltalk**

### Scanner
```smalltalk
Object subclass: #Scanner
    instanceVariableNames: 'stream currentChar token'
    classVariableNames: ''

Scanner>>nextToken
    "Return next token from stream"
```

### Parser
```smalltalk
Object subclass: #Parser
    instanceVariableNames: 'scanner currentToken'
    classVariableNames: ''

Parser>>parseMethod
    "Return MethodNode AST"
```

### AST Nodes
```smalltalk
Object subclass: #ASTNode
ASTNode subclass: #MethodNode
ASTNode subclass: #MessageNode
ASTNode subclass: #LiteralNode
ASTNode subclass: #VariableNode
ASTNode subclass: #BlockNode
```

### Verification
```smalltalk
| parser ast |
parser := Parser on: 'foo: x ^ x + 1'.
ast := parser parseMethod.
ast selector "Returns #foo:"
```

## Phase 5: Compiler in Smalltalk (Week 6-7)
**Goal: Bytecode compiler written in Smalltalk**

### Compiler
```smalltalk
Object subclass: #Compiler
    instanceVariableNames: 'methodNode bytecodes literals'
    classVariableNames: ''

Compiler>>compile: aString
    | parser ast |
    parser := Parser on: aString.
    ast := parser parseMethod.
    ^ self compileAST: ast
```

### BytecodeBuilder
```smalltalk
Object subclass: #BytecodeBuilder
    instanceVariableNames: 'bytecodes literals tempCount'
    classVariableNames: ''

BytecodeBuilder>>pushLiteral: anObject
    literals add: anObject.
    bytecodes add: PushLiteral.
    bytecodes add: literals size - 1
```

### Verification
```smalltalk
| method |
method := Compiler new compile: 'double: x ^ x * 2'.
3 perform: method selector with: 5 "Returns 10"
```

## Phase 6: Integration (Week 8)
**Goal: Replace C++ compiler with Smalltalk compiler**

### Class Creation
```smalltalk
Class>>compile: sourceString
    | method |
    method := Compiler new compile: sourceString.
    self addMethod: method

Object subclass: #Point 
    instanceVariableNames: 'x y'

Point compile: 'x ^ x'.
Point compile: 'y ^ y'.
Point compile: 'x: aNumber x := aNumber'.
```

### Interactive Development
```smalltalk
Workspace>>evaluate: aString
    | method result |
    method := Compiler new compile: 'doIt ^ ', aString.
    result := nil perform: #doIt.
    ^ result
```

### Verification
```smalltalk
"This should work from within Smalltalk:"
Object subclass: #Counter instanceVariableNames: 'count'.
Counter compile: 'initialize count := 0'.
Counter compile: 'increment count := count + 1'.
Counter compile: 'count ^ count'.

| c |
c := Counter new.
c increment.
c count "Returns 1"
```

## Technical Considerations

### Memory Management
- Parser and compiler will create many temporary objects
- Need efficient garbage collection
- Consider object pooling for AST nodes

### Error Handling  
- Parser must handle syntax errors gracefully
- Compiler must report semantic errors
- Need source position tracking for error messages

### Bootstrap Process
- Initial parser/compiler can be simple
- Once working, can recompile itself with improvements
- Keep C++ compiler as fallback during development

### Performance
- Parser/compiler performance not critical initially
- Can optimize later (memoization, better algorithms)
- Focus on correctness first

## Success Criteria

### Phase 1: Exception Handling
- [ ] ensure: blocks execute during normal and exception paths
- [ ] on:do: catches and handles exceptions
- [ ] Nested exception handlers work correctly
- [ ] 5/5 exception tests pass

### Phase 2: Core Classes  
- [ ] OrderedCollection supports parser needs
- [ ] Dictionary works for symbol tables
- [ ] Streams support scanning/parsing
- [ ] Can create classes from Smalltalk

### Phase 3: Testing Framework
- [ ] Can write and run TestCase subclasses
- [ ] Assertions work correctly  
- [ ] Test results are collected
- [ ] Can test parser/compiler incrementally

### Phase 4: Parser
- [ ] Parses method syntax correctly
- [ ] Produces correct AST
- [ ] Reports syntax errors  
- [ ] Handles all Smalltalk expressions

### Phase 5: Compiler
- [ ] Generates correct bytecode
- [ ] Handles variables properly
- [ ] Compiles blocks correctly
- [ ] Resulting methods execute properly

### Phase 6: Integration
- [ ] Can define new classes in Smalltalk
- [ ] Can add methods to classes
- [ ] Can evaluate arbitrary expressions
- [ ] C++ compiler no longer needed

## Benefits of Self-Hosted Compiler

1. **Dogfooding** - Living Objects compiles itself
2. **Flexibility** - Easy to experiment with language features
3. **Debugging** - Can debug compiler in Smalltalk
4. **Learning** - Great example of Living Objects capabilities
5. **Community** - Smalltalkers can contribute without C++ knowledge

## Next Steps After Completion

With self-hosted compilation working:
- Implement LSP server in Smalltalk
- Add incremental compilation
- Implement code formatting
- Add refactoring support
- Build class browser
- Create development tools

---

*This plan enables Living Objects to become a true Smalltalk system where everything, including the compiler, is an object that can be inspected, modified, and extended at runtime.*