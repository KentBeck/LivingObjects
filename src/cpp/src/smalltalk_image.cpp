#include "smalltalk_image.h"
#include "simple_parser.h"
#include "simple_compiler.h"
#include "simple_vm.h"
#include "smalltalk_string.h"
#include "symbol.h"
#include "primitive_methods.h"

#include <filesystem>
#include <fstream>
#include <iostream>
#include <sstream>
#include <chrono>
#include <algorithm>
#include <ctime>
#include <iomanip>

namespace smalltalk {

SmalltalkImage::SmalltalkImage() 
    : imageVersion_("1.0.0")
    , creationTime_(ImageUtils::getCurrentTimestamp())
    , modificationTime_(creationTime_) {
}

SmalltalkImage::~SmalltalkImage() {
    // Clean up any allocated resources
}

// === Source Code Loading ===

bool SmalltalkImage::loadSourceFile(const std::string& filename) {
    try {
        std::ifstream file(filename);
        if (!file.is_open()) {
            std::cerr << "Error: Could not open file: " << filename << std::endl;
            return false;
        }
        
        std::stringstream buffer;
        buffer << file.rdbuf();
        std::string content = buffer.str();
        
        // Store the source file
        sourceFiles_.emplace_back(filename, content, filename);
        
        // Parse and load the content
        bool success = parseSourceCode(content, filename);
        if (success) {
            touch();
        }
        
        return success;
    } catch (const std::exception& e) {
        std::cerr << "Error loading source file " << filename << ": " << e.what() << std::endl;
        return false;
    }
}

bool SmalltalkImage::loadSourceFromString(const std::string& content, const std::string& name) {
    try {
        std::string filename = name.empty() ? "<string>" : name;
        
        // Store the source
        sourceFiles_.emplace_back(filename, content, "");
        
        // Parse and load the content
        bool success = parseSourceCode(content, filename);
        if (success) {
            touch();
        }
        
        return success;
    } catch (const std::exception& e) {
        std::cerr << "Error loading source from string: " << e.what() << std::endl;
        return false;
    }
}

bool SmalltalkImage::loadSourceDirectory(const std::string& directory) {
    try {
        if (!std::filesystem::exists(directory) || !std::filesystem::is_directory(directory)) {
            std::cerr << "Error: Directory does not exist: " << directory << std::endl;
            return false;
        }
        
        std::vector<std::string> sourceFiles = ImageUtils::findSourceFiles(directory);
        
        bool allSuccess = true;
        for (const auto& filename : sourceFiles) {
            if (!loadSourceFile(filename)) {
                std::cerr << "Warning: Failed to load source file: " << filename << std::endl;
                allSuccess = false;
            }
        }
        
        return allSuccess;
    } catch (const std::exception& e) {
        std::cerr << "Error loading source directory " << directory << ": " << e.what() << std::endl;
        return false;
    }
}

bool SmalltalkImage::loadSourceFiles(const std::vector<std::string>& filenames) {
    bool allSuccess = true;
    
    for (const auto& filename : filenames) {
        if (!loadSourceFile(filename)) {
            allSuccess = false;
        }
    }
    
    return allSuccess;
}

// === Image Persistence ===

bool SmalltalkImage::saveImage(const std::string& filename) const {
    try {
        std::ofstream file(filename, std::ios::binary);
        if (!file.is_open()) {
            std::cerr << "Error: Could not create image file: " << filename << std::endl;
            return false;
        }
        
        // Write header
        if (!writeHeader(file)) {
            std::cerr << "Error: Failed to write image header" << std::endl;
            return false;
        }
        
        // Write classes
        if (!writeClasses(file)) {
            std::cerr << "Error: Failed to write classes" << std::endl;
            return false;
        }
        
        // Write globals
        if (!writeGlobals(file)) {
            std::cerr << "Error: Failed to write globals" << std::endl;
            return false;
        }
        
        // Write metadata
        if (!writeMetadata(file)) {
            std::cerr << "Error: Failed to write metadata" << std::endl;
            return false;
        }
        
        std::cout << "Image saved successfully to: " << filename << std::endl;
        return true;
        
    } catch (const std::exception& e) {
        std::cerr << "Error saving image: " << e.what() << std::endl;
        return false;
    }
}

bool SmalltalkImage::loadImage(const std::string& filename) {
    try {
        std::ifstream file(filename, std::ios::binary);
        if (!file.is_open()) {
            std::cerr << "Error: Could not open image file: " << filename << std::endl;
            return false;
        }
        
        // Read and validate header
        ImageHeader header;
        if (!readHeader(file, header)) {
            std::cerr << "Error: Invalid image file header" << std::endl;
            return false;
        }
        
        // Clear current state
        clearImage();
        
        // Load image data
        creationTime_ = header.creationTime;
        modificationTime_ = header.modificationTime;
        
        // Read classes
        if (!readClasses(file, header.classCount)) {
            std::cerr << "Error: Failed to read classes" << std::endl;
            return false;
        }
        
        // Read globals
        if (!readGlobals(file, header.globalCount)) {
            std::cerr << "Error: Failed to read globals" << std::endl;
            return false;
        }
        
        // Read metadata
        if (!readMetadata(file, header.metadataCount)) {
            std::cerr << "Error: Failed to read metadata" << std::endl;
            return false;
        }
        
        std::cout << "Image loaded successfully from: " << filename << std::endl;
        std::cout << "  Classes: " << header.classCount << std::endl;
        std::cout << "  Globals: " << header.globalCount << std::endl;
        std::cout << "  Created: " << ImageUtils::formatTimestamp(header.creationTime) << std::endl;
        
        return true;
        
    } catch (const std::exception& e) {
        std::cerr << "Error loading image: " << e.what() << std::endl;
        return false;
    }
}

// === Image Management ===

void SmalltalkImage::initializeFreshImage() {
    clearImage();
    
    // Initialize core classes
    ClassUtils::initializeCoreClasses();
    
    // Initialize primitives
    auto& primitiveRegistry = PrimitiveRegistry::getInstance();
    primitiveRegistry.initializeCorePrimitives();
    
    // Add primitive methods to Integer class
    Class* integerClass = ClassUtils::getIntegerClass();
    IntegerClassSetup::addPrimitiveMethods(integerClass);
    
    // Set up global bindings
    // Set up basic globals
    setGlobal("Object", TaggedValue(static_cast<Object*>(ClassUtils::getObjectClass())));
    setGlobal("Class", TaggedValue(static_cast<Object*>(ClassUtils::getClassClass())));
    setGlobal("Integer", TaggedValue(static_cast<Object*>(ClassUtils::getIntegerClass())));
    setGlobal("String", TaggedValue(static_cast<Object*>(ClassUtils::getStringClass())));
    setGlobal("Symbol", TaggedValue(static_cast<Object*>(ClassUtils::getSymbolClass())));
    setGlobal("Boolean", TaggedValue(static_cast<Object*>(ClassUtils::getBooleanClass())));
    
    // Set basic metadata
    setMetadata("description", "Smalltalk image");
    setMetadata("created_by", "SmalltalkLSP");
    
    touch();
}

void SmalltalkImage::clearImage() {
    sourceFiles_.clear();
    globals_.clear();
    metadata_.clear();
    
    // Note: We don't clear the class registry here as it's global
    // In a full implementation, we'd need a more sophisticated cleanup
}

size_t SmalltalkImage::getClassCount() const {
    auto& registry = ClassRegistry::getInstance();
    return registry.getAllClasses().size();
}

size_t SmalltalkImage::getMethodCount() const {
    size_t methodCount = 0;
    auto& registry = ClassRegistry::getInstance();
    
    for (Class* clazz : registry.getAllClasses()) {
        methodCount += clazz->getMethodDictionary().size();
    }
    
    return methodCount;
}

size_t SmalltalkImage::getGlobalCount() const {
    return globals_.size();
}

// === Class Management ===

void SmalltalkImage::addClass(Class* clazz) {
    if (clazz) {
        auto& registry = ClassRegistry::getInstance();
        registry.registerClass(clazz->getName(), clazz);
        touch();
    }
}

std::vector<Class*> SmalltalkImage::getAllClasses() const {
    auto& registry = ClassRegistry::getInstance();
    return registry.getAllClasses();
}

Class* SmalltalkImage::findClass(const std::string& name) const {
    auto& registry = ClassRegistry::getInstance();
    return registry.getClass(name);
}

// === Global Variables ===

void SmalltalkImage::setGlobal(const std::string& name, TaggedValue value) {
    globals_[name] = value;
    touch();
}

TaggedValue SmalltalkImage::getGlobal(const std::string& name) const {
    auto it = globals_.find(name);
    if (it != globals_.end()) {
        return it->second;
    }
    return TaggedValue::nil();
}

bool SmalltalkImage::hasGlobal(const std::string& name) const {
    return globals_.find(name) != globals_.end();
}

std::vector<std::string> SmalltalkImage::getGlobalNames() const {
    std::vector<std::string> names;
    for (const auto& pair : globals_) {
        names.push_back(pair.first);
    }
    std::sort(names.begin(), names.end());
    return names;
}

// === Image Execution ===

TaggedValue SmalltalkImage::evaluate(const std::string& code) {
    try {
        SimpleParser parser(code);
        auto methodAST = parser.parseMethod();
        
        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);
        
        SimpleVM vm;
        return vm.execute(*compiledMethod);
        
    } catch (const std::exception& e) {
        std::cerr << "Error evaluating code: " << e.what() << std::endl;
        return TaggedValue::nil();
    }
}

