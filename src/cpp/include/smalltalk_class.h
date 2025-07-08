#pragma once

#include "compiled_method.h"
#include "symbol.h"
#include "object.h"

#include <unordered_map>
#include <string>
#include <vector>
#include <memory>

namespace smalltalk {

// Forward declarations
// Object is now defined in object.h
class TaggedValue;

/**
 * MethodDictionary maps selector symbols to compiled methods.
 * This is the core data structure for method lookup.
 */
class MethodDictionary {
public:
    MethodDictionary() = default;
    
    // Add a method to the dictionary
    void addMethod(Symbol* selector, std::shared_ptr<CompiledMethod> method);
    
    // Look up a method by selector
    std::shared_ptr<CompiledMethod> lookupMethod(Symbol* selector) const;
    
    // Remove a method
    void removeMethod(Symbol* selector);
    
    // Check if a method exists
    bool hasMethod(Symbol* selector) const;
    
    // Get all selectors
    std::vector<Symbol*> getSelectors() const;
    
    // Get number of methods
    size_t size() const { return methods_.size(); }
    
    // Clear all methods
    void clear() { methods_.clear(); }
    
private:
    std::unordered_map<Symbol*, std::shared_ptr<CompiledMethod>> methods_;
};

/**
 * Object format flags for different types of objects
 */
enum class ObjectFormat : uint8_t {
    POINTER_OBJECTS = 0,    // Regular objects with named instance variables
    INDEXABLE_OBJECTS = 1,  // Objects with indexed slots (like Array)
    BYTE_INDEXABLE = 2,     // Objects with byte-indexed data (like ByteArray, String)
    COMPILED_METHOD = 3     // Special format for compiled methods
};

/**
 * Class represents a Smalltalk class.
 * Contains method dictionary, superclass, instance variables, etc.
 */
class Class : public Object {
public:
    Class(const std::string& name, Class* superclass = nullptr, Class* metaclass = nullptr);
    
    // Class identity
    const std::string& getName() const { return name_; }
    
    // Superclass chain
    Class* getSuperclass() const { return superclass_; }
    void setSuperclass(Class* superclass) { superclass_ = superclass; }
    
    // Metaclass (class of this class)
    Class* getMetaclass() const { return metaclass_; }
    void setMetaclass(Class* metaclass) { metaclass_ = metaclass; }
    
    // Method dictionary operations
    MethodDictionary& getMethodDictionary() { return methodDictionary_; }
    const MethodDictionary& getMethodDictionary() const { return methodDictionary_; }
    
    // Method lookup with inheritance
    std::shared_ptr<CompiledMethod> lookupMethod(Symbol* selector) const;
    
    // Add a method to this class
    void addMethod(Symbol* selector, std::shared_ptr<CompiledMethod> method);
    
    // Remove a method from this class
    void removeMethod(Symbol* selector);
    
    // Check if this class (not superclasses) has a method
    bool hasMethod(Symbol* selector) const;
    
    // Instance variables
    void addInstanceVariable(const std::string& name);
    const std::vector<std::string>& getInstanceVariables() const { return instanceVariables_; }
    size_t getInstanceVariableCount() const { return instanceVariables_.size(); }
    int getInstanceVariableIndex(const std::string& name) const;
    
    // Class variables
    void addClassVariable(const std::string& name);
    const std::vector<std::string>& getClassVariables() const { return classVariables_; }
    
    // Instance size and format
    size_t getInstanceSize() const { return instanceSize_; }
    void setInstanceSize(size_t size) { instanceSize_ = size; }
    
    ObjectFormat getFormat() const { return format_; }
    void setFormat(ObjectFormat format) { format_ = format; }
    
    bool isIndexable() const { 
        return format_ == ObjectFormat::INDEXABLE_OBJECTS || 
               format_ == ObjectFormat::BYTE_INDEXABLE; 
    }
    
    bool isByteIndexable() const { 
        return format_ == ObjectFormat::BYTE_INDEXABLE; 
    }
    
    bool isPointerFormat() const {
        return format_ == ObjectFormat::POINTER_OBJECTS;
    }
    
    // Object creation
    virtual Object* createInstance() const;
    virtual Object* createInstance(size_t indexedSize) const;
    
    // Inheritance testing
    bool isSubclassOf(const Class* other) const;
    bool isSuperclassOf(const Class* other) const;
    
    // String representation
    std::string toString() const override;
    
    // Class hierarchy utilities
    std::vector<Class*> getSuperclasses() const;
    std::vector<Class*> getAllSubclasses() const;
    
private:
    std::string name_;
    Class* superclass_;
    Class* metaclass_;
    
    MethodDictionary methodDictionary_;
    
    std::vector<std::string> instanceVariables_;
    std::vector<std::string> classVariables_;
    
    // Instance format information
    size_t instanceSize_;           // Number of named instance variables
    ObjectFormat format_;           // How instances are formatted
    
    // Keep track of direct subclasses
    mutable std::vector<Class*> subclasses_;
    
    // Register/unregister subclass relationships
    void addSubclass(Class* subclass) const;
    void removeSubclass(Class* subclass) const;
};

/**
 * Metaclass represents the class of a class.
 * Class methods are stored in the metaclass.
 */
class Metaclass : public Class {
public:
    Metaclass(const std::string& name, Class* instanceClass, Class* superclass = nullptr);
    
    // The class this is a metaclass of
    Class* getInstanceClass() const { return instanceClass_; }
    
    // Metaclasses create classes, not instances
    Object* createInstance() const override;
    
    std::string toString() const override;
    
private:
    Class* instanceClass_;
};

/**
 * ClassRegistry maintains all classes in the system.
 * Provides global access to classes by name.
 */
class ClassRegistry {
public:
    // Singleton access
    static ClassRegistry& getInstance();
    
    // Register a class
    void registerClass(const std::string& name, Class* clazz);
    
    // Look up a class by name
    Class* getClass(const std::string& name) const;
    
    // Check if a class exists
    bool hasClass(const std::string& name) const;
    
    // Get all registered classes
    std::vector<Class*> getAllClasses() const;
    
    // Get all class names
    std::vector<std::string> getAllClassNames() const;
    
    // Remove a class
    void removeClass(const std::string& name);
    
    // Clear all classes
    void clear();
    
private:
    ClassRegistry() = default;
    std::unordered_map<std::string, Class*> classes_;
};

// Convenience functions for common class operations
namespace ClassUtils {
    // Initialize the core class hierarchy (Object, Class, Metaclass, etc.)
    void initializeCoreClasses();
    
    // Get common classes
    Class* getObjectClass();
    Class* getClassClass();
    Class* getMetaclassClass();
    Class* getIntegerClass();
    Class* getBooleanClass();
    Class* getSymbolClass();
    Class* getStringClass();
    Class* getBlockClass();
    
    // Create a new class
    Class* createClass(const std::string& name, Class* superclass = nullptr);
    
    // Create a new metaclass
    Metaclass* createMetaclass(const std::string& name, Class* instanceClass, Class* superclass = nullptr);
    
    // Add a primitive method to a class
    void addPrimitiveMethod(Class* clazz, const std::string& selector, int primitiveNumber);
}

} // namespace smalltalk