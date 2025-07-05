#pragma once

namespace smalltalk {

class SmalltalkVM {
public:
    // Initialize the entire Smalltalk system
    // This MUST be called before any Smalltalk operations
    static void initialize();
    
    // Check if the VM has been initialized
    static bool isInitialized();
    
    // Shutdown the VM (cleanup)
    static void shutdown();

private:
    static bool initialized;
};

} // namespace smalltalk
