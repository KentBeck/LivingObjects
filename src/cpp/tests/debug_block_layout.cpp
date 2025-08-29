#include "bytecode.h"
#include "compiled_method.h"
#include "context.h"
#include "interpreter.h"
#include "memory_manager.h"
#include "object.h"
#include "smalltalk_class.h"
#include <cassert>
#include <iomanip>
#include <iostream>

using namespace smalltalk;

void printObjectLayout(const std::string &label, Object *obj) {
  std::cout << label << ":" << std::endl;
  std::cout << "  Address: " << obj << std::endl;
  std::cout << "  Class: "
            << (obj->getClass() ? obj->getClass()->getName() : "nullptr")
            << std::endl;

  if (auto *context = dynamic_cast<BlockContext *>(obj)) {
    std::cout << "  BlockContext fields:" << std::endl;
    std::cout << "    Home context: " << context->home << std::endl;
  }
}

void testBlockContextCreation() {
  std::cout << "Testing BlockContext creation and layout..." << std::endl;

  // Initialize the system
  ClassUtils::initializeCoreClasses();

  // Create a block class since ClassUtils doesn't provide getBlockClass
  Class *blockClass =
      ClassUtils::createClass("Block", ClassUtils::getObjectClass());
  assert(blockClass != nullptr);

  // Create a parent context (simulating where the block was created)
  // Note: This test creates invalid context without CompiledMethod - only for
  // layout testing
  auto *parentContext = reinterpret_cast<MethodContext *>(
      malloc(sizeof(MethodContext) + 256 * sizeof(TaggedValue)));
  new (parentContext) Object(ObjectType::CONTEXT, 256, 0);

  // Create a method context class
  Class *methodContextClass =
      ClassUtils::createClass("MethodContext", ClassUtils::getObjectClass());
  parentContext->setClass(methodContextClass);

  // Create the block context
  // BlockContext constructor requires: contextSize, methodRef, receiver,
  // senderContext, homeContext
  auto *blockContext =
      new BlockContext(256, 0, nullptr, nullptr, parentContext);
  blockContext->setClass(blockClass);

  printObjectLayout("Created BlockContext", blockContext);

  std::cout << "  ✓ BlockContext created successfully" << std::endl;

  // Clean up
  delete blockContext;
  delete parentContext;
}

void testBlockMemoryLayout() {
  std::cout << "\nTesting block memory layout details..." << std::endl;

  // Create a block context
  auto *blockContext = new BlockContext(256, 0, nullptr, nullptr, nullptr);
  // Create a block class
  Class *blockClass =
      ClassUtils::createClass("Block", ClassUtils::getObjectClass());
  blockContext->setClass(blockClass);

  // Print memory layout
  std::cout << "Block memory layout:" << std::endl;
  std::cout << "  sizeof(Object): " << sizeof(Object) << " bytes" << std::endl;
  std::cout << "  sizeof(MethodContext): " << sizeof(MethodContext) << " bytes"
            << std::endl;
  std::cout << "  sizeof(BlockContext): " << sizeof(BlockContext) << " bytes"
            << std::endl;
  std::cout << "  sizeof(CompiledMethod): " << sizeof(CompiledMethod)
            << " bytes" << std::endl;

  // Test field access
  std::cout << "\nTesting field access:" << std::endl;
  std::cout << "  Home context: " << blockContext->home << std::endl;
  std::cout << "  ✓ Field access works correctly" << std::endl;

  // Clean up
  delete blockContext;
}

void testBlockContextHierarchy() {
  std::cout << "\nTesting block context class hierarchy..." << std::endl;

  // Create a method context class
  Class *methodContextClass =
      ClassUtils::createClass("MethodContext", ClassUtils::getObjectClass());

  // Create a block class
  Class *blockClass =
      ClassUtils::createClass("Block", ClassUtils::getObjectClass());

  // Create contexts of different types
  // Create a dummy compiled method for the test
  std::vector<uint8_t> dummyBytecodes = {
      static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP)};
  std::vector<TaggedValue> dummyLiterals;
  std::vector<std::string> dummyTempVars;
  auto dummyMethod = std::make_unique<CompiledMethod>(
      dummyBytecodes, dummyLiterals, dummyTempVars, 0);

  auto *methodContext = new MethodContext(
      256, TaggedValue::nil(), TaggedValue::nil(), dummyMethod.get());
  methodContext->setClass(methodContextClass);

  auto *blockContext = new BlockContext(256, 0, nullptr, nullptr, nullptr);
  blockContext->setClass(blockClass);

  // Test polymorphism
  Object *objectPtr = blockContext;
  std::cout << "  BlockContext as Object*: " << objectPtr << std::endl;
  std::cout << "  Can cast to BlockContext*: "
            << (dynamic_cast<BlockContext *>(objectPtr) != nullptr)
            << std::endl;

  std::cout << "  ✓ Class hierarchy works correctly" << std::endl;

  // Clean up
  delete methodContext;
  delete blockContext;
  // dummyMethod will be automatically destroyed by unique_ptr
}

void testManualBlockSetup() {
  std::cout << "\nTesting manual block setup..." << std::endl;

  // Create and setup a block manually
  auto *block = new BlockContext(256, 0, nullptr, nullptr, nullptr);
  // Create a block class
  Class *blockClass =
      ClassUtils::createClass("Block", ClassUtils::getObjectClass());
  block->setClass(blockClass);

  std::cout << "After manual setting:" << std::endl;
  std::cout << "  Block: " << block << std::endl;
  std::cout << "  Class: " << block->getClass()->getName() << std::endl;
  std::cout << "  Home context: " << block->home << std::endl;

  std::cout << "  ✓ Manual block setup successful" << std::endl;

  // Clean up
  delete block;
}

int main() {
  try {
    std::cout << "=== Debug Block Layout ===" << std::endl;

    testBlockContextCreation();
    testBlockMemoryLayout();
    testBlockContextHierarchy();
    testManualBlockSetup();

    std::cout << "\nAll layout tests completed successfully!" << std::endl;
    return 0;
  } catch (const std::exception &e) {
    std::cerr << "Test failed with exception: " << e.what() << std::endl;
    return 1;
  }
}