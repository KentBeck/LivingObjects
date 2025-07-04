#include "interpreter.h"
#include "smalltalk_vm.h"
#include "primitives.h"
#include "smalltalk_class.h"
#include "symbol.h"
#include "simple_parser.h"

#include <cstring>
#include <stdexcept>
#include <vector>
#include <iostream>

namespace smalltalk
{

    Interpreter::Interpreter(MemoryManager &memory)
        : memoryManager(memory)
    {
        // Ensure VM is initialized before any operations
        if (!SmalltalkVM::isInitialized()) {
            SmalltalkVM::initialize();
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
        // Memory layout: [MethodContext][TaggedValue slots...]
        // The stackPointer will point into this TaggedValue array
        char *contextEnd = reinterpret_cast<char *>(context) + sizeof(MethodContext);
        TaggedValue *slots = reinterpret_cast<TaggedValue *>(contextEnd);

        // Validate alignment - TaggedValue must be properly aligned
        if (reinterpret_cast<uintptr_t>(slots) % alignof(TaggedValue) != 0)
        {
            throw std::runtime_error("Stack slots not properly aligned");
        }

        // Copy arguments to the context
        for (size_t i = 0; i < args.size(); i++)
        {
            slots[i] = TaggedValue(args[i]);
        }

        // Set up stack pointer to point to the first available slot after arguments
        TaggedValue *initialStackPos = slots + args.size();
        context->stackPointer = initialStackPos;

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
        TaggedValue result = TaggedValue::nil();
        if (activeContext != nullptr)
        {
            result = top();
        }

        // Restore previous context
        activeContext = previousContext;

        // Convert TaggedValue result to Object* for legacy compatibility
        if (result.isPointer())
        {
            return result.asObject();
        }
        else
        {
            // For immediate values, we need to box them
            // This is a temporary solution - ideally the whole system should use TaggedValue
            return nullptr; // TODO: Implement proper boxing
        }
    }

    TaggedValue Interpreter::executeCompiledMethod(const CompiledMethod &method)
    {
        const auto &bytecodes = method.getBytecodes();
        const auto &literals = method.getLiterals();

        // Create a method context for execution
        Object *self = memoryManager.allocateObject(ObjectType::OBJECT, 0); // Simple self object
        MethodContext *methodContext = memoryManager.allocateMethodContext(
            16,     // context size (enough for stack and temporaries)
            0,      // method reference (placeholder)
            self,   // self
            nullptr // sender
        );

        // Initialize the stack pointer properly
        char *contextEnd = reinterpret_cast<char *>(methodContext) + sizeof(MethodContext);
        TaggedValue *slots = reinterpret_cast<TaggedValue *>(contextEnd);
        methodContext->stackPointer = slots;

        // Set up execution state using context-based stack
        MethodContext *savedContext = activeContext;
        activeContext = methodContext;
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

                // Push TaggedValue directly
                TaggedValue literal = literals[literalIndex];
                push(literal);
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

                if (selectorIndex >= literals.size())
                {
                    throw std::runtime_error("Invalid selector index: " + std::to_string(selectorIndex));
                }

                // Get the selector from literals
                TaggedValue selectorValue = literals[selectorIndex];
                if (!selectorValue.isPointer())
                {
                    throw std::runtime_error("Selector is not a pointer");
                }

                // Try to get the selector as a Symbol
                Symbol *selector;
                try
                {
                    selector = selectorValue.asSymbol();
                }
                catch (const std::exception &)
                {
                    throw std::runtime_error("Selector is not a symbol");
                }

                std::string selectorString = selector->getName();

                // Pop arguments from stack
                std::vector<TaggedValue> args;
                args.reserve(argCount);
                for (uint32_t i = 0; i < argCount; i++)
                {
                    args.push_back(pop());
                }

                // Pop receiver from stack
                TaggedValue receiver = pop();

                // Send the message
                TaggedValue result = sendMessage(receiver, selectorString, args);

                // Push result directly
                push(result);
                break;
            }

            case Bytecode::CREATE_BLOCK:
            {
                ip++; // Skip opcode
                if (ip + 11 >= bytecodes.size())
                {
                    throw std::runtime_error("Invalid CREATE_BLOCK: not enough bytes for operands");
                }

                // Read block parameters (little-endian)
                uint32_t bytecodeSize = static_cast<uint32_t>(bytecodes[ip]) |
                                        (static_cast<uint32_t>(bytecodes[ip + 1]) << 8) |
                                        (static_cast<uint32_t>(bytecodes[ip + 2]) << 16) |
                                        (static_cast<uint32_t>(bytecodes[ip + 3]) << 24);
                ip += 4;
                uint32_t literalCount = static_cast<uint32_t>(bytecodes[ip]) |
                                        (static_cast<uint32_t>(bytecodes[ip + 1]) << 8) |
                                        (static_cast<uint32_t>(bytecodes[ip + 2]) << 16) |
                                        (static_cast<uint32_t>(bytecodes[ip + 3]) << 24);
                ip += 4;
                uint32_t tempVarCount = static_cast<uint32_t>(bytecodes[ip]) |
                                        (static_cast<uint32_t>(bytecodes[ip + 1]) << 8) |
                                        (static_cast<uint32_t>(bytecodes[ip + 2]) << 16) |
                                        (static_cast<uint32_t>(bytecodes[ip + 3]) << 24);
                ip += 4;

                // Execute CREATE_BLOCK handler using context-based stack
                handleCreateBlock(bytecodeSize, literalCount, tempVarCount);
                break;
            }

            case Bytecode::PUSH_TEMPORARY_VARIABLE:
            {
                ip++; // Skip opcode
                if (ip + 3 >= bytecodes.size())
                {
                    throw std::runtime_error("Invalid PUSH_TEMPORARY_VARIABLE: not enough bytes for operand");
                }

                uint32_t tempIndex =
                    static_cast<uint32_t>(bytecodes[ip]) |
                    (static_cast<uint32_t>(bytecodes[ip + 1]) << 8) |
                    (static_cast<uint32_t>(bytecodes[ip + 2]) << 16) |
                    (static_cast<uint32_t>(bytecodes[ip + 3]) << 24);
                ip += 4;

                // Use context-based temporary variable access
                handlePushTemporaryVariable(tempIndex);
                break;
            }

            case Bytecode::STORE_TEMPORARY_VARIABLE:
            {
                ip++; // Skip opcode
                if (ip + 3 >= bytecodes.size())
                {
                    throw std::runtime_error("Invalid STORE_TEMPORARY_VARIABLE: not enough bytes for operand");
                }

                uint32_t tempIndex =
                    static_cast<uint32_t>(bytecodes[ip]) |
                    (static_cast<uint32_t>(bytecodes[ip + 1]) << 8) |
                    (static_cast<uint32_t>(bytecodes[ip + 2]) << 16) |
                    (static_cast<uint32_t>(bytecodes[ip + 3]) << 24);
                ip += 4;

                // Use context-based temporary variable storage
                handleStoreTemporaryVariable(tempIndex);
                break;
            }

            case Bytecode::POP:
            {
                // Use context-based pop
                handlePop();
                break;
            }

            case Bytecode::DUPLICATE:
            {
                // Use context-based duplicate
                handleDuplicate();
                break;
            }

            case Bytecode::RETURN_STACK_TOP:
            {
                // Pop the return value from context-based stack
                TaggedValue returnValue = pop();

                // Restore the previous context
                activeContext = savedContext;

                // Return the TaggedValue directly
                return returnValue;
            }

            default:
                throw std::runtime_error("Unknown bytecode: " + std::to_string(static_cast<int>(instruction)));
            }
        }

        // If we reach here without explicit return, restore context and return nil
        activeContext = savedContext;
        return TaggedValue::nil();
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

    void Interpreter::push(TaggedValue value)
    {
        if (activeContext == nullptr)
        {
            throw std::runtime_error("No active context for push operation");
        }

        // Get current stack pointer as TaggedValue*
        TaggedValue *currentSP = activeContext->stackPointer;
        
        // Calculate stack bounds - use fixed size for now
        char *contextEnd = reinterpret_cast<char *>(activeContext) + sizeof(MethodContext);
        TaggedValue *stackStart = reinterpret_cast<TaggedValue *>(contextEnd);
        TaggedValue *stackEnd = stackStart + 16; // Fixed stack size

        // Check for stack overflow
        if (currentSP >= stackEnd)
        {
            throw std::runtime_error("Stack overflow");
        }

        // Push value and update stack pointer
        *currentSP = value;
        activeContext->stackPointer = currentSP + 1;
    }
    


    TaggedValue Interpreter::pop()
    {
        if (activeContext == nullptr)
        {
            throw std::runtime_error("No active context for pop operation");
        }

        // Get current stack pointer as TaggedValue*
        TaggedValue *currentSP = activeContext->stackPointer;
        
        // Calculate stack bounds
        char *contextEnd = reinterpret_cast<char *>(activeContext) + sizeof(MethodContext);
        TaggedValue *stackStart = reinterpret_cast<TaggedValue *>(contextEnd);

        // Check for stack underflow
        if (currentSP <= stackStart)
        {
            throw std::runtime_error("Stack underflow");
        }

        // Move stack pointer back and get value
        TaggedValue *newStackPos = currentSP - 1;
        TaggedValue value = *newStackPos;
        activeContext->stackPointer = newStackPos;
        return value;
    }
    


    TaggedValue Interpreter::top()
    {
        if (activeContext == nullptr)
        {
            throw std::runtime_error("No active context for top operation");
        }

        // Get current stack pointer as TaggedValue*
        TaggedValue *currentSP = activeContext->stackPointer;
        
        // Calculate stack bounds
        char *contextEnd = reinterpret_cast<char *>(activeContext) + sizeof(MethodContext);
        TaggedValue *stackStart = reinterpret_cast<TaggedValue *>(contextEnd);

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
        TaggedValue literal = TaggedValue::nil(); // Create a nil value
        push(literal);
    }

    void Interpreter::handlePushInstanceVariable(uint32_t offset)
    {
        // Get the instance variable from the receiver at the given offset
        // For now, just push a new object
        (void)offset; // Suppress unused parameter warning
        TaggedValue value = TaggedValue::nil();
        push(value);
    }

    void Interpreter::handlePushTemporaryVariable(uint32_t offset)
    {
        // Get the temporary variable at the given offset
        TaggedValue *slots = reinterpret_cast<TaggedValue *>(reinterpret_cast<char *>(activeContext) + sizeof(MethodContext));
        push(slots[offset]);
    }

    void Interpreter::handlePushSelf()
    {
        // Push the receiver onto the stack
        push(TaggedValue(activeContext->self));
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
        TaggedValue value = pop();
        TaggedValue *slots = reinterpret_cast<TaggedValue *>(reinterpret_cast<char *>(activeContext) + sizeof(MethodContext));
        slots[offset] = value;
        push(value); // Leave the value on the stack
    }

    void Interpreter::handleSendMessage(uint32_t selectorIndex, uint32_t argCount)
    {
        // Not fully implemented in this basic version
        // In a real implementation, this would look up the method and execute it
        (void)selectorIndex; // Suppress unused parameter warning

        // Pop arguments
        std::vector<TaggedValue> args;
        args.reserve(argCount);
        for (uint32_t i = 0; i < argCount; i++)
        {
            args.push_back(pop());
        }

        // Pop receiver
        TaggedValue receiver = pop();

        // Create a selector string (simplified for now)
        std::string selectorString = "unknownSelector";

        // Send the message using TaggedValue version
        TaggedValue result = sendMessage(receiver, selectorString, args);

        // Push the result
        push(result);
    }

    void Interpreter::handleReturnStackTop()
    {
        // Return from the current context
        TaggedValue result = top();

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
        TaggedValue condition = pop();

        // Check if the condition is true (simplified)
        bool isTrue = (!condition.isNil() && !condition.isFalse());

        // Jump if true
        if (isTrue)
        {
            activeContext->instructionPointer = target;
        }
    }

    void Interpreter::handleJumpIfFalse(uint32_t target)
    {
        // Pop the condition
        TaggedValue condition = pop();

        // Check if the condition is false (simplified)
        bool isFalse = (condition.isNil() || condition.isFalse());

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
        TaggedValue value = top();
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

        // Set the block's class to Block class
        Class *blockClass = ClassRegistry::getInstance().getClass("Block");
        if (blockClass != nullptr)
        {
            blockContext->setClass(blockClass);
        }

        // Push the block context onto the stack
        push(TaggedValue(blockContext));
    }

    void Interpreter::handleExecuteBlock(uint32_t argCount)
    {
        // For now, we only support `value` with no arguments
        (void)argCount;

        // Pop the block context from the stack
        TaggedValue blockValue = pop();
        if (!blockValue.isPointer())
        {
            throw std::runtime_error("Block value is not a pointer for EXECUTE_BLOCK");
        }
        Object *blockObj = blockValue.asObject();
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
