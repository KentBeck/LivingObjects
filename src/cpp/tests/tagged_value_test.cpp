#include "../include/tagged_value.h"
#include <iostream>
#include <cassert>

using namespace smalltalk;

// Simple test framework
#define TEST(name) void name()
#define EXPECT_EQ(expected, actual) assert((expected) == (actual))
#define EXPECT_NE(expected, actual) assert((expected) != (actual))

TEST(TestTaggedValueInteger) {
    // Test creating integer 3 - our target expression!
    TaggedValue three(3);
    
    // Verify it's recognized as an integer
    EXPECT_EQ(true, three.isInteger());
    EXPECT_EQ(false, three.isPointer());
    EXPECT_EQ(false, three.isSpecial());
    EXPECT_EQ(false, three.isFloat());
    
    // Verify the value can be extracted
    EXPECT_EQ(3, three.asInteger());
}

TEST(TestTaggedValueIntegerRange) {
    // Test various integer values
    TaggedValue zero(0);
    TaggedValue positive(42);
    TaggedValue negative(-17);
    TaggedValue large(1000000);
    
    EXPECT_EQ(0, zero.asInteger());
    EXPECT_EQ(42, positive.asInteger());
    EXPECT_EQ(-17, negative.asInteger());
    EXPECT_EQ(1000000, large.asInteger());
    
    // All should be integers
    EXPECT_EQ(true, zero.isInteger());
    EXPECT_EQ(true, positive.isInteger());
    EXPECT_EQ(true, negative.isInteger());
    EXPECT_EQ(true, large.isInteger());
}

TEST(TestTaggedValueSpecialValues) {
    // Test nil, true, false
    TaggedValue nil = TaggedValue::nil();
    TaggedValue trueVal = TaggedValue::trueValue();
    TaggedValue falseVal = TaggedValue::falseValue();
    
    EXPECT_EQ(true, nil.isNil());
    EXPECT_EQ(true, trueVal.isTrue());
    EXPECT_EQ(true, falseVal.isFalse());
    
    EXPECT_EQ(true, nil.isSpecial());
    EXPECT_EQ(true, trueVal.isSpecial());
    EXPECT_EQ(true, falseVal.isSpecial());
}

TEST(TestTaggedValueEquality) {
    // Test equality comparison
    TaggedValue three1(3);
    TaggedValue three2(3);
    TaggedValue four(4);
    
    EXPECT_EQ(true, three1 == three2);
    EXPECT_EQ(false, three1 == four);
    EXPECT_EQ(true, three1 != four);
    EXPECT_EQ(false, three1 != three2);
}

TEST(TestTaggedValueOutput) {
    // Test that we can print values (for debugging)
    TaggedValue three(3);
    TaggedValue nil = TaggedValue::nil();
    TaggedValue trueVal = TaggedValue::trueValue();
    
    // Just verify these don't crash
    std::cout << "Testing output: " << three << ", " << nil << ", " << trueVal << std::endl;
}

void runAllTests() {
    std::cout << "Running tagged value tests..." << std::endl;
    
    TestTaggedValueInteger();
    std::cout << "✓ Tagged value integer test passed (expression '3' works!)" << std::endl;
    
    TestTaggedValueIntegerRange();
    std::cout << "✓ Tagged value integer range test passed" << std::endl;
    
    TestTaggedValueSpecialValues();
    std::cout << "✓ Tagged value special values test passed" << std::endl;
    
    TestTaggedValueEquality();
    std::cout << "✓ Tagged value equality test passed" << std::endl;
    
    TestTaggedValueOutput();
    std::cout << "✓ Tagged value output test passed" << std::endl;
    
    std::cout << "All tagged value tests passed! 🎉" << std::endl;
    std::cout << "The Smalltalk expression '3' now works in C++!" << std::endl;
}

int main() {
    try {
        runAllTests();
        return 0;
    } catch (const std::exception& e) {
        std::cerr << "Test failed: " << e.what() << std::endl;
        return 1;
    } catch (...) {
        std::cerr << "Test failed with unknown exception" << std::endl;
        return 1;
    }
}
