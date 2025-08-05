Object subclass: #SUnitDemo
    instanceVariableNames: ''
    classVariableNames: ''
    package: 'SmalltalkInterpreter-Tests'

!SUnitDemo class methodsFor: 'demonstration'!
demo
    "Demonstrate the SUnit testing framework"
    
    Transcript show: '=== SUnit Framework Demonstration ==='; cr; cr.
    
    "1. Demonstrate basic test case creation and running"
    self demoBasicTestCase.
    
    "2. Demonstrate test assertions"
    self demoAssertions.
    
    "3. Demonstrate test suites"
    self demoTestSuite.
    
    "4. Demonstrate expression tests"
    self demoExpressionTests.
    
    Transcript show: '=== Demo Complete ==='; cr
!

demoBasicTestCase
    "Demonstrate creating and running a basic test case"
    | test result |
    
    Transcript show: '1. Basic Test Case:'; cr.
    
    "Create a simple test that should pass"
    test := TestCase selector: #testBasicArithmetic.
    result := test run.
    
    Transcript show: '   Created test: ', test printString; cr.
    Transcript show: '   Result: ', result printString; cr; cr
!

demoAssertions
    "Demonstrate different types of assertions"
    | test |
    
    Transcript show: '2. Test Assertions:'; cr.
    
    "Create a test case to demonstrate assertions"
    test := TestCase new.
    
    "Demonstrate successful assertions"
    Transcript show: '   Testing assert: true -> ';
    [test assert: true. Transcript show: 'PASS'] 
        on: TestFailure do: [:ex | Transcript show: 'FAIL: ', ex messageText].
    Transcript cr.
    
    Transcript show: '   Testing assert:equals: -> ';
    [test assert: 2 + 2 equals: 4. Transcript show: 'PASS'] 
        on: TestFailure do: [:ex | Transcript show: 'FAIL: ', ex messageText].
    Transcript cr.
    
    Transcript show: '   Testing deny: false -> ';
    [test deny: false. Transcript show: 'PASS'] 
        on: TestFailure do: [:ex | Transcript show: 'FAIL: ', ex messageText].
    Transcript cr.
    
    "Demonstrate failing assertion"
    Transcript show: '   Testing assert: false -> ';
    [test assert: false. Transcript show: 'PASS'] 
        on: TestFailure do: [:ex | Transcript show: 'FAIL: ', ex messageText].
    Transcript cr; cr
!

demoTestSuite
    "Demonstrate creating and running test suites"
    | suite result |
    
    Transcript show: '3. Test Suite:'; cr.
    
    "Create a test suite with multiple tests"
    suite := TestSuite named: 'Demo Suite'.
    suite addTest: (ExpressionTest selector: #testBasicArithmetic).
    suite addTest: (ExpressionTest selector: #testBooleanLiterals).
    suite addTest: (ExpressionTest selector: #testStringLiterals).
    
    Transcript show: '   Created suite: ', suite printString; cr.
    
    "Run the suite"
    result := suite run.
    Transcript show: '   Suite result: ', result printString; cr; cr
!

demoExpressionTests
    "Demonstrate running expression tests"
    
    Transcript show: '4. Expression Tests:'; cr.
    
    "Run a subset of expression tests"
    Transcript show: '   Running arithmetic tests...'; cr.
    RunExpressionTests runArithmeticTests.
    
    Transcript show: '   Running comparison tests...'; cr.
    RunExpressionTests runComparisonTests.
    
    Transcript cr
!