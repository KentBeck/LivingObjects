Object subclass: #SmalltalkMethod
    instanceVariableNames: 'selector bytecodes literals temporaryVariables class'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!SmalltalkMethod methodsFor: 'initialization'!
initialize
    bytecodes := OrderedCollection new.
    literals := OrderedCollection new.
    temporaryVariables := OrderedCollection new
!

selector: aSymbol
    selector := aSymbol
!

class: aClass
    class := aClass
!

addBytecode: aBytecode
    bytecodes add: aBytecode
!

addInt32: anInteger
    "Add a 32-bit integer to the bytecode array as 4 bytes"
    bytecodes add: ((anInteger bitShift: -24) bitAnd: 16rFF).
    bytecodes add: ((anInteger bitShift: -16) bitAnd: 16rFF).
    bytecodes add: ((anInteger bitShift: -8) bitAnd: 16rFF).
    bytecodes add: (anInteger bitAnd: 16rFF)
!

addPushInstanceVariable: offset
    "Add a bytecode to push an instance variable with the given offset"
    bytecodes add: BytecodeSet PUSH_INSTANCE_VARIABLE.
    self addInt32: offset
!

addPushTemporaryVariable: offset
    "Add a bytecode to push a temporary variable with the given offset"
    bytecodes add: BytecodeSet PUSH_TEMPORARY_VARIABLE.
    self addInt32: offset
!

addStoreInstanceVariable: offset
    "Add a bytecode to store into an instance variable with the given offset"
    bytecodes add: BytecodeSet STORE_INSTANCE_VARIABLE.
    self addInt32: offset
!

addStoreTemporaryVariable: offset
    "Add a bytecode to store into a temporary variable with the given offset"
    bytecodes add: BytecodeSet STORE_TEMPORARY_VARIABLE.
    self addInt32: offset
!

addSendMessage: selectorIndex argCount: argCount
    "Add a bytecode to send a message with the given selector index and argument count"
    bytecodes add: BytecodeSet SEND_MESSAGE.
    self addInt32: selectorIndex.
    self addInt32: argCount
!

addJump: target
    "Add a bytecode to jump to the given target"
    bytecodes add: BytecodeSet JUMP.
    self addInt32: target
!

addJumpIfTrue: target
    "Add a bytecode to jump to the given target if the top of the stack is true"
    bytecodes add: BytecodeSet JUMP_IF_TRUE.
    self addInt32: target
!

addJumpIfFalse: target
    "Add a bytecode to jump to the given target if the top of the stack is false"
    bytecodes add: BytecodeSet JUMP_IF_FALSE.
    self addInt32: target
!

addLiteral: aLiteral
    literals add: aLiteral.
    ^literals size - 1 "Return the index of the added literal"
!

addTemporaryVariable: aSymbol
    temporaryVariables add: aSymbol.
    ^temporaryVariables size - 1 "Return the index of the added temporary variable"
!

!SmalltalkMethod methodsFor: 'accessing'!
selector
    ^selector
!

bytecodes
    ^bytecodes
!

literals
    ^literals
!

temporaryVariables
    ^temporaryVariables
!

class
    ^class
!
