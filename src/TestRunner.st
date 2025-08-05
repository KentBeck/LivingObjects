Object subclass: #TestRunner
    instanceVariableNames: ''
    classVariableNames: ''
    package: 'SmalltalkInterpreter-Tests'

!TestRunner class methodsFor: 'running'!
run: aTestCase
    "Run a single test case and print the results"
    | result |
    result := aTestCase run.
    self printResult: result.
    ^result
!

runSuite: aTestSuite
    "Run a test suite and print the results"
    | result |
    result := aTestSuite run.
    self printResult: result.
    ^result
!

runClass: aTestClass
    "Run all tests in a test class"
    | suite |
    suite := TestSuite named: aTestClass name.
    suite addTestsFor: aTestClass.
    ^self runSuite: suite
!

printResult: aTestResult
    "Print the test result to the transcript"
    Transcript show: aTestResult printString; cr.
    
    aTestResult hasPassed 
        ifTrue: [Transcript show: 'ALL TESTS PASSED'; cr]
        ifFalse: [Transcript show: 'SOME TESTS FAILED'; cr]
!

!TestRunner methodsFor: 'running'!
run: aTestCase
    "Instance method to run a test case"
    ^self class run: aTestCase
!

runSuite: aTestSuite
    "Instance method to run a test suite"
    ^self class runSuite: aTestSuite
!