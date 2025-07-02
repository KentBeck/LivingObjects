#pragma once

#include "ast.h"

#include <memory>
#include <string>

namespace smalltalk {

/**
 * Simple recursive descent parser for basic expressions
 * Supports: integers, booleans, nil, +, -, *, /, <, >, =, ~=, <=, >=, (, )
 * 
 * Grammar:
 *   expression := comparison
 *   comparison := arithmetic (('<' | '>' | '=' | '~=' | '<=' | '>=') arithmetic)*
 *   arithmetic := term (('+' | '-') term)*
 *   term       := factor (('*' | '/') factor)*
 *   factor     := integer | literal | string | '(' expression ')'
 *   integer    := [0-9]+
 *   literal    := 'true' | 'false' | 'nil'
 *   string     := "'" [^']* "'"
 */
class SimpleParser {
public:
    SimpleParser(std::string input);
    
    // Parse the input and return a method AST
    std::unique_ptr<MethodNode> parseMethod();
    
private:
    // Parsing methods
    ASTNodePtr parseExpression();
    ASTNodePtr parseStatements();
    ASTNodePtr parseComparison();
    ASTNodePtr parseArithmetic();
    ASTNodePtr parseTerm();
    ASTNodePtr parseFactor();
    ASTNodePtr parseInteger();
    ASTNodePtr parseIdentifierOrLiteral();
    ASTNodePtr parseString();
    
    // Tokenization
    void skipWhitespace();
    char peek() const;
    char consume();
    bool isAtEnd() const;
    bool isDigit(char c) const;
    bool isAlpha(char c) const;
    
    // Error handling
    void error(const std::string& message);
    
    std::string input_;
    size_t pos_ = 0;
};

} // namespace smalltalk
