#include "../../include/globals.h"
#include "../../include/interpreter.h"
#include "../../include/primitives.h"
#include "../../include/smalltalk_class.h"
#include "../../include/smalltalk_exception.h"
#include "../../include/symbol.h"

#include <unordered_map>

namespace smalltalk {

// Per-instance dictionary store keyed by receiver Object*
namespace {
std::unordered_map<Object *, std::unordered_map<std::string, Object *>>
    kDictStore;

std::string keyToString(const TaggedValue &key) {
  if (key.isPointer()) {
    Object *obj = key.asObject();
    if (obj->header.getType() == ObjectType::SYMBOL) {
      Symbol *sym = static_cast<Symbol *>(obj);
      return sym->getName();
    }
    // Strings could be supported here; for now prefer symbols
  }
  throw PrimitiveFailure("Dictionary key must be a Symbol");
}

} // namespace

namespace DictionaryPrimitives {

TaggedValue at(TaggedValue receiver, const std::vector<TaggedValue> &args,
               Interpreter &interpreter) {
  (void)interpreter;
  Primitives::checkArgumentCount(args, 1, "Dictionary>>at:");
  Object *dict =
      Primitives::ensureReceiverIsObject(receiver, "Dictionary>>at:");
  auto &store = kDictStore[dict];
  std::string key = keyToString(args[0]);
  auto it = store.find(key);
  if (it == store.end()) {
    return TaggedValue::nil();
  }
  return TaggedValue::fromObject(it->second);
}

TaggedValue atPut(TaggedValue receiver, const std::vector<TaggedValue> &args,
                  Interpreter &interpreter) {
  (void)interpreter;
  Primitives::checkArgumentCount(args, 2, "Dictionary>>at:put:");
  Object *dict =
      Primitives::ensureReceiverIsObject(receiver, "Dictionary>>at:put:");
  std::string key = keyToString(args[0]);
  Object *valObj = args[1].isPointer() ? args[1].asObject() : nullptr;

  kDictStore[dict][key] = valObj;

  // If this is the global Smalltalk dictionary, mirror into Globals
  if (Globals::getSmalltalk() == dict && valObj != nullptr) {
    Globals::set(key, valObj);
  }

  return args[1];
}

TaggedValue keys(TaggedValue receiver, const std::vector<TaggedValue> &args,
                 Interpreter &interpreter) {
  (void)interpreter;
  Primitives::checkArgumentCount(args, 0, "Dictionary>>keys");
  Object *dict =
      Primitives::ensureReceiverIsObject(receiver, "Dictionary>>keys");
  auto &store = kDictStore[dict];

  // Return an Array of Symbols
  Class *arrayClass = ClassRegistry::getInstance().getClass("Array");
  if (!arrayClass) {
    throw PrimitiveFailure("Array class not found");
  }

  Object *arr = interpreter.getMemoryManager().allocateIndexableInstance(
      arrayClass, store.size());
  Object **slots = reinterpret_cast<Object **>(reinterpret_cast<char *>(arr) +
                                               sizeof(Object));
  size_t i = 0;
  for (auto &kv : store) {
    Symbol *sym = Symbol::intern(kv.first);
    slots[i++] = sym;
  }
  return TaggedValue::fromObject(arr);
}

TaggedValue size(TaggedValue receiver, const std::vector<TaggedValue> &args,
                 Interpreter &interpreter) {
  (void)interpreter;
  Primitives::checkArgumentCount(args, 0, "Dictionary>>size");
  Object *dict =
      Primitives::ensureReceiverIsObject(receiver, "Dictionary>>size");
  auto &store = kDictStore[dict];
  return TaggedValue::fromSmallInteger(static_cast<int32_t>(store.size()));
}

} // namespace DictionaryPrimitives

void registerDictionaryPrimitives() {
  PrimitiveRegistry &registry = PrimitiveRegistry::getInstance();
  registry.registerPrimitive(PrimitiveNumbers::DICT_AT,
                             DictionaryPrimitives::at);
  registry.registerPrimitive(PrimitiveNumbers::DICT_AT_PUT,
                             DictionaryPrimitives::atPut);
  registry.registerPrimitive(PrimitiveNumbers::DICT_KEYS,
                             DictionaryPrimitives::keys);
  registry.registerPrimitive(PrimitiveNumbers::DICT_SIZE,
                             DictionaryPrimitives::size);
}

} // namespace smalltalk
