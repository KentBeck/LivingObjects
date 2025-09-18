#include "bootstrap_builder.h"
#include "bootstrap_api.h"
#include "method_compiler.h"
#include "smalltalk_class.h"
#include "symbol.h"

namespace smalltalk
{

  namespace
  {
    struct PendingSmalltalk
    {
      Class *clazz;
      std::string source;
    };

    struct PendingPrimitive
    {
      Class *clazz;
      std::string selector;
      int primitive;
    };

    static std::vector<PendingSmalltalk> &smalltalkQueue()
    {
      static std::vector<PendingSmalltalk> q;
      return q;
    }

    static std::vector<PendingPrimitive> &primitiveQueue()
    {
      static std::vector<PendingPrimitive> q;
      return q;
    }

    static TaggedValue *arraySlots(Object *arr)
    {
      return reinterpret_cast<TaggedValue *>(reinterpret_cast<char *>(arr) +
                                             sizeof(Object));
    }

    static void installCompiledMethod(Class *cls, Symbol *selector,
                                      const std::shared_ptr<CompiledMethod> &m,
                                      MemoryManager &mm)
    {
      // std::cerr << "Installing method " << selector->getName() << " in " << cls->getName() << " with primitive " << m->primitiveNumber << std::endl;
      // Ensure dictionary object exists
      cls->ensureSmalltalkMethodDictionary(mm);
      Object *dict = cls->getMethodDictionaryObject();
      if (!dict)
        return;
      TaggedValue *slots = reinterpret_cast<TaggedValue *>(reinterpret_cast<char *>(dict) +
                                                           sizeof(Object));
      Class *arrayClass = ClassRegistry::getInstance().getClass("Array");
      if (!arrayClass)
        return;

      Object *keysArr = slots[0].isPointer() ? slots[0].asObject() : nullptr;
      Object *valsArr = slots[1].isPointer() ? slots[1].asObject() : nullptr;
      if (!keysArr || !valsArr)
      {
        // Allocate arrays with enough space for all expected methods (32 primitives + some Smalltalk methods)
        keysArr = mm.allocateIndexableInstance(arrayClass, 50);
        valsArr = mm.allocateIndexableInstance(arrayClass, 50);
        // Initialize size to 0 (actual used size)
        keysArr->header.size = 0;
        valsArr->header.size = 0;
        slots[0] = TaggedValue(keysArr);
        slots[1] = TaggedValue(valsArr);
      }

      size_t n = keysArr->header.size;
      TaggedValue *kSlots = arraySlots(keysArr);
      TaggedValue *vSlots = arraySlots(valsArr);

      // Replace if exists
      for (size_t i = 0; i < n; ++i)
      {
        if (kSlots[i].isPointer() && kSlots[i].asObject() == selector)
        {
          vSlots[i] = TaggedValue::fromObject(m.get());
          return;
        }
      }

      // We pre-allocated with capacity 50, so we should have space
      if (n < 50)
      {
        // We have space, just append
        kSlots[n] = TaggedValue::fromObject(selector);
        vSlots[n] = TaggedValue::fromObject(m.get());
        // std::cerr << "Stored method " << selector->getName() << " with primitive " << m->primitiveNumber << " at index " << n << " at address " << m.get() << std::endl;
        keysArr->header.size = n + 1;
        valsArr->header.size = n + 1;
      }
      else
      {
        // This shouldn't happen during bootstrap, but handle it gracefully
        throw std::runtime_error("Method dictionary capacity exceeded during bootstrap");
      }
    }
  } // namespace

  void BootstrapBuilder::buildMethodDictionaries(MemoryManager &mm)
  {
    // Iterate over all classes and ensure their Smalltalk MethodDictionary
    // mirrors are allocated and populated from the C++ maps. No sends.
    auto &registry = ClassRegistry::getInstance();
    for (Class *cls : registry.getAllClasses())
    {
      if (cls)
      {
        cls->ensureSmalltalkMethodDictionary(mm);
      }
    }
  }

  void BootstrapBuilder::registerSmalltalkMethod(Class *clazz,
                                                 const std::string &methodSource)
  {
    smalltalkQueue().push_back(PendingSmalltalk{clazz, methodSource});
  }

  void BootstrapBuilder::registerPrimitiveMethod(Class *clazz,
                                                 const std::string &selector,
                                                 int primitiveNumber)
  {
    primitiveQueue().push_back(
        PendingPrimitive{clazz, selector, primitiveNumber});
  }

  void BootstrapBuilder::prepareImage(MemoryManager &mm)
  {
    // Ensure each class has a dictionary object (empty is fine)
    buildMethodDictionaries(mm);

    // Install pending primitives (no sends)
    for (const auto &p : primitiveQueue())
    {
      // Allocate CompiledMethod in managed memory space
      Class *methodClass = ClassRegistry::getInstance().getClass("CompiledMethod");
      if (!methodClass)
      {
        // Create a basic CompiledMethod class if it doesn't exist
        methodClass = new Class("CompiledMethod", ClassUtils::getObjectClass(), nullptr);
        methodClass->setInstanceSize(sizeof(CompiledMethod));
        methodClass->setFormat(ObjectFormat::POINTER_OBJECTS);
        methodClass->setClass(ClassUtils::getClassClass());
        ClassRegistry::getInstance().registerClass("CompiledMethod", methodClass);
      }

      // Allocate in managed memory
      Object *methodObj = mm.allocateInstance(methodClass);
      CompiledMethod *method = static_cast<CompiledMethod *>(methodObj);

      // Initialize the method
      new (method) CompiledMethod(); // Placement new to call constructor
      method->primitiveNumber = p.primitive;
      method->bytecodes.clear();
      method->literals.clear();

      Symbol *sel = Symbol::intern(p.selector);

      // Create a shared_ptr that doesn't delete (managed by GC)
      std::shared_ptr<CompiledMethod> methodPtr(method, [](CompiledMethod *) {});
      installCompiledMethod(p.clazz, sel, methodPtr, mm);
    }

    // Install pending Smalltalk methods (no sends)
    for (const auto &m : smalltalkQueue())
    {
      MethodCompiler::addSmalltalkMethod(m.clazz, m.source, mm);
    }

    // Clear queues after installation
    primitiveQueue().clear();
    smalltalkQueue().clear();
  }

  // --- Free function API wrappers to avoid exposing builder in class init ---
  void registerSmalltalkMethod(Class *clazz, const std::string &methodSource)
  {
    BootstrapBuilder::registerSmalltalkMethod(clazz, methodSource);
  }

  void registerPrimitiveMethod(Class *clazz, const std::string &selector,
                               int primitiveNumber)
  {
    BootstrapBuilder::registerPrimitiveMethod(clazz, selector, primitiveNumber);
  }

  void prepareImage(MemoryManager &mm) { BootstrapBuilder::prepareImage(mm); }

} // namespace smalltalk
