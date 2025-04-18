Object subclass: #TestRunner
    instanceVariableNames: ''
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!TestRunner class methodsFor: 'running'!
run
    "Run the test examples"
    self testHelloWorld.
    self testFactorial.
    self testCounter
!

testHelloWorld
    "Test the hello world example"
    | source class method compiler result |
    
    Transcript show: 'Testing Hello World...'; cr.
    
    "Create a class for our test"
    class := SmalltalkClass new
        name: #TestHelloWorld;
        yourself.
        
    "Create a method"
    method := SmalltalkMethod new
        selector: #helloWorld;
        class: class;
        yourself.
        
    "Compile the method"
    source := Examples helloWorld.
    compiler := SmalltalkCompiler new
        source: source;
        method: method;
        yourself.
        
    method := compiler compile.
    class addMethod: method.
    
    "Create an instance and run the method"
    result := class new perform: #helloWorld.
    
    Transcript show: 'Result: ', result printString; cr
!

testFactorial
    "Test the factorial example"
    | source class method compiler result |
    
    Transcript show: 'Testing Factorial...'; cr.
    
    "Create a class for our test"
    class := SmalltalkClass new
        name: #TestFactorial;
        yourself.
        
    "Create a method"
    method := SmalltalkMethod new
        selector: #factorial:;
        class: class;
        yourself.
        
    "Compile the method"
    source := Examples factorial.
    compiler := SmalltalkCompiler new
        source: source;
        method: method;
        yourself.
        
    method := compiler compile.
    class addMethod: method.
    
    "Create an instance and run the method"
    result := class new perform: #factorial: with: 5.
    
    Transcript show: 'Factorial of 5: ', result printString; cr
!

testCounter
    "Test the counter example"
    | source class method compiler counter result |
    
    Transcript show: 'Testing Counter...'; cr.
    
    "Create a Counter class"
    class := SmalltalkClass new
        name: #Counter;
        yourself.
    class addInstanceVariableName: #count.
    
    "Create initialize method"
    method := SmalltalkMethod new
        selector: #initialize;
        class: class;
        yourself.
        
    "Compile the initialize method"
    source := 'initialize
        count := 0'.
    compiler := SmalltalkCompiler new
        source: source;
        method: method;
        yourself.
        
    method := compiler compile.
    class addMethod: method.
    
    "Create increment method"
    method := SmalltalkMethod new
        selector: #increment;
        class: class;
        yourself.
        
    "Compile the increment method"
    source := 'increment
        count := count + 1.
        ^count'.
    compiler := SmalltalkCompiler new
        source: source;
        method: method;
        yourself.
        
    method := compiler compile.
    class addMethod: method.
    
    "Create count accessor method"
    method := SmalltalkMethod new
        selector: #count;
        class: class;
        yourself.
        
    "Compile the count method"
    source := 'count
        ^count'.
    compiler := SmalltalkCompiler new
        source: source;
        method: method;
        yourself.
        
    method := compiler compile.
    class addMethod: method.
    
    "Create an instance and test the counter"
    counter := class new.
    counter perform: #initialize.
    
    result := counter perform: #count.
    Transcript show: 'Initial count: ', result printString; cr.
    
    result := counter perform: #increment.
    Transcript show: 'After increment: ', result printString; cr.
    
    result := counter perform: #increment.
    Transcript show: 'After another increment: ', result printString; cr
!
