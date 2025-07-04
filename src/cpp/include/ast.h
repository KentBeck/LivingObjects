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
     * Method node - represents a complete method with optional temporary variables and an expression body
     */
    class MethodNode : public ASTNode
    {
    public:
        explicit MethodNode(ASTNodePtr body) : tempVars_(), body_(std::move(body)) {}
        MethodNode(std::vector<std::string> tempVars, ASTNodePtr body)
            : tempVars_(std::move(tempVars)), body_(std::move(body)) {}

        const std::vector<std::string> &getTempVars() const { return tempVars_; }
        const ASTNode *getBody() const { return body_.get(); }

        std::string toString() const override
        {
            std::string result = "method ";
            if (!tempVars_.empty())
            {
                result += "| ";
                for (size_t i = 0; i < tempVars_.size(); i++)
                {
                    if (i > 0)
                        result += " ";
                    result += tempVars_[i];
                }
                result += " | ";
            }
            result += "{ " + body_->toString() + " }";
            return result;
        }

    private:
        std::vector<std::string> tempVars_;
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

    /**
     * Variable reference node for accessing variables
     */
    class VariableNode : public ASTNode
    {
    public:
        explicit VariableNode(std::string name) : name_(std::move(name)) {}

        const std::string &getName() const { return name_; }

        std::string toString() const override
        {
            return name_;
        }

    private:
        std::string name_;
    };

    /**
     * Assignment node for variable assignments like "x := 42"
     */
    class AssignmentNode : public ASTNode
    {
    public:
        AssignmentNode(std::string variable, ASTNodePtr value)
            : variable_(std::move(variable)), value_(std::move(value)) {}

        const std::string &getVariable() const { return variable_; }
        const ASTNode *getValue() const { return value_.get(); }

        std::string toString() const override
        {
            return variable_ + " := " + value_->toString();
        }

    private:
        std::string variable_;
        ASTNodePtr value_;
    };

    /**
     * Method node with temporary variables - represents a complete method
     */
    class MethodWithTempsNode : public ASTNode
    {
    public:
        MethodWithTempsNode(std::vector<std::string> tempVars, ASTNodePtr body)
            : tempVars_(std::move(tempVars)), body_(std::move(body)) {}

        const std::vector<std::string> &getTempVars() const { return tempVars_; }
        const ASTNode *getBody() const { return body_.get(); }

        std::string toString() const override
        {
            std::string result = "method ";
            if (!tempVars_.empty())
            {
                result += "| ";
                for (size_t i = 0; i < tempVars_.size(); i++)
                {
                    if (i > 0)
                        result += " ";
                    result += tempVars_[i];
                }
                result += " | ";
            }
            result += "{ " + body_->toString() + " }";
            return result;
        }

    private:
        std::vector<std::string> tempVars_;
        ASTNodePtr body_;
    };

} // namespace smalltalk
