#include "simple_interpreter.h"

#include <iostream>
#include <string>

using namespace smalltalk;

int main() {
    SimpleInterpreter interpreter;
    std::string input;
    
    std::cout << "🎯 Smalltalk C++ Interpreter v0.1" << '\n';
    std::cout << "Currently supports: integers (3, -42), special values (nil, true, false)" << '\n';
    std::cout << "Type 'quit' to exit." << '\n';
    std::cout << '\n';
    
    while (true) {
        std::cout << "st> ";
        std::getline(std::cin, input);
        
        // Check for quit command
        if (input == "quit" || input == "exit") {
            std::cout << "Goodbye! 👋" << '\n';
            break;
        }
        
        // Skip empty lines
        if (input.empty()) {
            continue;
        }
        
        try {
            TaggedValue result = interpreter.evaluate(input);
            std::cout << "=> " << result << '\n';
        } catch (const std::exception& e) {
            std::cout << "Error: " << e.what() << '\n';
        }
        
        std::cout << '\n';
    }
    
    return 0;
}