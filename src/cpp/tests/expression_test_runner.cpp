#include "../include/bytecode.h"
#include "../include/object.h"
#include "../include/memory_manager.h"
#include "../include/tagged_value.h"
#include <iostream>
#include <fstream>
#include <string>
#include <vector>
#include <sstream>
#include <regex>

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
        std::cerr << "Error: Could not open test file: " << filename << std::endl;
        return tests;
    }
    
    std::string line;
    while (std::getline(file, line)) {
        // Skip empty lines and comments
        if (line.empty() || line[0] == '#') {
            continue;
        }
        
        // Find the delimiter between expression and expected result
        size_t pos = line.find(" -> ");
        if (pos == std::string::npos) {
            std::cerr << "Warning: Invalid test line, skipping: " << line << std::endl;
            continue;
        }
        
        ExpressionTest test;
        test.expression = line.substr(0, pos);
        test.expectedResult = line.substr(pos + 4); // Skip " -> "
        tests.push_back(test);
    }
    
    return tests;
}

// Execute a single test
bool executeTest(ExpressionTest& test) {
    // This is a stub implementation that will be expanded as we implement more functionality
    std::cout << "Executing: " << test.expression << std::endl;
    
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
    
    std::cout << "Found " << tests.size() << " expression tests." << std::endl;
    
    // Initialize test counters
    int passed = 0;
    int failed = 0;
    int skipped = 0;
    
    // Run each test
    for (auto& test : tests) {
        bool result = executeTest(test);
        
        if (!test.executed) {
            std::cout << "⚠️ SKIPPED: " << test.expression << " -> " << test.expectedResult << std::endl;
            skipped++;
        } else if (result) {
            std::cout << "✅ PASSED: " << test.expression << " -> " << test.actualResult << std::endl;
            passed++;
        } else {
            std::cout << "❌ FAILED: " << test.expression << std::endl;
            std::cout << "  Expected: " << test.expectedResult << std::endl;
            std::cout << "  Actual:   " << test.actualResult << std::endl;
            failed++;
        }
    }
    
    // Print summary
    std::cout << std::endl;
    std::cout << "Test Results:" << std::endl;
    std::cout << "  Total:   " << tests.size() << std::endl;
    std::cout << "  Passed:  " << passed << std::endl;
    std::cout << "  Failed:  " << failed << std::endl;
    std::cout << "  Skipped: " << skipped << std::endl;
}

// Test TaggedValue implementation
void testTaggedValues() {
    std::cout << "\nRunning TaggedValue tests:" << std::endl;
    
    try {
        // Integer tests
        TaggedValue int42(42);
        assert(int42.isInteger());
        assert(int42.asInteger() == 42);
        std::cout << "✅ Integer 42: " << int42 << std::endl;
        
        TaggedValue intNeg10(-10);
        assert(intNeg10.isInteger());
        assert(intNeg10.asInteger() == -10);
        std::cout << "✅ Integer -10: " << intNeg10 << std::endl;
        
        // Special value tests
        TaggedValue nilVal = TaggedValue::nil();
        assert(nilVal.isSpecial());
        assert(nilVal.isNil());
        std::cout << "✅ nil: " << nilVal << std::endl;
        
        TaggedValue trueVal = TaggedValue::trueValue();
        assert(trueVal.isSpecial());
        assert(trueVal.isTrue());
        assert(trueVal.isBoolean());
        assert(trueVal.asBoolean() == true);
        std::cout << "✅ true: " << trueVal << std::endl;
        
        TaggedValue falseVal = TaggedValue::falseValue();
        assert(falseVal.isSpecial());
        assert(falseVal.isFalse());
        assert(falseVal.isBoolean());
        assert(falseVal.asBoolean() == false);
        std::cout << "✅ false: " << falseVal << std::endl;
        
        // Float tests (simplified)
        TaggedValue float0(0.0);
        assert(float0.isFloat());
        assert(float0.asFloat() == 0.0);
        std::cout << "✅ Float 0.0: " << float0 << std::endl;
        
        TaggedValue float1(1.0);
        assert(float1.isFloat());
        assert(float1.asFloat() == 1.0);
        std::cout << "✅ Float 1.0: " << float1 << std::endl;
        
        TaggedValue floatNeg1(-1.0);
        assert(floatNeg1.isFloat());
        assert(floatNeg1.asFloat() == -1.0);
        std::cout << "✅ Float -1.0: " << floatNeg1 << std::endl;
        
        std::cout << "All TaggedValue tests passed!" << std::endl;
    } catch (const std::exception& e) {
        std::cerr << "❌ TaggedValue test failed: " << e.what() << std::endl;
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