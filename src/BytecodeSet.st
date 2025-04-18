Object subclass: #BytecodeSet
    instanceVariableNames: ''
    classVariableNames: 'PUSH_LITERAL PUSH_INSTANCE_VARIABLE PUSH_TEMPORARY_VARIABLE PUSH_SELF STORE_INSTANCE_VARIABLE STORE_TEMPORARY_VARIABLE SEND_MESSAGE RETURN_STACK_TOP JUMP JUMP_IF_TRUE JUMP_IF_FALSE POP DUPLICATE'
    package: 'SmalltalkInterpreter'

!BytecodeSet class methodsFor: 'initialization'!
initialize
    "Initialize bytecode constants"
    PUSH_LITERAL := 0.            "0: Push a literal from the literals array (followed by 4-byte index)"
    PUSH_INSTANCE_VARIABLE := 1.  "1: Push an instance variable value (followed by 4-byte offset)"
    PUSH_TEMPORARY_VARIABLE := 2. "2: Push a temporary variable value (followed by 4-byte offset)"
    PUSH_SELF := 3.               "3: Push self onto the stack"
    STORE_INSTANCE_VARIABLE := 4. "4: Store a value into an instance variable (followed by 4-byte offset)"
    STORE_TEMPORARY_VARIABLE := 5. "5: Store a value into a temporary variable (followed by 4-byte offset)"
    SEND_MESSAGE := 6.            "6: Send a message (followed by 4-byte selector index and 4-byte arg count)"
    RETURN_STACK_TOP := 7.        "7: Return the value on top of the stack"
    JUMP := 8.                    "8: Jump to a different bytecode (followed by 4-byte target)"
    JUMP_IF_TRUE := 9.            "9: Jump if top of stack is true (followed by 4-byte target)"
    JUMP_IF_FALSE := 10.          "10: Jump if top of stack is false (followed by 4-byte target)"
    POP := 11.                    "11: Pop the top value from the stack"
    DUPLICATE := 12.              "12: Duplicate the top value on the stack"
!

!BytecodeSet class methodsFor: 'accessing'!
isPushLiteral: aBytecode
    ^aBytecode = PUSH_LITERAL
!

isPushInstanceVariable: aBytecode
    ^aBytecode = PUSH_INSTANCE_VARIABLE
!

isPushTemporaryVariable: aBytecode
    ^aBytecode = PUSH_TEMPORARY_VARIABLE
!

isPushSelf: aBytecode
    ^aBytecode = PUSH_SELF
!

isStoreInstanceVariable: aBytecode
    ^aBytecode = STORE_INSTANCE_VARIABLE
!

isStoreTemporaryVariable: aBytecode
    ^aBytecode = STORE_TEMPORARY_VARIABLE
!

isSendMessage: aBytecode
    ^aBytecode = SEND_MESSAGE
!

isReturn: aBytecode
    ^aBytecode = RETURN_STACK_TOP
!

isJump: aBytecode
    ^aBytecode = JUMP or: [aBytecode = JUMP_IF_TRUE or: [aBytecode = JUMP_IF_FALSE]]
!

isPop: aBytecode
    ^aBytecode = POP
!

isDuplicate: aBytecode
    ^aBytecode = DUPLICATE
!

!BytecodeSet class methodsFor: 'instruction size'!
instructionSize: aBytecode
    "Return the size of the instruction in bytes (including the opcode)"
    aBytecode = PUSH_LITERAL ifTrue: [^5]. "1 byte opcode + 4 byte index"
    aBytecode = PUSH_INSTANCE_VARIABLE ifTrue: [^5]. "1 byte opcode + 4 byte offset"
    aBytecode = PUSH_TEMPORARY_VARIABLE ifTrue: [^5]. "1 byte opcode + 4 byte offset"
    aBytecode = PUSH_SELF ifTrue: [^1]. "1 byte opcode"
    aBytecode = STORE_INSTANCE_VARIABLE ifTrue: [^5]. "1 byte opcode + 4 byte offset"
    aBytecode = STORE_TEMPORARY_VARIABLE ifTrue: [^5]. "1 byte opcode + 4 byte offset"
    aBytecode = SEND_MESSAGE ifTrue: [^9]. "1 byte opcode + 4 byte selector index + 4 byte arg count"
    aBytecode = RETURN_STACK_TOP ifTrue: [^1]. "1 byte opcode"
    aBytecode = JUMP ifTrue: [^5]. "1 byte opcode + 4 byte target"
    aBytecode = JUMP_IF_TRUE ifTrue: [^5]. "1 byte opcode + 4 byte target"
    aBytecode = JUMP_IF_FALSE ifTrue: [^5]. "1 byte opcode + 4 byte target"
    aBytecode = POP ifTrue: [^1]. "1 byte opcode"
    aBytecode = DUPLICATE ifTrue: [^1]. "1 byte opcode"

    self error: 'Unknown bytecode: ', aBytecode printString
!
