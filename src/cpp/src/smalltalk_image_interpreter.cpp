#include "interpreter.h"
#include "memory_manager.h"
#include "simple_compiler.h"
#include "simple_parser.h"
#include "smalltalk_image.h"
#include <iostream>

namespace smalltalk {

// Implementation of evaluate() method that creates temporary interpreter
// This replaces the implementation in smalltalk_image.cpp to break circular
// dependency
TaggedValue SmalltalkImage::evaluate(const std::string &code) {
  // For backward compatibility, create a temporary interpreter
  MemoryManager memoryManager;
  Interpreter interpreter(memoryManager, *this);
  return evaluate(code, interpreter);
}

// Implementation of evaluate() method that uses provided interpreter
TaggedValue SmalltalkImage::evaluate(const std::string &code,
                                     Interpreter &interpreter) {
  try {
    SimpleParser parser(code);
    auto methodAST = parser.parseMethod();

    SimpleCompiler compiler;
    auto compiledMethod = compiler.compile(*methodAST);

    return interpreter.executeCompiledMethod(*compiledMethod);

  } catch (const std::exception &e) {
    std::cerr << "Error evaluating code: " << e.what() << std::endl;
    return TaggedValue::nil();
  }
}

} // namespace smalltalk