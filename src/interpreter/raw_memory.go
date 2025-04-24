package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

// MAP_ANON is the macOS equivalent of MAP_ANONYMOUS
const MAP_ANON = 0x1000

// RawMemory manages raw memory allocation for Smalltalk objects
type RawMemory struct {
	// Memory regions
	FromSpace  []byte
	ToSpace    []byte
	FrameSpace []byte // For stack frames and other temporary allocations

	// Allocation pointers
	FromSpacePtr  uintptr
	ToSpacePtr    uintptr
	FrameSpacePtr uintptr

	// Space sizes
	SpaceSize int
	FrameSize int

	// GC threshold (percentage of space that triggers GC)
	GCThreshold float64

	// Stats
	GCCount     int
	Allocations int
}

// DefaultSpaceSize is the default size for each space (1MB)
const DefaultSpaceSize = 1024 * 1024

// NewRawMemory creates a new raw memory manager with the specified space sizes
func NewRawMemory(spaceSize, frameSize int) (*RawMemory, error) {
	if spaceSize <= 0 {
		spaceSize = DefaultSpaceSize
	}
	if frameSize <= 0 {
		frameSize = DefaultSpaceSize
	}

	// Allocate from-space
	fromSpace, err := mmapAnonymous(spaceSize)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate from-space: %w", err)
	}

	// Allocate to-space
	toSpace, err := mmapAnonymous(spaceSize)
	if err != nil {
		// Clean up from-space
		syscall.Munmap(fromSpace)
		return nil, fmt.Errorf("failed to allocate to-space: %w", err)
	}

	// Allocate frame space
	frameSpace, err := mmapAnonymous(frameSize)
	if err != nil {
		// Clean up from-space and to-space
		syscall.Munmap(fromSpace)
		syscall.Munmap(toSpace)
		return nil, fmt.Errorf("failed to allocate frame space: %w", err)
	}

	return &RawMemory{
		FromSpace:     fromSpace,
		ToSpace:       toSpace,
		FrameSpace:    frameSpace,
		FromSpacePtr:  uintptr(unsafe.Pointer(&fromSpace[0])),
		ToSpacePtr:    uintptr(unsafe.Pointer(&toSpace[0])),
		FrameSpacePtr: uintptr(unsafe.Pointer(&frameSpace[0])),
		SpaceSize:     spaceSize,
		FrameSize:     frameSize,
		GCThreshold:   0.8, // 80% threshold
		GCCount:       0,
		Allocations:   0,
	}, nil
}

// mmapAnonymous allocates memory using mmap with MAP_ANON
func mmapAnonymous(size int) ([]byte, error) {
	// Use mmap to allocate memory that won't be moved by Go's GC
	// MAP_PRIVATE | MAP_ANON ensures we get private memory not backed by a file
	mem, err := syscall.Mmap(
		-1, // File descriptor, -1 for anonymous mapping
		0,  // Offset
		size,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_PRIVATE|MAP_ANON,
	)
	if err != nil {
		return nil, err
	}
	return mem, nil
}

