#pragma once

#include "smalltalk_class.h"
#include "compiled_method.h"
#include <string>
#include <memory>

namespace smalltalk {

/**
 * Utility for parsing and compiling Smalltalk method source code
 */
class MethodCompiler {
public:
    /**
     * Parse and compile a Smalltalk method from source code.
     * 
     * @param methodSource The complete method source (e.g., "ensure: aBlock | result | ...")
     * @return A compiled method ready to be added to a class
     */
    static std::shared_ptr<CompiledMethod> compileMethod(const std::string& methodSource);
    
    /**
     * Parse and compile a Smalltalk method and add it to a class.
     * 
     * @param clazz The class to add the method to
     * @param methodSource The complete method source
     */
    static void addSmalltalkMethod(Class* clazz, const std::string& methodSource);
    
private:
    /**
     * Parse the method selector and parameters from source.
     * Returns the selector string and modifies the source to remove the signature.
     */
    static std::string parseMethodSignature(std::string& methodBody);
};

} // namespace smalltalk