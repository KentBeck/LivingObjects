#include "interpreter.h"
#include "smalltalk_vm.h"
#include "primitives.h"
#include "smalltalk_class.h"
#include "symbol.h"
#include "simple_parser.h"
#include "smalltalk_image.h"
#include "smalltalk_exception.h"
#include "logger.h"
#include "vm_debugger.h"

#include <cstring>
#include <stdexcept>
#include <vector>
#include <iostream>
#include <algorithm>

namespace smalltalk
{

    Interpreter::Interpreter(MemoryManager &memory, SmalltalkImage &image)
        : memoryManager(memory), image(image)
    {
        // Ensure VM is initialized before any operations
        if (!SmalltalkVM::isInitialized())
        {
            SmalltalkVM::initialize();
        }
        // Initialize the stack chunk
        currentChunk = memoryManager.allocateStackChunk(1024);
    }

    Object *Interpreter::executeMethod(CompiledMethod *method, Object *receiver, std::vector<Object *> &args)
    {
        // Save current context
        MethodContext *previousContext = activeContext;

        // Create a new method context - convert Object* to TaggedValue
        TaggedValue receiverValue = TaggedValue::fromObject(receiver);
        TaggedValue senderValue = previousContext ? TaggedValue::fromObject(previousContext) : TaggedValue::nil();
        MethodContext *context = memoryManager.allocateMethodContext(10 + args.size(), method->getHash(), receiverValue, senderValue, method);

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

        // Execute the context using direct method execution (fixes architectural issue)
        TaggedValue result = executeMethodContext(context, method);

        // Convert TaggedValue result to Object* for legacy compatibility
        Object *resultObj;
        if (result.isPointer())
        {
            resultObj = result.asObject();
        }
        else if (result.isInteger())
        {
            resultObj = memoryManager.allocateInteger(result.asInteger());
        }
        else if (result.isBoolean())
        {
            resultObj = memoryManager.allocateBoolean(result.asBoolean());
        }
        else if (result.isNil())
        {
            // Nil is a special object, not an immediate value that needs boxing
            resultObj = ClassRegistry::getInstance().getClass("UndefinedObject");
        }
        else
        {
            // Should not happen for now, but handle other immediate types if they arise
            throw std::runtime_error("Unhandled immediate value type in executeMethod");
        }

        // Restore previous context
        activeContext = previousContext;

        return resultObj;
    }

    TaggedValue Interpreter::executeCompiledMethod(const CompiledMethod &method)
    {
        // Create a method context for execution
        Object *self = memoryManager.allocateObject(ObjectType::OBJECT, 0); // Simple self object
        TaggedValue selfValue = TaggedValue::fromObject(self);
        TaggedValue senderValue = TaggedValue::nil(); // No sender
        MethodContext *methodContext = memoryManager.allocateMethodContext(
            16,                                   // context size (enough for stack and temporaries)
            method.getHash(),                     // method reference
            selfValue,                            // self as TaggedValue
            senderValue,                          // sender as TaggedValue
            const_cast<CompiledMethod *>(&method) // compiled method pointer
        );

        // Initialize the stack pointer properly - start after temporary variables
        char *contextEnd = reinterpret_cast<char *>(methodContext) + sizeof(MethodContext);
        TaggedValue *slots = reinterpret_cast<TaggedValue *>(contextEnd);

        // Get the number of temporary variables to ensure stack starts after them
        size_t tempVarCount = method.getTempVars().size();

        // Initialize temporary variables to nil
        for (size_t i = 0; i < tempVarCount; i++)
        {
            slots[i] = TaggedValue::nil();
        }

        // Set stack pointer to start after temporary variables
        methodContext->stackPointer = slots + tempVarCount;

        // Use the new architectural approach - direct method execution with currentMethod set
        return executeMethodContext(methodContext, const_cast<CompiledMethod *>(&method));
    }

    TaggedValue Interpreter::executeCompiledMethod(const CompiledMethod &method, MethodContext *context)
    {
        // Execute the method context directly with the provided context and method
        return executeMethodContext(context, const_cast<CompiledMethod *>(&method));
    }

    TaggedValue Interpreter::executeMethodContext(MethodContext *context)
    {
        if (!context->method)
        {
            throw std::runtime_error("No method associated with context");
        }

        // Delegate to the main execution method
        return executeMethodContext(context, context->method);
    }

    TaggedValue Interpreter::executeMethodContext(MethodContext *context, CompiledMethod *method)
    {
        // Set up execution state using context-based stack
        MethodContext *savedContext = activeContext;
        activeContext = context;
        context->method = method;

        // Execute bytecode using the core execution engine
        TaggedValue result = execute();

        // Restore previous context
        activeContext = savedContext;
        return result;
    }

