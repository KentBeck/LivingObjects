#include "bytecode.h"
#include "memory_manager.h"
#include "object.h"
#include "tagged_value.h"

#include <fstream>
#include <iostream>
#include <regex>
#include <sstream>
#include <string>
#include <vector>
#include <cassert>

using namespace smalltalk;

// A simple representation of a Smalltalk expression test
struct ExpressionTest {
    std::string expression;
    std::string expectedResult;
    bool executed = false;
    bool passed = false;
    std::string actualResult;
};

// Parse the expression test file
std::vector<ExpressionTest> parseExpressionTests(const std::string& filename) {
    std::vector<ExpressionTest> tests;
    std::ifstream file(filename);
    
    if (!file.is_open()) {
        std::cerr << "Error: Could not open test file: " << filename << '\n';
        return tests;
    }
    
    std::string line;
    while (std::getline(file, line)) {
        // Skip empty lines and comments
        if (line.empty() || line[0] == '#') {
            continue;
        }
        
        // Find the delimiter between expression and expected result
        size_t pos = line.find(" -> ");        if (pos == std::string::npos) {            std::cerr << "Warning: Invalid test line, skipping: " << line << '\n';            continue;        }                const size_t DELIMITER_LENGTH = 4;        ExpressionTest test;        test.expression = line.substr(0, pos);        test.expectedResult = line.substr(pos + DELIMITER_LENGTH); // Skip " -> "
        tests.push_back(test);
    }
    
    return tests;
}

// Execute a single test
bool executeTest(ExpressionTest& test) {
    // This is a stub implementation that will be expanded as we implement more functionality
    std::cout << "Executing: " << test.expression << '\n';
    
    // Simple pattern matching for test expressions
    // This is just a placeholder - in reality we would parse and execute the expression
    
    if (std::regex_match(test.expression, std::regex("\\d+ \\+ \\d+"))) {
        std::regex pattern("(\\d+) \\+ (\\d+)");
        std::smatch matches;
        std::regex_search(test.expression, matches, pattern);
        
        int a = std::stoi(matches[1].str());
        int b = std::stoi(matches[2].str());
        int result = a + b;
        
        test.actualResult = std::to_string(result);
        test.executed = true;
        test.passed = (test.actualResult == test.expectedResult);
        
        return test.passed;
    }
    
    if (std::regex_match(test.expression, std::regex("\\d+ - \\d+"))) {
        std::regex pattern("(\\d+) - (\\d+)");
        std::smatch matches;
        std::regex_search(test.expression, matches, pattern);
        
        int a = std::stoi(matches[1].str());
        int b = std::stoi(matches[2].str());
        int result = a - b;
        
        test.actualResult = std::to_string(result);
        test.executed = true;
        test.passed = (test.actualResult == test.expectedResult);
        
        return test.passed;
    }
    
    if (std::regex_match(test.expression, std::regex("\\d+ \\* \\d+"))) {
        std::regex pattern("(\\d+) \\* (\\d+)");
        std::smatch matches;
        std::regex_search(test.expression, matches, pattern);
        
        int a = std::stoi(matches[1].str());
        int b = std::stoi(matches[2].str());
        int result = a * b;
        
        test.actualResult = std::to_string(result);
        test.executed = true;
        test.passed = (test.actualResult == test.expectedResult);
        
        return test.passed;
    }
    
    if (std::regex_match(test.expression, std::regex("\\d+ / \\d+"))) {
        std::regex pattern("(\\d+) / (\\d+)");
        std::smatch matches;
        std::regex_search(test.expression, matches, pattern);
        
        int a = std::stoi(matches[1].str());
        int b = std::stoi(matches[2].str());
        int result = a / b;
        
        test.actualResult = std::to_string(result);
        test.executed = true;
        test.passed = (test.actualResult == test.expectedResult);
        
        return test.passed;
    }
    
    // For tests we're not ready to implement yet
    test.executed = false;
    test.actualResult = "NOT IMPLEMENTED";
    return false;
}

