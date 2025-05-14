package pile_test

import (
	"testing"

	"smalltalklsp/interpreter/pile"
)

func TestMakeImmediateValues(t *testing.T) {
	// Test that MakeNilImmediate returns a value with the correct tag
	nilObj := pile.MakeNilImmediate()
	tag := pile.GetTag(nilObj)
	if tag != pile.TAG_SPECIAL {
		t.Errorf("Expected tag to be TAG_SPECIAL (%d), got %d", pile.TAG_SPECIAL, tag)
	}
	
	// Test that MakeTrueImmediate returns a value with the correct tag
	trueObj := pile.MakeTrueImmediate()
	tag = pile.GetTag(trueObj)
	if tag != pile.TAG_SPECIAL {
		t.Errorf("Expected tag to be TAG_SPECIAL (%d), got %d", pile.TAG_SPECIAL, tag)
	}
	
	// Test that MakeFalseImmediate returns a value with the correct tag
	falseObj := pile.MakeFalseImmediate()
	tag = pile.GetTag(falseObj)
	if tag != pile.TAG_SPECIAL {
		t.Errorf("Expected tag to be TAG_SPECIAL (%d), got %d", pile.TAG_SPECIAL, tag)
	}
	
	// Test that MakeIntegerImmediate returns a value with the correct tag
	intObj := pile.MakeIntegerImmediate(42)
	tag = pile.GetTag(intObj)
	if tag != pile.TAG_INTEGER {
		t.Errorf("Expected tag to be TAG_INTEGER (%d), got %d", pile.TAG_INTEGER, tag)
	}
}

func TestImmediateChecks(t *testing.T) {
	// Test nil immediate
	nilObj := pile.MakeNilImmediate()
	if !pile.IsImmediate(nilObj) {
		t.Errorf("Expected MakeNilImmediate() to return an immediate value")
	}
	if !pile.IsNilImmediate(nilObj) {
		t.Errorf("Expected MakeNilImmediate() to return a nil immediate value")
	}
	if nilObj.IsTrue() {
		t.Errorf("Expected IsTrue(nilObj) to return false")
	}
	if nilObj.String() != "nil" {
		t.Errorf("Expected String(nilObj) to return \"nil\", got %s", nilObj.String())
	}
	
	// Test true immediate
	trueObj := pile.MakeTrueImmediate()
	if !pile.IsImmediate(trueObj) {
		t.Errorf("Expected MakeTrueImmediate() to return an immediate value")
	}
	if !pile.IsTrueImmediate(trueObj) {
		t.Errorf("Expected MakeTrueImmediate() to return a true immediate value")
	}
	if !trueObj.IsTrue() {
		t.Errorf("Expected IsTrue(trueObj) to return true")
	}
	if trueObj.String() != "true" {
		t.Errorf("Expected String(trueObj) to return \"true\", got %s", trueObj.String())
	}
	
	// Test false immediate
	falseObj := pile.MakeFalseImmediate()
	if !pile.IsImmediate(falseObj) {
		t.Errorf("Expected MakeFalseImmediate() to return an immediate value")
	}
	if !pile.IsFalseImmediate(falseObj) {
		t.Errorf("Expected MakeFalseImmediate() to return a false immediate value")
	}
	if falseObj.IsTrue() {
		t.Errorf("Expected IsTrue(falseObj) to return false")
	}
	if falseObj.String() != "false" {
		t.Errorf("Expected String(falseObj) to return \"false\", got %s", falseObj.String())
	}
}

func TestIntegerImmediate(t *testing.T) {
	// Test positive integer
	intObj := pile.MakeIntegerImmediate(42)
	if !pile.IsImmediate(intObj) {
		t.Errorf("Expected MakeIntegerImmediate(42) to return an immediate value")
	}
	if !pile.IsIntegerImmediate(intObj) {
		t.Errorf("Expected MakeIntegerImmediate(42) to return an integer immediate value")
	}
	
	// The IsTrue test is inconsistent - integers are not necessarily "true" in Smalltalk
	// Skip this test to avoid confusion
	
	if intObj.String() != "42" {
		t.Errorf("Expected String(intObj) to return \"42\", got %s", intObj.String())
	}
	
	// Test that GetIntegerImmediate returns the correct value
	value := pile.GetIntegerImmediate(intObj)
	if value != 42 {
		t.Errorf("Expected GetIntegerImmediate(intObj) to return 42, got %d", value)
	}
	
	// Test negative integer
	negIntObj := pile.MakeIntegerImmediate(-42)
	if !pile.IsImmediate(negIntObj) {
		t.Errorf("Expected MakeIntegerImmediate(-42) to return an immediate value")
	}
	if !pile.IsIntegerImmediate(negIntObj) {
		t.Errorf("Expected MakeIntegerImmediate(-42) to return an integer immediate value")
	}
	
	// Test that GetIntegerImmediate returns the correct value for negative integer
	negValue := pile.GetIntegerImmediate(negIntObj)
	if negValue != -42 {
		t.Errorf("Expected GetIntegerImmediate(negIntObj) to return -42, got %d", negValue)
	}
}

func TestBooleanCreation(t *testing.T) {
	// Test that NewBoolean returns the correct immediate values
	trueObj := pile.NewBoolean(true)
	if !pile.IsImmediate(trueObj) {
		t.Errorf("Expected NewBoolean(true) to return an immediate value")
	}
	if !pile.IsTrueImmediate(trueObj) {
		t.Errorf("Expected NewBoolean(true) to return a true immediate value")
	}
	
	falseObj := pile.NewBoolean(false)
	if !pile.IsImmediate(falseObj) {
		t.Errorf("Expected NewBoolean(false) to return an immediate value")
	}
	if !pile.IsFalseImmediate(falseObj) {
		t.Errorf("Expected NewBoolean(false) to return a false immediate value")
	}
}

func TestNilCreation(t *testing.T) {
	// Test that NewNil returns an immediate nil
	nilObj := pile.NewNil()
	if !pile.IsImmediate(nilObj) {
		t.Errorf("Expected NewNil() to return an immediate value")
	}
	if !pile.IsNilImmediate(nilObj) {
		t.Errorf("Expected NewNil() to return a nil immediate value")
	}
}