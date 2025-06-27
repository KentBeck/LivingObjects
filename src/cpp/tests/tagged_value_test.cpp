#include "tagged_value.h"

#include <cassert>
#include <iostream>

using namespace smalltalk;

// Simple test framework
#define TEST(name) void name()
#define EXPECT_EQ(expected, actual) assert((expected) == (actual))
#define EXPECT_NE(expected, actual) assert((expected) != (actual))

TEST(TestTaggedValueInteger) {
    const int TEST_INTEGER_THREE = 3;
    // Test creating integer 3 - our target expression!
    TaggedValue three(TEST_INTEGER_THREE);
    
    // Verify it's recognized as an integer
    EXPECT_EQ(true, three.isInteger());
    EXPECT_EQ(false, three.isPointer());
    EXPECT_EQ(false, three.isSpecial());
    EXPECT_EQ(false, three.isFloat());
    
    // Verify the value can be extracted
    EXPECT_EQ(TEST_INTEGER_THREE, three.asInteger());
}

TEST(TestTaggedValueIntegerRange) {
    const int TEST_INTEGER_ZERO = 0;
    const int TEST_INTEGER_POSITIVE = 42;
    const int TEST_INTEGER_NEGATIVE = -17;
    const int TEST_INTEGER_LARGE = 1000000;
    // Test various integer values
    TaggedValue zero(TEST_INTEGER_ZERO);
    TaggedValue positive(TEST_INTEGER_POSITIVE);
    TaggedValue negative(TEST_INTEGER_NEGATIVE);
    TaggedValue large(TEST_INTEGER_LARGE);
    
    EXPECT_EQ(TEST_INTEGER_ZERO, zero.asInteger());
    EXPECT_EQ(TEST_INTEGER_POSITIVE, positive.asInteger());
    EXPECT_EQ(TEST_INTEGER_NEGATIVE, negative.asInteger());
    EXPECT_EQ(TEST_INTEGER_LARGE, large.asInteger());
    
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
    const int TEST_EQUALITY_THREE = 3;
    const int TEST_EQUALITY_FOUR = 4;
    // Test equality comparison
    TaggedValue three1(TEST_EQUALITY_THREE);
    TaggedValue three2(TEST_EQUALITY_THREE);
    TaggedValue four(TEST_EQUALITY_FOUR);
    
    EXPECT_EQ(true, three1 == three2);
    EXPECT_EQ(false, three1 == four);
    EXPECT_EQ(true, three1 != four);
    EXPECT_EQ(false, three1 != three2);
}

TEST(TestTaggedValueOutput) {
    const int TEST_OUTPUT_THREE = 3;
    // Test that we can print values (for debugging)
    TaggedValue three(TEST_OUTPUT_THREE);
    TaggedValue nil = TaggedValue::nil();
    TaggedValue trueVal = TaggedValue::trueValue();
    
    // Just verify these don't crash
    std::cout << "Testing output: " << three << ", " << nil << ", " << trueVal << '\n';
}

void runAllTests() {
    std::cout << "Running tagged value tests..." << '\n';
    
    TestTaggedValueInteger();
    std::cout << "✓ Tagged value integer test passed (expression '3' works!)" << '\n';
    
    TestTaggedValueIntegerRange();
    std::cout << "✓ Tagged value integer range test passed" << '\n';
    
    TestTaggedValueSpecialValues();
    std::cout << "✓ Tagged value special values test passed" << '\n';
    
    TestTaggedValueEquality();
    std::cout << "✓ Tagged value equality test passed" << '\n';
    
    TestTaggedValueOutput();
    std::cout << "✓ Tagged value output test passed" << '\n';
    
    std::cout << "All tagged value tests passed! 🎉" << '\n';
    std::cout << "The Smalltalk expression '3' now works in C++!" << '\n';
}

int main() {
    try {
        runAllTests();
        return 0;
    } catch (const std::exception& e) {
        std::cerr << "Test failed: " << e.what() << '\n';
        return 1;
    } catch (...) {
        std::cerr << "Test failed with unknown exception" << '\n';
        return 1;
    }
}