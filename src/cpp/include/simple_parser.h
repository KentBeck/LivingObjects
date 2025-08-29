#pragma once

#include "ast.h"

#include <memory>
#include <string>

namespace smalltalk {

/**
 * Recursive descent parser for Smalltalk expressions following proper
 * precedence
 *
 * Grammar (following Smalltalk-80 specification):
 *   method := temporaries? statements
 *   temporaries := '|' identifier* '|'
 *   statements := statement ('.' statement)* '.'?
 *   statement := expression | '^' expression
 *   expression := assignment | keywordExpression
 *   assignment := identifier ':=' expression
 *   keywordExpression := binaryExpression (keyword binaryExpression)*
 *   binaryExpression := unaryExpression (binarySelector unaryExpression)*
 *   unaryExpression := primary unarySelector*
 *   primary := identifier | literal | block | '(' expression ')' | arrayLiteral
 *   keyword := identifier ':'
 *   binarySelector := '+' | '-' | '*' | '/' | '=' | '~=' | '<' | '>' | '<=' |
 * '>=' | ',' unarySelector := identifier literal := integer | string | symbol |
 * 'true' | 'false' | 'nil' block := '[' blockBody ']' arrayLiteral := '#('
 * literalElement* ')'
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
  ASTNodePtr parseSymbol();
  ASTNodePtr parseArrayLiteral();
  ASTNodePtr parseVariable();
  ASTNodePtr parseBlock();
  std::string parseIdentifier();
  ASTNodePtr parseKeywordMessage();

  // Temporary variable parsing
  std::vector<std::string> parseTemporaryVariables();
  bool isTemporaryVariableDeclaration();

  // Primitive parsing
  int parsePrimitive();

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
