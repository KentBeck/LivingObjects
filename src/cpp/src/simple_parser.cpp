#include "simple_parser.h"
#include <stdexcept>
#include <cctype>
#include <iostream>

namespace smalltalk {

SimpleParser::SimpleParser(const std::string& input) 
    : input_(input), pos_(0) {
}

std::unique_ptr<MethodNode> SimpleParser::parseMethod() {
    skipWhitespace();
    auto body = parseExpression();
    skipWhitespace();
    
    if (!isAtEnd()) {
        error("Unexpected characters at end of input");
    }
    
    return std::make_unique<MethodNode>(std::move(body));
}

ASTNodePtr SimpleParser::parseExpression() {
    auto left = parseTerm();
    
    while (!isAtEnd()) {
        skipWhitespace();
        char op = peek();
        
        if (op == '+' || op == '-') {
            consume(); // consume operator
            auto right = parseTerm();
            
            BinaryOpNode::Operator binOp = (op == '+') ? 
                BinaryOpNode::Operator::Add : BinaryOpNode::Operator::Subtract;
            
            left = std::make_unique<BinaryOpNode>(std::move(left), binOp, std::move(right));
        } else {
            break;
        }
    }
    
    return left;
}

ASTNodePtr SimpleParser::parseTerm() {
    auto left = parseFactor();
    
    while (!isAtEnd()) {
        skipWhitespace();
        char op = peek();
        
        if (op == '*' || op == '/') {
            consume(); // consume operator
            auto right = parseFactor();
            
            BinaryOpNode::Operator binOp = (op == '*') ? 
                BinaryOpNode::Operator::Multiply : BinaryOpNode::Operator::Divide;
            
            left = std::make_unique<BinaryOpNode>(std::move(left), binOp, std::move(right));
        } else {
            break;
        }
    }
    
    return left;
}

ASTNodePtr SimpleParser::parseFactor() {
    skipWhitespace();
    
    if (peek() == '(') {
        consume(); // consume '('
        auto expr = parseExpression();
        skipWhitespace();
        
        if (peek() != ')') {
            error("Expected ')' after expression");
        }
        consume(); // consume ')'
        
        return expr;
    } else if (isDigit(peek())) {
        return parseInteger();
    } else {
        error(std::string("Unexpected character: ") + peek());
        return nullptr; // Never reached
    }
}

ASTNodePtr SimpleParser::parseInteger() {
    std::string numStr;
    
    while (!isAtEnd() && isDigit(peek())) {
        numStr += consume();
    }
    
    if (numStr.empty()) {
        error("Expected integer");
    }
    
    try {
        int value = std::stoi(numStr);
        return std::make_unique<LiteralNode>(TaggedValue(static_cast<int32_t>(value)));
    } catch (const std::exception& e) {
        error("Invalid integer: " + numStr);
        return nullptr; // Never reached
    }
}

void SimpleParser::skipWhitespace() {
    while (!isAtEnd() && (std::isspace(peek()) != 0)) {
        consume();
    }
}

char SimpleParser::peek() const {
    if (isAtEnd()) {
        return '\0';
    }
    return input_[pos_];
}

char SimpleParser::consume() {
    if (isAtEnd()) {
        error("Unexpected end of input");
    }
    return input_[pos_++];
}

bool SimpleParser::isAtEnd() const {
    return pos_ >= input_.size();
}

bool SimpleParser::isDigit(char c) const {
    return c >= '0' && c <= '9';
}

void SimpleParser::error(const std::string& message) {
    throw std::runtime_error("Parse error at position " + std::to_string(pos_) + ": " + message);
}

} // namespace smalltalk
