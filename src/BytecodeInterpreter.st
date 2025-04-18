Object subclass: #BytecodeInterpreter
    instanceVariableNames: 'stack method pc context result'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!BytecodeInterpreter methodsFor: 'initialization'!
initialize
    stack := OrderedCollection new.
    pc := 1.
    context := nil.
    result := nil
!

method: aBytecodeArray
    method := aBytecodeArray
!

context: aContext
    context := aContext
!

!BytecodeInterpreter methodsFor: 'execution'!
execute
    "Execute the method bytecodes until a return"
    [self hasMoreBytecodes and: [result isNil]] whileTrue: [self step].
    ^result
!

hasMoreBytecodes
    ^pc <= method size
!

step
    | bytecode instructionSize |
    bytecode := method at: pc.
    self interpretBytecode: bytecode.
    instructionSize := BytecodeSet instructionSize: bytecode.
    pc := pc + instructionSize
!

readInt32At: offset
    "Read a 32-bit integer from the bytecode array at the given offset"
    | value |
    value := 0.
    value := value + ((method at: offset) bitShift: 24).
    value := value + ((method at: offset + 1) bitShift: 16).
    value := value + ((method at: offset + 2) bitShift: 8).
    value := value + (method at: offset + 3).
    ^value
!

interpretBytecode: aBytecode
    "Interpret a bytecode based on our bytecode set"
    (BytecodeSet isPushLiteral: aBytecode)
        ifTrue: [^self pushLiteral: (self readInt32At: pc + 1)].

    (BytecodeSet isPushInstanceVariable: aBytecode)
        ifTrue: [^self pushInstanceVariable: (self readInt32At: pc + 1)].

    (BytecodeSet isPushTemporaryVariable: aBytecode)
        ifTrue: [^self pushTemporaryVariable: (self readInt32At: pc + 1)].

    (BytecodeSet isPushSelf: aBytecode)
        ifTrue: [^self pushSelf].

    (BytecodeSet isStoreInstanceVariable: aBytecode)
        ifTrue: [^self storeInstanceVariable: (self readInt32At: pc + 1)].

    (BytecodeSet isStoreTemporaryVariable: aBytecode)
        ifTrue: [^self storeTemporaryVariable: (self readInt32At: pc + 1)].

    (BytecodeSet isSendMessage: aBytecode)
        ifTrue: [
            | selectorIndex argCount |
            selectorIndex := self readInt32At: pc + 1.
            argCount := self readInt32At: pc + 5.
            ^self sendMessageWithSelector: selectorIndex argCount: argCount].

    aBytecode = BytecodeSet RETURN_STACK_TOP ifTrue: [^self returnStackTop].

    aBytecode = BytecodeSet JUMP ifTrue: [^self jumpTo: (self readInt32At: pc + 1)].
    aBytecode = BytecodeSet JUMP_IF_TRUE ifTrue: [^self jumpIfTrueTo: (self readInt32At: pc + 1)].
    aBytecode = BytecodeSet JUMP_IF_FALSE ifTrue: [^self jumpIfFalseTo: (self readInt32At: pc + 1)].

    aBytecode = BytecodeSet POP ifTrue: [^self pop].
    aBytecode = BytecodeSet DUPLICATE ifTrue: [^self duplicate].

    self error: 'Unknown bytecode: ', aBytecode printString
!

!BytecodeInterpreter methodsFor: 'stack operations'!

pushLiteral: index
    "Push a literal from the method's literal frame"
    stack addLast: (context method literals at: index + 1)
!

pushInstanceVariable: index
    "Push an instance variable from the receiver"
    | name |
    name := context receiver class instanceVariableNames at: index + 1.
    stack addLast: (context receiver at: name)
!

pushTemporaryVariable: index
    "Push a temporary variable from the context"
    stack addLast: (context temporaryAt: index)
!

pushSelf
    "Push the receiver"
    stack addLast: context receiver
!

storeInstanceVariable: index
    "Store into an instance variable of the receiver"
    | value name |
    value := stack removeLast.
    name := context receiver class instanceVariableNames at: index + 1.
    context receiver at: name put: value.
    stack addLast: value
!

storeTemporaryVariable: index
    "Store into a temporary variable"
    | value |
    value := stack removeLast.
    context temporaryAt: index put: value.
    stack addLast: value
!

pop
    "Pop the top value from the stack"
    stack removeLast
!

duplicate
    "Duplicate the top value on the stack"
    stack addLast: stack last
!

!BytecodeInterpreter methodsFor: 'control flow'!
jumpTo: target
    "Jump to a target bytecode index"
    pc := target
!

jumpIfTrueTo: target
    "Jump if the top of the stack is true"
    | condition |
    condition := stack removeLast.
    condition ifTrue: [pc := target]
!

jumpIfFalseTo: target
    "Jump if the top of the stack is false"
    | condition |
    condition := stack removeLast.
    condition ifFalse: [pc := target]
!

!BytecodeInterpreter methodsFor: 'message sending'!
sendMessageWithSelector: selectorIndex argCount: argCount
    "Send a message with the given selector index and argument count"
    | args selector receiver method newContext |

    "Pop arguments from the stack in reverse order"
    args := OrderedCollection new.
    argCount timesRepeat: [args addFirst: stack removeLast].

    "Get the selector from the literals array and pop the receiver"
    selector := context method literals at: selectorIndex + 1.
    receiver := stack removeLast.

    "Look up the method"
    method := receiver class lookupMethod: selector.
    method ifNil: [self error: 'Method not found: ', selector].

    "Create a new context"
    newContext := SmalltalkContext new
        receiver: receiver;
        method: method;
        arguments: args;
        sender: context;
        yourself.

    "Create a new interpreter to execute the method"
    result := (BytecodeInterpreter new)
        method: method bytecodes;
        context: newContext;
        execute
!

!BytecodeInterpreter methodsFor: 'returns'!
returnStackTop
    "Return the value on top of the stack"
    result := stack removeLast