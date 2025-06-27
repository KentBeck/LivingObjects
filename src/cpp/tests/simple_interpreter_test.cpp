#include "simple_interpreter.h"

#include <cassert>
#include <iostream>

using namespace smalltalk;

// Simple test framework
#define TEST(name) void name()
#define EXPECT_EQ(expected, actual) assert((expected) == (actual))
#define EXPECT_TRUE(condition) assert(condition)
#define EXPECT_FALSE(condition) assert(!(condition))

TEST(TestEvaluateInteger3) {
    const int TEST_INTEGER_THREE = 3;
    // This is our target: make "3" work!
    SimpleInterpreter interpreter;
    
    TaggedValue result = interpreter.evaluate("3");
    
    EXPECT_TRUE(result.isInteger());
    EXPECT_EQ(TEST_INTEGER_THREE, result.asInteger());
    
    std::cout << "✨ SUCCESS: '3' evaluates to " << result << '\n';
}

TEST(TestEvaluateVariousIntegers) {
    SimpleInterpreter interpreter;
    
    const int TEST_INTEGER_ZERO = 0;
    const int TEST_INTEGER_POSITIVE = 42;
    const int TEST_INTEGER_NEGATIVE = -17;
    const int TEST_INTEGER_LARGE = 1000000;

    // Test various integer expressions
    TaggedValue zero = interpreter.evaluate("0");
    TaggedValue positive = interpreter.evaluate("42");
    TaggedValue negative = interpreter.evaluate("-17");
    TaggedValue large = interpreter.evaluate("1000000");
    
    EXPECT_TRUE(zero.isInteger());
    EXPECT_EQ(TEST_INTEGER_ZERO, zero.asInteger());
    
    EXPECT_TRUE(positive.isInteger());
    EXPECT_EQ(TEST_INTEGER_POSITIVE, positive.asInteger());
    
    EXPECT_TRUE(negative.isInteger());
    EXPECT_EQ(TEST_INTEGER_NEGATIVE, negative.asInteger());
    
    EXPECT_TRUE(large.isInteger());
    EXPECT_EQ(TEST_INTEGER_LARGE, large.asInteger());
}

TEST(TestEvaluateSpecialValues) {
    SimpleInterpreter interpreter;
    
    // Test special values
    TaggedValue nil = interpreter.evaluate("nil");
    TaggedValue trueVal = interpreter.evaluate("true");
    TaggedValue falseVal = interpreter.evaluate("false");
    
    EXPECT_TRUE(nil.isNil());
    EXPECT_TRUE(trueVal.isTrue());
    EXPECT_TRUE(falseVal.isFalse());
}

TEST(TestEvaluateWithWhitespace) {
    SimpleInterpreter interpreter;
    
    const int WHITESPACE_TEST_THREE = 3;
    const int WHITESPACE_TEST_FORTY_TWO = 42;

    // Test that whitespace is handled correctly
    TaggedValue result1 = interpreter.evaluate("  3  ");
    TaggedValue result2 = interpreter.evaluate("\t42\n");
    TaggedValue result3 = interpreter.evaluate(" nil ");
    
    EXPECT_TRUE(result1.isInteger());
    EXPECT_EQ(WHITESPACE_TEST_THREE, result1.asInteger());
    
    EXPECT_TRUE(result2.isInteger());
    EXPECT_EQ(WHITESPACE_TEST_FORTY_TWO, result2.asInteger());

    
    EXPECT_TRUE(result3.isNil());
}

TEST(TestEvaluateInvalidExpression) {
    SimpleInterpreter interpreter;
    
    // Test that invalid expressions throw exceptions
    bool caught = false;
    try {
        interpreter.evaluate("invalid");
    } catch (const std::runtime_error& e) {
        caught = true;
        std::cout << "Expected error for 'invalid': " << e.what() << '\n';
    }
    EXPECT_TRUE(caught);
}

void runAllTests() {
    std::cout << "Running simple interpreter tests..." << '\n';
    
    TestEvaluateInteger3();
    std::cout << "✓ Evaluate integer '3' test passed" << '\n';
    
    TestEvaluateVariousIntegers();
    std::cout << "✓ Evaluate various integers test passed" << '\n';
    
    TestEvaluateSpecialValues();
    std::cout << "✓ Evaluate special values test passed" << '\n';
    
    TestEvaluateWithWhitespace();
    std::cout << "✓ Evaluate with whitespace test passed" << '\n';
    
    TestEvaluateInvalidExpression();
    std::cout << "✓ Evaluate invalid expression test passed" << '\n';
    
    std::cout << "All simple interpreter tests passed! 🚀" << '\n';
    std::cout << "The C++ Smalltalk interpreter can now evaluate '3'!" << '\n';
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