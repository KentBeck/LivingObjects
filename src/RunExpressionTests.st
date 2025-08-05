Object subclass: #RunExpressionTests
    instanceVariableNames: ''
    classVariableNames: ''
    package: 'SmalltalkInterpreter-Tests'

!RunExpressionTests class methodsFor: 'running'!
run
    "Run all expression tests and demonstrate SUnit framework"
    | suite result |
    
    Transcript show: '=== SUnit Expression Tests ==='; cr; cr.
    
    "Create a test suite for expression tests"
    suite := TestSuite named: 'Expression Tests'.
    
    "Add individual test methods"
    suite addTest: (ExpressionTest selector: #testBasicArithmetic).
    suite addTest: (ExpressionTest selector: #testArithmeticPrecedence).
    suite addTest: (ExpressionTest selector: #testComplexArithmetic).
    suite addTest: (ExpressionTest selector: #testBasicComparison).
    suite addTest: (ExpressionTest selector: #testComparisonResults).
    suite addTest: (ExpressionTest selector: #testComplexComparison).
    suite addTest: (ExpressionTest selector: #testBooleanLiterals).
    suite addTest: (ExpressionTest selector: #testNilLiteral).
    suite addTest: (ExpressionTest selector: #testStringLiterals).
    suite addTest: (ExpressionTest selector: #testTemporaryVariables).
    suite addTest: (ExpressionTest selector: #testMultipleTemporaryVariables).
    suite addTest: (ExpressionTest selector: #testTemporaryVariableAssignmentExpression).
    suite addTest: (ExpressionTest selector: #testStringConcatenation).
    suite addTest: (ExpressionTest selector: #testStringSize).
    suite addTest: (ExpressionTest selector: #testArrayCreation).
    suite addTest: (ExpressionTest selector: #testArrayLiterals).
    suite addTest: (ExpressionTest selector: #testArrayAccess).
    suite addTest: (ExpressionTest selector: #testSimpleBlockValue).
    suite addTest: (ExpressionTest selector: #testBlockWithArgument).
    suite addTest: (ExpressionTest selector: #testBlockWithTemporaries).
    suite addTest: (ExpressionTest selector: #testSymbolLiterals).
    suite addTest: (ExpressionTest selector: #testSymbolCreation).
    
    "Run the test suite"
    result := TestRunner runSuite: suite.
    
    Transcript cr; show: '=== Test Summary ==='; cr.
    Transcript show: 'Total tests: ', suite size printString; cr.
    Transcript show: 'Passed: ', result passedCount printString; cr.
    Transcript show: 'Failed: ', result failureCount printString; cr.
    Transcript show: 'Errors: ', result errorCount printString; cr.
    
    result hasPassed
        ifTrue: [Transcript show: 'ALL TESTS PASSED!'; cr]
        ifFalse: [Transcript show: 'SOME TESTS FAILED - see details above'; cr].
    
    ^result
!

runSingle: testSelector
    "Run a single test method"
    | test result |
    test := ExpressionTest selector: testSelector.
    result := TestRunner run: test.
    ^result
!

runArithmeticTests
    "Run just the arithmetic tests"
    | suite |
    suite := TestSuite named: 'Arithmetic Tests'.
    suite addTest: (ExpressionTest selector: #testBasicArithmetic).
    suite addTest: (ExpressionTest selector: #testArithmeticPrecedence).
    suite addTest: (ExpressionTest selector: #testComplexArithmetic).
    ^TestRunner runSuite: suite
!

runComparisonTests
    "Run just the comparison tests"
    | suite |
    suite := TestSuite named: 'Comparison Tests'.
    suite addTest: (ExpressionTest selector: #testBasicComparison).
    suite addTest: (ExpressionTest selector: #testComparisonResults).
    suite addTest: (ExpressionTest selector: #testComplexComparison).
    ^TestRunner runSuite: suite
!

runStringTests
    "Run just the string tests"
    | suite |
    suite := TestSuite named: 'String Tests'.
    suite addTest: (ExpressionTest selector: #testStringLiterals).
    suite addTest: (ExpressionTest selector: #testStringConcatenation).
    suite addTest: (ExpressionTest selector: #testStringSize).
    ^TestRunner runSuite: suite
!

runBlockTests
    "Run just the block tests"
    | suite |
    suite := TestSuite named: 'Block Tests'.
    suite addTest: (ExpressionTest selector: #testSimpleBlockValue).
    suite addTest: (ExpressionTest selector: #testBlockWithArgument).
    suite addTest: (ExpressionTest selector: #testBlockWithTemporaries).
    ^TestRunner runSuite: suite
!