#include "smalltalk_vm.h"
#include "globals.h"
#include "primitives.h"
#include "smalltalk_class.h"
#include "symbol.h"
#include <iostream>

namespace smalltalk {

bool SmalltalkVM::initialized = false;

void SmalltalkVM::initialize() {
  if (initialized)
    return;

  // Initialize core classes first
  ClassUtils::initializeCoreClasses();

  // Initialize primitive registry
  PrimitiveRegistry::getInstance().initializeCorePrimitives();

  // Initialize primitive methods
  Primitives::initialize();

  // Create the global Smalltalk dictionary and populate with core classes
  if (!Globals::isInitialized()) {
    Class *dictClass = ClassRegistry::getInstance().getClass("Dictionary");
    if (dictClass) {
      // Allocate a Dictionary instance and install as Smalltalk
      Object *smalltalkDict = dictClass->createInstance();
      Globals::setSmalltalk(smalltalkDict);

      // Make Smalltalk refer to itself and register classes
      Globals::set("Smalltalk", smalltalkDict);
      for (Class *cls : ClassRegistry::getInstance().getAllClasses()) {
        Globals::set(cls->getName(), cls);
      }
    }
  }

  initialized = true;
}

bool SmalltalkVM::isInitialized() { return initialized; }

void SmalltalkVM::shutdown() {
  // Cleanup code would go here
  initialized = false;
}

} // namespace smalltalk
