#include "simple_compiler.h"

#include <stdexcept>

namespace smalltalk {

std::unique_ptr<CompiledMethod> SimpleCompiler::compile(const MethodNode& method) {
    auto compiledMethod = std::make_unique<CompiledMethod>();
    
    // Compile the method body
    compileNode(*method.getBody(), *compiledMethod);
    
    // Add return instruction
    compiledMethod->addBytecode(static_cast<uint8_t>(Bytecode::RETURN_STACK_TOP));
    
    return compiledMethod;
}

void SimpleCompiler::compileNode(const ASTNode& node, CompiledMethod& method) {
    // Use dynamic_cast to determine node type
    if (const auto* literal = dynamic_cast<const LiteralNode*>(&node)) {
        compileLiteral(*literal, method);
    } else if (const auto* binOp = dynamic_cast<const BinaryOpNode*>(&node)) {
        compileBinaryOp(*binOp, method);
    } else {
        throw std::runtime_error("Unknown AST node type");
    }
}

void SimpleCompiler::compileLiteral(const LiteralNode& node, CompiledMethod& method) {
    // Add the literal value to the method's literal table
    uint32_t literalIndex = method.addLiteral(node.getValue());
    
    // Generate PUSH_LITERAL bytecode
    method.addBytecode(static_cast<uint8_t>(Bytecode::PUSH_LITERAL));
    method.addOperand(literalIndex);
}

void SimpleCompiler::compileBinaryOp(const BinaryOpNode& node, CompiledMethod& method) {
    // Compile receiver (left operand)
    compileNode(*node.getLeft(), method);
    
    // Compile argument (right operand)  
    compileNode(*node.getRight(), method);
    
    // For now, we'll use primitive operations instead of message sends
    // In a full Smalltalk implementation, these would be message sends like "+" to the receiver
    
    // Generate inline arithmetic instead of message send for simplicity
    switch (node.getOperator()) {
        case BinaryOpNode::Operator::Add:
            // In a real VM, we'd have primitive add instruction or send "+" message
            // For now, we'll implement this in the interpreter directly
            method.addBytecode(static_cast<uint8_t>(Bytecode::SEND_MESSAGE));
            
            // Add "+" selector as literal and get its index
            // For simplicity, we'll use a special literal index 999 to mean "add"
            method.addOperand(999); // selector index (special value for "+")
            method.addOperand(1);   // argument count
            break;
            
        case BinaryOpNode::Operator::Subtract:
            method.addBytecode(static_cast<uint8_t>(Bytecode::SEND_MESSAGE));
            method.addOperand(998); // special value for "-"
            method.addOperand(1);
            break;
            
        case BinaryOpNode::Operator::Multiply:
            method.addBytecode(static_cast<uint8_t>(Bytecode::SEND_MESSAGE));
            method.addOperand(997); // special value for "*"
            method.addOperand(1);
            break;
            
        case BinaryOpNode::Operator::Divide:
            method.addBytecode(static_cast<uint8_t>(Bytecode::SEND_MESSAGE));
            method.addOperand(996); // special value for "/"
            method.addOperand(1);
            break;
            
        default:
            throw std::runtime_error("Unsupported binary operator");
    }
}

std::string SimpleCompiler::getSelectorForOperator(BinaryOpNode::Operator op) {
    switch (op) {
        case BinaryOpNode::Operator::Add:       return "+";
        case BinaryOpNode::Operator::Subtract:  return "-";
        case BinaryOpNode::Operator::Multiply:  return "*";
        case BinaryOpNode::Operator::Divide:    return "/";
        case BinaryOpNode::Operator::LessThan:  return "<";
        case BinaryOpNode::Operator::GreaterThan: return ">";
        case BinaryOpNode::Operator::Equal:     return "=";
        default:                                return "unknown";
    }
}

} // namespace smalltalk