TaggedValue SmalltalkImage::doIt(const std::string& expression) {
    return evaluate(expression);
}

// === Image Introspection ===

std::unordered_map<std::string, std::string> SmalltalkImage::getMetadata() const {
    return metadata_;
}

void SmalltalkImage::setMetadata(const std::string& key, const std::string& value) {
    metadata_[key] = value;
    touch();
}

// === Binary Serialization ===

bool SmalltalkImage::writeHeader(std::ofstream& file) const {
    ImageHeader header;
    header.magic = IMAGE_MAGIC;
    header.version = IMAGE_VERSION;
    header.creationTime = creationTime_;
    header.modificationTime = modificationTime_;
    header.classCount = static_cast<uint32_t>(getClassCount());
    header.methodCount = static_cast<uint32_t>(getMethodCount());
    header.globalCount = static_cast<uint32_t>(globals_.size());
    header.metadataCount = static_cast<uint32_t>(metadata_.size());
    header.dataOffset = sizeof(ImageHeader);
    
    file.write(reinterpret_cast<const char*>(&header), sizeof(header));
    return file.good();
}

bool SmalltalkImage::readHeader(std::ifstream& file, ImageHeader& header) {
    file.read(reinterpret_cast<char*>(&header), sizeof(header));
    
    if (!file.good()) {
        return false;
    }
    
    if (header.magic != IMAGE_MAGIC) {
        std::cerr << "Error: Invalid image file magic number" << std::endl;
        return false;
    }
    
    if (header.version != IMAGE_VERSION) {
        std::cerr << "Error: Unsupported image file version: " << header.version << std::endl;
        return false;
    }
    
    return true;
}

