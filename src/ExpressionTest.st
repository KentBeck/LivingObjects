TestCase subclass: #ExpressionTest
    instanceVariableNames: ''
    classVariableNames: ''
    package: 'SmalltalkInterpreter-Tests'

!ExpressionTest methodsFor: 'testing - arithmetic'!
testBasicArithmetic
    "Test basic arithmetic operations"
    self assert: 3 + 4 equals: 7.
    self assert: 10 - 3 equals: 7.
    self assert: 6 * 7 equals: 42.
    self assert: 20 / 4 equals: 5
!

testArithmeticPrecedence
    "Test arithmetic operator precedence"
    self assert: 2 + 3 * 4 equals: 14.
    self assert: (2 + 3) * 4 equals: 20.
    self assert: 10 - 2 * 3 equals: 4.
    self assert: (10 - 2) / 4 equals: 2
!

testComplexArithmetic
    "Test complex arithmetic expressions"
    self assert: (2 + 3) * (4 + 1) equals: 25.
    self assert: ((10 + 5) * 2) / 6 equals: 5
!

!ExpressionTest methodsFor: 'testing - comparison'!
testBasicComparison
    "Test basic comparison operations"
    self assert: 3 < 5.
    self assert: 7 > 2.
    self assert: 3 = 3.
    self assert: 4 ~= 5.
    self assert: 4 <= 4.
    self assert: 5 >= 3
!

testComparisonResults
    "Test that comparisons return proper boolean values"
    self assert: (5 < 3) equals: false.
    self assert: (2 > 7) equals: false.
    self assert: (3 = 4) equals: false.
    self assert: (3 < 5) equals: true.
    self assert: (7 > 2) equals: true.
    self assert: (3 = 3) equals: true
!

testComplexComparison
    "Test complex comparison expressions"
    self assert: (3 + 2) < (4 * 2).
    self assert: (10 - 3) > (2 * 3).
    self assert: (6 / 2) = (1 + 2)
!

!ExpressionTest methodsFor: 'testing - literals'!
testBooleanLiterals
    "Test boolean literal values"
    self assert: true equals: true.
    self assert: false equals: false.
    self deny: true = false.
    self deny: false = true
!

testNilLiteral
    "Test nil literal"
    self assert: nil equals: nil
!

testStringLiterals
    "Test string literals"
    self assert: 'hello' equals: 'hello'.
    self assert: 'world' equals: 'world'.
    self deny: 'hello' = 'world'
!

!ExpressionTest methodsFor: 'testing - variables'!
testTemporaryVariables
    "Test temporary variable assignment and access"
    | x |
    x := 42.
    self assert: x equals: 42
!

testMultipleTemporaryVariables
    "Test multiple temporary variables"
    | x y |
    x := 10.
    y := 20.
    self assert: x + y equals: 30
!

testTemporaryVariableAssignmentExpression
    "Test that assignment returns the assigned value"
    | x result |
    result := (x := 5) + 1.
    self assert: result equals: 6.
    self assert: x equals: 5
!

!ExpressionTest methodsFor: 'testing - strings'!
testStringConcatenation
    "Test string concatenation"
    self assert: 'hello' , ' world' equals: 'hello world'.
    self assert: 'a' , 'b' , 'c' equals: 'abc'
!

testStringSize
    "Test string size method"
    self assert: 'hello' size equals: 5.
    self assert: '' size equals: 0.
    self assert: 'a' size equals: 1
!

!ExpressionTest methodsFor: 'testing - collections'!
testArrayCreation
    "Test array creation"
    | array |
    array := Array new: 3.
    self assert: array size equals: 3
!

testArrayLiterals
    "Test array literal syntax"
    | array |
    array := #(1 2 3).
    self assert: array size equals: 3.
    self assert: (array at: 1) equals: 1.
    self assert: (array at: 2) equals: 2.
    self assert: (array at: 3) equals: 3
!

testArrayAccess
    "Test array element access"
    | array |
    array := #(1 2 3).
    self assert: (array at: 2) equals: 2
!

!ExpressionTest methodsFor: 'testing - blocks'!
testSimpleBlockValue
    "Test simple block evaluation"
    | result |
    result := [3 + 4] value.
    self assert: result equals: 7
!

testBlockWithArgument
    "Test block with argument"
    | result |
    result := [:x | x + 1] value: 5.
    self assert: result equals: 6
!

testBlockWithTemporaries
    "Test block with temporary variables"
    | result |
    result := [| x | x := 5. x + 1] value.
    self assert: result equals: 6
!

!ExpressionTest methodsFor: 'testing - symbols'!
testSymbolLiterals
    "Test symbol literals"
    self assert: #abc equals: #abc.
    self deny: #abc = #def
!

testSymbolCreation
    "Test that symbols are properly created"
    | sym |
    sym := #test.
    self assert: sym class name equals: #Symbol
!