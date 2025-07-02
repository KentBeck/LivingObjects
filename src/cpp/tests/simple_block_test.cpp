#include "simple_parser.h"
#include "simple_compiler.h"
#include "ast.h"
#include "bytecode.h"

#include <iostream>
#include <cassert>

using namespace smalltalk;

int main() {
    try {
        std::cout << "Testing Block parsing and compilation (simple test)..." << std::endl;
        
        // Test 1: Parse a simple block expression
        std::cout << "\n=== Test 1: Parsing [3 + 4] ===" << std::endl;
        
        SimpleParser parser("[3 + 4]");
        auto methodAST = parser.parseMethod();
        
        std::cout << "Parsed AST: " << methodAST->toString() << std::endl;
        
        // The AST should be a method containing a block
        const MethodNode* method = dynamic_cast<const MethodNode*>(methodAST.get());
        assert(method != nullptr);
        std::cout << "✓ Successfully parsed as MethodNode" << std::endl;
        
        const BlockNode* block = dynamic_cast<const BlockNode*>(method->getBody());
        assert(block != nullptr);
        std::cout << "✓ Method body is a BlockNode" << std::endl;
        
        std::cout << "Block AST: " << block->toString() << std::endl;
        
        // Test 2: Compile the block
        std::cout << "\n=== Test 2: Compiling block ===" << std::endl;
        
        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*method);
        
        std::cout << "Compiled method: " << compiledMethod->toString() << std::endl;
        
        // Check that CREATE_BLOCK bytecode was generated
        const auto& bytecodes = compiledMethod->getBytecodes();
        bool hasCreateBlock = false;
        for (size_t i = 0; i < bytecodes.size(); i++) {
            std::cout << "Bytecode[" << i << "]: " << static_cast<int>(bytecodes[i]) 
                     << " (" << getBytecodeString(static_cast<Bytecode>(bytecodes[i])) << ")" << std::endl;
            if (bytecodes[i] == static_cast<uint8_t>(Bytecode::CREATE_BLOCK)) {
                hasCreateBlock = true;
                std::cout << "✓ Found CREATE_BLOCK bytecode at position " << i << std::endl;
            }
        }
        assert(hasCreateBlock);
        
        // Test 3: Test nested expressions in blocks
        std::cout << "\n=== Test 3: Parsing nested block [2 * (3 + 1)] ===" << std::endl;
        
        SimpleParser parser2("[2 * (3 + 1)]");
        auto methodAST2 = parser2.parseMethod();
        
        std::cout << "Parsed AST: " << methodAST2->toString() << std::endl;
        
        const MethodNode* method2 = dynamic_cast<const MethodNode*>(methodAST2.get());
        assert(method2 != nullptr);
        
        const BlockNode* block2 = dynamic_cast<const BlockNode*>(method2->getBody());
        assert(block2 != nullptr);
        std::cout << "✓ Nested expression in block parsed correctly" << std::endl;
        
        // Test 4: Compile nested block
        auto compiledMethod2 = compiler.compile(*method2);
        std::cout << "Compiled nested block: " << compiledMethod2->toString() << std::endl;
        
        // Test 5: Test simple literal blocks
        std::cout << "\n=== Test 4: Parsing literal blocks ===" << std::endl;
        
        SimpleParser parser3("[42]");
        auto methodAST3 = parser3.parseMethod();
        std::cout << "Parsed [42]: " << methodAST3->toString() << std::endl;
        
        SimpleParser parser4("[true]");
        auto methodAST4 = parser4.parseMethod();
        std::cout << "Parsed [true]: " << methodAST4->toString() << std::endl;
        
        SimpleParser parser5("[nil]");
        auto methodAST5 = parser5.parseMethod();
        std::cout << "Parsed [nil]: " << methodAST5->toString() << std::endl;
        
        // Test 6: Test error handling for malformed blocks
        std::cout << "\n=== Test 5: Error handling ===" << std::endl;
        
        try {
            SimpleParser parser6("[3 + 4"); // Missing closing bracket
            auto methodAST6 = parser6.parseMethod();
            assert(false && "Should have thrown an error for malformed block");
        } catch (const std::exception& e) {
            std::cout << "✓ Correctly caught error for malformed block: " << e.what() << std::endl;
        }
        
        try {
            SimpleParser parser7("["); // Empty block with missing closing bracket
            auto methodAST7 = parser7.parseMethod();
            assert(false && "Should have thrown an error for incomplete block");
        } catch (const std::exception& e) {
            std::cout << "✓ Correctly caught error for incomplete block: " << e.what() << std::endl;
        }
        
        // Test 7: Test multiple blocks (although not in the same expression for now)
        std::cout << "\n=== Test 6: Multiple separate block expressions ===" << std::endl;
        
        std::vector<std::string> blockExpressions = {
            "[1]",
            "[1 + 2]", 
            "[3 * 4]",
            "[true]",
            "[false]",
            "[nil]",
            "[(1 + 2) * 3]"
        };
        
        for (const auto& expr : blockExpressions) {
            try {
                SimpleParser parser(expr);
                auto ast = parser.parseMethod();
                auto compiled = compiler.compile(*dynamic_cast<const MethodNode*>(ast.get()));
                std::cout << "✓ " << expr << " -> " << ast->toString() << std::endl;
            } catch (const std::exception& e) {
                std::cout << "✗ " << expr << " failed: " << e.what() << std::endl;
            }
        }
        
        std::cout << "\n=== All block parsing tests completed successfully! ===" << std::endl;
        
        return 0;
        
    } catch (const std::exception& e) {
        std::cerr << "Test FAILED with exception: " << e.what() << std::endl;
        return 1;
    }
}