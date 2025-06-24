#include "../include/simple_interpreter.h"
#include <iostream>
#include <string>

using namespace smalltalk;

int main() {
    SimpleInterpreter interpreter;
    std::string input;
    
    std::cout << "🎯 Smalltalk C++ Interpreter v0.1" << std::endl;
    std::cout << "Currently supports: integers (3, -42), special values (nil, true, false)" << std::endl;
    std::cout << "Type 'quit' to exit." << std::endl;
    std::cout << std::endl;
    
    while (true) {
        std::cout << "st> ";
        std::getline(std::cin, input);
        
        // Check for quit command
        if (input == "quit" || input == "exit") {
            std::cout << "Goodbye! 👋" << std::endl;
            break;
        }
        
        // Skip empty lines
        if (input.empty()) {
            continue;
        }
        
        try {
            TaggedValue result = interpreter.evaluate(input);
            std::cout << "=> " << result << std::endl;
        } catch (const std::exception& e) {
            std::cout << "Error: " << e.what() << std::endl;
        }
        
        std::cout << std::endl;
    }
    
    return 0;
}
