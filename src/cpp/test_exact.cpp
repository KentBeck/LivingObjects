#include <iostream>
#include <string>

int main() {
  // Test the exact string from the test file
  std::string test = "'hello' , ' world'";
  std::cout << "String: '" << test << "'" << std::endl;
  std::cout << "Length: " << test.length() << std::endl;

  // Check each character
  for (size_t i = 0; i < test.length(); ++i) {
    char c = test[i];
    std::cout << "pos[" << i << "] = '" << c << "' (ASCII " << (int)c << ")";
    if (c == ' ')
      std::cout << " [SPACE]";
    if (c == '\t')
      std::cout << " [TAB]";
    if (c == '\n')
      std::cout << " [NEWLINE]";
    std::cout << std::endl;
  }

  return 0;
}