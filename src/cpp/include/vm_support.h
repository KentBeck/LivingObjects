#pragma once

#include "simple_object.h"
#include "simple_tagged_value.h"
#include <vector>
#include <string>

namespace smalltalk {

// Forward declarations
class SmalltalkClass;
class MemoryManager;

/**
 * VM Support Functions for Simplified Object Model
 * 
 * This shows how the VM works with the new unified approach:
 * - TaggedValue for immediates (nil, true, false, integers, floats)
 * - Object* for ALL heap objects regardless of Smalltalk class
 * - VM behavior based on Smalltalk class, not C++ type
 * - No parallel C++ hierarchy - just uniform Object struct
 */

// === CORE VM VALUE TYPE ===

/**
 * VMValue represents any Smalltalk value the VM can work with
 * Either an immediate (TaggedValue) or reference to heap object (Object*)
 */
class VMValue {
public:
    // Constructors
    VMValue() : tagged_value_(TaggedValue::nil()), heap_object_(nullptr), is_immediate_(true) {}
    VMValue(TaggedValue tv) : tagged_value_(tv), heap_object_(nullptr), is_immediate_(true) {}
    VMValue(Object* obj) : tagged_value_(TaggedValue::nil()), heap_object_(obj), is_immediate_(false) {}
    
    // Type checking
    bool is_immediate() const { return is_immediate_; }
    bool is_heap_object() const { return !is_immediate_; }
    
    // Access
    TaggedValue as_immediate() const { 
        assert(is_immediate_); 
        return tagged_value_; 
    }
    
    Object* as_object() const { 
        assert(is_heap_object()); 
        return heap_object_; 
    }
    
    // Convenience checks
    bool is_nil() const { return is_immediate_ && tagged_value_.is_nil(); }
    bool is_integer() const { return is_immediate_ && tagged_value_.is_integer(); }
    bool is_boolean() const { return is_immediate_ && tagged_value_.is_boolean(); }
    
private:
    TaggedValue tagged_value_;
    Object* heap_object_;
    bool is_immediate_;
};

// === WELL-KNOWN SMALLTALK CLASSES ===

/**
 * The VM needs to know about certain Smalltalk classes for special handling.
 * These are discovered at startup by looking up class names, NOT by C++ type.
 */
struct WellKnownClasses {
    SmalltalkClass* Object;
    SmalltalkClass* Class; 
    SmalltalkClass* Array;
    SmalltalkClass* String;
    SmalltalkClass* Symbol;
    SmalltalkClass* ByteArray;
    SmalltalkClass* CompiledMethod;
    SmalltalkClass* BlockClosure;
    SmalltalkClass* Dictionary;
    SmalltalkClass* SmallInteger;  // For boxing when needed
    SmalltalkClass* Float;         // For boxing when needed
    SmalltalkClass* True;
    SmalltalkClass* False;
    
