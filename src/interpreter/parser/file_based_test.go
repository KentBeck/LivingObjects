package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/pile"
	"smalltalklsp/interpreter/vm"
)

// parseSmalltalkExpression parses a Smalltalk expression or method and returns the AST as JSON
func parseSmalltalkExpression(expression string, isMethod bool) (string, error) {
	// Create a class for context
	objectClass := pile.NewClass("Object", nil)
	objectClass.ClassField = objectClass // Set class's class to itself for proper checks
	classObj := pile.ClassToObject(objectClass)

	// Create a VM
	vmInstance := vm.NewVM()

	// Create a parser with the VM
	p := NewParser(expression, classObj, vmInstance)

	// Parse the expression
	var node ast.Node
	var err error
	if isMethod {
		node, err = p.Parse()
	} else {
		node, err = p.ParseExpression()
	}
	if err != nil {
		return "", err
	}

	// Convert the node to JSON using our visitor
	visitor := &jsonVisitor{}
	jsonResult := visitor.visitNode(node)

	// Return the JSON string
	return jsonResult, nil
}

// jsonVisitor is a simplified visitor that converts AST nodes to JSON
type jsonVisitor struct{}

func (v *jsonVisitor) visitNode(node ast.Node) string {
	switch n := node.(type) {
	case *ast.MethodNode:
		return v.visitMethodNode(n)
	case *ast.ReturnNode:
		return v.visitReturnNode(n)
	case *ast.SelfNode:
		return v.visitSelfNode(n)
	case *ast.LiteralNode:
		return v.visitLiteralNode(n)
	case *ast.VariableNode:
		return v.visitVariableNode(n)
	case *ast.AssignmentNode:
		return v.visitAssignmentNode(n)
	case *ast.MessageSendNode:
		return v.visitMessageSendNode(n)
	case *ast.BlockNode:
		return v.visitBlockNode(n)
	default:
		return fmt.Sprintf(`{"type": "Unknown", "value": "%T"}`, n)
	}
}

func (v *jsonVisitor) visitMethodNode(node *ast.MethodNode) string {
	// Convert body to JSON
	bodyJSON := "null"
	if node.Body != nil {
		bodyJSON = v.visitNode(node.Body)
	}

	// Convert parameters to JSON array
	paramsJSON := fmt.Sprintf(`["%s"]`, strings.Join(node.Parameters, `", "`))
	if len(node.Parameters) == 0 {
		paramsJSON = "[]"
	}

	// Convert temporaries to JSON array
	tempsJSON := fmt.Sprintf(`["%s"]`, strings.Join(node.Temporaries, `", "`))
	if len(node.Temporaries) == 0 {
		tempsJSON = "[]"
	}

	return fmt.Sprintf(`{"type":"MethodNode","selector":"%s","parameters":%s,"temporaries":%s,"body":%s}`,
		node.Selector, paramsJSON, tempsJSON, bodyJSON)
}

func (v *jsonVisitor) visitReturnNode(node *ast.ReturnNode) string {
	// Convert expression to JSON
	exprJSON := "null"
	if node.Expression != nil {
		exprJSON = v.visitNode(node.Expression)
	}

	return fmt.Sprintf(`{"type":"ReturnNode","expression":%s}`, exprJSON)
}

func (v *jsonVisitor) visitSelfNode(node *ast.SelfNode) string {
	return `{"type":"SelfNode"}`
}

