#include <iostream>

// Simple test to verify our understanding
int main() {
  std::cout << "=== Smalltalk BNF Parsing Analysis ===" << std::endl;
  std::cout << std::endl;

  std::cout << "Input: 'Array new: 3'" << std::endl;
  std::cout << "Expected parsing according to BNF:" << std::endl;
  std::cout << "1. parseKeywordMessage() calls parseBinaryMessage()"
            << std::endl;
  std::cout << "2. parseBinaryMessage() calls parseUnary()" << std::endl;
  std::cout
      << "3. parseUnary() parses 'Array', sees 'new:', restores pos to 'n'"
      << std::endl;
  std::cout
      << "4. parseBinaryMessage() sees 'n' (not binary op), returns 'Array'"
      << std::endl;
  std::cout << "5. parseKeywordMessage() continues from 'n', parses 'new: 3'"
            << std::endl;
  std::cout
      << "6. Result: MessageSend(receiver='Array', selector='new:', args=['3'])"
      << std::endl;
  std::cout << std::endl;

  std::cout << "Current error: 'Parse error at position 6: Unexpected "
               "characters at end of input'"
            << std::endl;
  std::cout << "Position 6 = 'n' in 'new'" << std::endl;
  std::cout
      << "This suggests parsing stopped after 'Array ' and didn't continue."
      << std::endl;
  std::cout << std::endl;

  std::cout << "Possible issues:" << std::endl;
  std::cout << "1. parseKeywordMessage() while loop not entered" << std::endl;
  std::cout
      << "2. parseKeywordMessage() while loop entered but exits prematurely"
      << std::endl;
  std::cout << "3. Some other method in chain not working as expected"
            << std::endl;

  return 0;
}