bool SmalltalkImage::writeClasses(std::ofstream& file) const {
    auto& registry = ClassRegistry::getInstance();
    auto classes = registry.getAllClasses();
    
    for (Class* clazz : classes) {
        // Write class name
        if (!writeString(file, clazz->getName())) {
            return false;
        }
        
        // Write superclass name (or empty string if none)
        std::string superclassName = "";
        if (clazz->getSuperclass()) {
            superclassName = clazz->getSuperclass()->getName();
        }
        if (!writeString(file, superclassName)) {
            return false;
        }
        
        // Write instance variables
        auto instanceVars = clazz->getInstanceVariables();
        uint32_t varCount = static_cast<uint32_t>(instanceVars.size());
        file.write(reinterpret_cast<const char*>(&varCount), sizeof(varCount));
        
        for (const auto& varName : instanceVars) {
            if (!writeString(file, varName)) {
                return false;
            }
        }
        
        // For now, we don't serialize methods (they would need special handling)
        uint32_t methodCount = 0;
        file.write(reinterpret_cast<const char*>(&methodCount), sizeof(methodCount));
    }
    
    return file.good();
}

bool SmalltalkImage::readClasses(std::ifstream& file, uint32_t count) {
    for (uint32_t i = 0; i < count; ++i) {
        std::string className;
        if (!readString(file, className)) {
            return false;
        }
        
        std::string superclassName;
        if (!readString(file, superclassName)) {
            return false;
        }
        
        // Read instance variables
        uint32_t varCount;
        file.read(reinterpret_cast<char*>(&varCount), sizeof(varCount));
        
        std::vector<std::string> instanceVars;
        for (uint32_t j = 0; j < varCount; ++j) {
            std::string varName;
            if (!readString(file, varName)) {
                return false;
            }
            instanceVars.push_back(varName);
        }
        
        // Skip methods for now
        uint32_t methodCount;
        file.read(reinterpret_cast<char*>(&methodCount), sizeof(methodCount));
        // In a full implementation, we'd read the methods here
        
        // Create the class if it doesn't exist
        auto& registry = ClassRegistry::getInstance();
        if (!registry.hasClass(className)) {
            Class* superclass = nullptr;
            if (!superclassName.empty()) {
                superclass = registry.getClass(superclassName);
            }
            
            Class* newClass = new Class(className, superclass);
            for (const auto& varName : instanceVars) {
                newClass->addInstanceVariable(varName);
            }
            
            registry.registerClass(className, newClass);
        }
    }
    
    return file.good();
}

