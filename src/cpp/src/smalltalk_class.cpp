#include "smalltalk_class.h"
#include "bootstrap_api.h"
#include "memory_manager.h"
#include "method_compiler.h"
#include "primitives.h"
#include <algorithm>
#include <iostream>
#include <sstream>

namespace smalltalk
{

  // MethodDictionary implementation
  void MethodDictionary::addMethod(Symbol *selector,
                                   std::shared_ptr<CompiledMethod> method)
  {
    methods_[selector] = method;
  }

  std::shared_ptr<CompiledMethod>
  MethodDictionary::lookupMethod(Symbol *selector) const
  {
    auto it = methods_.find(selector);
    if (it != methods_.end())
    {
      return it->second;
    }
    return nullptr;
  }

  void MethodDictionary::removeMethod(Symbol *selector)
  {
    methods_.erase(selector);
  }

  bool MethodDictionary::hasMethod(Symbol *selector) const
  {
    return methods_.find(selector) != methods_.end();
  }

  std::vector<Symbol *> MethodDictionary::getSelectors() const
  {
    std::vector<Symbol *> selectors;
    for (const auto &pair : methods_)
    {
      selectors.push_back(pair.first);
    }
    return selectors;
  }

  // Class implementation
  Class::Class(const std::string &name, Class *superclass, Class *metaclass)
      : Object(ObjectType::CLASS, sizeof(Class), nullptr), name_(name),
        superclass_(superclass), metaclass_(metaclass), instanceSize_(0),
        format_(ObjectFormat::POINTER_OBJECTS)
  {

    // Calculate instance size from superclass chain
    if (superclass_ != nullptr)
    {
      instanceSize_ = superclass_->getInstanceSize();
      superclass_->addSubclass(this);
    }
  }

  std::shared_ptr<CompiledMethod> Class::lookupMethod(Symbol *selector) const
  {
    // Consult only the Smalltalk MethodDictionary mirror for this class
    if (methodDictObject_)
    {
      TaggedValue *slots = reinterpret_cast<TaggedValue *>(
          reinterpret_cast<char *>(methodDictObject_) + sizeof(Object));
      Object *keysArr = slots[0].isPointer() ? slots[0].asObject() : nullptr;
      Object *valsArr = slots[1].isPointer() ? slots[1].asObject() : nullptr;
      if (keysArr && valsArr)
      {
        size_t n = keysArr->header.size;
        TaggedValue *kSlots = reinterpret_cast<TaggedValue *>(
            reinterpret_cast<char *>(keysArr) + sizeof(Object));
        TaggedValue *vSlots = reinterpret_cast<TaggedValue *>(
            reinterpret_cast<char *>(valsArr) + sizeof(Object));
        for (size_t i = 0; i < n; ++i)
        {
          if (kSlots[i].isPointer() && kSlots[i].asObject() == selector)
          {
            if (vSlots[i].isPointer())
            {
              CompiledMethod *method = static_cast<CompiledMethod *>(vSlots[i].asObject());
              // std::cerr << "Retrieved method " << selector->getName() << " at address " << method << std::endl;
              // std::cerr << "Method primitive number: " << method->primitiveNumber << std::endl;
              // std::cerr << "Method bytecode size: " << method->getBytecodes().size() << std::endl;
              return std::shared_ptr<CompiledMethod>(method, [](CompiledMethod *) {});
            }
            break;
          }
        }
      }
    }

    // Not found here; check superclass chain
    if (superclass_ != nullptr)
    {
      return superclass_->lookupMethod(selector);
    }

    return nullptr;
  }

