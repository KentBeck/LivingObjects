#include "simple_parser.h"
#include "smalltalk_string.h"
#include "smalltalk_class.h"

#include <cctype>
#include <iostream>
#include <stdexcept>

namespace smalltalk
{

    SimpleParser::SimpleParser(std::string input)
        : input_(std::move(input))
    {
    }

    std::unique_ptr<MethodNode> SimpleParser::parseMethod()
    {
        skipWhitespace();

        // Check for temporary variable declarations
        std::vector<std::string> tempVars;
        if (isTemporaryVariableDeclaration())
        {
            tempVars = parseTemporaryVariables();
            skipWhitespace();
        }

        auto body = parseStatements();
        skipWhitespace();

        if (!isAtEnd())
        {
            error("Unexpected characters at end of input");
        }

        // Always return a MethodNode, with empty tempVars if none were declared
        return std::make_unique<MethodNode>(std::move(tempVars), std::move(body));
    }

    ASTNodePtr SimpleParser::parseExpression()
    {
        return parseAssignmentExpression();
    }

    ASTNodePtr SimpleParser::parseStatements()
    {
        auto statements = std::make_unique<SequenceNode>();

        // Parse first statement
        auto stmt = parseStatement();
        statements->addStatement(std::move(stmt));

        // Parse additional statements separated by periods
        while (!isAtEnd())
        {
            skipWhitespace();
            if (peek() == '.')
            {
                consume(); // consume '.'
                skipWhitespace();

                if (isAtEnd())
                {
                    break; // Allow trailing period
                }

                auto nextStmt = parseStatement();
                statements->addStatement(std::move(nextStmt));
            }
            else
            {
                break;
            }
        }

        // If only one statement, return it directly
        if (statements->getStatements().size() == 1)
        {
            return std::move(const_cast<std::vector<ASTNodePtr> &>(statements->getStatements())[0]);
        }

        return std::move(statements);
    }

    ASTNodePtr SimpleParser::parseStatement()
    {
        skipWhitespace();
        if (peek() == '^') {
            return parseReturn();
        }

        // For statements, just parse as expression (assignment is now part of expression parsing)
        return parseExpression();
    }

    ASTNodePtr SimpleParser::parseReturn() {
        consume(); // consume '^'
        skipWhitespace();
        auto value = parseExpression();
        return std::make_unique<ReturnNode>(std::move(value));
    }

    ASTNodePtr SimpleParser::parseAssignmentExpression()
    {
        // Check for assignment (lowest precedence)
        size_t savedPos = pos_;
        skipWhitespace();

        // Look ahead for identifier followed by :=
        if (isAlpha(peek()))
        {
            std::string identifier;
            while (!isAtEnd() && (isAlpha(peek()) || isDigit(peek())))
            {
                identifier += consume();
            }

            skipWhitespace();
            if (!isAtEnd() && peek() == ':' && pos_ + 1 < input_.size() && input_[pos_ + 1] == '=')
            {
                // This is an assignment
                consume(); // consume ':'
                consume(); // consume '='
                skipWhitespace();

                auto value = parseAssignmentExpression(); // Right-associative
                return std::make_unique<AssignmentNode>(std::move(identifier), std::move(value));
            }
        }

        // Not an assignment, restore position and parse binary messages
        pos_ = savedPos;
        return parseBinaryMessage();
    }

    // Helper to check if current position starts a binary selector
    bool SimpleParser::isBinarySelector()
    {
        if (isAtEnd())
            return false;

        char c = peek();
        // Check for single-character binary selectors
        if (c == '+' || c == '-' || c == '*' || c == '/' || c == ',' ||
            c == '<' || c == '>' || c == '=' || c == '~')
        {
            return true;
        }
        return false;
    }

    // Parse a binary selector and return it as a string
    std::string SimpleParser::parseBinarySelector()
    {
        std::string selector;

        char c = peek();
        if (c == '<')
        {
            consume();
            selector = "<";
            if (!isAtEnd() && peek() == '=')
            {
                consume();
                selector = "<=";
            }
        }
        else if (c == '>')
        {
            consume();
            selector = ">";
            if (!isAtEnd() && peek() == '=')
            {
                consume();
                selector = ">=";
            }
        }
        else if (c == '~')
        {
            consume();
            if (!isAtEnd() && peek() == '=')
            {
                consume();
                selector = "~=";
            }
            else
            {
                error("Expected '=' after '~'");
            }
        }
        else if (c == '+' || c == '-' || c == '*' || c == '/' || c == ',' || c == '=')
        {
            consume();
            selector = c;
        }
        else
        {
            error("Invalid binary selector");
        }

        return selector;
    }

