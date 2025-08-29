#include <iostream>
#include <string>

int main() {
  std::string input = "Array new: 3";
  std::cout << "Input: '" << input << "'" << std::endl;
  std::cout << "Length: " << input.size() << std::endl;

  for (size_t i = 0; i < input.size(); ++i) {
    std::cout << "Position " << i << ": '" << input[i] << "'" << std::endl;
  }

  std::cout << "Position 6 should be: '" << input[6] << "'" << std::endl;
  return 0;
}