func (v *jsonVisitor) visitLiteralNode(node *ast.LiteralNode) string {
	literalJSON := "null"
	if node.Value != nil {
		// Try to convert the literal value based on its type
		if pile.IsIntegerImmediate(node.Value) {
			literalJSON = fmt.Sprintf(`{"type":"Integer","value":%d}`,
				pile.GetIntegerImmediate(node.Value))
		} else if pile.IsTrueImmediate(node.Value) {
			literalJSON = `{"type":"Boolean","value":true}`
		} else if pile.IsFalseImmediate(node.Value) {
			literalJSON = `{"type":"Boolean","value":false}`
		} else if pile.IsNilImmediate(node.Value) {
			literalJSON = `{"type":"Nil"}`
		} else if pile.IsFloatImmediate(node.Value) {
			literalJSON = fmt.Sprintf(`{"type":"Float","value":%f}`,
				pile.GetFloatImmediate(node.Value))
		} else if node.Value.Type() == pile.OBJ_STRING {
			str := pile.ObjectToString(node.Value)
			literalJSON = fmt.Sprintf(`{"type":"String","value":"%s"}`, escapeString(str.GetValue()))
		} else if node.Value.Type() == pile.OBJ_SYMBOL {
			sym := pile.ObjectToSymbol(node.Value)
			literalJSON = fmt.Sprintf(`{"type":"Symbol","value":"%s"}`, escapeString(sym.GetValue()))
		} else if node.Value.Type() == pile.OBJ_ARRAY {
			array := pile.ObjectToArray(node.Value)
			elements := make([]string, array.Size())
			for i := 0; i < array.Size(); i++ {
				elem := array.At(i)
				if pile.IsIntegerImmediate(elem) {
					elements[i] = fmt.Sprintf(`{"type":"Integer","value":%d}`,
						pile.GetIntegerImmediate(elem))
				} else if pile.IsTrueImmediate(elem) {
					elements[i] = `{"type":"Boolean","value":true}`
				} else if pile.IsFalseImmediate(elem) {
					elements[i] = `{"type":"Boolean","value":false}`
				} else if pile.IsNilImmediate(elem) {
					elements[i] = `{"type":"Nil"}`
				} else if elem.Type() == pile.OBJ_STRING {
					str := pile.ObjectToString(elem)
					escapedStr := escapeString(str.GetValue())
					elements[i] = fmt.Sprintf(`{"type":"String","value":"%s"}`, escapedStr)
				} else {
					elements[i] = "null"
				}
			}
			literalJSON = fmt.Sprintf(`{"type":"Array","elements":[%s]}`, strings.Join(elements, ","))
		} else {
			// For unknown types, just use a generic description
			literalJSON = fmt.Sprintf(`{"type":"Object","objectType":%d}`, node.Value.Type())
		}
	}

	return fmt.Sprintf(`{"type":"LiteralNode","value":%s}`, literalJSON)
}

func (v *jsonVisitor) visitVariableNode(node *ast.VariableNode) string {
	return fmt.Sprintf(`{"type":"VariableNode","name":"%s"}`, node.Name)
}

func (v *jsonVisitor) visitAssignmentNode(node *ast.AssignmentNode) string {
	// Convert expression to JSON
	exprJSON := "null"
	if node.Expression != nil {
		exprJSON = v.visitNode(node.Expression)
	}

	return fmt.Sprintf(`{"type":"AssignmentNode","variable":"%s","expression":%s}`,
		node.Variable, exprJSON)
}

func (v *jsonVisitor) visitMessageSendNode(node *ast.MessageSendNode) string {
	// Convert receiver to JSON
	receiverJSON := "null"
	if node.Receiver != nil {
		receiverJSON = v.visitNode(node.Receiver)
	}

	// Convert arguments to JSON array
	argJSONs := make([]string, 0, len(node.Arguments))
	for _, arg := range node.Arguments {
		if arg != nil {
			argJSONs = append(argJSONs, v.visitNode(arg))
		} else {
			argJSONs = append(argJSONs, "null")
		}
	}

	argsJSON := "[]"
	if len(argJSONs) > 0 {
		argsJSON = fmt.Sprintf(`[%s]`, strings.Join(argJSONs, ","))
	}

	return fmt.Sprintf(`{"type":"MessageSendNode","receiver":%s,"selector":"%s","arguments":%s}`,
		receiverJSON, node.Selector, argsJSON)
}

func (v *jsonVisitor) visitBlockNode(node *ast.BlockNode) string {
	// Convert body to JSON
	bodyJSON := "null"
	if node.Body != nil {
		bodyJSON = v.visitNode(node.Body)
	}

	// Convert parameters to JSON array
	paramsJSON := fmt.Sprintf(`["%s"]`, strings.Join(node.Parameters, `", "`))
	if len(node.Parameters) == 0 {
		paramsJSON = "[]"
	}

	// Convert temporaries to JSON array
	tempsJSON := fmt.Sprintf(`["%s"]`, strings.Join(node.Temporaries, `", "`))
	if len(node.Temporaries) == 0 {
		tempsJSON = "[]"
	}

	return fmt.Sprintf(`{"type":"BlockNode","parameters":%s,"temporaries":%s,"body":%s}`,
		paramsJSON, tempsJSON, bodyJSON)
}