  void Class::addMethod(Symbol *selector,
                        std::shared_ptr<CompiledMethod> method)
  {
    // Update Smalltalk mirror if present; otherwise invalidate to rebuild lazily
    if (methodDictObject_)
    {
      TaggedValue *slots = reinterpret_cast<TaggedValue *>(
          reinterpret_cast<char *>(methodDictObject_) + sizeof(Object));
      Object *keysArr = slots[0].isPointer() ? slots[0].asObject() : nullptr;
      Object *valsArr = slots[1].isPointer() ? slots[1].asObject() : nullptr;
      if (keysArr && valsArr)
      {
        size_t n = keysArr->header.size;
        TaggedValue *kSlots = reinterpret_cast<TaggedValue *>(
            reinterpret_cast<char *>(keysArr) + sizeof(Object));
        TaggedValue *vSlots = reinterpret_cast<TaggedValue *>(
            reinterpret_cast<char *>(valsArr) + sizeof(Object));
        for (size_t i = 0; i < n; ++i)
        {
          if (kSlots[i].isPointer() && kSlots[i].asObject() == selector)
          {
            vSlots[i] = TaggedValue::fromObject(method.get());
            // Also update C++ map to keep dual-consistency until removed
            methodDictionary_.addMethod(selector, method);
            return;
          }
        }
        // If key not found, invalidate to rebuild lazily
        methodDictObject_ = nullptr;
      }
      else
      {
        methodDictObject_ = nullptr;
      }
    }
    methodDictionary_.addMethod(selector, method);
  }

  void Class::removeMethod(Symbol *selector)
  {
    // Update Smalltalk dictionary mirror by compacting arrays if present
    if (methodDictObject_)
    {
      TaggedValue *slots = reinterpret_cast<TaggedValue *>(
          reinterpret_cast<char *>(methodDictObject_) + sizeof(Object));
      Object *keysArr = slots[0].isPointer() ? slots[0].asObject() : nullptr;
      Object *valsArr = slots[1].isPointer() ? slots[1].asObject() : nullptr;
      auto arraySlots = [](Object *arr) -> Object **
      {
        return reinterpret_cast<Object **>(reinterpret_cast<char *>(arr) +
                                           sizeof(Object));
      };
      if (keysArr && valsArr)
      {
        size_t n = keysArr->header.size;
        Object **kSlots = arraySlots(keysArr);
        Object **vSlots = arraySlots(valsArr);
        ssize_t idx = -1;
        for (size_t i = 0; i < n; ++i)
        {
          if (kSlots[i] == selector)
          {
            idx = static_cast<ssize_t>(i);
            break;
          }
        }
        if (idx >= 0)
        {
          if (n == 1)
          {
            // Replace with empty arrays
            auto allocArray0 = [&]() -> Object *
            {
              size_t total = sizeof(Object);
              Object *obj = reinterpret_cast<Object *>(malloc(total));
              if (!obj)
                throw std::runtime_error("Failed to allocate array");
              new (obj) Object(ObjectType::ARRAY, total,
                               ClassRegistry::getInstance().getClass("Array"));
              return obj;
            };
            slots[0] = TaggedValue(allocArray0());
            slots[1] = TaggedValue(allocArray0());
          }
          else
          {
            // Allocate n-1 arrays and copy excluding idx
            auto allocArrayN = [&](size_t sz) -> Object *
            {
              size_t total = sizeof(Object) + sz * sizeof(Object *);
              Object *obj = reinterpret_cast<Object *>(malloc(total));
              if (!obj)
                throw std::runtime_error("Failed to allocate array");
              new (obj) Object(ObjectType::ARRAY, total,
                               ClassRegistry::getInstance().getClass("Array"));
              Object **es = arraySlots(obj);
              for (size_t i = 0; i < sz; ++i)
                es[i] = nullptr;
              return obj;
            };
            Object *newKeys = allocArrayN(n - 1);
            Object *newVals = allocArrayN(n - 1);
            Object **nk = arraySlots(newKeys);
            Object **nv = arraySlots(newVals);
            size_t j = 0;
            for (size_t i = 0; i < n; ++i)
            {
              if (static_cast<ssize_t>(i) == idx)
                continue;
              nk[j] = kSlots[i];
              nv[j] = vSlots[i];
              ++j;
            }
            slots[0] = TaggedValue(newKeys);
            slots[1] = TaggedValue(newVals);
          }
        }
      }
    }
  }

