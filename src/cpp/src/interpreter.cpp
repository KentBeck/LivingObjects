#include "interpreter.h"
#include "primitives.h"
#include "smalltalk_class.h"
#include "symbol.h"
#include "simple_parser.h"

#include <cstring>
#include <stdexcept>
#include <vector>

namespace smalltalk
{

    Interpreter::Interpreter(MemoryManager &memory)
        : memoryManager(memory)
    {
        // Initialize core classes and primitives on first use
        static bool systemInitialized = false;
        if (!systemInitialized)
        {
            ClassUtils::initializeCoreClasses();
            Primitives::initialize();
            systemInitialized = true;
        }
        // Initialize the stack chunk
        currentChunk = memoryManager.allocateStackChunk(1024);
    }

    Object *Interpreter::executeMethod(Object *method, Object *receiver, std::vector<Object *> &args)
    {
        // Create a new method context
        uint32_t methodObj = static_cast<uint32_t>(reinterpret_cast<uintptr_t>(method));
        MethodContext *context = memoryManager.allocateMethodContext(10 + args.size(), methodObj, receiver, nullptr);

        // Get the variable-sized storage area safely
        // Memory layout: [MethodContext][Object* slots...]
        // The stackPointer will point into this Object* array
        char *contextEnd = reinterpret_cast<char *>(context) + sizeof(MethodContext);
        Object **slots = reinterpret_cast<Object **>(contextEnd);

        // Validate alignment - Object** must be properly aligned
        if (reinterpret_cast<uintptr_t>(slots) % alignof(Object *) != 0)
        {
            throw std::runtime_error("Stack slots not properly aligned");
        }

        // Copy arguments to the context
        for (size_t i = 0; i < args.size(); i++)
        {
            slots[i] = args[i];
        }

        // Set up stack pointer to point to the first available slot after arguments
        // Note: stackPointer is stored as Object* but represents a position in Object** array
        Object **initialStackPos = slots + args.size();
        context->stackPointer = reinterpret_cast<Object *>(initialStackPos);

        // Execute the context
        return executeContext(context);
    }

    Object *Interpreter::executeContext(MethodContext *context)
    {
        // Save current context
        MethodContext *previousContext = activeContext;

        // Set new active context
        activeContext = context;

        // Execute bytecodes until context returns
        executeLoop();

        // Get the return value (top of stack)
        Object *result = nullptr;
        if (activeContext != nullptr)
        {
            result = top();
        }

        // Restore previous context
        activeContext = previousContext;

        return result;
    }

