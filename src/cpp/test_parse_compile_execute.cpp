#include <iostream>
#include <cassert>
#include "include/interpreter.h"
#include "include/memory_manager.h"
#include "include/smalltalk_class.h"
#include "include/primitives.h"
#include "include/simple_parser.h"
#include "include/simple_compiler.h"
#include "include/compiled_method.h"

using namespace smalltalk;

void test_parse_literal() {
    std::cout << "Testing literal parsing..." << std::endl;
    
    SimpleParser parser("42");
    auto method = parser.parseMethod();
    assert(method != nullptr);
    
    auto literal = dynamic_cast<const LiteralNode*>(method->getBody());
    assert(literal != nullptr);
    assert(literal->getValue().isSmallInteger());
    assert(literal->getValue().getSmallInteger() == 42);
    
    std::cout << "✓ Literal parsing works" << std::endl;
}

void test_parse_object_new() {
    std::cout << "Testing 'Object new' parsing..." << std::endl;
    
    // Initialize system first so classes are registered
    MemoryManager memory;
    Interpreter interpreter(memory);
    
    SimpleParser parser("Object new");
    auto method = parser.parseMethod();
    assert(method != nullptr);
    
    auto messageSend = dynamic_cast<const MessageSendNode*>(method->getBody());
    assert(messageSend != nullptr);
    
    // Check receiver is "Object"
    auto receiver = dynamic_cast<const LiteralNode*>(messageSend->getReceiver());
    assert(receiver != nullptr);
    
    // Check selector is "new"
    assert(messageSend->getSelector() == "new");
    
    // Check no arguments
    assert(messageSend->getArguments().empty());
    
    std::cout << "✓ 'Object new' parsing works" << std::endl;
}

void test_compile_literal() {
    std::cout << "Testing literal compilation..." << std::endl;
    
    SimpleParser parser("42");
    auto method = parser.parseMethod();
    
    SimpleCompiler compiler;
    auto compiledMethod = compiler.compile(*method);
    
    assert(compiledMethod != nullptr);
    assert(!compiledMethod->getBytecodes().empty());
    assert(!compiledMethod->getLiterals().empty());
    
    // Should have PUSH_LITERAL and RETURN_STACK_TOP
    const auto& bytecodes = compiledMethod->getBytecodes();
    assert(bytecodes[0] == static_cast<uint8_t>(Bytecode::PUSH_LITERAL));
    
    std::cout << "✓ Literal compilation works" << std::endl;
}

void test_compile_object_new() {
    std::cout << "Testing 'Object new' compilation..." << std::endl;
    
    // Initialize system first so classes are registered
    MemoryManager memory;
    Interpreter interpreter(memory);
    
    SimpleParser parser("Object new");
    auto method = parser.parseMethod();
    
    SimpleCompiler compiler;
    auto compiledMethod = compiler.compile(*method);
    
    assert(compiledMethod != nullptr);
    
    const auto& bytecodes = compiledMethod->getBytecodes();
    assert(!bytecodes.empty());
    
    // Should contain PUSH_LITERAL (for Object), SEND_MESSAGE, RETURN_STACK_TOP
    bool foundPushLiteral = false;
    bool foundSendMessage = false;
    bool foundReturn = false;
    
    for (size_t i = 0; i < bytecodes.size(); i++) {
        if (bytecodes[i] == static_cast<uint8_t>(Bytecode::PUSH_LITERAL)) {
            foundPushLiteral = true;
        } else if (bytecodes[i] == static_cast<uint8_t>(Bytecode::SEND_MESSAGE)) {
            foundSendMessage = true;
        } else if (bytecodes[i] == static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP)) {
            foundReturn = true;
        }
    }
    
    assert(foundPushLiteral);
    assert(foundSendMessage);
    assert(foundReturn);
    
    std::cout << "✓ 'Object new' compilation works" << std::endl;
}

