// BootstrapBuilder builds Smalltalk-side structures without sending messages.
#pragma once

#include "memory_manager.h"
#include <string>

namespace smalltalk {

class BootstrapBuilder {
public:
  // Build real Smalltalk MethodDictionaries for every registered class.
  // Allocation uses MemoryManager and performs no message sends.
  static void buildMethodDictionaries(MemoryManager &mm);

  // Register methods to be installed during prepareImage (no sends now).
  static void registerSmalltalkMethod(Class *clazz,
                                      const std::string &methodSource);
  static void registerPrimitiveMethod(Class *clazz, const std::string &selector,
                                      int primitiveNumber);

  // Perform explicit image preparation: ensure dictionaries exist and install
  // all pending primitive and Smalltalk methods via MemoryManager without
  // message sends.
  static void prepareImage(MemoryManager &mm);
};

} // namespace smalltalk
