Object subclass: #WriteStream
    instanceVariableNames: 'collection position'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!WriteStream class methodsFor: 'instance creation'!
on: aCollection
    "Create a new WriteStream on the given collection"
    ^self basicNew initialize: aCollection
!

!WriteStream methodsFor: 'initialization'!
initialize: aCollection
    "Initialize the stream with the given collection"
    collection := aCollection.
    position := 0
!

!WriteStream methodsFor: 'accessing'!
contents
    "Return the contents of the stream"
    ^collection copyFrom: 1 to: position
!

position
    "Return the current position in the stream"
    ^position
!

position: anInteger
    "Set the current position in the stream"
    (anInteger < 0 or: [anInteger > collection size])
        ifTrue: [self error: 'Position out of bounds'].
    position := anInteger
!

!WriteStream methodsFor: 'writing'!
nextPut: anObject
    "Append anObject to the stream"
    position := position + 1.
    position > collection size
        ifTrue: [self growCollection].
    collection at: position put: anObject.
    ^anObject
!

nextPutAll: aCollection
    "Append all elements of aCollection to the stream"
    aCollection do: [:each | self nextPut: each].
    ^aCollection
!

!WriteStream methodsFor: 'private'!
growCollection
    "Grow the collection to accommodate more elements"
    | newCollection |
    newCollection := collection class new: collection size * 2.
    1 to: collection size do: [:i |
        newCollection at: i put: (collection at: i)].
    collection := newCollection
!

!WriteStream methodsFor: 'positioning'!
reset
    "Reset the stream position to the beginning"
    position := 0
!

skip: anInteger
    "Skip forward or backward by anInteger positions"
    self position: position + anInteger
!

!WriteStream methodsFor: 'testing'!
isEmpty
    "Return true if the stream is empty"
    ^position = 0
!

atEnd
    "Return true if the stream is at the end"
    ^position = collection size
!