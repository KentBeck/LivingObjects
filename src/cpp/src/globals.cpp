#include "globals.h"
#include "smalltalk_class.h"
#include "symbol.h"

#include <unordered_map>

namespace smalltalk {

namespace {
Object *g_smalltalk = nullptr;
std::unordered_map<std::string, Object *> g_globals;
} // namespace

namespace Globals {

Object *getSmalltalk() { return g_smalltalk; }

void setSmalltalk(Object *dict) { g_smalltalk = dict; }

bool isInitialized() { return g_smalltalk != nullptr; }

Object *get(const std::string &name) {
  auto it = g_globals.find(name);
  if (it != g_globals.end())
    return it->second;
  return nullptr;
}

void set(const std::string &name, Object *obj) { g_globals[name] = obj; }

} // namespace Globals

} // namespace smalltalk
