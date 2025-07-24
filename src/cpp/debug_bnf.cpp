#include "simple_parser.h"
#include <iostream>

using namespace smalltalk;

void testExpression(const std::string& expr) {
    std::cout << "\n=== Testing: '" << expr << "' ===" << std::endl;
    try {
        SimpleParser parser(expr);
        auto method = parser.parseMethod();
        std::cout << "✓ SUCCESS: Parsed successfully!" << std::endl;
    } catch (const std::exception& e) {
        std::cout << "✗ FAILED: " << e.what() << std::endl;
    }
}

int main() {
    std::cout << "Testing Smalltalk BNF compliance..." << std::endl;
    
    // Test basic cases first
    testExpression("42");           // primary -> integer
    testExpression("Array");        // primary -> identifier  
    testExpression("Array new");    // unaryExpression -> primary unarySelector
    testExpression("Array new: 3"); // keywordExpression -> binaryExpression keyword binaryExpression
    
    return 0;
}