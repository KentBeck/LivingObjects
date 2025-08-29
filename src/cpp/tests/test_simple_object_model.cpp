#include <cassert>
#include <iostream>
#include <string>
#include <vector>

// Include our new simplified headers
#include "../include/simple_object.h"
#include "../include/simple_tagged_value.h"
#include "../include/vm_support.h"

using namespace smalltalk;

// Mock SmalltalkClass for testing
class MockSmalltalkClass : public SmalltalkClass {
public:
  std::string class_name;
  ObjectFormat object_format;
  uint32_t inst_size;

  MockSmalltalkClass(const std::string &class_name, ObjectFormat fmt,
                     uint32_t size = 0)
      : class_name(class_name), object_format(fmt), inst_size(size) {}

  // Implement pure virtual methods
  const char *name() const override { return class_name.c_str(); }
  ObjectFormat format() const override { return object_format; }
  uint32_t instance_size() const override { return inst_size; }
  SmalltalkClass *superclass() const override { return nullptr; }
  bool is_subclass_of(const SmalltalkClass * /*other*/) const override {
    return false;
  }
  void *lookup_method(const char * /*selector*/) const override {
    return nullptr;
  }
};

// Simple memory allocator for testing
class TestMemoryManager {
public:
  std::vector<void *> allocated_objects;

  void *allocate(size_t bytes) {
    void *ptr = std::malloc(bytes);
    allocated_objects.push_back(ptr);
    return ptr;
  }

  ~TestMemoryManager() {
    for (void *ptr : allocated_objects) {
      std::free(ptr);
    }
  }
};

// Forward declarations (after TestMemoryManager is defined)
VMValue create_test_array(uint32_t size, TestMemoryManager &memory,
                          const WellKnownClasses &classes);
Object *create_simple_object(SmalltalkClass *st_class, uint32_t size,
                             TestMemoryManager &memory);

// Test the new TaggedValue system
void test_tagged_values() {
  std::cout << "=== Testing TaggedValue System ===" << std::endl;

  // Create immediates - these are NEVER heap objects
  TaggedValue nil_val = TaggedValue::nil();
  TaggedValue true_val = TaggedValue::true_value();
  TaggedValue false_val = TaggedValue::false_value();
  TaggedValue int_val = TaggedValue::integer(42);
  TaggedValue float_val = TaggedValue::float_value(3.14f);

  // Test type checking
  assert(nil_val.is_nil());
  assert(true_val.is_boolean() && true_val.is_true());
  assert(false_val.is_boolean() && false_val.is_false());
  assert(int_val.is_integer());
  assert(float_val.is_float());

  // Test value extraction
  assert(true_val.as_boolean() == true);
  assert(false_val.as_boolean() == false);
  assert(int_val.as_integer() == 42);
  assert(float_val.as_float() == 3.14f);

  // Test comparison
  assert(TaggedValue::integer(42) == TaggedValue::integer(42));
  assert(TaggedValue::integer(42) != TaggedValue::integer(43));
  assert(TaggedValue::true_value() == TaggedValue::true_value());

  // Test string representation for debugging
  assert(nil_val.to_string() == "nil");
  assert(true_val.to_string() == "true");
  assert(int_val.to_string() == "42");

  std::cout << "✓ TaggedValue tests passed" << std::endl;
}

