#include "simple_parser.h"
#include "simple_compiler.h"
#include <iostream>

using namespace smalltalk;

int main()
{
    try
    {
        // Test simple block
        std::cout << "Testing [3 + 4]..." << std::endl;
        SimpleParser parser1("[3 + 4]");
        auto ast1 = parser1.parseMethod();
        std::cout << "Parsed: " << ast1->toString() << std::endl;

        // Test compilation
        SimpleCompiler compiler1;
        auto compiled1 = compiler1.compile(*ast1);
        std::cout << "Compiled: " << compiled1->toString() << std::endl;

        // Test block with parameter
        std::cout << "\nTesting [:x | x + 1]..." << std::endl;
        SimpleParser parser2("[:x | x + 1]");
        auto ast2 = parser2.parseMethod();
        std::cout << "Parsed: " << ast2->toString() << std::endl;

        // Test compilation
        SimpleCompiler compiler2;
        auto compiled2 = compiler2.compile(*ast2);
        std::cout << "Compiled: " << compiled2->toString() << std::endl;

        std::cout << "\nParser and compiler tests passed!" << std::endl;
        return 0;
    }
    catch (const std::exception &e)
    {
        std::cout << "Error: " << e.what() << std::endl;
        return 1;
    }
}