    ASTNodePtr SimpleParser::parseBinaryMessage()
    {
        auto left = parseUnary();

        while (!isAtEnd())
        {
            skipWhitespace();

            if (isBinarySelector())
            {
                std::string selector = parseBinarySelector();
                skipWhitespace();
                auto right = parseUnary();

                std::vector<ASTNodePtr> args;
                args.push_back(std::move(right));
                left = std::make_unique<MessageSendNode>(std::move(left), selector, std::move(args));
            }
            else
            {
                break;
            }
        }

        return left;
    }

    ASTNodePtr SimpleParser::parseComparison()
    {
        // This method is now unused but kept for compatibility
        return parseBinaryMessage();
    }

    ASTNodePtr SimpleParser::parseFactor()
    {
        return parseUnary();
    }

    ASTNodePtr SimpleParser::parseUnary()
    {
        auto receiver = parsePrimary();

        // Parse chain of unary messages
        while (!isAtEnd())
        {
            skipWhitespace();

            // Check if next token is an identifier (unary message)
            if (isAlpha(peek()))
            {
                // Save current position to check if this is actually a message
                size_t savedPos = pos_;

                // Parse the identifier
                std::string selector;
                while (!isAtEnd() && (isAlpha(peek()) || isDigit(peek())))
                {
                    selector += consume();
                }

                // Check if this is followed by something that would make it NOT a unary message
                skipWhitespace();
                char next = isAtEnd() ? '\0' : peek();

                // If followed by operators or end of input, it's a unary message
                if (next == '\0' || next == '+' || next == '-' || next == '*' || next == '/' ||
                    next == '<' || next == '>' || next == '=' || next == '~' ||
                    next == ')' || next == ']' || next == '.')
                {

                    // It's a unary message
                    std::vector<ASTNodePtr> args; // No arguments for unary messages
                    receiver = std::make_unique<MessageSendNode>(std::move(receiver), selector, std::move(args));
                }
                else
                {
                    // Not a unary message, restore position
                    pos_ = savedPos;
                    break;
                }
            }
            else
            {
                break;
            }
        }

        return receiver;
    }

    ASTNodePtr SimpleParser::parsePrimary()
    {
        skipWhitespace();

        if (peek() == '(')
        {
            consume(); // consume '('
            auto expr = parseExpression(); // Allow assignments in parentheses
            skipWhitespace();

            if (peek() != ')')
            {
                error("Expected ')' after expression");
            }
            consume(); // consume ')'

            return expr;
        }
        else if (peek() == '[')
        {
            return parseBlock();
        }
        else if (isDigit(peek()))
        {
            return parseInteger();
        }
        else if (isAlpha(peek()))
        {
            return parseIdentifierOrLiteral();
        }
        else if (peek() == '\'')
        {
            return parseString();
        }
        else
        {
            error(std::string("Unexpected character: ") + peek());
            return nullptr; // Never reached
        }
    }

    ASTNodePtr SimpleParser::parseInteger()
    {
        std::string numStr;

        while (!isAtEnd() && isDigit(peek()))
        {
            numStr += consume();
        }

        if (numStr.empty())
        {
            error("Expected integer");
        }

        try
        {
            int value = std::stoi(numStr);
            return std::make_unique<LiteralNode>(TaggedValue(static_cast<int32_t>(value)));
        }
        catch (const std::exception &e)
        {
            error("Invalid integer: " + numStr);
            return nullptr; // Never reached
        }
    }

    ASTNodePtr SimpleParser::parseIdentifierOrLiteral()
    {
        std::string identifier;

        while (!isAtEnd() && (isAlpha(peek()) || isDigit(peek())))
        {
            identifier += consume();
        }

        if (identifier.empty())
        {
            error("Expected identifier");
        }

        // Check for boolean and nil literals
        if (identifier == "true")
        {
            return std::make_unique<LiteralNode>(TaggedValue::trueValue());
        }
        else if (identifier == "false")
        {
            return std::make_unique<LiteralNode>(TaggedValue::falseValue());
        }
        else if (identifier == "nil")
        {
            return std::make_unique<LiteralNode>(TaggedValue::nil());
        }

        // Check for global class names
        ClassRegistry &registry = ClassRegistry::getInstance();
        if (registry.hasClass(identifier))
        {
            Class *clazz = registry.getClass(identifier);
            return std::make_unique<LiteralNode>(TaggedValue::fromObject(clazz));
        }

        // Treat unrecognized identifiers as variable references
        return std::make_unique<VariableNode>(identifier);
    }