  bool Class::hasMethod(Symbol *selector) const
  {
    if (methodDictObject_)
    {
      TaggedValue *slots = reinterpret_cast<TaggedValue *>(
          reinterpret_cast<char *>(methodDictObject_) + sizeof(Object));
      Object *keysArr = slots[0].isPointer() ? slots[0].asObject() : nullptr;
      if (keysArr)
      {
        size_t n = keysArr->header.size;
        TaggedValue *kSlots = reinterpret_cast<TaggedValue *>(
            reinterpret_cast<char *>(keysArr) + sizeof(Object));
        for (size_t i = 0; i < n; ++i)
        {
          if (kSlots[i].isPointer() && kSlots[i].asObject() == selector)
            return true;
        }
        return false;
      }
    }
    return false;
  }

  void Class::addInstanceVariable(const std::string &name)
  {
    instanceVariables_.push_back(name);
    instanceSize_++;
  }

  int Class::getInstanceVariableIndex(const std::string &name) const
  {
    auto it =
        std::find(instanceVariables_.begin(), instanceVariables_.end(), name);
    if (it != instanceVariables_.end())
    {
      return static_cast<int>(it - instanceVariables_.begin());
    }
    return -1;
  }

  void Class::addClassVariable(const std::string &name)
  {
    classVariables_.push_back(name);
  }

  Object *Class::createInstance() const { return createInstance(0); }

  Object *Class::createInstance(size_t indexedSize) const
  {
    // Calculate total size needed
    size_t totalSize = sizeof(Object);

    if (format_ == ObjectFormat::POINTER_OBJECTS)
    {
      // Named instance variables + indexed slots (if any)
      totalSize += (instanceSize_ + indexedSize) * sizeof(Object *);
    }
    else if (format_ == ObjectFormat::INDEXABLE_OBJECTS)
    {
      // Named instance variables + indexed object slots
      totalSize +=
          instanceSize_ * sizeof(Object *) + indexedSize * sizeof(Object *);
    }
    else if (format_ == ObjectFormat::BYTE_INDEXABLE)
    {
      // Named instance variables + byte data
      totalSize += instanceSize_ * sizeof(Object *) + indexedSize;
    }

    // Allocate the object
    Object *obj = reinterpret_cast<Object *>(malloc(totalSize));
    if (obj == nullptr)
    {
      throw std::runtime_error("Failed to allocate memory for object");
    }

    // Initialize the object
    new (obj) Object(ObjectType::OBJECT, totalSize, const_cast<Class *>(this));

    // Clear instance variable slots
    if (instanceSize_ > 0 || indexedSize > 0)
    {
      Object **slots = reinterpret_cast<Object **>(reinterpret_cast<char *>(obj) +
                                                   sizeof(Object));
      for (size_t i = 0; i < instanceSize_ + indexedSize; ++i)
      {
        slots[i] = nullptr;
      }
    }

    return obj;
  }

  bool Class::isSubclassOf(const Class *other) const
  {
    if (other == nullptr)
      return false;

    Class *current = superclass_;
    while (current != nullptr)
    {
      if (current == other)
      {
        return true;
      }
      current = current->getSuperclass();
    }
    return false;
  }

  bool Class::isSuperclassOf(const Class *other) const
  {
    return other != nullptr && other->isSubclassOf(this);
  }

  std::string Class::toString() const { return name_; }

