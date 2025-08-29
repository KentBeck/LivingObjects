#include "compiled_method.h"
#include "interpreter.h"
#include "memory_manager.h"
#include "primitives.h"
#include "simple_compiler.h"
#include "simple_parser.h"
#include "smalltalk_class.h"
#include "smalltalk_image.h"
#include "smalltalk_vm.h"
#include "tagged_value.h"

#include <cassert>
#include <iostream>

#define TEST(name) void name()
#define EXPECT_EQ(expected, actual) assert((expected) == (actual))
#define EXPECT_TRUE(expr) assert(expr)
#define EXPECT_FALSE(expr) assert(!(expr))

using namespace smalltalk;

TEST(testIntegerAdditionWithPrimitive) {
  std::cout << "Testing Integer + with primitive..." << std::endl;

  MemoryManager memoryManager;
  SmalltalkVM vm;
  vm.initialize();

  // Method source with primitive 1 (SmallInteger +) and fallback
  std::string methodSource = "<primitive: 1> | result | result := 0. ^result";

  try {
    SimpleParser parser(methodSource);
    auto methodAST = parser.parseMethod();

    SimpleCompiler compiler;
    auto compiledMethod = compiler.compile(*methodAST);

    // Verify primitive was set
    EXPECT_EQ(1, compiledMethod->primitiveNumber);

    // Create interpreter and test execution
    SmalltalkImage image;
    Interpreter interpreter(memoryManager, image);

    // Test with integers that should use the primitive
    TaggedValue receiver(10);
    std::vector<TaggedValue> args = {TaggedValue(15)};

    auto &primitiveRegistry = PrimitiveRegistry::getInstance();
    primitiveRegistry.initializeCorePrimitives();

    TaggedValue result =
        primitiveRegistry.callPrimitive(1, receiver, args, interpreter);
    EXPECT_EQ(25, result.asInteger());

    std::cout << "✓ Integer addition primitive works: 10 + 15 = "
              << result.asInteger() << std::endl;
  } catch (const std::exception &e) {
    std::cerr << "Error: " << e.what() << std::endl;
    assert(false);
  }
}

TEST(testMultiplePrimitiveMethods) {
  std::cout << "Testing multiple methods with different primitives..."
            << std::endl;

  MemoryManager memoryManager;
  SmalltalkVM vm;
  vm.initialize();

  auto &primitiveRegistry = PrimitiveRegistry::getInstance();
  primitiveRegistry.initializeCorePrimitives();

  SmalltalkImage image;
  Interpreter interpreter(memoryManager, image);

  // Test different arithmetic primitives
  struct PrimitiveTest {
    int primitiveNumber;
    int receiver;
    int arg;
    int expected;
    std::string operation;
  };

  std::vector<PrimitiveTest> tests = {
      {1, 10, 5, 15, "+"}, // Addition
      {2, 10, 3, 7, "-"},  // Subtraction
      {9, 4, 5, 20, "*"},  // Multiplication
      {3, 10, 5, 1, "<"},  // Less than (returns true=1)
      {4, 10, 5, 1, ">"},  // Greater than (returns true=1)
      {7, 10, 10, 1, "="}, // Equality (returns true=1)
  };

  for (const auto &test : tests) {
    std::string methodSource =
        "<primitive: " + std::to_string(test.primitiveNumber) + "> ^nil";

    SimpleParser parser(methodSource);
    auto methodAST = parser.parseMethod();

    SimpleCompiler compiler;
    auto compiledMethod = compiler.compile(*methodAST);

    EXPECT_EQ(test.primitiveNumber, compiledMethod->primitiveNumber);

    TaggedValue receiver(test.receiver);
    std::vector<TaggedValue> args = {TaggedValue(test.arg)};

    TaggedValue result = primitiveRegistry.callPrimitive(
        test.primitiveNumber, receiver, args, interpreter);

    // For comparison operations, check boolean result
    if (test.primitiveNumber >= 3 && test.primitiveNumber <= 8) {
      bool boolResult =
          (test.primitiveNumber == 3)   ? (test.receiver < test.arg)
          : (test.primitiveNumber == 4) ? (test.receiver > test.arg)
          : (test.primitiveNumber == 7) ? (test.receiver == test.arg)
                                        : false;
      EXPECT_TRUE(boolResult ? result.isTrue() : result.isFalse());
      std::cout << "✓ Primitive " << test.primitiveNumber << " ("
                << test.operation << "): " << test.receiver << " "
                << test.operation << " " << test.arg << " = "
                << (boolResult ? "true" : "false") << std::endl;
    } else {
      EXPECT_EQ(test.expected, result.asInteger());
      std::cout << "✓ Primitive " << test.primitiveNumber << " ("
                << test.operation << "): " << test.receiver << " "
                << test.operation << " " << test.arg << " = "
                << result.asInteger() << std::endl;
    }
  }
}

TEST(testPrimitiveWithComplexFallback) {
  std::cout << "Testing primitive with complex fallback code..." << std::endl;

  MemoryManager memoryManager;

  // Method with non-existent primitive that will always fail
  // This tests that the fallback code is properly compiled
  std::string methodSource =
      "<primitive: 9999> | a b | a := 100. b := 200. ^a + b";

  try {
    SimpleParser parser(methodSource);
    auto methodAST = parser.parseMethod();

    SimpleCompiler compiler;
    auto compiledMethod = compiler.compile(*methodAST);

    // Check primitive number was set
    EXPECT_EQ(9999, compiledMethod->primitiveNumber);

    // Check that fallback bytecode was generated (should have more than just
    // RETURN)
    EXPECT_TRUE(compiledMethod->bytecodes.size() > 1);

    // Count the number of literal constants (should have 100 and 200)
    int literalCount = 0;
    for (const auto &literal : compiledMethod->literals) {
      if (literal.isInteger()) {
        int value = literal.asInteger();
        if (value == 100 || value == 200) {
          literalCount++;
        }
      }
    }
    EXPECT_EQ(2, literalCount);

    std::cout << "✓ Complex fallback code compiled with "
              << compiledMethod->bytecodes.size() << " bytecodes and "
              << compiledMethod->literals.size() << " literals" << std::endl;
  } catch (const std::exception &e) {
    std::cerr << "Error: " << e.what() << std::endl;
    assert(false);
  }
}

int main() {
  std::cout << "Running primitive method integration tests..." << std::endl
            << std::endl;

  testIntegerAdditionWithPrimitive();
  testMultiplePrimitiveMethods();
  testPrimitiveWithComplexFallback();

  std::cout << std::endl
            << "All primitive method integration tests passed!" << std::endl;
  return 0;
}