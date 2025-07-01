#include "tagged_value.h"
#include <iostream>

using namespace smalltalk;

int main() {
    std::cout << "=== TaggedValue Boolean Debug Test ===" << std::endl;
    
    // Create boolean values
    TaggedValue trueVal = TaggedValue::trueValue();
    TaggedValue falseVal = TaggedValue::falseValue();
    TaggedValue nilVal = TaggedValue::nil();
    TaggedValue intVal = TaggedValue(42);
    
    std::cout << "Testing trueValue():" << std::endl;
    std::cout << "  isBoolean(): " << trueVal.isBoolean() << std::endl;
    std::cout << "  isTrue(): " << trueVal.isTrue() << std::endl;
    std::cout << "  isFalse(): " << trueVal.isFalse() << std::endl;
    std::cout << "  isNil(): " << trueVal.isNil() << std::endl;
    std::cout << "  isInteger(): " << trueVal.isInteger() << std::endl;
    std::cout << "  isPointer(): " << trueVal.isPointer() << std::endl;
    std::cout << "  isSpecial(): " << trueVal.isSpecial() << std::endl;
    std::cout << "  rawValue(): " << std::hex << trueVal.rawValue() << std::dec << std::endl;
    
    std::cout << "\nTesting falseValue():" << std::endl;
    std::cout << "  isBoolean(): " << falseVal.isBoolean() << std::endl;
    std::cout << "  isTrue(): " << falseVal.isTrue() << std::endl;
    std::cout << "  isFalse(): " << falseVal.isFalse() << std::endl;
    std::cout << "  isNil(): " << falseVal.isNil() << std::endl;
    std::cout << "  isInteger(): " << falseVal.isInteger() << std::endl;
    std::cout << "  isPointer(): " << falseVal.isPointer() << std::endl;
    std::cout << "  isSpecial(): " << falseVal.isSpecial() << std::endl;
    std::cout << "  rawValue(): " << std::hex << falseVal.rawValue() << std::dec << std::endl;
    
    std::cout << "\nTesting nil():" << std::endl;
    std::cout << "  isBoolean(): " << nilVal.isBoolean() << std::endl;
    std::cout << "  isTrue(): " << nilVal.isTrue() << std::endl;
    std::cout << "  isFalse(): " << nilVal.isFalse() << std::endl;
    std::cout << "  isNil(): " << nilVal.isNil() << std::endl;
    std::cout << "  isInteger(): " << nilVal.isInteger() << std::endl;
    std::cout << "  isPointer(): " << nilVal.isPointer() << std::endl;
    std::cout << "  isSpecial(): " << nilVal.isSpecial() << std::endl;
    std::cout << "  rawValue(): " << std::hex << nilVal.rawValue() << std::dec << std::endl;
    
    std::cout << "\nTesting integer(42):" << std::endl;
    std::cout << "  isBoolean(): " << intVal.isBoolean() << std::endl;
    std::cout << "  isTrue(): " << intVal.isTrue() << std::endl;
    std::cout << "  isFalse(): " << intVal.isFalse() << std::endl;
    std::cout << "  isNil(): " << intVal.isNil() << std::endl;
    std::cout << "  isInteger(): " << intVal.isInteger() << std::endl;
    std::cout << "  isPointer(): " << intVal.isPointer() << std::endl;
    std::cout << "  isSpecial(): " << intVal.isSpecial() << std::endl;
    std::cout << "  rawValue(): " << std::hex << intVal.rawValue() << std::dec << std::endl;
    
    std::cout << "\n=== Constants Check ===" << std::endl;
    std::cout << "SPECIAL_NIL: " << std::hex << TaggedValue::SPECIAL_NIL << std::dec << std::endl;
    std::cout << "SPECIAL_TRUE: " << std::hex << TaggedValue::SPECIAL_TRUE << std::dec << std::endl;
    std::cout << "SPECIAL_FALSE: " << std::hex << TaggedValue::SPECIAL_FALSE << std::dec << std::endl;
    
    return 0;
}