#pragma once

#include "object.h"
#include "context.h"
#include <cstddef>
#include <vector>

namespace smalltalk {

class MemoryManager {
public:
    // Constructor - initializes memory spaces
    MemoryManager(size_t initialSpaceSize = 1024 * 1024);
    
    // Destructor - cleans up memory
    ~MemoryManager();
    
    // Allocation methods
    Object* allocateObject(ObjectType type, size_t size);
    Object* allocateBytes(size_t byteSize);
    Object* allocateArray(size_t length);
    
    // Context allocation
    MethodContext* allocateMethodContext(size_t size, uint32_t method, Object* self, Object* sender);
    BlockContext* allocateBlockContext(size_t size, uint32_t method, Object* self, Object* sender, Object* home);
    
    // Stack chunk allocation
    StackChunk* allocateStackChunk(size_t size);
    
    // Garbage collection
    void collectGarbage();
    
    // Memory statistics
    size_t getFreeSpace() const;
    size_t getTotalSpace() const;
    size_t getUsedSpace() const;
    
private:
    // Memory spaces
    void* fromSpace;
    void* toSpace;
    size_t spaceSize;
    
    // Current allocation pointer
    void* currentAllocation;
    
    // Root set for GC
    std::vector<Object**> roots;
    
    // Register a root for GC
    void addRoot(Object** root);
    void removeRoot(Object** root);
    
    // Forward an object during GC
    Object* forwardObject(Object* obj);
    
    // Copy an object during GC
    Object* copyObject(Object* obj);
    
    // Scan an object during GC
    void scanObject(Object* obj);
    
    // Flip spaces after GC
    void flipSpaces();
};

} // namespace smalltalk