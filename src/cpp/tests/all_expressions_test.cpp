#include "simple_parser.h"
#include "simple_compiler.h"
#include "simple_vm.h"
#include "tagged_value.h"
#include "smalltalk_string.h"
#include "smalltalk_class.h"
#include "primitive_methods.h"

#include <iostream>
#include <string>
#include <vector>

using namespace smalltalk;

struct ExpressionTest {
    std::string expression;
    std::string expectedResult;
    bool shouldPass;
    std::string category;
};

void testExpression(const ExpressionTest& test) {
    std::cout << "Testing: " << test.expression << " -> " << test.expectedResult;
    
    try {
        // Parse, compile, and execute the expression
        SimpleParser parser(test.expression);
        auto methodAST = parser.parseMethod();
        
        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);
        
        SimpleVM vm;
        TaggedValue result = vm.execute(*compiledMethod);
        
        // Convert result to string for comparison
        std::string resultStr;
        if (result.isInteger()) {
            resultStr = std::to_string(result.asInteger());
        } else if (result.isBoolean()) {
            resultStr = result.asBoolean() ? "true" : "false";
        } else if (result.isNil()) {
            resultStr = "nil";
        } else if (StringUtils::isString(result)) {
            String* str = StringUtils::asString(result);
            resultStr = str->getContent(); // Get content without quotes for comparison
        } else {
            resultStr = "Object";
        }
        
        if (test.shouldPass && resultStr == test.expectedResult) {
            std::cout << " ✅ PASS" << std::endl;
        } else if (test.shouldPass) {
            std::cout << " ❌ FAIL (got: " << resultStr << ")" << std::endl;
        } else {
            std::cout << " ❌ FAIL (should have failed but got: " << resultStr << ")" << std::endl;
        }
        
    } catch (const std::exception& e) {
        if (test.shouldPass) {
            std::cout << " ❌ FAIL (exception: " << e.what() << ")" << std::endl;
        } else {
            std::cout << " ✅ EXPECTED FAIL (" << e.what() << ")" << std::endl;
        }
    }
}