    TaggedValue Interpreter::executeCompiledMethod(const CompiledMethod &method)
    {
        const auto &bytecodes = method.getBytecodes();
        const auto &literals = method.getLiterals();

        // Set up execution state
        std::vector<TaggedValue> stack;
        size_t ip = 0;

        // Main bytecode execution loop - process one instruction at a time
        while (ip < bytecodes.size())
        {
            uint8_t opcode = bytecodes[ip];
            Bytecode instruction = static_cast<Bytecode>(opcode);

            switch (instruction)
            {
            case Bytecode::PUSH_LITERAL:
            {
                ip++; // Skip opcode
                if (ip + 3 >= bytecodes.size())
                {
                    throw std::runtime_error("Invalid PUSH_LITERAL: not enough bytes for operand");
                }

                uint32_t literalIndex =
                    static_cast<uint32_t>(bytecodes[ip]) |
                    (static_cast<uint32_t>(bytecodes[ip + 1]) << 8) |
                    (static_cast<uint32_t>(bytecodes[ip + 2]) << 16) |
                    (static_cast<uint32_t>(bytecodes[ip + 3]) << 24);
                ip += 4;

                if (literalIndex >= literals.size())
                {
                    throw std::runtime_error("Invalid literal index: " + std::to_string(literalIndex));
                }

                stack.push_back(literals[literalIndex]);
                break;
            }

            case Bytecode::SEND_MESSAGE:
            {
                ip++; // Skip opcode
                if (ip + 7 >= bytecodes.size())
                {
                    throw std::runtime_error("Invalid SEND_MESSAGE: not enough bytes for operands");
                }

                // Read selector index
                uint32_t selectorIndex =
                    static_cast<uint32_t>(bytecodes[ip]) |
                    (static_cast<uint32_t>(bytecodes[ip + 1]) << 8) |
                    (static_cast<uint32_t>(bytecodes[ip + 2]) << 16) |
                    (static_cast<uint32_t>(bytecodes[ip + 3]) << 24);
                ip += 4;

                // Read argument count
                uint32_t argCount =
                    static_cast<uint32_t>(bytecodes[ip]) |
                    (static_cast<uint32_t>(bytecodes[ip + 1]) << 8) |
                    (static_cast<uint32_t>(bytecodes[ip + 2]) << 16) |
                    (static_cast<uint32_t>(bytecodes[ip + 3]) << 24);
                ip += 4;

                // Validate arguments
                if (stack.size() < argCount + 1)
                {
                    throw std::runtime_error("Stack underflow in SEND_MESSAGE");
                }
                if (selectorIndex >= literals.size())
                {
                    throw std::runtime_error("Invalid selector index: " + std::to_string(selectorIndex));
                }

                // Pop arguments and receiver
                std::vector<TaggedValue> args;
                for (uint32_t i = 0; i < argCount; i++)
                {
                    args.insert(args.begin(), stack.back());
                    stack.pop_back();
                }
                TaggedValue receiver = stack.back();
                stack.pop_back();

                // Perform the operation
                TaggedValue selector = literals[selectorIndex];
                TaggedValue result;

                if (argCount == 0)
                {
                    // Unary message - use new message sending system
                    if (selector.isPointer())
                    {
                        Object *selectorObj = selector.asObject();
                        Symbol *selectorSymbol = reinterpret_cast<Symbol *>(selectorObj);
                        result = sendMessage(receiver, selectorSymbol->getName(), args);
                    }
                    else
                    {
                        result = TaggedValue(); // nil
                    }
                }
                else if (argCount == 1)
                {
                    result = performOperation(receiver, args[0], selector);
                }
                else
                {
                    // For now, only support unary and binary operations
                    result = TaggedValue(); // nil
                }

                stack.push_back(result);
                break;
            }

            case Bytecode::CREATE_BLOCK:
            {
                ip++; // Skip opcode
                if (ip + 11 >= bytecodes.size())
                {
                    throw std::runtime_error("Invalid CREATE_BLOCK: not enough bytes for operands");
                }

                // Read block parameters (we skip them for now)
                ip += 12; // Skip 3 * 4-byte parameters

                // Create a temporary context to hold the block
                MethodContext *tempContext = memoryManager.allocateMethodContext(4, 0, nullptr, nullptr);
                tempContext->self = reinterpret_cast<Object *>(0x1000);
                MethodContext *oldContext = activeContext;
                activeContext = tempContext;

                // Execute CREATE_BLOCK handler
                handleCreateBlock(0, 0, 0);
                Object *blockObj = pop();

                // Restore context
                activeContext = oldContext;

                // Push block as TaggedValue
                stack.push_back(TaggedValue(blockObj));
                break;
            }

            case Bytecode::RETURN_STACK_TOP:
            {
                if (stack.empty())
                {
                    return TaggedValue(); // Return nil if stack is empty
                }
                return stack.back();
            }

            default:
                // Skip unknown instructions
                ip++;
                break;
            }
        }

        // If we reach here, return top of stack or nil
        return stack.empty() ? TaggedValue() : stack.back();
    }

