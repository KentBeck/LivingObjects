#pragma once

#include "ast.h"
#include "bytecode.h"
#include "compiled_method.h"

#include <memory>

namespace smalltalk {

/**
 * Simple compiler that converts AST to bytecode
 */
class SimpleCompiler {
public:
    SimpleCompiler() = default;
    
    // Compile a method AST to bytecode
    std::unique_ptr<CompiledMethod> compile(const MethodNode& method);
    
private:
    // Compile different node types
    void compileNode(const ASTNode& node, CompiledMethod& method);
    void compileLiteral(const LiteralNode& node, CompiledMethod& method);
    void compileBinaryOp(const BinaryOpNode& node, CompiledMethod& method);
    void compileMessageSend(const MessageSendNode& node, CompiledMethod& method);
    void compileBlock(const BlockNode& node, CompiledMethod& method);
    void compileSequence(const SequenceNode& node, CompiledMethod& method);
    
    // Get the selector name for a binary operator
    std::string getSelectorForOperator(BinaryOpNode::Operator op);
};

} // namespace smalltalk
