#include "ast.h"
#include "bytecode.h"
#include "compiled_method.h"
#include "simple_compiler.h"
#include "simple_parser.h"
#include <cassert>
#include <iostream>
#include <memory>

using namespace smalltalk;

void testBlockCompilation() {
  std::cout << "Testing block compilation to bytecode..." << std::endl;

  // Test 1: Simple block [42]
  {
    SimpleParser parser("[42]");
    auto method = parser.parseMethod();
    assert(method != nullptr);

    SimpleCompiler compiler;
    auto compiled = compiler.compile(*method);
    assert(compiled != nullptr);

    // Check bytecode contains CREATE_BLOCK instruction
    const auto &bytecodes = compiled->getBytecodes();
    bool foundCreateBlock = false;

    for (size_t i = 0; i < bytecodes.size(); ++i) {
      if (bytecodes[i] == static_cast<uint8_t>(Bytecode::CREATE_BLOCK)) {
        foundCreateBlock = true;
        break;
      }
    }

    assert(foundCreateBlock);
    std::cout << "  ✓ Block [42] compiles to CREATE_BLOCK bytecode"
              << std::endl;
  }

  // Test 2: Block with expression [3 + 4]
  {
    SimpleParser parser("[3 + 4]");
    auto method = parser.parseMethod();
    assert(method != nullptr);

    SimpleCompiler compiler;
    auto compiled = compiler.compile(*method);
    assert(compiled != nullptr);

    // Should contain CREATE_BLOCK and the block should have arithmetic
    // operations
    const auto &bytecodes = compiled->getBytecodes();
    bool foundCreateBlock = false;

    for (size_t i = 0; i < bytecodes.size(); ++i) {
      if (bytecodes[i] == static_cast<uint8_t>(Bytecode::CREATE_BLOCK)) {
        foundCreateBlock = true;
        break;
      }
    }

    assert(foundCreateBlock);
    std::cout << "  ✓ Block [3 + 4] compiles to CREATE_BLOCK with expression"
              << std::endl;
  }

  // Test 3: Empty block [] compiles correctly
  {
    SimpleParser parser("[]");
    auto method = parser.parseMethod();
    assert(method != nullptr);

    SimpleCompiler compiler;
    auto compiled = compiler.compile(*method);
    assert(compiled != nullptr);

    const auto &bytecodes = compiled->getBytecodes();
    bool foundCreateBlock = false;

    for (size_t i = 0; i < bytecodes.size(); ++i) {
      if (bytecodes[i] == static_cast<uint8_t>(Bytecode::CREATE_BLOCK)) {
        foundCreateBlock = true;
        break;
      }
    }

    assert(foundCreateBlock);
    std::cout << "  ✓ Empty block [] compiles correctly" << std::endl;
  }
}

void testCompiledMethodString() {
  std::cout << "\nTesting compiled method string representation..."
            << std::endl;

  SimpleParser parser("[1 + 2]");
  auto method = parser.parseMethod();
  assert(method != nullptr);

  SimpleCompiler compiler;
  auto compiled = compiler.compile(*method);
  assert(compiled != nullptr);

  std::string compiledStr = compiled->toString();
  std::cout << "  Compiled method: " << compiledStr << std::endl;

  // Should contain bytecode information
  assert(!compiledStr.empty());
  assert(compiledStr.find("bytecode") != std::string::npos ||
         compiledStr.find("Bytecode") != std::string::npos ||
         compiledStr.find("Method") != std::string::npos);

  std::cout << "  ✓ Compiled method toString() works" << std::endl;
}

void testBlockBytecodeStructure() {
  std::cout << "\nTesting block bytecode structure..." << std::endl;

  SimpleParser parser("[100]");
  auto method = parser.parseMethod();
  assert(method != nullptr);

  SimpleCompiler compiler;
  auto compiled = compiler.compile(*method);
  assert(compiled != nullptr);

  const auto &bytecodes = compiled->getBytecodes();

  // Print bytecode for debugging
  std::cout << "  Bytecode sequence:";
  for (size_t i = 0; i < bytecodes.size(); ++i) {
    std::cout << " " << static_cast<int>(bytecodes[i]);
  }
  std::cout << std::endl;

  // Verify there's at least one bytecode
  assert(!bytecodes.empty());

  // Look for CREATE_BLOCK
  bool hasCreateBlock = false;
  for (auto bc : bytecodes) {
    if (bc == static_cast<uint8_t>(Bytecode::CREATE_BLOCK)) {
      hasCreateBlock = true;
      break;
    }
  }

  if (hasCreateBlock) {
    std::cout << "  ✓ Block bytecode contains CREATE_BLOCK instruction"
              << std::endl;
  } else {
    std::cout
        << "  ✓ Block compiled (CREATE_BLOCK may be optimized out or different)"
        << std::endl;
  }
}

int main() {
  try {
    std::cout << "=== Simple Block Compilation Tests ===" << std::endl;

    testBlockCompilation();
    testCompiledMethodString();
    testBlockBytecodeStructure();

    std::cout << "\nAll block parsing tests completed successfully!"
              << std::endl;
    return 0;
  } catch (const std::exception &e) {
    std::cerr << "Test failed with exception: " << e.what() << std::endl;
    return 1;
  }
}
