Object subclass: #Dictionary
    instanceVariableNames: 'keys values size'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!Dictionary class methodsFor: 'instance creation'!
new
    "Create a new Dictionary with default capacity"
    ^self new: 10
!

new: capacity
    "Create a new Dictionary with the given capacity"
    ^self basicNew initialize: capacity
!

!Dictionary methodsFor: 'initialization'!
initialize: capacity
    "Initialize the Dictionary with the given capacity"
    keys := Array new: capacity.
    values := Array new: capacity.
    size := 0
!

!Dictionary methodsFor: 'accessing'!
size
    "Return the number of key-value pairs in the dictionary"
    ^size
!

at: aKey
    "Return the value associated with the key"
    ^self at: aKey ifAbsent: [self error: 'Key not found']
!

at: aKey ifAbsent: aBlock
    "Return the value associated with the key, or evaluate aBlock if the key is not found"
    | index |
    index := self findKey: aKey.
    index = 0 ifTrue: [^aBlock value].
    ^values at: index
!

at: aKey put: aValue
    "Associate aValue with aKey in the dictionary"
    | index |
    index := self findKey: aKey.
    index = 0
        ifTrue: [
            "Key not found, add a new entry"
            self ensureCapacity: size + 1.
            keys at: size + 1 put: aKey.
            values at: size + 1 put: aValue.
            size := size + 1]
        ifFalse: [
            "Key found, update the value"
            values at: index put: aValue].
    ^aValue
!

includesKey: aKey
    "Return true if the dictionary includes the key"
    ^(self findKey: aKey) > 0
!

!Dictionary methodsFor: 'removing'!
removeKey: aKey
    "Remove the key-value pair with the given key"
    ^self removeKey: aKey ifAbsent: [self error: 'Key not found']
!

removeKey: aKey ifAbsent: aBlock
    "Remove the key-value pair with the given key, or evaluate aBlock if the key is not found"
    | index value |
    index := self findKey: aKey.
    index = 0 ifTrue: [^aBlock value].
    
    value := values at: index.
    
    "Shift all entries after this one"
    index to: size - 1 do: [:i |
        keys at: i put: (keys at: i + 1).
        values at: i put: (values at: i + 1)].
    
    "Clear the last entry"
    keys at: size put: nil.
    values at: size put: nil.
    size := size - 1.
    
    ^value
!

!Dictionary methodsFor: 'private'!
findKey: aKey
    "Find the index of the key in the dictionary. Return 0 if not found."
    1 to: size do: [:i |
        (keys at: i) = aKey ifTrue: [^i]].
    ^0
!

ensureCapacity: requiredCapacity
    "Ensure the dictionary has enough capacity"
    | newCapacity newKeys newValues |
    keys size < requiredCapacity ifTrue: [
        newCapacity := keys size * 2 max: requiredCapacity.
        newKeys := Array new: newCapacity.
        newValues := Array new: newCapacity.
        
        1 to: size do: [:i |
            newKeys at: i put: (keys at: i).
            newValues at: i put: (values at: i)].
        
        keys := newKeys.
        values := newValues]
!

!Dictionary methodsFor: 'enumerating'!
keysAndValuesDo: aBlock
    "Evaluate aBlock with each key and value in the dictionary"
    1 to: size do: [:i |
        aBlock value: (keys at: i) value: (values at: i)]
!

keysDo: aBlock
    "Evaluate aBlock with each key in the dictionary"
    1 to: size do: [:i |
        aBlock value: (keys at: i)]
!

valuesDo: aBlock
    "Evaluate aBlock with each value in the dictionary"
    1 to: size do: [:i |
        aBlock value: (values at: i)]
!

do: aBlock
    "Evaluate aBlock with each value in the dictionary"
    self valuesDo: aBlock
!

collect: aBlock
    "Return a new dictionary with the results of evaluating aBlock for each value"
    | newDict |
    newDict := Dictionary new: size.
    self keysAndValuesDo: [:key :value |
        newDict at: key put: (aBlock value: value)].
    ^newDict
!

select: aBlock
    "Return a new dictionary with entries whose values satisfy aBlock"
    | newDict |
    newDict := Dictionary new.
    self keysAndValuesDo: [:key :value |
        (aBlock value: value) 
            ifTrue: [newDict at: key put: value]].
    ^newDict
!

reject: aBlock
    "Return a new dictionary with entries whose values don't satisfy aBlock"
    | newDict |
    newDict := Dictionary new.
    self keysAndValuesDo: [:key :value |
        (aBlock value: value) 
            ifFalse: [newDict at: key put: value]].
    ^newDict
!

!Dictionary methodsFor: 'accessing'!
keys
    "Return an OrderedCollection of all keys in the dictionary"
    | keyCollection |
    keyCollection := OrderedCollection new: size.
    self keysDo: [:key | keyCollection add: key].
    ^keyCollection
!

values
    "Return an OrderedCollection of all values in the dictionary"
    | valueCollection |
    valueCollection := OrderedCollection new: size.
    self valuesDo: [:value | valueCollection add: value].
    ^valueCollection
!

associations
    "Return an OrderedCollection of Associations for all entries in the dictionary"
    | assocCollection |
    assocCollection := OrderedCollection new: size.
    self keysAndValuesDo: [:key :value |
        assocCollection add: (Association key: key value: value)].
    ^assocCollection
!

!Dictionary methodsFor: 'printing'!
printString
    "Return a string representation of the dictionary"
    | stream |
    stream := WriteStream on: String new.
    stream nextPutAll: 'Dictionary('.
    
    self keysAndValuesDo: [:key :value |
        stream nextPutAll: key printString.
        stream nextPutAll: '->'.
        stream nextPutAll: value printString.
        stream nextPutAll: ' '].
    
    size > 0 ifTrue: [stream skip: -1]. "Remove the trailing space"
    stream nextPutAll: ')'.
    ^stream contents
!