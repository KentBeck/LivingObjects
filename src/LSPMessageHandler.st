Object subclass: #LSPMessageHandler
    instanceVariableNames: 'lspServer'
    classVariableNames: ''
    package: 'SmalltalkLSP'

!LSPMessageHandler methodsFor: 'initialization'!
initialize
    lspServer := SmalltalkLSP new
!

!LSPMessageHandler methodsFor: 'processing'!
handleMessage: aJsonString
    | message method params |
    message := Json readFrom: aJsonString readStream.
    method := message at: 'method'.
    params := message at: 'params'.
    
    method = 'textDocument/didOpen' 
        ifTrue: [^lspServer textDocumentDidOpen: params].
    method = 'textDocument/didChange'
        ifTrue: [^lspServer textDocumentDidChange: params].
    method = 'textDocument/completion'
        ifTrue: [^lspServer completion: params]