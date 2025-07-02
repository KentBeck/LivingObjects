#include "simple_memory_manager.h"
#include <algorithm>
#include <cstring>
#include <stdexcept>
#include <iostream>
#include <fstream>
#include <cassert>

namespace smalltalk {

SimpleMemoryManager::SimpleMemoryManager(size_t heap_size)
    : heap_size_(heap_size),
      from_space_(std::make_unique<char[]>(heap_size)),
      to_space_(std::make_unique<char[]>(heap_size)),
      gc_threshold_(0.8),
      auto_gc_enabled_(true),
      gc_count_(0),
      total_collected_(0) {
    
    // Initialize space pointers
    from_start_ = from_space_.get();
    from_end_ = from_start_ + heap_size_;
    allocation_ptr_ = from_start_;
    
    to_start_ = to_space_.get();
    to_end_ = to_start_ + heap_size_;
    
    // Clear the heap
    std::memset(from_start_, 0, heap_size_);
    std::memset(to_start_, 0, heap_size_);
}

SimpleMemoryManager::~SimpleMemoryManager() = default;

// === CORE ALLOCATION ===

Object* SimpleMemoryManager::allocate_object(SmalltalkClass* smalltalk_class, uint32_t data_size, bool is_byte_object) {
    // Calculate total size needed
    size_t total_size = object_size_bytes(data_size, is_byte_object);
    
    // Align to ALIGNMENT_BYTES boundary
    total_size = (total_size + ALIGNMENT_BYTES - 1) & ~(ALIGNMENT_BYTES - 1);
    
    // Check if GC is needed
    if (auto_gc_enabled_) {
        gc_if_needed(total_size);
    }
    
    // Allocate raw memory
    void* raw_ptr = allocate_raw(total_size);
    if (!raw_ptr) {
        throw std::runtime_error("Out of memory");
    }
    
    // Construct Object in-place
    Object* obj = new (raw_ptr) Object(smalltalk_class, data_size);
    
    // Initialize data area to zero
    std::memset(obj->data(), 0, is_byte_object ? data_size : data_size * sizeof(void*));
    
    // Set flags based on object type
    if (is_byte_object) {
        obj->header.clear_flag(ObjectFlag::HAS_POINTERS);
    } else {
        obj->header.set_flag(ObjectFlag::HAS_POINTERS);
    }
    
    return obj;
}

Object* SimpleMemoryManager::allocate_regular_object(SmalltalkClass* smalltalk_class, uint32_t instance_var_count) {
    return allocate_object(smalltalk_class, instance_var_count, false);
}

Object* SimpleMemoryManager::allocate_array(SmalltalkClass* array_class, uint32_t element_count) {
    return allocate_object(array_class, element_count, false);
}

Object* SimpleMemoryManager::allocate_byte_array(SmalltalkClass* byte_array_class, uint32_t byte_count) {
    return allocate_object(byte_array_class, byte_count, true);
}

Object* SimpleMemoryManager::allocate_string(SmalltalkClass* string_class, uint32_t byte_count) {
    return allocate_object(string_class, byte_count, true);
}

// === OBJECT BUILDER IMPLEMENTATION ===

SimpleMemoryManager::ObjectBuilder::ObjectBuilder(SimpleMemoryManager& manager, SmalltalkClass* st_class)
    : manager_(manager), class_(st_class), size_(0), is_byte_object_(false), 
      identity_hash_(0), immutable_(false), pinned_(false) {}

SimpleMemoryManager::ObjectBuilder& SimpleMemoryManager::ObjectBuilder::with_instance_variables(uint32_t count) {
    size_ = count;
    is_byte_object_ = false;
    return *this;
}

SimpleMemoryManager::ObjectBuilder& SimpleMemoryManager::ObjectBuilder::with_array_elements(uint32_t count) {
    size_ = count;
    is_byte_object_ = false;
    return *this;
}

SimpleMemoryManager::ObjectBuilder& SimpleMemoryManager::ObjectBuilder::with_byte_data(uint32_t count) {
    size_ = count;
    is_byte_object_ = true;
    return *this;
}

SimpleMemoryManager::ObjectBuilder& SimpleMemoryManager::ObjectBuilder::with_identity_hash(uint16_t hash) {
    identity_hash_ = hash;
    return *this;
}

SimpleMemoryManager::ObjectBuilder& SimpleMemoryManager::ObjectBuilder::immutable() {
    immutable_ = true;
    return *this;
}

SimpleMemoryManager::ObjectBuilder& SimpleMemoryManager::ObjectBuilder::pinned() {
    pinned_ = true;
    return *this;
}

Object* SimpleMemoryManager::ObjectBuilder::build() {
    Object* obj = manager_.allocate_object(class_, size_, is_byte_object_);
    
    if (identity_hash_ != 0) {
        obj->set_identity_hash(identity_hash_);
    }
    
    if (immutable_) {
        obj->header.set_flag(ObjectFlag::IMMUTABLE);
    }
    
    if (pinned_) {
        obj->header.set_flag(ObjectFlag::PINNED);
    }
    
    return obj;
}

// === GARBAGE COLLECTION ===

size_t SimpleMemoryManager::collect_garbage() {
    size_t bytes_freed = perform_gc();
    
    gc_count_++;
    total_collected_ += bytes_freed;
    
    return bytes_freed;
}

void SimpleMemoryManager::add_root(Object** root_ptr) {
    object_roots_.push_back(root_ptr);
}

void SimpleMemoryManager::remove_root(Object** root_ptr) {
    auto it = std::find(object_roots_.begin(), object_roots_.end(), root_ptr);
    if (it != object_roots_.end()) {
        object_roots_.erase(it);
    }
}

void SimpleMemoryManager::add_root(VMValue* root_value) {
    vmvalue_roots_.push_back(root_value);
}

void SimpleMemoryManager::remove_root(VMValue* root_value) {
    auto it = std::find(vmvalue_roots_.begin(), vmvalue_roots_.end(), root_value);
    if (it != vmvalue_roots_.end()) {
        vmvalue_roots_.erase(it);
    }
}

// === MEMORY STATISTICS ===

size_t SimpleMemoryManager::bytes_allocated() const {
    return static_cast<size_t>(allocation_ptr_ - from_start_);
}

size_t SimpleMemoryManager::bytes_free() const {
    return heap_size_ - bytes_allocated();
}

double SimpleMemoryManager::heap_utilization() const {
    return static_cast<double>(bytes_allocated()) / static_cast<double>(heap_size_);
}

void SimpleMemoryManager::set_gc_threshold(double threshold) {
    if (threshold < 0.0 || threshold > 1.0) {
        throw std::invalid_argument("GC threshold must be between 0.0 and 1.0");
    }
    gc_threshold_ = threshold;
}

// === DEBUGGING SUPPORT ===

SimpleMemoryManager::HeapStats SimpleMemoryManager::get_heap_stats() const {
    HeapStats stats = {0, 0, 0, 0, 0, 0};
    
    // Walk through all objects in from-space
    char* current = from_start_;
    while (current < allocation_ptr_) {
        Object* obj = reinterpret_cast<Object*>(current);
        
        stats.total_objects++;
        
        // Classify object type based on format
        ObjectFormat format = get_object_format(obj);
        switch (format) {
            case ObjectFormat::REGULAR:
                stats.regular_objects++;
                break;
            case ObjectFormat::ARRAY:
                stats.array_objects++;
                break;
            case ObjectFormat::BYTE_ARRAY:
                stats.byte_objects++;
                break;
            case ObjectFormat::COMPILED_METHOD:
                stats.byte_objects++; // Methods are byte objects
                break;
        }
        
        // Calculate object size
        size_t obj_size = object_size_bytes(obj->size(), format == ObjectFormat::BYTE_ARRAY);
        obj_size = (obj_size + ALIGNMENT_BYTES - 1) & ~(ALIGNMENT_BYTES - 1);
        stats.total_bytes += obj_size;
        
        current += obj_size;
    }
    
    // Calculate fragmentation (internal fragmentation due to alignment)
    stats.fragmentation_bytes = bytes_allocated() - stats.total_bytes;
    
    return stats;
}

void SimpleMemoryManager::dump_heap(const char* filename) const {
    std::ofstream file(filename);
    if (!file.is_open()) {
        throw std::runtime_error("Could not open file for heap dump");
    }
    
    file << "=== Heap Dump ===" << std::endl;
    file << "Heap size: " << heap_size_ << " bytes" << std::endl;
    file << "Bytes allocated: " << bytes_allocated() << " bytes" << std::endl;
    file << "Heap utilization: " << (heap_utilization() * 100.0) << "%" << std::endl;
    file << "GC count: " << gc_count_ << std::endl;
    file << std::endl;
    
    HeapStats stats = get_heap_stats();
    file << "Objects: " << stats.total_objects << std::endl;
    file << "  Regular: " << stats.regular_objects << std::endl;
    file << "  Arrays: " << stats.array_objects << std::endl;
    file << "  Byte objects: " << stats.byte_objects << std::endl;
    file << "Total object bytes: " << stats.total_bytes << std::endl;
    file << "Fragmentation: " << stats.fragmentation_bytes << " bytes" << std::endl;
    file << std::endl;
    
    // Dump individual objects
    char* current = from_start_;
    size_t object_index = 0;
    while (current < allocation_ptr_) {
        Object* obj = reinterpret_cast<Object*>(current);
        
        file << "Object " << object_index << " @ " << static_cast<void*>(obj) << std::endl;
        file << "  Class: " << (obj->get_class() ? obj->get_class()->name() : "NULL") << std::endl;
        file << "  Size: " << obj->size() << std::endl;
        file << "  Hash: " << obj->identity_hash() << std::endl;
        file << "  Flags: ";
        if (obj->header.has_flag(ObjectFlag::MARKED)) file << "MARKED ";
        if (obj->header.has_flag(ObjectFlag::IMMUTABLE)) file << "IMMUTABLE ";
        if (obj->header.has_flag(ObjectFlag::PINNED)) file << "PINNED ";
        if (obj->header.has_flag(ObjectFlag::HAS_POINTERS)) file << "HAS_POINTERS ";
        file << std::endl;
        
        size_t obj_size = object_size_bytes(obj->size(), !obj->header.has_flag(ObjectFlag::HAS_POINTERS));
        obj_size = (obj_size + ALIGNMENT_BYTES - 1) & ~(ALIGNMENT_BYTES - 1);
        current += obj_size;
        object_index++;
    }
}

bool SimpleMemoryManager::validate_heap() const {
    char* current = from_start_;
    while (current < allocation_ptr_) {
        Object* obj = reinterpret_cast<Object*>(current);
        
        // Check object pointer is aligned
        if (reinterpret_cast<uintptr_t>(obj) % ALIGNMENT_BYTES != 0) {
            std::cerr << "Heap validation failed: Object at " << static_cast<void*>(obj) 
                      << " is not properly aligned" << std::endl;
            return false;
        }
        
        // Check object has valid class
        if (!obj->get_class()) {
            std::cerr << "Heap validation failed: Object at " << static_cast<void*>(obj) 
                      << " has NULL class" << std::endl;
            return false;
        }
        
        // Check object size is reasonable
        if (obj->size() > heap_size_) {
            std::cerr << "Heap validation failed: Object at " << static_cast<void*>(obj) 
                      << " has unreasonable size " << obj->size() << std::endl;
            return false;
        }
        
        // Move to next object
        size_t obj_size = object_size_bytes(obj->size(), !obj->header.has_flag(ObjectFlag::HAS_POINTERS));
        obj_size = (obj_size + ALIGNMENT_BYTES - 1) & ~(ALIGNMENT_BYTES - 1);
        current += obj_size;
    }
    
    return true;
}

// === PRIVATE IMPLEMENTATION ===

void* SimpleMemoryManager::allocate_raw(size_t bytes) {
    // Check if we have enough space
    if (allocation_ptr_ + bytes > from_end_) {
        return nullptr; // Out of memory
    }
    
    void* result = allocation_ptr_;
    allocation_ptr_ += bytes;
    
    return result;
}

void SimpleMemoryManager::gc_if_needed(size_t requested_bytes) {
    // Check if allocation would exceed threshold
    size_t bytes_after_alloc = bytes_allocated() + requested_bytes;
    double utilization_after = static_cast<double>(bytes_after_alloc) / static_cast<double>(heap_size_);
    
    if (utilization_after > gc_threshold_) {
        collect_garbage();
    }
}

size_t SimpleMemoryManager::perform_gc() {
    size_t bytes_before = bytes_allocated();
    
    // Clear to-space
    std::memset(to_start_, 0, heap_size_);
    
    // Reset to-space allocation pointer
    char* to_allocation_ptr = to_start_;
    
    // Forward all root objects
    for (Object** root : object_roots_) {
        if (*root) {
            *root = forward_object(*root, to_allocation_ptr);
        }
    }
    
    // Forward all VMValue roots
    for (VMValue* root : vmvalue_roots_) {
        scan_vmvalue_references(root, to_allocation_ptr);
    }
    
    // Scan forwarded objects for references
    char* scan_ptr = to_start_;
    while (scan_ptr < to_allocation_ptr) {
        Object* obj = reinterpret_cast<Object*>(scan_ptr);
        scan_object_references(obj, to_allocation_ptr);
        
        // Move to next object
        size_t obj_size = object_size_bytes(obj->size(), !obj->header.has_flag(ObjectFlag::HAS_POINTERS));
        obj_size = (obj_size + ALIGNMENT_BYTES - 1) & ~(ALIGNMENT_BYTES - 1);
        scan_ptr += obj_size;
    }
    
    // Update allocation pointer in to-space
    to_allocation_ptr = scan_ptr;
    
    // Flip spaces
    flip_spaces();
    
    // Update allocation pointer
    allocation_ptr_ = to_allocation_ptr;
    
    size_t bytes_after = bytes_allocated();
    return bytes_before - bytes_after;
}

Object* SimpleMemoryManager::copy_object(Object* obj, char*& to_alloc_ptr) {
    if (!obj) return nullptr;
    
    // Calculate object size
    bool is_byte_obj = !obj->header.has_flag(ObjectFlag::HAS_POINTERS);
    size_t obj_size = object_size_bytes(obj->size(), is_byte_obj);
    obj_size = (obj_size + ALIGNMENT_BYTES - 1) & ~(ALIGNMENT_BYTES - 1);
    
    // Allocate space in to-space
    char* new_location = to_alloc_ptr;
    to_alloc_ptr += obj_size;
    
    if (to_alloc_ptr > to_end_) {
        throw std::runtime_error("Out of memory during GC");
    }
    
    // Copy the object
    std::memcpy(new_location, obj, obj_size);
    
    // Clear GC flags in the new copy
    Object* new_obj = reinterpret_cast<Object*>(new_location);
    new_obj->header.clear_flag(ObjectFlag::MARKED);
    new_obj->header.clear_flag(ObjectFlag::FORWARDED);
    
    return new_obj;
}

Object* SimpleMemoryManager::forward_object(Object* obj, char*& to_alloc_ptr) {
    if (!obj) return nullptr;
    
    // Check if object is already forwarded
    if (obj->is_forwarded()) {
        return obj->forwarding_address();
    }
    
    // Copy object to to-space
    Object* new_obj = copy_object(obj, to_alloc_ptr);
    
    // Set up forwarding pointer
    obj->set_forwarding_address(new_obj);
    
    return new_obj;
}

void SimpleMemoryManager::scan_object_references(Object* obj, char*& to_alloc_ptr) {
    if (!obj || !obj->has_pointers()) {
        return;
    }
    
    ObjectFormat format = get_object_format(obj);
    
    switch (format) {
        case ObjectFormat::REGULAR:
            scan_regular_object(obj, to_alloc_ptr);
            break;
        case ObjectFormat::ARRAY:
            scan_array_object(obj, to_alloc_ptr);
            break;
        case ObjectFormat::BYTE_ARRAY:
            scan_byte_object(obj);
            break;
        case ObjectFormat::COMPILED_METHOD:
            scan_compiled_method(obj);
            break;
    }
}

void SimpleMemoryManager::scan_vmvalue_references(VMValue* value, char*& to_alloc_ptr) {
    if (value->is_heap_object()) {
        Object* obj = value->as_object();
        if (obj) {
            // Forward the object and update the VMValue
            Object* forwarded = forward_object(obj, to_alloc_ptr);
            *value = VMValue(forwarded);
        }
    }
    // Immediates don't need forwarding
}

void SimpleMemoryManager::flip_spaces() {
    // Swap the spaces
    std::swap(from_space_, to_space_);
    std::swap(from_start_, to_start_);
    std::swap(from_end_, to_end_);
    
    // Reset allocation pointer to start of new from-space
    allocation_ptr_ = from_start_;
}

void SimpleMemoryManager::scan_regular_object(Object* obj, char*& to_alloc_ptr) {
    // Scan all instance variable slots
    for (uint32_t i = 0; i < obj->size(); i++) {
        Object** slot_ptr = reinterpret_cast<Object**>(obj->slots()) + i;
        if (*slot_ptr) {
            *slot_ptr = forward_object(*slot_ptr, to_alloc_ptr);
        }
    }
}

void SimpleMemoryManager::scan_array_object(Object* obj, char*& to_alloc_ptr) {
    // Same as regular object - scan all element slots
    scan_regular_object(obj, to_alloc_ptr);
}

void SimpleMemoryManager::scan_byte_object(Object* obj) {
    // Byte objects typically don't contain pointers
    // But we should check the class to be sure
    (void)obj; // Suppress unused parameter warning
}

void SimpleMemoryManager::scan_compiled_method(Object* obj) {
    // CompiledMethod has a complex structure with both bytecode and literals
    // For now, treat it as having no pointers (literals would be in a separate area)
    (void)obj; // Suppress unused parameter warning
}

bool SimpleMemoryManager::object_has_pointers(Object* obj) {
    return obj->header.has_flag(ObjectFlag::HAS_POINTERS);
}

bool SimpleMemoryManager::is_valid_object_pointer(void* ptr) const {
    return is_in_from_space(ptr) || is_in_to_space(ptr);
}

bool SimpleMemoryManager::is_in_from_space(void* ptr) const {
    return ptr >= from_start_ && ptr < from_end_;
}

bool SimpleMemoryManager::is_in_to_space(void* ptr) const {
    return ptr >= to_start_ && ptr < to_end_;
}

} // namespace smalltalk