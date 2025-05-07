package main

import (
	"fmt"
	"os"

	"smalltalklsp/interpreter/tests"
)

func main() {
	// Check if a filename was provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <test_file>")
		os.Exit(1)
	}

	// Get the filename
	filename := os.Args[1]

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: File %s does not exist\n", filename)
		os.Exit(1)
	}

	// Run the tests
	fmt.Printf("Running tests from %s\n", filename)
	results, err := tests.RunTests(filename)
	if err != nil {
		fmt.Printf("Error running tests: %v\n", err)
		os.Exit(1)
	}

	// Print the results
	fmt.Printf("\nResults:\n")
	passed := 0
	for i, result := range results {
		fmt.Printf("Test %d: %s\n", i+1, result.Expression)
		fmt.Printf("  Expected: %s\n", result.ExpectedResult)
		fmt.Printf("  Actual:   %s\n", result.ActualResult)
		if result.Passed {
			fmt.Printf("  Result:   PASS\n")
			passed++
		} else {
			fmt.Printf("  Result:   FAIL\n")
		}
		fmt.Println()
	}

	// Print the summary
	fmt.Printf("Summary: %d/%d tests passed\n", passed, len(results))
	if passed < len(results) {
		os.Exit(1)
	}
}
