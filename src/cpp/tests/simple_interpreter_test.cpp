#include "../include/simple_interpreter.h"
#include <iostream>
#include <cassert>

using namespace smalltalk;

// Simple test framework
#define TEST(name) void name()
#define EXPECT_EQ(expected, actual) assert((expected) == (actual))
#define EXPECT_TRUE(condition) assert(condition)
#define EXPECT_FALSE(condition) assert(!(condition))

TEST(TestEvaluateInteger3) {
    // This is our target: make "3" work!
    SimpleInterpreter interpreter;
    
    TaggedValue result = interpreter.evaluate("3");
    
    EXPECT_TRUE(result.isInteger());
    EXPECT_EQ(3, result.asInteger());
    
    std::cout << "✨ SUCCESS: '3' evaluates to " << result << std::endl;
}

TEST(TestEvaluateVariousIntegers) {
    SimpleInterpreter interpreter;
    
    // Test various integer expressions
    TaggedValue zero = interpreter.evaluate("0");
    TaggedValue positive = interpreter.evaluate("42");
    TaggedValue negative = interpreter.evaluate("-17");
    TaggedValue large = interpreter.evaluate("1000000");
    
    EXPECT_TRUE(zero.isInteger());
    EXPECT_EQ(0, zero.asInteger());
    
    EXPECT_TRUE(positive.isInteger());
    EXPECT_EQ(42, positive.asInteger());
    
    EXPECT_TRUE(negative.isInteger());
    EXPECT_EQ(-17, negative.asInteger());
    
    EXPECT_TRUE(large.isInteger());
    EXPECT_EQ(1000000, large.asInteger());
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
    
    // Test that whitespace is handled correctly
    TaggedValue result1 = interpreter.evaluate("  3  ");
    TaggedValue result2 = interpreter.evaluate("\t42\n");
    TaggedValue result3 = interpreter.evaluate(" nil ");
    
    EXPECT_TRUE(result1.isInteger());
    EXPECT_EQ(3, result1.asInteger());
    
    EXPECT_TRUE(result2.isInteger());
    EXPECT_EQ(42, result2.asInteger());
    
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
        std::cout << "Expected error for 'invalid': " << e.what() << std::endl;
    }
    EXPECT_TRUE(caught);
}

void runAllTests() {
    std::cout << "Running simple interpreter tests..." << std::endl;
    
    TestEvaluateInteger3();
    std::cout << "✓ Evaluate integer '3' test passed" << std::endl;
    
    TestEvaluateVariousIntegers();
    std::cout << "✓ Evaluate various integers test passed" << std::endl;
    
    TestEvaluateSpecialValues();
    std::cout << "✓ Evaluate special values test passed" << std::endl;
    
    TestEvaluateWithWhitespace();
    std::cout << "✓ Evaluate with whitespace test passed" << std::endl;
    
    TestEvaluateInvalidExpression();
    std::cout << "✓ Evaluate invalid expression test passed" << std::endl;
    
    std::cout << "All simple interpreter tests passed! 🚀" << std::endl;
    std::cout << "The C++ Smalltalk interpreter can now evaluate '3'!" << std::endl;
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
