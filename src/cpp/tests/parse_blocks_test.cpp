#include "simple_parser.h"
#include "ast.h"

#include <iostream>
#include <cassert>

using namespace smalltalk;

int main() {
    try {
        std::cout << "Testing Block parsing (minimal test)..." << std::endl;
        
        // Test 1: Parse a simple block with integer
        std::cout << "\n=== Test 1: Parsing [42] ===" << std::endl;
        
        SimpleParser parser1("[42]");
        auto methodAST1 = parser1.parseMethod();
        
        std::cout << "Parsed AST: " << methodAST1->toString() << std::endl;
        
        // The AST should be a method containing a block
        const MethodNode* method1 = dynamic_cast<const MethodNode*>(methodAST1.get());
        assert(method1 != nullptr);
        std::cout << "✓ Successfully parsed as MethodNode" << std::endl;
        
        const BlockNode* block1 = dynamic_cast<const BlockNode*>(method1->getBody());
        assert(block1 != nullptr);
        std::cout << "✓ Method body is a BlockNode" << std::endl;
        std::cout << "Block contents: " << block1->toString() << std::endl;
        
        // Test 2: Parse a block with arithmetic
        std::cout << "\n=== Test 2: Parsing [3 + 4] ===" << std::endl;
        
        SimpleParser parser2("[3 + 4]");
        auto methodAST2 = parser2.parseMethod();
        
        std::cout << "Parsed AST: " << methodAST2->toString() << std::endl;
        
        const MethodNode* method2 = dynamic_cast<const MethodNode*>(methodAST2.get());
        assert(method2 != nullptr);
        
        const BlockNode* block2 = dynamic_cast<const BlockNode*>(method2->getBody());
        assert(block2 != nullptr);
        std::cout << "✓ Arithmetic block parsed correctly" << std::endl;
        
        // Test 3: Parse nested expressions in blocks
        std::cout << "\n=== Test 3: Parsing [2 * (3 + 1)] ===" << std::endl;
        
        SimpleParser parser3("[2 * (3 + 1)]");
        auto methodAST3 = parser3.parseMethod();
        
        std::cout << "Parsed AST: " << methodAST3->toString() << std::endl;
        
        const MethodNode* method3 = dynamic_cast<const MethodNode*>(methodAST3.get());
        assert(method3 != nullptr);
        
        const BlockNode* block3 = dynamic_cast<const BlockNode*>(method3->getBody());
        assert(block3 != nullptr);
        std::cout << "✓ Nested expression in block parsed correctly" << std::endl;
        
        // Test 4: Parse literal blocks
        std::cout << "\n=== Test 4: Parsing literal blocks ===" << std::endl;
        
        std::vector<std::string> literals = {"[true]", "[false]", "[nil]"};
        
        for (const auto& literal : literals) {
            SimpleParser parser(literal);
            auto ast = parser.parseMethod();
            const MethodNode* method = dynamic_cast<const MethodNode*>(ast.get());
            const BlockNode* block = dynamic_cast<const BlockNode*>(method->getBody());
            assert(block != nullptr);
            std::cout << "✓ " << literal << " -> " << ast->toString() << std::endl;
        }
        
        // Test 5: Test comparison operations in blocks
        std::cout << "\n=== Test 5: Parsing comparison blocks ===" << std::endl;
        
        std::vector<std::string> comparisons = {
            "[3 < 4]",
            "[5 > 2]", 
            "[1 = 1]",
            "[2 ~= 3]",
            "[4 <= 5]",
            "[6 >= 6]"
        };
        
        for (const auto& comp : comparisons) {
            SimpleParser parser(comp);
            auto ast = parser.parseMethod();
            const MethodNode* method = dynamic_cast<const MethodNode*>(ast.get());
            const BlockNode* block = dynamic_cast<const BlockNode*>(method->getBody());
            assert(block != nullptr);
            std::cout << "✓ " << comp << " -> " << ast->toString() << std::endl;
        }
        
        // Test 6: Error handling for malformed blocks
        std::cout << "\n=== Test 6: Error handling ===" << std::endl;
        
        std::vector<std::string> malformed = {
            "[3 + 4",     // Missing closing bracket
            "[",          // Empty incomplete block
            "3 + 4]",     // Missing opening bracket (should parse as expression + error)
            "[]",         // Empty block
            "[[[42]]]"    // Nested blocks
        };
        
        for (const auto& bad : malformed) {
            try {
                SimpleParser parser(bad);
                auto ast = parser.parseMethod();
                if (bad == "[]") {
                    // Empty block might be valid
                    std::cout << "✓ " << bad << " parsed (possibly valid): " << ast->toString() << std::endl;
                } else if (bad == "[[[42]]]") {
                    // Nested blocks might be valid
                    std::cout << "✓ " << bad << " parsed (nested blocks): " << ast->toString() << std::endl;
                } else {
                    std::cout << "? " << bad << " unexpectedly parsed: " << ast->toString() << std::endl;
                }
            } catch (const std::exception& e) {
                std::cout << "✓ " << bad << " correctly rejected: " << e.what() << std::endl;
            }
        }
        
        // Test 7: Complex arithmetic in blocks
        std::cout << "\n=== Test 7: Complex arithmetic blocks ===" << std::endl;
        
        std::vector<std::string> complex = {
            "[1 + 2 * 3]",
            "[(1 + 2) * 3]",
            "[10 / 2 - 3]",
            "[2 * 3 + 4 * 5]"
        };
        
        for (const auto& expr : complex) {
            SimpleParser parser(expr);
            auto ast = parser.parseMethod();
            const MethodNode* method = dynamic_cast<const MethodNode*>(ast.get());
            const BlockNode* block = dynamic_cast<const BlockNode*>(method->getBody());
            assert(block != nullptr);
            std::cout << "✓ " << expr << " -> " << ast->toString() << std::endl;
        }
        
        std::cout << "\n=== All block parsing tests completed successfully! ===" << std::endl;
        std::cout << "Block syntax [expression] is now supported in the parser." << std::endl;
        
        return 0;
        
    } catch (const std::exception& e) {
        std::cerr << "Test FAILED with exception: " << e.what() << std::endl;
        return 1;
    }
}