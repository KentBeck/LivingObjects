Object subclass: #TestByteArray
    instanceVariableNames: ''
    classVariableNames: ''
    package: 'SmalltalkInterpreter-Tests'

!TestByteArray class methodsFor: 'running'!
run
    "Run the ByteArray tests"
    | test |
    Transcript show: 'Running ByteArray tests...'; cr.
    
    test := ByteArrayTest new.
    
    self runTest: #testCreation on: test.
    self runTest: #testAccessing on: test.
    self runTest: #testBoundsChecking on: test.
    self runTest: #testUint32Operations on: test.
    self runTest: #testCopying on: test.
    self runTest: #testFromOrderedCollection on: test.
    self runTest: #testAsOrderedCollection on: test.
    self runTest: #testEquality on: test.
    
    Transcript show: 'All ByteArray tests completed successfully!'; cr
!

runTest: selector on: testInstance
    "Run a single test method and report the result"
    | result |
    Transcript show: 'Running ', selector; cr.
    
    [
        result := testInstance perform: selector.
        Transcript show: '  PASSED'; cr
    ] on: Error do: [:ex |
        Transcript show: '  FAILED: ', ex description; cr.
        Transcript show: '  ', ex messageText; cr
    ]
!