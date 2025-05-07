package main

import (
	"flag"
	"fmt"
	"os"

	testing "smalltalklsp/interpreter/testing"
)

func main() {
	// Define command-line flags
	fileFlag := flag.String("file", "", "Path to a test file")
	inputFlag := flag.String("input", "", "Smalltalk code to execute")
	expectedFlag := flag.String("expected", "", "Expected result")
	stringFlag := flag.String("string", "", "String containing test cases")

	// Parse the flags
	flag.Parse()

	// Check which mode to run in
	if *fileFlag != "" {
		// Run tests from a file
		err := testing.RunStringTestsFromFile(*fileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else if *inputFlag != "" && *expectedFlag != "" {
		// Run a single test
		err := testing.RunSingleStringTest(*inputFlag, *expectedFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else if *stringFlag != "" {
		// Run tests from a string
		err := testing.RunStringTestsFromString(*stringFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Print usage
		fmt.Println("String->String Testing Framework for Smalltalk")
		fmt.Println("Usage:")
		fmt.Println("  -file <path>: Run tests from a file")
		fmt.Println("  -input <code> -expected <result>: Run a single test")
		fmt.Println("  -string <tests>: Run tests from a string")
		fmt.Println("")
		fmt.Println("File/String Format:")
		fmt.Println("  # Comment")
		fmt.Println("  Input code ! Expected result")
		os.Exit(1)
	}
}
