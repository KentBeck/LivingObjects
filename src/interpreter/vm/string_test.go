package vm_test

import (
	"testing"

	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

func TestVMNewString(t *testing.T) {
	// Create a VM for testing
	virtualMachine := vm.NewVM()

	tests := []struct {
		name  string
		value string
	}{
		{"Empty string", ""},
		{"Simple string", "hello"},
		{"String with spaces", "hello world"},
		{"String with special chars", "hello\nworld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a string using the VM's NewString method
			strObj := virtualMachine.NewString(tt.value)

			// Check that the object is not nil
			if strObj == nil {
				t.Errorf("NewString(%q) returned nil", tt.value)
				return
			}

			// Check that the object has the correct type
			if strObj.Type() != pile.OBJ_STRING {
				t.Errorf("NewString(%q).Type() = %d, want %d", tt.value, strObj.Type(), pile.OBJ_STRING)
			}

			// Check that the object has the correct class
			class := virtualMachine.GetClass(strObj)
			if class != pile.ObjectToClass(virtualMachine.Globals["String"]) {
				t.Errorf("NewString(%q) has class %v, want %v", tt.value, class, pile.ObjectToClass(virtualMachine.Globals["String"]))
			}

			// Check that the object has the correct value
			str := pile.ObjectToString(strObj)
			if str.GetValue() != tt.value {
				t.Errorf("NewString(%q).GetValue() = %q, want %q", tt.value, str.GetValue(), tt.value)
			}
		})
	}
}