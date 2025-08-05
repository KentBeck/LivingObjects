Object subclass: #TestSuite
    instanceVariableNames: 'tests name'
    classVariableNames: ''
    package: 'SmalltalkInterpreter-Tests'

!TestSuite class methodsFor: 'instance creation'!
new
    ^self basicNew initialize
!

named: aString
    ^self new name: aString
!

!TestSuite methodsFor: 'initialization'!
initialize
    super initialize.
    tests := OrderedCollection new.
    name := 'TestSuite'
!

name: aString
    name := aString
!

!TestSuite methodsFor: 'accessing'!
name
    ^name
!

tests
    ^tests
!

size
    ^tests size
!

!TestSuite methodsFor: 'adding'!
addTest: aTestCase
    tests add: aTestCase
!

addTestsFor: aTestClass
    "Add all test methods from the given test class"
    | selectors |
    selectors := aTestClass selectors select: [:each | each beginsWith: 'test'].
    selectors do: [:selector |
        self addTest: (aTestClass selector: selector)
    ]
!

!TestSuite methodsFor: 'running'!
run
    "Run all tests in this suite and return a TestResult"
    | result |
    result := TestResult new.
    tests do: [:test |
        test runCase: result
    ].
    ^result
!

!TestSuite methodsFor: 'printing'!
printString
    ^name, ' (', tests size printString, ' tests)'
!