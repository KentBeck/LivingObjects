#include "simple_parser.h"
#include "simple_compiler.h"
#include "interpreter.h"
#include "memory_manager.h"
#include "tagged_value.h"
#include "smalltalk_class.h"
#include "smalltalk_vm.h"
#include "smalltalk_image.h"
#include "compiled_method.h"
#include "primitives.h"

#include <iostream>
#include <cassert>

#define TEST(name) void name()
#define EXPECT_EQ(expected, actual) assert((expected) == (actual))
#define EXPECT_TRUE(expr) assert(expr)
#define EXPECT_FALSE(expr) assert(!(expr))

using namespace smalltalk;

TEST(testPrimitiveMethodCompilation) {
    std::cout << "Testing primitive method compilation..." << std::endl;
    
    MemoryManager memoryManager;
    
    // Test parsing method with primitive syntax
    // Smalltalk syntax: <primitive: 1>
    std::string methodSource = "<primitive: 1> ^self";
    
    try {
        SimpleParser parser(methodSource);
        auto methodAST = parser.parseMethod();
        
        // The method should have been parsed
        EXPECT_TRUE(methodAST != nullptr);
        
        // Compile the method
        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);
        
        // Check that the primitive number was set
        EXPECT_EQ(1, compiledMethod->primitiveNumber);
        
        // The method should still have bytecode for fallback
        EXPECT_TRUE(compiledMethod->bytecodes.size() > 0);
        
        std::cout << "✓ Method with primitive compiled successfully" << std::endl;
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        assert(false);
    }
}

TEST(testPrimitiveMethodExecution) {
    std::cout << "Testing primitive method execution..." << std::endl;
    
    MemoryManager memoryManager;
    SmalltalkVM vm;
    
    // Initialize the VM and primitive registry
    vm.initialize();
    auto& primitiveRegistry = PrimitiveRegistry::getInstance();
    primitiveRegistry.initializeCorePrimitives();
    
    // Create a method with primitive for integer addition
    std::string methodSource = "<primitive: 1> ^self";
    
    try {
        SimpleParser parser(methodSource);
        auto methodAST = parser.parseMethod();
        
        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);
        
        // Create an interpreter
        SmalltalkImage image;
        Interpreter interpreter(memoryManager, image);
        
        // Test primitive execution
        TaggedValue receiver(5);  // Integer constructor
        std::vector<TaggedValue> args = { TaggedValue(3) };
        
        // Check that the primitive exists
        EXPECT_TRUE(primitiveRegistry.hasPrimitive(1));
        
        // Execute the primitive using the registry
        TaggedValue result = primitiveRegistry.callPrimitive(1, receiver, args, interpreter);
        EXPECT_EQ(8, result.asInteger());
        
        std::cout << "✓ Primitive method executed successfully" << std::endl;
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        assert(false);
    }
}

TEST(testPrimitiveMethodWithFallback) {
    std::cout << "Testing primitive method with fallback code..." << std::endl;
    
    MemoryManager memoryManager;
    
    // Method with primitive and fallback code
    std::string methodSource = "<primitive: 999> ^self";
    
    try {
        SimpleParser parser(methodSource);
        auto methodAST = parser.parseMethod();
        
        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);
        
        // Check primitive number
        EXPECT_EQ(999, compiledMethod->primitiveNumber);
        
        // Check that fallback bytecode was generated
        EXPECT_TRUE(compiledMethod->bytecodes.size() > 0);
        
        std::cout << "✓ Primitive method with fallback compiled successfully" << std::endl;
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        assert(false);
    }
}

int main() {
    std::cout << "Running primitive method tests..." << std::endl << std::endl;
    
    testPrimitiveMethodCompilation();
    testPrimitiveMethodExecution();
    testPrimitiveMethodWithFallback();
    
    std::cout << std::endl << "All primitive method tests passed!" << std::endl;
    return 0;
}