// Test the new uniform Object system
void test_uniform_objects() {
  std::cout << "=== Testing Uniform Object System ===" << std::endl;

  TestMemoryManager memory;

  // Create mock Smalltalk classes
  MockSmalltalkClass *array_class =
      new MockSmalltalkClass("Array", ObjectFormat::ARRAY);
  MockSmalltalkClass *string_class =
      new MockSmalltalkClass("String", ObjectFormat::BYTE_ARRAY);
  MockSmalltalkClass *point_class =
      new MockSmalltalkClass("Point", ObjectFormat::REGULAR, 2);

  // Test 1: Create an Array object (holds pointers)
  size_t array_size = 3;
  size_t array_bytes =
      object_size_bytes(array_size, false); // false = slot object
  Object *array_obj = static_cast<Object *>(memory.allocate(array_bytes));
  new (array_obj) Object(array_class, array_size);

  // Arrays hold pointers to other objects/values
  array_obj->set_slot(0, nullptr);                           // nil equivalent
  array_obj->set_slot(1, array_obj);                         // self-reference
  array_obj->set_slot(2, reinterpret_cast<void *>(0x12345)); // some pointer

  assert(array_obj->get_class() == array_class);
  assert(array_obj->size() == 3);
  assert(array_obj->slot(1) == array_obj);

  // Test 2: Create a String object (holds bytes)
  std::string content = "Hello";
  size_t string_bytes =
      object_size_bytes(content.length(), true); // true = byte object
  Object *string_obj = static_cast<Object *>(memory.allocate(string_bytes));
  new (string_obj) Object(string_class, content.length());

  // Copy string content into the object
  std::memcpy(string_obj->bytes(), content.data(), content.length());

  assert(string_obj->get_class() == string_class);
  assert(string_obj->size() == 5);
  assert(string_obj->byte_at(0) == 'H');
  assert(string_obj->byte_at(4) == 'o');

  // Test 3: Create a regular object with instance variables (Point with x, y)
  size_t point_bytes = object_size_bytes(2, false); // 2 instance variables
  Object *point_obj = static_cast<Object *>(memory.allocate(point_bytes));
  new (point_obj) Object(point_class, 2);

  // Set instance variables (x=10, y=20) - stored as pointers to TaggedValues
  // In real implementation, these would be proper references
  point_obj->set_slot(0, reinterpret_cast<void *>(10)); // x
  point_obj->set_slot(1, reinterpret_cast<void *>(20)); // y

  assert(point_obj->get_class() == point_class);
  assert(point_obj->size() == 2);
  assert(reinterpret_cast<intptr_t>(point_obj->slot(0)) == 10);
  assert(reinterpret_cast<intptr_t>(point_obj->slot(1)) == 20);

  std::cout << "✓ Uniform Object tests passed" << std::endl;

  delete array_class;
  delete string_class;
  delete point_class;
}

// Test VMValue - the unified value type
void test_vm_values() {
  std::cout << "=== Testing VMValue System ===" << std::endl;

  TestMemoryManager memory;
  MockSmalltalkClass *array_class =
      new MockSmalltalkClass("Array", ObjectFormat::ARRAY);

  // Create some test values
  VMValue nil_value(TaggedValue::nil());
  VMValue int_value(TaggedValue::integer(123));
  VMValue bool_value(TaggedValue::true_value());

  // Create a heap object
  size_t obj_bytes = object_size_bytes(1, false);
  Object *heap_obj = static_cast<Object *>(memory.allocate(obj_bytes));
  new (heap_obj) Object(array_class, 1);
  VMValue heap_value(heap_obj);

  // Test type discrimination
  assert(nil_value.is_immediate());
  assert(int_value.is_immediate());
  assert(bool_value.is_immediate());
  assert(heap_value.is_heap_object());

  // Test value extraction
  assert(nil_value.as_immediate().is_nil());
  assert(int_value.as_immediate().as_integer() == 123);
  assert(bool_value.as_immediate().is_true());
  assert(heap_value.as_object() == heap_obj);

  // Test convenience methods
  assert(nil_value.is_nil());
  assert(int_value.is_integer());
  assert(bool_value.is_boolean());

  std::cout << "✓ VMValue tests passed" << std::endl;

  delete array_class;
}

// Demonstrate how VM operations work without C++ inheritance
void test_vm_operations() {
  std::cout << "=== Testing VM Operations ===" << std::endl;

  TestMemoryManager memory;

  // Create mock well-known classes
  WellKnownClasses classes;
  classes.Array = new MockSmalltalkClass("Array", ObjectFormat::ARRAY);
  classes.String = new MockSmalltalkClass("String", ObjectFormat::BYTE_ARRAY);
  classes.SmallInteger =
      new MockSmalltalkClass("SmallInteger", ObjectFormat::REGULAR);

  // Test: VM handles array operations based on Smalltalk class, not C++ type
  VMValue array_val = create_test_array(3, memory, classes);

  // In the real VM, these operations check the Smalltalk class:
  // if (is_instance_of(array_val, classes.Array)) { ... }
  // NOT: if (dynamic_cast<ArrayObject*>(...))

  assert(array_val.is_heap_object());
  Object *array_obj = array_val.as_object();
  assert(array_obj->get_class() == classes.Array);

  // VM operations work uniformly on Object* regardless of what it represents
  array_obj->set_slot(0, reinterpret_cast<void *>(1)); // Store "1"
  array_obj->set_slot(1, reinterpret_cast<void *>(2)); // Store "2"
  array_obj->set_slot(2, reinterpret_cast<void *>(3)); // Store "3"

  assert(reinterpret_cast<intptr_t>(array_obj->slot(0)) == 1);
  assert(reinterpret_cast<intptr_t>(array_obj->slot(1)) == 2);
  assert(reinterpret_cast<intptr_t>(array_obj->slot(2)) == 3);

  std::cout << "✓ VM Operations tests passed" << std::endl;

  delete classes.Array;
  delete classes.String;
  delete classes.SmallInteger;
}

