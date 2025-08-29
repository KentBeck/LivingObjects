#include "ast.h"
#include "simple_parser.h"
#include "smalltalk_string.h"
#include <cassert>
#include <iostream>
#include <memory>

using namespace smalltalk;

void testBasicBlockParsing() {
  std::cout << "Testing basic block parsing..." << std::endl;

  // Test 1: Simple block with integer
  {
    SimpleParser parser("[42]");
    auto method = parser.parseMethod();
    assert(method != nullptr);

    // The method should contain a single statement
    auto body = method->getBody();
    assert(body != nullptr);

    // Check if it's a block node
    if (auto *seqNode = dynamic_cast<const SequenceNode *>(body)) {
      auto &statements = seqNode->getStatements();
      assert(statements.size() == 1);
      auto *blockNode = dynamic_cast<const BlockNode *>(statements[0].get());
      assert(blockNode != nullptr);
      std::cout << "  ✓ Simple block [42] parsed correctly" << std::endl;
    } else if (auto *blockNode = dynamic_cast<const BlockNode *>(body)) {
      (void)blockNode; // silence unused variable warning in release builds
      std::cout << "  ✓ Simple block [42] parsed correctly" << std::endl;
    } else {
      std::cerr << "  ✗ Failed to parse simple block" << std::endl;
      assert(false);
    }
  }

  // Test 2: Block with arithmetic expression
  {
    SimpleParser parser("[3 + 4]");
    auto method = parser.parseMethod();
    assert(method != nullptr);

    auto body = method->getBody();
    assert(body != nullptr);

    // Should be a block containing a binary operation
    bool foundBlock = false;
    if (auto *seqNode = dynamic_cast<const SequenceNode *>(body)) {
      auto &statements = seqNode->getStatements();
      if (!statements.empty()) {
        if (auto *blockNode =
                dynamic_cast<const BlockNode *>(statements[0].get())) {
          foundBlock = true;
          auto blockBody = blockNode->getBody();
          assert(blockBody != nullptr);
          std::cout << "  ✓ Block with expression [3 + 4] parsed correctly"
                    << std::endl;
        }
      }
    } else if (auto *blockNode = dynamic_cast<const BlockNode *>(body)) {
      foundBlock = true;
      (void)blockNode; // silence unused variable warning in release builds
      auto blockBody = blockNode->getBody();
      assert(blockBody != nullptr);
      std::cout << "  ✓ Block with expression [3 + 4] parsed correctly"
                << std::endl;
    }

    assert(foundBlock);
  }

  // Test 3: Nested blocks
  {
    SimpleParser parser("[[1]]");
    auto method = parser.parseMethod();
    assert(method != nullptr);

    auto body = method->getBody();
    assert(body != nullptr);

    // Should be a block containing another block
    bool foundNestedBlock = false;
    if (auto *outerBlock = dynamic_cast<const BlockNode *>(body)) {
      auto innerBody = outerBlock->getBody();
      if (auto *innerBlock = dynamic_cast<const BlockNode *>(innerBody)) {
        (void)innerBlock; // silence unused variable warning in release builds
        foundNestedBlock = true;
        std::cout << "  ✓ Nested blocks [[1]] parsed correctly" << std::endl;
      } else if (auto *seqNode =
                     dynamic_cast<const SequenceNode *>(innerBody)) {
        auto &statements = seqNode->getStatements();
        if (!statements.empty() &&
            dynamic_cast<const BlockNode *>(statements[0].get())) {
          foundNestedBlock = true;
          std::cout << "  ✓ Nested blocks [[1]] parsed correctly" << std::endl;
        }
      }
    } else if (auto *seqNode = dynamic_cast<const SequenceNode *>(body)) {
      auto &statements = seqNode->getStatements();
      if (!statements.empty()) {
        if (auto *outerBlock =
                dynamic_cast<const BlockNode *>(statements[0].get())) {
          auto innerBody = outerBlock->getBody();
          if (dynamic_cast<const BlockNode *>(innerBody)) {
            foundNestedBlock = true;
            std::cout << "  ✓ Nested blocks [[1]] parsed correctly"
                      << std::endl;
          }
        }
      }
    }

    assert(foundNestedBlock);
  }

  std::cout << "All block parsing tests passed!" << std::endl;
}

void testBlockToString() {
  std::cout << "\nTesting block toString() representation..." << std::endl;

  SimpleParser parser("[5 * 6]");
  auto method = parser.parseMethod();
  assert(method != nullptr);

  std::string methodStr = method->toString();
  std::cout << "  Parsed block representation: " << methodStr << std::endl;

  // The string should contain "Block" somewhere
  assert(methodStr.find("Block") != std::string::npos ||
         methodStr.find("block") != std::string::npos ||
         methodStr.find("[") != std::string::npos);

  std::cout << "  ✓ Block toString() works correctly" << std::endl;
}

int main() {
  try {
    std::cout << "=== Minimal Block Parsing Tests ===" << std::endl;

    testBasicBlockParsing();
    testBlockToString();

    std::cout << "\nTest completed successfully!" << std::endl;
    return 0;
  } catch (const std::exception &e) {
    std::cerr << "Test failed with exception: " << e.what() << std::endl;
    return 1;
  }
}