    ASTNodePtr SimpleParser::parseString()
    {
        if (peek() != '\'')
        {
            error("Expected string to start with '");
        }
        consume(); // consume opening '

        std::string content;

        while (!isAtEnd() && peek() != '\'')
        {
            if (peek() == '\\')
            {
                // Handle escape sequences
                consume(); // consume backslash
                if (isAtEnd())
                {
                    error("Unexpected end of input in string literal");
                }

                char escaped = consume();
                switch (escaped)
                {
                case 'n':
                    content += '\n';
                    break;
                case 't':
                    content += '\t';
                    break;
                case 'r':
                    content += '\r';
                    break;
                case '\\':
                    content += '\\';
                    break;
                case '\'':
                    content += '\'';
                    break;
                default:
                    content += escaped; // Keep unknown escapes as-is
                    break;
                }
            }
            else
            {
                content += consume();
            }
        }

        if (isAtEnd() || peek() != '\'')
        {
            error("Unterminated string literal");
        }
        consume(); // consume closing '

        // Create a String object and wrap it in a TaggedValue
        TaggedValue stringValue = StringUtils::createTaggedString(content);
        return std::make_unique<LiteralNode>(stringValue);
    }

    void SimpleParser::skipWhitespace()
    {
        while (!isAtEnd() && (std::isspace(peek()) != 0))
        {
            consume();
        }
    }

    char SimpleParser::peek() const
    {
        if (isAtEnd())
        {
            return '\0';
        }
        return input_[pos_];
    }

    char SimpleParser::consume()
    {
        if (isAtEnd())
        {
            error("Unexpected end of input");
        }
        return input_[pos_++];
    }

    bool SimpleParser::isAtEnd() const
    {
        return pos_ >= input_.size();
    }

    bool SimpleParser::isDigit(char c) const
    {
        return c >= '0' && c <= '9';
    }

    bool SimpleParser::isAlpha(char c) const
    {
        return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z');
    }

    ASTNodePtr SimpleParser::parseBlock()
    {
        consume(); // consume '['
        skipWhitespace();

        std::vector<std::string> parameters;

        // Check for block parameters starting with ':'
        if (peek() == ':')
        {
            // Parse block parameters
            while (peek() == ':')
            {
                consume(); // consume ':'
                skipWhitespace();

                if (!isAlpha(peek()))
                {
                    error("Expected identifier after ':' in block parameter");
                }

                std::string paramName;
                while (!isAtEnd() && (isAlpha(peek()) || isDigit(peek())))
                {
                    paramName += consume();
                }
                parameters.push_back(paramName);
                skipWhitespace();
            }

            // Expect '|' after parameters
            if (peek() != '|')
            {
                error("Expected '|' after block parameters");
            }
            consume(); // consume '|'
            skipWhitespace();
        }

        // Parse block body
        auto body = parseStatements();
        skipWhitespace();

        if (peek() != ']')
        {
            error("Expected ']' after block expression");
        }
        consume(); // consume ']'

        if (parameters.empty())
        {
            return std::make_unique<BlockNode>(std::move(body));
        }
        else
        {
            return std::make_unique<BlockNode>(std::move(parameters), std::move(body));
        }
    }

    void SimpleParser::error(const std::string &message)
    {
        throw std::runtime_error("Parse error at position " + std::to_string(pos_) + ": " + message);
    }

    std::vector<std::string> SimpleParser::parseTemporaryVariables()
    {
        std::vector<std::string> tempVars;

        if (peek() != '|')
        {
            error("Expected '|' to start temporary variable declaration");
        }
        consume(); // consume '|'

        skipWhitespace();

        // Parse variable names
        while (!isAtEnd() && peek() != '|')
        {
            if (!isAlpha(peek()))
            {
                error("Expected variable name in temporary variable declaration");
            }

            std::string varName;
            while (!isAtEnd() && (isAlpha(peek()) || isDigit(peek())))
            {
                varName += consume();
            }

            tempVars.push_back(varName);
            skipWhitespace();
        }

        if (peek() != '|')
        {
            error("Expected '|' to end temporary variable declaration");
        }
        consume(); // consume '|'

        return tempVars;
    }

    bool SimpleParser::isTemporaryVariableDeclaration()
    {
        return !isAtEnd() && peek() == '|';
    }

    ASTNodePtr SimpleParser::parseVariable()
    {
        if (!isAlpha(peek()))
        {
            error("Expected variable name");
        }

        std::string varName;
        while (!isAtEnd() && (isAlpha(peek()) || isDigit(peek())))
        {
            varName += consume();
        }

        return std::make_unique<VariableNode>(std::move(varName));
    }

    std::string SimpleParser::parseIdentifier()
    {
        if (!isAlpha(peek()))
        {
            error("Expected identifier");
        }

        std::string identifier;
        while (!isAtEnd() && (isAlpha(peek()) || isDigit(peek())))
        {
            identifier += consume();
        }

        return identifier;
    }

} // namespace smalltalk
