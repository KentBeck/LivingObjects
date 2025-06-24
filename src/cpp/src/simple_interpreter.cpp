#include "../include/simple_interpreter.h"
#include <algorithm>
#include <cctype>
#include <sstream>
#include <limits>

namespace smalltalk
{

    TaggedValue SimpleInterpreter::evaluate(const std::string &expression)
    {
        std::string trimmed = trim(expression);

        // Handle special values
        if (trimmed == "nil")
        {
            return TaggedValue::nil();
        }
        if (trimmed == "true")
        {
            return TaggedValue::trueValue();
        }
        if (trimmed == "false")
        {
            return TaggedValue::falseValue();
        }

        // Try to parse as integer
        int32_t intValue;
        if (tryParseInteger(trimmed, intValue))
        {
            return TaggedValue(intValue);
        }

        // If we get here, the expression is not supported
        throw std::runtime_error("Unsupported expression: " + expression);
    }

    bool SimpleInterpreter::tryParseInteger(const std::string &str, int32_t &result)
    {
        if (str.empty())
        {
            return false;
        }

        try
        {
            // Use stringstream for parsing
            std::stringstream ss(str);
            long long temp;
            ss >> temp;

            // Check if the entire string was consumed
            if (!ss.eof())
            {
                return false;
            }

            // Check if the value fits in int32_t
            if (temp < std::numeric_limits<int32_t>::min() ||
                temp > std::numeric_limits<int32_t>::max())
            {
                return false;
            }

            result = static_cast<int32_t>(temp);
            return true;
        }
        catch (...)
        {
            return false;
        }
    }

    std::string SimpleInterpreter::trim(const std::string &str)
    {
        // Find first non-whitespace character
        auto start = std::find_if(str.begin(), str.end(), [](unsigned char ch)
                                  { return !std::isspace(ch); });

        // Find last non-whitespace character
        auto end = std::find_if(str.rbegin(), str.rend(), [](unsigned char ch)
                                { return !std::isspace(ch); })
                       .base();

        // Return trimmed string
        return (start < end) ? std::string(start, end) : std::string();
    }

} // namespace smalltalk
