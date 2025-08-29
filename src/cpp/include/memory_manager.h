#pragma once

#include "context.h"
#include "object.h"

#include <cstddef>
#include <memory>
#include <vector>

namespace smalltalk {

// Forward declaration
class Class;

/**
 * MemoryManager handles object allocation and garbage collection.
 *
 * This class manages memory spaces for Smalltalk objects using a
 * stop-and-copy garbage collection algorithm. It provides methods
 * for allocating objects, contexts, and other Smalltalk structures,
 * as well as performing garbage collection when necessary.
 */
class MemoryManager {
public:
  /**
   * Constructor that initializes memory spaces.
   *
   * @param initialSpaceSize The size of each memory space in bytes.
   */
  static constexpr size_t DEFAULT_INITIAL_SPACE_SIZE =
      static_cast<size_t>(1024) * 1024; // 1MB

  /**
   * Constructor that initializes memory spaces.
   *
   * @param initialSpaceSize The size of each memory space in bytes.
   */
  MemoryManager(size_t initialSpaceSize = DEFAULT_INITIAL_SPACE_SIZE);

  /**
   * Destructor that cleans up memory.
   */
  ~MemoryManager();

  // Allocation methods

  /**
   * Allocates a new object of the specified type and size.
   *
   * @param type The type of object to allocate.
   * @param size The size of the object in slots.
   * @return A pointer to the allocated object.
   */
  Object *allocateObject(ObjectType type, size_t size);

  /**
   * Allocates a new instance of the specified class.
   *
   * @param clazz The class to create an instance of.
   * @return A pointer to the allocated instance.
   */
  Object *allocateInstance(Class *clazz);

  /**
   * Allocates a new indexable instance of the specified class.
   *
   * @param clazz The class to create an instance of.
   * @param indexedSize The number of indexed slots.
   * @return A pointer to the allocated instance.
   */
  Object *allocateIndexableInstance(Class *clazz, size_t indexedSize);

  /**
   * Allocates a new byte-indexable instance of the specified class.
   *
   * @param clazz The class to create an instance of.
   * @param byteSize The number of bytes in the indexed portion.
   * @return A pointer to the allocated instance.
   */
  Object *allocateByteIndexableInstance(Class *clazz, size_t byteSize);

  /**
   * Allocates a byte array of the specified size.
   *
   * @param byteSize The size of the byte array in bytes.
   * @return A pointer to the allocated byte array.
   */
  Object *allocateBytes(size_t byteSize);

  /**
   * Allocates an array of the specified length.
   *
   * @param length The number of slots in the array.
   * @return A pointer to the allocated array.
   */
  Object *allocateArray(size_t length);

  // Allocate boxed primitive values
  Object *allocateInteger(int32_t value);
  Object *allocateBoolean(bool value);

  // Context allocation

  /**
   * Allocates a method context.
   *
   * @param size The size of the context in slots.
   * @param method The method reference.
   * @param self The receiver object.
   * @param sender The sender context.
   * @return A pointer to the allocated method context.
   */
  MethodContext *allocateMethodContext(size_t size, TaggedValue self,
                                       TaggedValue sender, TaggedValue home,
                                       CompiledMethod *compiledMethod);

  /**
   * Allocates a block context.
   *
   * @param size The size of the context in slots.
   * @param method The method reference.
   * @param self The receiver object.
   * @param sender The sender context.
   * @param home The home context.
   * @return A pointer to the allocated block context.
   */
  BlockContext *allocateBlockContext(size_t size, TaggedValue self,
                                     TaggedValue sender, TaggedValue home);

  /**
   * Allocates a stack chunk.
   *
   * @param size The size of the chunk in slots.
   * @return A pointer to the allocated stack chunk.
   */
  StackChunk *allocateStackChunk(size_t size);

  /**
   * Performs garbage collection.
   */
  void collectGarbage();

  // Memory statistics

  /**
   * Gets the amount of free space available.
   *
   * @return The amount of free space in bytes.
   */
  size_t getFreeSpace() const;

  /**
   * Gets the total space available.
   *
   * @return The total space in bytes.
   */
  size_t getTotalSpace() const;

  /**
   * Gets the amount of used space.
   *
   * @return The amount of used space in bytes.
   */
  size_t getUsedSpace() const;

private:
  // Size must be declared before unique_ptr members that use it
  size_t spaceSize; ///< Size of each space in bytes

  // Memory spaces (RAII-managed)
  std::unique_ptr<void, decltype(&std::free)> fromSpacePtr;
  std::unique_ptr<void, decltype(&std::free)> toSpacePtr;

  // Raw pointers for compatibility
  void *fromSpace; ///< Current allocation space
  void *toSpace;   ///< Space for copying during GC

  // Current allocation pointer
  void *currentAllocation;

  // Stack chunk management
  std::vector<std::unique_ptr<void, decltype(&std::free)>> stackChunks;

  // Root set for GC
  std::vector<Object **> roots;

  /**
   * Registers a root for garbage collection.
   *
   * @param root Pointer to the root object reference.
   */
  void addRoot(Object **root);

  /**
   * Unregisters a root from garbage collection.
   *
   * @param root Pointer to the root object reference.
   */
  void removeRoot(Object **root);

  /**
   * Forwards an object during garbage collection.
   *
   * If the object has already been forwarded, returns the forwarding
   * address. Otherwise, copies the object to the to-space and returns
   * the new address.
   *
   * @param obj The object to forward.
   * @return The forwarded object address.
   */
  Object *forwardObject(Object *obj);

  /**
   * Copies an object during garbage collection.
   *
   * @param obj The object to copy.
   * @return The copied object address.
   */
  Object *copyObject(Object *obj);

  /**
   * Scans an object during garbage collection.
   *
   * Iterates through all references in the object and forwards them.
   *
   * @param obj The object to scan.
   */
  void scanObject(Object *obj);

  /**
   * Flips the from-space and to-space after garbage collection.
   */
  void flipSpaces();
};

} // namespace smalltalk
