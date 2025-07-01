#include "simple_parser.h"
#include "simple_compiler.h"
#include "simple_vm.h"
#include "tagged_value.h"

#include <iostream>
#include <string>
#include <cassert>

using namespace smalltalk;

void testComparison(const std::string& expression, bool expectedResult) {
    try {
        // Parse, compile, and execute the expression
        SimpleParser parser(expression);
        auto methodAST = parser.parseMethod();
        
        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);
        
        SimpleVM vm;
        TaggedValue result = vm.execute(*compiledMethod);
        
        // Check that result is a boolean
        if (!result.isBoolean()) {
            std::cerr << "FAIL: " << expression << " -> Expected boolean, got " << result << std::endl;
            assert(false);
        }
        
        // Check the actual boolean value
        bool actualResult = result.asBoolean();
        if (actualResult != expectedResult) {
            std::cerr << "FAIL: " << expression << " -> Expected " << (expectedResult ? "true" : "false") 
                      << ", got " << (actualResult ? "true" : "false") << std::endl;
            assert(false);
        }
        
        std::cout << "PASS: " << expression << " -> " << (actualResult ? "true" : "false") << std::endl;
        
    } catch (const std::exception& e) {
        std::cerr << "FAIL: " << expression << " -> Exception: " << e.what() << std::endl;
        assert(false);
    }
}

void runComparisonTests() {
    std::cout << "Running comparison operation tests..." << std::endl;
    
    // Less than tests
    testComparison("3 < 5", true);
    testComparison("5 < 3", false);
    testComparison("4 < 4", false);
    
    // Greater than tests
    testComparison("7 > 2", true);
    testComparison("2 > 7", false);
    testComparison("5 > 5", false);
    
    // Equal tests
    testComparison("3 = 3", true);
    testComparison("3 = 4", false);
    testComparison("0 = 0", true);
    testComparison("42 = 42", true);
    
    // Not equal tests
    testComparison("4 ~= 5", true);
    testComparison("3 ~= 3", false);
    testComparison("0 ~= 1", true);
    
    // Less than or equal tests
    testComparison("4 <= 4", true);
    testComparison("3 <= 5", true);
    testComparison("6 <= 4", false);
    
    // Greater than or equal tests
    testComparison("5 >= 3", true);
    testComparison("4 >= 4", true);
    testComparison("2 >= 5", false);
    
    // Complex expressions with comparisons
    testComparison("(3 + 2) < (4 * 2)", true);  // 5 < 8
    testComparison("(10 - 3) > (2 * 3)", true); // 7 > 6
    testComparison("(6 / 2) = (1 + 2)", true);  // 3 = 3
    
    std::cout << "All comparison tests passed!" << std::endl;
}

int main() {
    runComparisonTests();
    return 0;
}