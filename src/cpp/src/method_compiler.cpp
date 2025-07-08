#include "method_compiler.h"
#include "simple_parser.h"
#include "simple_compiler.h"
#include "symbol.h"
#include <regex>
#include <sstream>

namespace smalltalk {

std::shared_ptr<CompiledMethod> MethodCompiler::compileMethod(const std::string& methodSource) {
    // Make a copy of the source to modify
    std::string source = methodSource;
    
    // Parse the method signature to get the selector and modify the source
    std::string selector = parseMethodSignature(source);
    
    // The remaining source is the method body
    // Parse it as a method body (not an expression)
    SimpleParser parser(source);
    auto methodAST = parser.parseMethod();
    
    // Compile to bytecode
    SimpleCompiler compiler;
    auto compiledMethod = compiler.compile(*methodAST);
    
    return compiledMethod;
}

void MethodCompiler::addSmalltalkMethod(Class* clazz, const std::string& methodSource) {
    // Make a copy of the source to modify
    std::string source = methodSource;
    
    // Parse the method signature to get the selector
    std::string selector = parseMethodSignature(source);
    
    // Compile the method
    auto compiledMethod = compileMethod(methodSource);
    
    // Create selector symbol
    Symbol* selectorSymbol = Symbol::intern(selector);
    
    // Add method to class
    clazz->addMethod(selectorSymbol, compiledMethod);
}

std::string MethodCompiler::parseMethodSignature(std::string& methodBody) {
    // Simple parsing for method signatures
    // Handle patterns like:
    // "ensure: aBlock"
    // "on: exceptionClass do: handlerBlock"
    // "value"
    // "value: anArg"
    
    std::istringstream iss(methodBody);
    std::string line;
    std::getline(iss, line);
    
    // Find the first line that contains the method signature
    std::string selector;
    std::vector<std::string> parameters;
    size_t bodyStart = 0;
    
    // Look for a line that contains method signature patterns
    if (line.find(':') != std::string::npos) {
        // Keyword message - extract selector parts and parameters
        std::regex pattern(R"((\w+):\s*(\w+))");
        std::sregex_iterator iter(line.begin(), line.end(), pattern);
        std::sregex_iterator end;
        
        std::string fullSelector;
        for (; iter != end; ++iter) {
            if (!fullSelector.empty()) fullSelector += ":";
            fullSelector += iter->str(1) + ":";
            parameters.push_back(iter->str(2));
        }
        selector = fullSelector;
        bodyStart = line.length() + 1; // +1 for newline
    } else {
        // Unary message - just the method name
        std::regex pattern(R"(^\s*(\w+))");
        std::smatch match;
        if (std::regex_search(line, match, pattern)) {
            selector = match[1];
            bodyStart = line.length() + 1;
        }
    }
    
    // Remove the signature line from methodBody
    if (bodyStart < methodBody.length()) {
        methodBody = methodBody.substr(bodyStart);
    } else {
        methodBody = "";
    }
    
    // If we have parameters, add them as temporary variables to the method body
    if (!parameters.empty()) {
        std::string parameterDecl = "| ";
        for (const auto& param : parameters) {
            parameterDecl += param + " ";
        }
        parameterDecl += "|\n";
        
        // Check if method body already has temporary variables
        if (methodBody.find('|') != std::string::npos) {
            // Method already has temp vars, we need to merge them
            size_t firstPipe = methodBody.find('|');
            size_t secondPipe = methodBody.find('|', firstPipe + 1);
            if (secondPipe != std::string::npos) {
                // Extract existing temp vars
                std::string existingTemps = methodBody.substr(firstPipe + 1, secondPipe - firstPipe - 1);
                std::string newTemps = "| ";
                for (const auto& param : parameters) {
                    newTemps += param + " ";
                }
                // Trim leading/trailing spaces from existing temps and normalize spaces
                size_t start = existingTemps.find_first_not_of(' ');
                size_t end = existingTemps.find_last_not_of(' ');
                if (start != std::string::npos && end != std::string::npos) {
                    existingTemps = existingTemps.substr(start, end - start + 1);
                    // Replace multiple spaces with single spaces
                    std::string normalizedTemps;
                    bool lastWasSpace = false;
                    for (char c : existingTemps) {
                        if (c == ' ') {
                            if (!lastWasSpace) {
                                normalizedTemps += c;
                                lastWasSpace = true;
                            }
                        } else {
                            normalizedTemps += c;
                            lastWasSpace = false;
                        }
                    }
                    newTemps += normalizedTemps + " ";
                }
                newTemps += "|";
                
                // Replace the temp var declaration
                methodBody = newTemps + methodBody.substr(secondPipe + 1);
            } else {
                // Malformed temp vars, just prepend parameters
                methodBody = parameterDecl + methodBody;
            }
        } else {
            // No existing temp vars, just add parameters
            methodBody = parameterDecl + methodBody;
        }
    }
    
    return selector;
}

} // namespace smalltalk