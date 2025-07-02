#include "simple_object.h"
#include <cstring>

namespace smalltalk {

/**
 * Check if an object is an instance of a specific Smalltalk class
 * This walks the class hierarchy to check for inheritance
 */
bool is_instance_of(const Object* obj, const SmalltalkClass* target_class) {
    if (!obj || !target_class) {
        return false;
    }
    
    SmalltalkClass* obj_class = obj->get_class();
    if (!obj_class) {
        return false;
    }
    
    // Check direct class match
    if (obj_class == target_class) {
        return true;
    }
    
    // Check inheritance hierarchy
    return obj_class->is_subclass_of(target_class);
}

/**
 * Check if an object "understands" a message (has the method)
 * This checks the object's class method dictionary
 */
bool understands(const Object* obj, const char* selector) {
    if (!obj || !selector) {
        return false;
    }
    
    SmalltalkClass* obj_class = obj->get_class();
    if (!obj_class) {
        return false;
    }
    
    // Try to lookup the method
    void* method = obj_class->lookup_method(selector);
    return method != nullptr;
}

/**
 * Get the format of an object based on its Smalltalk class
 * This determines how the object's data should be interpreted
 */
ObjectFormat get_object_format(const Object* obj) {
    if (!obj) {
        return ObjectFormat::REGULAR;
    }
    
    SmalltalkClass* obj_class = obj->get_class();
    if (!obj_class) {
        return ObjectFormat::REGULAR;
    }
    
    return obj_class->format();
}

} // namespace smalltalk