// Close releases all allocated memory
func (rm *RawMemory) Close() error {
	var errs []error

	if err := syscall.Munmap(rm.FromSpace); err != nil {
		errs = append(errs, fmt.Errorf("failed to unmap from-space: %w", err))
	}

	if err := syscall.Munmap(rm.ToSpace); err != nil {
		errs = append(errs, fmt.Errorf("failed to unmap to-space: %w", err))
	}

	if err := syscall.Munmap(rm.FrameSpace); err != nil {
		errs = append(errs, fmt.Errorf("failed to unmap frame space: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing raw memory: %v", errs)
	}

	return nil
}

// GetFromSpaceUsage returns the percentage of from-space that is used
func (rm *RawMemory) GetFromSpaceUsage() float64 {
	used := rm.FromSpacePtr - uintptr(unsafe.Pointer(&rm.FromSpace[0]))
	return float64(used) / float64(rm.SpaceSize)
}

// ShouldCollect returns true if garbage collection should be triggered
func (rm *RawMemory) ShouldCollect() bool {
	return rm.GetFromSpaceUsage() >= rm.GCThreshold
}

// AlignUp aligns the given address up to the specified alignment
func AlignUp(addr uintptr, alignment uintptr) uintptr {
	return (addr + alignment - 1) & ^(alignment - 1)
}

// Allocate allocates memory for an object of the given size
// Returns a pointer to the allocated memory
func (rm *RawMemory) Allocate(size int) unsafe.Pointer {
	// Check if we need to collect garbage
	if rm.ShouldCollect() {
		// We'll let the VM handle collection
		return nil
	}

	// Align the size to 8 bytes (64-bit alignment)
	alignedSize := ((size + 7) / 8) * 8

	// Align the allocation pointer
	rm.FromSpacePtr = AlignUp(rm.FromSpacePtr, 8)

	// Check if we have enough space
	if int(rm.FromSpacePtr-uintptr(unsafe.Pointer(&rm.FromSpace[0])))+alignedSize > rm.SpaceSize {
		// Not enough space
		return nil
	}

	// Get the current allocation pointer
	ptr := unsafe.Pointer(rm.FromSpacePtr)

	// Update the allocation pointer
	rm.FromSpacePtr += uintptr(alignedSize)

	// Increment allocation count
	rm.Allocations++

	return ptr
}

// AllocateObject allocates memory for an Object struct
func (rm *RawMemory) AllocateObject() *Object {
	ptr := rm.Allocate(int(unsafe.Sizeof(Object{})))
	if ptr == nil {
		return nil
	}

	// Initialize the memory to zero
	obj := (*Object)(ptr)
	*obj = Object{} // Zero initialization

	return obj
}

// AllocateString allocates memory for a String struct
func (rm *RawMemory) AllocateString() *String {
	ptr := rm.Allocate(int(unsafe.Sizeof(String{})))
	if ptr == nil {
		return nil
	}

	// Initialize the memory to zero
	str := (*String)(ptr)
	*str = String{} // Zero initialization

	return str
}

// AllocateSymbol allocates memory for a Symbol struct
func (rm *RawMemory) AllocateSymbol() *Symbol {
	ptr := rm.Allocate(int(unsafe.Sizeof(Symbol{})))
	if ptr == nil {
		return nil
	}

	// Initialize the memory to zero
	sym := (*Symbol)(ptr)
	*sym = Symbol{} // Zero initialization

	return sym
}

// AllocateClass allocates memory for a Class struct
func (rm *RawMemory) AllocateClass() *Class {
	ptr := rm.Allocate(int(unsafe.Sizeof(Class{})))
	if ptr == nil {
		return nil
	}

	// Initialize the memory to zero
	class := (*Class)(ptr)
	*class = Class{} // Zero initialization

	return class
}

// AllocateMethod allocates memory for a Method struct
func (rm *RawMemory) AllocateMethod() *Method {
	ptr := rm.Allocate(int(unsafe.Sizeof(Method{})))
	if ptr == nil {
		return nil
	}

	// Initialize the memory to zero
	method := (*Method)(ptr)
	*method = Method{} // Zero initialization

	return method
}

// AllocateObjectArray allocates memory for an array of Object pointers
func (rm *RawMemory) AllocateObjectArray(size int) []*Object {
	if size <= 0 {
		return nil
	}

	// Calculate the size of the array
	arraySize := size * int(unsafe.Sizeof((*Object)(nil)))

	// Allocate memory for the array
	ptr := rm.Allocate(arraySize)
	if ptr == nil {
		return nil
	}

	// Create a slice header pointing to the allocated memory
	slice := &sliceHeader{
		Data: uintptr(ptr),
		Len:  size,
		Cap:  size,
	}

	// Convert the slice header to a slice
	array := *(*[]*Object)(unsafe.Pointer(slice))

	// Initialize the array to nil
	for i := range array {
		array[i] = nil
	}

	return array
}

// sliceHeader represents the runtime structure of a slice
type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

// AllocateByteArray allocates memory for a byte array
func (rm *RawMemory) AllocateByteArray(size int) []byte {
	if size <= 0 {
		return nil
	}

	// Allocate memory for the array
	ptr := rm.Allocate(size)
	if ptr == nil {
		return nil
	}

	// Create a slice header pointing to the allocated memory
	slice := &sliceHeader{
		Data: uintptr(ptr),
		Len:  size,
		Cap:  size,
	}

	// Convert the slice header to a slice
	array := *(*[]byte)(unsafe.Pointer(slice))

	// Initialize the array to zero
	for i := range array {
		array[i] = 0
	}

	return array
}

// AllocateStringMap allocates memory for a string map
func (rm *RawMemory) AllocateStringMap() map[string]*Object {
	// For maps, we can't directly allocate memory
	// We'll use Go's map allocation for now
	// In a more advanced implementation, we could create our own map implementation
	return make(map[string]*Object)
}

// IsPointerInFromSpace checks if a pointer is in the from-space
func (rm *RawMemory) IsPointerInFromSpace(ptr unsafe.Pointer) bool {
	addr := uintptr(ptr)
	fromSpaceStart := uintptr(unsafe.Pointer(&rm.FromSpace[0]))
	fromSpaceEnd := fromSpaceStart + uintptr(rm.SpaceSize)

	return addr >= fromSpaceStart && addr < fromSpaceEnd
}

// IsPointerInToSpace checks if a pointer is in the to-space
func (rm *RawMemory) IsPointerInToSpace(ptr unsafe.Pointer) bool {
	addr := uintptr(ptr)
	toSpaceStart := uintptr(unsafe.Pointer(&rm.ToSpace[0]))
	toSpaceEnd := toSpaceStart + uintptr(rm.SpaceSize)

	return addr >= toSpaceStart && addr < toSpaceEnd
}

// SwapSpaces swaps the from-space and to-space
func (rm *RawMemory) SwapSpaces() {
	rm.FromSpace, rm.ToSpace = rm.ToSpace, rm.FromSpace
	rm.FromSpacePtr = uintptr(unsafe.Pointer(&rm.FromSpace[0]))
	rm.ToSpacePtr = uintptr(unsafe.Pointer(&rm.ToSpace[0]))
}

// CopyObject copies an object from from-space to to-space
// Returns a pointer to the copied object in to-space
func (rm *RawMemory) CopyObject(obj *Object) *Object {
	if obj == nil {
		return nil
	}

	// Check if it's an immediate value
	if IsImmediate(obj) {
		// Immediate values don't need to be copied
		return obj
	}

	// Check if the object has already been moved
	if obj.Moved && obj.ForwardingPtr != nil {
		return obj.ForwardingPtr
	}

	// Determine the size of the object based on its type
	var size int
	switch obj.Type {
	case OBJ_STRING:
		size = int(unsafe.Sizeof(String{}))
	case OBJ_SYMBOL:
		size = int(unsafe.Sizeof(Symbol{}))
	case OBJ_CLASS:
		size = int(unsafe.Sizeof(Class{}))
	case OBJ_METHOD:
		size = int(unsafe.Sizeof(Method{}))
	default:
		size = int(unsafe.Sizeof(Object{}))
	}

	// Align the to-space pointer
	rm.ToSpacePtr = AlignUp(rm.ToSpacePtr, 8)

	// Copy the object to to-space
	dst := unsafe.Pointer(rm.ToSpacePtr)
	src := unsafe.Pointer(obj)

	// Copy the memory
	for i := 0; i < size; i++ {
		*(*byte)(unsafe.Pointer(uintptr(dst) + uintptr(i))) = *(*byte)(unsafe.Pointer(uintptr(src) + uintptr(i)))
	}

	// Update the to-space pointer
	rm.ToSpacePtr += uintptr(size)

	// Set the forwarding pointer
	obj.Moved = true
	obj.ForwardingPtr = (*Object)(dst)

	return (*Object)(dst)
}
