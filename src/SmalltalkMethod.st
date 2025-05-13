Object subclass: #SmalltalkMethod
    instanceVariableNames: 'selector bytecodes literals temporaryVariables class bytecodePosition'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!SmalltalkMethod methodsFor: 'initialization'!
initialize
    bytecodes := ByteArray new: 64.  "Start with a reasonable size"
    bytecodePosition := 1.
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
    "Add a single bytecode to the bytecode array"
    self ensureSpace: 1.
    bytecodes at: bytecodePosition put: aBytecode.
    bytecodePosition := bytecodePosition + 1
!

addInt32: anInteger
    "Add a 32-bit integer to the bytecode array as 4 bytes"
    self ensureSpace: 4.
    bytecodes at: bytecodePosition put: (anInteger bitAnd: 16rFF).
    bytecodes at: bytecodePosition + 1 put: ((anInteger bitShift: -8) bitAnd: 16rFF).
    bytecodes at: bytecodePosition + 2 put: ((anInteger bitShift: -16) bitAnd: 16rFF).
    bytecodes at: bytecodePosition + 3 put: ((anInteger bitShift: -24) bitAnd: 16rFF).
    bytecodePosition := bytecodePosition + 4
!

ensureSpace: additionalBytes
    "Ensure there's enough space in the bytecode array"
    | newSize newBytecodes |
    bytecodePosition + additionalBytes - 1 > bytecodes size ifTrue: [
        newSize := bytecodes size * 2 max: bytecodePosition + additionalBytes.
        newBytecodes := ByteArray new: newSize.
        1 to: bytecodes size do: [:i | 
            newBytecodes at: i put: (bytecodes at: i)].
        bytecodes := newBytecodes
    ]
!

finalizeBytecodes
    "Trim the bytecode array to the actual size used"
    bytecodes := bytecodes copyFrom: 1 to: bytecodePosition - 1
!

addPushInstanceVariable: offset
    "Add a bytecode to push an instance variable with the given offset"
    self addBytecode: BytecodeSet PUSH_INSTANCE_VARIABLE.
    self addInt32: offset
!

addPushTemporaryVariable: offset
    "Add a bytecode to push a temporary variable with the given offset"
    self addBytecode: BytecodeSet PUSH_TEMPORARY_VARIABLE.
    self addInt32: offset
!

addStoreInstanceVariable: offset
    "Add a bytecode to store into an instance variable with the given offset"
    self addBytecode: BytecodeSet STORE_INSTANCE_VARIABLE.
    self addInt32: offset
!

addStoreTemporaryVariable: offset
    "Add a bytecode to store into a temporary variable with the given offset"
    self addBytecode: BytecodeSet STORE_TEMPORARY_VARIABLE.
    self addInt32: offset
!

addSendMessage: selectorIndex argCount: argCount
    "Add a bytecode to send a message with the given selector index and argument count"
    self addBytecode: BytecodeSet SEND_MESSAGE.
    self addInt32: selectorIndex.
    self addInt32: argCount
!

addJump: target
    "Add a bytecode to jump to the given target"
    self addBytecode: BytecodeSet JUMP.
    self addInt32: target
!

addJumpIfTrue: target
    "Add a bytecode to jump to the given target if the top of the stack is true"
    self addBytecode: BytecodeSet JUMP_IF_TRUE.
    self addInt32: target
!

addJumpIfFalse: target
    "Add a bytecode to jump to the given target if the top of the stack is false"
    self addBytecode: BytecodeSet JUMP_IF_FALSE.
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
    "Ensure the bytecode array is finalized before returning it"
    bytecodePosition > 1 ifTrue: [self finalizeBytecodes].
    ^bytecodes
!

bytecodeAt: anIndex
    "Return the bytecode at the given index"
    ^bytecodes at: anIndex
!

bytecodeUint32At: anIndex
    "Return the 32-bit unsigned integer at the given index"
    ^bytecodes uint32At: anIndex
!

bytecodeSize
    "Return the size of the bytecode array"
    ^bytecodePosition - 1
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

!SmalltalkMethod methodsFor: 'converting'!
asString
    "Return a string representation of the method"
    | stream |
    stream := WriteStream on: String new.
    stream nextPutAll: class name.
    stream nextPutAll: '>>'.
    stream nextPutAll: selector.
    stream nextPutAll: ' ['.
    stream nextPutAll: bytecodes size printString.
    stream nextPutAll: ' bytes, '.
    stream nextPutAll: literals size printString.
    stream nextPutAll: ' literals]'.
    ^stream contents
!

!SmalltalkMethod methodsFor: 'execution'!
valueWithReceiver: anObject arguments: anArray
    "Execute the method with the given receiver and arguments"
    | context interpreter |
    context := SmalltalkContext new
        receiver: anObject;
        arguments: anArray;
        method: self;
        yourself.
        
    interpreter := BytecodeInterpreter new
        method: self bytecodes;
        context: context;
        yourself.
        
    ^interpreter execute
!
