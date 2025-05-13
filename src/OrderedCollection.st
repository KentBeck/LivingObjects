Object subclass: #OrderedCollection
    instanceVariableNames: 'array size'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!OrderedCollection class methodsFor: 'instance creation'!
new
    "Create a new OrderedCollection with default capacity"
    ^self new: 10
!

new: capacity
    "Create a new OrderedCollection with the given capacity"
    ^self basicNew initialize: capacity
!

!OrderedCollection methodsFor: 'initialization'!
initialize: capacity
    "Initialize the OrderedCollection with the given capacity"
    array := Array new: capacity.
    size := 0
!

!OrderedCollection methodsFor: 'accessing'!
size
    "Return the number of elements in the collection"
    ^size
!

at: anIndex
    "Return the element at the given index"
    (anIndex < 1 or: [anIndex > size]) 
        ifTrue: [self error: 'Index out of bounds'].
    ^array at: anIndex
!

at: anIndex put: anObject
    "Set the element at the given index"
    (anIndex < 1 or: [anIndex > size]) 
        ifTrue: [self error: 'Index out of bounds'].
    array at: anIndex put: anObject.
    ^anObject
!

first
    "Return the first element in the collection"
    size = 0 ifTrue: [self error: 'Collection is empty'].
    ^array at: 1
!

last
    "Return the last element in the collection"
    size = 0 ifTrue: [self error: 'Collection is empty'].
    ^array at: size
!

!OrderedCollection methodsFor: 'adding'!
add: anObject
    "Add an object to the end of the collection"
    self ensureCapacity: size + 1.
    size := size + 1.
    array at: size put: anObject.
    ^anObject
!

addFirst: anObject
    "Add an object to the beginning of the collection"
    self ensureCapacity: size + 1.
    size + 1 to: 2 by: -1 do: [:i |
        array at: i put: (array at: i - 1)].
    array at: 1 put: anObject.
    size := size + 1.
    ^anObject
!

addLast: anObject
    "Add an object to the end of the collection"
    ^self add: anObject
!

addAll: aCollection
    "Add all elements from the given collection"
    aCollection do: [:each | self add: each].
    ^aCollection
!

!OrderedCollection methodsFor: 'removing'!
removeAt: anIndex
    "Remove the element at the given index"
    | removed |
    (anIndex < 1 or: [anIndex > size]) 
        ifTrue: [self error: 'Index out of bounds'].
    
    removed := array at: anIndex.
    anIndex to: size - 1 do: [:i |
        array at: i put: (array at: i + 1)].
    array at: size put: nil.
    size := size - 1.
    ^removed
!

removeFirst
    "Remove the first element in the collection"
    size = 0 ifTrue: [self error: 'Collection is empty'].
    ^self removeAt: 1
!

removeLast
    "Remove the last element in the collection"
    size = 0 ifTrue: [self error: 'Collection is empty'].
    ^self removeAt: size
!

remove: anObject ifAbsent: exceptionBlock
    "Remove the given object from the collection"
    1 to: size do: [:i |
        (array at: i) = anObject ifTrue: [
            ^self removeAt: i]].
    ^exceptionBlock value
!

remove: anObject
    "Remove the given object from the collection"
    ^self remove: anObject ifAbsent: [self error: 'Object not found']
!

!OrderedCollection methodsFor: 'private'!
ensureCapacity: requiredCapacity
    "Ensure the collection has enough capacity"
    | newCapacity newArray |
    array size < requiredCapacity ifTrue: [
        newCapacity := array size * 2 max: requiredCapacity.
        newArray := Array new: newCapacity.
        1 to: size do: [:i | newArray at: i put: (array at: i)].
        array := newArray]
!

!OrderedCollection methodsFor: 'enumerating'!
do: aBlock
    "Evaluate the block for each element in the collection"
    1 to: size do: [:i | aBlock value: (array at: i)]
!

collect: aBlock
    "Return a new collection with the results of evaluating the block for each element"
    | newCollection |
    newCollection := self class new: size.
    1 to: size do: [:i |
        newCollection add: (aBlock value: (array at: i))].
    ^newCollection
!

select: aBlock
    "Return a new collection with elements that satisfy the condition"
    | newCollection |
    newCollection := self class new.
    1 to: size do: [:i |
        (aBlock value: (array at: i)) 
            ifTrue: [newCollection add: (array at: i)]].
    ^newCollection
!

reject: aBlock
    "Return a new collection with elements that don't satisfy the condition"
    | newCollection |
    newCollection := self class new.
    1 to: size do: [:i |
        (aBlock value: (array at: i)) 
            ifFalse: [newCollection add: (array at: i)]].
    ^newCollection
!

detect: aBlock ifNone: exceptionBlock
    "Return the first element that satisfies the condition"
    1 to: size do: [:i |
        (aBlock value: (array at: i)) 
            ifTrue: [^array at: i]].
    ^exceptionBlock value
!

detect: aBlock
    "Return the first element that satisfies the condition"
    ^self detect: aBlock ifNone: [self error: 'Element not found']
!

!OrderedCollection methodsFor: 'converting'!
asArray
    "Convert the collection to an Array"
    | newArray |
    newArray := Array new: size.
    1 to: size do: [:i | newArray at: i put: (array at: i)].
    ^newArray
!

asOrderedCollection
    "Return self since we're already an OrderedCollection"
    ^self
!

!OrderedCollection methodsFor: 'printing'!
printString
    "Return a string representation of the collection"
    | stream |
    stream := WriteStream on: String new.
    stream nextPutAll: 'OrderedCollection('.
    1 to: size do: [:i |
        i > 1 ifTrue: [stream nextPutAll: ' '].
        stream nextPutAll: (array at: i) printString].
    stream nextPutAll: ')'.
    ^stream contents
!