// Helper function for testing
VMValue create_test_array(uint32_t size, TestMemoryManager &memory,
                          const WellKnownClasses &classes) {
  size_t bytes = object_size_bytes(size, false);
  Object *obj = static_cast<Object *>(memory.allocate(bytes));
  new (obj) Object(classes.Array, size);
  return VMValue(obj);
}

Object *create_simple_object(SmalltalkClass *st_class, uint32_t size,
                             TestMemoryManager &memory) {
  size_t bytes = object_size_bytes(size, false);
  Object *obj = static_cast<Object *>(memory.allocate(bytes));
  new (obj) Object(st_class, size);
  return obj;
}

// Demonstrate the key benefit: No more parallel C++ hierarchies!
void test_no_parallel_hierarchy() {
  std::cout << "=== Testing: No Parallel C++ Hierarchy ===" << std::endl;

  TestMemoryManager memory;

  // OLD WAY (what we're replacing):
  // SmallInteger* old_int = new SmallInteger(42, integer_class);
  // Boolean* old_bool = new Boolean(true, boolean_class);
  // Array* old_array = new Array(3, array_class);
  // String* old_string = new String("hello", string_class);
  // Each has different C++ type, complex virtual dispatch, etc.

  // NEW WAY: Everything is just Object* + Smalltalk class
  MockSmalltalkClass *integer_class =
      new MockSmalltalkClass("SmallInteger", ObjectFormat::REGULAR);
  MockSmalltalkClass *boolean_class =
      new MockSmalltalkClass("Boolean", ObjectFormat::REGULAR);
  MockSmalltalkClass *array_class =
      new MockSmalltalkClass("Array", ObjectFormat::ARRAY);
  MockSmalltalkClass *string_class =
      new MockSmalltalkClass("String", ObjectFormat::BYTE_ARRAY);

  // But integers and booleans are immediates - no heap objects needed!
  TaggedValue int_val = TaggedValue::integer(42);
  TaggedValue bool_val = TaggedValue::boolean(true);

  // Only collections and complex objects become heap objects
  Object *array_obj = create_simple_object(array_class, 3, memory);
  Object *string_obj = create_simple_object(string_class, 5, memory);

  // VM treats them all uniformly - checks Smalltalk class when needed
  std::vector<VMValue> values = {VMValue(int_val), VMValue(bool_val),
                                 VMValue(array_obj), VMValue(string_obj)};

  // Process all values uniformly
  for (const auto &value : values) {
    if (value.is_immediate()) {
      std::cout << "  Immediate: " << value.as_immediate().to_string()
                << std::endl;
    } else {
      Object *obj = value.as_object();
      MockSmalltalkClass *mock_class =
          static_cast<MockSmalltalkClass *>(obj->get_class());
      std::cout << "  Heap Object: " << mock_class->class_name
                << " (size=" << obj->size() << ")" << std::endl;
    }
  }

  std::cout << "✓ No Parallel Hierarchy test passed" << std::endl;

  delete integer_class;
  delete boolean_class;
  delete array_class;
  delete string_class;
}

int main() {
  std::cout << "Testing Simplified Smalltalk Object Model" << std::endl;
  std::cout << "=========================================" << std::endl;

  try {
    test_tagged_values();
    test_uniform_objects();
    test_vm_values();
    test_vm_operations();
    test_no_parallel_hierarchy();

    std::cout << std::endl;
    std::cout << "🎉 All tests passed!" << std::endl;
    std::cout << std::endl;
    std::cout << "KEY BENEFITS DEMONSTRATED:" << std::endl;
    std::cout << "• Immediates are TaggedValues only - never heap allocated"
              << std::endl;
    std::cout << "• All heap objects use uniform Object struct" << std::endl;
    std::cout << "• VM behavior based on Smalltalk class, not C++ type"
              << std::endl;
    std::cout << "• No complex inheritance hierarchy to maintain" << std::endl;
    std::cout << "• Easier debugging, GC, and serialization" << std::endl;
    std::cout << "• Ready for JIT compilation" << std::endl;

  } catch (const std::exception &e) {
    std::cerr << "Test failed: " << e.what() << std::endl;
    return 1;
  }

  return 0;
}