// Minimal bootstrap API to avoid circular dependencies
#pragma once

#include <string>

namespace smalltalk {

class Class;
class MemoryManager;

// Register methods to be installed during prepareImage (no sends now).
void registerSmalltalkMethod(Class *clazz, const std::string &methodSource);
void registerPrimitiveMethod(Class *clazz, const std::string &selector,
                             int primitiveNumber);

// Perform explicit image preparation: ensure dictionaries exist and install
// all pending primitive and Smalltalk methods via MemoryManager without sends.
void prepareImage(MemoryManager &mm);

} // namespace smalltalk

