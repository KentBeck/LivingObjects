#include "simple_parser.h"
#include "simple_compiler.h"
#include "simple_vm.h"

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
        // Step 1: Parse the expression
        SimpleParser parser(expression);
        auto methodAST = parser.parseMethod();
        
        std::cout << "Parsed: " << methodAST->toString() << '\n';
        
        // Step 2: Compile to bytecode
        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);
        
        std::cout << "Compiled: " << compiledMethod->toString() << '\n';
        
        // Step 3: Execute the bytecode
        SimpleVM vm;
        TaggedValue result = vm.execute(*compiledMethod);
        
        // Step 4: Print the result
        std::cout << "Result: " << result << '\n';
        
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << '\n';
        return 1;
    }
    
    return 0;
}