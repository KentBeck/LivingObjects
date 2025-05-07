package tests

import (
	"path/filepath"
	"testing"
)

// TestExpressions runs all the expression tests from string_tests.txt
func TestExpressions(t *testing.T) {
	// Get the path to the test file
	testFile := filepath.Join(".", "string_tests.txt")

	// Run the tests
	results, err := RunTests(testFile)
	if err != nil {
		t.Fatalf("Error running expression tests: %v", err)
	}

	// Check the results
	for i, result := range results {
		t.Run(result.Expression, func(t *testing.T) {
			if !result.Passed {
				t.Errorf("Test %d failed: %s\nExpected: %s\nActual: %s",
					i+1, result.Expression, result.ExpectedResult, result.ActualResult)
			}
		})
	}
}