int main() {
    // Initialize class system and primitives before running tests
    ClassUtils::initializeCoreClasses();
    
    // Initialize primitive registry
    auto& primitiveRegistry = PrimitiveRegistry::getInstance();
    primitiveRegistry.initializeCorePrimitives();
    
    // Add primitive methods to Integer class
    Class* integerClass = ClassUtils::getIntegerClass();
    IntegerClassSetup::addPrimitiveMethods(integerClass);
    
    std::vector<ExpressionTest> tests = {
        // Basic arithmetic - SHOULD PASS
        {"3 + 4", "7", true, "arithmetic"},
        {"5 - 2", "3", true, "arithmetic"},
        {"2 * 3", "6", true, "arithmetic"},
        {"10 / 2", "5", true, "arithmetic"},
        
        // Complex arithmetic - SHOULD PASS
        {"(3 + 2) * 4", "20", true, "arithmetic"},
        {"10 - 2 * 3", "4", true, "arithmetic"},
        {"(10 - 2) / 4", "2", true, "arithmetic"},
        
        // Integer comparisons - SHOULD PASS
        {"3 < 5", "true", true, "comparison"},
        {"7 > 2", "true", true, "comparison"},
        {"3 = 3", "true", true, "comparison"},
        {"4 ~= 5", "true", true, "comparison"},
        {"4 <= 4", "true", true, "comparison"},
        {"5 >= 3", "true", true, "comparison"},
        {"5 < 3", "false", true, "comparison"},
        {"2 > 7", "false", true, "comparison"},
        {"3 = 4", "false", true, "comparison"},
        
        // Complex comparisons - SHOULD PASS
        {"(3 + 2) < (4 * 2)", "true", true, "comparison"},
        {"(10 - 3) > (2 * 3)", "true", true, "comparison"},
        {"(6 / 2) = (1 + 2)", "true", true, "comparison"},
        
        // Basic object creation - SHOULD FAIL (not implemented)
        {"Object new", "<Object>", false, "object_creation"},
        {"Array new: 3", "<Array size: 3>", false, "object_creation"},
        
        // String literals - SHOULD PASS (basic string parsing)
        {"'hello'", "hello", true, "strings"},
        {"'world'", "world", true, "strings"},
        
        // String operations - SHOULD FAIL (not implemented)
        {"'hello' , ' world'", "hello world", false, "string_operations"},
        {"'hello' size", "5", false, "string_operations"},
        
        // Literals - SHOULD PASS (now implemented)
        {"true", "true", true, "literals"},
        {"false", "false", true, "literals"},
        {"nil", "nil", true, "literals"},
        
        // Variable assignment - SHOULD FAIL (not implemented)
        {"| x | x := 42. x", "42", false, "variables"},
        
        // Blocks - SHOULD FAIL (not implemented)
        {"[3 + 4] value", "7", false, "blocks"},
        {"[:x | x + 1] value: 5", "6", false, "blocks"},
        
        // Conditionals - SHOULD FAIL (not implemented)
        {"(3 < 4) ifTrue: [10] ifFalse: [20]", "10", false, "conditionals"},
        {"true ifTrue: [42]", "42", false, "conditionals"},
        
        // Collections - SHOULD FAIL (not implemented)
        {"#(1 2 3) at: 2", "2", false, "collections"},
        {"#(1 2 3) size", "3", false, "collections"},
        
        // Dictionary operations - SHOULD FAIL (not implemented)
        {"Dictionary new", "<Dictionary>", false, "dictionaries"},
        
        // Class creation - SHOULD FAIL (not implemented)
        {"Object subclass: #Point", "<Class: Point>", false, "class_creation"}
    };
    
    std::cout << "=== Smalltalk Expression Test Suite ===" << std::endl;
    std::cout << "Testing " << tests.size() << " expressions..." << std::endl << std::endl;
    
    int passCount = 0;
    int totalCount = 0;
    std::string currentCategory = "";
    
    for (const auto& test : tests) {
        if (test.category != currentCategory) {
            currentCategory = test.category;
            std::cout << std::endl << "=== " << currentCategory << " ===" << std::endl;
        }
        
        testExpression(test);
        
        // Count as pass if result matches expectation (either should pass and did, or should fail and did)
        try {
            SimpleParser parser(test.expression);
            auto methodAST = parser.parseMethod();
            SimpleCompiler compiler;
            auto compiledMethod = compiler.compile(*methodAST);
            SimpleVM vm;
            TaggedValue result = vm.execute(*compiledMethod);
                    
            std::string resultStr;
            if (result.isInteger()) {
                resultStr = std::to_string(result.asInteger());
            } else if (result.isBoolean()) {
                resultStr = result.asBoolean() ? "true" : "false";
            } else if (result.isNil()) {
                resultStr = "nil";
            } else if (StringUtils::isString(result)) {
                String* str = StringUtils::asString(result);
                resultStr = str->getContent(); // Get content without quotes for comparison
            } else {
                resultStr = "Object";
            }
            
            if (test.shouldPass && resultStr == test.expectedResult) {
                passCount++;
            } else if (!test.shouldPass) {
                // This should have failed but didn't - that's actually bad
            }
        } catch (const std::exception&) {
            if (!test.shouldPass) {
                passCount++; // Expected to fail and did fail
            }
        }
        totalCount++;
    }
    
    std::cout << std::endl << "=== SUMMARY ===" << std::endl;
    std::cout << "Expressions that work correctly: " << passCount << "/" << totalCount << std::endl;
    
    // Count by category
    std::cout << std::endl << "By category:" << std::endl;
    std::vector<std::string> categories = {"arithmetic", "comparison", "object_creation", "strings", "string_operations", "literals", "variables", "blocks", "conditionals", "collections", "dictionaries", "class_creation"};
    
    for (const auto& category : categories) {
        int categoryPass = 0;
        int categoryTotal = 0;
        
        for (const auto& test : tests) {
            if (test.category == category) {
                categoryTotal++;
                
                bool actuallyPassed = false;
                try {
                    SimpleParser parser(test.expression);
                    auto methodAST = parser.parseMethod();
                    SimpleCompiler compiler;
                    auto compiledMethod = compiler.compile(*methodAST);
                    SimpleVM vm;
                    TaggedValue result = vm.execute(*compiledMethod);
                    
                    std::string resultStr;
                    if (result.isInteger()) {
                        resultStr = std::to_string(result.asInteger());
                    } else if (result.isBoolean()) {
                        resultStr = result.asBoolean() ? "true" : "false";
                    } else if (result.isNil()) {
                        resultStr = "nil";
                    } else {
                        resultStr = "Object";
                    }
                    
                    if (test.shouldPass && resultStr == test.expectedResult) {
                        actuallyPassed = true;
                    }
                } catch (const std::exception&) {
                    if (!test.shouldPass) {
                        actuallyPassed = true;
                    }
                }
                
                if (actuallyPassed) categoryPass++;
            }
        }
        
        if (categoryTotal > 0) {
            std::cout << "  " << category << ": " << categoryPass << "/" << categoryTotal;
            if (categoryPass == categoryTotal) {
                std::cout << " ✅";
            } else {
                std::cout << " ❌";
            }
            std::cout << std::endl;
        }
    }
    
    return 0;
}