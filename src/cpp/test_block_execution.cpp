#include "simple_parser.h"
#include "simple_compiler.h"
#include "interpreter.h"
#include "memory_manager.h"
#include "smalltalk_class.h"
#include "primitive_methods.h"
#include <iostream>

using namespace smalltalk;

void testExpression(const std::string& expr) {
    std::cout << "Testing: " << expr << std::endl;
    
    try {
        // Parse, compile, and execute the expression
        SimpleParser parser(expr);
        auto methodAST = parser.parseMethod();
        std::cout << "  Parsed: " << methodAST->toString() << std::endl;

        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);
        std::cout << "  Compiled: " << compiledMethod->toString() << std::endl;

        MemoryManager memoryManager;
        Interpreter interpreter(memoryManager);
        TaggedValue result = interpreter.executeCompiledMethod(*compiledMethod);

        // Convert result to string for display
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

        std::cout << "  Result: " << resultStr << std::endl;
        std::cout << "  ✅ SUCCESS" << std::endl;
    } catch (const std::exception& e) {
        std::cout << "  ❌ ERROR: " << e.what() << std::endl;
    }
    std::cout << std::endl;
}

int main() {
    // Initialize class system and primitives
    ClassUtils::initializeCoreClasses();
    auto& primitiveRegistry = PrimitiveRegistry::getInstance();
    primitiveRegistry.initializeCorePrimitives();

    std::cout << "=== Block Execution Tests ===" << std::endl;
    
    // Test simple block creation (should work)
    testExpression("[3 + 4]");
    
    // Test block with parameter creation (should work)
    testExpression("[:x | x + 1]");
    
    // Test block execution (this is where we expect issues)
    testExpression("[3 + 4] value");
    
    return 0;
}
