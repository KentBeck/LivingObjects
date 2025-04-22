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

func TestImmediateBoolean(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Test true immediate
	// Test that TrueObject is an immediate value
	if !IsImmediate(vm.TrueObject) {
		t.Errorf("Expected TrueObject to be an immediate value")
	}

	// Test that TrueObject is specifically a true immediate
	if !IsTrueImmediate(vm.TrueObject) {
		t.Errorf("Expected TrueObject to be a true immediate value")
	}

	// Test that GetClass returns the correct class for immediate true
	trueClass := vm.GetClass(vm.TrueObject)
	if trueClass != vm.TrueClass {
		t.Errorf("Expected GetClass(TrueObject) to return TrueClass, got %v", trueClass)
	}

	// Test that IsTrue returns true for immediate true
	if !vm.TrueObject.IsTrue() {
		t.Errorf("Expected IsTrue(TrueObject) to return true")
	}

	// Test that String returns "true" for immediate true
	if vm.TrueObject.String() != "true" {
		t.Errorf("Expected String(TrueObject) to return \"true\", got %s", vm.TrueObject.String())
	}

	// Test false immediate
	// Test that FalseObject is an immediate value
	if !IsImmediate(vm.FalseObject) {
		t.Errorf("Expected FalseObject to be an immediate value")
	}

	// Test that FalseObject is specifically a false immediate
	if !IsFalseImmediate(vm.FalseObject) {
		t.Errorf("Expected FalseObject to be a false immediate value")
	}

	// Test that GetClass returns the correct class for immediate false
	falseClass := vm.GetClass(vm.FalseObject)
	if falseClass != vm.FalseClass {
		t.Errorf("Expected GetClass(FalseObject) to return FalseClass, got %v", falseClass)
	}

	// Test that IsTrue returns false for immediate false
	if vm.FalseObject.IsTrue() {
		t.Errorf("Expected IsTrue(FalseObject) to return false")
	}

	// Test that String returns "false" for immediate false
	if vm.FalseObject.String() != "false" {
		t.Errorf("Expected String(FalseObject) to return \"false\", got %s", vm.FalseObject.String())
	}

	// Test that NewBoolean returns the correct immediate values
	trueObj := NewBoolean(true)
	if !IsImmediate(trueObj) {
		t.Errorf("Expected NewBoolean(true) to return an immediate value")
	}
	if !IsTrueImmediate(trueObj) {
		t.Errorf("Expected NewBoolean(true) to return a true immediate value")
	}

	falseObj := NewBoolean(false)
	if !IsImmediate(falseObj) {
		t.Errorf("Expected NewBoolean(false) to return an immediate value")
	}
	if !IsFalseImmediate(falseObj) {
		t.Errorf("Expected NewBoolean(false) to return a false immediate value")
	}
}

func TestMakeImmediateValues(t *testing.T) {
	// Test that MakeNilImmediate returns a value with the correct tag
	nilObj := MakeNilImmediate()
	tag := GetTag(nilObj)
	if tag != TAG_SPECIAL {
		t.Errorf("Expected tag to be TAG_SPECIAL (%d), got %d", TAG_SPECIAL, tag)
	}

	// Test that MakeTrueImmediate returns a value with the correct tag
	trueObj := MakeTrueImmediate()
	tag = GetTag(trueObj)
	if tag != TAG_SPECIAL {
		t.Errorf("Expected tag to be TAG_SPECIAL (%d), got %d", TAG_SPECIAL, tag)
	}

	// Test that MakeFalseImmediate returns a value with the correct tag
	falseObj := MakeFalseImmediate()
	tag = GetTag(falseObj)
	if tag != TAG_SPECIAL {
		t.Errorf("Expected tag to be TAG_SPECIAL (%d), got %d", TAG_SPECIAL, tag)
	}
}
