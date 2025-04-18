Object subclass: #Amanuensis
    instanceVariableNames: 'isTranscribing transcriptionBuffer'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!Amanuensis methodsFor: 'initialization'!
initialize
    isTranscribing := false.
    transcriptionBuffer := OrderedCollection new
!

!Amanuensis methodsFor: 'transcription'!
toggleTranscription
    "Toggle the transcription state"
    isTranscribing := isTranscribing not.
    
    isTranscribing 
        ifTrue: [
            Transcript show: 'Transcription started...'; cr.
            self startTranscribing]
        ifFalse: [
            self stopTranscribing.
            Transcript show: 'Transcription stopped.'; cr.
            Transcript show: 'Transcribed content:'; cr.
            Transcript show: self getTranscribedContent; cr]
!

startTranscribing
    "Start the transcription process"
    transcriptionBuffer := OrderedCollection new
!

stopTranscribing
    "Stop the transcription process"
    "Nothing special to do here, just stop adding to the buffer"
!

isTranscribing
    "Return whether transcription is active"
    ^isTranscribing
!

transcribe: aString
    "Add text to the transcription buffer if transcribing is active"
    isTranscribing ifTrue: [
        transcriptionBuffer add: aString.
        Transcript show: 'Transcribed: ', aString; cr]
!

getTranscribedContent
    "Return the transcribed content as a string"
    | content |
    content := ''.
    transcriptionBuffer do: [:line | 
        content := content, line, String cr].
    ^content
!

!Amanuensis class methodsFor: 'instance creation'!
new
    "Create a new instance of Amanuensis"
    ^super new initialize
!
