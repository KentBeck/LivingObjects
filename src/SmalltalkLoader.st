Object subclass: #SmalltalkLoader
    instanceVariableNames: 'loadedFiles'
    classVariableNames: ''
    package: 'SmalltalkInterpreter'

!SmalltalkLoader class methodsFor: 'instance creation'!
new
    "Create a new SmalltalkLoader"
    ^self basicNew initialize
!

!SmalltalkLoader methodsFor: 'initialization'!
initialize
    "Initialize the loader"
    loadedFiles := OrderedCollection new
!

!SmalltalkLoader methodsFor: 'loading'!
loadFile: filename
    "Load a single Smalltalk file"
    | content |
    Transcript show: 'Loading file: ', filename; cr.
    
    "Check if the file has already been loaded"
    (loadedFiles includes: filename) ifTrue: [
        Transcript show: '  Already loaded, skipping'; cr.
        ^self].
    
    "Read the file content"
    content := self readFile: filename.
    
    "Compile and evaluate the file content"
    self evaluateFileContent: content.
    
    "Mark the file as loaded"
    loadedFiles add: filename.
    
    Transcript show: '  Successfully loaded: ', filename; cr
!

loadDirectory: directoryPath
    "Load all Smalltalk files in a directory"
    | files |
    Transcript show: 'Loading directory: ', directoryPath; cr.
    
    "Get all .st files in the directory"
    files := self listFilesInDirectory: directoryPath.
    
    "Load each file"
    files do: [:filename |
        (filename endsWith: '.st') ifTrue: [
            self loadFile: directoryPath, '/', filename]]
!

loadFiles: filenameArray
    "Load multiple Smalltalk files in the given order"
    filenameArray do: [:filename |
        self loadFile: filename]
!

!SmalltalkLoader methodsFor: 'file operations'!
readFile: filename
    "Read the content of a file"
    "This is a primitive operation that will be implemented in Go"
    self primitiveFailed
!

listFilesInDirectory: directoryPath
    "List all files in a directory"
    "This is a primitive operation that will be implemented in Go"
    self primitiveFailed
!

!SmalltalkLoader methodsFor: 'evaluation'!
evaluateFileContent: content
    "Evaluate the content of a Smalltalk file"
    "This will be implemented to parse and evaluate Smalltalk code"
    self primitiveFailed
!

!SmalltalkLoader methodsFor: 'accessing'!
loadedFiles
    "Return the list of loaded files"
    ^loadedFiles
!

!SmalltalkLoader methodsFor: 'testing'!
isFileLoaded: filename
    "Check if a file has been loaded"
    ^loadedFiles includes: filename
!