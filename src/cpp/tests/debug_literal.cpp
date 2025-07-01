#include "simple_parser.h"
#include "simple_compiler.h"
#include "tagged_value.h"
#include <iostream>

using namespace smalltalk;

int main() {
    std::cout << "=== Boolean Literal Compilation Debug ===" << std::endl;
    
    try {
        // Test parsing and compiling "true"
        std::cout << "Testing 'true':" << std::endl;
        SimpleParser parser("true");
        auto methodAST = parser.parseMethod();
        std::cout << "  Parsed successfully" << std::endl;
        
        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);
        std::cout << "  Compiled successfully" << std::endl;
        
        // Check the literals
        auto literals = compiledMethod->getLiterals();
        std::cout << "  Number of literals: " << literals.size() << std::endl;
        
        if (literals.size() > 0) {
            TaggedValue literal = literals[0];
            std::cout << "  First literal properties:" << std::endl;
            std::cout << "    isBoolean(): " << literal.isBoolean() << std::endl;
            std::cout << "    isTrue(): " << literal.isTrue() << std::endl;
            std::cout << "    isFalse(): " << literal.isFalse() << std::endl;
            std::cout << "    isNil(): " << literal.isNil() << std::endl;
            std::cout << "    isInteger(): " << literal.isInteger() << std::endl;
            std::cout << "    isSpecial(): " << literal.isSpecial() << std::endl;
            std::cout << "    rawValue(): " << std::hex << literal.rawValue() << std::dec << std::endl;
            
            // Test the toString logic directly
            if (literal.isNil()) {
                std::cout << "  ToString would be: nil" << std::endl;
            } else if (literal.isBoolean()) {
                std::cout << "  ToString would be: " << (literal.asBoolean() ? "true" : "false") << std::endl;
            } else if (literal.isInteger()) {
                std::cout << "  ToString would be: " << literal.asInteger() << std::endl;
            } else {
                std::cout << "  ToString would be: ?" << std::endl;
            }
        }
        
        std::cout << "  Full method toString: " << compiledMethod->toString() << std::endl;
        
    } catch (const std::exception& e) {
        std::cout << "Error: " << e.what() << std::endl;
    }
    
    return 0;
}