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
        auto body = parseExpression();
        skipWhitespace();

        if (!isAtEnd())
        {
            error("Unexpected characters at end of input");
        }

        return std::make_unique<MethodNode>(std::move(body));
    }

    ASTNodePtr SimpleParser::parseExpression()
    {
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
            auto expr = parseExpression();
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
            consume(); // consume '['
            auto expr = parseStatements();
            skipWhitespace();

            if (peek() != ']')
            {
                error("Expected ']' after block expression");
            }
            consume(); // consume ']'

            return std::make_unique<BlockNode>(std::move(expr));
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

        // For now, treat unrecognized identifiers as variables (temporary)
        // In a full implementation, this would handle local variables, instance variables, etc.
        error("Unknown identifier: " + identifier);
        return nullptr; // Never reached
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

    ASTNodePtr SimpleParser::parseStatements()
    {
        skipWhitespace();

        // Handle empty blocks
        if (peek() == ']' || isAtEnd())
        {
            // Return a nil literal for empty blocks
            return std::make_unique<LiteralNode>(TaggedValue());
        }

        auto sequence = std::make_unique<SequenceNode>();

        // Parse first statement
        auto statement = parseExpression();
        sequence->addStatement(std::move(statement));

        // Parse additional statements separated by periods
        while (!isAtEnd())
        {
            skipWhitespace();

            if (peek() == '.')
            {
                consume(); // consume '.'
                skipWhitespace();

                // Check if this is the end (trailing period)
                if (peek() == ']' || isAtEnd())
                {
                    break;
                }

                // Parse next statement
                auto nextStatement = parseExpression();
                sequence->addStatement(std::move(nextStatement));
            }
            else
            {
                // No more periods, we're done
                break;
            }
        }

        // If only one statement, return it directly instead of a sequence
        const auto &statements = sequence->getStatements();
        if (statements.size() == 1)
        {
            // Create a copy of the single statement
            const auto *statement = statements[0].get();
            if (const auto *literal = dynamic_cast<const LiteralNode *>(statement))
            {
                return std::make_unique<LiteralNode>(literal->getValue());
            }
            // For other types, return the sequence
            return std::move(sequence);
        }

        return std::move(sequence);
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

    void SimpleParser::error(const std::string &message)
    {
        throw std::runtime_error("Parse error at position " + std::to_string(pos_) + ": " + message);
    }

} // namespace smalltalk
