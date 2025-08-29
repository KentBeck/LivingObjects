#pragma once

#include <cassert>
#include <cstddef>
#include <cstdint>

namespace smalltalk {

// Forward declarations
class SmalltalkClass;

/**
 * Object format types - how objects store their data
 */
enum class ObjectFormat {
  REGULAR,        // Normal object with instance variables
  ARRAY,          // Indexable array of pointers
  BYTE_ARRAY,     // Indexable array of bytes
  COMPILED_METHOD // Special format for executable methods
};

/**
 * Base SmalltalkClass definition
 * This represents a Smalltalk class object that defines the behavior of
 * instances
 */
class SmalltalkClass {
public:
  virtual ~SmalltalkClass() = default;

  // Basic class information
  virtual const char *name() const = 0;
  virtual ObjectFormat format() const = 0;
  virtual uint32_t instance_size() const = 0;

  // Class hierarchy
  virtual SmalltalkClass *superclass() const = 0;
  virtual bool is_subclass_of(const SmalltalkClass *other) const = 0;

  // Method lookup (to be implemented by VM)
  virtual void *lookup_method(const char *selector) const = 0;
};

/**
 * Simplified Object Model for Smalltalk VM
 *
 * Key Design Principles:
 * 1. NO C++ inheritance hierarchy - single Object struct for all heap objects
 * 2. Immediates (integers, booleans, nil) are TaggedValues only, never heap
 * objects
 * 3. VM special handling based on Smalltalk class, not C++ type
 * 4. Uniform memory layout: [Object header][variable data...]
 * 5. Object behavior determined by Smalltalk class, not C++ type
 */

/**
 * Object header flags for GC and VM
 */
enum class ObjectFlag : uint8_t {
  MARKED = 0,      // GC mark bit
  REMEMBERED = 1,  // In remembered set (generational GC)
  IMMUTABLE = 2,   // Object cannot be modified
  FORWARDED = 3,   // Object has forwarding pointer (during GC)
  PINNED = 4,      // Cannot be moved by GC
  HAS_POINTERS = 5 // Contains references to other objects
};

/**
 * Uniform object header (64 bits total)
 * Used by ALL heap-allocated objects regardless of Smalltalk class
 */
struct ObjectHeader {
  uint32_t size;  // Size in slots (for arrays) or bytes (for byte objects)
  uint16_t flags; // ObjectFlag bits
  uint16_t hash;  // Identity hash (16 bits sufficient for most cases)

  // Flag operations
  void set_flag(ObjectFlag flag) { flags |= (1 << static_cast<uint8_t>(flag)); }

  bool has_flag(ObjectFlag flag) const {
    return (flags & (1 << static_cast<uint8_t>(flag))) != 0;
  }

  void clear_flag(ObjectFlag flag) {
    flags &= ~(1 << static_cast<uint8_t>(flag));
  }

  void clear_all_flags() { flags = 0; }
};

/**
 * Universal Object structure - used for ALL heap-allocated Smalltalk objects
 *
 * Memory Layout:
 * [ObjectHeader: 8 bytes][SmalltalkClass*: 8 bytes][variable data...]
 *
 * The variable data area interpretation depends on the Smalltalk class:
 * - Regular objects: instance variable slots
 * - Arrays: element pointers
 * - ByteArrays: raw bytes
 * - Strings: UTF-8 bytes
 * - CompiledMethods: bytecode + literals
 * - Classes: instance variables + method dictionary
 */
struct Object {
  ObjectHeader header;
  SmalltalkClass *smalltalk_class; // Every object knows its Smalltalk class

  // Variable data follows immediately after this struct
  // Accessed via helper methods below

  /**
   * Constructor for heap objects
   * @param st_class The Smalltalk class of this object
   * @param size Size in slots or bytes depending on object type
   * @param identity_hash Optional identity hash (0 = generate later)
   */
  Object(SmalltalkClass *st_class, uint32_t size, uint16_t identity_hash = 0)
      : header{size, 0, identity_hash}, smalltalk_class(st_class) {

    // Set HAS_POINTERS flag based on class (implementation will check class
    // format) For now, assume most objects have pointers
    header.set_flag(ObjectFlag::HAS_POINTERS);
  }

