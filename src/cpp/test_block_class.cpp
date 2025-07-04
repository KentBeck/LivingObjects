#include "smalltalk_class.h"
#include "symbol.h"
#include <iostream>

using namespace smalltalk;

int main() {
    // Initialize class system
    ClassUtils::initializeCoreClasses();
    
    // Get the Block class
    Class* blockClass = ClassRegistry::getInstance().getClass("Block");
    if (blockClass == nullptr) {
        std::cout << "❌ Block class not found!" << std::endl;
        return 1;
    }
    
    std::cout << "✅ Block class found: " << blockClass->getName() << std::endl;
    
    // Check if the 'value' method exists
    Symbol* valueSymbol = Symbol::intern("value");
    auto valueMethod = blockClass->lookupMethod(valueSymbol);
    
    if (valueMethod == nullptr) {
        std::cout << "❌ 'value' method not found in Block class!" << std::endl;
        return 1;
    }
    
    std::cout << "✅ 'value' method found with primitive number: " << valueMethod->primitiveNumber << std::endl;
    
    return 0;
}