  void Class::ensureSmalltalkMetadata(MemoryManager &mm)
  {
    // Helper to allocate and fill an Array of Symbols from a vector<string>
    auto makeSymbolArray = [&](const std::vector<std::string> &names,
                               Object *&cache)
    {
      if (cache)
        return;
      Class *arrayClass = ClassRegistry::getInstance().getClass("Array");
      if (!arrayClass)
        return;
      cache = mm.allocateIndexableInstance(arrayClass, names.size());
      TaggedValue *slots = reinterpret_cast<TaggedValue *>(
          reinterpret_cast<char *>(cache) + sizeof(Object));
      for (size_t i = 0; i < names.size(); ++i)
      {
        Symbol *sym = Symbol::intern(names[i]);
        slots[i] = TaggedValue::fromObject(sym);
      }
    };

    makeSymbolArray(instanceVariables_, instanceVarNamesArray_);
    makeSymbolArray(classVariables_, classVarNamesArray_);
  }

  void Class::ensureSmalltalkMethodDictionary(MemoryManager &mm)
  {
    if (methodDictObject_)
      return;
    Class *dictClass = ClassRegistry::getInstance().getClass("Dictionary");
    if (!dictClass)
      return;
    // Allocate the Dictionary instance via MemoryManager (no direct malloc)
    methodDictObject_ = mm.allocateInstance(dictClass);

    // Create empty arrays that can grow during bootstrap
    // The installCompiledMethod function will populate them
    TaggedValue *slots = reinterpret_cast<TaggedValue *>(
        reinterpret_cast<char *>(methodDictObject_) + sizeof(Object));
    Class *arrayClass = ClassRegistry::getInstance().getClass("Array");
    if (!arrayClass)
      return;

    // Create empty arrays with capacity for growth
    Object *keysArr = mm.allocateIndexableInstance(arrayClass, 50);
    Object *valsArr = mm.allocateIndexableInstance(arrayClass, 50);
    keysArr->header.size = 0; // Start with 0 elements
    valsArr->header.size = 0; // Start with 0 elements
    slots[0] = TaggedValue(keysArr);
    slots[1] = TaggedValue(valsArr);

    // Don't populate here - let installCompiledMethod handle it
    // This avoids double allocation during bootstrap
  }

  std::vector<Class *> Class::getSuperclasses() const
  {
    std::vector<Class *> superclasses;
    Class *current = superclass_;
    while (current != nullptr)
    {
      superclasses.push_back(current);
      current = current->getSuperclass();
    }
    return superclasses;
  }

  std::vector<Class *> Class::getAllSubclasses() const
  {
    std::vector<Class *> allSubclasses;

    // Add direct subclasses
    for (Class *subclass : subclasses_)
    {
      allSubclasses.push_back(subclass);

      // Recursively add subclasses of subclasses
      auto subSubclasses = subclass->getAllSubclasses();
      allSubclasses.insert(allSubclasses.end(), subSubclasses.begin(),
                           subSubclasses.end());
    }

    return allSubclasses;
  }

  void Class::addSubclass(Class *subclass) const
  {
    if (subclass != nullptr)
    {
      subclasses_.push_back(subclass);
    }
  }

  void Class::removeSubclass(Class *subclass) const
  {
    auto it = std::find(subclasses_.begin(), subclasses_.end(), subclass);
    if (it != subclasses_.end())
    {
      subclasses_.erase(it);
    }
  }

  // Metaclass implementation
  Metaclass::Metaclass(const std::string &name, Class *instanceClass,
                       Class *superclass)
      : Class(name + " class", superclass, nullptr),
        instanceClass_(instanceClass) {}

  Object *Metaclass::createInstance() const
  {
    // Metaclasses create new instances of their instance class
    if (instanceClass_ != nullptr)
    {
      return instanceClass_->createInstance();
    }
    return nullptr;
  }

  std::string Metaclass::toString() const { return getName(); }

  // ClassRegistry implementation
  ClassRegistry &ClassRegistry::getInstance()
  {
    static ClassRegistry instance;
    return instance;
  }

  void ClassRegistry::registerClass(const std::string &name, Class *clazz)
  {
    classes_[name] = clazz;
  }

