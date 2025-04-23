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

func TestImmediateInteger(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Test small integer
	intObj := vm.NewInteger(42)

	// Test that it's an immediate value
	if !IsImmediate(intObj) {
		t.Errorf("Expected NewInteger(42) to return an immediate value")
	}

	// Test that it's specifically an integer immediate
	if !IsIntegerImmediate(intObj) {
		t.Errorf("Expected NewInteger(42) to return an integer immediate value")
	}

	// Test that GetClass returns the correct class for immediate integer
	intClass := vm.GetClass(intObj)
	if intClass != vm.IntegerClass {
		t.Errorf("Expected GetClass(intObj) to return IntegerClass, got %v", intClass)
	}

	// Test that IsTrue returns true for immediate integer
	if !intObj.IsTrue() {
		t.Errorf("Expected IsTrue(intObj) to return true")
	}

	// Test that String returns the correct string for immediate integer
	if intObj.String() != "42" {
		t.Errorf("Expected String(intObj) to return \"42\", got %s", intObj.String())
	}

	// Test that GetIntegerImmediate returns the correct value
	value := GetIntegerImmediate(intObj)
	if value != 42 {
		t.Errorf("Expected GetIntegerImmediate(intObj) to return 42, got %d", value)
	}

	// Test negative integer
	negIntObj := vm.NewInteger(-42)

	// Test that it's an immediate value
	if !IsImmediate(negIntObj) {
		t.Errorf("Expected NewInteger(-42) to return an immediate value")
	}

	// Test that it's specifically an integer immediate
	if !IsIntegerImmediate(negIntObj) {
		t.Errorf("Expected NewInteger(-42) to return an integer immediate value")
	}

	// Test that GetIntegerImmediate returns the correct value for negative integer
	negValue := GetIntegerImmediate(negIntObj)
	if negValue != -42 {
		t.Errorf("Expected GetIntegerImmediate(negIntObj) to return -42, got %d", negValue)
	}

	// Skip testing large integers since they now panic
	// We're now using only immediate integers
}

func TestImmediateIntegerPrimitives(t *testing.T) {
	// Create a VM
	vm := NewVM()

	// Create immediate integers
	int1 := vm.NewInteger(42)
	int2 := vm.NewInteger(10)

	// Test addition
	plusSelector := NewSymbol("+")
	plusMethod := vm.lookupMethod(vm.IntegerClass, plusSelector)
	result := vm.executePrimitive(int1, plusSelector, []*Object{int2}, plusMethod)

	// Check that the result is an immediate integer
	if !IsIntegerImmediate(result) {
		t.Errorf("Expected result of addition to be an immediate integer")
	}

	// Check that the value is correct
	value := GetIntegerImmediate(result)
	if value != 52 {
		t.Errorf("Expected result of 42 + 10 to be 52, got %d", value)
	}

	// Test subtraction
	minusSelector := NewSymbol("-")
	minusMethod := vm.lookupMethod(vm.IntegerClass, minusSelector)
	result = vm.executePrimitive(int1, minusSelector, []*Object{int2}, minusMethod)

	// Check that the result is an immediate integer
	if !IsIntegerImmediate(result) {
		t.Errorf("Expected result of subtraction to be an immediate integer")
	}

	// Check that the value is correct
	value = GetIntegerImmediate(result)
	if value != 32 {
		t.Errorf("Expected result of 42 - 10 to be 32, got %d", value)
	}

	// Test multiplication
	timesSelector := NewSymbol("*")
	timesMethod := vm.lookupMethod(vm.IntegerClass, timesSelector)
	result = vm.executePrimitive(int1, timesSelector, []*Object{int2}, timesMethod)

	// Check that the result is an immediate integer
	if !IsIntegerImmediate(result) {
		t.Errorf("Expected result of multiplication to be an immediate integer")
	}

	// Check that the value is correct
	value = GetIntegerImmediate(result)
	if value != 420 {
		t.Errorf("Expected result of 42 * 10 to be 420, got %d", value)
	}

	// Test equality
	equalsSelector := NewSymbol("=")
	equalsMethod := vm.lookupMethod(vm.IntegerClass, equalsSelector)
	result = vm.executePrimitive(int1, equalsSelector, []*Object{int1}, equalsMethod)

	// Check that the result is a boolean
	if !IsTrueImmediate(result) {
		t.Errorf("Expected result of equality to be true")
	}

	// Test inequality
	result = vm.executePrimitive(int1, equalsSelector, []*Object{int2}, equalsMethod)

	// Check that the result is a boolean
	if !IsFalseImmediate(result) {
		t.Errorf("Expected result of inequality to be false")
	}

	// Test less than
	lessSelector := NewSymbol("<")
	lessMethod := vm.lookupMethod(vm.IntegerClass, lessSelector)
	result = vm.executePrimitive(int2, lessSelector, []*Object{int1}, lessMethod)

	// Check that the result is a boolean
	if !IsTrueImmediate(result) {
		t.Errorf("Expected result of less than to be true")
	}

	// Test greater than
	greaterSelector := NewSymbol(">")
	greaterMethod := vm.lookupMethod(vm.IntegerClass, greaterSelector)
	result = vm.executePrimitive(int1, greaterSelector, []*Object{int2}, greaterMethod)

	// Check that the result is a boolean
	if !IsTrueImmediate(result) {
		t.Errorf("Expected result of greater than to be true")
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

	// Test that MakeIntegerImmediate returns a value with the correct tag
	intObj := MakeIntegerImmediate(42)
	tag = GetTag(intObj)
	if tag != TAG_INTEGER {
		t.Errorf("Expected tag to be TAG_INTEGER (%d), got %d", TAG_INTEGER, tag)
	}
}
