#include "memory_manager.h"
#include "interpreter.h"
#include <iostream>
#include <string>

using namespace smalltalk;

void printBanner() {
    std::cout << "===============================" << std::endl;
    std::cout << "Smalltalk+ VM (C++ Implementation)" << std::endl;
    std::cout << "Version 0.1.0" << std::endl;
    std::cout << "===============================" << std::endl;
}

void printUsage() {
    std::cout << "Usage:" << std::endl;
    std::cout << "  smalltalk-vm [options] [image-file]" << std::endl;
    std::cout << std::endl;
    std::cout << "Options:" << std::endl;
    std::cout << "  --help    Show this help message" << std::endl;
    std::cout << "  --version Show version information" << std::endl;
}

int main(int argc, char** argv) {
    // Parse command line arguments
    std::string imageFile = "default.image";
    
    for (int i = 1; i < argc; i++) {
        std::string arg = argv[i];
        
        if (arg == "--help") {
            printBanner();
            printUsage();
            return 0;
        } else if (arg == "--version") {
            printBanner();
            return 0;
        } else {
            // Assume it's an image file
            imageFile = arg;
        }
    }
    
    try {
        printBanner();
        
        // Create memory manager
        MemoryManager memoryManager;
        
        std::cout << "Memory initialized with " << memoryManager.getTotalSpace() 
                  << " bytes total space" << std::endl;
        
        // Create interpreter
        Interpreter interpreter(memoryManager);
        
        std::cout << "Interpreter initialized" << std::endl;
        
        // In the future, we'd load an image file here
        std::cout << "Image file: " << imageFile << " (not loaded yet)" << std::endl;
        
        // For now, just create and execute a simple test method
        // This will be implemented in the full version
        
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        return 1;
    }
    
    return 0;
}