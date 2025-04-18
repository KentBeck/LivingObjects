SmalltalkObject subclass: #SmalltalkClass
    instanceVariableNames: 'name superclass methods instanceVariableNames'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!SmalltalkClass methodsFor: 'initialization'!
initialize
    super initialize.
    methods := Dictionary new.
    instanceVariableNames := OrderedCollection new
!

name: aSymbol
    name := aSymbol
!

superclass: aClass
    superclass := aClass
!

addInstanceVariableName: aSymbol
    instanceVariableNames add: aSymbol
!

!SmalltalkClass methodsFor: 'accessing'!
name
    ^name
!

superclass
    ^superclass
!

instanceVariableNames
    ^instanceVariableNames
!

!SmalltalkClass methodsFor: 'methods'!
addMethod: aMethod
    methods at: aMethod selector put: aMethod
!

lookupMethod: aSymbol
    "Look up a method in this class or its superclasses"
    ^methods at: aSymbol ifAbsent: [
        superclass 
            ifNotNil: [superclass lookupMethod: aSymbol]
            ifNil: [nil]
    ]
!

!SmalltalkClass methodsFor: 'instance creation'!
new
    "Create a new instance of this class"
    | instance |
    instance := SmalltalkObject new.
    instance initializeWithClass: self.
    ^instance
!
