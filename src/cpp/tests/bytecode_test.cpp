#include "bytecode.h"
#include <gtest/gtest.h>

using namespace smalltalk;

TEST(BytecodeTest, InstructionSizes) {
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

TEST(BytecodeTest, BytecodeNames) {
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

TEST(BytecodeTest, BytecodeValues) {
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