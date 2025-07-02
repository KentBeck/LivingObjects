#include "context.h"
#include "memory_manager.h"
#include "object.h"

#include <iostream>
#include <cstring>

using namespace smalltalk;

int main() {
    std::cout << "=== BlockContext Memory Layout Debug ===" << std::endl;
    
    // Print sizes
    std::cout << "sizeof(Object): " << sizeof(Object) << std::endl;
    std::cout << "sizeof(MethodContext): " << sizeof(MethodContext) << std::endl;
    std::cout << "sizeof(BlockContext): " << sizeof(BlockContext) << std::endl;
    std::cout << "sizeof(Object*): " << sizeof(Object*) << std::endl;
    
    // Print field offsets
    std::cout << "\n=== Field Offsets ===" << std::endl;
    BlockContext dummy(0, 0, nullptr, nullptr, nullptr);
    
    // Calculate offsets
    char* base = reinterpret_cast<char*>(&dummy);
    char* home_ptr = reinterpret_cast<char*>(&dummy.home);
    char* header_ptr = reinterpret_cast<char*>(&dummy.header);
    
    std::cout << "Object.header offset: " << (header_ptr - base) << std::endl;
    std::cout << "BlockContext.home offset: " << (home_ptr - base) << std::endl;
    
    // Test constructor behavior
    std::cout << "\n=== Constructor Test ===" << std::endl;
    
    MemoryManager memoryManager;
    
    // Create a dummy home context
    MethodContext* homeContext = memoryManager.allocateMethodContext(2, 123, nullptr, nullptr);
    std::cout << "Created home context at: " << homeContext << std::endl;
    
    // Test direct construction
    std::cout << "\n--- Direct Construction Test ---" << std::endl;
    char buffer[sizeof(BlockContext) + 64]; // Extra space
    memset(buffer, 0, sizeof(buffer));
    
    BlockContext* directBlock = reinterpret_cast<BlockContext*>(buffer);
    new (directBlock) BlockContext(2, 456, nullptr, nullptr, homeContext);
    
    std::cout << "Direct construction:" << std::endl;
    std::cout << "  Block at: " << directBlock << std::endl;
    std::cout << "  Block.home: " << directBlock->home << std::endl;
    std::cout << "  Expected home: " << homeContext << std::endl;
    std::cout << "  Block.header.type: " << directBlock->header.type << std::endl;
    std::cout << "  Block.header.hash: " << directBlock->header.hash << std::endl;
    
    // Test memory manager allocation
    std::cout << "\n--- Memory Manager Allocation Test ---" << std::endl;
    BlockContext* allocatedBlock = memoryManager.allocateBlockContext(2, 789, nullptr, nullptr, homeContext);
    
    std::cout << "Memory manager allocation:" << std::endl;
    std::cout << "  Block at: " << allocatedBlock << std::endl;
    std::cout << "  Block.home: " << allocatedBlock->home << std::endl;
    std::cout << "  Expected home: " << homeContext << std::endl;
    std::cout << "  Block.header.type: " << allocatedBlock->header.type << std::endl;
    std::cout << "  Block.header.hash: " << allocatedBlock->header.hash << std::endl;
    
    // Examine memory contents
    std::cout << "\n=== Memory Contents ===" << std::endl;
    char* blockPtr = reinterpret_cast<char*>(allocatedBlock);
    std::cout << "First 64 bytes of block memory:" << std::endl;
    for (int i = 0; i < 64; i += 8) {
        std::cout << "  [" << i << "]: ";
        for (int j = 0; j < 8 && i + j < 64; j++) {
            printf("%02x ", static_cast<unsigned char>(blockPtr[i + j]));
        }
        std::cout << std::endl;
    }
    
    // Test what happens if we manually set the home field
    std::cout << "\n=== Manual Field Setting Test ===" << std::endl;
    allocatedBlock->home = homeContext;
    std::cout << "After manual setting:" << std::endl;
    std::cout << "  Block.home: " << allocatedBlock->home << std::endl;
    
    // Check if it persists
    std::cout << "Checking if it persists..." << std::endl;
    std::cout << "  Block.home: " << allocatedBlock->home << std::endl;
    
    return 0;
}