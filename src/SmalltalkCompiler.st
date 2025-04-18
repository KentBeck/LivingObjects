Object subclass: #SmalltalkCompiler
    instanceVariableNames: 'source tokens currentToken position method'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!SmalltalkCompiler methodsFor: 'initialization'!
initialize
    tokens := OrderedCollection new.
    position := 1
!

source: aString
    source := aString.
    self tokenize
!

method: aMethod
    method := aMethod
!

!SmalltalkCompiler methodsFor: 'tokenizing'!
tokenize
    "A very simple tokenizer for demonstration purposes"
    | i char token inString |
    tokens := OrderedCollection new.
    i := 1.
    token := ''.
    inString := false.

    [i <= source size] whileTrue: [
        char := source at: i.

        inString
            ifTrue: [
                token := token, char asString.
                char = $' ifTrue: [
                    (i < source size and: [(source at: i+1) = $'])
                        ifTrue: [
                            token := token, (source at: i+1) asString.
                            i := i + 1]
                        ifFalse: [inString := false]].
            ]
            ifFalse: [
                (char isSeparator or: [char = $. or: [char = $[ or: [char = $] or: [char = $: or: [char = $^]]]]])
                    ifTrue: [
                        token isEmpty ifFalse: [
                            tokens add: token.
                            token := ''].
                        (char isSeparator not and: [char ~= $.]) ifTrue: [
                            tokens add: char asString].
                    ]
                    ifFalse: [
                        char = $'
                            ifTrue: [
                                token isEmpty ifFalse: [
                                    tokens add: token.
                                    token := ''].
                                token := char asString.
                                inString := true]
                            ifFalse: [token := token, char asString]
                    ]
            ].
        i := i + 1.
    ].

    token isEmpty ifFalse: [tokens add: token].
    position := 1
!

!SmalltalkCompiler methodsFor: 'parsing'!
compile
    "Compile the source code into bytecodes"
    | newMethod |
    newMethod := method.
    self parseMethod.
    ^newMethod
!

parseMethod
    "Parse a method definition"
    | selector body |
    selector := self parseMethodHeader.
    method selector: selector.
    body := self parseMethodBody.
    ^method
!

parseMethodHeader
    "Parse the method header to get the selector"
    | selector |
    "This is a simplified version that just returns the first token as the selector"
    selector := tokens at: position.
    position := position + 1.
    ^selector asSymbol
!

parseMethodBody
    "Parse the method body"
    | statements |
    "Skip any temporary variable declarations for simplicity"
    statements := self parseStatements.
    ^statements
!

parseStatements
    "Parse a sequence of statements"
    | statements |
    statements := OrderedCollection new.

    [position <= tokens size] whileTrue: [
        statements add: self parseStatement.
        (position <= tokens size and: [(tokens at: position) = '.'])
            ifTrue: [position := position + 1]
            ifFalse: [^statements]
    ].

    ^statements
!

parseStatement
    "Parse a single statement"
    (position <= tokens size and: [(tokens at: position) = '^'])
        ifTrue: [
            position := position + 1.
            ^self parseReturn].

    ^self parseExpression
!

parseExpression
    "Parse an expression"
    | receiver messages |
    receiver := self parsePrimary.

    "Parse message sends"
    [(position <= tokens size) and: [(tokens at: position) ~= '.' and: [(tokens at: position) ~= ']']]]
        whileTrue: [
            receiver := self parseMessageSend: receiver].

    ^receiver
!

parsePrimary
    "Parse a primary expression (literal, variable, block, or nested expression)"
    | token |
    token := tokens at: position.
    position := position + 1.

    "Handle literals, variables, blocks, and nested expressions"
    (token first isDigit)
        ifTrue: [
            "Number literal"
            | value index |
            value := token asNumber.
            index := method addLiteral: value.
            method addBytecode: BytecodeSet PUSH_LITERAL.
            method addInt32: index.
            ^value].

    (token first = $')
        ifTrue: [
            "String literal"
            | value index |
            value := token copyFrom: 2 to: token size - 1.
            index := method addLiteral: value.
            method addBytecode: BytecodeSet PUSH_LITERAL.
            method addInt32: index.
            ^value].

    token = 'self'
        ifTrue: [
            method addBytecode: BytecodeSet PUSH_SELF.
            ^token].

    token = 'true'
        ifTrue: [
            | index |
            index := method addLiteral: true.
            method addBytecode: BytecodeSet PUSH_LITERAL.
            method addInt32: index.
            ^true].

    token = 'false'
        ifTrue: [
            | index |
            index := method addLiteral: false.
            method addBytecode: BytecodeSet PUSH_LITERAL.
            method addInt32: index.
            ^false].

    token = 'nil'
        ifTrue: [
            | index |
            index := method addLiteral: nil.
            method addBytecode: BytecodeSet PUSH_LITERAL.
            method addInt32: index.
            ^nil].

    token = '['
        ifTrue: [
            "Block"
            | block |
            block := self parseBlock.
            ^block].

    token = '('
        ifTrue: [
            "Nested expression"
            | expr |
            expr := self parseExpression.
            (tokens at: position) = ')' ifTrue: [position := position + 1].
            ^expr].

    "Assume it's a variable reference"
    | index |
    index := method addLiteral: token asSymbol.
    method addBytecode: BytecodeSet PUSH_LITERAL.
    method addInt32: index.
    ^token
!

parseBlock
    "Parse a block"
    | statements |
    statements := self parseStatements.
    (tokens at: position) = ']' ifTrue: [position := position + 1].
    ^statements
!

parseMessageSend: receiver
    "Parse a message send to the receiver"
    | selector args index |
    selector := tokens at: position.
    position := position + 1.

    args := OrderedCollection new.

    "Parse arguments if any"
    [(position <= tokens size) and: [(tokens at: position) ~= '.' and: [(tokens at: position) ~= ']']]]
        whileTrue: [
            args add: self parseExpression].

    "Add bytecodes for the message send"
    index := method addLiteral: selector asSymbol.
    method addBytecode: BytecodeSet PUSH_LITERAL.
    method addInt32: index.
    method addSendMessage: index argCount: args size.

    ^selector
!

parseReturn
    "Parse a return statement"
    | expr |
    expr := self parseExpression.
    method addBytecode: BytecodeSet RETURN_STACK_TOP.
    ^expr
!
