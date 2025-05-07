package testing

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// RunStringTestsFromFile runs tests from a file
// The file format is:
// ```
// # Comment
// Input code ! Expected result
// ```
func RunStringTestsFromFile(filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create a test runner
	runner := NewStringTestRunner()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	lineNum := 0
	testCases := []StringTestCase{}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the test case
		parts := strings.Split(line, "!")
		if len(parts) != 2 {
			return fmt.Errorf("invalid test case format at line %d: %s", lineNum, line)
		}

		input := strings.TrimSpace(parts[0])
		expected := strings.TrimSpace(parts[1])

		// Add the test case
		testCases = append(testCases, StringTestCase{
			Input:       input,
			Expected:    expected,
			Description: fmt.Sprintf("Line %d", lineNum),
		})
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Run the tests
	runner.RunTests(testCases)

	// Print the results
	runner.PrintResults()

	return nil
}

// RunStringTestsFromString runs tests from a string
// The string format is the same as the file format
func RunStringTestsFromString(input string) error {
	// Create a test runner
	runner := NewStringTestRunner()

	// Split the input into lines
	lines := strings.Split(input, "\n")
	lineNum := 0
	testCases := []StringTestCase{}

	for _, line := range lines {
		lineNum++

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the test case
		parts := strings.Split(line, "!")
		if len(parts) != 2 {
			return fmt.Errorf("invalid test case format at line %d: %s", lineNum, line)
		}

		input := strings.TrimSpace(parts[0])
		expected := strings.TrimSpace(parts[1])

		// Add the test case
		testCases = append(testCases, StringTestCase{
			Input:       input,
			Expected:    expected,
			Description: fmt.Sprintf("Line %d", lineNum),
		})
	}

	// Run the tests
	runner.RunTests(testCases)

	// Print the results
	runner.PrintResults()

	return nil
}

// RunSingleStringTest runs a single test
func RunSingleStringTest(input string, expected string) error {
	// Create a test runner
	runner := NewStringTestRunner()

	// Create a test case
	testCase := StringTestCase{
		Input:    input,
		Expected: expected,
	}

	// Run the test
	result := runner.RunTest(testCase)

	// Print the result
	if result.Passed {
		fmt.Println("Test PASSED")
	} else if result.Error != nil {
		fmt.Printf("Test ERROR: %v\n", result.Error)
	} else {
		fmt.Println("Test FAILED")
		fmt.Printf("  Input:    %s\n", result.TestCase.Input)
		fmt.Printf("  Expected: %s\n", result.TestCase.Expected)
		fmt.Printf("  Actual:   %s\n", result.Actual)
	}

	return nil
}
