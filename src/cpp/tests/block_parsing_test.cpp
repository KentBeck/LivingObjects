#include "simple_parser.h"
#include "simple_compiler.h"
#include "simple_vm.h"
#include "ast.h"
#include "tagged_value.h"

#include <iostream>
#include <cassert>

using namespace smalltalk;

int main() {
    try {
        std::cout << "Testing Block parsing and compilation..." << std::endl;
        
        // Test 1: Parse a simple block expression
        std::cout << "\n=== Test 1: Parsing [3 + 4] ===" << std::endl;
        
        SimpleParser parser("[3 + 4]");
        auto methodAST = parser.parseMethod();
        
        std::cout << "Parsed AST: " << methodAST->toString() << std::endl;
        
        // The AST should be a method containing a block
        const MethodNode* method = dynamic_cast<const MethodNode*>(methodAST.get());
        assert(method != nullptr);
        
        const BlockNode* block = dynamic_cast<const BlockNode*>(method->getBody());
        assert(block != nullptr);
        
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
            if (bytecodes[i] == static_cast<uint8_t>(Bytecode::CREATE_BLOCK)) {
                hasCreateBlock = true;
                std::cout << "Found CREATE_BLOCK bytecode at position " << i << std::endl;
                break;
            }
        }
        assert(hasCreateBlock);
        
        // Test 3: Execute the compiled method (this should create a block object)
        std::cout << "\n=== Test 3: Executing compiled method ===" << std::endl;
        
        SimpleVM vm;
        TaggedValue result = vm.execute(*compiledMethod);
        
        std::cout << "Execution result: ";
        if (result.isPointer()) {
            std::cout << "Object@" << result.asPointer() << std::endl;
            
            // Check if the result is actually a BlockContext
            Object* obj = result.asObject();
            if (obj->header.type == static_cast<uint64_t>(ContextType::BLOCK_CONTEXT)) {
                std::cout << "Success! Created a BlockContext object" << std::endl;
                
                BlockContext* blockCtx = static_cast<BlockContext*>(obj);
                std::cout << "Block home: " << blockCtx->home << std::endl;
                std::cout << "Block method ref: " << blockCtx->header.hash << std::endl;
            } else {
                std::cout << "Warning: Object is not a BlockContext (type: " << obj->header.type << ")" << std::endl;
            }
        } else if (result.isInteger()) {
            std::cout << result.asInteger() << std::endl;
        } else if (result.isNil()) {
            std::cout << "nil" << std::endl;
        } else {
            std::cout << "Unknown type" << std::endl;
        }
        
        // Test 4: Test nested expressions in blocks
        std::cout << "\n=== Test 4: Parsing nested block [2 * (3 + 1)] ===" << std::endl;
        
        SimpleParser parser2("[2 * (3 + 1)]");
        auto methodAST2 = parser2.parseMethod();
        
        std::cout << "Parsed AST: " << methodAST2->toString() << std::endl;
        
        // Test 5: Test error handling for malformed blocks
        std::cout << "\n=== Test 5: Error handling ===" << std::endl;
        
        try {
            SimpleParser parser3("[3 + 4"); // Missing closing bracket
            auto methodAST3 = parser3.parseMethod();
            assert(false && "Should have thrown an error for malformed block");
        } catch (const std::exception& e) {
            std::cout << "Correctly caught error for malformed block: " << e.what() << std::endl;
        }
        
        try {
            SimpleParser parser4("3 + 4]"); // Missing opening bracket
            auto methodAST4 = parser4.parseMethod();
            // This should parse as "3 + 4" followed by an unexpected ']'
            std::cout << "Parsing '3 + 4]': " << methodAST4->toString() << std::endl;
        } catch (const std::exception& e) {
            std::cout << "Error parsing '3 + 4]': " << e.what() << std::endl;
        }
        
        std::cout << "\n=== All block parsing tests completed! ===" << std::endl;
        
        return 0;
        
    } catch (const std::exception& e) {
        std::cerr << "Test FAILED with exception: " << e.what() << std::endl;
        return 1;
    }
}