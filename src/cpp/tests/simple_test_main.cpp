#include "../include/bytecode.h"
#include "../include/object.h"
#include "../include/memory_manager.h"
#include <iostream>
#include <cassert>
#include <cstring>

using namespace smalltalk;

// Simple test framework
#define TEST(name) void name()
#define EXPECT_EQ(expected, actual) assert((expected) == (actual))
#define EXPECT_NE(expected, actual) assert((expected) != (actual))
#define EXPECT_STREQ(expected, actual) assert(strcmp((expected), (actual)) == 0)
#define EXPECT_LT(a, b) assert((a) < (b))
#define EXPECT_GT(a, b) assert((a) > (b))

TEST(TestBytecodeInstructionSizes) {
    // Test instruction sizes match the Go implementation
    EXPECT_EQ(5, getInstructionSize(Bytecode::PUSH_LITERAL));
    EXPECT_EQ(5, getInstructionSize(Bytecode::PUSH_INSTANCE_VARIABLE));
    EXPECT_EQ(5, getInstructionSize(Bytecode::PUSH_TEMPORARY_VARIABLE));
    EXPECT_EQ(1, getInstructionSize(Bytecode::PUSH_SELF));
    EXPECT_EQ(5, getInstructionSize(Bytecode::STORE_INSTANCE_VARIABLE));
    EXPECT_EQ(5, getInstructionSize(Bytecode::STORE_TEMPORARY_VARIABLE));
    EXPECT_EQ(9, getInstructionSize(Bytecode::SEND_MESSAGE));
    EXPECT_EQ(1, getInstructionSize(Bytecode::RETURN_STACK_TOP));
    EXPECT_EQ(5, getInstructionSize(Bytecode::JUMP));
    EXPECT_EQ(5, getInstructionSize(Bytecode::JUMP_IF_TRUE));
    EXPECT_EQ(5, getInstructionSize(Bytecode::JUMP_IF_FALSE));
    EXPECT_EQ(1, getInstructionSize(Bytecode::POP));
    EXPECT_EQ(1, getInstructionSize(Bytecode::DUPLICATE));
    EXPECT_EQ(13, getInstructionSize(Bytecode::CREATE_BLOCK));
    EXPECT_EQ(5, getInstructionSize(Bytecode::EXECUTE_BLOCK));
}

TEST(TestBytecodeNames) {
    // Test bytecode names
    EXPECT_STREQ("PUSH_LITERAL", getBytecodeString(Bytecode::PUSH_LITERAL));
    EXPECT_STREQ("PUSH_INSTANCE_VARIABLE", getBytecodeString(Bytecode::PUSH_INSTANCE_VARIABLE));
    EXPECT_STREQ("PUSH_TEMPORARY_VARIABLE", getBytecodeString(Bytecode::PUSH_TEMPORARY_VARIABLE));
    EXPECT_STREQ("PUSH_SELF", getBytecodeString(Bytecode::PUSH_SELF));
    EXPECT_STREQ("STORE_INSTANCE_VARIABLE", getBytecodeString(Bytecode::STORE_INSTANCE_VARIABLE));
    EXPECT_STREQ("STORE_TEMPORARY_VARIABLE", getBytecodeString(Bytecode::STORE_TEMPORARY_VARIABLE));
    EXPECT_STREQ("SEND_MESSAGE", getBytecodeString(Bytecode::SEND_MESSAGE));
    EXPECT_STREQ("RETURN_STACK_TOP", getBytecodeString(Bytecode::RETURN_STACK_TOP));
    EXPECT_STREQ("JUMP", getBytecodeString(Bytecode::JUMP));
    EXPECT_STREQ("JUMP_IF_TRUE", getBytecodeString(Bytecode::JUMP_IF_TRUE));
    EXPECT_STREQ("JUMP_IF_FALSE", getBytecodeString(Bytecode::JUMP_IF_FALSE));
    EXPECT_STREQ("POP", getBytecodeString(Bytecode::POP));
    EXPECT_STREQ("DUPLICATE", getBytecodeString(Bytecode::DUPLICATE));
    EXPECT_STREQ("CREATE_BLOCK", getBytecodeString(Bytecode::CREATE_BLOCK));
    EXPECT_STREQ("EXECUTE_BLOCK", getBytecodeString(Bytecode::EXECUTE_BLOCK));
}

