"True class methods for conditional execution"

True>>ifTrue: trueBlock
    ^trueBlock value

True>>ifFalse: falseBlock
    ^nil

True>>ifTrue: trueBlock ifFalse: falseBlock
    ^trueBlock value

True>>ifFalse: falseBlock ifTrue: trueBlock
    ^trueBlock value

True>>not
    ^false

True>>& aBoolean
    "Eager AND"
    ^aBoolean

True>>| aBoolean
    "Eager OR"
    ^true

True>>and: aBlock
    "Lazy AND"
    ^aBlock value

True>>or: aBlock
    "Lazy OR - block not evaluated"
    ^true

True>>xor: aBoolean
    ^aBoolean not

True>>printOn: aStream
    aStream nextPutAll: 'true'

True>>asString
    ^'true'
