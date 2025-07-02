#pragma once

#include "object.h"
#include "tagged_value.h"
#include <string>
#include <memory>

namespace smalltalk {

/**
 * String represents a heap-allocated string object in Smalltalk.
 * Unlike immediate values, strings are stored on the heap and
 * referenced via tagged pointers.
 */
class String : public Object {
public:
    String(const std::string& content, Class* stringClass);
    String(const char* content, Class* stringClass);
    
    // Get the string content
    const std::string& getContent() const { return content_; }
    
    // Get the size of the string
    size_t size() const { return content_.size(); }
    size_t length() const { return content_.size(); }
    
    // Check if string is empty
    bool isEmpty() const { return content_.empty(); }
    
    // String concatenation
    String* concatenate(const String* other) const;
    String* concatenate(const std::string& other) const;
    
    // String comparison
    bool equals(const String* other) const;
    bool equals(const std::string& other) const;
    
    // Character access
    char at(size_t index) const;
    
    // String operations
    String* substring(size_t start, size_t length) const;
    String* toLowerCase() const;
    String* toUpperCase() const;
    
    // Find operations
    int indexOf(char ch) const;
    int indexOf(const std::string& substr) const;
    bool contains(const std::string& substr) const;
    
    // String representation
    std::string toString() const override;
    
    // Conversion to C++ string
    const std::string& toCppString() const { return content_; }
    
    // Hash for use in collections
    size_t hash() const;
    
    // Comparison operators
    bool operator==(const String& other) const;
    bool operator!=(const String& other) const;
    bool operator<(const String& other) const;
    bool operator>(const String& other) const;
    bool operator<=(const String& other) const;
    bool operator>=(const String& other) const;
    
private:
    std::string content_;
};

/**
 * StringUtils provides utility functions for string operations
 */
namespace StringUtils {
    // Create a new String object
    String* createString(const std::string& content);
    String* createString(const char* content);
    
    // Convert TaggedValue to String (if it contains a string)
    String* asString(const TaggedValue& value);
    
    // Check if TaggedValue contains a string
    bool isString(const TaggedValue& value);
    
    // Create TaggedValue from String
    TaggedValue createTaggedString(const std::string& content);
    TaggedValue createTaggedString(String* str);
}

} // namespace smalltalk