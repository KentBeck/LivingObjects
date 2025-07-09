#include "logger.h"
#include "vm_debugger.h"
#include <iostream>

int main() {
    using namespace smalltalk;
    
    std::cout << "Testing logging infrastructure..." << std::endl;
    
    // Test basic logging
    Logger::getInstance().setLevel(LogLevel::DEBUG_LEVEL);
    Logger::getInstance().setConsoleOutput(true);
    
    LOG_INFO("Basic logging test");
    LOG_DEBUG("Debug message");
    LOG_WARN("Warning message");
    LOG_ERROR("Error message");
    
    // Test VM-specific logging
    LOG_VM_INFO("VM initialization");
    LOG_VM_DEBUG("VM debug message");
    LOG_BYTECODE_DEBUG("Bytecode execution");
    LOG_MEMORY_DEBUG("Memory allocation");
    LOG_GC_DEBUG("GC debug");
    
    // Test VM debugger
    VMDebugger::getInstance().setDebugLevel(LogLevel::DEBUG_LEVEL);
    VMDebugger::getInstance().enableAllTracing();
    
    std::vector<smalltalk::TaggedValue> args;
    args.push_back(TaggedValue(42));
    args.push_back(TaggedValue(true));
    
    VM_DEBUG_METHOD_ENTRY("testMethod", "TestClass", args);
    VM_DEBUG_METHOD_EXIT("testMethod", "TestClass", TaggedValue(100));
    
    VM_DEBUG_EXCEPTION("TestException", "Test exception message", "TestContext");
    
    VM_DEBUG_ALLOC("TestObject", 64, reinterpret_cast<void*>(0x12345678));
    VM_DEBUG_DEALLOC("TestObject", reinterpret_cast<void*>(0x12345678));
    
    VM_DEBUG_PERF("TestOperation", 15.5);
    
    std::cout << "Logging test completed!" << std::endl;
    
    return 0;
}