Object subclass: #Exception
    instanceVariableNames: 'messageText'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!Exception class methodsFor: 'instance creation'!
signal
    "Signal an exception with no message"
    ^self new signal
!

signal: aString
    "Signal an exception with the given message"
    ^self new messageText: aString; signal
!

!Exception methodsFor: 'initialization'!
initialize
    super initialize.
    messageText := ''
!

messageText: aString
    messageText := aString
!

!Exception methodsFor: 'accessing'!
messageText
    ^messageText
!

!Exception methodsFor: 'signaling'!
signal
    "Signal this exception - for now, just print and halt"
    Transcript show: 'Exception: ', messageText; cr.
    ^self
!

!Exception methodsFor: 'printing'!
printString
    ^'Exception: ', messageText
!