void test_execute_literal() {
    std::cout << "Testing literal execution..." << std::endl;
    
    // Create memory manager and interpreter to initialize system
    MemoryManager memory;
    Interpreter interpreter(memory);
    
    SimpleParser parser("42");
    auto method = parser.parseMethod();
    
    SimpleCompiler compiler;
    auto compiledMethod = compiler.compile(*method);
    
    TaggedValue result = interpreter.executeCompiledMethod(*compiledMethod);
    
    assert(result.isSmallInteger());
    assert(result.getSmallInteger() == 42);
    
    std::cout << "✓ Literal execution works" << std::endl;
}

void test_execute_object_new() {
    std::cout << "Testing 'Object new' execution..." << std::endl;
    
    // Create memory manager and interpreter to initialize system
    MemoryManager memory;
    Interpreter interpreter(memory);
    
    SimpleParser parser("Object new");
    auto method = parser.parseMethod();
    
    SimpleCompiler compiler;
    auto compiledMethod = compiler.compile(*method);
    
    try {
        TaggedValue result = interpreter.executeCompiledMethod(*compiledMethod);
        
        assert(result.isPointer());
        Object* obj = result.asObject();
        assert(obj != nullptr);
        assert(obj->header.getType() == ObjectType::OBJECT);
        
        // Get Object class and verify
        ClassRegistry& registry = ClassRegistry::getInstance();
        Class* objectClass = registry.getClass("Object");
        assert(obj->getClass() == objectClass);
        
        std::cout << "✓ 'Object new' execution works!" << std::endl;
        std::cout << "🎯 FULL PARSE → COMPILE → EXECUTE PIPELINE WORKING!" << std::endl;
        
    } catch (const std::exception& e) {
        std::cerr << "❌ Error executing 'Object new': " << e.what() << std::endl;
        throw;
    }
}

void test_arithmetic_through_pipeline() {
    std::cout << "Testing arithmetic through full pipeline..." << std::endl;
    
    MemoryManager memory;
    Interpreter interpreter(memory);
    
    SimpleParser parser("3 + 4");
    auto method = parser.parseMethod();
    
    SimpleCompiler compiler;
    auto compiledMethod = compiler.compile(*method);
    
    try {
        TaggedValue result = interpreter.executeCompiledMethod(*compiledMethod);
        
        assert(result.isSmallInteger());
        assert(result.getSmallInteger() == 7);
        
        std::cout << "✓ Arithmetic pipeline works (3 + 4 = 7)" << std::endl;
        
    } catch (const std::exception& e) {
        std::cerr << "❌ Error executing '3 + 4': " << e.what() << std::endl;
        throw;
    }
}

int main() {
    std::cout << "=== Parse → Compile → Execute Pipeline Test ===" << std::endl;
    std::cout << "===============================================" << std::endl;
    
    try {
        // Test each stage individually
        test_parse_literal();
        test_parse_object_new();
        test_compile_literal();
        test_compile_object_new();
        test_execute_literal();
        
        // Test full pipeline
        test_arithmetic_through_pipeline();
        test_execute_object_new();
        
        std::cout << "\n🎉 All pipeline tests passed!" << std::endl;
        std::cout << "\nACHIEVEMENTS:" << std::endl;
        std::cout << "✅ Parse: Handles literals and 'Object new' message sends" << std::endl;
        std::cout << "✅ Compile: Generates correct bytecode for message sends" << std::endl;
        std::cout << "✅ Execute: Runs bytecode and calls primitives correctly" << std::endl;
        std::cout << "✅ Pipeline: Full parse → compile → execute flow working" << std::endl;
        std::cout << "✅ Objects: 'Object new' creates real objects with correct types" << std::endl;
        std::cout << "✅ Integration: Parser, compiler, and VM work together" << std::endl;
        
        std::cout << "\n🚀 Parser extension successful!" << std::endl;
        
    } catch (const std::exception& e) {
        std::cerr << "\n❌ Test failed: " << e.what() << std::endl;
        return 1;
    }
    
    return 0;
}