#include "simple_parser.h"
#include "simple_compiler.h"
#include "interpreter.h"
#include "memory_manager.h"
#include "smalltalk_string.h"
#include "smalltalk_class.h"
#include "primitive_methods.h"
#include "primitives/block.h"
#include "compiled_method.h"

#include <iostream>
#include <string>

using namespace smalltalk;

void printUsage() {
    std::cout << "Usage:" << '\n';
    std::cout << "  smalltalk-vm <expression>" << '\n';
    std::cout << '\n';
    std::cout << "Examples:" << '\n';
    std::cout << "  smalltalk-vm \"42\"" << '\n';
    std::cout << "  smalltalk-vm \"3 + 4\"" << '\n';
    std::cout << "  smalltalk-vm \"(10 - 2) * 3\"" << '\n';
}

int main(int argc, char** argv) {
    if (argc != 2) {
        printUsage();
        return 1;
    }
    
    std::string expression = argv[1];
    
    try {
        // Step 1: Initialize class system and primitives before parsing
        ClassUtils::initializeCoreClasses();
        
        // Initialize primitive registry
        auto& primitiveRegistry = PrimitiveRegistry::getInstance();
        primitiveRegistry.initializeCorePrimitives();
        
        // Add primitive methods to Integer class
        Class* integerClass = ClassUtils::getIntegerClass();
        IntegerClassSetup::addPrimitiveMethods(integerClass);
        
        // Register block primitive
        primitiveRegistry.registerPrimitive(PrimitiveNumbers::BLOCK_VALUE, BlockPrimitives::value);
        
        // Step 2: Parse the expression
        SimpleParser parser(expression);
        auto methodAST = parser.parseMethod();
        
        std::cout << "Parsed: " << methodAST->toString() << '\n';
        
        // Step 3: Compile to bytecode
        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);
        
        std::cout << "Compiled: " << compiledMethod->toString() << '\n';
        
        // Step 4: Execute using unified interpreter
        MemoryManager memoryManager;
        Interpreter interpreter(memoryManager);
        
        // Register block primitive
        primitiveRegistry.registerPrimitive(PrimitiveNumbers::BLOCK_VALUE, BlockPrimitives::value);
        
        // Execute the compiled method
        TaggedValue result = interpreter.executeCompiledMethod(*compiledMethod);
        
        // Step 5: Print the result
        if (StringUtils::isString(result)) {
            String* str = StringUtils::asString(result);
            std::cout << "Result: " << str->toString() << '\n';
        } else {
            std::cout << "Result: " << result << '\n';
        }
        
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << '\n';
        return 1;
    }
    
    return 0;
}