    TaggedValue Interpreter::performOperation(const TaggedValue &left, const TaggedValue &right, const TaggedValue &selector)
    {
        if (!selector.isPointer())
            return TaggedValue();

        // Try to get the selector as a symbol and extract the operator string
        std::string op;
        try
        {
            Symbol *sym = selector.asSymbol();
            op = sym->toString();
        }
        catch (...)
        {
            // If it's not a symbol, try as a string or return nil
            return TaggedValue();
        }

        if (left.isInteger() && right.isInteger())
        {
            int32_t l = left.asInteger();
            int32_t r = right.asInteger();

            if (op == "+" || op == "#+")
                return TaggedValue(l + r);
            if (op == "-" || op == "#-")
                return TaggedValue(l - r);
            if (op == "*" || op == "#*")
                return TaggedValue(l * r);
            if (op == "/" || op == "#/")
                return TaggedValue(r != 0 ? l / r : 0);
            if (op == "<" || op == "#<")
                return l < r ? TaggedValue::trueValue() : TaggedValue::falseValue();
            if (op == ">" || op == "#>")
                return l > r ? TaggedValue::trueValue() : TaggedValue::falseValue();
            if (op == "=" || op == "#=")
                return l == r ? TaggedValue::trueValue() : TaggedValue::falseValue();
            if (op == "~=" || op == "#~=")
                return l != r ? TaggedValue::trueValue() : TaggedValue::falseValue();
            if (op == "<=" || op == "#<=")
                return l <= r ? TaggedValue::trueValue() : TaggedValue::falseValue();
            if (op == ">=" || op == "#>=")
                return l >= r ? TaggedValue::trueValue() : TaggedValue::falseValue();
        }

        return TaggedValue();
    }

    void Interpreter::executeLoop()
    {
        executing = true;

        while (executing && (activeContext != nullptr))
        {
            // For a basic implementation, we'll just pretend we're executing bytecodes
            // In a real implementation, we would fetch bytecodes from the method

            // For the minimal example, let's just execute a simple sequence
            // that pushes a value and returns
            handlePushSelf();
            handleReturnStackTop();

            // Stop execution after this simple sequence
            executing = false;
        }
    }

    void Interpreter::dispatch(Bytecode bytecode)
    {
        // Get current instruction pointer
        size_t ip = activeContext->instructionPointer;

        // Update instruction pointer for next instruction
        activeContext->instructionPointer += getInstructionSize(bytecode);

        // Dispatch based on bytecode
        switch (bytecode)
        {
        case Bytecode::PUSH_LITERAL:
            handlePushLiteral(readUInt32(ip + 1));
            break;
        case Bytecode::PUSH_INSTANCE_VARIABLE:
            handlePushInstanceVariable(readUInt32(ip + 1));
            break;
        case Bytecode::PUSH_TEMPORARY_VARIABLE:
            handlePushTemporaryVariable(readUInt32(ip + 1));
            break;
        case Bytecode::PUSH_SELF:
            handlePushSelf();
            break;
        case Bytecode::STORE_INSTANCE_VARIABLE:
            handleStoreInstanceVariable(readUInt32(ip + 1));
            break;
        case Bytecode::STORE_TEMPORARY_VARIABLE:
            handleStoreTemporaryVariable(readUInt32(ip + 1));
            break;
        case Bytecode::SEND_MESSAGE:
            handleSendMessage(readUInt32(ip + 1), readUInt32(ip + 5));
            break;
        case Bytecode::RETURN_STACK_TOP:
            handleReturnStackTop();
            break;
        case Bytecode::JUMP:
            handleJump(readUInt32(ip + 1));
            break;
        case Bytecode::JUMP_IF_TRUE:
            handleJumpIfTrue(readUInt32(ip + 1));
            break;
        case Bytecode::JUMP_IF_FALSE:
            handleJumpIfFalse(readUInt32(ip + 1));
            break;
        case Bytecode::POP:
            handlePop();
            break;
        case Bytecode::DUPLICATE:
            handleDuplicate();
            break;
        case Bytecode::CREATE_BLOCK:
            handleCreateBlock(readUInt32(ip + 1), readUInt32(ip + 5), readUInt32(ip + 9));
            break;
        case Bytecode::EXECUTE_BLOCK:
            handleExecuteBlock(readUInt32(ip + 1));
            break;
        default:
            throw std::runtime_error("Unknown bytecode");
        }
    }

    uint32_t Interpreter::readUInt32(size_t offset)
    {
        // In a real implementation, this would read from actual bytecodes
        // For this stub, just return a placeholder value
        (void)offset; // Suppress unused parameter warning
        return 0;
    }

