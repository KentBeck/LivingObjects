#include "simple_parser.h"
#include <iostream>

using namespace smalltalk;

int main() {
  std::cout << "Testing string concatenation parsing..." << std::endl;

  try {
    SimpleParser parser("'hello' , ' world'");
    auto method = parser.parseMethod();
    std::cout << "✓ Parse successful!" << std::endl;

    // Test that we can parse without "Unexpected characters" error
    std::cout << "✓ No 'Unexpected characters' error" << std::endl;

    return 0;
  } catch (const std::exception &e) {
    std::cout << "✗ Parse error: " << e.what() << std::endl;
    return 1;
  }
}