#include "memory_manager.h"
#include "object.h"
#include "smalltalk_class.h"
#include <cassert>
#include <iostream>

using namespace smalltalk;

void testBasicClassAllocation() {
  std::cout << "Testing basic class allocation..." << std::endl;

  // Initialize the class system
  ClassUtils::initializeCoreClasses();

  MemoryManager memoryManager;

  // Get the Object class
  Class *objectClass = ClassUtils::getObjectClass();
  assert(objectClass != nullptr);

  // Allocate an instance of Object
  Object *obj = memoryManager.allocateInstance(objectClass);
  assert(obj != nullptr);
  assert(obj->getClass() == objectClass);

  std::cout << "  ✓ Object instance allocated successfully" << std::endl;
  std::cout << "    Instance size: " << objectClass->getInstanceSize()
            << std::endl;
  std::cout << "    Format: " << static_cast<int>(objectClass->getFormat())
            << std::endl;
}

void testCustomClassWithInstanceVariables() {
  std::cout << "\nTesting custom class with instance variables..." << std::endl;

  // Create a custom Point class
  Class *objectClass = ClassUtils::getObjectClass();
  Class *pointClass = ClassUtils::createClass("Point", objectClass);

  // Add instance variables
  pointClass->addInstanceVariable("x");
  pointClass->addInstanceVariable("y");

  std::cout << "  Point class created with " << pointClass->getInstanceSize()
            << " instance variables" << std::endl;

  MemoryManager memoryManager;

  // Allocate an instance of Point
  Object *point = memoryManager.allocateInstance(pointClass);
  assert(point != nullptr);
  assert(point->getClass() == pointClass);

  std::cout << "  ✓ Point instance allocated successfully" << std::endl;

  // Verify the instance has the correct number of slots
  assert(pointClass->getInstanceSize() == 2);
  assert(pointClass->getInstanceVariableIndex("x") == 0);
  assert(pointClass->getInstanceVariableIndex("y") == 1);

  std::cout << "  ✓ Instance variables accessible" << std::endl;
}

void testArrayClassAllocation() {
  std::cout << "\nTesting Array class allocation..." << std::endl;

  // Create an Array class
  Class *objectClass = ClassUtils::getObjectClass();
  Class *arrayClass = ClassUtils::createClass("Array", objectClass);
  arrayClass->setFormat(ObjectFormat::INDEXABLE_OBJECTS);

  MemoryManager memoryManager;

  // Allocate an array with 5 elements
  Object *array = memoryManager.allocateIndexableInstance(arrayClass, 5);
  assert(array != nullptr);
  assert(array->getClass() == arrayClass);

  std::cout << "  ✓ Array instance allocated successfully" << std::endl;
  std::cout << "    Indexable format: " << arrayClass->isIndexable()
            << std::endl;
  std::cout << "    Is pointer format: " << arrayClass->isPointerFormat()
            << std::endl;
}

void testByteArrayClassAllocation() {
  std::cout << "\nTesting ByteArray class allocation..." << std::endl;

  // Create a ByteArray class
  Class *objectClass = ClassUtils::getObjectClass();
  Class *byteArrayClass = ClassUtils::createClass("ByteArray", objectClass);
  byteArrayClass->setFormat(ObjectFormat::BYTE_INDEXABLE);

  MemoryManager memoryManager;

  // Allocate a byte array with 100 bytes
  Object *byteArray =
      memoryManager.allocateByteIndexableInstance(byteArrayClass, 100);
  assert(byteArray != nullptr);
  assert(byteArray->getClass() == byteArrayClass);

  std::cout << "  ✓ ByteArray instance allocated successfully" << std::endl;
  std::cout << "    Byte indexable: " << byteArrayClass->isByteIndexable()
            << std::endl;
  std::cout << "    Is indexable: " << byteArrayClass->isIndexable()
            << std::endl;
}

