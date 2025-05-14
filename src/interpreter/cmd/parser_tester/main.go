package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/parser"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

// SimplifiedVM implements the minimal VM interface needed for parsing
type SimplifiedVM struct {
	vm *vm.VM
}

func NewSimplifiedVM() *SimplifiedVM {
	return &SimplifiedVM{
		vm: vm.NewVM(),
	}
}

func (s *SimplifiedVM) NewInteger(value int64) *pile.Object {
	return s.vm.NewInteger(value)
}

func (s *SimplifiedVM) NewString(value string) *pile.Object {
	return s.vm.NewString(value)
}

func (s *SimplifiedVM) NewArray(size int) *pile.Object {
	return s.vm.NewArray(size)
}

// parseCode parses a string as either a method or an expression
func parseCode(input string, methodMode bool) (ast.Node, error) {
	// Create a VM for the parser
	virtualMachine := NewSimplifiedVM()

	// Create a dummy class for the method
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass // Set class's class to itself for proper checks
	classObj := pile.ClassToObject(objectClass)

	// Create a parser
	p := parser.NewParser(input, classObj, virtualMachine)

	// Parse based on whether we're in method or expression mode
	if methodMode {
		return p.Parse()
	} else {
		return p.ParseExpression()
	}
}

// prettyJSON converts a JSON string to a formatted, indented JSON string
func prettyJSON(input string) string {
	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return fmt.Sprintf("Error formatting JSON: %v\nOriginal JSON: %s", err, input)
	}

	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error indenting JSON: %v\nOriginal JSON: %s", err, input)
	}

	return string(pretty)
}

func main() {
	// Check if we have the right number of arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: parser_tester [code to parse | -f filename] [--method]")
		fmt.Println("\nCurrently supported expressions:")
		fmt.Println("  - Unary message sends: \"self factorial\"")
		fmt.Println("  - Binary message sends: \"1 + 2\"")
		fmt.Println("  - Return statements: \"^self factorial\"")
		fmt.Println("  - Literals: integers, booleans, nil, strings, symbols")
		fmt.Println("\nExamples:")
		fmt.Println("  parser_tester \"1 + 2\"")
		fmt.Println("  parser_tester \"self factorial\"")
		fmt.Println("  parser_tester \"^self factorial\"") 
		fmt.Println("  parser_tester -f mycode.st")
		fmt.Println("  parser_tester \"yourself ^self\" --method")
		fmt.Println("\nNote: The parser is still under development and doesn't yet support:")
		fmt.Println("  - Assignment expressions (x := 5)")
		fmt.Println("  - Block literals [...]")
		fmt.Println("  - Keyword messages with block arguments")
		os.Exit(1)
	}

	// Determine if we're parsing a method or an expression
	methodMode := false
	for _, arg := range os.Args {
		if arg == "--method" {
			methodMode = true
			break
		}
	}

	// Get the code to parse - either from a file or directly from arguments
	var code string
	if os.Args[1] == "-f" {
		if len(os.Args) < 3 {
			fmt.Println("Error: No file specified after -f")
			os.Exit(1)
		}
		fileContent, err := ioutil.ReadFile(os.Args[2])
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		code = string(fileContent)
	} else {
		code = strings.Join(os.Args[1:], " ")
		// Remove --method flag if present
		code = strings.ReplaceAll(code, " --method", "")
	}

	// Parse the code
	node, err := parseCode(code, methodMode)
	if err != nil {
		fmt.Printf("Error parsing code: %v\n", err)
		fmt.Println("\nThe parser is still under development and doesn't yet support:")
		fmt.Println("  - Assignment expressions (x := 5)")
		fmt.Println("  - Block literals [...]")
		fmt.Println("  - Keyword messages with block arguments")
		fmt.Println("\nTry simpler expressions like:")
		fmt.Println("  - \"1 + 2\"")
		fmt.Println("  - \"self factorial\"")
		os.Exit(1)
	}

	// Convert the AST to JSON
	visitor := &JSONVisitor{}
	jsonResult := node.Accept(visitor).(string)

	// Print the pretty JSON
	fmt.Println(prettyJSON(jsonResult))
}