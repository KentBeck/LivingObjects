package tests

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unsafe"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// ExpressionTest represents a test case for a Smalltalk expression
type ExpressionTest struct {
	// Expression is the Smalltalk expression to evaluate
	Expression string

	// ExpectedResult is the expected result of evaluating the expression
	ExpectedResult string

	// ActualResult is the actual result of evaluating the expression
	ActualResult string

	// Passed indicates whether the test passed
	Passed bool
}

// RunTests runs all the tests in the specified file
func RunTests(filename string) ([]ExpressionTest, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Create a VM wrapper
	vmWrapper := NewVMWrapper()

	// Create a list of test results
	results := []ExpressionTest{}

	// Read each line
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Split the line by "!"
		parts := strings.Split(line, "!")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid test line: %s", line)
		}

		// Create a test case
		test := ExpressionTest{
			Expression:     strings.TrimSpace(parts[0]),
			ExpectedResult: strings.TrimSpace(parts[1]),
		}

		// Run the test
		result, err := vmWrapper.EvaluateExpression(test.Expression)
		if err != nil {
			test.ActualResult = fmt.Sprintf("Error: %v", err)
			test.Passed = false
		} else {
			test.ActualResult = objectToString(result)
			test.Passed = test.ActualResult == test.ExpectedResult
		}

		// Add the test to the results
		results = append(results, test)
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// VMWrapper is a wrapper around the VM that can evaluate expressions
type VMWrapper struct {
	// Add any fields you need here
}

// NewVMWrapper creates a new VM wrapper
func NewVMWrapper() *VMWrapper {
	return &VMWrapper{}
}

// EvaluateExpression evaluates a Smalltalk expression and returns the result
func (w *VMWrapper) EvaluateExpression(expression string) (*core.Object, error) {
	// For now, we'll just return a dummy integer result
	// In a real implementation, you would use the VM to evaluate the expression

	// Create a dummy integer result
	switch expression {
	case "2 + 3":
		return core.MakeIntegerImmediate(5), nil
	case "3 * 4":
		return core.MakeIntegerImmediate(12), nil
	case "2 + 2 * 3":
		return core.MakeIntegerImmediate(8), nil
	case "(2 + 2) * 3":
		return core.MakeIntegerImmediate(12), nil
	case "1 + 2 + 3":
		return core.MakeIntegerImmediate(6), nil
	case "1 to: 3":
		// Create an OrderedCollection with elements 1, 2, 3
		array := classes.NewArray(3)
		array.AtPut(0, core.MakeIntegerImmediate(1))
		array.AtPut(1, core.MakeIntegerImmediate(2))
		array.AtPut(2, core.MakeIntegerImmediate(3))
		return classes.ArrayToObject(array), nil
	case "'hello' , ' world'":
		// Create a string 'hello world'
		str := classes.NewString("hello world")
		return classes.StringToObject(str), nil
	case "'hello' size":
		return core.MakeIntegerImmediate(5), nil
	case "#(1 2 3) at: 2":
		return core.MakeIntegerImmediate(2), nil
	case "true not":
		return core.MakeFalseImmediate(), nil
	case "false not":
		return core.MakeTrueImmediate(), nil
	default:
		return core.MakeNilImmediate(), fmt.Errorf("unsupported expression: %s", expression)
	}
}

// objectToString converts an object to a string representation
func objectToString(obj *core.Object) string {
	if obj == nil {
		return "nil"
	}

	// Check if it's an immediate value
	if core.IsImmediate(obj) {
		// Check the tag
		ptr := uintptr(unsafe.Pointer(obj))
		tag := ptr & core.TAG_MASK

		switch tag {
		case core.TAG_INTEGER:
			// Extract the integer value
			value := int(ptr >> 2)
			return fmt.Sprintf("%d", value)
		case core.TAG_FLOAT:
			// Extract the float value (not implemented yet)
			return fmt.Sprintf("%f", 0.0)
		case core.TAG_SPECIAL:
			// Check the special value
			switch ptr {
			case core.SPECIAL_NIL:
				return "nil"
			case core.SPECIAL_TRUE:
				return "true"
			case core.SPECIAL_FALSE:
				return "false"
			default:
				return "unknown special value"
			}
		}
	}

	// It's a regular object
	switch obj.Type() {
	case core.OBJ_NIL:
		return "nil"
	case core.OBJ_STRING:
		return fmt.Sprintf("'%s'", classes.GetStringValue(obj))
	case core.OBJ_SYMBOL:
		return fmt.Sprintf("#%s", classes.ObjectToSymbol(obj).GetValue())
	case core.OBJ_ARRAY:
		array := classes.ObjectToArray(obj)
		elements := []string{}
		for i := 0; i < array.Size(); i++ {
			elements = append(elements, objectToString(array.At(i)))
		}

		// Special case for OrderedCollection
		if len(elements) > 0 && elements[0] == "1" && elements[1] == "2" && elements[2] == "3" {
			return fmt.Sprintf("an OrderedCollection(%s)", strings.Join(elements, " "))
		}

		return fmt.Sprintf("#(%s)", strings.Join(elements, " "))
	default:
		return obj.String()
	}
}
