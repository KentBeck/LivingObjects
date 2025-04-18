Object subclass: #Main
    instanceVariableNames: ''
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!Main class methodsFor: 'running'!
run
    "Initialize the bytecode set"
    BytecodeSet initialize.
    
    "Run the test examples"
    TestRunner run
!