    TaggedValue Interpreter::execute()
    {
        if (!activeContext)
        {
            throw std::runtime_error("No active context for execution");
        }

        if (!activeContext->method)
        {
            throw std::runtime_error("Active context must have a method to execute");
        }

        // Main bytecode execution loop - continue until no more contexts
        while (activeContext != nullptr)
        {
            // Check if we've reached the end of the current method's bytecodes
            // There should be no such thing as an implict return--maybe just from blocks?
            if (activeContext->instructionPointer >= activeContext->method->getBytecodes().size())
            {
                // Implicit return - if stack is empty, return self; otherwise return top of stack
                char *contextEnd = reinterpret_cast<char *>(activeContext) + sizeof(MethodContext);
                TaggedValue *slots = reinterpret_cast<TaggedValue *>(contextEnd);
                TaggedValue *stackStart = slots + activeContext->method->getTempVars().size(); // Stack starts after temp vars
                TaggedValue *currentSP = activeContext->stackPointer;

                if (currentSP <= stackStart)
                {
                    // Stack is empty - return self
                    push(activeContext->self);
                }
                // else: there's already a value on stack to return

                returnStackTop();
                continue;
            }

            uint8_t opcode = activeContext->method->getBytecodes()[activeContext->instructionPointer];
            Bytecode instruction = static_cast<Bytecode>(opcode);

            // Skip opcode - move to first operand byte
            activeContext->instructionPointer++;

            // Debug bytecode tracing
            if (Logger::getInstance().getLevel() <= LogLevel::DEBUG_LEVEL)
            {
                std::vector<TaggedValue> stack;
                // Add current stack contents for debugging (simplified)
                VM_DEBUG_BYTECODE(std::to_string(static_cast<int>(instruction)),
                                  activeContext->instructionPointer, stack);
            }

            switch (instruction)
            {
            case Bytecode::PUSH_LITERAL:
                pushLiteral();
                break;

            case Bytecode::PUSH_SELF:
                pushSelf();
                break;

            case Bytecode::SEND_MESSAGE:
                sendMessageBytecode();
                break;

            case Bytecode::CREATE_BLOCK:
                createBlock();
                break;

            case Bytecode::PUSH_TEMPORARY_VARIABLE:
                pushTemporaryVariable();
                break;

            case Bytecode::STORE_TEMPORARY_VARIABLE:
                storeTemporaryVariable();
                break;

            case Bytecode::POP:
                popStack();
                break;

            case Bytecode::DUPLICATE:
                duplicate();
                break;

            case Bytecode::RETURN_STACK_TOP:
                returnStackTop();
                break;

            default:
                throw std::runtime_error("Unknown bytecode: " + std::to_string(static_cast<int>(instruction)));
            }
        }

        // Execution completed - return the last return value
        return lastReturnValue;
    }

    // Bytecode operation helpers

    void Interpreter::pushLiteral()
    {
        uint32_t literalIndex = readUint32FromBytecode(activeContext->method->getBytecodes(), activeContext);

        if (literalIndex >= activeContext->method->getLiterals().size())
        {
            throw std::runtime_error("Invalid literal index: " + std::to_string(literalIndex));
        }

        // Push TaggedValue directly
        TaggedValue literal = activeContext->method->getLiterals()[literalIndex];
        push(literal);
    }

    void Interpreter::pushSelf()
    {
        push(activeContext->self);
    }

    void Interpreter::sendMessageBytecode()
    {
        if (activeContext->instructionPointer + 7 >= activeContext->method->getBytecodes().size())
        {
            throw std::runtime_error("Invalid SEND_MESSAGE: not enough bytes for operands");
        }

        // Read selector index
        uint32_t selectorIndex = readUint32FromBytecode(activeContext->method->getBytecodes(), activeContext);

        // Read argument count
        uint32_t argCount = readUint32FromBytecode(activeContext->method->getBytecodes(), activeContext);

        if (selectorIndex >= activeContext->method->getLiterals().size())
        {
            throw std::runtime_error("Invalid selector index: " + std::to_string(selectorIndex));
        }

        // Get the selector from literals
        TaggedValue selectorValue = activeContext->method->getLiterals()[selectorIndex];
        if (!selectorValue.isPointer())
        {
            throw std::runtime_error("Selector is not a pointer");
        }

        // Get selector string
        Object *selectorObj = selectorValue.asObject();
        std::string selectorString;
        if (selectorObj->header.getType() == ObjectType::SYMBOL)
        {
            Symbol *symbol = reinterpret_cast<Symbol *>(selectorObj);
            selectorString = symbol->getName();
        }
        else
        {
            throw std::runtime_error("Selector is not a symbol");
        }

        // Pop arguments from stack (in reverse order)
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
    }

