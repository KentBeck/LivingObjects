#include "simple_compiler.h"
#include "symbol.h"

#include <stdexcept>

namespace smalltalk
{

    std::unique_ptr<CompiledMethod> SimpleCompiler::compile(const MethodNode &method)
    {
        auto compiledMethod = std::make_unique<CompiledMethod>();

        // Set up temporary variables (may be empty)
        tempVars_ = method.getTempVars();

        // Compile the method body
        compileNode(*method.getBody(), *compiledMethod);

        // Add return instruction
        compiledMethod->addBytecode(static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP));

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
        else if (const auto *assignment = dynamic_cast<const AssignmentNode *>(&node))
        {
            compileAssignment(*assignment, method);
        }
        else
        {
            throw std::runtime_error("Unknown AST node type");
        }
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
        // Simplified block compilation for now
        // Just generate CREATE_BLOCK with placeholder values
        // The interpreter will need to handle block creation

        // For demonstration, use a simple approach:
        // Store the block's AST representation or compile it inline

        // Generate CREATE_BLOCK bytecode with minimal parameters
        method.addBytecode(static_cast<uint8_t>(Bytecode::CREATE_BLOCK));
        method.addOperand(0); // bytecode size (placeholder)
        method.addOperand(0); // literal count (placeholder)
        method.addOperand(0); // temp var count (placeholder)

        // TODO: Properly compile block body and store method reference
        // For now, this creates an empty block that can be executed
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

        // Variable not found - this is an error
        throw std::runtime_error("Undefined variable: " + varName);
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
                // Generate STORE_TEMPORARY_VARIABLE bytecode
                method.addBytecode(static_cast<uint8_t>(Bytecode::STORE_TEMPORARY_VARIABLE));
                method.addOperand(static_cast<uint32_t>(i));
                return;
            }
        }

        // Variable not found - this is an error
        throw std::runtime_error("Cannot assign to undefined variable: " + varName);
    }

} // namespace smalltalk