bool SmalltalkImage::writeGlobals(std::ofstream& file) const {
    for (const auto& pair : globals_) {
        if (!writeString(file, pair.first)) {
            return false;
        }
        if (!writeTaggedValue(file, pair.second)) {
            return false;
        }
    }
    return file.good();
}

bool SmalltalkImage::readGlobals(std::ifstream& file, uint32_t count) {
    for (uint32_t i = 0; i < count; ++i) {
        std::string name;
        if (!readString(file, name)) {
            return false;
        }
        
        TaggedValue value;
        if (!readTaggedValue(file, value)) {
            return false;
        }
        
        globals_[name] = value;
    }
    return file.good();
}

bool SmalltalkImage::writeMetadata(std::ofstream& file) const {
    for (const auto& pair : metadata_) {
        if (!writeString(file, pair.first)) {
            return false;
        }
        if (!writeString(file, pair.second)) {
            return false;
        }
    }
    return file.good();
}

bool SmalltalkImage::readMetadata(std::ifstream& file, uint32_t count) {
    for (uint32_t i = 0; i < count; ++i) {
        std::string key;
        if (!readString(file, key)) {
            return false;
        }
        
        std::string value;
        if (!readString(file, value)) {
            return false;
        }
        
        metadata_[key] = value;
    }
    return file.good();
}

bool SmalltalkImage::writeTaggedValue(std::ofstream& file, const TaggedValue& value) const {
    // Write the raw tagged value for simple serialization
    // In a full implementation, this would need to handle object references
    uint64_t rawValue = value.rawValue();
    file.write(reinterpret_cast<const char*>(&rawValue), sizeof(rawValue));
    return file.good();
}

bool SmalltalkImage::readTaggedValue(std::ifstream& file, TaggedValue& value) {
    // Read the raw tagged value
    // In a full implementation, this would need to resolve object references
    uint64_t rawValue;
    file.read(reinterpret_cast<char*>(&rawValue), sizeof(rawValue));
    
    if (file.good()) {
        // This is a simplified reconstruction - in reality we'd need proper object reconstruction
        if ((rawValue & 0x03) == 0x03) { // Integer
            int32_t intValue = static_cast<int32_t>(rawValue >> 2);
            value = TaggedValue(intValue);
        } else if ((rawValue & 0x03) == 0x01) { // Special
            if (rawValue == TaggedValue::SPECIAL_NIL) {
                value = TaggedValue::nil();
            } else if (rawValue == TaggedValue::SPECIAL_TRUE) {
                value = TaggedValue::trueValue();
            } else if (rawValue == TaggedValue::SPECIAL_FALSE) {
                value = TaggedValue::falseValue();
            }
        } else {
            // For now, just set to nil for pointer values
            value = TaggedValue::nil();
        }
    }
    
    return file.good();
}

bool SmalltalkImage::writeString(std::ofstream& file, const std::string& str) const {
    uint32_t length = static_cast<uint32_t>(str.length());
    file.write(reinterpret_cast<const char*>(&length), sizeof(length));
    file.write(str.c_str(), length);
    return file.good();
}

bool SmalltalkImage::readString(std::ifstream& file, std::string& str) {
    uint32_t length;
    file.read(reinterpret_cast<char*>(&length), sizeof(length));
    
    if (!file.good() || length > 1000000) { // Sanity check
        return false;
    }
    
    str.resize(length);
    file.read(&str[0], length);
    return file.good();
}

// === Source Parsing ===

bool SmalltalkImage::parseSourceCode(const std::string& source, const std::string& filename) {
    // For now, just evaluate the source as expressions
    // In a full implementation, this would parse class definitions, method definitions, etc.
    
    try {
        // Split source into lines and evaluate each non-empty line
        std::istringstream stream(source);
        std::string line;
        int lineNumber = 0;
        
        while (std::getline(stream, line)) {
            lineNumber++;
            
            // Skip empty lines and comments
            if (line.empty() || line[0] == '"' || line.find_first_not_of(" \t") == std::string::npos) {
                continue;
            }
            
            // Try to evaluate the line
            try {
                TaggedValue result = evaluate(line);
                // For now, just ignore the result
                (void)result;
            } catch (const std::exception& e) {
                std::cerr << "Warning: Error in " << filename << " line " << lineNumber 
                         << ": " << e.what() << std::endl;
                // Continue processing other lines
            }
        }
        
        return true;
        
    } catch (const std::exception& e) {
        std::cerr << "Error parsing source " << filename << ": " << e.what() << std::endl;
        return false;
    }
}

