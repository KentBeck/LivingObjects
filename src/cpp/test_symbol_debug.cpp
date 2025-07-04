#include "symbol.h"
#include "tagged_value.h"
#include "object.h"
#include <iostream>

using namespace smalltalk;

int main() {
    try {
        // Create a symbol
        std::cout << "Creating symbol 'value'..." << std::endl;
        Symbol* valueSymbol = Symbol::intern("value");
        std::cout << "Symbol created: " << valueSymbol->toString() << std::endl;
        std::cout << "Symbol name: " << valueSymbol->getName() << std::endl;
        std::cout << "Symbol object type: " << static_cast<int>(valueSymbol->header.getType()) << std::endl;
        std::cout << "Expected SYMBOL type: " << static_cast<int>(ObjectType::SYMBOL) << std::endl;
        
        // Create a TaggedValue from the symbol
        std::cout << "\nCreating TaggedValue from symbol..." << std::endl;
        TaggedValue symbolValue(valueSymbol);
        std::cout << "TaggedValue created" << std::endl;
        std::cout << "TaggedValue isPointer: " << symbolValue.isPointer() << std::endl;
        
        if (symbolValue.isPointer()) {
            Object* obj = symbolValue.asObject();
            std::cout << "Object type from TaggedValue: " << static_cast<int>(obj->header.getType()) << std::endl;
            std::cout << "Is SYMBOL type? " << (obj->header.getType() == ObjectType::SYMBOL) << std::endl;
            
            if (obj->header.getType() == ObjectType::SYMBOL) {
                Symbol* sym = reinterpret_cast<Symbol*>(obj);
                std::cout << "Symbol name from TaggedValue: " << sym->getName() << std::endl;
            }
        }
        
        return 0;
    } catch (const std::exception& e) {
        std::cout << "Error: " << e.what() << std::endl;
        return 1;
    }
}
