#pragma once

#include "object.h"
#include "smalltalk_class.h"
#include "tagged_value.h"
#include <string>
#include <vector>
#include <unordered_map>
#include <memory>
#include <fstream>

namespace smalltalk {

// Forward declarations
class CompiledMethod;
class Symbol;

/**
 * Represents a single Smalltalk source file to be loaded into the image
 */
struct SourceFile {
    std::string filename;
    std::string content;
    std::string relativePath;
    
    SourceFile(const std::string& name, const std::string& text, const std::string& path = "")
        : filename(name), content(text), relativePath(path) {}
};

/**
 * Represents the current state of the Smalltalk system including all classes,
 * methods, and global variables. Can be saved to and loaded from binary files.
 */
class SmalltalkImage {
public:
    SmalltalkImage();
    ~SmalltalkImage();
    
    // === Source Code Loading ===
    
    // Load a single Smalltalk source file
    bool loadSourceFile(const std::string& filename);
    bool loadSourceFromString(const std::string& content, const std::string& name = "");
    
    // Load all .st files from a directory
    bool loadSourceDirectory(const std::string& directory);
    
    // Load multiple source files
    bool loadSourceFiles(const std::vector<std::string>& filenames);
    
    // Get list of loaded source files
    const std::vector<SourceFile>& getSourceFiles() const { return sourceFiles_; }
    
    // === Image Persistence ===
    
    // Save current image state to binary file
    bool saveImage(const std::string& filename) const;
    
    // Load image state from binary file
    bool loadImage(const std::string& filename);
    
    // === Image Management ===
    
    // Initialize a fresh image with core classes
    void initializeFreshImage();
    
    // Clear all loaded code (reset to empty state)
    void clearImage();
    
    // Get image statistics
    size_t getClassCount() const;
    size_t getMethodCount() const;
    size_t getGlobalCount() const;
    
    // === Class Management ===
    
    // Register a class in the image
    void addClass(Class* clazz);
    
    // Get all classes in the image
    std::vector<Class*> getAllClasses() const;
    
    // Find a class by name
    Class* findClass(const std::string& name) const;
    
    // === Global Variables ===
    
    // Set a global variable
    void setGlobal(const std::string& name, TaggedValue value);
    
    // Get a global variable
    TaggedValue getGlobal(const std::string& name) const;
    
    // Check if global exists
    bool hasGlobal(const std::string& name) const;
    
    // Get all global names
    std::vector<std::string> getGlobalNames() const;
    
    // === Image Execution ===
    
    // Evaluate Smalltalk code in the context of this image
    TaggedValue evaluate(const std::string& code);
    
    // Execute a "do it" (arbitrary Smalltalk expression)
    TaggedValue doIt(const std::string& expression);
    
    // === Image Introspection ===
    
    // Get image version information
    std::string getVersion() const { return imageVersion_; }
    void setVersion(const std::string& version) { imageVersion_ = version; }
    
    // Get creation timestamp
    uint64_t getCreationTime() const { return creationTime_; }
    
    // Get last modification time
    uint64_t getModificationTime() const { return modificationTime_; }
    
    // Get image metadata
    std::unordered_map<std::string, std::string> getMetadata() const;
    void setMetadata(const std::string& key, const std::string& value);
    
public:
    // === Binary Serialization Constants ===
    
    // Magic number for image files
    static constexpr uint32_t IMAGE_MAGIC = 0x53544C4B; // "STLK" in hex
    static constexpr uint32_t IMAGE_VERSION = 1;
    
    // Binary format structure
    struct ImageHeader {
        uint32_t magic;           // Magic number for validation
        uint32_t version;         // Image format version
        uint64_t creationTime;    // When image was created
        uint64_t modificationTime; // When image was last modified
        uint32_t classCount;      // Number of classes
        uint32_t methodCount;     // Number of methods
        uint32_t globalCount;     // Number of globals
        uint32_t metadataCount;   // Number of metadata entries
        uint64_t dataOffset;      // Offset to start of object data
    };

private:
    // === Internal State ===
    
    std::vector<SourceFile> sourceFiles_;
    std::unordered_map<std::string, TaggedValue> globals_;
    std::unordered_map<std::string, std::string> metadata_;
    
    std::string imageVersion_;
    uint64_t creationTime_;
    uint64_t modificationTime_;
    
    // Serialization helpers
    bool writeHeader(std::ofstream& file) const;
    bool readHeader(std::ifstream& file, ImageHeader& header);
    
    bool writeClasses(std::ofstream& file) const;
    bool readClasses(std::ifstream& file, uint32_t count);
    
    bool writeGlobals(std::ofstream& file) const;
    bool readGlobals(std::ifstream& file, uint32_t count);
    
    bool writeMetadata(std::ofstream& file) const;
    bool readMetadata(std::ifstream& file, uint32_t count);
    
    bool writeTaggedValue(std::ofstream& file, const TaggedValue& value) const;
    bool readTaggedValue(std::ifstream& file, TaggedValue& value);
    
    bool writeString(std::ofstream& file, const std::string& str) const;
    bool readString(std::ifstream& file, std::string& str);
    
    // === Source Parsing ===
    
    // Parse Smalltalk source code and add to image
    bool parseSourceCode(const std::string& source, const std::string& filename);
    
    // Parse a class definition
    bool parseClassDefinition(const std::string& source);
    
    // Parse a method definition
    bool parseMethodDefinition(const std::string& source, Class* targetClass);
    
    // Update modification time
    void touch();
};

/**
 * ImageManager provides global access to the current Smalltalk image
 */
class ImageManager {
public:
    // Singleton access
    static ImageManager& getInstance();
    
    // Get the current image
    SmalltalkImage* getCurrentImage() { return currentImage_.get(); }
    
    // Set the current image
    void setCurrentImage(std::unique_ptr<SmalltalkImage> image);
    
    // Create a new fresh image
    void createFreshImage();
    
    // Load an image from file
    bool loadImageFromFile(const std::string& filename);
    
    // Save current image to file
    bool saveImageToFile(const std::string& filename);
    
    // Load source files into current image
    bool loadSourceFiles(const std::vector<std::string>& filenames);
    bool loadSourceDirectory(const std::string& directory);
    
private:
    ImageManager() = default;
    std::unique_ptr<SmalltalkImage> currentImage_;
};

/**
 * Utility functions for working with Smalltalk images
 */
namespace ImageUtils {
    // Create a standard image with common classes loaded
    std::unique_ptr<SmalltalkImage> createStandardImage();
    
    // Load all .st files from a directory recursively
    std::vector<std::string> findSourceFiles(const std::string& directory);
    
    // Validate that a file is a valid Smalltalk image
    bool isValidImageFile(const std::string& filename);
    
    // Get image file information without fully loading it
    bool getImageInfo(const std::string& filename, std::string& version, 
                     uint64_t& creationTime, uint32_t& classCount);
    
    // Convert timestamp to human-readable string
    std::string formatTimestamp(uint64_t timestamp);
    
    // Get current timestamp
    uint64_t getCurrentTimestamp();
}

} // namespace smalltalk