package tests

import (
	"bufio"
	"fmt"
	"os"
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/parser"
	"smalltalklsp/interpreter/vm"
	"strings"
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

		// Create a new VM instance for each test
		vmInstance := vm.NewVM()

		// Run the test
		result, err := evaluateExpression(vmInstance, test.Expression)
		if err != nil {
			test.ActualResult = fmt.Sprintf("Error: %v", err)
			test.Passed = false
		} else {
			test.ActualResult = result.String()
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

func evaluateExpression(vmInstance *vm.VM, expression string) (*core.Object, error) {
	// Wrap the expression in a method with a return statement
	wrappedExpression := "^ " + expression

	// Parse the wrapped expression
	parsed, err := parser.NewParser(wrappedExpression, classes.ClassToObject(vmInstance.ObjectClass)).ParseExpression()
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression: %s - %v", expression, err)
	}
	if parsed == nil {
		return nil, fmt.Errorf("failed to parse expression: %s", expression)
	}

	// Compile the parsed expression
	method := compiler.NewBytecodeCompiler(classes.ClassToObject(vmInstance.ObjectClass)).Compile(parsed)
	methodObj := classes.MethodToObject(method)

	// No debug printing

	// Create a context for execution
	context := vm.NewContext(methodObj, classes.ClassToObject(vmInstance.ObjectClass), []*core.Object{}, nil)

	// Execute through VM.Execute()
	result, err := vmInstance.ExecuteContext(context)
	if err != nil {
		return nil, err
	}

	return result.(*core.Object), nil
}
