#include "simple_parser.h"
#include <iostream>

using namespace smalltalk;

void testParsing(const std::string& input) {
    try {
        std::cout << "Testing: '" << input << "'" << std::endl;
        SimpleParser parser(input);
        auto method = parser.parseMethod();
        std::cout << "✓ Parse successful!" << std::endl;
    } catch (const std::exception& e) {
        std::cout << "✗ Parse error: " << e.what() << std::endl;
    }
    std::cout << std::endl;
}

int main() {
    std::cout << "Testing keyword message parsing..." << std::endl;
    
    // Test simple cases first
    testParsing("Array");
    testParsing("Object new");
    testParsing("Array new: 3");
    
    return 0;
}