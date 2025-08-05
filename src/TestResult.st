Object subclass: #TestResult
    instanceVariableNames: 'passed failures errors runCount'
    classVariableNames: ''
    package: 'SmalltalkInterpreter-Tests'

!TestResult methodsFor: 'initialization'!
initialize
    super initialize.
    passed := OrderedCollection new.
    failures := OrderedCollection new.
    errors := OrderedCollection new.
    runCount := 0
!

!TestResult methodsFor: 'accessing'!
runCount
    ^runCount
!

passedCount
    ^passed size
!

failureCount
    ^failures size
!

errorCount
    ^errors size
!

passed
    ^passed
!

failures
    ^failures
!

errors
    ^errors
!

!TestResult methodsFor: 'recording'!
addPass: aTestCase
    passed add: aTestCase.
    runCount := runCount + 1
!

addFailure: aTestCase with: aMessage
    failures add: (Array with: aTestCase with: aMessage).
    runCount := runCount + 1
!

addError: aTestCase with: aMessage
    errors add: (Array with: aTestCase with: aMessage).
    runCount := runCount + 1
!

!TestResult methodsFor: 'testing'!
hasErrors
    ^errors size > 0
!

hasFailures
    ^failures size > 0
!

hasPassed
    ^self hasErrors not and: [self hasFailures not]
!

!TestResult methodsFor: 'printing'!
printOn: aStream
    aStream nextPutAll: runCount printString, ' tests run, '.
    aStream nextPutAll: self passedCount printString, ' passed, '.
    aStream nextPutAll: self failureCount printString, ' failures, '.
    aStream nextPutAll: self errorCount printString, ' errors'.
    
    self hasFailures ifTrue: [
        aStream nextPutAll: String cr, 'FAILURES:', String cr.
        failures do: [:each |
            aStream nextPutAll: (each at: 1) printString, ': ', (each at: 2), String cr
        ]
    ].
    
    self hasErrors ifTrue: [
        aStream nextPutAll: String cr, 'ERRORS:', String cr.
        errors do: [:each |
            aStream nextPutAll: (each at: 1) printString, ': ', (each at: 2), String cr
        ]
    ]
!

printString
    | stream |
    stream := WriteStream on: String new.
    self printOn: stream.
    ^stream contents
!