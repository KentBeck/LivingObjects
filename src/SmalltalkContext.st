Object subclass: #SmalltalkContext
    instanceVariableNames: 'receiver method arguments temporaryVariables sender pc'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!SmalltalkContext methodsFor: 'initialization'!
initialize
    temporaryVariables := Dictionary new.
    pc := 1
!

receiver: anObject
    receiver := anObject
!

method: aMethod
    method := aMethod
!

arguments: anArray
    arguments := anArray
!

sender: aContext
    sender := aContext
!

!SmalltalkContext methodsFor: 'accessing'!
receiver
    ^receiver
!

method
    ^method
!

arguments
    ^arguments
!

sender
    ^sender
!

pc
    ^pc
!

pc: anInteger
    pc := anInteger
!

temporaryAt: aSymbol
    ^temporaryVariables at: aSymbol ifAbsent: [nil]
!

temporaryAt: aSymbol put: aValue
    ^temporaryVariables at: aSymbol put: aValue
!

temporaryAt: anIndex
    | symbol |
    symbol := method temporaryVariables at: anIndex.
    ^self temporaryAt: symbol
!

temporaryAt: anIndex put: aValue
    | symbol |
    symbol := method temporaryVariables at: anIndex.
    ^self temporaryAt: symbol put: aValue
!

argumentAt: anIndex
    ^arguments at: anIndex
!
