#include "memory_manager.h"
#include <cstring>
#include <algorithm>
#include <stdexcept>

namespace smalltalk {

MemoryManager::MemoryManager(size_t initialSpaceSize)
    : spaceSize(initialSpaceSize) {
    // Allocate memory spaces
    fromSpace = malloc(spaceSize);
    toSpace = malloc(spaceSize);
    
    if (!fromSpace || !toSpace) {
        throw std::runtime_error("Failed to allocate memory spaces");
    }
    
    // Initialize allocation pointer to the start of fromSpace
    currentAllocation = fromSpace;
}

MemoryManager::~MemoryManager() {
    // Free memory spaces
    free(fromSpace);
    free(toSpace);
}

Object* MemoryManager::allocateObject(ObjectType type, size_t size) {
    // Check if there's enough space
    size_t requiredBytes = sizeof(Object) + (size * sizeof(Object*));
    if (static_cast<char*>(fromSpace) + spaceSize - static_cast<char*>(currentAllocation) < requiredBytes) {
        // Not enough space, trigger garbage collection
        collectGarbage();
        
        // Check if there's still not enough space after GC
        if (static_cast<char*>(fromSpace) + spaceSize - static_cast<char*>(currentAllocation) < requiredBytes) {
            throw std::runtime_error("Out of memory");
        }
    }
    
    // Allocate the object
    Object* obj = static_cast<Object*>(currentAllocation);
    new (obj) Object(type, size);
    
    // Update allocation pointer
    currentAllocation = static_cast<char*>(currentAllocation) + requiredBytes;
    
    return obj;
}

Object* MemoryManager::allocateBytes(size_t byteSize) {
    // Calculate size with proper alignment
    size_t alignedSize = (byteSize + 7) & ~7; // Align to 8 bytes
    
    // Allocate as a byte array
    Object* obj = allocateObject(ObjectType::BYTE_ARRAY, alignedSize);
    
    // Initialize the byte array
    memset(reinterpret_cast<char*>(obj) + sizeof(Object), 0, alignedSize);
    
    return obj;
}

Object* MemoryManager::allocateArray(size_t length) {
    // Allocate as an array
    Object* obj = allocateObject(ObjectType::ARRAY, length);
    
    // Initialize the array slots to null
    Object** slots = reinterpret_cast<Object**>(reinterpret_cast<char*>(obj) + sizeof(Object));
    for (size_t i = 0; i < length; i++) {
        slots[i] = nullptr;
    }
    
    return obj;
}

MethodContext* MemoryManager::allocateMethodContext(size_t size, uint32_t method, Object* self, Object* sender) {
    // Check if there's enough space
    size_t requiredBytes = sizeof(MethodContext) + (size * sizeof(Object*));
    if (static_cast<char*>(fromSpace) + spaceSize - static_cast<char*>(currentAllocation) < requiredBytes) {
        // Not enough space, trigger garbage collection
        collectGarbage();
        
        // Check if there's still not enough space after GC
        if (static_cast<char*>(fromSpace) + spaceSize - static_cast<char*>(currentAllocation) < requiredBytes) {
            throw std::runtime_error("Out of memory");
        }
    }
    
    // Allocate the context
    MethodContext* context = static_cast<MethodContext*>(currentAllocation);
    new (context) MethodContext(size, method, self, sender);
    
    // Update allocation pointer
    currentAllocation = static_cast<char*>(currentAllocation) + requiredBytes;
    
    return context;
}

BlockContext* MemoryManager::allocateBlockContext(size_t size, uint32_t method, Object* self, Object* sender, Object* home) {
    // Check if there's enough space
    size_t requiredBytes = sizeof(BlockContext) + (size * sizeof(Object*));
    if (static_cast<char*>(fromSpace) + spaceSize - static_cast<char*>(currentAllocation) < requiredBytes) {
        // Not enough space, trigger garbage collection
        collectGarbage();
        
        // Check if there's still not enough space after GC
        if (static_cast<char*>(fromSpace) + spaceSize - static_cast<char*>(currentAllocation) < requiredBytes) {
            throw std::runtime_error("Out of memory");
        }
    }
    
    // Allocate the context
    BlockContext* context = static_cast<BlockContext*>(currentAllocation);
    new (context) BlockContext(size, method, self, sender, home);
    
    // Update allocation pointer
    currentAllocation = static_cast<char*>(currentAllocation) + requiredBytes;
    
    return context;
}

StackChunk* MemoryManager::allocateStackChunk(size_t size) {
    // For now, allocate stack chunks directly with malloc
    // In a real implementation, this would be managed differently
    size_t requiredBytes = sizeof(StackChunk) + (size * sizeof(Object*));
    StackChunk* chunk = static_cast<StackChunk*>(malloc(requiredBytes));
    
    if (!chunk) {
        throw std::runtime_error("Failed to allocate stack chunk");
    }
    
    new (chunk) StackChunk(size);
    chunk->allocationPointer = reinterpret_cast<char*>(chunk) + sizeof(StackChunk);
    
    return chunk;
}

void MemoryManager::collectGarbage() {
    // Simple implementation of a stop & copy garbage collector
    
    // Reset toSpace
    memset(toSpace, 0, spaceSize);
    
    // Initialize scan and free pointers in toSpace
    void* scanPtr = toSpace;
    void* freePtr = toSpace;
    
    // Forward all root objects
    for (Object** root : roots) {
        if (*root) {
            *root = forwardObject(*root);
        }
    }
    
    // Scan phase - scan all objects in toSpace
    while (scanPtr < freePtr) {
        Object* obj = static_cast<Object*>(scanPtr);
        scanObject(obj);
        
        // Move to next object
        size_t objSize;
        switch (static_cast<ObjectType>(obj->header.type)) {
            case ObjectType::BYTE_ARRAY:
                objSize = sizeof(Object) + obj->header.size;
                break;
            default:
                objSize = sizeof(Object) + (obj->header.size * sizeof(Object*));
                break;
        }
        
        // Align to 8 bytes
        objSize = (objSize + 7) & ~7;
        
        scanPtr = static_cast<char*>(scanPtr) + objSize;
    }
    
    // Flip spaces
    flipSpaces();
    
    // Reset allocation pointer
    currentAllocation = fromSpace;
}

Object* MemoryManager::forwardObject(Object* obj) {
    // Check if object is already forwarded
    if (obj->header.hasFlag(ObjectFlag::FORWARDED)) {
        // Get forwarding address from size field
        uintptr_t forwardingAddress = static_cast<uintptr_t>(obj->header.size);
        return reinterpret_cast<Object*>(forwardingAddress);
    }
    
    // Copy object to toSpace
    Object* newObj = copyObject(obj);
    
    // Mark original as forwarded and store forwarding address
    obj->header.setFlag(ObjectFlag::FORWARDED);
    obj->header.size = reinterpret_cast<uintptr_t>(newObj);
    
    return newObj;
}

Object* MemoryManager::copyObject(Object* obj) {
    // Calculate object size
    size_t objSize;
    switch (static_cast<ObjectType>(obj->header.type)) {
        case ObjectType::BYTE_ARRAY:
            objSize = sizeof(Object) + obj->header.size;
            break;
        default:
            objSize = sizeof(Object) + (obj->header.size * sizeof(Object*));
            break;
    }
    
    // Align to 8 bytes
    objSize = (objSize + 7) & ~7;
    
    // Check if there's enough space in toSpace
    if (static_cast<char*>(toSpace) + spaceSize - static_cast<char*>(currentAllocation) < objSize) {
        throw std::runtime_error("Out of memory during garbage collection");
    }
    
    // Copy object to toSpace
    Object* newObj = static_cast<Object*>(currentAllocation);
    memcpy(newObj, obj, objSize);
    
    // Clear forwarded flag in copy
    newObj->header.clearFlag(ObjectFlag::FORWARDED);
    
    // Update allocation pointer
    currentAllocation = static_cast<char*>(currentAllocation) + objSize;
    
    return newObj;
}

void MemoryManager::scanObject(Object* obj) {
    // Scan object fields for references
    if (static_cast<ObjectType>(obj->header.type) != ObjectType::BYTE_ARRAY) {
        // Get object slots
        Object** slots = reinterpret_cast<Object**>(reinterpret_cast<char*>(obj) + sizeof(Object));
        
        // Forward all slot references
        for (size_t i = 0; i < obj->header.size; i++) {
            if (slots[i]) {
                slots[i] = forwardObject(slots[i]);
            }
        }
    }
}

void MemoryManager::flipSpaces() {
    // Swap fromSpace and toSpace
    void* temp = fromSpace;
    fromSpace = toSpace;
    toSpace = temp;
}

void MemoryManager::addRoot(Object** root) {
    roots.push_back(root);
}

void MemoryManager::removeRoot(Object** root) {
    auto it = std::find(roots.begin(), roots.end(), root);
    if (it != roots.end()) {
        roots.erase(it);
    }
}

size_t MemoryManager::getFreeSpace() const {
    return static_cast<char*>(fromSpace) + spaceSize - static_cast<char*>(currentAllocation);
}

size_t MemoryManager::getTotalSpace() const {
    return spaceSize;
}

size_t MemoryManager::getUsedSpace() const {
    return static_cast<char*>(currentAllocation) - static_cast<char*>(fromSpace);
}

} // namespace smalltalk