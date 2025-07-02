#pragma once

#include "simple_object.h"
#include "simple_tagged_value.h"
#include "vm_support.h"
#include <memory>
#include <vector>
#include <cstddef>

namespace smalltalk {

/**
 * Simplified Memory Manager for Unified Object Model
 * 
 * Key improvements over the old memory manager:
 * 1. Single allocation method for all Object types
 * 2. No special case methods for each object type
 * 3. Understands TaggedValues (never allocates heap memory for immediates)
 * 4. GC works with uniform Object layout
 * 5. Builder pattern for complex object construction
 */
class SimpleMemoryManager {
public:
    static constexpr size_t DEFAULT_HEAP_SIZE = 64 * 1024 * 1024; // 64MB
    static constexpr size_t ALIGNMENT_BYTES = 8; // 8-byte alignment
    
    explicit SimpleMemoryManager(size_t heap_size = DEFAULT_HEAP_SIZE);
    ~SimpleMemoryManager();
    
    // === CORE ALLOCATION ===
    
    /**
     * Allocate a new Object with the specified Smalltalk class and size
     * @param smalltalk_class The Smalltalk class for this object
     * @param data_size Size in slots (for regular/array objects) or bytes (for byte objects)
     * @param is_byte_object true for ByteArray/String, false for regular/array objects
     * @return Pointer to allocated Object
     */
    Object* allocate_object(SmalltalkClass* smalltalk_class, uint32_t data_size, bool is_byte_object = false);
    
    /**
     * Convenience methods for common object types
     */
    Object* allocate_regular_object(SmalltalkClass* smalltalk_class, uint32_t instance_var_count);
    Object* allocate_array(SmalltalkClass* array_class, uint32_t element_count);
    Object* allocate_byte_array(SmalltalkClass* byte_array_class, uint32_t byte_count);
    Object* allocate_string(SmalltalkClass* string_class, uint32_t byte_count);
    
    // === OBJECT BUILDER PATTERN ===
    
    class ObjectBuilder {
    public:
        ObjectBuilder(SimpleMemoryManager& manager, SmalltalkClass* st_class);
        
        ObjectBuilder& with_instance_variables(uint32_t count);
        ObjectBuilder& with_array_elements(uint32_t count);
        ObjectBuilder& with_byte_data(uint32_t count);
        ObjectBuilder& with_identity_hash(uint16_t hash);
        ObjectBuilder& immutable();
        ObjectBuilder& pinned();
        
        Object* build();
        
    private:
        SimpleMemoryManager& manager_;
        SmalltalkClass* class_;
        uint32_t size_;
        bool is_byte_object_;
        uint16_t identity_hash_;
        bool immutable_;
        bool pinned_;
    };
    
    ObjectBuilder new_object(SmalltalkClass* smalltalk_class) {
        return ObjectBuilder(*this, smalltalk_class);
    }
    
    // === GARBAGE COLLECTION ===
    
    /**
     * Trigger garbage collection
     * @return Number of bytes freed
     */
    size_t collect_garbage();
    
    /**
     * Add a root object for GC (object that should never be collected)
     */
    void add_root(Object** root_ptr);
    void remove_root(Object** root_ptr);
    
    /**
     * Add a VMValue root (handles both immediates and heap objects)
     */
    void add_root(VMValue* root_value);
    void remove_root(VMValue* root_value);
    
    // === MEMORY STATISTICS ===
    
    size_t heap_size() const { return heap_size_; }
    size_t bytes_allocated() const;
    size_t bytes_free() const;
    double heap_utilization() const;
    
    // GC statistics
    size_t collection_count() const { return gc_count_; }
    size_t total_bytes_collected() const { return total_collected_; }
    
    // === CONFIGURATION ===
    
    /**
     * Set GC trigger threshold (0.0 - 1.0)
     * GC triggers when heap utilization exceeds this threshold
     */
    void set_gc_threshold(double threshold);
    
    /**
     * Enable/disable automatic GC on allocation
     */
    void set_auto_gc(bool enabled) { auto_gc_enabled_ = enabled; }
    
    // === DEBUGGING SUPPORT ===
    
