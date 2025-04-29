package main

import (
	"testing"
	"unsafe"
)

// TestRawMemoryAllocation tests the allocation of memory in the raw memory system
func TestRawMemoryAllocation(t *testing.T) {
	// Create a new raw memory
	rawMem, err := NewRawMemory(1024*1024, 1024*1024)
	if err != nil {
		t.Fatalf("Failed to create raw memory: %v", err)
	}
	defer rawMem.Close()

	// Allocate some memory
	ptr1 := rawMem.Allocate(100)
	if ptr1 == nil {
		t.Fatalf("Failed to allocate memory")
	}

	// Allocate more memory
	ptr2 := rawMem.Allocate(200)
	if ptr2 == nil {
		t.Fatalf("Failed to allocate memory")
	}

	// Check that the pointers are different
	if ptr1 == ptr2 {
		t.Errorf("Expected different pointers, got the same pointer")
	}

	// Check that the pointers are in the from-space
	if !rawMem.IsPointerInFromSpace(ptr1) {
		t.Errorf("Expected ptr1 to be in from-space")
	}
	if !rawMem.IsPointerInFromSpace(ptr2) {
		t.Errorf("Expected ptr2 to be in from-space")
	}

	// Check that the pointers are aligned
	if uintptr(ptr1)%8 != 0 {
		t.Errorf("Expected ptr1 to be 8-byte aligned, got %v", ptr1)
	}
	if uintptr(ptr2)%8 != 0 {
		t.Errorf("Expected ptr2 to be 8-byte aligned, got %v", ptr2)
	}

	// Check that the second pointer is after the first
	if uintptr(ptr2) <= uintptr(ptr1) {
		t.Errorf("Expected ptr2 to be after ptr1")
	}

	// Check that the distance between the pointers is at least 100 bytes (aligned to 8)
	distance := uintptr(ptr2) - uintptr(ptr1)
	if distance < 104 { // 100 rounded up to the next multiple of 8 (104)
		t.Errorf("Expected distance between pointers to be at least 104 bytes, got %d", distance)
	}
}

// TestRawMemoryObjectAllocation tests the allocation of objects in the raw memory system
func TestRawMemoryObjectAllocation(t *testing.T) {
	// Create a new raw memory
	rawMem, err := NewRawMemory(1024*1024, 1024*1024)
	if err != nil {
		t.Fatalf("Failed to create raw memory: %v", err)
	}
	defer rawMem.Close()

	// Allocate an object
	obj := rawMem.AllocateObject()
	if obj == nil {
		t.Fatalf("Failed to allocate object")
	}

	// Check that the object is in the from-space
	if !rawMem.IsPointerInFromSpace(unsafe.Pointer(obj)) {
		t.Errorf("Expected object to be in from-space")
	}

	// Check that the object is properly initialized
	if obj.Type() != 0 {
		t.Errorf("Expected object type to be 0, got %d", obj.Type())
	}
	if obj.Class != nil {
		t.Errorf("Expected object class to be nil, got %v", obj.Class)
	}
	if obj.Moved != false {
		t.Errorf("Expected object moved to be false, got %v", obj.Moved)
	}
	if obj.ForwardingPtr != nil {
		t.Errorf("Expected object forwarding pointer to be nil, got %v", obj.ForwardingPtr)
	}
}

// TestRawMemoryStringAllocation tests the allocation of strings in the raw memory system
func TestRawMemoryStringAllocation(t *testing.T) {
	// Create a new raw memory
	rawMem, err := NewRawMemory(1024*1024, 1024*1024)
	if err != nil {
		t.Fatalf("Failed to create raw memory: %v", err)
	}
	defer rawMem.Close()

	// Allocate a string
	str := rawMem.AllocateString()
	if str == nil {
		t.Fatalf("Failed to allocate string")
	}

	// Check that the string is in the from-space
	if !rawMem.IsPointerInFromSpace(unsafe.Pointer(str)) {
		t.Errorf("Expected string to be in from-space")
	}

	// Check that the string is properly initialized
	if str.Type() != 0 {
		t.Errorf("Expected string type to be 0, got %d", str.Type())
	}
	if str.Value != "" {
		t.Errorf("Expected string value to be empty, got %s", str.Value)
	}
}

