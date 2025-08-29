#pragma once

#include "object.h"
#include <string>

namespace smalltalk {

namespace Globals {
// Accessor for the global Smalltalk dictionary object (may be nullptr early)
Object *getSmalltalk();
void setSmalltalk(Object *dict);

// Get/set globals by name in the Smalltalk dictionary
Object *get(const std::string &name);
void set(const std::string &name, Object *obj);

// Utility: check if Smalltalk dictionary is initialized
bool isInitialized();
} // namespace Globals

} // namespace smalltalk
