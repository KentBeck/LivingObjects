Object subclass: #TestCase
    instanceVariableNames: 'testSelector'
    classVariableNames: ''
    package: 'SmalltalkInterpreter-Tests'

!TestCase class methodsFor: 'instance creation'!
selector: aSymbol
    "Create a test case with the given test selector"
    ^self new testSelector: aSymbol
!

!TestCase methodsFor: 'initialization'!
testSelector: aSymbol
    testSelector := aSymbol
!

!TestCase methodsFor: 'accessing'!
testSelector
    ^testSelector
!

!TestCase methodsFor: 'running'!
setUp
    "Override in subclasses to set up test fixtures"
!

tearDown
    "Override in subclasses to clean up after tests"
!

run
    "Run this test case and return a TestResult"
    | result |
    result := TestResult new.
    self runCase: result.
    ^result
!

runCase: aTestResult
    "Run this test case and record the result"
    [
        self setUp.
        [
            self perform: testSelector.
            aTestResult addPass: self
        ] ensure: [
            self tearDown
        ]
    ] on: TestFailure do: [:failure |
        aTestResult addFailure: self with: failure messageText
    ] on: Error do: [:error |
        aTestResult addError: self with: error messageText
    ]
!

!TestCase methodsFor: 'assertions'!
assert: aBoolean
    "Assert that aBoolean is true"
    aBoolean ifFalse: [
        TestFailure signal: 'Assertion failed'
    ]
!

assert: actual equals: expected
    "Assert that actual equals expected"
    actual = expected ifFalse: [
        TestFailure signal: 'Expected ', expected printString, ' but got ', actual printString
    ]
!

deny: aBoolean
    "Assert that aBoolean is false"
    aBoolean ifTrue: [
        TestFailure signal: 'Denial failed'
    ]
!

should: aBlock raise: anExceptionClass
    "Assert that evaluating aBlock raises an exception of class anExceptionClass"
    | raised |
    raised := false.
    
    [aBlock value]
        on: anExceptionClass
        do: [:ex | raised := true].
    
    raised ifFalse: [
        TestFailure signal: 'Expected ', anExceptionClass name, ' but no exception was raised'
    ]
!

shouldnt: aBlock raise: anExceptionClass
    "Assert that evaluating aBlock does not raise an exception of class anExceptionClass"
    [aBlock value]
        on: anExceptionClass
        do: [:ex | 
            TestFailure signal: 'Did not expect ', anExceptionClass name, ' but it was raised: ', ex messageText
        ]
!

fail
    "Unconditionally fail this test"
    TestFailure signal: 'Test failed'
!

fail: aString
    "Unconditionally fail this test with the given message"
    TestFailure signal: aString
!

!TestCase methodsFor: 'printing'!
printString
    ^self class name, '>>', testSelector printString
!