// Run all tests
void runExpressionTests(const std::string& filename) {
    std::vector<ExpressionTest> tests = parseExpressionTests(filename);
    
    std::cout << "Found " << tests.size() << " expression tests." << '\n';
    
    // Initialize test counters
    int passed = 0;
    int failed = 0;
    int skipped = 0;
    
    // Run each test
    for (auto& test : tests) {
        bool result = executeTest(test);
        
        if (!test.executed) {
            std::cout << "⚠️ SKIPPED: " << test.expression << " -> " << test.expectedResult << '\n';
            skipped++;
        } else if (result) {
            std::cout << "✅ PASSED: " << test.expression << " -> " << test.actualResult << '\n';
            passed++;
        } else {
            std::cout << "❌ FAILED: " << test.expression << '\n';
            std::cout << "  Expected: " << test.expectedResult << '\n';
            std::cout << "  Actual:   " << test.actualResult << '\n';
            failed++;
        }
    }
    
    // Print summary
    std::cout << '\n';
    std::cout << "Test Results:" << '\n';
    std::cout << "  Total:   " << tests.size() << '\n';
    std::cout << "  Passed:  " << passed << '\n';
    std::cout << "  Failed:  " << failed << '\n';
    std::cout << "  Skipped: " << skipped << '\n';
}

// Test TaggedValue implementation
void testTaggedValues() {
    std::cout << "\nRunning TaggedValue tests:" << '\n';
    
    try {
        // Integer tests
        const int TEST_INT_42 = 42;
        TaggedValue int42(TEST_INT_42);
        assert(int42.isInteger());
        assert(int42.asInteger() == TEST_INT_42);
        std::cout << "✅ Integer 42: " << int42 << '\n';
        
        const int TEST_INT_NEG10 = -10;
        TaggedValue intNeg10(TEST_INT_NEG10);
        assert(intNeg10.isInteger());
        assert(intNeg10.asInteger() == TEST_INT_NEG10);
        std::cout << "✅ Integer -10: " << intNeg10 << '\n';
        
        // Special value tests
        TaggedValue nilVal = TaggedValue::nil();
        assert(nilVal.isSpecial());
        assert(nilVal.isNil());
        std::cout << "✅ nil: " << nilVal << '\n';
        
        TaggedValue trueVal = TaggedValue::trueValue();
        assert(trueVal.isSpecial());
        assert(trueVal.isTrue());
        assert(trueVal.isBoolean());
        assert(trueVal.asBoolean() == true);
        std::cout << "✅ true: " << trueVal << '\n';
        
        TaggedValue falseVal = TaggedValue::falseValue();
        assert(falseVal.isSpecial());
        assert(falseVal.isFalse());
        assert(falseVal.isBoolean());
        assert(falseVal.asBoolean() == false);
        std::cout << "✅ false: " << falseVal << '\n';
        
        // Float tests (simplified)
        TaggedValue float0(0.0);
        assert(float0.isFloat());
        assert(float0.asFloat() == 0.0);
        std::cout << "✅ Float 0.0: " << float0 << '\n';
        
        TaggedValue float1(1.0);
        assert(float1.isFloat());
        assert(float1.asFloat() == 1.0);
        std::cout << "✅ Float 1.0: " << float1 << '\n';
        
        TaggedValue floatNeg1(-1.0);
        assert(floatNeg1.isFloat());
        assert(floatNeg1.asFloat() == -1.0);
        std::cout << "✅ Float -1.0: " << floatNeg1 << '\n';
        
        std::cout << "All TaggedValue tests passed!" << '\n';
    } catch (const std::exception& e) {
        std::cerr << "❌ TaggedValue test failed: " << e.what() << '\n';
        exit(1);
    }
}

int main(int argc, char** argv) {
    // Test the TaggedValue implementation
    testTaggedValues();
    
    // Default test file
    std::string testFile = "tests/expression_tests.txt";
    
    // Allow specifying a test file on the command line
    if (argc > 1) {
        testFile = argv[1];
    }
    
    // Run the expression tests
    runExpressionTests(testFile);
    
    return 0;
}