TEST(TestBytecodeValues) {
    // Test bytecode values match the Go implementation
    EXPECT_EQ(0, static_cast<uint8_t>(Bytecode::PUSH_LITERAL));
    EXPECT_EQ(1, static_cast<uint8_t>(Bytecode::PUSH_INSTANCE_VARIABLE));
    EXPECT_EQ(2, static_cast<uint8_t>(Bytecode::PUSH_TEMPORARY_VARIABLE));
    EXPECT_EQ(3, static_cast<uint8_t>(Bytecode::PUSH_SELF));
    EXPECT_EQ(4, static_cast<uint8_t>(Bytecode::STORE_INSTANCE_VARIABLE));
    EXPECT_EQ(5, static_cast<uint8_t>(Bytecode::STORE_TEMPORARY_VARIABLE));
    EXPECT_EQ(6, static_cast<uint8_t>(Bytecode::SEND_MESSAGE));
    EXPECT_EQ(7, static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP));
    EXPECT_EQ(8, static_cast<uint8_t>(Bytecode::JUMP));
    EXPECT_EQ(9, static_cast<uint8_t>(Bytecode::JUMP_IF_TRUE));
    EXPECT_EQ(10, static_cast<uint8_t>(Bytecode::JUMP_IF_FALSE));
    EXPECT_EQ(11, static_cast<uint8_t>(Bytecode::POP));
    EXPECT_EQ(12, static_cast<uint8_t>(Bytecode::DUPLICATE));
    EXPECT_EQ(13, static_cast<uint8_t>(Bytecode::CREATE_BLOCK));
    EXPECT_EQ(14, static_cast<uint8_t>(Bytecode::EXECUTE_BLOCK));
}

TEST(TestMemoryObjectAllocation) {
    MemoryManager memory;
    
    // Test allocating a basic object
    Object* obj = memory.allocateObject(ObjectType::OBJECT, 10);
    EXPECT_NE(nullptr, obj);
    EXPECT_EQ(static_cast<uint64_t>(ObjectType::OBJECT), obj->header.type);
    EXPECT_EQ(10UL, obj->header.size);
    
    // Test the free space decreased
    EXPECT_LT(memory.getFreeSpace(), memory.getTotalSpace());
    EXPECT_GT(memory.getUsedSpace(), 0UL);
}

TEST(TestMemoryByteArrayAllocation) {
    MemoryManager memory;
    
    // Test allocating a byte array
    Object* bytes = memory.allocateBytes(100);
    EXPECT_NE(nullptr, bytes);
    EXPECT_EQ(static_cast<uint64_t>(ObjectType::BYTE_ARRAY), bytes->header.type);
    
    // Check the allocated size is properly aligned
    size_t alignedSize = (100 + 7) & ~7; // Align to 8 bytes
    EXPECT_EQ(alignedSize, bytes->header.size);
}

void runAllTests() {
    std::cout << "Running tests..." << std::endl;
    
    TestBytecodeInstructionSizes();
    std::cout << "✓ Bytecode instruction sizes test passed" << std::endl;
    
    TestBytecodeNames();
    std::cout << "✓ Bytecode names test passed" << std::endl;
    
    TestBytecodeValues();
    std::cout << "✓ Bytecode values test passed" << std::endl;
    
    TestMemoryObjectAllocation();
    std::cout << "✓ Memory object allocation test passed" << std::endl;
    
    TestMemoryByteArrayAllocation();
    std::cout << "✓ Memory byte array allocation test passed" << std::endl;
    
    std::cout << "All tests passed!" << std::endl;
}

int main() {
    try {
        runAllTests();
        return 0;
    } catch (const std::exception& e) {
        std::cerr << "Test failed: " << e.what() << std::endl;
        return 1;
    } catch (...) {
        std::cerr << "Test failed with unknown exception" << std::endl;
        return 1;
    }
}