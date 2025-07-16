#include "simple_parser.h"
#include <iostream>

using namespace smalltalk;

int main() {
    try {
        SimpleParser parser("'hello' , ' world'");
        auto method = parser.parseMethod();
        std::cout << "Parse successful!" << std::endl;
    } catch (const std::exception& e) {
        std::cout << "Parse error: " << e.what() << std::endl;
    }
    return 0;
}