#include <iostream>
#include <string>

// Simple test to verify parsing works
int main() {
    std::string testExpression = "'hello' , ' world'";
    std::cout << "Testing parsing of: " << testExpression << std::endl;
    
    // For now, just verify the string is well-formed
    std::cout << "String length: " << testExpression.length() << std::endl;
    std::cout << "Characters:" << std::endl;
    for (size_t i = 0; i < testExpression.length(); ++i) {
        char c = testExpression[i];
        std::cout << "  [" << i << "] = '" << c << "' (ASCII " << (int)c << ")";
        if (c == ' ') std::cout << " [SPACE]";
        if (c == ',') std::cout << " [COMMA]";
        if (c == '\'') std::cout << " [QUOTE]";
        std::cout << std::endl;
    }
    
    // Check for the comma operator specifically
    bool foundComma = false;
    for (char c : testExpression) {
        if (c == ',') {
            foundComma = true;
            break;
        }
    }
    
    if (foundComma) {
        std::cout << "✓ Comma operator found in expression" << std::endl;
    } else {
        std::cout << "✗ Comma operator NOT found in expression" << std::endl;
    }
    
    return 0;
}