    void Interpreter::push(Object *value)
    {
        if (activeContext == nullptr)
        {
            throw std::runtime_error("No active context for push operation");
        }

        // Get current stack pointer and validate bounds
        Object **currentSP = getCurrentStackPointer(activeContext);
        Object **stackEnd = getStackEnd(activeContext);

        // Check for stack overflow
        if (currentSP >= stackEnd)
        {
            throw std::runtime_error("Stack overflow");
        }

        // Validate that we're writing to a valid Object* slot
        validateStackBounds(activeContext, currentSP);

        // Push value and update stack pointer
        *currentSP = value;
        Object **newStackPos = currentSP + 1;
        activeContext->stackPointer = reinterpret_cast<Object *>(newStackPos);
    }

    Object *Interpreter::pop()
    {
        if (activeContext == nullptr)
        {
            throw std::runtime_error("No active context for pop operation");
        }

        // Get current stack pointer and validate bounds
        Object **currentSP = getCurrentStackPointer(activeContext);
        Object **stackStart = getStackStart(activeContext);

        // Check for stack underflow
        if (currentSP <= stackStart)
        {
            throw std::runtime_error("Stack underflow");
        }

        // Move stack pointer back and validate
        Object **newStackPos = currentSP - 1;
        validateStackBounds(activeContext, newStackPos);

        // Get value and update stack pointer
        Object *value = *newStackPos;
        activeContext->stackPointer = reinterpret_cast<Object *>(newStackPos);
        return value;
    }

    Object *Interpreter::top()
    {
        if (activeContext == nullptr)
        {
            throw std::runtime_error("No active context for top operation");
        }

        // Get current stack pointer and validate bounds
        Object **currentSP = getCurrentStackPointer(activeContext);
        Object **stackStart = getStackStart(activeContext);

        // Check for empty stack
        if (currentSP <= stackStart)
        {
            throw std::runtime_error("Stack is empty");
        }

        // Return top value without modifying stack pointer
        return *(currentSP - 1);
    }

    // Bytecode handler implementations
    void Interpreter::handlePushLiteral(uint32_t index)
    {
        // In a real implementation, this would get the literal from the method's literal array
        // For now, we just create a new object
        (void)index; // Suppress unused parameter warning
        Object *literal = memoryManager.allocateObject(ObjectType::OBJECT, 0);
        push(literal);
    }

    void Interpreter::handlePushInstanceVariable(uint32_t offset)
    {
        // Get the instance variable from the receiver at the given offset
        // For now, just push a new object
        (void)offset; // Suppress unused parameter warning
        Object *value = memoryManager.allocateObject(ObjectType::OBJECT, 0);
        push(value);
    }

    void Interpreter::handlePushTemporaryVariable(uint32_t offset)
    {
        // Get the temporary variable at the given offset
        Object **slots = reinterpret_cast<Object **>(reinterpret_cast<char *>(activeContext) + sizeof(MethodContext));
        push(slots[offset]);
    }

    void Interpreter::handlePushSelf()
    {
        // Push the receiver onto the stack
        push(activeContext->self);
    }

    void Interpreter::handleStoreInstanceVariable(uint32_t offset)
    {
        // Store the top of the stack into the instance variable at the given offset
        // Not implemented in this basic version
        (void)offset; // Suppress unused parameter warning
    }

    void Interpreter::handleStoreTemporaryVariable(uint32_t offset)
    {
        // Store the top of the stack into the temporary variable at the given offset
        Object *value = pop();
        Object **slots = reinterpret_cast<Object **>(reinterpret_cast<char *>(activeContext) + sizeof(MethodContext));
        slots[offset] = value;
        push(value); // Leave the value on the stack
    }

    void Interpreter::handleSendMessage(uint32_t selectorIndex, uint32_t argCount)
    {
        // Not fully implemented in this basic version
        // In a real implementation, this would look up the method and execute it
        (void)selectorIndex; // Suppress unused parameter warning

        // Pop arguments
        std::vector<Object *> args;
        args.reserve(argCount);
        for (uint32_t i = 0; i < argCount; i++)
        {
            args.push_back(pop());
        }

        // Pop receiver
        Object *receiver = pop();

        // Create a selector object
        Object *selector = memoryManager.allocateObject(ObjectType::SYMBOL, 0);

        // Send the message
        Object *result = sendMessage(receiver, selector, args);

        // Push the result
        push(result);
    }

