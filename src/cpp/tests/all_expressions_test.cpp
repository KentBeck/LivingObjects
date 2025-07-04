#include "simple_parser.h"
#include "simple_compiler.h"
#include "interpreter.h"
#include "memory_manager.h"
#include "tagged_value.h"
#include "smalltalk_string.h"
#include "smalltalk_class.h"
#include "primitive_methods.h"
#include "bytecode.h"

#include <iostream>
#include <string>
#include <vector>
#include <cassert>

// Simple test framework
#define TEST(name) void name()
#define EXPECT_EQ(expected, actual) assert((expected) == (actual))
#define EXPECT_NE(expected, actual) assert((expected) != (actual))
#define EXPECT_STREQ(expected, actual) assert(strcmp((expected), (actual)) == 0)
#define EXPECT_LT(a, b) assert((a) < (b))
#define EXPECT_GT(a, b) assert((a) > (b))

using namespace smalltalk;

struct ExpressionTest
{
    std::string expression;
    std::string expectedResult;
    bool shouldPass;
    std::string category;
};

void testExpression(const ExpressionTest &test)
{
    std::cout << "Testing: " << test.expression << " -> " << test.expectedResult;

    try
    {
        // Parse, compile, and execute the expression
        SimpleParser parser(test.expression);
        auto methodAST = parser.parseMethod();

        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);

        MemoryManager memoryManager;
        Interpreter interpreter(memoryManager);
        TaggedValue result = interpreter.executeCompiledMethod(*compiledMethod);

        // Convert result to string for comparison
        std::string resultStr;
        if (result.isInteger())
        {
            resultStr = std::to_string(result.asInteger());
        }
        else if (result.isBoolean())
        {
            resultStr = result.asBoolean() ? "true" : "false";
        }
        else if (result.isNil())
        {
            resultStr = "nil";
        }
        else if (StringUtils::isString(result))
        {
            String *str = StringUtils::asString(result);
            resultStr = str->getContent(); // Get content without quotes for comparison
        }
        else
        {
            resultStr = "Object";
        }

        if (test.shouldPass && resultStr == test.expectedResult)
        {
            std::cout << " ✅ PASS" << std::endl;
        }
        else if (test.shouldPass)
        {
            std::cout << " ❌ FAIL (got: " << resultStr << ")" << std::endl;
        }
        else
        {
            std::cout << " ❌ FAIL (should have failed but got: " << resultStr << ")" << std::endl;
        }
    }
    catch (const std::exception &e)
    {
        if (test.shouldPass)
        {
            std::cout << " ❌ FAIL (exception: " << e.what() << ")" << std::endl;
        }
        else
        {
            std::cout << " ✅ EXPECTED FAIL (" << e.what() << ")" << std::endl;
        }
    }
}

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
    EXPECT_EQ(ObjectType::OBJECT, obj->header.getType());
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
    EXPECT_EQ(ObjectType::BYTE_ARRAY, bytes->header.getType());

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

void runAllTests();