    void Interpreter::createBlock()
    {
        if (activeContext->instructionPointer + 7 >= activeContext->method->getBytecodes().size())
        {
            throw std::runtime_error("Invalid CREATE_BLOCK: not enough bytes for operands");
        }

        // Read block parameters (little-endian)
        uint32_t literalIndex = readUint32FromBytecode(activeContext->method->getBytecodes(), activeContext);
        readUint32FromBytecode(activeContext->method->getBytecodes(), activeContext); // Skip parameter count (not used)

        // The block's compiled method should be in the current method's literals
        if (!activeContext || !activeContext->method)
        {
            throw std::runtime_error("No current method for block creation");
        }

        // Get the block's compiled method directly from the literal
        if (literalIndex >= activeContext->method->getLiterals().size())
        {
            throw std::runtime_error("Invalid literal index for block: " + std::to_string(literalIndex));
        }

        TaggedValue blockMethodValue = activeContext->method->getLiterals()[literalIndex];
        if (!blockMethodValue.isPointer())
        {
            throw std::runtime_error("Block method literal is not a pointer");
        }

        // Cast to CompiledMethod
        CompiledMethod *blockMethod = reinterpret_cast<CompiledMethod *>(blockMethodValue.asObject());

        // Create a simple block context that can be used with the Block class
        MethodContext *homeContext = activeContext;
        if (homeContext == nullptr)
        {
            throw std::runtime_error("Cannot create block without an active context");
        }

        // Create block object that stores the compiled method directly
        TaggedValue receiverValue = homeContext->self;
        TaggedValue senderValue = TaggedValue::nil();
        TaggedValue homeValue = TaggedValue::fromObject(homeContext);

        // Allocate space for block context plus one slot for the compiled method
        BlockContext *blockContext = memoryManager.allocateBlockContext(
            8,             // context size (includes space for method pointer)
            0,             // hash not used anymore
            receiverValue, // receiver (same as home context)
            senderValue,   // sender (will be set when block is executed)
            homeValue      // home context
        );

        // Store the compiled method directly in the block context's variable area
        char *contextEnd = reinterpret_cast<char *>(blockContext) + sizeof(BlockContext);
        TaggedValue *slots = reinterpret_cast<TaggedValue *>(contextEnd);
        slots[0] = TaggedValue::fromObject(reinterpret_cast<Object *>(blockMethod));

        // Set the block's class to Block class
        Class *blockClass = ClassRegistry::getInstance().getClass("Block");
        if (blockClass != nullptr)
        {
            blockContext->setClass(blockClass);
        }

        // Push the block context onto the stack
        push(TaggedValue::fromObject(blockContext));
    }

    void Interpreter::pushTemporaryVariable()
    {
        uint32_t tempIndex = readUint32FromBytecode(activeContext->method->getBytecodes(), activeContext);
        TaggedValue *slots = reinterpret_cast<TaggedValue *>(reinterpret_cast<char *>(activeContext) + sizeof(MethodContext));
        push(slots[tempIndex]);
    }

    void Interpreter::storeTemporaryVariable()
    {
        uint32_t tempIndex = readUint32FromBytecode(activeContext->method->getBytecodes(), activeContext);
        TaggedValue value = pop();
        TaggedValue *slots = reinterpret_cast<TaggedValue *>(reinterpret_cast<char *>(activeContext) + sizeof(MethodContext));
        slots[tempIndex] = value;
        push(value); // Leave the value on the stack
    }

    void Interpreter::popStack()
    {
        // Pop the top value from the stack
        pop();
    }

    void Interpreter::duplicate()
    {
        // Duplicate the top value on the stack
        TaggedValue value = top();
        push(value);
    }

    void Interpreter::returnStackTop()
    {
        // Pop the return value from current context's stack
        TaggedValue returnValue = pop();

        // Get the sender context
        if (!activeContext->sender.isPointer())
        {
            // No sender - this is the top-level method, end execution with result
            lastReturnValue = returnValue;
            activeContext = nullptr;
            return;
        }

        // Switch to sender context
        MethodContext *senderContext = static_cast<MethodContext *>(activeContext->sender.asObject());
        activeContext = senderContext;

        // Push return value onto sender's stack
        push(returnValue);
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
            VM_DEBUG_EXCEPTION("MessageSend", "Cannot determine receiver class", selector);
            throw std::runtime_error("Cannot determine receiver class");
        }

        // Debug tracing
        VM_DEBUG_METHOD_ENTRY(selector, receiverClass->getName(), args);

        // Create selector symbol
        Symbol *selectorSymbol = Symbol::intern(selector);

