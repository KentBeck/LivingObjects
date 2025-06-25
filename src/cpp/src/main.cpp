#include "simple_parser.h"
#include "simple_compiler.h"
#include "simple_vm.h"
#include <iostream>
#include <string>

using namespace smalltalk;

void printUsage() {
    std::cout << "Usage:" << std::endl;
    std::cout << "  smalltalk-vm <expression>" << std::endl;
    std::cout << std::endl;
    std::cout << "Examples:" << std::endl;
    std::cout << "  smalltalk-vm \"42\"" << std::endl;
    std::cout << "  smalltalk-vm \"3 + 4\"" << std::endl;
    std::cout << "  smalltalk-vm \"(10 - 2) * 3\"" << std::endl;
}

int main(int argc, char** argv) {
    if (argc != 2) {
        printUsage();
        return 1;
    }
    
    std::string expression = argv[1];
    
    try {
        // Step 1: Parse the expression
        SimpleParser parser(expression);
        auto methodAST = parser.parseMethod();
        
        std::cout << "Parsed: " << methodAST->toString() << std::endl;
        
        // Step 2: Compile to bytecode
        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);
        
        std::cout << "Compiled: " << compiledMethod->toString() << std::endl;
        
        // Step 3: Execute the bytecode
        SimpleVM vm;
        TaggedValue result = vm.execute(*compiledMethod);
        
        // Step 4: Print the result
        std::cout << "Result: " << result << std::endl;
        
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        return 1;
    }
    
    return 0;
}