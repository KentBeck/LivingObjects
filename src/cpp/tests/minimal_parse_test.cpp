#include "simple_parser.h"
#include "ast.h"

#include <iostream>
#include <memory>

using namespace smalltalk;

// Minimal test to avoid dependencies
int main() {
    std::cout << "Minimal block parsing test..." << std::endl;
    
    try {
        // Test 1: Simple integer block
        std::cout << "Testing [42]..." << std::endl;
        SimpleParser parser1("[42]");
        auto ast1 = parser1.parseMethod();
        std::cout << "Result: " << ast1->toString() << std::endl;
        
        // Test 2: Arithmetic block  
        std::cout << "Testing [3 + 4]..." << std::endl;
        SimpleParser parser2("[3 + 4]");
        auto ast2 = parser2.parseMethod();
        std::cout << "Result: " << ast2->toString() << std::endl;
        
        // Test 3: Check structure
        const MethodNode* method = dynamic_cast<const MethodNode*>(ast1.get());
        if (method) {
            std::cout << "✓ Parsed as MethodNode" << std::endl;
            const BlockNode* block = dynamic_cast<const BlockNode*>(method->getBody());
            if (block) {
                std::cout << "✓ Contains BlockNode" << std::endl;
            } else {
                std::cout << "✗ No BlockNode found" << std::endl;
            }
        } else {
            std::cout << "✗ Not a MethodNode" << std::endl;
        }
        
        std::cout << "Test completed successfully!" << std::endl;
        return 0;
        
    } catch (const std::exception& e) {
        std::cout << "Error: " << e.what() << std::endl;
        return 1;
    }
}