        // Look up method
        std::shared_ptr<CompiledMethod> method = receiverClass->lookupMethod(selectorSymbol);

        if (method)
        {
            if (method->primitiveNumber != 0)
            {
                // Try primitive first
                try
                {
                    LOG_VM_DEBUG("Calling primitive " + std::to_string(method->primitiveNumber) + " for " + selector);
                    TaggedValue result = Primitives::callPrimitive(method->primitiveNumber, receiver, args, *this);
                    VM_DEBUG_METHOD_EXIT(selector, receiverClass->getName(), result);
                    return result;
                }
                catch (const PrimitiveFailure &e)
                {
                    LOG_VM_DEBUG("Primitive " + std::to_string(method->primitiveNumber) + " failed, falling back to Smalltalk code");
                    // Fall back to Smalltalk code below
                }
            }

            // Execute the compiled method (either non-primitive or primitive fallback)
            // Create a new context for the method execution
            MethodContext *newContext = memoryManager.allocateMethodContext(
                16 + method->getTempVars().size(),      // stack size + temp vars
                method->getHash(),                      // method hash (stored in header)
                receiver,                               // self
                TaggedValue::fromObject(activeContext), // sender
                method.get()                            // compiled method pointer
            );

            // Initialize stack pointer and temporary variables
            char *contextEnd = reinterpret_cast<char *>(newContext) + sizeof(MethodContext);
            TaggedValue *slots = reinterpret_cast<TaggedValue *>(contextEnd);

            // Store arguments as temporary variables (method parameters come first)
            // In Smalltalk, method parameters are the first temporary variables
            for (size_t i = 0; i < args.size(); i++)
            {
                slots[i] = args[i]; // Store in correct order, not reversed
            }

            // Initialize remaining temp vars to nil
            for (size_t i = args.size(); i < method->getTempVars().size(); i++)
            {
                slots[i] = TaggedValue::nil();
            }

            // Set stack pointer to start after all temporary variables
            newContext->stackPointer = slots + method->getTempVars().size();

            // Execute the method directly with the compiled method object
            TaggedValue result = executeCompiledMethod(*method, newContext);

            // Debug tracing
            VM_DEBUG_METHOD_EXIT(selector, receiverClass->getName(), result);

            return result;
        }

        // No method found - throw proper exception
        VM_DEBUG_EXCEPTION("MessageNotUnderstood", "Method not found: " + selector, receiverClass->getName());
        // later push a context to execute MNU
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

    bool Interpreter::findExceptionHandler(MethodContext *&handlerContext, int &handlerPC)
    {
        // Walk up the context chain looking for exception handlers
        MethodContext *context = activeContext;

        while (context != nullptr)
        {
            // Get the compiled method directly from the context
            CompiledMethod *method = context->method;

            if (method && method->primitiveNumber == PrimitiveNumbers::EXCEPTION_MARK)
            {
                // This method has an exception handler marker
                // The handler starts right after the primitive failure
                handlerContext = context;
                handlerPC = 0; // Start from beginning since primitive failed
                return true;
            }

            // Move to sender context
            if (context->sender.isPointer())
            {
                context = static_cast<MethodContext *>(context->sender.asObject());
            }
            else
            {
                break;
            }
        }

        return false;
    }

    void Interpreter::unwindToContext(MethodContext *targetContext)
    {
        // Unwind the stack to the target context
        while (activeContext != targetContext && activeContext != nullptr)
        {
            // Get sender before we potentially invalidate the context
            TaggedValue sender = activeContext->sender;

            // TODO: In a full implementation, we would run unwind blocks here
            // For now, just move to the sender

            if (sender.isPointer())
            {
                activeContext = static_cast<MethodContext *>(sender.asObject());
            }
            else
            {
                activeContext = nullptr;
            }
        }

        if (activeContext == nullptr && targetContext != nullptr)
        {
            // We couldn't reach the target context - this is an error
            throw std::runtime_error("Failed to unwind to exception handler context");
        }
    }

    uint32_t Interpreter::readUint32FromBytecode(const std::vector<uint8_t> &bytecodes, MethodContext *context)
    {
        if (context->instructionPointer + 3 >= bytecodes.size())
        {
            throw std::runtime_error("Invalid bytecode: not enough bytes for 32-bit operand");
        }

        uint32_t value = static_cast<uint32_t>(bytecodes[context->instructionPointer]) |
                         (static_cast<uint32_t>(bytecodes[context->instructionPointer + 1]) << 8) |
                         (static_cast<uint32_t>(bytecodes[context->instructionPointer + 2]) << 16) |
                         (static_cast<uint32_t>(bytecodes[context->instructionPointer + 3]) << 24);
        context->instructionPointer += 4;
        return value;
    }

} // namespace smalltalk
