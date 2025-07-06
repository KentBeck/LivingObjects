#pragma once

#include "ast.h"

#include <memory>
#include <string>

namespace smalltalk
{

    /**
     * Simple recursive descent parser for basic expressions and unary messages
     * Supports: integers, booleans, nil, +, -, *, /, <, >, =, ~=, <=, >=, (, ), unary messages
     *
     * Grammar:
     *   expression := comparison
     *   comparison := arithmetic (('<' | '>' | '=' | '~=' | '<=' | '>=') arithmetic)*
     *   arithmetic := term (('+' | '-') term)*
     *   term       := factor (('*' | '/') factor)*
     *   factor     := unary
     *   unary      := primary (identifier)*
     *   primary    := integer | literal | string | identifier | '(' expression ')'
     *   integer    := [0-9]+
     *   literal    := 'true' | 'false' | 'nil'
     *   string     := "'" [^']* "'"
     *   identifier := [a-zA-Z][a-zA-Z0-9]*
     */
    class SimpleParser
    {
    public:
        SimpleParser(std::string input);

        // Parse the input and return a method AST
        std::unique_ptr<MethodNode> parseMethod();

    private:
        // Parsing methods
        ASTNodePtr parseExpression();
        ASTNodePtr parseStatements();
        ASTNodePtr parseStatement();
        ASTNodePtr parseAssignmentExpression();
        ASTNodePtr parseReturn();
        ASTNodePtr parseBinaryMessage();
        ASTNodePtr parseComparison();
        ASTNodePtr parseFactor();
        ASTNodePtr parseUnary();
        ASTNodePtr parsePrimary();
        ASTNodePtr parseInteger();
        ASTNodePtr parseIdentifierOrLiteral();
        ASTNodePtr parseString();
        ASTNodePtr parseVariable();
        ASTNodePtr parseBlock();
        std::string parseIdentifier();

        // Temporary variable parsing
        std::vector<std::string> parseTemporaryVariables();
        bool isTemporaryVariableDeclaration();

        // Tokenization
        void skipWhitespace();
        char peek() const;
        char consume();
        bool isAtEnd() const;
        bool isDigit(char c) const;
        bool isAlpha(char c) const;

        // Binary selector helpers
        bool isBinarySelector();
        std::string parseBinarySelector();

        // Error handling
        void error(const std::string &message);

        std::string input_;
        size_t pos_ = 0;
    };

} // namespace smalltalk