    // Initialize by looking up class names in the system
    void initialize_from_system();
};

// === VM OPERATIONS ON VALUES ===

/**
 * Get the Smalltalk class of any value (immediate or heap object)
 */
SmalltalkClass* get_smalltalk_class(const VMValue& value, const WellKnownClasses& classes);

/**
 * Check if a value is an instance of a specific Smalltalk class
 * Works for both immediates and heap objects
 */
bool is_instance_of(const VMValue& value, SmalltalkClass* target_class, const WellKnownClasses& classes);

/**
 * Convert immediate to heap object when needed (boxing)
 * Returns the existing object if already a heap object
 */
Object* box_if_needed(const VMValue& value, MemoryManager& memory, const WellKnownClasses& classes);

/**
 * Try to convert heap object to immediate (unboxing)
 * Returns original VMValue if cannot unbox
 */
VMValue unbox_if_possible(const VMValue& value, const WellKnownClasses& classes);

// === OBJECT CREATION HELPERS ===

/**
 * Create objects based on Smalltalk class, not C++ type
 * The VM determines layout from the class format specification
 */
Object* create_object(SmalltalkClass* smalltalk_class, uint32_t size, MemoryManager& memory);

/**
 * Create specific object types - VM knows how to format them
 */
Object* create_array(uint32_t size, MemoryManager& memory, const WellKnownClasses& classes);
Object* create_string(const std::string& content, MemoryManager& memory, const WellKnownClasses& classes);
Object* create_symbol(const std::string& content, MemoryManager& memory, const WellKnownClasses& classes);
Object* create_byte_array(uint32_t size, MemoryManager& memory, const WellKnownClasses& classes);

// === OBJECT ACCESS HELPERS ===

/**
 * Array operations - VM checks that object is actually an array
 */
VMValue array_at(Object* array_obj, uint32_t index, const WellKnownClasses& classes);
void array_at_put(Object* array_obj, uint32_t index, const VMValue& value, const WellKnownClasses& classes);
uint32_t array_size(Object* array_obj, const WellKnownClasses& classes);

/**
 * String operations - VM checks that object is actually a string  
 */
std::string string_content(Object* string_obj, const WellKnownClasses& classes);
void string_set_content(Object* string_obj, const std::string& content, const WellKnownClasses& classes);

/**
 * ByteArray operations
 */
uint8_t byte_array_at(Object* byte_array_obj, uint32_t index, const WellKnownClasses& classes);
void byte_array_at_put(Object* byte_array_obj, uint32_t index, uint8_t value, const WellKnownClasses& classes);

/**
 * Instance variable access - works for any object
 */
VMValue get_instance_variable(Object* obj, uint32_t index);
void set_instance_variable(Object* obj, uint32_t index, const VMValue& value);

// === DEBUGGING AND INSPECTION ===

/**
 * Get string representation of any value for debugging
 */
std::string value_to_string(const VMValue& value, const WellKnownClasses& classes);

/**
 * Get detailed object information for debugging
 */
struct ObjectInfo {
    std::string class_name;
    uint32_t size;
    ObjectFormat format;
    std::vector<std::string> instance_variables;
    bool is_immediate;
};

ObjectInfo inspect_object(const VMValue& value, const WellKnownClasses& classes);

// === EXAMPLE VM PRIMITIVE IMPLEMENTATION ===

/**
 * Example of how VM primitives work with the new model
 * No more switching on C++ types - just check Smalltalk classes
 */
class VMPrimitives {
public:
    VMPrimitives(const WellKnownClasses& classes) : classes_(classes) {}
    
    // Arithmetic primitive - handles both immediates and heap objects
    VMValue primitive_add(const VMValue& left, const VMValue& right, MemoryManager& memory);
    
    // Comparison primitive 
    VMValue primitive_less_than(const VMValue& left, const VMValue& right);
    
    // Object creation
    VMValue primitive_new(SmalltalkClass* class_obj, MemoryManager& memory);
    VMValue primitive_new_with_size(SmalltalkClass* class_obj, uint32_t size, MemoryManager& memory);
    
    // Array primitives
    VMValue primitive_array_at(const VMValue& array, const VMValue& index);
    VMValue primitive_array_at_put(const VMValue& array, const VMValue& index, const VMValue& value);
    VMValue primitive_array_size(const VMValue& array);
    
private:
    const WellKnownClasses& classes_;
    
    // Helper to ensure numeric arguments
    bool ensure_numeric(const VMValue& value, double& result);
};

// === MIGRATION HELPERS ===

/**
 * These functions help migrate from the old inheritance-based model
 * They will be removed once migration is complete
 */
namespace migration {
    
    // Convert old SmallInteger* to new model
    VMValue from_old_small_integer(void* old_small_integer);
    
    // Convert old Boolean* to new model  
    VMValue from_old_boolean(void* old_boolean);
    
    // Convert old Object* to new model (for heap objects)
    VMValue from_old_object(void* old_object);
}

} // namespace smalltalk