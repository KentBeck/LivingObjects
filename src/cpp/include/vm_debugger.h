#pragma once

#include "logger.h"
#include "context.h"
#include "compiled_method.h"
#include "tagged_value.h"
#include <string>
#include <vector>
#include <sstream>

namespace smalltalk {

class VMDebugger {
public:
    static VMDebugger& getInstance() {
        static VMDebugger instance;
        return instance;
    }

    void setDebugLevel(LogLevel level) {
        debugLevel = level;
        Logger::getInstance().setLevel(level);
    }

    void enableBytecodeTracing(bool enabled) {
        traceBytecode = enabled;
    }

    void enableStackTracing(bool enabled) {
        traceStack = enabled;
    }

    void enableMethodCalls(bool enabled) {
        traceMethodCalls = enabled;
    }

    void enableMemoryDebug(bool enabled) {
        traceMemory = enabled;
    }

    // Bytecode execution tracing
    void traceBytecodeExecution(const std::string& bytecode, uint32_t ip, 
                               const std::vector<TaggedValue>& stack) {
        if (!traceBytecode || debugLevel > LogLevel::DEBUG_LEVEL) return;

        std::stringstream ss;
        ss << "IP:" << ip << " " << bytecode;
        
        if (traceStack && !stack.empty()) {
            ss << " | Stack: [";
            for (size_t i = 0; i < stack.size(); ++i) {
                if (i > 0) ss << ", ";
                ss << taggedValueToString(stack[i]);
            }
            ss << "]";
        }
        
        LOG_BYTECODE_DEBUG(ss.str());
    }

    // Method call tracing
    void traceMethodEntry(const std::string& methodName, const std::string& className,
                         const std::vector<TaggedValue>& args) {
        if (!traceMethodCalls || debugLevel > LogLevel::DEBUG_LEVEL) return;

        std::stringstream ss;
        ss << "CALL: " << className << ">>" << methodName;
        
        if (!args.empty()) {
            ss << " with args: [";
            for (size_t i = 0; i < args.size(); ++i) {
                if (i > 0) ss << ", ";
                ss << taggedValueToString(args[i]);
            }
            ss << "]";
        }
        
        LOG_VM_DEBUG(ss.str());
    }

    void traceMethodExit(const std::string& methodName, const std::string& className,
                        TaggedValue result) {
        if (!traceMethodCalls || debugLevel > LogLevel::DEBUG_LEVEL) return;

        std::stringstream ss;
        ss << "RETURN: " << className << ">>" << methodName;
        ss << " -> " << taggedValueToString(result);
        
        LOG_VM_DEBUG(ss.str());
    }

    // Stack frame debugging
    void dumpStackFrame(MethodContext* context) {
        if (!context) return;

        std::stringstream ss;
        ss << "=== Stack Frame Dump ===" << std::endl;
        ss << "  IP: " << context->instructionPointer << std::endl;
        ss << "  Hash: " << context->header.hash << std::endl;
        ss << "  Self: " << taggedValueToString(context->self) << std::endl;
        ss << "  Sender: " << taggedValueToString(context->sender) << std::endl;
        
        LOG_VM_DEBUG(ss.str());
    }

    // Memory debugging
    void traceAllocation(const std::string& objectType, size_t size, void* address) {
        if (!traceMemory || debugLevel > LogLevel::DEBUG_LEVEL) return;

        std::stringstream ss;
        ss << "ALLOC: " << objectType << " (" << size << " bytes) at " << address;
        LOG_MEMORY_DEBUG(ss.str());
    }

    void traceDeallocation(const std::string& objectType, void* address) {
        if (!traceMemory || debugLevel > LogLevel::DEBUG_LEVEL) return;

        std::stringstream ss;
        ss << "DEALLOC: " << objectType << " at " << address;
        LOG_MEMORY_DEBUG(ss.str());
    }

    // Exception debugging
    void traceException(const std::string& exceptionType, const std::string& message,
                       const std::string& context) {
        std::stringstream ss;
        ss << "EXCEPTION: " << exceptionType << " - " << message;
        if (!context.empty()) {
            ss << " (in " << context << ")";
        }
        LOG_VM_ERROR(ss.str());
    }

    // Performance debugging
    void tracePerformance(const std::string& operation, double durationMs) {
        if (debugLevel > LogLevel::DEBUG_LEVEL) return;

        std::stringstream ss;
        ss << "PERF: " << operation << " took " << durationMs << "ms";
        LOG_VM_DEBUG(ss.str());
    }

    // General debugging utilities
    std::string taggedValueToString(TaggedValue value) {
        std::stringstream ss;
        
        if (value.isInteger()) {
            ss << value.asInteger();
        } else if (value.isBoolean()) {
            ss << (value.asBoolean() ? "true" : "false");
        } else if (value.isNil()) {
            ss << "nil";
        } else if (value.isPointer()) {
            ss << "<Object@" << std::hex << reinterpret_cast<uintptr_t>(value.asObject()) << ">";
        } else {
            ss << "<Unknown TaggedValue>";
        }
        
        return ss.str();
    }

    void enableAllTracing() {
        traceBytecode = true;
        traceStack = true;
        traceMethodCalls = true;
        traceMemory = true;
        setDebugLevel(LogLevel::DEBUG_LEVEL);
    }

    void disableAllTracing() {
        traceBytecode = false;
        traceStack = false;
        traceMethodCalls = false;
        traceMemory = false;
        setDebugLevel(LogLevel::INFO);
    }

private:
    VMDebugger() : debugLevel(LogLevel::INFO), traceBytecode(false), 
                   traceStack(false), traceMethodCalls(false), traceMemory(false) {}

    LogLevel debugLevel;
    bool traceBytecode;
    bool traceStack;
    bool traceMethodCalls;
    bool traceMemory;
};

// Convenient macros for VM debugging
#define VM_DEBUG_BYTECODE(bytecode, ip, stack) \
    VMDebugger::getInstance().traceBytecodeExecution(bytecode, ip, stack)

#define VM_DEBUG_METHOD_ENTRY(method, className, args) \
    VMDebugger::getInstance().traceMethodEntry(method, className, args)

#define VM_DEBUG_METHOD_EXIT(method, className, result) \
    VMDebugger::getInstance().traceMethodExit(method, className, result)

#define VM_DEBUG_EXCEPTION(type, message, context) \
    VMDebugger::getInstance().traceException(type, message, context)

#define VM_DEBUG_ALLOC(type, size, addr) \
    VMDebugger::getInstance().traceAllocation(type, size, addr)

#define VM_DEBUG_DEALLOC(type, addr) \
    VMDebugger::getInstance().traceDeallocation(type, addr)

#define VM_DEBUG_PERF(operation, duration) \
    VMDebugger::getInstance().tracePerformance(operation, duration)

} // namespace smalltalk