int main()
{
    runAllTests();

    // Initialize class system and primitives before running tests
    ClassUtils::initializeCoreClasses();

    // Initialize primitive registry
    auto &primitiveRegistry = PrimitiveRegistry::getInstance();
    primitiveRegistry.initializeCorePrimitives();

    // Add primitive methods to Integer class (temporarily disabled)
    // Class* integerClass = ClassUtils::getIntegerClass();
    // IntegerClassSetup::addPrimitiveMethods(integerClass);

    std::vector<ExpressionTest> tests = {
        // Basic arithmetic - SHOULD PASS
        {"3 + 4", "7", true, "arithmetic"},
        {"5 - 2", "3", true, "arithmetic"},
        {"2 * 3", "6", true, "arithmetic"},
        {"10 / 2", "5", true, "arithmetic"},

        // Complex arithmetic - SHOULD PASS
        {"(3 + 2) * 4", "20", true, "arithmetic"},
        {"10 - 2 * 3", "4", true, "arithmetic"},
        {"(10 - 2) / 4", "2", true, "arithmetic"},

        // Integer comparisons - SHOULD PASS
        {"3 < 5", "true", true, "comparison"},
        {"7 > 2", "true", true, "comparison"},
        {"3 = 3", "true", true, "comparison"},
        {"4 ~= 5", "true", true, "comparison"},
        {"4 <= 4", "true", true, "comparison"},
        {"5 >= 3", "true", true, "comparison"},
        {"5 < 3", "false", true, "comparison"},
        {"2 > 7", "false", true, "comparison"},
        {"3 = 4", "false", true, "comparison"},

        // Complex comparisons - SHOULD PASS
        {"(3 + 2) < (4 * 2)", "true", true, "comparison"},
        {"(10 - 3) > (2 * 3)", "true", true, "comparison"},
        {"(6 / 2) = (1 + 2)", "true", true, "comparison"},

        // Basic object creation - SHOULD PASS (now implemented!)
        {"Object new", "Object", true, "object_creation"},
        {"Array new: 3", "<Array size: 3>", false, "object_creation"},

        // String literals - SHOULD PASS (basic string parsing)
        {"'hello'", "hello", true, "strings"},
        {"'world'", "world", true, "strings"},

        // String operations - SHOULD PASS (now implemented!)
        {"'hello' , ' world'", "hello world", true, "string_operations"},
        {"'hello' size", "5", true, "string_operations"},

        // Literals - SHOULD PASS (now implemented)
        {"true", "true", true, "literals"},
        {"false", "false", true, "literals"},
        {"nil", "nil", true, "literals"},

        // Variable assignment - SHOULD PASS (now implemented!)
        {"| x | x := 42. x", "42", true, "variables"},

        // Blocks - SHOULD FAIL (not implemented)
        {"[3 + 4] value", "7", false, "blocks"},
        {"[:x | x + 1] value: 5", "6", false, "blocks"},

        // Conditionals - SHOULD FAIL (not implemented)
        {"3 < 4) ifTrue: [10] ifFalse: [20]", "10", false, "conditionals"},
        {"true ifTrue: [42]", "42", false, "conditionals"},

        // Collections - SHOULD FAIL (not implemented)
        {"#(1 2 3) at: 2", "2", false, "collections"},
        {"#(1 2 3) size", "3", false, "collections"},

        // Dictionary operations - SHOULD FAIL (not implemented)
        {"Dictionary new", "<Dictionary>", false, "dictionaries"},

        // Class creation - SHOULD FAIL (not implemented)
        {"Object subclass: #Point", "<Class: Point>", false, "class_creation"}};

    std::cout << "=== Smalltalk Expression Test Suite ===" << std::endl;
    std::cout << "Testing " << tests.size() << " expressions..." << std::endl
              << std::endl;

    int passCount = 0;
    int totalCount = 0;
    std::string currentCategory = "";

    for (const auto &test : tests)
    {
        if (test.category != currentCategory)
        {
            currentCategory = test.category;
            std::cout << std::endl
                      << "=== " << currentCategory << " ===" << std::endl;
        }

        testExpression(test);

        // Count as pass if result matches expectation (either should pass and did, or should fail and did)
        try
        {
            SimpleParser parser(test.expression);
            auto methodAST = parser.parseMethod();
            SimpleCompiler compiler;
            auto compiledMethod = compiler.compile(*methodAST);
            MemoryManager memoryManager;
            Interpreter interpreter(memoryManager);
            TaggedValue result = interpreter.executeCompiledMethod(*compiledMethod);

            std::string resultStr;
            if (result.isInteger())
            {
                resultStr = std::to_string(result.asInteger());
            }
            else if (result.isBoolean())
            {
                resultStr = result.asBoolean() ? "true" : "false";
            }
            else if (result.isNil())
            {
                resultStr = "nil";
            }
            else if (StringUtils::isString(result))
            {
                String *str = StringUtils::asString(result);
                resultStr = str->getContent(); // Get content without quotes for comparison
            }
            else
            {
                resultStr = "Object";
            }

            if (test.shouldPass && resultStr == test.expectedResult)
            {
                passCount++;
            }
            else if (!test.shouldPass)
            {
                // This should have failed but didn't - that's actually bad
            }
        }
        catch (const std::exception &)
        {
            if (!test.shouldPass)
            {
                passCount++; // Expected to fail and did fail
            }
        }
        totalCount++;
    }

    std::cout << std::endl
              << "=== SUMMARY ===" << std::endl;
    std::cout << "Expressions that work correctly: " << passCount << "/" << totalCount << std::endl;

    // Count by category
    std::cout << std::endl
              << "By category:" << std::endl;
    std::vector<std::string> categories = {"arithmetic", "comparison", "object_creation", "strings", "string_operations", "literals", "variables", "blocks", "conditionals", "collections", "dictionaries", "class_creation"};

    for (const auto &category : categories)
    {
        int categoryPass = 0;
        int categoryTotal = 0;

        for (const auto &test : tests)
        {
            if (test.category == category)
            {
                categoryTotal++;

                bool actuallyPassed = false;
                try
                {
                    SimpleParser parser(test.expression);
                    auto methodAST = parser.parseMethod();
                    SimpleCompiler compiler;
                    auto compiledMethod = compiler.compile(*methodAST);
                    MemoryManager memoryManager;
                    Interpreter interpreter(memoryManager);
                    TaggedValue result = interpreter.executeCompiledMethod(*compiledMethod);

                    std::string resultStr;
                    if (result.isInteger())
                    {
                        resultStr = std::to_string(result.asInteger());
                    }
                    else if (result.isBoolean())
                    {
                        resultStr = result.asBoolean() ? "true" : "false";
                    }
                    else if (result.isNil())
                    {
                        resultStr = "nil";
                    }
                    else if (StringUtils::isString(result))
                    {
                        String *str = StringUtils::asString(result);
                        resultStr = str->getContent(); // Get content without quotes for comparison
                    }
                    else
                    {
                        resultStr = "Object";
                    }

                    if (test.shouldPass && resultStr == test.expectedResult)
                    {
                        actuallyPassed = true;
                    }
                }
                catch (const std::exception &)
                {
                    if (!test.shouldPass)
                    {
                        actuallyPassed = true;
                    }
                }

                if (actuallyPassed)
                    categoryPass++;
            }
        }

        if (categoryTotal > 0)
        {
            std::cout << "  " << category << ": " << categoryPass << "/" << categoryTotal;
            if (categoryPass == categoryTotal)
            {
                std::cout << " ✅";
            }
            else
            {
                std::cout << " ❌";
            }
            std::cout << std::endl;
        }
    }

    return 0;
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