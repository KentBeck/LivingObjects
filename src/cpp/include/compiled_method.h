#pragma once

#include "tagged_value.h"
#include "symbol.h"
#include "object.h"

#include <cstdint>
#include <string>
#include <vector>

namespace smalltalk
{

    /**
     * A compiled method containing bytecode and literal values
     * Now properly inherits from Object to be a real Smalltalk object
     */
    class CompiledMethod : public Object
    {
    public:
        // Constructor that properly initializes Object base class
        CompiledMethod() : Object(ObjectType::METHOD, sizeof(CompiledMethod)) {}
        virtual ~CompiledMethod() = default;

        // Primitive number (0 if no primitive)
        int primitiveNumber = 0;

        // Bytecodes, literals, and temporary variables (public for easy access)
        std::vector<uint8_t> bytecodes;
        std::vector<TaggedValue> literals;
        std::vector<std::string> tempVars;

        // Add bytecode
        void addBytecode(uint8_t bytecode)
        {
            bytecodes.push_back(bytecode);
        }

        // Add a 4-byte operand (little-endian)
        void addOperand(uint32_t operand)
        {
            bytecodes.push_back(static_cast<uint8_t>(operand & 0xFF));
            bytecodes.push_back(static_cast<uint8_t>((operand >> 8) & 0xFF));
            bytecodes.push_back(static_cast<uint8_t>((operand >> 16) & 0xFF));
            bytecodes.push_back(static_cast<uint8_t>((operand >> 24) & 0xFF));
        }

        // Add a literal value and return its index
        uint32_t addLiteral(TaggedValue value)
        {
            uint32_t index = static_cast<uint32_t>(literals.size());
            literals.push_back(value);
            return index;
        }

        // Add a temporary variable and return its index
        uint32_t addTempVar(const std::string &name)
        {
            uint32_t index = static_cast<uint32_t>(tempVars.size());
            tempVars.push_back(name);
            return index;
        }

        // Getters
        const std::vector<uint8_t> &getBytecodes() const { return bytecodes; }
        const std::vector<TaggedValue> &getLiterals() const { return literals; }
        const std::vector<std::string> &getTempVars() const { return tempVars; }

        // Get hash of the method
        uint32_t getHash() const {
            // Simple hash for now, can be improved later
            uint32_t hash = 0;
            for (uint8_t b : bytecodes) {
                hash = (hash << 5) + hash + b;
            }
            for (const auto& lit : literals) {
                hash = (hash << 5) + hash + lit.rawValue();
            }
            return hash;
        }

        // Get literal by index
        TaggedValue getLiteral(uint32_t index) const
        {
            if (index >= literals.size())
            {
                throw std::runtime_error("Literal index out of bounds");
            }
            return literals[index];
        }

        // Debug output
        virtual std::string toString() const
        {
            std::string result = "CompiledMethod {\n";
            result += "  Bytecodes: [";
            for (size_t i = 0; i < bytecodes.size(); ++i)
            {
                if (i > 0)
                    result += ", ";
                result += std::to_string(static_cast<int>(bytecodes[i]));
            }
            result += "]\n";

            result += "  Literals: [";
            for (size_t i = 0; i < literals.size(); ++i)
            {
                if (i > 0)
                    result += ", ";
                if (literals[i].isNil())
                {
                    result += "nil";
                }
                else if (literals[i].isBoolean())
                {
                    result += literals[i].asBoolean() ? "true" : "false";
                }
                else if (literals[i].isInteger())
                {
                    result += std::to_string(literals[i].asInteger());
                }
                else if (literals[i].isPointer())
                {
                    try
                    {
                        Symbol *symbol = literals[i].asSymbol();
                        result += symbol->toString();
                    }
                    catch (...)
                    {
                        result += "Object@" + std::to_string(reinterpret_cast<uintptr_t>(literals[i].asPointer()));
                    }
                }
                else
                {
                    result += "?";
                }
            }
            result += "]\n";
            result += "}";

            return result;
        }
    };

} // namespace smalltalk
