#include "bytecode.h"
#include "compiled_method.h"
#include "interpreter.h"
#include "memory_manager.h"
#include "primitive_methods.h"
#include "primitives/block.h"
#include "smalltalk_class.h"
#include "smalltalk_image.h"
#include "smalltalk_vm.h"
#include "symbol.h"
#include "tagged_value.h"
#include <cassert>
#include <iostream>
#include <vector>

using namespace smalltalk;

void testBlockValuePrimitive() {
  std::cout << "Testing block value primitive..." << std::endl;

  // Initialize the VM (classes + primitives)
  SmalltalkVM::initialize();

  // Create a simple block that returns 42
  MemoryManager memoryManager;
  SmalltalkImage image;
  Interpreter interpreter(memoryManager, image);

  // Create a compiled method for the block body
  auto blockMethod = std::make_unique<CompiledMethod>();

  // Block body: push 42, return
  blockMethod->addLiteral(TaggedValue(42));
  size_t literalIndex = blockMethod->getLiterals().size() - 1;

  blockMethod->addBytecode(static_cast<uint8_t>(Bytecode::PUSH_LITERAL));
  // Add 4 bytes for the literal index (PUSH_LITERAL expects a 32-bit operand)
  blockMethod->addBytecode(literalIndex & 0xFF);
  blockMethod->addBytecode((literalIndex >> 8) & 0xFF);
  blockMethod->addBytecode((literalIndex >> 16) & 0xFF);
  blockMethod->addBytecode((literalIndex >> 24) & 0xFF);

  blockMethod->addBytecode(static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP));

  // Test the block execution
  try {
    // Execute the block method directly
    TaggedValue result = interpreter.executeCompiledMethod(*blockMethod);

    assert(result.isInteger());
    assert(result.asInteger() == 42);

    std::cout << "  ✓ Block value primitive returns correct result: " << result
              << std::endl;
  } catch (const std::exception &e) {
    std::cerr << "  ✗ Block execution failed: " << e.what() << std::endl;
    assert(false);
  }
}

void testBlockWithExpression() {
  std::cout << "Testing block with arithmetic expression..." << std::endl;

  // Create a block that computes 3 + 4
  MemoryManager memoryManager;
  SmalltalkImage image;
  Interpreter interpreter(memoryManager, image);

  auto blockMethod = std::make_unique<CompiledMethod>();

  // Block body: push 3, push 4, send +, return
  blockMethod->addLiteral(TaggedValue(3));
  size_t literal3Index = blockMethod->getLiterals().size() - 1;

  blockMethod->addLiteral(TaggedValue(4));
  size_t literal4Index = blockMethod->getLiterals().size() - 1;

  blockMethod->addBytecode(static_cast<uint8_t>(Bytecode::PUSH_LITERAL));
  blockMethod->addBytecode(literal3Index & 0xFF);
  blockMethod->addBytecode((literal3Index >> 8) & 0xFF);
  blockMethod->addBytecode((literal3Index >> 16) & 0xFF);
  blockMethod->addBytecode((literal3Index >> 24) & 0xFF);

  blockMethod->addBytecode(static_cast<uint8_t>(Bytecode::PUSH_LITERAL));
  blockMethod->addBytecode(literal4Index & 0xFF);
  blockMethod->addBytecode((literal4Index >> 8) & 0xFF);
  blockMethod->addBytecode((literal4Index >> 16) & 0xFF);
  blockMethod->addBytecode((literal4Index >> 24) & 0xFF);

  // Add the + selector as a literal
  blockMethod->addLiteral(TaggedValue(Symbol::intern("+")));
  size_t selectorIndex = blockMethod->getLiterals().size() - 1;

  blockMethod->addBytecode(static_cast<uint8_t>(Bytecode::SEND_MESSAGE));
  blockMethod->addBytecode(selectorIndex); // selector index
  blockMethod->addBytecode(0);             // high byte of selector index
  blockMethod->addBytecode(0);             // unused
  blockMethod->addBytecode(0);             // unused
  blockMethod->addBytecode(1);             // arg count
  blockMethod->addBytecode(0);             // high byte of arg count
  blockMethod->addBytecode(0);             // unused
  blockMethod->addBytecode(0);             // unused

  blockMethod->addBytecode(static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP));

  try {
    TaggedValue result = interpreter.executeCompiledMethod(*blockMethod);

    assert(result.isInteger());
    assert(result.asInteger() == 7);

    std::cout << "  ✓ Block with expression [3 + 4] returns: " << result
              << std::endl;
  } catch (const std::exception &e) {
    std::cerr << "  ✗ Block expression execution failed: " << e.what()
              << std::endl;
    assert(false);
  }
}

void testEmptyBlock() {
  std::cout << "Testing empty block..." << std::endl;

  MemoryManager memoryManager;
  SmalltalkImage image;
  Interpreter interpreter(memoryManager, image);

  auto blockMethod = std::make_unique<CompiledMethod>();

  // Empty block should return nil
  blockMethod->addLiteral(TaggedValue::nil());
  size_t nilIndex = blockMethod->getLiterals().size() - 1;

  blockMethod->addBytecode(static_cast<uint8_t>(Bytecode::PUSH_LITERAL));
  blockMethod->addBytecode(nilIndex & 0xFF);
  blockMethod->addBytecode((nilIndex >> 8) & 0xFF);
  blockMethod->addBytecode((nilIndex >> 16) & 0xFF);
  blockMethod->addBytecode((nilIndex >> 24) & 0xFF);

  blockMethod->addBytecode(static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP));

  try {
    TaggedValue result = interpreter.executeCompiledMethod(*blockMethod);

    assert(result.isNil());

    std::cout << "  ✓ Empty block returns nil: " << result << std::endl;
  } catch (const std::exception &e) {
    std::cerr << "  ✗ Empty block execution failed: " << e.what() << std::endl;
    assert(false);
  }
}

int main() {
  try {
    std::cout << "=== Minimal Block Runtime Tests ===" << std::endl;

    testBlockValuePrimitive();
    testBlockWithExpression();
    testEmptyBlock();

    std::cout << "\nAll tests PASSED!" << std::endl;
    return 0;
  } catch (const std::exception &e) {
    std::cerr << "Test suite failed with exception: " << e.what() << std::endl;
    return 1;
  }
}
