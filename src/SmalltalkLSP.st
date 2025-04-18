Object subclass: #SmalltalkLSP
    instanceVariableNames: 'documents workspaceFolders'
    classVariableNames: ''
    package: 'SmalltalkLSP'

!SmalltalkLSP methodsFor: 'initialization'!
initialize
    documents := Dictionary new.
    workspaceFolders := OrderedCollection new
!

!SmalltalkLSP methodsFor: 'lsp-handlers'!
textDocumentDidOpen: params
    | uri text |
    uri := params at: 'textDocument' at: 'uri'.
    text := params at: 'textDocument' at: 'text'.
    documents at: uri put: text
!

textDocumentDidChange: params
    | uri changes |
    uri := params at: 'textDocument' at: 'uri'.
    changes := params at: 'contentChanges'.
    self updateDocument: uri withChanges: changes
!

completion: params
    | uri position |
    uri := params at: 'textDocument' at: 'uri'.
    position := params at: 'position'.
    ^self computeCompletionsAt: position inDocument: uri
!

!SmalltalkLSP methodsFor: 'private'!
computeCompletionsAt: position inDocument: uri
    | text completions |
    text := documents at: uri.
    completions := OrderedCollection new.
    "Add completion logic here using the bytecode interpreter for analysis"
    ^completions