  Class *ClassRegistry::getClass(const std::string &name) const
  {
    auto it = classes_.find(name);
    if (it != classes_.end())
    {
      return it->second;
    }
    return nullptr;
  }

  bool ClassRegistry::hasClass(const std::string &name) const
  {
    return classes_.find(name) != classes_.end();
  }

  std::vector<Class *> ClassRegistry::getAllClasses() const
  {
    std::vector<Class *> classes;
    for (const auto &pair : classes_)
    {
      classes.push_back(pair.second);
    }
    return classes;
  }

  std::vector<std::string> ClassRegistry::getAllClassNames() const
  {
    std::vector<std::string> names;
    for (const auto &pair : classes_)
    {
      names.push_back(pair.first);
    }
    return names;
  }

  void ClassRegistry::removeClass(const std::string &name)
  {
    classes_.erase(name);
  }

  void ClassRegistry::clear() { classes_.clear(); }

  // ClassUtils implementation
  namespace ClassUtils
  {

    // Core class storage as a singleton
    struct CoreClasses
    {
      Class *objectClass = nullptr;
      Class *classClass = nullptr;
      Metaclass *metaclassClass = nullptr;
      Class *integerClass = nullptr;
      Class *booleanClass = nullptr;
      Class *trueClass = nullptr;
      Class *falseClass = nullptr;
      Class *undefinedObjectClass = nullptr;
      Class *symbolClass = nullptr;
      Class *stringClass = nullptr;
      Class *blockClass = nullptr;

      static CoreClasses &getInstance()
      {
        static CoreClasses instance;
        return instance;
      }
    };

