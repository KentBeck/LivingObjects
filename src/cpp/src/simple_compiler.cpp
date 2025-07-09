#include "simple_compiler.h"
#include "symbol.h"
#include "smalltalk_exception.h"
#include "smalltalk_class.h"

#include <stdexcept>

namespace smalltalk
{

    std::unique_ptr<CompiledMethod> SimpleCompiler::compile(const MethodNode &method)
    {
        auto compiledMethod = std::make_unique<CompiledMethod>();

        // Set up temporary variables (may be empty)
        tempVars_ = method.getTempVars();
        
        // Add temp vars to the compiled method
        for (const auto& tempVar : tempVars_) {
            compiledMethod->addTempVar(tempVar);
        }

        // Compile the method body
        compileNode(*method.getBody(), *compiledMethod);

        // Add implicit return of last expression (if no explicit return was generated)
        if (compiledMethod->getBytecodes().empty() ||
            compiledMethod->getBytecodes().back() != static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP))
        {
            compiledMethod->addBytecode(static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP));
        }

        return compiledMethod;
    }

    void SimpleCompiler::compileNode(const ASTNode &node, CompiledMethod &method)
    {
        // Use dynamic_cast to determine node type
        if (const auto *literal = dynamic_cast<const LiteralNode *>(&node))
        {
            compileLiteral(*literal, method);
        }

        else if (const auto *messageSend = dynamic_cast<const MessageSendNode *>(&node))
        {
            compileMessageSend(*messageSend, method);
        }
        else if (const auto *block = dynamic_cast<const BlockNode *>(&node))
        {
            compileBlock(*block, method);
        }
        else if (const auto *sequence = dynamic_cast<const SequenceNode *>(&node))
        {
            compileSequence(*sequence, method);
        }
        else if (const auto *variable = dynamic_cast<const VariableNode *>(&node))
        {
            compileVariable(*variable, method);
        }
        else if (const auto *self = dynamic_cast<const SelfNode *>(&node))
        {
            compileSelf(*self, method);
        }
        else if (const auto *assignment = dynamic_cast<const AssignmentNode *>(&node))
        {
            compileAssignment(*assignment, method);
        }
        else if (const auto *ret = dynamic_cast<const ReturnNode *>(&node))
        {
            compileReturn(*ret, method);
        }
        else
        {
            throw std::runtime_error("Unknown AST node type");
        }
    }

    void SimpleCompiler::compileReturn(const ReturnNode &node, CompiledMethod &method)
    {
        compileNode(*node.getValue(), method);
        method.addBytecode(static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP));
    }

    void SimpleCompiler::compileLiteral(const LiteralNode &node, CompiledMethod &method)
    {
        // Add the literal value to the method's literal table
        uint32_t literalIndex = method.addLiteral(node.getValue());

        // Generate PUSH_LITERAL bytecode
        method.addBytecode(static_cast<uint8_t>(Bytecode::PUSH_LITERAL));
        method.addOperand(literalIndex);
    }

    void SimpleCompiler::compileMessageSend(const MessageSendNode &node, CompiledMethod &method)
    {
        // Compile receiver
        compileNode(*node.getReceiver(), method);

        // Compile arguments
        const auto &arguments = node.getArguments();
        for (const auto &arg : arguments)
        {
            compileNode(*arg, method);
        }

        // Create a symbol for the selector and add it to the literals table
        Symbol *selectorSymbol = Symbol::intern(node.getSelector());
        uint32_t selectorIndex = method.addLiteral(TaggedValue(selectorSymbol));

        // Generate SEND_MESSAGE bytecode
        method.addBytecode(static_cast<uint8_t>(Bytecode::SEND_MESSAGE));
        method.addOperand(selectorIndex);                           // selector index from literal table
        method.addOperand(static_cast<uint32_t>(arguments.size())); // argument count
    }

    void SimpleCompiler::compileBlock(const BlockNode &node, CompiledMethod &method)
    {
        // Create a separate compiled method for the block on the heap
        auto blockMethod = std::make_unique<CompiledMethod>();

        // Create a separate compiler instance for the block with its own temp var context
        SimpleCompiler blockCompiler;

        // LEXICAL SCOPING: Include outer method's temporary variables first
        // This allows blocks to access variables from their lexical environment
        blockCompiler.tempVars_ = tempVars_;  // Copy outer method's temp vars
        
        // Add block parameters to the compiler's temp var list
        for (const auto &param : node.getParameters())
        {
            blockCompiler.tempVars_.push_back(param);
        }
        
        // Add block temporaries to the compiler's temp var list
        for (const auto &temp : node.getTemporaries())
        {
            blockCompiler.tempVars_.push_back(temp);
        }

        // LEXICAL SCOPING: Add outer method's temporary variables first
        // This allows blocks to access variables from their lexical environment
        for (const auto &outerVar : tempVars_)
        {
            blockMethod->addTempVar(outerVar);
        }
        
        // Store the count of home variables for runtime lexical scoping
        blockMethod->homeVarCount = tempVars_.size();
        
        // Add block parameters as temporary variables to the block method
        for (const auto &param : node.getParameters())
        {
            blockMethod->addTempVar(param);
        }
        
        // Add block temporaries as temporary variables to the block method
        for (const auto &temp : node.getTemporaries())
        {
            blockMethod->addTempVar(temp);
        }

        // Compile the block body using the block compiler
        blockCompiler.compileNode(*node.getBody(), *blockMethod);

        // Add return instruction if not already present
        if (blockMethod->getBytecodes().empty() ||
            blockMethod->getBytecodes().back() != static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP))
        {
            blockMethod->addBytecode(static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP));
        }

        // Store the block method as a literal in the main method
        // The pointer will be valid as long as the parent method exists  
        CompiledMethod* blockMethodPtr = blockMethod.release();
        
        // IMPORTANT: We also need to add the block method to the image so it can be found during execution
        // For now, we'll store the pointer directly as a literal
        // TODO: Add proper block method registration to image
        int blockMethodIndex = method.addLiteral(TaggedValue::fromObject(reinterpret_cast<Object*>(blockMethodPtr)));

        // Generate CREATE_BLOCK bytecode with the literal index
        method.addBytecode(static_cast<uint8_t>(Bytecode::CREATE_BLOCK));
        method.addOperand(blockMethodIndex); // index of the block method in literals
        method.addOperand(static_cast<uint32_t>(node.getParameters().size())); // parameter count
    }

    void SimpleCompiler::compileSequence(const SequenceNode &node, CompiledMethod &method)
    {
        const auto &statements = node.getStatements();

        // Compile each statement
        for (size_t i = 0; i < statements.size(); i++)
        {
            compileNode(*statements[i], method);

            // Pop intermediate results except for the last statement
            // The last statement's result should remain on the stack as the sequence result
            if (i < statements.size() - 1)
            {
                method.addBytecode(static_cast<uint8_t>(Bytecode::POP));
            }
        }
    }

    void SimpleCompiler::compileVariable(const VariableNode &node, CompiledMethod &method)
    {
        const std::string &varName = node.getName();

        // Find the variable in temporary variables
        for (size_t i = 0; i < tempVars_.size(); i++)
        {
            if (tempVars_[i] == varName)
            {
                // Generate PUSH_TEMPORARY_VARIABLE bytecode
                method.addBytecode(static_cast<uint8_t>(Bytecode::PUSH_TEMPORARY_VARIABLE));
                method.addOperand(static_cast<uint32_t>(i));
                return;
            }
        }

        // Check if it's a global class name
        Class* globalClass = ClassRegistry::getInstance().getClass(varName);
        if (globalClass != nullptr)
        {
            // Add the class as a literal and push it
            uint32_t literalIndex = method.addLiteral(TaggedValue::fromObject(globalClass));
            method.addBytecode(static_cast<uint8_t>(Bytecode::PUSH_LITERAL));
            method.addOperand(literalIndex);
            return;
        }

        // Variable not found - this is an error
        // Throw proper Smalltalk exception
        ExceptionHandler::throwException(std::make_unique<NameError>(varName));
    }

    void SimpleCompiler::compileSelf(const SelfNode &node, CompiledMethod &method)
    {
        // Generate PUSH_SELF bytecode
        method.addBytecode(static_cast<uint8_t>(Bytecode::PUSH_SELF));
    }

    void SimpleCompiler::compileAssignment(const AssignmentNode &node, CompiledMethod &method)
    {
        const std::string &varName = node.getVariable();

        // Compile the value expression first
        compileNode(*node.getValue(), method);

        // Find the variable in temporary variables
        for (size_t i = 0; i < tempVars_.size(); i++)
        {
            if (tempVars_[i] == varName)
            {
                // Duplicate the value so assignment returns the assigned value
                method.addBytecode(static_cast<uint8_t>(Bytecode::DUPLICATE));

                // Generate STORE_TEMPORARY_VARIABLE bytecode (this will pop one copy)
                method.addBytecode(static_cast<uint8_t>(Bytecode::STORE_TEMPORARY_VARIABLE));
                method.addOperand(static_cast<uint32_t>(i));
                return;
            }
        }

        // Variable not found - this is an error
        // Throw proper Smalltalk exception
        ExceptionHandler::throwException(std::make_unique<NameError>(varName));
    }

} // namespace smalltalk
