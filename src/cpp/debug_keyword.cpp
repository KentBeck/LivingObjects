#include "simple_parser.h"
#include <iostream>

using namespace smalltalk;

int main() {
    try {
        std::cout << "Testing simple keyword message parsing..." << std::endl;
        
        // Test 1: Simple case
        SimpleParser parser1("Object new");
        auto method1 = parser1.parseMethod();
        std::cout << "✓ 'Object new' parsed successfully" << std::endl;
        
        // Test 2: Keyword with argument
        SimpleParser parser2("Array new: 3");
        auto method2 = parser2.parseMethod();
        std::cout << "✓ 'Array new: 3' parsed successfully" << std::endl;
        
    } catch (const std::exception& e) {
        std::cout << "✗ Parse error: " << e.what() << std::endl;
    }
    return 0;
}