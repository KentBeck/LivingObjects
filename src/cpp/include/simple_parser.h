#pragma once

#include "ast.h"

#include <memory>
#include <string>

namespace smalltalk {

/**
 * Simple recursive descent parser for basic expressions
 * Supports: integers, +, -, *, /, (, )
 * 
 * Grammar:
 *   expression := term (('+' | '-') term)*
 *   term       := factor (('*' | '/') factor)*
 *   factor     := integer | '(' expression ')'
 *   integer    := [0-9]+
 */
class SimpleParser {
public:
    SimpleParser(std::string input);
    
    // Parse the input and return a method AST
    std::unique_ptr<MethodNode> parseMethod();
    
private:
    // Parsing methods
    ASTNodePtr parseExpression();
    ASTNodePtr parseTerm();
    ASTNodePtr parseFactor();
    ASTNodePtr parseInteger();
    
    // Tokenization
    void skipWhitespace();
    char peek() const;
    char consume();
    bool isAtEnd() const;
    bool isDigit(char c) const;
    
    // Error handling
    void error(const std::string& message);
    
    std::string input_;
    size_t pos_ = 0;
};

} // namespace smalltalk
