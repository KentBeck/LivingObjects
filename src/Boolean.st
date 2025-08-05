Object subclass: #Boolean
    instanceVariableNames: ''
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!Boolean methodsFor: 'logical operations'!
& aBoolean
    "Logical AND - both operands are always evaluated"
    ^self and: [aBoolean]
!

| aBoolean
    "Logical OR - both operands are always evaluated"
    ^self or: [aBoolean]
!

not
    "Logical NOT"
    ^self ifTrue: [false] ifFalse: [true]
!

!Boolean methodsFor: 'short-circuit operations'!
and: aBlock
    "Short-circuit AND - aBlock is only evaluated if receiver is true"
    ^self ifTrue: aBlock ifFalse: [false]
!

or: aBlock
    "Short-circuit OR - aBlock is only evaluated if receiver is false"
    ^self ifTrue: [true] ifFalse: aBlock
!

!Boolean methodsFor: 'converting'!
asString
    ^self printString
!

!Boolean methodsFor: 'printing'!
printString
    "Subclasses should override this"
    ^'Boolean'
!