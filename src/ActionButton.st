Object subclass: #ActionButton
    instanceVariableNames: 'label action isPressed'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!ActionButton methodsFor: 'initialization'!
initialize
    isPressed := false
!

label: aString
    label := aString
!

action: aBlock
    action := aBlock
!

!ActionButton methodsFor: 'actions'!
press
    "Simulate pressing the button"
    isPressed := true.
    self performAction.
    isPressed := false
!

performAction
    "Perform the associated action"
    action value
!

!ActionButton class methodsFor: 'instance creation'!
label: aString action: aBlock
    "Create a new button with the given label and action"
    ^self new
        label: aString;
        action: aBlock;
        yourself
!
