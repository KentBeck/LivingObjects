Object subclass: #AmanuensisUI
    instanceVariableNames: 'amanuensis actionButton inputField'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!AmanuensisUI methodsFor: 'initialization'!
initialize
    amanuensis := Amanuensis new.
    
    actionButton := ActionButton 
        label: 'Toggle Transcription' 
        action: [self toggleTranscription].
        
    inputField := ''
!

!AmanuensisUI methodsFor: 'actions'!
toggleTranscription
    "Toggle the transcription state"
    amanuensis toggleTranscription.
    
    "Update button label based on transcription state"
    amanuensis isTranscribing
        ifTrue: [actionButton label: 'Stop Transcription']
        ifFalse: [actionButton label: 'Start Transcription']
!

submitText: aString
    "Submit text to be transcribed"
    inputField := aString.
    amanuensis transcribe: aString
!

pressActionButton
    "Simulate pressing the action button"
    actionButton press
!

!AmanuensisUI methodsFor: 'accessing'!
amanuensis
    ^amanuensis
!

actionButton
    ^actionButton
!

!AmanuensisUI class methodsFor: 'instance creation'!
new
    "Create a new AmanuensisUI"
    ^super new initialize
!
