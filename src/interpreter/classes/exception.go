package classes

import (
	"unsafe"

	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/runtime"
)

// Exception represents a Smalltalk exception
type Exception struct {
	core.Object
	MessageText *core.Object
	Tag         *core.Object
}

// NewException creates a new exception object
func NewException(class *core.Object) *core.Object {
	exception := &Exception{
		Object: core.Object{
			TypeField: core.OBJ_EXCEPTION,
		},
		MessageText: core.MakeNilImmediate(),
		Tag:         core.MakeNilImmediate(),
	}

	exception.SetClass(class)
	return ExceptionToObject(exception)
}

// ExceptionToObject converts an Exception to an Object
func ExceptionToObject(e *Exception) *core.Object {
	return (*core.Object)(unsafe.Pointer(e))
}

// ObjectToException converts an Object to an Exception
func ObjectToException(o *core.Object) *Exception {
	return (*Exception)(unsafe.Pointer(o))
}

// String returns a string representation of the exception
func (e *Exception) String() string {
	return "Exception"
}

// GetMessageText returns the message text of the exception
func (e *Exception) GetMessageText() *core.Object {
	return e.MessageText
}

// SetMessageText sets the message text of the exception
func (e *Exception) SetMessageText(messageText *core.Object) {
	e.MessageText = messageText
}

// GetTag returns the tag of the exception
func (e *Exception) GetTag() *core.Object {
	return e.Tag
}

// SetTag sets the tag of the exception
func (e *Exception) SetTag(tag *core.Object) {
	e.Tag = tag
}

// Signal signals the exception
func (e *Exception) Signal() *core.Object {
	// Use the runtime package to signal the exception
	return runtime.SignalException(ExceptionToObject(e))
}
