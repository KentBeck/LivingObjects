Exception subclass: #TestFailure
    instanceVariableNames: ''
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!TestFailure class methodsFor: 'instance creation'!
signal: aString
    "Signal a test failure with the given message"
    ^self new messageText: aString; signal
!