    void initializeCoreClasses()
    {
      auto &core = CoreClasses::getInstance();

      // Check if already initialized by checking if core classes exist
      if (core.objectClass != nullptr)
        return;

      auto &registry = ClassRegistry::getInstance();

      // Create Object class first (root of hierarchy)
      core.objectClass = new Class("Object", nullptr, nullptr);
      core.objectClass->setInstanceSize(0); // Object has no instance variables
      core.objectClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      registry.registerClass("Object", core.objectClass);

      // Create Class class
      core.classClass = new Class("Class", core.objectClass, nullptr);
      core.classClass->setInstanceSize(0); // Class metadata is handled internally
      core.classClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      registry.registerClass("Class", core.classClass);

      // Set Object's class to be Class
      core.objectClass->setClass(core.classClass);

      // Create Metaclass class
      core.metaclassClass =
          new Metaclass("Metaclass", core.classClass, core.classClass);
      core.metaclassClass->setInstanceSize(0);
      core.metaclassClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      registry.registerClass("Metaclass", core.metaclassClass);

      // Set Class's class to be Metaclass
      core.classClass->setClass(core.metaclassClass);

      // Create Integer class
      core.integerClass = new Class("Integer", core.objectClass, nullptr);
      core.integerClass->setInstanceSize(0); // Integers are immediate values
      core.integerClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      core.integerClass->setClass(core.classClass);
      registry.registerClass("Integer", core.integerClass);

      // Create Boolean hierarchy
      core.booleanClass = new Class("Boolean", core.objectClass, nullptr);
      core.booleanClass->setInstanceSize(0); // Booleans are immediate values
      core.booleanClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      core.booleanClass->setClass(core.classClass);
      registry.registerClass("Boolean", core.booleanClass);

      // Create True/False as subclasses of Boolean
      core.trueClass = new Class("True", core.booleanClass, nullptr);
      core.trueClass->setInstanceSize(0);
      core.trueClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      core.trueClass->setClass(core.classClass);
      registry.registerClass("True", core.trueClass);

      core.falseClass = new Class("False", core.booleanClass, nullptr);
      core.falseClass->setInstanceSize(0);
      core.falseClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      core.falseClass->setClass(core.classClass);
      registry.registerClass("False", core.falseClass);

      // Create UndefinedObject class (class of nil)
      core.undefinedObjectClass =
          new Class("UndefinedObject", core.objectClass, nullptr);
      core.undefinedObjectClass->setInstanceSize(0);
      core.undefinedObjectClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      core.undefinedObjectClass->setClass(core.classClass);
      registry.registerClass("UndefinedObject", core.undefinedObjectClass);
      // Register minimal methods for UndefinedObject (installed during prepare)
      registerSmalltalkMethod(core.undefinedObjectClass,
                              "printString\n^ 'nil'");
      registerSmalltalkMethod(core.undefinedObjectClass,
                              "asString\n^ 'nil'");
      registerSmalltalkMethod(core.undefinedObjectClass,
                              "isNil\n^ true");
      registerSmalltalkMethod(core.undefinedObjectClass,
                              "ifNil: block\n^ block value");
      registerSmalltalkMethod(core.undefinedObjectClass,
                              "ifNotNil: block\n^ nil");

      // Create Symbol class
      core.symbolClass = new Class("Symbol", core.objectClass, nullptr);
      core.symbolClass->setInstanceSize(0); // Symbol data is managed internally
      core.symbolClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      core.symbolClass->setClass(core.classClass);
      registry.registerClass("Symbol", core.symbolClass);

      // Create String class - byte indexable for character data
      core.stringClass = new Class("String", core.objectClass, nullptr);
      core.stringClass->setInstanceSize(0); // No named instance variables
      core.stringClass->setFormat(ObjectFormat::BYTE_INDEXABLE);
      core.stringClass->setClass(core.classClass);
      registry.registerClass("String", core.stringClass);

      // Create Array class - indexable for elements
      Class *arrayClass = new Class("Array", core.objectClass, nullptr);
      arrayClass->setInstanceSize(0); // No named instance variables
      arrayClass->setFormat(ObjectFormat::INDEXABLE_OBJECTS);
      arrayClass->setClass(core.classClass);
      registry.registerClass("Array", arrayClass);

      // Create ByteArray class - byte indexable
      Class *byteArrayClass = new Class("ByteArray", core.objectClass, nullptr);
      byteArrayClass->setInstanceSize(0); // No named instance variables
      byteArrayClass->setFormat(ObjectFormat::BYTE_INDEXABLE);
      byteArrayClass->setClass(core.classClass);
      registry.registerClass("ByteArray", byteArrayClass);

      // Create Block class
      core.blockClass = new Class("Block", core.objectClass, nullptr);
      core.blockClass->setInstanceSize(0); // No named instance variables
      core.blockClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      core.blockClass->setClass(core.classClass);
      registry.registerClass("Block", core.blockClass);

      // Create Dictionary class (two named instance variables: keys, values)
      Class *dictionaryClass = new Class("Dictionary", core.objectClass, nullptr);
      dictionaryClass->setInstanceSize(2);
      dictionaryClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      dictionaryClass->setClass(core.classClass);
      registry.registerClass("Dictionary", dictionaryClass);

      // Add primitive methods to Object class
      addPrimitiveMethod(core.objectClass, "new", 70);      // Object new
      addPrimitiveMethod(core.objectClass, "basicNew", 71); // Object basicNew
      addPrimitiveMethod(core.objectClass,
                         "basicNew:", 72); // Object basicNew: size
      addPrimitiveMethod(core.objectClass, "identityHash",
                         75);                             // Object identityHash
      addPrimitiveMethod(core.objectClass, "class", 111); // Object class

      // Default nil-testing and conditional behavior on Object (non-nil)
      // isNil -> false; ifNil: -> nil; ifNotNil: -> evaluate block
      registerSmalltalkMethod(core.objectClass,
                              "isNil\n^ false");
      registerSmalltalkMethod(core.objectClass,
                              "ifNil: block\n^ nil");
      registerSmalltalkMethod(core.objectClass,
                              "ifNotNil: block\n^ block value");

      // Add class methods to Class class (for all classes)
      registerPrimitiveMethod(core.classClass, "new",
                              PrimitiveNumbers::NEW); // Class new
      registerPrimitiveMethod(core.classClass, "new:",
                              72); // Class new: size

      // Add primitive methods to Array class (instance methods)
      registerPrimitiveMethod(arrayClass, "at:", 60);
      registerPrimitiveMethod(arrayClass, "at:put:", 61);
      registerPrimitiveMethod(arrayClass, "size", 62);

      // Add primitive methods to Integer class (arithmetic and comparison)
      registerPrimitiveMethod(core.integerClass, "+",
                              PrimitiveNumbers::SMALL_INT_ADD);
      registerPrimitiveMethod(core.integerClass, "-",
                              PrimitiveNumbers::SMALL_INT_SUB);
      registerPrimitiveMethod(core.integerClass, "*",
                              PrimitiveNumbers::SMALL_INT_MUL);
      registerPrimitiveMethod(core.integerClass, "/",
                              PrimitiveNumbers::SMALL_INT_DIV);
      registerPrimitiveMethod(core.integerClass, "<",
                              PrimitiveNumbers::SMALL_INT_LT);
      registerPrimitiveMethod(core.integerClass, ">",
                              PrimitiveNumbers::SMALL_INT_GT);
      registerPrimitiveMethod(core.integerClass, "=",
                              PrimitiveNumbers::SMALL_INT_EQ);
      registerPrimitiveMethod(core.integerClass, "~=",
                              PrimitiveNumbers::SMALL_INT_NE);
      registerPrimitiveMethod(core.integerClass, "<=",
                              PrimitiveNumbers::SMALL_INT_LE);
      registerPrimitiveMethod(core.integerClass, ">=",
                              PrimitiveNumbers::SMALL_INT_GE);

      // Add primitive methods to String class
      registerPrimitiveMethod(core.stringClass, "at:",
                              PrimitiveNumbers::STRING_AT);
      registerPrimitiveMethod(core.stringClass, ",",
                              PrimitiveNumbers::STRING_CONCAT);
      registerPrimitiveMethod(core.stringClass, "size",
                              PrimitiveNumbers::STRING_SIZE);
      registerPrimitiveMethod(core.stringClass, "asSymbol",
                              PrimitiveNumbers::STRING_AS_SYMBOL);

      // Add primitive methods to Block class
      registerPrimitiveMethod(core.blockClass, "value",
                              PrimitiveNumbers::BLOCK_VALUE);
      registerPrimitiveMethod(core.blockClass, "value:",
                              PrimitiveNumbers::BLOCK_VALUE_ARG);

      // Add primitive methods to Dictionary class
      registerPrimitiveMethod(dictionaryClass, "at:",
                              PrimitiveNumbers::DICT_AT);
      registerPrimitiveMethod(dictionaryClass, "at:put:",
                              PrimitiveNumbers::DICT_AT_PUT);
      registerPrimitiveMethod(dictionaryClass, "keys",
                              PrimitiveNumbers::DICT_KEYS);
      registerPrimitiveMethod(dictionaryClass, "size",
                              PrimitiveNumbers::DICT_SIZE);

      // Install minimal boolean control-flow methods in True and False (prepare)
      // True
      registerSmalltalkMethod(core.trueClass,
                              "ifTrue: block\n^ block value");
      registerSmalltalkMethod(core.trueClass,
                              "ifFalse: block\n^ nil");
      registerSmalltalkMethod(
          core.trueClass, "ifTrue: t ifFalse: f\n^ t value");
      registerSmalltalkMethod(
          core.trueClass, "ifFalse: f ifTrue: t\n^ t value");
      // False
      registerSmalltalkMethod(core.falseClass,
                              "ifTrue: block\n^ nil");
      registerSmalltalkMethod(
          core.falseClass, "ifFalse: block\n^ block value");
      registerSmalltalkMethod(
          core.falseClass, "ifTrue: t ifFalse: f\n^ f value");
      registerSmalltalkMethod(
          core.falseClass, "ifFalse: f ifTrue: t\n^ f value");

      // Create SystemLoader class and add minimal start: primitive
      Class *systemLoaderClass =
          new Class("SystemLoader", core.objectClass, nullptr);
      systemLoaderClass->setInstanceSize(0);
      systemLoaderClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      systemLoaderClass->setClass(core.classClass);
      registry.registerClass("SystemLoader", systemLoaderClass);
      registerPrimitiveMethod(
          systemLoaderClass, "start:", PrimitiveNumbers::SYSTEM_LOADER_START);

      // Create Compiler class and add compile:in: bridge primitive
      Class *compilerClass = new Class("Compiler", core.objectClass, nullptr);
      compilerClass->setInstanceSize(0);
      compilerClass->setFormat(ObjectFormat::POINTER_OBJECTS);
      compilerClass->setClass(core.classClass);
      registry.registerClass("Compiler", compilerClass);
      registerPrimitiveMethod(
          compilerClass, "compile:in:", PrimitiveNumbers::COMPILER_COMPILE_IN);
    }

