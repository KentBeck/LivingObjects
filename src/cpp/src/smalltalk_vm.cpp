#include "smalltalk_vm.h"
#include "smalltalk_class.h"
#include "primitives.h"
#include <iostream>

namespace smalltalk {

bool SmalltalkVM::initialized = false;

void SmalltalkVM::initialize() {
    if (initialized) return;
    
    // Initialize core classes first
    ClassUtils::initializeCoreClasses();
    
    // Initialize primitive registry
    PrimitiveRegistry::getInstance().initializeCorePrimitives();
    
    // Initialize primitive methods
    Primitives::initialize();
    
    initialized = true;
}

bool SmalltalkVM::isInitialized() {
    return initialized;
}

void SmalltalkVM::shutdown() {
    // Cleanup code would go here
    initialized = false;
}

} // namespace smalltalk
