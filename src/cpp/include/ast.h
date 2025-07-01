#pragma once

#include "tagged_value.h"

#include <memory>
#include <string>
#include <vector>

namespace smalltalk {

// Forward declarations
class ASTNode;
using ASTNodePtr = std::unique_ptr<ASTNode>;

/**
 * Base class for AST nodes
 */
class ASTNode {
public:
    virtual ~ASTNode() = default;
    virtual std::string toString() const = 0;
};

/**
 * Literal value node (numbers, nil, true, false)
 */
class LiteralNode : public ASTNode {
public:
    explicit LiteralNode(TaggedValue value) : value_(value) {}
    
    TaggedValue getValue() const { return value_; }
    
    std::string toString() const override {
        if (value_.isInteger()) {
            return std::to_string(value_.asInteger());
        } else if (value_.isFloat()) {
            return std::to_string(value_.asFloat());
        } else if (value_.isNil()) {
            return "nil";
        } else if (value_.isTrue()) {
            return "true";
        } else if (value_.isFalse()) {
            return "false";
        }
        return "unknown";
    }
    
private:
    TaggedValue value_;
};

/**
 * Binary operation node (+, -, *, /, etc.)
 */
class BinaryOpNode : public ASTNode {
public:
    enum class Operator {
        Add, Subtract, Multiply, Divide,
        LessThan, GreaterThan, Equal, NotEqual,
        LessThanOrEqual, GreaterThanOrEqual
    };
    
    BinaryOpNode(ASTNodePtr left, Operator op, ASTNodePtr right)
        : left_(std::move(left)), op_(op), right_(std::move(right)) {}
    
    const ASTNode* getLeft() const { return left_.get(); }
    const ASTNode* getRight() const { return right_.get(); }
    Operator getOperator() const { return op_; }
    
    std::string toString() const override {
        std::string opStr;
        switch (op_) {
            case Operator::Add: opStr = "+"; break;
            case Operator::Subtract: opStr = "-"; break;
            case Operator::Multiply: opStr = "*"; break;
            case Operator::Divide: opStr = "/"; break;
            case Operator::LessThan: opStr = "<"; break;
            case Operator::GreaterThan: opStr = ">"; break;
            case Operator::Equal: opStr = "="; break;
            case Operator::NotEqual: opStr = "~="; break;
            case Operator::LessThanOrEqual: opStr = "<="; break;
            case Operator::GreaterThanOrEqual: opStr = ">="; break;
            default: opStr = "?"; break;
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
class MethodNode : public ASTNode {
public:
    explicit MethodNode(ASTNodePtr body) : body_(std::move(body)) {}
    
    const ASTNode* getBody() const { return body_.get(); }
    
    std::string toString() const override {
        return "method { " + body_->toString() + " }";
    }
    
private:
    ASTNodePtr body_;
};

} // namespace smalltalk
