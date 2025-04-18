Object subclass: #Examples
    instanceVariableNames: ''
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!Examples class methodsFor: 'examples'!
factorial
    "Example factorial method"
    ^ '
    factorial: n
        n = 0 ifTrue: [^1].
        n = 1 ifTrue: [^1].
        ^n * (self factorial: n - 1)
    '
!

fibonacci
    "Example fibonacci method"
    ^ '
    fibonacci: n
        n <= 1 ifTrue: [^n].
        ^(self fibonacci: n - 1) + (self fibonacci: n - 2)
    '
!

counter
    "Example counter class"
    ^ '
    Object subclass: #Counter
        instanceVariableNames: ''count''
        classVariableNames: ''''
        package: ''Examples''
        
    !Counter methodsFor: ''initialization''!
    initialize
        count := 0
    !
    
    !Counter methodsFor: ''accessing''!
    count
        ^count
    !
    
    increment
        count := count + 1.
        ^count
    !
    
    decrement
        count := count - 1.
        ^count
    !
    '
!

helloWorld
    "Example hello world method"
    ^ '
    helloWorld
        ^''Hello, World!''
    '
!

sumTo
    "Example sum to n method"
    ^ '
    sumTo: n
        | sum |
        sum := 0.
        1 to: n do: [:i | sum := sum + i].
        ^sum
    '
!
