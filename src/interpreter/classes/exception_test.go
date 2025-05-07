package classes

import (
	"testing"

	"smalltalklsp/interpreter/core"
)

func TestNewException(t *testing.T) {
	// Create a class for the exception
	objectClass := NewClass("Object", nil)
	exceptionClass := NewClass("Exception", objectClass)

	// Create a new exception
	exception := NewException(ClassToObject(exceptionClass))

	// Check that the exception is of the correct type
	if exception.Type() != core.OBJ_EXCEPTION {
		t.Errorf("NewException(class).Type() = %d, want %d", exception.Type(), core.OBJ_EXCEPTION)
	}

	// Convert to Exception
	exceptionObj := ObjectToException(exception)

	// Check that the exception has the correct class
	if exception.Class() != ClassToObject(exceptionClass) {
		t.Errorf("exception.Class() = %v, want %v", exception.Class(), ClassToObject(exceptionClass))
	}

	// Check that the message text is nil
	if !core.IsNilImmediate(exceptionObj.GetMessageText()) {
		t.Errorf("exceptionObj.GetMessageText() = %v, want nil", exceptionObj.GetMessageText())
	}

	// Check that the tag is nil
	if !core.IsNilImmediate(exceptionObj.GetTag()) {
		t.Errorf("exceptionObj.GetTag() = %v, want nil", exceptionObj.GetTag())
	}
}

func TestExceptionAccessors(t *testing.T) {
	// Create a class for the exception
	objectClass := NewClass("Object", nil)
	exceptionClass := NewClass("Exception", objectClass)

	// Create a new exception
	exception := NewException(ClassToObject(exceptionClass))
	exceptionObj := ObjectToException(exception)

	// Create a message text
	messageText := NewString("Test exception")

	// Set the message text
	exceptionObj.SetMessageText(StringToObject(messageText))

	// Check that the message text is set correctly
	if exceptionObj.GetMessageText() != StringToObject(messageText) {
		t.Errorf("exceptionObj.GetMessageText() = %v, want %v", exceptionObj.GetMessageText(), StringToObject(messageText))
	}

	// Create a tag
	tag := NewString("Test tag")

	// Set the tag
	exceptionObj.SetTag(StringToObject(tag))

	// Check that the tag is set correctly
	if exceptionObj.GetTag() != StringToObject(tag) {
		t.Errorf("exceptionObj.GetTag() = %v, want %v", exceptionObj.GetTag(), StringToObject(tag))
	}
}

func TestExceptionToObjectAndBack(t *testing.T) {
	// Create a class for the exception
	objectClass := NewClass("Object", nil)
	exceptionClass := NewClass("Exception", objectClass)

	// Create a new exception
	exception := NewException(ClassToObject(exceptionClass))

	// Convert to Object and back
	obj := ExceptionToObject(ObjectToException(exception))

	// Check that the object is the same
	if obj != exception {
		t.Errorf("ExceptionToObject(ObjectToException(exception)) = %v, want %v", obj, exception)
	}

	// Check that the type is preserved
	if obj.Type() != core.OBJ_EXCEPTION {
		t.Errorf("obj.Type() = %d, want %d", obj.Type(), core.OBJ_EXCEPTION)
	}
}
