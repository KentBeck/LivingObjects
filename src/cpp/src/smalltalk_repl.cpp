#include "ast.h"
#include "compiled_method.h"
#include "interpreter.h"
#include "memory_manager.h"
#include "simple_compiler.h"
#include "simple_parser.h"
#include "tagged_value.h"
#include "primitive_methods.h"
#include "primitives/block.h"
#include "smalltalk_class.h"
#include "smalltalk_string.h"
#include "smalltalk_image.h"

#include <iostream>
#include <string>
#include <memory>

using namespace smalltalk;

int main()
{
    // Initialize class system and primitives
    ClassUtils::initializeCoreClasses();

    // Initialize primitive registry
    auto &primitiveRegistry = PrimitiveRegistry::getInstance();
    primitiveRegistry.initializeCorePrimitives();

    // Add primitive methods to Integer class
    Class* integerClass = ClassUtils::getIntegerClass();
    IntegerClassSetup::addPrimitiveMethods(integerClass);

    // Register block primitive
    primitiveRegistry.registerPrimitive(PrimitiveNumbers::BLOCK_VALUE, BlockPrimitives::value);

    // Setup the Smalltalk environment
    MemoryManager memory_manager;
    SmalltalkImage image;
    Interpreter interpreter(memory_manager, image);
    SimpleCompiler compiler;
    std::string input;

    std::cout << "🎯 Smalltalk C++ Bytecode Interpreter v0.2" << '\n';
    std::cout << "Currently supports arithmetic, comparisons, and blocks." << '\n';
    std::cout << "Type 'quit' to exit." << '\n';
    std::cout << '\n';

    while (true)
    {
        std::cout << "st> ";
        std::getline(std::cin, input);

        // Check for quit command
        if (input == "quit" || input == "exit")
        {
            std::cout << "Goodbye! 👋" << '\n';
            break;
        }

        // Skip empty lines
        if (input.empty())
        {
            continue;
        }

        try
        {
            // Parse
            SimpleParser parser(input);
            std::unique_ptr<MethodNode> methodAST = parser.parseMethod();

            // Compile
            std::unique_ptr<CompiledMethod> method = compiler.compile(*methodAST);

            // Execute
            TaggedValue result = interpreter.executeCompiledMethod(*method);

            // Print result
            if (StringUtils::isString(result))
            {
                String *str = StringUtils::asString(result);
                std::cout << "=> " << str->toString() << '\n';
            }
            else
            {
                std::cout << "=> " << result << '\n';
            }
        }
        catch (const std::exception &e)
        {
            std::cout << "Error: " << e.what() << '\n';
        }

        std::cout << '\n';
    }

    return 0;
}