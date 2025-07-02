#include "simple_parser.h"
#include "smalltalk_string.h"

#include <cctype>
#include <iostream>
#include <stdexcept>

namespace smalltalk {

SimpleParser::SimpleParser(std::string input) 
    : input_(std::move(input)) {
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
    return parseComparison();
}

ASTNodePtr SimpleParser::parseComparison() {
    auto left = parseArithmetic();
    
    while (!isAtEnd()) {
        skipWhitespace();
        
        // Check for comparison operators
        BinaryOpNode::Operator compOp;
        bool foundComparison = false;
        
        if (peek() == '<') {
            consume();
            if (peek() == '=') {
                consume();
                compOp = BinaryOpNode::Operator::LessThanOrEqual;
            } else {
                compOp = BinaryOpNode::Operator::LessThan;
            }
            foundComparison = true;
        } else if (peek() == '>') {
            consume();
            if (peek() == '=') {
                consume();
                compOp = BinaryOpNode::Operator::GreaterThanOrEqual;
            } else {
                compOp = BinaryOpNode::Operator::GreaterThan;
            }
            foundComparison = true;
        } else if (peek() == '=') {
            consume();
            compOp = BinaryOpNode::Operator::Equal;
            foundComparison = true;
        } else if (peek() == '~') {
            consume();
            if (peek() == '=') {
                consume();
                compOp = BinaryOpNode::Operator::NotEqual;
                foundComparison = true;
            } else {
                error("Expected '=' after '~'");
            }
        }
        
        if (foundComparison) {
            auto right = parseArithmetic();
            left = std::make_unique<BinaryOpNode>(std::move(left), compOp, std::move(right));
        } else {
            break;
        }
    }
    
    return left;
}

ASTNodePtr SimpleParser::parseArithmetic() {
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
    } else if (peek() == '[') {
        consume(); // consume '['
        auto expr = parseExpression();
        skipWhitespace();
        
        if (peek() != ']') {
            error("Expected ']' after block expression");
        }
        consume(); // consume ']'
        
        return std::make_unique<BlockNode>(std::move(expr));
    } else if (isDigit(peek())) {
        return parseInteger();
    } else if (isAlpha(peek())) {
        return parseIdentifierOrLiteral();
    } else if (peek() == '\'') {
        return parseString();
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

ASTNodePtr SimpleParser::parseIdentifierOrLiteral() {
    std::string identifier;
    
    while (!isAtEnd() && (isAlpha(peek()) || isDigit(peek()))) {
        identifier += consume();
    }
    
    if (identifier.empty()) {
        error("Expected identifier");
    }
    
    // Check for boolean and nil literals
    if (identifier == "true") {
        return std::make_unique<LiteralNode>(TaggedValue::trueValue());
    } else if (identifier == "false") {
        return std::make_unique<LiteralNode>(TaggedValue::falseValue());
    } else if (identifier == "nil") {
        return std::make_unique<LiteralNode>(TaggedValue::nil());
    } else {
        error("Unknown identifier: " + identifier);
        return nullptr; // Never reached
    }
}

ASTNodePtr SimpleParser::parseString() {
    if (peek() != '\'') {
        error("Expected string to start with '");
    }
    consume(); // consume opening '
    
    std::string content;
    
    while (!isAtEnd() && peek() != '\'') {
        if (peek() == '\\') {
            // Handle escape sequences
            consume(); // consume backslash
            if (isAtEnd()) {
                error("Unexpected end of input in string literal");
            }
            
            char escaped = consume();
            switch (escaped) {
                case 'n': content += '\n'; break;
                case 't': content += '\t'; break;
                case 'r': content += '\r'; break;
                case '\\': content += '\\'; break;
                case '\'': content += '\''; break;
                default:
                    content += escaped; // Keep unknown escapes as-is
                    break;
            }
        } else {
            content += consume();
        }
    }
    
    if (isAtEnd() || peek() != '\'') {
        error("Unterminated string literal");
    }
    consume(); // consume closing '
    
    // Create a String object and wrap it in a TaggedValue
    TaggedValue stringValue = StringUtils::createTaggedString(content);
    return std::make_unique<LiteralNode>(stringValue);
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

bool SimpleParser::isAlpha(char c) const {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z');
}

void SimpleParser::error(const std::string& message) {
    throw std::runtime_error("Parse error at position " + std::to_string(pos_) + ": " + message);
}

} // namespace smalltalk
