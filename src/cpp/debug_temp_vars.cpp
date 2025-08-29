#include "interpreter.h"
#include "memory_manager.h"
#include "simple_compiler.h"
#include "simple_parser.h"
#include "smalltalk_image.h"
#include "smalltalk_vm.h"

#include <iostream>

using namespace smalltalk;

int main() {
  // Initialize VM
  SmalltalkVM::initialize();

  MemoryManager memoryManager;
  SmalltalkImage image;

  try {
    // Test parsing temporary variables
    std::cout << "Testing temporary variable parsing..." << std::endl;
    SimpleParser parser("| x | x := 42. x");
    auto methodAST = parser.parseMethod();

    std::cout << "Parsing successful!" << std::endl;
    std::cout << "Temp vars: ";
    for (const auto &var : methodAST->getTempVars()) {
      std::cout << var << " ";
    }
    std::cout << std::endl;

    // Test compilation
    std::cout << "Testing compilation..." << std::endl;
    SimpleCompiler compiler;
    auto compiledMethod = compiler.compile(*methodAST);

    std::cout << "Compilation successful!" << std::endl;
    std::cout << "Temp vars in compiled method: ";
    for (const auto &var : compiledMethod->getTempVars()) {
      std::cout << var << " ";
    }
    std::cout << std::endl;

    // Test execution
    std::cout << "Testing execution..." << std::endl;
    CompiledMethod *rawCompiledMethod = compiledMethod.get();
    image.addCompiledMethod(std::move(compiledMethod));

    Interpreter interpreter(memoryManager, image);
    TaggedValue result = interpreter.executeCompiledMethod(*rawCompiledMethod);

    std::cout << "Execution successful!" << std::endl;
    if (result.isInteger()) {
      std::cout << "Result: " << result.asInteger() << std::endl;
    } else {
      std::cout << "Result is not an integer" << std::endl;
    }

  } catch (const std::exception &e) {
    std::cout << "Error: " << e.what() << std::endl;
    return 1;
  }

  return 0;
}