    void Interpreter::handleReturnStackTop()
    {
        // Return from the current context
        Object *result = top();

        // Set the sender as the new active context
        MethodContext *sender = reinterpret_cast<MethodContext *>(activeContext->sender);
        activeContext = sender;

        // If there is a sender, push the result onto its stack
        if (activeContext != nullptr)
        {
            push(result);
        }
    }

    void Interpreter::handleJump(uint32_t target)
    {
        // Jump to the target instruction
        activeContext->instructionPointer = target;
    }

    void Interpreter::handleJumpIfTrue(uint32_t target)
    {
        // Pop the condition
        Object *condition = pop();

        // Check if the condition is true (simplified)
        bool isTrue = (condition != nullptr);

        // Jump if true
        if (isTrue)
        {
            activeContext->instructionPointer = target;
        }
    }

    void Interpreter::handleJumpIfFalse(uint32_t target)
    {
        // Pop the condition
        Object *condition = pop();

        // Check if the condition is false (simplified)
        bool isFalse = (condition == nullptr);

        // Jump if false
        if (isFalse)
        {
            activeContext->instructionPointer = target;
        }
    }

    void Interpreter::handlePop()
    {
        // Pop the top value from the stack
        pop();
    }

    void Interpreter::handleDuplicate()
    {
        // Duplicate the top value on the stack
        Object *value = top();
        push(value);
    }

    void Interpreter::handleCreateBlock(uint32_t bytecodeSize, uint32_t literalCount, uint32_t tempVarCount)
    {
        // Create a proper block context
        (void)bytecodeSize; // Will be used when we have proper block compilation
        (void)literalCount; // Will be used when we have proper block compilation
        (void)tempVarCount; // Will be used when we have proper block compilation

        // Get the current context as the home context for the block
        MethodContext *homeContext = activeContext;
        if (homeContext == nullptr)
        {
            throw std::runtime_error("Cannot create block without an active context");
        }

        // Create a block context
        // For now, use a placeholder method reference (will be improved later)
        uint32_t blockMethodRef = 999; // Placeholder method reference

        BlockContext *blockContext = memoryManager.allocateBlockContext(
            4,                 // context size (enough for basic temporaries)
            blockMethodRef,    // method reference for the block's code
            homeContext->self, // receiver (same as home context)
            nullptr,           // sender (will be set when block is executed)
            homeContext        // home context
        );

        // Push the block context onto the stack
        push(blockContext);
    }

    void Interpreter::handleExecuteBlock(uint32_t argCount)
    {
        // For now, we only support `value` with no arguments
        (void)argCount;

        // Pop the block context from the stack
        Object *blockObj = pop();
        if (blockObj->header.getType() != ObjectType::CONTEXT ||
            blockObj->header.getContextType() != static_cast<uint8_t>(ContextType::BLOCK_CONTEXT))
        {
            throw std::runtime_error("Object on stack is not a block for EXECUTE_BLOCK");
        }
        BlockContext *block = static_cast<BlockContext *>(blockObj);

        // Get the home context
        MethodContext *home = static_cast<MethodContext *>(block->home);

        // Get block's compiled method from the image using the reference in the block
        uint32_t methodRef = block->header.hash;
        // TODO: This depends on `image` being available and `getCompiledMethod`
        // and `CompiledMethod::tempCount` being correct.
        // CompiledMethod* compiledBlock = image->getCompiledMethod(methodRef);
        // if (!compiledBlock) {
        //     throw std::runtime_error("Could not find compiled method for block");
        // }

        // A block executes in the scope of its home context, so it uses the home's receiver
        Object *receiver = home->self;

        // The sender of the new context is the currently active context
        MethodContext *sender = activeContext;

        // Allocate a new context for the block
        // HACK: Assuming tempCount of 16 until CompiledMethod is available
        MethodContext *newContext = memoryManager.allocateMethodContext(
            16, methodRef, receiver, sender);
        newContext->instructionPointer = 0; // Start at the beginning of the block's bytecodes

        // Activate the new context by making it the current one
        activeContext = newContext;
    }

