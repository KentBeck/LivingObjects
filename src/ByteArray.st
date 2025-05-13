Object subclass: #ByteArray
    instanceVariableNames: 'bytes size'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!ByteArray class methodsFor: 'instance creation'!
new: anInteger
    "Create a new ByteArray with the given size"
    ^self basicNew initialize: anInteger
!

fromOrderedCollection: aCollection
    "Create a ByteArray from an OrderedCollection of bytes"
    | array |
    array := self new: aCollection size.
    1 to: aCollection size do: [:i |
        array at: i put: (aCollection at: i)].
    ^array
!

!ByteArray methodsFor: 'initialization'!
initialize: anInteger
    "Initialize the ByteArray with the given size"
    size := anInteger.
    bytes := Array new: size.
    1 to: size do: [:i | bytes at: i put: 0]
!

!ByteArray methodsFor: 'accessing'!
size
    "Return the size of the ByteArray"
    ^size
!

at: anIndex
    "Return the byte at the given index"
    (anIndex < 1 or: [anIndex > size]) 
        ifTrue: [self error: 'Index out of bounds'].
    ^bytes at: anIndex
!

at: anIndex put: aByte
    "Set the byte at the given index"
    (anIndex < 1 or: [anIndex > size]) 
        ifTrue: [self error: 'Index out of bounds'].
    (aByte < 0 or: [aByte > 255]) 
        ifTrue: [self error: 'Value out of byte range (0-255)'].
    bytes at: anIndex put: aByte.
    ^aByte
!

!ByteArray methodsFor: 'converting'!
asOrderedCollection
    "Convert the ByteArray to an OrderedCollection"
    | collection |
    collection := OrderedCollection new: size.
    1 to: size do: [:i |
        collection add: (bytes at: i)].
    ^collection
!

!ByteArray methodsFor: 'printing'!
printString
    "Return a string representation of the ByteArray"
    | stream |
    stream := WriteStream on: String new.
    stream nextPutAll: 'ByteArray('.
    stream nextPutAll: size printString.
    stream nextPutAll: ')'.
    ^stream contents
!

!ByteArray methodsFor: 'integer operations'!
uint32At: anIndex
    "Read a 32-bit unsigned integer from the ByteArray at the given index"
    | value |
    (anIndex < 1 or: [anIndex + 3 > size]) 
        ifTrue: [self error: 'Index out of bounds for uint32'].
    
    value := 0.
    value := value + ((self at: anIndex) bitShift: 0).
    value := value + ((self at: anIndex + 1) bitShift: 8).
    value := value + ((self at: anIndex + 2) bitShift: 16).
    value := value + ((self at: anIndex + 3) bitShift: 24).
    ^value
!

uint32At: anIndex put: anInteger
    "Write a 32-bit unsigned integer to the ByteArray at the given index"
    (anIndex < 1 or: [anIndex + 3 > size]) 
        ifTrue: [self error: 'Index out of bounds for uint32'].
    
    self at: anIndex put: (anInteger bitAnd: 16rFF).
    self at: anIndex + 1 put: ((anInteger bitShift: -8) bitAnd: 16rFF).
    self at: anIndex + 2 put: ((anInteger bitShift: -16) bitAnd: 16rFF).
    self at: anIndex + 3 put: ((anInteger bitShift: -24) bitAnd: 16rFF).
    ^anInteger
!

!ByteArray methodsFor: 'copying'!
copyFrom: startIndex to: endIndex
    "Return a new ByteArray containing the bytes from startIndex to endIndex"
    | newSize newArray |
    (startIndex < 1 or: [startIndex > size]) 
        ifTrue: [self error: 'Start index out of bounds'].
    (endIndex < startIndex or: [endIndex > size]) 
        ifTrue: [self error: 'End index out of bounds'].
    
    newSize := endIndex - startIndex + 1.
    newArray := ByteArray new: newSize.
    
    1 to: newSize do: [:i |
        newArray at: i put: (self at: startIndex + i - 1)].
    
    ^newArray
!

copy
    "Return a copy of the ByteArray"
    ^self copyFrom: 1 to: size
!

!ByteArray methodsFor: 'comparing'!
= anObject
    "Return true if the receiver is equal to anObject"
    (anObject isKindOf: ByteArray) ifFalse: [^false].
    self size = anObject size ifFalse: [^false].
    
    1 to: size do: [:i |
        (self at: i) = (anObject at: i) ifFalse: [^false]].
    
    ^true
!

hash
    "Return a hash value for the receiver"
    | hash |
    hash := size hash.
    1 to: (size min: 10) do: [:i |
        hash := hash bitXor: (self at: i) hash].
    ^hash
!