#pragma once

#include "tagged_value.h"
#include "smalltalk_string.h"

#include <memory>
#include <string>
#include <vector>

namespace smalltalk
{

    // Forward declarations
    class ASTNode;
    using ASTNodePtr = std::unique_ptr<ASTNode>;

    /**
     * Base class for AST nodes
     */
    class ASTNode
    {
    public:
        virtual ~ASTNode() = default;
        virtual std::string toString() const = 0;
    };

    /**
     * Literal value node (numbers, nil, true, false)
     */
    class LiteralNode : public ASTNode
    {
    public:
        explicit LiteralNode(TaggedValue value) : value_(value) {}

        TaggedValue getValue() const { return value_; }

        std::string toString() const override
        {
            if (value_.isInteger())
            {
                return std::to_string(value_.asInteger());
            }
            else if (value_.isFloat())
            {
                return std::to_string(value_.asFloat());
            }
            else if (value_.isNil())
            {
                return "nil";
            }
            else if (value_.isTrue())
            {
                return "true";
            }
            else if (value_.isFalse())
            {
                return "false";
            }
            else if (value_.isPointer())
            {
                // Check if it's a string
                try
                {
                    Object *obj = value_.asObject();
                    if (obj && obj->header.getType() == ObjectType::OBJECT)
                    {
                        // Try to cast to String
                        String *str = static_cast<String *>(obj);
                        return "'" + str->getContent() + "'";
                    }
                }
                catch (...)
                {
                    // Fall through to unknown
                }
            }
            return "unknown";
        }

    private:
        TaggedValue value_;
    };

    /**
     * Binary operation node (+, -, *, /, etc.)
     */
    class BinaryOpNode : public ASTNode
    {
    public:
        enum class Operator
        {
            Add,
            Subtract,
            Multiply,
            Divide,
            Concatenate,
            LessThan,
            GreaterThan,
            Equal,
            NotEqual,
            LessThanOrEqual,
            GreaterThanOrEqual
        };

        BinaryOpNode(ASTNodePtr left, Operator op, ASTNodePtr right)
            : left_(std::move(left)), op_(op), right_(std::move(right)) {}

        const ASTNode *getLeft() const { return left_.get(); }
        const ASTNode *getRight() const { return right_.get(); }
        Operator getOperator() const { return op_; }

        std::string toString() const override
        {
            std::string opStr;
            switch (op_)
            {
            case Operator::Add:
                opStr = "+";
                break;
            case Operator::Subtract:
                opStr = "-";
                break;
            case Operator::Multiply:
                opStr = "*";
                break;
            case Operator::Divide:
                opStr = "/";
                break;
            case Operator::Concatenate:
                opStr = ",";
                break;
            case Operator::LessThan:
                opStr = "<";
                break;
            case Operator::GreaterThan:
                opStr = ">";
                break;
            case Operator::Equal:
                opStr = "=";
                break;
            case Operator::NotEqual:
                opStr = "~=";
                break;
            case Operator::LessThanOrEqual:
                opStr = "<=";
                break;
            case Operator::GreaterThanOrEqual:
                opStr = ">=";
                break;
            default:
                opStr = "?";
                break;
            }
            return "(" + left_->toString() + " " + opStr + " " + right_->toString() + ")";
        }

    private:
        ASTNodePtr left_;
        Operator op_;
        ASTNodePtr right_;
    };

    /**
     * Method node - represents a complete method with an expression body
     */
    class MethodNode : public ASTNode
    {
    public:
        explicit MethodNode(ASTNodePtr body) : body_(std::move(body)) {}

        const ASTNode *getBody() const { return body_.get(); }

        std::string toString() const override
        {
            return "method { " + body_->toString() + " }";
        }

    private:
        ASTNodePtr body_;
    };

    /**
     * Block node for expressions like [3 + 4]
     */
    class BlockNode : public ASTNode
    {
    public:
        explicit BlockNode(ASTNodePtr body) : body_(std::move(body)) {}

        const ASTNode *getBody() const { return body_.get(); }

        std::string toString() const override
        {
            return "[" + body_->toString() + "]";
        }

    private:
        ASTNodePtr body_;
    };

    /**
     * Sequence node for multiple statements separated by periods
     */
    class SequenceNode : public ASTNode
    {
    public:
        SequenceNode() = default;

        void addStatement(ASTNodePtr statement)
        {
            statements_.push_back(std::move(statement));
        }

        const std::vector<ASTNodePtr> &getStatements() const { return statements_; }

        std::string toString() const override
        {
            if (statements_.empty())
            {
                return "";
            }

            std::string result = statements_[0]->toString();
            for (size_t i = 1; i < statements_.size(); i++)
            {
                result += ". " + statements_[i]->toString();
            }
            return result;
        }

    private:
        std::vector<ASTNodePtr> statements_;
    };

    /**
     * Message send node for expressions like "Object new" or "array at: 1"
     */
    class MessageSendNode : public ASTNode
    {
    public:
        MessageSendNode(ASTNodePtr receiver, std::string selector, std::vector<ASTNodePtr> arguments)
            : receiver_(std::move(receiver)), selector_(std::move(selector)), arguments_(std::move(arguments)) {}

        const ASTNode *getReceiver() const { return receiver_.get(); }
        const std::string &getSelector() const { return selector_; }
        const std::vector<ASTNodePtr> &getArguments() const { return arguments_; }

        std::string toString() const override
        {
            std::string result = receiver_->toString() + " " + selector_;
            for (const auto &arg : arguments_)
            {
                result += " " + arg->toString();
            }
            return result;
        }

    private:
        ASTNodePtr receiver_;
        std::string selector_;
        std::vector<ASTNodePtr> arguments_;
    };

} // namespace smalltalk