    Object *Interpreter::sendMessage(Object *receiver, Object *selector, std::vector<Object *> &args)
    {
        // Convert to TaggedValue for new message sending
        TaggedValue tvReceiver = TaggedValue::fromObject(receiver);
        std::vector<TaggedValue> tvArgs;
        tvArgs.reserve(args.size());
        for (Object *arg : args)
        {
            tvArgs.push_back(TaggedValue::fromObject(arg));
        }

        // Get selector string
        std::string selectorString;
        if (selector != nullptr && selector->header.getType() == ObjectType::SYMBOL)
        {
            // Symbol inherits from Object, so this cast is safe
            Symbol *sym = reinterpret_cast<Symbol *>(selector);
            selectorString = sym->getName();
        }
        else
        {
            throw std::runtime_error("Invalid selector in message send");
        }

        TaggedValue result = sendMessage(tvReceiver, selectorString, tvArgs);
        return result.asObject();
    }

    TaggedValue Interpreter::sendMessage(TaggedValue receiver, const std::string &selector, const std::vector<TaggedValue> &args)
    {
        // Get receiver's class
        Class *receiverClass = getObjectClass(receiver);
        if (receiverClass == nullptr)
        {
            throw std::runtime_error("Cannot determine receiver class");
        }

        // Create selector symbol
        Symbol *selectorSymbol = Symbol::intern(selector);

        // Look up method
        std::shared_ptr<CompiledMethod> method = receiverClass->lookupMethod(selectorSymbol);

        if (method && method->primitiveNumber != 0)
        {
            // Try primitive first
            try
            {
                return Primitives::callPrimitive(method->primitiveNumber, receiver, args, *this);
            }
            catch (const PrimitiveFailure &e)
            {
                // Fall back to Smalltalk code (not implemented yet)
                throw std::runtime_error("Primitive failed and fallback not implemented: " + std::string(e.what()));
            }
        }

        // No method found or no primitive - error for now
        throw std::runtime_error("Method not found: " + selector);
    }

    Class *Interpreter::getObjectClass(TaggedValue value)
    {
        if (value.isSmallInteger())
        {
            return ClassUtils::getIntegerClass();
        }
        if (value.isBoolean())
        {
            return ClassUtils::getBooleanClass();
        }
        if (value.isNil())
        {
            return ClassRegistry::getInstance().getClass("UndefinedObject");
        }
        if (value.isPointer())
        {
            return value.asObject()->getClass();
        }

        throw std::runtime_error("Unknown value type");
    }

    // Temporarily removed - will implement proper message send parsing later

    void Interpreter::switchContext(MethodContext *newContext)
    {
        // Set the new active context
        activeContext = newContext;
    }

    // Stack bounds checking helper methods
    Object **Interpreter::getStackStart(MethodContext *context)
    {
        if (context == nullptr)
        {
            throw std::runtime_error("Cannot get stack start for null context");
        }
        char *contextStart = reinterpret_cast<char *>(context) + sizeof(MethodContext);
        return reinterpret_cast<Object **>(contextStart);
    }

    Object **Interpreter::getStackEnd(MethodContext *context)
    {
        if (context == nullptr)
        {
            throw std::runtime_error("Cannot get stack end for null context");
        }
        Object **stackStart = getStackStart(context);
        return stackStart + context->header.size;
    }

    Object **Interpreter::getCurrentStackPointer(MethodContext *context)
    {
        if (context == nullptr)
        {
            throw std::runtime_error("Cannot get stack pointer for null context");
        }

        // Convert stored Object* back to Object** with validation
        Object **stackPointer = reinterpret_cast<Object **>(context->stackPointer);

        // Validate alignment
        if (reinterpret_cast<uintptr_t>(stackPointer) % alignof(Object *) != 0)
        {
            throw std::runtime_error("Stack pointer is not properly aligned");
        }

        return stackPointer;
    }

    void Interpreter::validateStackBounds(MethodContext *context, Object **stackPointer)
    {
        if (context == nullptr)
        {
            throw std::runtime_error("Cannot validate bounds for null context");
        }

        Object **stackStart = getStackStart(context);
        Object **stackEnd = getStackEnd(context);

        if (stackPointer < stackStart)
        {
            throw std::runtime_error("Stack pointer below stack start");
        }

        if (stackPointer > stackEnd)
        {
            throw std::runtime_error("Stack pointer above stack end");
        }
    }

} // namespace smalltalk
