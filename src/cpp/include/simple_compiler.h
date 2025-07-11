#pragma once

#include "ast.h"
#include "bytecode.h"
#include "compiled_method.h"

#include <memory>

namespace smalltalk
{

    /**
     * Simple compiler that converts AST to bytecode
     */
    class SimpleCompiler
    {
    public:
        SimpleCompiler() = default;

        // Compile a method AST to bytecode
        std::unique_ptr<CompiledMethod> compile(const MethodNode &method);

    private:
        // Compile different node types
        void compileNode(const ASTNode &node, CompiledMethod &method);
        void compileLiteral(const LiteralNode &node, CompiledMethod &method);
        void compileMessageSend(const MessageSendNode &node, CompiledMethod &method);
        void compileBlock(const BlockNode &node, CompiledMethod &method);
        void compileSequence(const SequenceNode &node, CompiledMethod &method);
        void compileVariable(const VariableNode &node, CompiledMethod &method);
        void compileSelf(const SelfNode &node, CompiledMethod &method);
        void compileAssignment(const AssignmentNode &node, CompiledMethod &method);
        void compileReturn(const ReturnNode &node, CompiledMethod &method);
        void compileArrayLiteral(const ArrayLiteralNode &node, CompiledMethod &method);

        // Temporary variables for the current method being compiled
        std::vector<std::string> tempVars_;
    };

} // namespace smalltalk
