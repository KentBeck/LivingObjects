#include "smalltalk_string.h"
#include "smalltalk_class.h"
#include <algorithm>
#include <cctype>
#include <functional>

namespace smalltalk {

String::String(const std::string& content, Class* stringClass) 
    : Object(ObjectType::OBJECT, sizeof(String), stringClass), content_(content) {
}

String::String(const char* content, Class* stringClass) 
    : Object(ObjectType::OBJECT, sizeof(String), stringClass), content_(content ? content : "") {
}

String* String::concatenate(const String* other) const {
    if (other == nullptr) {
        return new String(content_, getClass());
    }
    return new String(content_ + other->content_, getClass());
}

String* String::concatenate(const std::string& other) const {
    return new String(content_ + other, getClass());
}

bool String::equals(const String* other) const {
    if (other == nullptr) {
        return false;
    }
    return content_ == other->content_;
}

bool String::equals(const std::string& other) const {
    return content_ == other;
}

char String::at(size_t index) const {
    if (index >= content_.size()) {
        throw std::runtime_error("String index out of bounds: " + std::to_string(index));
    }
    return content_[index];
}

String* String::substring(size_t start, size_t length) const {
    if (start >= content_.size()) {
        return new String("", getClass());
    }
    
    size_t actualLength = std::min(length, content_.size() - start);
    return new String(content_.substr(start, actualLength), getClass());
}

String* String::toLowerCase() const {
    std::string lower = content_;
    std::transform(lower.begin(), lower.end(), lower.begin(), ::tolower);
    return new String(lower, getClass());
}

String* String::toUpperCase() const {
    std::string upper = content_;
    std::transform(upper.begin(), upper.end(), upper.begin(), ::toupper);
    return new String(upper, getClass());
}

int String::indexOf(char ch) const {
    size_t pos = content_.find(ch);
    return pos == std::string::npos ? -1 : static_cast<int>(pos);
}

int String::indexOf(const std::string& substr) const {
    size_t pos = content_.find(substr);
    return pos == std::string::npos ? -1 : static_cast<int>(pos);
}

bool String::contains(const std::string& substr) const {
    return content_.find(substr) != std::string::npos;
}

std::string String::toString() const {
    return "'" + content_ + "'";
}

size_t String::hash() const {
    return std::hash<std::string>{}(content_);
}

bool String::operator==(const String& other) const {
    return content_ == other.content_;
}

bool String::operator!=(const String& other) const {
    return content_ != other.content_;
}

bool String::operator<(const String& other) const {
    return content_ < other.content_;
}

bool String::operator>(const String& other) const {
    return content_ > other.content_;
}

bool String::operator<=(const String& other) const {
    return content_ <= other.content_;
}

bool String::operator>=(const String& other) const {
    return content_ >= other.content_;
}

// StringUtils implementation
namespace StringUtils {
    
    String* createString(const std::string& content) {
        Class* stringClass = ClassUtils::getStringClass();
        return new String(content, stringClass);
    }
    
    String* createString(const char* content) {
        Class* stringClass = ClassUtils::getStringClass();
        return new String(content, stringClass);
    }
    
    String* asString(const TaggedValue& value) {
        if (!value.isPointer()) {
            return nullptr;
        }
        
        try {
            return static_cast<String*>(value.asObject());
        } catch (...) {
            return nullptr;
        }
    }
    
    bool isString(const TaggedValue& value) {
        if (!value.isPointer()) {
            return false;
        }
        
        try {
            Object* obj = value.asObject();
            Class* stringClass = ClassUtils::getStringClass();
            return obj->getClass() == stringClass;
        } catch (...) {
            return false;
        }
    }
    
    TaggedValue createTaggedString(const std::string& content) {
        String* str = createString(content);
        return TaggedValue(static_cast<Object*>(str));
    }
    
    TaggedValue createTaggedString(String* str) {
        return TaggedValue(static_cast<Object*>(str));
    }
}

} // namespace smalltalk