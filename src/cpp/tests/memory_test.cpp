#include "memory_manager.h"

#include <gtest/gtest.h>

using namespace smalltalk;

TEST(MemoryTest, ObjectAllocation) {
    MemoryManager memory;
    
    // Test allocating a basic object
    Object* obj = memory.allocateObject(ObjectType::OBJECT, 10);
    ASSERT_NE(nullptr, obj);
    EXPECT_EQ(static_cast<uint64_t>(ObjectType::OBJECT), obj->header.type);
    EXPECT_EQ(10UL, obj->header.size);
    
    // Test the free space decreased
    EXPECT_LT(memory.getFreeSpace(), memory.getTotalSpace());
    EXPECT_GT(memory.getUsedSpace(), 0UL);
}

TEST(MemoryTest, ByteArrayAllocation) {
    MemoryManager memory;
    
    // Test allocating a byte array
    Object* bytes = memory.allocateBytes(100);
    ASSERT_NE(nullptr, bytes);
    EXPECT_EQ(static_cast<uint64_t>(ObjectType::BYTE_ARRAY), bytes->header.type);
    
    // Check the allocated size is properly aligned
    size_t alignedSize = (100 + 7) & ~7; // Align to 8 bytes
    EXPECT_EQ(alignedSize, bytes->header.size);
}

TEST(MemoryTest, ArrayAllocation) {
    MemoryManager memory;
    
    // Test allocating an array
    Object* array = memory.allocateArray(5);
    ASSERT_NE(nullptr, array);
    EXPECT_EQ(static_cast<uint64_t>(ObjectType::ARRAY), array->header.type);
    EXPECT_EQ(5UL, array->header.size);
    
    // Check that the array slots are initialized to null
    Object** slots = reinterpret_cast<Object**>(reinterpret_cast<char*>(array) + sizeof(Object));
    for (size_t i = 0; i < 5; i++) {
        EXPECT_EQ(nullptr, slots[i]);
    }
}

TEST(MemoryTest, ContextAllocation) {
    MemoryManager memory;
    
    // Create a self object
    Object* self = memory.allocateObject(ObjectType::OBJECT, 0);
    
    // Test allocating a method context
    MethodContext* context = memory.allocateMethodContext(5, 123, self, nullptr);
    ASSERT_NE(nullptr, context);
    EXPECT_EQ(static_cast<uint64_t>(ContextType::METHOD_CONTEXT), context->header.type);
    EXPECT_EQ(5UL, context->header.size);
    EXPECT_EQ(123U, context->header.hash);
    EXPECT_EQ(self, reinterpret_cast<Object**>(reinterpret_cast<char*>(context) + sizeof(Object))[1]);
    EXPECT_EQ(nullptr, reinterpret_cast<Object**>(reinterpret_cast<char*>(context) + sizeof(Object))[0]);
    
    // Test allocating a block context
    BlockContext* blockContext = memory.allocateBlockContext(3, 456, self, nullptr, static_cast<Object*>(context));
    ASSERT_NE(nullptr, blockContext);
    EXPECT_EQ(static_cast<uint64_t>(ContextType::BLOCK_CONTEXT), blockContext->header.type);
    EXPECT_EQ(3UL, blockContext->header.size);
    EXPECT_EQ(456U, blockContext->header.hash);
    EXPECT_EQ(self, reinterpret_cast<Object**>(reinterpret_cast<char*>(blockContext) + sizeof(Object))[1]);
    EXPECT_EQ(nullptr, reinterpret_cast<Object**>(reinterpret_cast<char*>(blockContext) + sizeof(Object))[0]);
    EXPECT_EQ(static_cast<Object*>(context), blockContext->home);
}

TEST(MemoryTest, StackChunkAllocation) {
    MemoryManager memory;
    
    // Test allocating a stack chunk
    StackChunk* chunk = memory.allocateStackChunk(100);
    ASSERT_NE(nullptr, chunk);
    EXPECT_EQ(static_cast<uint64_t>(ContextType::STACK_CHUNK_BOUNDARY), chunk->header.type);
    EXPECT_EQ(100UL, chunk->header.size);
    EXPECT_EQ(nullptr, chunk->previousChunk);
    EXPECT_EQ(nullptr, chunk->nextChunk);
    EXPECT_NE(nullptr, chunk->allocationPointer);
    
    // Clean up
    free(chunk);
}