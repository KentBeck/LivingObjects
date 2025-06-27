#include "../include/bytecode.h"
#include "../include/object.h"
#include "../include/memory_manager.h"
#include "../include/tagged_value.h"
#include <iostream>
#include <cassert>

using namespace smalltalk;

// Simple test framework
#define TEST(name) void name()
#define EXPECT_EQ(expected, actual) assert((expected) == (actual))
#define EXPECT_NE(expected, actual) assert((expected) != (actual))
#define EXPECT_STREQ(expected, actual) assert(strcmp((expected), (actual)) == 0)
#define EXPECT_LT(a, b) assert((a) < (b))
#define EXPECT_GT(a, b) assert((a) > (b))

TEST(TestBytecodeInstructionSizes)
{
    // Test instruction sizes match the Go implementation
    EXPECT_EQ(INSTRUCTION_SIZE_FOUR_BYTE_OPERAND, getInstructionSize(Bytecode::PUSH_LITERAL));
    EXPECT_EQ(INSTRUCTION_SIZE_FOUR_BYTE_OPERAND, getInstructionSize(Bytecode::PUSH_INSTANCE_VARIABLE));
    EXPECT_EQ(INSTRUCTION_SIZE_FOUR_BYTE_OPERAND, getInstructionSize(Bytecode::PUSH_TEMPORARY_VARIABLE));
    EXPECT_EQ(INSTRUCTION_SIZE_ONE_BYTE_OPCODE, getInstructionSize(Bytecode::PUSH_SELF));
    EXPECT_EQ(INSTRUCTION_SIZE_FOUR_BYTE_OPERAND, getInstructionSize(Bytecode::STORE_INSTANCE_VARIABLE));
    EXPECT_EQ(INSTRUCTION_SIZE_FOUR_BYTE_OPERAND, getInstructionSize(Bytecode::STORE_TEMPORARY_VARIABLE));
    EXPECT_EQ(INSTRUCTION_SIZE_SEND_MESSAGE, getInstructionSize(Bytecode::SEND_MESSAGE));
    EXPECT_EQ(INSTRUCTION_SIZE_ONE_BYTE_OPCODE, getInstructionSize(Bytecode::RETURN_STACK_TOP));
    EXPECT_EQ(INSTRUCTION_SIZE_FOUR_BYTE_OPERAND, getInstructionSize(Bytecode::JUMP));
    EXPECT_EQ(INSTRUCTION_SIZE_FOUR_BYTE_OPERAND, getInstructionSize(Bytecode::JUMP_IF_TRUE));
    EXPECT_EQ(INSTRUCTION_SIZE_FOUR_BYTE_OPERAND, getInstructionSize(Bytecode::JUMP_IF_FALSE));
    EXPECT_EQ(INSTRUCTION_SIZE_ONE_BYTE_OPCODE, getInstructionSize(Bytecode::POP));
    EXPECT_EQ(INSTRUCTION_SIZE_ONE_BYTE_OPCODE, getInstructionSize(Bytecode::DUPLICATE));
    EXPECT_EQ(INSTRUCTION_SIZE_CREATE_BLOCK, getInstructionSize(Bytecode::CREATE_BLOCK));
    EXPECT_EQ(INSTRUCTION_SIZE_FOUR_BYTE_OPERAND, getInstructionSize(Bytecode::EXECUTE_BLOCK));
}

TEST(TestBytecodeNames)
{
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

TEST(TestBytecodeValues)
{
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

TEST(TestMemoryObjectAllocation)
{
    MemoryManager memory;

    // Test allocating a basic object
    Object *obj = memory.allocateObject(ObjectType::OBJECT, 10);
    EXPECT_NE(nullptr, obj);
    EXPECT_EQ(static_cast<uint64_t>(ObjectType::OBJECT), obj->header.type);
    EXPECT_EQ(10UL, obj->header.size);

    // Test the free space decreased
    EXPECT_LT(memory.getFreeSpace(), memory.getTotalSpace());
    EXPECT_GT(memory.getUsedSpace(), 0UL);
}

TEST(TestMemoryByteArrayAllocation)
{
    MemoryManager memory;

    // Test allocating a byte array
    Object *bytes = memory.allocateBytes(100);
    EXPECT_NE(nullptr, bytes);
    EXPECT_EQ(static_cast<uint64_t>(ObjectType::BYTE_ARRAY), bytes->header.type);

    // Check the allocated size is properly aligned
    size_t alignedSize = (100 + 7) & ~7; // Align to 8 bytes
    EXPECT_EQ(alignedSize, bytes->header.size);
}

TEST(TestTaggedValueInteger)
{
    // Test creating integer 3
    TaggedValue three(3);

    // Verify it's recognized as an integer
    EXPECT_EQ(true, three.isInteger());
    EXPECT_EQ(false, three.isPointer());
    EXPECT_EQ(false, three.isSpecial());
    EXPECT_EQ(false, three.isFloat());

    // Verify the value can be extracted
    EXPECT_EQ(3, three.asInteger());
}

TEST(TestTaggedValueIntegerRange)
{
    // Test various integer values
    TaggedValue zero(0);
    TaggedValue positive(42);
    TaggedValue negative(-17);
    TaggedValue large(1000000);

    EXPECT_EQ(0, zero.asInteger());
    EXPECT_EQ(42, positive.asInteger());
    EXPECT_EQ(-17, negative.asInteger());
    EXPECT_EQ(1000000, large.asInteger());

    // All should be integers
    EXPECT_EQ(true, zero.isInteger());
    EXPECT_EQ(true, positive.isInteger());
    EXPECT_EQ(true, negative.isInteger());
    EXPECT_EQ(true, large.isInteger());
}

TEST(TestTaggedValueSpecialValues)
{
    // Test nil, true, false
    TaggedValue nil = TaggedValue::nil();
    TaggedValue trueVal = TaggedValue::trueValue();
    TaggedValue falseVal = TaggedValue::falseValue();

    EXPECT_EQ(true, nil.isNil());
    EXPECT_EQ(true, trueVal.isTrue());
    EXPECT_EQ(true, falseVal.isFalse());

    EXPECT_EQ(true, nil.isSpecial());
    EXPECT_EQ(true, trueVal.isSpecial());
    EXPECT_EQ(true, falseVal.isSpecial());
}

void runAllTests()
{
    std::cout << "Running tests..." << '\n';

    TestBytecodeInstructionSizes();
    std::cout << "✓ Bytecode instruction sizes test passed" << '\n';

    TestBytecodeNames();
    std::cout << "✓ Bytecode names test passed" << '\n';

    TestBytecodeValues();
    std::cout << "✓ Bytecode values test passed" << '\n';

    TestMemoryObjectAllocation();
    std::cout << "✓ Memory object allocation test passed" << '\n';

    TestMemoryByteArrayAllocation();
    std::cout << "✓ Memory byte array allocation test passed" << '\n';

    TestTaggedValueInteger();
    std::cout << "✓ Tagged value integer test passed" << '\n';

    TestTaggedValueIntegerRange();
    std::cout << "✓ Tagged value integer range test passed" << '\n';

    TestTaggedValueSpecialValues();
    std::cout << "✓ Tagged value special values test passed" << '\n';

    std::cout << "All tests passed!" << '\n';
}

int main()
{
    try
    {
        runAllTests();
        return 0;
    }
    catch (const std::exception &e)
    {
        std::cerr << "Test failed: " << e.what() << '\n';
        return 1;
    }
    catch (...)
    {
        std::cerr << "Test failed with unknown exception" << '\n';
        return 1;
    }
}
