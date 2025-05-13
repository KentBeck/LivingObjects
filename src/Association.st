Object subclass: #Association
    instanceVariableNames: 'key value'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!Association class methodsFor: 'instance creation'!
key: aKey value: aValue
    "Create a new Association with the given key and value"
    ^self basicNew initialize: aKey value: aValue
!

!Association methodsFor: 'initialization'!
initialize: aKey value: aValue
    "Initialize the Association with the given key and value"
    key := aKey.
    value := aValue
!

!Association methodsFor: 'accessing'!
key
    "Return the key of the association"
    ^key
!

key: aKey
    "Set the key of the association"
    key := aKey
!

value
    "Return the value of the association"
    ^value
!

value: aValue
    "Set the value of the association"
    value := aValue
!

!Association methodsFor: 'comparing'!
= anObject
    "Return true if the receiver is equal to anObject"
    (anObject isKindOf: Association) ifFalse: [^false].
    ^key = anObject key and: [value = anObject value]
!

hash
    "Return a hash value for the receiver"
    ^key hash bitXor: value hash
!

!Association methodsFor: 'printing'!
printString
    "Return a string representation of the association"
    | stream |
    stream := WriteStream on: String new.
    stream nextPutAll: key printString.
    stream nextPutAll: '->'.
    stream nextPutAll: value printString.
    ^stream contents
!