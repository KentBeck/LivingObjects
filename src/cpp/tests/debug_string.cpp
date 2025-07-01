#include "smalltalk_string.h"
#include "smalltalk_class.h"
#include "tagged_value.h"
#include <iostream>

using namespace smalltalk;

int main() {
    std::cout << "=== String Debug Test ===" << std::endl;
    
    try {
        // Initialize the class system
        ClassUtils::initializeCoreClasses();
        
        std::cout << "Class system initialized" << std::endl;
        
        // Get String class
        Class* stringClass = ClassUtils::getStringClass();
        std::cout << "String class: " << (stringClass ? stringClass->getName() : "null") << std::endl;
        
        // Create a String object directly
        String* testString = new String("hello", stringClass);
        std::cout << "Created String object: " << testString->toString() << std::endl;
        std::cout << "String content: '" << testString->getContent() << "'" << std::endl;
        std::cout << "String size: " << testString->size() << std::endl;
        
        // Test StringUtils
        std::cout << "\nTesting StringUtils:" << std::endl;
        String* utilString = StringUtils::createString("world");
        std::cout << "StringUtils::createString: " << utilString->toString() << std::endl;
        
        // Test TaggedValue creation
        std::cout << "\nTesting TaggedValue creation:" << std::endl;
        TaggedValue taggedString = StringUtils::createTaggedString("test");
        std::cout << "TaggedValue created" << std::endl;
        
        // Test TaggedValue properties
        std::cout << "taggedString.isPointer(): " << taggedString.isPointer() << std::endl;
        std::cout << "taggedString.isInteger(): " << taggedString.isInteger() << std::endl;
        std::cout << "taggedString.isBoolean(): " << taggedString.isBoolean() << std::endl;
        std::cout << "taggedString.isNil(): " << taggedString.isNil() << std::endl;
        
        // Test StringUtils recognition
        std::cout << "StringUtils::isString(taggedString): " << StringUtils::isString(taggedString) << std::endl;
        
        if (StringUtils::isString(taggedString)) {
            String* extractedString = StringUtils::asString(taggedString);
            std::cout << "Extracted string: " << extractedString->toString() << std::endl;
        } else {
            std::cout << "TaggedValue is not recognized as a string" << std::endl;
        }
        
        // Test class comparison
        Object* obj = taggedString.asObject();
        Class* objClass = obj->getClass();
        std::cout << "Object class: " << (objClass ? objClass->getName() : "null") << std::endl;
        std::cout << "String class == Object class: " << (stringClass == objClass) << std::endl;
        
    } catch (const std::exception& e) {
        std::cout << "Error: " << e.what() << std::endl;
    }
    
    return 0;
}