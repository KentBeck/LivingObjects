Object subclass: #SmalltalkObject
    instanceVariableNames: 'class instanceVariables'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!SmalltalkObject methodsFor: 'initialization'!
initialize
    instanceVariables := Dictionary new
!

initializeWithClass: aClass
    self initialize.
    class := aClass
!

!SmalltalkObject methodsFor: 'accessing'!
class
    ^class
!

at: aSymbol
    ^instanceVariables at: aSymbol ifAbsent: [nil]
!

at: aSymbol put: aValue
    ^instanceVariables at: aSymbol put: aValue
!

perform: aSymbol withArguments: anArray
    "Find and execute the method for the given selector"
    | method context interpreter |
    method := self class lookupMethod: aSymbol.
    method ifNil: [self error: 'Method not found: ', aSymbol].
    
    context := SmalltalkContext new
        receiver: self;
        arguments: anArray;
        method: method;
        yourself.
        
    interpreter := BytecodeInterpreter new
        method: method bytecodes;
        context: context;
        yourself.
        
    ^interpreter execute
!

perform: aSymbol
    ^self perform: aSymbol withArguments: #()
!

perform: aSymbol with: arg1
    ^self perform: aSymbol withArguments: {arg1}
!

perform: aSymbol with: arg1 with: arg2
    ^self perform: aSymbol withArguments: {arg1. arg2}
!
