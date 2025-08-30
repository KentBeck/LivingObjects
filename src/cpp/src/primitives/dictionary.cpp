#include "../../include/globals.h"
#include "../../include/interpreter.h"
#include "../../include/primitives.h"
#include "../../include/smalltalk_class.h"
#include "../../include/smalltalk_exception.h"
#include "../../include/symbol.h"

namespace smalltalk {

namespace {
inline Object **namedSlots(Object *obj) {
  return reinterpret_cast<Object **>(reinterpret_cast<char *>(obj) +
                                     sizeof(Object));
}

Object *ensureArray(Object *dict, size_t slotIndex, Interpreter &interpreter) {
  Object **slots = namedSlots(dict);
  Object *arr = slots[slotIndex];
  if (arr == nullptr) {
    Class *arrayClass = ClassRegistry::getInstance().getClass("Array");
    if (!arrayClass)
      throw PrimitiveFailure("Array class not found");
    arr =
        interpreter.getMemoryManager().allocateIndexableInstance(arrayClass, 0);
    slots[slotIndex] = arr;
  }
  return arr;
}

inline size_t arraySize(Object *arr) { return arr->header.size; }
inline Object **arraySlots(Object *arr) {
  return reinterpret_cast<Object **>(reinterpret_cast<char *>(arr) +
                                     sizeof(Object));
}

ssize_t indexOfKey(Object *keysArr, Symbol *keySym) {
  size_t n = arraySize(keysArr);
  Object **ks = arraySlots(keysArr);
  for (size_t i = 0; i < n; ++i) {
    if (ks[i] == keySym)
      return static_cast<ssize_t>(i);
  }
  return -1;
}

void appendKV(Object *dict, Object *&keysArr, Object *&valsArr, Symbol *keySym,
              Object *valObj, Interpreter &interpreter) {
  size_t n = arraySize(keysArr);
  Class *arrayClass = ClassRegistry::getInstance().getClass("Array");
  Object *newKeys = interpreter.getMemoryManager().allocateIndexableInstance(
      arrayClass, n + 1);
  Object *newVals = interpreter.getMemoryManager().allocateIndexableInstance(
      arrayClass, n + 1);
  Object **oldK = arraySlots(keysArr);
  Object **oldV = arraySlots(valsArr);
  Object **newK = arraySlots(newKeys);
  Object **newV = arraySlots(newVals);
  for (size_t i = 0; i < n; ++i) {
    newK[i] = oldK[i];
    newV[i] = oldV[i];
  }
  newK[n] = keySym;
  newV[n] = valObj;
  Object **slots = namedSlots(dict);
  slots[0] = newKeys;
  slots[1] = newVals;
}
} // namespace

namespace DictionaryPrimitives {

TaggedValue at(TaggedValue receiver, const std::vector<TaggedValue> &args,
               Interpreter &interpreter) {
  Primitives::checkArgumentCount(args, 1, "Dictionary>>at:");
  Object *dict =
      Primitives::ensureReceiverIsObject(receiver, "Dictionary>>at:");
  if (!args[0].isPointer() ||
      args[0].asObject()->header.getType() != ObjectType::SYMBOL) {
    throw PrimitiveFailure("Dictionary key must be a Symbol");
  }
  Symbol *keySym = static_cast<Symbol *>(args[0].asObject());
  Object *keysArr = ensureArray(dict, 0, interpreter);
  Object *valsArr = ensureArray(dict, 1, interpreter);
  ssize_t idx = indexOfKey(keysArr, keySym);
  if (idx < 0)
    return TaggedValue::nil();
  Object **vals = arraySlots(valsArr);
  Object *valObj = vals[static_cast<size_t>(idx)];
  return valObj ? TaggedValue::fromObject(valObj) : TaggedValue::nil();
}

TaggedValue atPut(TaggedValue receiver, const std::vector<TaggedValue> &args,
                  Interpreter &interpreter) {
  Primitives::checkArgumentCount(args, 2, "Dictionary>>at:put:");
  Object *dict =
      Primitives::ensureReceiverIsObject(receiver, "Dictionary>>at:put:");
  if (!args[0].isPointer() ||
      args[0].asObject()->header.getType() != ObjectType::SYMBOL) {
    throw PrimitiveFailure("Dictionary key must be a Symbol");
  }
  Symbol *keySym = static_cast<Symbol *>(args[0].asObject());
  Object *valObj = args[1].isPointer() ? args[1].asObject() : nullptr;

  Object *keysArr = ensureArray(dict, 0, interpreter);
  Object *valsArr = ensureArray(dict, 1, interpreter);
  ssize_t idx = indexOfKey(keysArr, keySym);
  if (idx >= 0) {
    Object **vals = arraySlots(valsArr);
    vals[static_cast<size_t>(idx)] = valObj;
  } else {
    appendKV(dict, keysArr, valsArr, keySym, valObj, interpreter);
  }

  if (Globals::getSmalltalk() == dict && valObj != nullptr) {
    Globals::set(keySym->getName(), valObj);
  }
  return args[1];
}

TaggedValue keys(TaggedValue receiver, const std::vector<TaggedValue> &args,
                 Interpreter &interpreter) {
  Primitives::checkArgumentCount(args, 0, "Dictionary>>keys");
  Object *dict =
      Primitives::ensureReceiverIsObject(receiver, "Dictionary>>keys");
  Object *keysArr = ensureArray(dict, 0, interpreter);
  return TaggedValue::fromObject(keysArr);
}

TaggedValue size(TaggedValue receiver, const std::vector<TaggedValue> &args,
                 Interpreter &interpreter) {
  Primitives::checkArgumentCount(args, 0, "Dictionary>>size");
  Object *dict =
      Primitives::ensureReceiverIsObject(receiver, "Dictionary>>size");
  Object *keysArr = ensureArray(dict, 0, interpreter);
  return TaggedValue::fromSmallInteger(
      static_cast<int32_t>(arraySize(keysArr)));
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