    Class *getObjectClass() { return CoreClasses::getInstance().objectClass; }
    Class *getClassClass() { return CoreClasses::getInstance().classClass; }
    Class *getMetaclassClass() { return CoreClasses::getInstance().metaclassClass; }
    Class *getIntegerClass() { return CoreClasses::getInstance().integerClass; }
    Class *getBooleanClass() { return CoreClasses::getInstance().booleanClass; }
    Class *getTrueClass() { return CoreClasses::getInstance().trueClass; }
    Class *getFalseClass() { return CoreClasses::getInstance().falseClass; }
    Class *getUndefinedObjectClass()
    {
      return CoreClasses::getInstance().undefinedObjectClass;
    }
    Class *getSymbolClass() { return CoreClasses::getInstance().symbolClass; }
    Class *getStringClass() { return CoreClasses::getInstance().stringClass; }
    Class *getBlockClass() { return CoreClasses::getInstance().blockClass; }

    void addPrimitiveMethod(Class *clazz, const std::string &selector,
                            int primitiveNumber);

    void addPrimitiveMethod(Class *clazz, const std::string &selector,
                            int primitiveNumber)
    {
      // Register for installation during prepareImage
      registerPrimitiveMethod(clazz, selector, primitiveNumber);
    }

    void buildAllMethodDictionaries(MemoryManager &mm)
    {
      auto &registry = ClassRegistry::getInstance();
      for (Class *cls : registry.getAllClasses())
      {
        if (cls)
        {
          cls->ensureSmalltalkMethodDictionary(mm);
        }
      }
    }

    void buildAllClassMetadata(MemoryManager &mm)
    {
      auto &registry = ClassRegistry::getInstance();
      for (Class *cls : registry.getAllClasses())
      {
        if (cls)
        {
          cls->ensureSmalltalkMetadata(mm);
        }
      }
    }

    Class *createClass(const std::string &name, Class *superclass)
    {
      auto &core = CoreClasses::getInstance();
      if (superclass == nullptr)
      {
        superclass = core.objectClass;
      }

      Class *newClass = new Class(name, superclass, nullptr);
      newClass->setClass(core.classClass);

      ClassRegistry::getInstance().registerClass(name, newClass);
      return newClass;
    }

    Metaclass *createMetaclass(const std::string &name, Class *instanceClass,
                               Class *superclass)
    {
      auto &core = CoreClasses::getInstance();
      Metaclass *newMetaclass = new Metaclass(name, instanceClass, superclass);
      newMetaclass->setClass(core.metaclassClass);

      ClassRegistry::getInstance().registerClass(name + " class", newMetaclass);
      return newMetaclass;
    }
  } // namespace ClassUtils

} // namespace smalltalk