void testStringClassAllocation() {
  std::cout << "\nTesting String class allocation..." << std::endl;

  // Get the String class (already created in initializeCoreClasses)
  Class *stringClass = ClassUtils::getStringClass();
  assert(stringClass != nullptr);
  assert(stringClass->isByteIndexable());

  MemoryManager memoryManager;

  // Allocate a string with 20 character capacity
  Object *string = memoryManager.allocateByteIndexableInstance(stringClass, 20);
  assert(string != nullptr);
  assert(string->getClass() == stringClass);

  std::cout << "  ✓ String instance allocated successfully" << std::endl;
  std::cout << "    Format: " << static_cast<int>(stringClass->getFormat())
            << std::endl;
}

void testClassHierarchyWithInstanceSizes() {
  std::cout << "\nTesting class hierarchy with instance sizes..." << std::endl;

  // Create a hierarchy: Object -> Shape -> Rectangle
  Class *objectClass = ClassUtils::getObjectClass();

  Class *shapeClass = ClassUtils::createClass("Shape", objectClass);
  shapeClass->addInstanceVariable("color");

  Class *rectangleClass = ClassUtils::createClass("Rectangle", shapeClass);
  rectangleClass->addInstanceVariable("width");
  rectangleClass->addInstanceVariable("height");

  std::cout << "  Object instance size: " << objectClass->getInstanceSize()
            << std::endl;
  std::cout << "  Shape instance size: " << shapeClass->getInstanceSize()
            << std::endl;
  std::cout << "  Rectangle instance size: "
            << rectangleClass->getInstanceSize() << std::endl;

  // Rectangle should inherit Shape's instance variables plus its own
  assert(objectClass->getInstanceSize() == 0);
  assert(shapeClass->getInstanceSize() == 1);     // color
  assert(rectangleClass->getInstanceSize() == 3); // color + width + height

  MemoryManager memoryManager;

  // Allocate instances
  Object *shape = memoryManager.allocateInstance(shapeClass);
  Object *rectangle = memoryManager.allocateInstance(rectangleClass);

  assert(shape->getClass() == shapeClass);
  assert(rectangle->getClass() == rectangleClass);

  std::cout << "  ✓ Class hierarchy allocation works correctly" << std::endl;
}

void testErrorConditions() {
  std::cout << "\nTesting error conditions..." << std::endl;

  MemoryManager memoryManager;

  // Test null class
  try {
    memoryManager.allocateInstance(nullptr);
    assert(false && "Should have thrown exception");
  } catch (const std::runtime_error &e) {
    std::cout << "  ✓ Null class allocation properly throws: " << e.what()
              << std::endl;
  }

  // Test indexable allocation on non-indexable class
  Class *objectClass = ClassUtils::getObjectClass();
  try {
    memoryManager.allocateIndexableInstance(objectClass, 10);
    assert(false && "Should have thrown exception");
  } catch (const std::runtime_error &e) {
    std::cout
        << "  ✓ Non-indexable class indexable allocation properly throws: "
        << e.what() << std::endl;
  }

  // Test byte indexable allocation on non-byte-indexable class
  try {
    memoryManager.allocateByteIndexableInstance(objectClass, 100);
    assert(false && "Should have thrown exception");
  } catch (const std::runtime_error &e) {
    std::cout
        << "  ✓ Non-byte-indexable class byte allocation properly throws: "
        << e.what() << std::endl;
  }
}

int main() {
  try {
    std::cout << "=== Class-Based Allocation Tests ===" << std::endl;

    testBasicClassAllocation();
    testCustomClassWithInstanceVariables();
    testArrayClassAllocation();
    testByteArrayClassAllocation();
    testStringClassAllocation();
    testClassHierarchyWithInstanceSizes();
    testErrorConditions();

    std::cout << "\n🎉 All class allocation tests passed!" << std::endl;
    return 0;
  } catch (const std::exception &e) {
    std::cerr << "Test failed with exception: " << e.what() << std::endl;
    return 1;
  }
}