bool SmalltalkImage::parseClassDefinition(const std::string& source) {
    // TODO: Implement class definition parsing
    // Format: "Object subclass: #MyClass instanceVariableNames: 'var1 var2'"
    (void)source;
    return false;
}

bool SmalltalkImage::parseMethodDefinition(const std::string& source, Class* targetClass) {
    // TODO: Implement method definition parsing
    // Format: "methodName: arg | temp | temp := arg. ^temp"
    (void)source;
    (void)targetClass;
    return false;
}

void SmalltalkImage::touch() {
    modificationTime_ = ImageUtils::getCurrentTimestamp();
}

// === ImageManager Implementation ===

ImageManager& ImageManager::getInstance() {
    static ImageManager instance;
    return instance;
}

void ImageManager::setCurrentImage(std::unique_ptr<SmalltalkImage> image) {
    currentImage_ = std::move(image);
}

void ImageManager::createFreshImage() {
    currentImage_ = std::make_unique<SmalltalkImage>();
    currentImage_->initializeFreshImage();
}

bool ImageManager::loadImageFromFile(const std::string& filename) {
    auto newImage = std::make_unique<SmalltalkImage>();
    if (newImage->loadImage(filename)) {
        currentImage_ = std::move(newImage);
        return true;
    }
    return false;
}

bool ImageManager::saveImageToFile(const std::string& filename) {
    if (currentImage_) {
        return currentImage_->saveImage(filename);
    }
    return false;
}

bool ImageManager::loadSourceFiles(const std::vector<std::string>& filenames) {
    if (!currentImage_) {
        createFreshImage();
    }
    return currentImage_->loadSourceFiles(filenames);
}

bool ImageManager::loadSourceDirectory(const std::string& directory) {
    if (!currentImage_) {
        createFreshImage();
    }
    return currentImage_->loadSourceDirectory(directory);
}

// === ImageUtils Implementation ===

namespace ImageUtils {
    
    std::unique_ptr<SmalltalkImage> createStandardImage() {
        auto image = std::make_unique<SmalltalkImage>();
        image->initializeFreshImage();
        
        // Load standard library if available
        if (std::filesystem::exists("src") && std::filesystem::is_directory("src")) {
            image->loadSourceDirectory("src");
        }
        
        return image;
    }
    
    std::vector<std::string> findSourceFiles(const std::string& directory) {
        std::vector<std::string> sourceFiles;
        
        try {
            for (const auto& entry : std::filesystem::recursive_directory_iterator(directory)) {
                if (entry.is_regular_file() && entry.path().extension() == ".st") {
                    sourceFiles.push_back(entry.path().string());
                }
            }
            
            // Sort for consistent loading order
            std::sort(sourceFiles.begin(), sourceFiles.end());
            
        } catch (const std::exception& e) {
            std::cerr << "Error scanning directory " << directory << ": " << e.what() << std::endl;
        }
        
        return sourceFiles;
    }
    
    bool isValidImageFile(const std::string& filename) {
        try {
            std::ifstream file(filename, std::ios::binary);
            if (!file.is_open()) {
                return false;
            }
            
            uint32_t magic;
            file.read(reinterpret_cast<char*>(&magic), sizeof(magic));
            
            return file.good() && magic == SmalltalkImage::IMAGE_MAGIC;
            
        } catch (...) {
            return false;
        }
    }
    
    bool getImageInfo(const std::string& filename, std::string& version, 
                     uint64_t& creationTime, uint32_t& classCount) {
        try {
            std::ifstream file(filename, std::ios::binary);
            if (!file.is_open()) {
                return false;
            }
            
            SmalltalkImage::ImageHeader header;
            file.read(reinterpret_cast<char*>(&header), sizeof(header));
            
            if (!file.good() || header.magic != SmalltalkImage::IMAGE_MAGIC) {
                return false;
            }
            
            version = "1.0.0"; // TODO: Read from metadata
            creationTime = header.creationTime;
            classCount = header.classCount;
            
            return true;
            
        } catch (...) {
            return false;
        }
    }
    
    std::string formatTimestamp(uint64_t timestamp) {
        std::time_t time_t = static_cast<std::time_t>(timestamp);
        
        std::stringstream ss;
        ss << std::put_time(std::localtime(&time_t), "%Y-%m-%d %H:%M:%S");
        return ss.str();
    }
    
    uint64_t getCurrentTimestamp() {
        auto now = std::chrono::system_clock::now();
        return std::chrono::system_clock::to_time_t(now);
    }
}

} // namespace smalltalk