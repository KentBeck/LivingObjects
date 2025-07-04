#include "smalltalk_class.h"
#include "primitives.h"
#include <algorithm>
#include <sstream>

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

        // Global class pointers for core classes
        static Class *objectClass = nullptr;
        static Class *classClass = nullptr;
        static Metaclass *metaclassClass = nullptr;
        static Class *integerClass = nullptr;
        static Class *booleanClass = nullptr;
        static Class *symbolClass = nullptr;
        static Class *stringClass = nullptr;

        void initializeCoreClasses()
        {
            auto &registry = ClassRegistry::getInstance();

            // Create Object class first (root of hierarchy)
            objectClass = new Class("Object", nullptr, nullptr);
            objectClass->setInstanceSize(0); // Object has no instance variables
            objectClass->setFormat(ObjectFormat::POINTER_OBJECTS);
            registry.registerClass("Object", objectClass);

            // Create Class class
            classClass = new Class("Class", objectClass, nullptr);
            classClass->setInstanceSize(0); // Class metadata is handled internally
            classClass->setFormat(ObjectFormat::POINTER_OBJECTS);
            registry.registerClass("Class", classClass);

            // Set Object's class to be Class
            objectClass->setClass(classClass);

            // Create Metaclass class
            metaclassClass = new Metaclass("Metaclass", classClass, classClass);
            metaclassClass->setInstanceSize(0);
            metaclassClass->setFormat(ObjectFormat::POINTER_OBJECTS);
            registry.registerClass("Metaclass", metaclassClass);

            // Set Class's class to be Metaclass
            classClass->setClass(metaclassClass);

            // Create Integer class
            integerClass = new Class("Integer", objectClass, nullptr);
            integerClass->setInstanceSize(0); // Integers are immediate values
            integerClass->setFormat(ObjectFormat::POINTER_OBJECTS);
            integerClass->setClass(classClass);
            registry.registerClass("Integer", integerClass);

            // Create Boolean class
            booleanClass = new Class("Boolean", objectClass, nullptr);
            booleanClass->setInstanceSize(0); // Booleans are immediate values
            booleanClass->setFormat(ObjectFormat::POINTER_OBJECTS);
            booleanClass->setClass(classClass);
            registry.registerClass("Boolean", booleanClass);

            // Create Symbol class
            symbolClass = new Class("Symbol", objectClass, nullptr);
            symbolClass->setInstanceSize(0); // Symbol data is managed internally
            symbolClass->setFormat(ObjectFormat::POINTER_OBJECTS);
            symbolClass->setClass(classClass);
            registry.registerClass("Symbol", symbolClass);

            // Create String class - byte indexable for character data
            stringClass = new Class("String", objectClass, nullptr);
            stringClass->setInstanceSize(0); // No named instance variables
            stringClass->setFormat(ObjectFormat::BYTE_INDEXABLE);
            stringClass->setClass(classClass);
            registry.registerClass("String", stringClass);

            // Create Array class - indexable for elements
            Class *arrayClass = new Class("Array", objectClass, nullptr);
            arrayClass->setInstanceSize(0); // No named instance variables
            arrayClass->setFormat(ObjectFormat::INDEXABLE_OBJECTS);
            arrayClass->setClass(classClass);
            registry.registerClass("Array", arrayClass);

            // Create ByteArray class - byte indexable
            Class *byteArrayClass = new Class("ByteArray", objectClass, nullptr);
            byteArrayClass->setInstanceSize(0); // No named instance variables
            byteArrayClass->setFormat(ObjectFormat::BYTE_INDEXABLE);
            byteArrayClass->setClass(classClass);
            registry.registerClass("ByteArray", byteArrayClass);

            // Add primitive methods to Object class
            addPrimitiveMethod(objectClass, "new", 70);          // Object new
            addPrimitiveMethod(objectClass, "basicNew", 71);     // Object basicNew
            addPrimitiveMethod(objectClass, "basicNew:", 72);    // Object basicNew: size
            addPrimitiveMethod(objectClass, "identityHash", 75); // Object identityHash
            addPrimitiveMethod(objectClass, "class", 111);       // Object class

            // Add class methods to Class class (for all classes)
            addPrimitiveMethod(classClass, "new:", 72); // Class new: size (for Array new:, etc.)

            // Add primitive methods to Array class (instance methods)
            addPrimitiveMethod(arrayClass, "at:", 60);     // Array at: index
            addPrimitiveMethod(arrayClass, "at:put:", 61); // Array at: index put: value
            addPrimitiveMethod(arrayClass, "size", 62);    // Array size

            // Add primitive methods to Integer class (arithmetic and comparison)
            addPrimitiveMethod(integerClass, "+", PrimitiveNumbers::SMALL_INT_ADD); // Integer +
            addPrimitiveMethod(integerClass, "-", PrimitiveNumbers::SMALL_INT_SUB); // Integer -
            addPrimitiveMethod(integerClass, "*", PrimitiveNumbers::SMALL_INT_MUL); // Integer *
            addPrimitiveMethod(integerClass, "/", PrimitiveNumbers::SMALL_INT_DIV); // Integer /
            addPrimitiveMethod(integerClass, "<", PrimitiveNumbers::SMALL_INT_LT);  // Integer <
            addPrimitiveMethod(integerClass, ">", PrimitiveNumbers::SMALL_INT_GT);  // Integer >
            addPrimitiveMethod(integerClass, "=", PrimitiveNumbers::SMALL_INT_EQ);  // Integer =
            addPrimitiveMethod(integerClass, "~=", PrimitiveNumbers::SMALL_INT_NE); // Integer ~=
            addPrimitiveMethod(integerClass, "<=", PrimitiveNumbers::SMALL_INT_LE); // Integer <=
            addPrimitiveMethod(integerClass, ">=", PrimitiveNumbers::SMALL_INT_GE); // Integer >=

            // Add primitive methods to String class
            addPrimitiveMethod(stringClass, ",", PrimitiveNumbers::STRING_CONCAT);  // String ,
            addPrimitiveMethod(stringClass, "size", PrimitiveNumbers::STRING_SIZE); // String size
        }

        Class *getObjectClass() { return objectClass; }
        Class *getClassClass() { return classClass; }
        Class *getMetaclassClass() { return metaclassClass; }
        Class *getIntegerClass() { return integerClass; }
        Class *getBooleanClass() { return booleanClass; }
        Class *getSymbolClass() { return symbolClass; }
        Class *getStringClass() { return stringClass; }

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
            if (superclass == nullptr)
            {
                superclass = objectClass;
            }

            Class *newClass = new Class(name, superclass, nullptr);
            newClass->setClass(classClass);

            ClassRegistry::getInstance().registerClass(name, newClass);
            return newClass;
        }

        Metaclass *createMetaclass(const std::string &name, Class *instanceClass, Class *superclass)
        {
            Metaclass *newMetaclass = new Metaclass(name, instanceClass, superclass);
            newMetaclass->setClass(metaclassClass);

            ClassRegistry::getInstance().registerClass(name + " class", newMetaclass);
            return newMetaclass;
        }
    }

} // namespace smalltalk