  // === ACCESSORS FOR OBJECT METADATA ===

  SmalltalkClass *get_class() const { return smalltalk_class; }

  uint32_t size() const { return header.size; }

  uint16_t identity_hash() const { return header.hash; }

  void set_identity_hash(uint16_t hash) { header.hash = hash; }

  // === ACCESSORS FOR VARIABLE DATA ===

  /**
   * Get pointer to variable data area
   * Use with caution - caller must know the object format
   */
  void *data() { return reinterpret_cast<char *>(this) + sizeof(Object); }

  const void *data() const {
    return reinterpret_cast<const char *>(this) + sizeof(Object);
  }

  /**
   * Access as slot array (for regular objects and arrays)
   * Each slot holds a pointer to another Object or TaggedValue
   */
  void **slots() { return static_cast<void **>(data()); }

  const void *const *slots() const {
    return static_cast<const void *const *>(data());
  }

  const void *slot(size_t index) const {
    assert(index < header.size);
    return slots()[index];
  }

  void *slot(size_t index) {
    assert(index < header.size);
    return slots()[index];
  }

  void set_slot(size_t index, void *value) {
    assert(index < header.size);
    slots()[index] = value;
  }

  /**
   * Access as byte array (for ByteArray, String, CompiledMethod bytecode)
   */
  uint8_t *bytes() { return static_cast<uint8_t *>(data()); }

  const uint8_t *bytes() const { return static_cast<const uint8_t *>(data()); }

  uint8_t byte_at(size_t index) const {
    assert(index < header.size);
    return bytes()[index];
  }

  void set_byte_at(size_t index, uint8_t value) {
    assert(index < header.size);
    bytes()[index] = value;
  }

  // === OBJECT IDENTITY ===

  bool operator==(const Object &other) const {
    return this == &other; // Smalltalk object identity
  }

  bool operator!=(const Object &other) const { return this != &other; }

  // === GC SUPPORT ===

  bool is_marked() const { return header.has_flag(ObjectFlag::MARKED); }

  void mark() { header.set_flag(ObjectFlag::MARKED); }

  void unmark() { header.clear_flag(ObjectFlag::MARKED); }

  bool has_pointers() const {
    return header.has_flag(ObjectFlag::HAS_POINTERS);
  }

  bool is_forwarded() const { return header.has_flag(ObjectFlag::FORWARDED); }

  void set_forwarding_address(Object *new_address) {
    header.set_flag(ObjectFlag::FORWARDED);
    // Store forwarding address in first slot
    set_slot(0, new_address);
  }

  Object *forwarding_address() {
    assert(is_forwarded());
    return static_cast<Object *>(slot(0));
  }
};

// === VM HELPER FUNCTIONS ===

/**
 * Calculate total object size including header and variable data
 */
inline size_t object_size_bytes(uint32_t data_size,
                                bool is_byte_object = false) {
  size_t base_size = sizeof(Object);
  if (is_byte_object) {
    // Byte objects: size is in bytes, add padding for alignment
    return base_size + ((data_size + 7) & ~7); // 8-byte align
  } else {
    // Slot objects: size is in pointer slots
    return base_size + (data_size * sizeof(void *));
  }
}

/**
 * Check if an object is an instance of a specific Smalltalk class
 * The VM will implement this by walking the class hierarchy
 */
bool is_instance_of(const Object *obj, const SmalltalkClass *target_class);

/**
 * Check if an object "understands" a message (has the method)
 * The VM will implement this by checking method dictionaries
 */
bool understands(const Object *obj, const char *selector);

/**
 * Get the format of an object (regular, byte array, etc.)
 * Based on the Smalltalk class format specification
 */
ObjectFormat get_object_format(const Object *obj);

} // namespace smalltalk