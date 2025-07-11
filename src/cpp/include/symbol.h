#pragma once

#include "object.h"
#include <string>
#include <unordered_map>
#include <memory>

namespace smalltalk
{

    /**
     * Symbol represents an interned string used for selectors and identifiers.
     * Symbols with the same string content are guaranteed to be identical objects.
     */
    class Symbol : public Object
    {
    public:
        // Get or create a symbol with the given name
        static Symbol *intern(const std::string &name);

        // Get the symbol's name
        const std::string &getName() const { return name_; }

        // Symbols can be compared by pointer equality since they're interned
        bool operator==(const Symbol &other) const
        {
            return this == &other;
        }

        bool operator!=(const Symbol &other) const
        {
            return this != &other;
        }

        // Hash function for use in maps
        size_t hash() const
        {
            return std::hash<std::string>{}(name_);
        }

        // String representation
        virtual std::string toString() const override
        {
            return "Symbol(" + name_ + ")";
        }

    private:
        explicit Symbol(std::string name) 
            : Object(ObjectType::SYMBOL, sizeof(Symbol)), 
              name_(std::move(name)) {}

        std::string name_;

        // Global symbol table for interning
        static std::unordered_map<std::string, std::unique_ptr<Symbol>> symbolTable_;
    };

} // namespace smalltalk