// escapeString escapes special characters in a string for JSON
func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// areJSONEqual compares two JSON strings by parsing them into Go objects
func areJSONEqual(json1, json2 string) (bool, error) {
	var obj1, obj2 interface{}

	// Try to unmarshal the first JSON string
	if err := json.Unmarshal([]byte(json1), &obj1); err != nil {
		return false, fmt.Errorf("error unmarshalling first JSON: %v", err)
	}

	// Try to unmarshal the second JSON string
	if err := json.Unmarshal([]byte(json2), &obj2); err != nil {
		return false, fmt.Errorf("error unmarshalling second JSON: %v", err)
	}

	// Compare the two objects
	equal := jsonEqual(obj1, obj2)
	return equal, nil
}

// jsonEqual recursively compares two JSON objects
func jsonEqual(obj1, obj2 interface{}) bool {
	switch v1 := obj1.(type) {
	case map[string]interface{}:
		v2, ok := obj2.(map[string]interface{})
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		for key, value1 := range v1 {
			value2, ok := v2[key]
			if !ok {
				return false
			}
			if !jsonEqual(value1, value2) {
				return false
			}
		}
		return true
	case []interface{}:
		v2, ok := obj2.([]interface{})
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		for i, value1 := range v1 {
			if !jsonEqual(value1, v2[i]) {
				return false
			}
		}
		return true
	case string:
		v2, ok := obj2.(string)
		return ok && v1 == v2
	case float64:
		v2, ok := obj2.(float64)
		return ok && v1 == v2
	case bool:
		v2, ok := obj2.(bool)
		return ok && v1 == v2
	case nil:
		return obj2 == nil
	default:
		return false
	}
}

// runFileBasedTests runs tests from a specified test file
func runFileBasedTests(t *testing.T, testFilePath string) {
	// Open the test file
	file, err := os.Open(testFilePath)
	if err != nil {
		t.Fatalf("Error opening test file: %v", err)
	}
	defer file.Close()

	// Parse the test file
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Split the line by !
		parts := strings.Split(line, "!")
		if len(parts) != 4 {
			t.Errorf("Line %d: Invalid format, expected 4 parts separated by !, got %d", lineNumber, len(parts))
			continue
		}
		
		testName := parts[0]
		expression := parts[1]
		testType := parts[2] // "expression" or "method"
		expectedJSON := parts[3]
		
		// Run the test
		t.Run(testName, func(t *testing.T) {
			// Parse the expression
			isMethod := testType == "method"
			actualJSON, err := parseSmalltalkExpression(expression, isMethod)
			if err != nil {
				t.Fatalf("Error parsing expression '%s': %v", expression, err)
			}
			
			// Compare the actual JSON with the expected JSON
			equal, err := areJSONEqual(actualJSON, expectedJSON)
			if err != nil {
				t.Errorf("Error comparing JSON: %v", err)
				t.Logf("Expected: %s", expectedJSON)
				t.Logf("Actual: %s", actualJSON)
				return
			}
			
			if !equal {
				t.Errorf("Unexpected JSON result")
				t.Logf("Expression: %s", expression)
				t.Logf("Expected: %s", expectedJSON)
				t.Logf("Actual: %s", actualJSON)
			}
		})
	}
	
	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading test file: %v", err)
	}
}

// TestFileBasedExpressions runs tests from the expression_tests.txt file
func TestFileBasedExpressions(t *testing.T) {
	testFilePath := filepath.Join("testdata", "expression_tests.txt")
	runFileBasedTests(t, testFilePath)
}

// TestFileBasedMethods runs tests from the method_tests.txt file
func TestFileBasedMethods(t *testing.T) {
	testFilePath := filepath.Join("testdata", "method_tests.txt")
	runFileBasedTests(t, testFilePath)
}