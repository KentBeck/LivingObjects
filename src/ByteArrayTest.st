Object subclass: #ByteArrayTest
    instanceVariableNames: ''
    classVariableNames: ''
    package: 'SmalltalkInterpreter-Tests'

!ByteArrayTest methodsFor: 'testing'!
testCreation
    "Test creating a ByteArray"
    | array |
    array := ByteArray new: 10.
    self assert: array size equals: 10.
    
    "Test that all elements are initialized to zero"
    1 to: 10 do: [:i |
        self assert: (array at: i) equals: 0]
!

testAccessing
    "Test accessing elements in a ByteArray"
    | array |
    array := ByteArray new: 5.
    
    "Test setting and getting values"
    array at: 1 put: 42.
    array at: 2 put: 255.
    array at: 3 put: 0.
    array at: 4 put: 127.
    array at: 5 put: 128.
    
    self assert: (array at: 1) equals: 42.
    self assert: (array at: 2) equals: 255.
    self assert: (array at: 3) equals: 0.
    self assert: (array at: 4) equals: 127.
    self assert: (array at: 5) equals: 128
!

testBoundsChecking
    "Test bounds checking in ByteArray"
    | array |
    array := ByteArray new: 3.
    
    "Test accessing out of bounds"
    self should: [array at: 0] raise: Error.
    self should: [array at: 4] raise: Error.
    self should: [array at: 0 put: 1] raise: Error.
    self should: [array at: 4 put: 1] raise: Error.
    
    "Test putting values out of byte range"
    self should: [array at: 1 put: -1] raise: Error.
    self should: [array at: 1 put: 256] raise: Error
!

testUint32Operations
    "Test reading and writing 32-bit integers"
    | array |
    array := ByteArray new: 8.
    
    "Test writing and reading uint32 values"
    array uint32At: 1 put: 42.
    self assert: (array uint32At: 1) equals: 42.
    
    array uint32At: 5 put: 16rFFFFFFFF.
    self assert: (array uint32At: 5) equals: 16rFFFFFFFF.
    
    "Test bounds checking for uint32"
    self should: [array uint32At: 0] raise: Error.
    self should: [array uint32At: 6] raise: Error.
    self should: [array uint32At: 0 put: 1] raise: Error.
    self should: [array uint32At: 6 put: 1] raise: Error
!

testCopying
    "Test copying ByteArrays"
    | array copy partial |
    array := ByteArray new: 5.
    1 to: 5 do: [:i | array at: i put: i].
    
    "Test full copy"
    copy := array copy.
    self assert: copy size equals: 5.
    1 to: 5 do: [:i | self assert: (copy at: i) equals: i].
    
    "Test partial copy"
    partial := array copyFrom: 2 to: 4.
    self assert: partial size equals: 3.
    1 to: 3 do: [:i | self assert: (partial at: i) equals: i + 1]
!

testFromOrderedCollection
    "Test creating a ByteArray from an OrderedCollection"
    | collection array |
    collection := OrderedCollection new.
    collection add: 10; add: 20; add: 30; add: 40; add: 50.
    
    array := ByteArray fromOrderedCollection: collection.
    self assert: array size equals: 5.
    self assert: (array at: 1) equals: 10.
    self assert: (array at: 2) equals: 20.
    self assert: (array at: 3) equals: 30.
    self assert: (array at: 4) equals: 40.
    self assert: (array at: 5) equals: 50
!

testAsOrderedCollection
    "Test converting a ByteArray to an OrderedCollection"
    | array collection |
    array := ByteArray new: 5.
    1 to: 5 do: [:i | array at: i put: i * 10].
    
    collection := array asOrderedCollection.
    self assert: collection size equals: 5.
    self assert: (collection at: 1) equals: 10.
    self assert: (collection at: 2) equals: 20.
    self assert: (collection at: 3) equals: 30.
    self assert: (collection at: 4) equals: 40.
    self assert: (collection at: 5) equals: 50
!

testEquality
    "Test ByteArray equality"
    | array1 array2 array3 |
    array1 := ByteArray new: 3.
    array1 at: 1 put: 1; at: 2 put: 2; at: 3 put: 3.
    
    array2 := ByteArray new: 3.
    array2 at: 1 put: 1; at: 2 put: 2; at: 3 put: 3.
    
    array3 := ByteArray new: 3.
    array3 at: 1 put: 1; at: 2 put: 2; at: 3 put: 4.
    
    self assert: array1 = array2.
    self deny: array1 = array3.
    self deny: array2 = array3
!

!ByteArrayTest methodsFor: 'assertions'!
assert: actual equals: expected
    "Assert that actual equals expected"
    actual = expected ifFalse: [
        self error: 'Expected ', expected printString, ' but got ', actual printString]
!

assert: aBoolean
    "Assert that aBoolean is true"
    aBoolean ifFalse: [self error: 'Assertion failed']
!

deny: aBoolean
    "Assert that aBoolean is false"
    aBoolean ifTrue: [self error: 'Denial failed']
!

should: aBlock raise: anExceptionClass
    "Assert that evaluating aBlock raises an exception of class anExceptionClass"
    | raised exception |
    raised := false.
    
    [aBlock value]
        on: Error
        do: [:ex | 
            raised := true.
            exception := ex].
    
    raised ifFalse: [self error: 'Expected exception but none was raised']
!