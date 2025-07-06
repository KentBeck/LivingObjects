#include "tagged_value.h"
#include "object.h"
#include "smalltalk_class.h"
#include "smalltalk_string.h"
#include "symbol.h"

namespace smalltalk
{

    bool TaggedValue::isObjectOfClass(Class *clazz) const
    {
        if (!isPointer())
        {
            return false;
        }

        try
        {
            Object *obj = asObject();
            return obj->getClass() == clazz;
        }
        catch (...)
        {
            return false;
        }
    }

    Class *TaggedValue::getClass() const
    {
        if (isInteger())
        {
            return ClassUtils::getIntegerClass();
        }
        else if (isBoolean())
        {
            return ClassUtils::getBooleanClass();
        }
        else if (isNil())
        {
            // nil is a special case - it's the sole instance of UndefinedObject
            // For now, we'll return Object class
            return ClassUtils::getObjectClass();
        }
        else if (isPointer())
        {
            try
            {
                // Try as Symbol first
                asSymbol(); // Just check if it's a symbol
                return ClassUtils::getSymbolClass();
            }
            catch (...)
            {
                try
                {
                    // Try as String
                    Object *obj = asObject();
                    if (StringUtils::isString(*this))
                    {
                        return ClassUtils::getStringClass();
                    }
                    // Try as generic Object
                    return obj->getClass();
                }
                catch (...)
                {
                    return nullptr;
                }
            }
        }
        return nullptr;
    }

    TaggedValue TaggedValue::fromObject(Object *object)
    {
        if (object == nullptr)
        {
            return TaggedValue::nil();
        }

        // Check if this is a tagged value wrapper
        if (object->header.hasFlag(ObjectFlag::TAGGED_VALUE_WRAPPER))
        {
            // Check if it's a boxed integer
            if (object->getClass() == ClassUtils::getIntegerClass()) {
                TaggedValue* valueSlot = reinterpret_cast<TaggedValue*>(reinterpret_cast<char*>(object) + sizeof(Object));
                return *valueSlot;
            }
            // Check if it's a boxed boolean
            if (object->getClass() == ClassUtils::getBooleanClass()) {
                TaggedValue* valueSlot = reinterpret_cast<TaggedValue*>(reinterpret_cast<char*>(object) + sizeof(Object));
                return *valueSlot;
            }
            // If it's a generic tagged value wrapper, extract the original TaggedValue
            TaggedValue *storedValue = reinterpret_cast<TaggedValue *>
                (reinterpret_cast<char *>(object) + sizeof(Object));
            return *storedValue;
        }

        // Regular object pointer
        return TaggedValue(object);
    }

} // namespace smalltalk