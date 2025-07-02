#include "simple_compiler.h"
#include "symbol.h"

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
    } else if (const auto* block = dynamic_cast<const BlockNode*>(&node)) {
        compileBlock(*block, method);
    } else if (const auto* sequence = dynamic_cast<const SequenceNode*>(&node)) {
        compileSequence(*sequence, method);
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
    
    // Get the selector string for this operator
    std::string selectorString = getSelectorForOperator(node.getOperator());
    
    // Create a symbol for the selector and add it to the literals table
    Symbol* selectorSymbol = Symbol::intern(selectorString);
    uint32_t selectorIndex = method.addLiteral(TaggedValue(selectorSymbol));
    
    // Generate SEND_MESSAGE bytecode with proper selector index
    method.addBytecode(static_cast<uint8_t>(Bytecode::SEND_MESSAGE));
    method.addOperand(selectorIndex); // selector index from literal table
    method.addOperand(1);             // argument count
    
    switch (node.getOperator()) {
        case BinaryOpNode::Operator::Add:
        case BinaryOpNode::Operator::Subtract:
        case BinaryOpNode::Operator::Multiply:
        case BinaryOpNode::Operator::Divide:
        case BinaryOpNode::Operator::LessThan:
        case BinaryOpNode::Operator::GreaterThan:
        case BinaryOpNode::Operator::Equal:
        case BinaryOpNode::Operator::NotEqual:
        case BinaryOpNode::Operator::LessThanOrEqual:
        case BinaryOpNode::Operator::GreaterThanOrEqual:
            // All handled by the common code above
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
        case BinaryOpNode::Operator::NotEqual:  return "~=";
        case BinaryOpNode::Operator::LessThanOrEqual: return "<=";
        case BinaryOpNode::Operator::GreaterThanOrEqual: return ">=";
        default:                                return "unknown";
    }
}

void SimpleCompiler::compileBlock(const BlockNode& node, CompiledMethod& method) {
    // Simplified block compilation for now
    // Just generate CREATE_BLOCK with placeholder values
    // The interpreter will need to handle block creation
    
    // For demonstration, use a simple approach:
    // Store the block's AST representation or compile it inline
    
    // Generate CREATE_BLOCK bytecode with minimal parameters
    method.addBytecode(static_cast<uint8_t>(Bytecode::CREATE_BLOCK));
    method.addOperand(0);  // bytecode size (placeholder)
    method.addOperand(0);  // literal count (placeholder)  
    method.addOperand(0);  // temp var count (placeholder)
    
    // TODO: Properly compile block body and store method reference
    // For now, this creates an empty block that can be executed
}

void SimpleCompiler::compileSequence(const SequenceNode& node, CompiledMethod& method) {
    const auto& statements = node.getStatements();
    
    // Compile each statement
    for (size_t i = 0; i < statements.size(); i++) {
        compileNode(*statements[i], method);
        
        // Pop intermediate results except for the last statement
        // The last statement's result should remain on the stack as the sequence result
        if (i < statements.size() - 1) {
            method.addBytecode(static_cast<uint8_t>(Bytecode::POP));
        }
    }
}

} // namespace smalltalk
