#include "smalltalk_class.h"
#include "primitives.h"
#include <algorithm>
#include <sstream>
#include <iostream>

namespace smalltalk
{

    // MethodDictionary implementation
    void MethodDictionary::addMethod(Symbol *selector, std::shared_ptr<CompiledMethod> method)
    {
        methods_[selector] = method;
    }

    std::shared_ptr<CompiledMethod> MethodDictionary::lookupMethod(Symbol *selector) const
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
        : Object(ObjectType::CLASS, sizeof(Class), nullptr),
          name_(name),
          superclass_(superclass),
          metaclass_(metaclass),
          instanceSize_(0),
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
        // Look in this class first
        auto method = methodDictionary_.lookupMethod(selector);
        if (method != nullptr)
        {
            return method;
        }

        // Look in superclass chain
        if (superclass_ != nullptr)
        {
            return superclass_->lookupMethod(selector);
        }

        return nullptr;
    }

    void Class::addMethod(Symbol *selector, std::shared_ptr<CompiledMethod> method)
    {
        methodDictionary_.addMethod(selector, method);
    }

    void Class::removeMethod(Symbol *selector)
    {
        methodDictionary_.removeMethod(selector);
    }

    bool Class::hasMethod(Symbol *selector) const
    {
        return methodDictionary_.hasMethod(selector);
    }

    void Class::addInstanceVariable(const std::string &name)
    {
        instanceVariables_.push_back(name);
        instanceSize_++;
    }

    int Class::getInstanceVariableIndex(const std::string &name) const
    {
        auto it = std::find(instanceVariables_.begin(), instanceVariables_.end(), name);
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

    Object *Class::createInstance() const
    {
        return createInstance(0);
    }

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
            totalSize += instanceSize_ * sizeof(Object *) + indexedSize * sizeof(Object *);
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
            Object **slots = reinterpret_cast<Object **>(reinterpret_cast<char *>(obj) + sizeof(Object));
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

    std::string Class::toString() const
    {
        return name_;
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
            allSubclasses.insert(allSubclasses.end(), subSubclasses.begin(), subSubclasses.end());
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
    Metaclass::Metaclass(const std::string &name, Class *instanceClass, Class *superclass)
        : Class(name + " class", superclass, nullptr), instanceClass_(instanceClass)
    {
    }

    Object *Metaclass::createInstance() const
    {
        // Metaclasses create new instances of their instance class
        if (instanceClass_ != nullptr)
        {
            return instanceClass_->createInstance();
        }
        return nullptr;
    }

    std::string Metaclass::toString() const
    {
        return getName();
    }

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

    void ClassRegistry::clear()
    {
        classes_.clear();
    }

    // ClassUtils implementation
    namespace ClassUtils
    {

        // Core class storage as a singleton
        struct CoreClasses {
            Class *objectClass = nullptr;
            Class *classClass = nullptr;
            Metaclass *metaclassClass = nullptr;
            Class *integerClass = nullptr;
            Class *booleanClass = nullptr;
            Class *symbolClass = nullptr;
            Class *stringClass = nullptr;
            Class *blockClass = nullptr;
            
            static CoreClasses& getInstance() {
                static CoreClasses instance;
                return instance;
            }
        };

        void initializeCoreClasses()
        {
            auto& core = CoreClasses::getInstance();
            
            // Check if already initialized by checking if core classes exist
            if (core.objectClass != nullptr) return;
            
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
            core.metaclassClass = new Metaclass("Metaclass", core.classClass, core.classClass);
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

            // Create Boolean class
            core.booleanClass = new Class("Boolean", core.objectClass, nullptr);
            core.booleanClass->setInstanceSize(0); // Booleans are immediate values
            core.booleanClass->setFormat(ObjectFormat::POINTER_OBJECTS);
            core.booleanClass->setClass(core.classClass);
            registry.registerClass("Boolean", core.booleanClass);

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

            // Add primitive methods to Object class
            addPrimitiveMethod(core.objectClass, "new", 70);          // Object new
            addPrimitiveMethod(core.objectClass, "basicNew", 71);     // Object basicNew
            addPrimitiveMethod(core.objectClass, "basicNew:", 72);    // Object basicNew: size
            addPrimitiveMethod(core.objectClass, "identityHash", 75); // Object identityHash
            addPrimitiveMethod(core.objectClass, "class", 111);       // Object class

            // Add class methods to Class class (for all classes)
            addPrimitiveMethod(core.classClass, "new", PrimitiveNumbers::NEW);   // Class new
            addPrimitiveMethod(core.classClass, "new:", 72); // Class new: size (for Array new:, etc.)

            // Add primitive methods to Array class (instance methods)
            addPrimitiveMethod(arrayClass, "at:", 60);     // Array at: index
            addPrimitiveMethod(arrayClass, "at:put:", 61); // Array at: index put: value
            addPrimitiveMethod(arrayClass, "size", 62);    // Array size

            // Add primitive methods to Integer class (arithmetic and comparison)
            addPrimitiveMethod(core.integerClass, "+", PrimitiveNumbers::SMALL_INT_ADD); // Integer +
            addPrimitiveMethod(core.integerClass, "-", PrimitiveNumbers::SMALL_INT_SUB); // Integer -
            addPrimitiveMethod(core.integerClass, "*", PrimitiveNumbers::SMALL_INT_MUL); // Integer *
            addPrimitiveMethod(core.integerClass, "/", PrimitiveNumbers::SMALL_INT_DIV); // Integer /
            addPrimitiveMethod(core.integerClass, "<", PrimitiveNumbers::SMALL_INT_LT);  // Integer <
            addPrimitiveMethod(core.integerClass, ">", PrimitiveNumbers::SMALL_INT_GT);  // Integer >
            addPrimitiveMethod(core.integerClass, "=", PrimitiveNumbers::SMALL_INT_EQ);  // Integer =
            addPrimitiveMethod(core.integerClass, "~=", PrimitiveNumbers::SMALL_INT_NE); // Integer ~=
            addPrimitiveMethod(core.integerClass, "<=", PrimitiveNumbers::SMALL_INT_LE); // Integer <=
            addPrimitiveMethod(core.integerClass, ">=", PrimitiveNumbers::SMALL_INT_GE); // Integer >=

            // Add primitive methods to String class
            addPrimitiveMethod(core.stringClass, "at:", PrimitiveNumbers::STRING_AT);    // String at:
            addPrimitiveMethod(core.stringClass, ",", PrimitiveNumbers::STRING_CONCAT);  // String ,
            addPrimitiveMethod(core.stringClass, "size", PrimitiveNumbers::STRING_SIZE); // String size
            addPrimitiveMethod(core.stringClass, "asSymbol", PrimitiveNumbers::STRING_AS_SYMBOL); // String asSymbol

            // Add primitive methods to Block class
            addPrimitiveMethod(core.blockClass, "value", PrimitiveNumbers::BLOCK_VALUE); // Block value
            addPrimitiveMethod(core.blockClass, "value:", PrimitiveNumbers::BLOCK_VALUE_ARG); // Block value:

            // Create SystemLoader class and add minimal start: primitive
            Class* systemLoaderClass = new Class("SystemLoader", core.objectClass, nullptr);
            systemLoaderClass->setInstanceSize(0);
            systemLoaderClass->setFormat(ObjectFormat::POINTER_OBJECTS);
            systemLoaderClass->setClass(core.classClass);
            registry.registerClass("SystemLoader", systemLoaderClass);
            addPrimitiveMethod(systemLoaderClass, "start:", PrimitiveNumbers::SYSTEM_LOADER_START);

            // Create Compiler class and add compile:in: bridge primitive
            Class* compilerClass = new Class("Compiler", core.objectClass, nullptr);
            compilerClass->setInstanceSize(0);
            compilerClass->setFormat(ObjectFormat::POINTER_OBJECTS);
            compilerClass->setClass(core.classClass);
            registry.registerClass("Compiler", compilerClass);
            addPrimitiveMethod(compilerClass, "compile:in:", PrimitiveNumbers::COMPILER_COMPILE_IN);
        }

        Class *getObjectClass() { return CoreClasses::getInstance().objectClass; }
        Class *getClassClass() { return CoreClasses::getInstance().classClass; }
        Class *getMetaclassClass() { return CoreClasses::getInstance().metaclassClass; }
        Class *getIntegerClass() { return CoreClasses::getInstance().integerClass; }
        Class *getBooleanClass() { return CoreClasses::getInstance().booleanClass; }
        Class *getSymbolClass() { return CoreClasses::getInstance().symbolClass; }
        Class *getStringClass() { return CoreClasses::getInstance().stringClass; }
        Class *getBlockClass() { return CoreClasses::getInstance().blockClass; }

        void addPrimitiveMethod(Class *clazz, const std::string &selector, int primitiveNumber);

        void addPrimitiveMethod(Class *clazz, const std::string &selector, int primitiveNumber)
        {
            // Create a compiled method with the primitive
            auto method = std::make_shared<CompiledMethod>();
            method->primitiveNumber = primitiveNumber;
            method->bytecodes.clear(); // No bytecode needed for primitives
            method->literals.clear();

            // Create selector symbol
            Symbol *selectorSymbol = Symbol::intern(selector);

            // Add method to class
            clazz->addMethod(selectorSymbol, method);
        }

        Class *createClass(const std::string &name, Class *superclass)
        {
            auto& core = CoreClasses::getInstance();
            if (superclass == nullptr)
            {
                superclass = core.objectClass;
            }

            Class *newClass = new Class(name, superclass, nullptr);
            newClass->setClass(core.classClass);

            ClassRegistry::getInstance().registerClass(name, newClass);
            return newClass;
        }

        Metaclass *createMetaclass(const std::string &name, Class *instanceClass, Class *superclass)
        {
            auto& core = CoreClasses::getInstance();
            Metaclass *newMetaclass = new Metaclass(name, instanceClass, superclass);
            newMetaclass->setClass(core.metaclassClass);

            ClassRegistry::getInstance().registerClass(name + " class", newMetaclass);
            return newMetaclass;
        }
    }

} // namespace smalltalk
