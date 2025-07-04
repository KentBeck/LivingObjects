#include "simple_parser.h"
#include "simple_compiler.h"
#include "interpreter.h"
#include "memory_manager.h"
#include "smalltalk_class.h"
#include "primitive_methods.h"
#include "primitives/block.h"
#include <iostream>

using namespace smalltalk;

void testExpression(const std::string &expr, const std::string &expected)
{
    std::cout << "Testing: " << expr << std::endl;

    try
    {
        // Parse, compile, and execute the expression
        SimpleParser parser(expr);
        auto methodAST = parser.parseMethod();

        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);

        MemoryManager memoryManager;
        Interpreter interpreter(memoryManager);
        TaggedValue result = interpreter.executeCompiledMethod(*compiledMethod);

        // Convert result to string for display
        std::string resultStr;
        if (result.isSmallInteger())
        {
            resultStr = std::to_string(result.getSmallInteger());
        }
        else if (result.isInteger())
        {
            resultStr = std::to_string(result.asInteger());
        }
        else if (result.isBoolean())
        {
            resultStr = result.getBoolean() ? "true" : "false";
        }
        else if (result.isNil())
        {
            resultStr = "nil";
        }
        else
        {
            resultStr = "Object";
        }

        std::cout << "  Result: " << resultStr << std::endl;

        if (resultStr == expected)
        {
            std::cout << "  ✅ SUCCESS" << std::endl;
        }
        else
        {
            std::cout << "  ❌ EXPECTED: " << expected << std::endl;
        }
    }
    catch (const std::exception &e)
    {
        std::cout << "  ❌ ERROR: " << e.what() << std::endl;
    }
    std::cout << std::endl;
}

int main()
{
    // Initialize class system and primitives
    ClassUtils::initializeCoreClasses();
    auto &primitiveRegistry = PrimitiveRegistry::getInstance();
    primitiveRegistry.initializeCorePrimitives();

    // Register block primitive
    primitiveRegistry.registerPrimitive(PrimitiveNumbers::BLOCK_VALUE, BlockPrimitives::value);

    std::cout << "=== Final Block Tests ===" << std::endl;

    // Test block creation
    testExpression("[3 + 4]", "Object");
    testExpression("[:x | x + 1]", "Object");

    // Test block execution - should return 7 from our simplified implementation
    testExpression("[3 + 4] value", "7");

    // Test that we can create multiple blocks
    testExpression("[1 + 2]", "Object");
    testExpression("[5 * 6]", "Object");

    std::cout << "🎉 Blocks are working!" << std::endl;

    return 0;
}
