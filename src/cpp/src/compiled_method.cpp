#include "compiled_method.h"
#include "memory_manager.h"
#include "smalltalk_class.h"
#include "symbol.h"

namespace smalltalk {

static void writeBytes(Object *obj, const std::vector<uint8_t> &data) {
  // Assumes obj is byte-indexable and has enough capacity for data.size()
  char *base = reinterpret_cast<char *>(obj) + sizeof(Object);
  uint8_t *bytes = reinterpret_cast<uint8_t *>(base);
  for (size_t i = 0; i < data.size(); ++i) {
    bytes[i] = data[i];
  }
}

static void setArraySlot(Object *arr, size_t index, Object *value) {
  Object **slots = reinterpret_cast<Object **>(reinterpret_cast<char *>(arr) +
                                               sizeof(Object));
  slots[index] = value;
}

void CompiledMethod::ensureSmalltalkBacking(MemoryManager &mm) {
  // ByteArray for bytecodes
  if (!bytecodesBytes) {
    Class *byteArrayClass = ClassRegistry::getInstance().getClass("ByteArray");
    if (byteArrayClass) {
      bytecodesBytes =
          mm.allocateByteIndexableInstance(byteArrayClass, bytecodes.size());
      writeBytes(bytecodesBytes, bytecodes);
    }
  }

  // Array for literals (boxed where necessary)
  if (!literalsArray) {
    Class *arrayClass = ClassRegistry::getInstance().getClass("Array");
    if (arrayClass) {
      literalsArray = mm.allocateIndexableInstance(arrayClass, literals.size());
      for (size_t i = 0; i < literals.size(); ++i) {
        Object *elem = nullptr;
        if (literals[i].isPointer()) {
          elem = literals[i].asObject();
        } else {
          elem = literals[i].toObject(mm);
        }
        setArraySlot(literalsArray, i, elem);
      }
    }
  }

  // Array of Symbols for temp var names
  if (!tempNamesArray) {
    Class *arrayClass = ClassRegistry::getInstance().getClass("Array");
    if (arrayClass) {
      tempNamesArray =
          mm.allocateIndexableInstance(arrayClass, tempVars.size());
      for (size_t i = 0; i < tempVars.size(); ++i) {
        Symbol *sym = Symbol::intern(tempVars[i]);
        setArraySlot(tempNamesArray, i, sym);
      }
    }
  }
}

} // namespace smalltalk
