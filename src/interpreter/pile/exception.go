package pile

import (
	"unsafe"
)

// Exception represents a Smalltalk exception
type Exception struct {
	Object
	MessageText *Object
	Tag         *Object
}

// NewException creates a new exception object
func NewException(class *Object) *Object {
	exception := &Exception{
		Object: Object{
			TypeField: OBJ_EXCEPTION,
		},
		MessageText: MakeNilImmediate(),
		Tag:         MakeNilImmediate(),
	}

	exception.SetClass(class)
	return ExceptionToObject(exception)
}

// ExceptionToObject converts an Exception to an Object
func ExceptionToObject(e *Exception) *Object {
	return (*Object)(unsafe.Pointer(e))
}

// ObjectToException converts an Object to an Exception
func ObjectToException(o *Object) *Exception {
	return (*Exception)(unsafe.Pointer(o))
}

// String returns a string representation of the exception
func (e *Exception) String() string {
	return "Exception"
}

// GetMessageText returns the message text of the exception
func (e *Exception) GetMessageText() *Object {
	return e.MessageText
}

// SetMessageText sets the message text of the exception
func (e *Exception) SetMessageText(messageText *Object) {
	e.MessageText = messageText
}

// GetTag returns the tag of the exception
func (e *Exception) GetTag() *Object {
	return e.Tag
}

// SetTag sets the tag of the exception
func (e *Exception) SetTag(tag *Object) {
	e.Tag = tag
}

// Signal signals the exception
func (e *Exception) Signal() *Object {
	// Use the SignalException function to signal the exception
	return SignalException(ExceptionToObject(e))
}