#include "simple_parser.h"
#include "simple_compiler.h"
#include "interpreter.h"
#include "memory_manager.h"
#include "smalltalk_vm.h"
#include "smalltalk_image.h"
#include <iostream>

using namespace smalltalk;

int main() {
    // Initialize VM
    SmalltalkVM::initialize();
    
    MemoryManager memoryManager;
    SmalltalkImage image;
    
    // Test #(1 2 3) size
    SimpleParser parser("#(1 2 3) size");
    auto methodAST = parser.parseMethod();
    
    SimpleCompiler compiler;
    auto compiledMethod = compiler.compile(*methodAST);
    CompiledMethod* rawCompiledMethod = compiledMethod.get();
    
    image.addCompiledMethod(std::move(compiledMethod));
    
    Interpreter interpreter(memoryManager, image);
    TaggedValue result = interpreter.executeCompiledMethod(*rawCompiledMethod);
    
    std::cout << "Result: " << result.asInteger() << std::endl;
    
    // Now debug what's in the array
    parser = SimpleParser("#(1 2 3)");
    methodAST = parser.parseMethod();
    compiler = SimpleCompiler();
    compiledMethod = compiler.compile(*methodAST);
    rawCompiledMethod = compiledMethod.get();
    image.addCompiledMethod(std::move(compiledMethod));
    
    Interpreter interpreter2(memoryManager, image);
    TaggedValue arrayResult = interpreter2.executeCompiledMethod(*rawCompiledMethod);
    
    Object* arrayObj = arrayResult.asObject();
    std::cout << "Array object header size: " << arrayObj->header.size << std::endl;
    
    Class* arrayClass = arrayObj->getClass();
    if (arrayClass) {
        std::cout << "Array class instance size: " << arrayClass->getInstanceSize() << std::endl;
        std::cout << "Array class name: " << arrayClass->getName() << std::endl;
    }
    
    return 0;
}