    struct HeapStats {
        size_t total_objects;
        size_t regular_objects;
        size_t array_objects;
        size_t byte_objects;
        size_t total_bytes;
        size_t fragmentation_bytes;
    };
    
    HeapStats get_heap_stats() const;
    void dump_heap(const char* filename) const;
    bool validate_heap() const;
    
private:
    // Memory spaces for stop-and-copy GC
    size_t heap_size_;
    std::unique_ptr<char[]> from_space_;
    std::unique_ptr<char[]> to_space_;
    
    char* from_start_;
    char* from_end_;
    char* allocation_ptr_;  // Current allocation position in from-space
    
    char* to_start_;
    char* to_end_;
    
    // Root set for GC
    std::vector<Object**> object_roots_;
    std::vector<VMValue*> vmvalue_roots_;
    
    // GC configuration
    double gc_threshold_;
    bool auto_gc_enabled_;
    
    // Statistics
    size_t gc_count_;
    size_t total_collected_;
    
    // === INTERNAL ALLOCATION ===
    
    /**
     * Low-level allocation that handles GC triggering
     */
    void* allocate_raw(size_t bytes);
    
    /**
     * Check if GC should be triggered and do it if needed
     */
    void gc_if_needed(size_t requested_bytes);
    
    // === GARBAGE COLLECTION IMPLEMENTATION ===
    
    /**
     * Stop-and-copy garbage collection implementation
     */
    size_t perform_gc();
    
    /**
     * Copy an object from from-space to to-space
     */
    Object* copy_object(Object* obj, char*& to_alloc_ptr);
    
    /**
     * Forward an object reference (returns forwarding address if already copied)
     */
    Object* forward_object(Object* obj, char*& to_alloc_ptr);
    
    /**
     * Scan an object's references and forward them
     */
    void scan_object_references(Object* obj, char*& to_alloc_ptr);
    
    /**
     * Scan a VMValue and forward any object references
     */
    void scan_vmvalue_references(VMValue* value, char*& to_alloc_ptr);
    
    /**
     * Flip from-space and to-space after GC
     */
    void flip_spaces();
    
    // === OBJECT SCANNING HELPERS ===
    
    /**
     * Scan object based on its format
     */
    void scan_regular_object(Object* obj, char*& to_alloc_ptr);
    void scan_array_object(Object* obj, char*& to_alloc_ptr);
    void scan_byte_object(Object* obj); // Usually no references to scan
    void scan_compiled_method(Object* obj);
    
    /**
     * Determine if an object contains pointers based on its class
     */
    bool object_has_pointers(Object* obj);
    
    // === HEAP VALIDATION ===
    
    bool is_valid_object_pointer(void* ptr) const;
    bool is_in_from_space(void* ptr) const;
    bool is_in_to_space(void* ptr) const;
};

// === RAII ROOT MANAGEMENT ===

/**
 * RAII wrapper for automatically managing GC roots
 */
template<typename T>
class GCRoot {
public:
    GCRoot(SimpleMemoryManager& mm, T* ptr) : mm_(mm), ptr_(ptr) {
        if constexpr (std::is_same_v<T, Object*>) {
            mm_.add_root(reinterpret_cast<Object**>(&ptr_));
        } else if constexpr (std::is_same_v<T, VMValue>) {
            mm_.add_root(ptr_);
        }
    }
    
    ~GCRoot() {
        if constexpr (std::is_same_v<T, Object*>) {
            mm_.remove_root(reinterpret_cast<Object**>(&ptr_));
        } else if constexpr (std::is_same_v<T, VMValue>) {
            mm_.remove_root(ptr_);
        }
    }
    
    T& get() { return *ptr_; }
    const T& get() const { return *ptr_; }
    
    T* operator->() { return ptr_; }
    const T* operator->() const { return ptr_; }
    
    // Non-copyable, movable
    GCRoot(const GCRoot&) = delete;
    GCRoot& operator=(const GCRoot&) = delete;
    GCRoot(GCRoot&&) = default;
    GCRoot& operator=(GCRoot&&) = default;
    
private:
    SimpleMemoryManager& mm_;
    T* ptr_;
};

// Convenience macros for GC root management
#define GC_ROOT(mm, ptr) GCRoot _gc_root_##__LINE__(mm, &(ptr))

} // namespace smalltalk