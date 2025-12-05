"False class methods for conditional execution"

False>>ifTrue: trueBlock
    ^nil

False>>ifFalse: falseBlock
    ^falseBlock value

False>>ifTrue: trueBlock ifFalse: falseBlock
    ^falseBlock value

False>>ifFalse: falseBlock ifTrue: trueBlock
    ^falseBlock value

False>>not
    ^true

False>>& aBoolean
    "Eager AND"
    ^false

False>>| aBoolean
    "Eager OR"
    ^aBoolean

False>>and: aBlock
    "Lazy AND - block not evaluated"
    ^false

False>>or: aBlock
    "Lazy OR"
    ^aBlock value

False>>xor: aBoolean
    ^aBoolean

False>>printOn: aStream
    aStream nextPutAll: 'false'

False>>asString
    ^'false'
