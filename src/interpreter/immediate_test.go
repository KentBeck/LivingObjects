package main

import (
	"testing"
)

func TestImmediateNil(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Test that NilObject is an immediate value
	if !IsImmediate(vm.NilObject) {
		t.Errorf("Expected NilObject to be an immediate value")
	}

	// Test that NilObject is specifically a nil immediate
	if !IsNilImmediate(vm.NilObject) {
		t.Errorf("Expected NilObject to be a nil immediate value")
	}

	// Test that GetClass returns the correct class for immediate nil
	nilClass := vm.GetClass(vm.NilObject)
	if nilClass != vm.NilClass {
		t.Errorf("Expected GetClass(NilObject) to return NilClass, got %v", nilClass)
	}

	// Test that IsTrue returns false for immediate nil
	if vm.NilObject.IsTrue() {
		t.Errorf("Expected IsTrue(NilObject) to return false")
	}

	// Test that String returns "nil" for immediate nil
	if vm.NilObject.String() != "nil" {
		t.Errorf("Expected String(NilObject) to return \"nil\", got %s", vm.NilObject.String())
	}

	// Test that NewNil returns an immediate nil
	nilObj := NewNil()
	if !IsImmediate(nilObj) {
		t.Errorf("Expected NewNil() to return an immediate value")
	}
	if !IsNilImmediate(nilObj) {
		t.Errorf("Expected NewNil() to return a nil immediate value")
	}
}

func TestMakeNilImmediate(t *testing.T) {
	// Test that MakeNilImmediate returns a value with the correct tag
	nilObj := MakeNilImmediate()
	tag := GetTag(nilObj)
	if tag != TAG_SPECIAL {
		t.Errorf("Expected tag to be TAG_SPECIAL (%d), got %d", TAG_SPECIAL, tag)
	}
}