// TestRawMemorySymbolAllocation tests the allocation of symbols in the raw memory system
func TestRawMemorySymbolAllocation(t *testing.T) {
	// Create a new raw memory
	rawMem, err := NewRawMemory(1024*1024, 1024*1024)
	if err != nil {
		t.Fatalf("Failed to create raw memory: %v", err)
	}
	defer rawMem.Close()

	// Allocate a symbol
	sym := rawMem.AllocateSymbol()
	if sym == nil {
		t.Fatalf("Failed to allocate symbol")
	}

	// Check that the symbol is in the from-space
	if !rawMem.IsPointerInFromSpace(unsafe.Pointer(sym)) {
		t.Errorf("Expected symbol to be in from-space")
	}

	// Check that the symbol is properly initialized
	if sym.Type() != 0 {
		t.Errorf("Expected symbol type to be 0, got %d", sym.Type())
	}
	if sym.Value != "" {
		t.Errorf("Expected symbol value to be empty, got %s", sym.Value)
	}
}

// TestRawMemoryArrayAllocation tests the allocation of arrays in the raw memory system
func TestRawMemoryArrayAllocation(t *testing.T) {
	// Create a new raw memory
	rawMem, err := NewRawMemory(1024*1024, 1024*1024)
	if err != nil {
		t.Fatalf("Failed to create raw memory: %v", err)
	}
	defer rawMem.Close()

	// Allocate an array
	array := rawMem.AllocateObjectArray(10)
	if array == nil {
		t.Fatalf("Failed to allocate array")
	}

	// Check that the array has the correct length
	if len(array) != 10 {
		t.Errorf("Expected array length to be 10, got %d", len(array))
	}

	// Check that the array elements are initialized to nil
	for i, elem := range array {
		if elem != nil {
			t.Errorf("Expected array element %d to be nil, got %v", i, elem)
		}
	}
}

// TestRawMemoryManagerAllocation tests the allocation of objects using the raw memory manager
func TestRawMemoryManagerAllocation(t *testing.T) {
	// Create a new raw memory manager
	memoryManager, err := NewRawMemoryManager()
	if err != nil {
		t.Fatalf("Failed to create raw memory manager: %v", err)
	}
	defer memoryManager.Close()

	// Allocate an object
	obj := memoryManager.AllocateObject()
	if obj == nil {
		t.Fatalf("Failed to allocate object")
	}

	// Check that the object is properly initialized
	if obj.Type() != 0 {
		t.Errorf("Expected object type to be 0, got %d", obj.Type())
	}
	if obj.Class != nil {
		t.Errorf("Expected object class to be nil, got %v", obj.Class)
	}
	if obj.Moved != false {
		t.Errorf("Expected object moved to be false, got %v", obj.Moved)
	}
	if obj.ForwardingPtr != nil {
		t.Errorf("Expected object forwarding pointer to be nil, got %v", obj.ForwardingPtr)
	}
}

// TestVMWithRawMemory tests the creation of a VM with raw memory
func TestVMWithRawMemory(t *testing.T) {
	// Create a new VM with raw memory
	vm, err := NewVMWithRawMemory()
	if err != nil {
		t.Fatalf("Failed to create VM with raw memory: %v", err)
	}
	defer vm.Close()

	// Check that the VM is properly initialized
	if vm.ObjectClass == nil {
		t.Errorf("Expected ObjectClass to be non-nil")
	}
	if vm.NilClass == nil {
		t.Errorf("Expected NilClass to be non-nil")
	}
	if vm.TrueClass == nil {
		t.Errorf("Expected TrueClass to be non-nil")
	}
	if vm.FalseClass == nil {
		t.Errorf("Expected FalseClass to be non-nil")
	}
	if vm.IntegerClass == nil {
		t.Errorf("Expected IntegerClass to be non-nil")
	}
	if vm.FloatClass == nil {
		t.Errorf("Expected FloatClass to be non-nil")
	}

	// Check that the special objects are properly initialized
	if !IsNilImmediate(vm.NilObject) {
		t.Errorf("Expected NilObject to be an immediate nil value")
	}
	if !IsTrueImmediate(vm.TrueObject) {
		t.Errorf("Expected TrueObject to be an immediate true value")
	}
	if !IsFalseImmediate(vm.FalseObject) {
		t.Errorf("Expected FalseObject to be an immediate false value")
	}
}
