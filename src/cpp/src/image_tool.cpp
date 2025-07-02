#include "smalltalk_image.h"
#include "smalltalk_string.h"
#include "smalltalk_class.h"
#include <iostream>
#include <vector>
#include <string>

using namespace smalltalk;

void printUsage(const char* programName) {
    std::cout << "Smalltalk Image Tool\n";
    std::cout << "Usage: " << programName << " <command> [options]\n\n";
    std::cout << "Commands:\n";
    std::cout << "  create <image_file>                    Create a fresh image\n";
    std::cout << "  load <source_dir> <image_file>         Load source files and save image\n";
    std::cout << "  loadfiles <file1> <file2> ... <image> Load specific files and save image\n";
    std::cout << "  info <image_file>                      Show image information\n";
    std::cout << "  run <image_file> <expression>          Load image and evaluate expression\n";
    std::cout << "  eval <expression>                      Evaluate expression in fresh image\n\n";
    std::cout << "Examples:\n";
    std::cout << "  " << programName << " create my_image.st\n";
    std::cout << "  " << programName << " load src/ my_image.st\n";
    std::cout << "  " << programName << " loadfiles Point.st Rectangle.st my_image.st\n";
    std::cout << "  " << programName << " info my_image.st\n";
    std::cout << "  " << programName << " run my_image.st \"3 + 4\"\n";
    std::cout << "  " << programName << " eval \"'hello world'\"\n";
}

void printResult(const TaggedValue& result) {
    if (result.isInteger()) {
        std::cout << "=> " << result.asInteger() << std::endl;
    } else if (result.isBoolean()) {
        std::cout << "=> " << (result.asBoolean() ? "true" : "false") << std::endl;
    } else if (result.isNil()) {
        std::cout << "=> nil" << std::endl;
    } else if (StringUtils::isString(result)) {
        String* str = StringUtils::asString(result);
        std::cout << "=> " << str->toString() << std::endl;
    } else {
        std::cout << "=> " << result << std::endl;
    }
}

int createImage(const std::string& imageFile) {
    std::cout << "Creating fresh image: " << imageFile << std::endl;
    
    auto& manager = ImageManager::getInstance();
    manager.createFreshImage();
    
    if (manager.saveImageToFile(imageFile)) {
        auto* image = manager.getCurrentImage();
        std::cout << "Image created successfully!" << std::endl;
        std::cout << "  Classes: " << image->getClassCount() << std::endl;
        std::cout << "  Globals: " << image->getGlobalCount() << std::endl;
        return 0;
    } else {
        std::cerr << "Failed to create image" << std::endl;
        return 1;
    }
}

int loadSourceDirectory(const std::string& sourceDir, const std::string& imageFile) {
    std::cout << "Loading source directory: " << sourceDir << std::endl;
    std::cout << "Creating image: " << imageFile << std::endl;
    
    auto& manager = ImageManager::getInstance();
    manager.createFreshImage();
    
    if (!manager.loadSourceDirectory(sourceDir)) {
        std::cerr << "Failed to load source directory" << std::endl;
        return 1;
    }
    
    if (manager.saveImageToFile(imageFile)) {
        auto* image = manager.getCurrentImage();
        std::cout << "Image created successfully!" << std::endl;
        std::cout << "  Source files: " << image->getSourceFiles().size() << std::endl;
        std::cout << "  Classes: " << image->getClassCount() << std::endl;
        std::cout << "  Methods: " << image->getMethodCount() << std::endl;
        std::cout << "  Globals: " << image->getGlobalCount() << std::endl;
        return 0;
    } else {
        std::cerr << "Failed to save image" << std::endl;
        return 1;
    }
}

int loadSourceFiles(const std::vector<std::string>& sourceFiles, const std::string& imageFile) {
    std::cout << "Loading " << sourceFiles.size() << " source files..." << std::endl;
    for (const auto& file : sourceFiles) {
        std::cout << "  " << file << std::endl;
    }
    std::cout << "Creating image: " << imageFile << std::endl;
    
    auto& manager = ImageManager::getInstance();
    manager.createFreshImage();
    
    if (!manager.loadSourceFiles(sourceFiles)) {
        std::cerr << "Failed to load source files" << std::endl;
        return 1;
    }
    
    if (manager.saveImageToFile(imageFile)) {
        auto* image = manager.getCurrentImage();
        std::cout << "Image created successfully!" << std::endl;
        std::cout << "  Source files: " << image->getSourceFiles().size() << std::endl;
        std::cout << "  Classes: " << image->getClassCount() << std::endl;
        std::cout << "  Methods: " << image->getMethodCount() << std::endl;
        std::cout << "  Globals: " << image->getGlobalCount() << std::endl;
        return 0;
    } else {
        std::cerr << "Failed to save image" << std::endl;
        return 1;
    }
}

int showImageInfo(const std::string& imageFile) {
    std::cout << "Image information: " << imageFile << std::endl;
    
    if (!ImageUtils::isValidImageFile(imageFile)) {
        std::cerr << "Error: Not a valid Smalltalk image file" << std::endl;
        return 1;
    }
    
    std::string version;
    uint64_t creationTime;
    uint32_t classCount;
    
    if (ImageUtils::getImageInfo(imageFile, version, creationTime, classCount)) {
        std::cout << "  Version: " << version << std::endl;
        std::cout << "  Created: " << ImageUtils::formatTimestamp(creationTime) << std::endl;
        std::cout << "  Classes: " << classCount << std::endl;
        
        // Load the image to get more detailed info
        auto& manager = ImageManager::getInstance();
        if (manager.loadImageFromFile(imageFile)) {
            auto* image = manager.getCurrentImage();
            std::cout << "  Methods: " << image->getMethodCount() << std::endl;
            std::cout << "  Globals: " << image->getGlobalCount() << std::endl;
            std::cout << "  Modified: " << ImageUtils::formatTimestamp(image->getModificationTime()) << std::endl;
            
            // Show loaded source files
            auto sourceFiles = image->getSourceFiles();
            if (!sourceFiles.empty()) {
                std::cout << "  Source files (" << sourceFiles.size() << "):" << std::endl;
                for (const auto& file : sourceFiles) {
                    std::cout << "    " << file.filename << std::endl;
                }
            }
            
            // Show globals
            auto globalNames = image->getGlobalNames();
            if (!globalNames.empty()) {
                std::cout << "  Globals:" << std::endl;
                for (const auto& name : globalNames) {
                    std::cout << "    " << name << std::endl;
                }
            }
            
            // Show metadata
            auto metadata = image->getMetadata();
            if (!metadata.empty()) {
                std::cout << "  Metadata:" << std::endl;
                for (const auto& pair : metadata) {
                    std::cout << "    " << pair.first << ": " << pair.second << std::endl;
                }
            }
        }
        
        return 0;
    } else {
        std::cerr << "Failed to read image information" << std::endl;
        return 1;
    }
}

int runExpression(const std::string& imageFile, const std::string& expression) {
    std::cout << "Loading image: " << imageFile << std::endl;
    
    auto& manager = ImageManager::getInstance();
    if (!manager.loadImageFromFile(imageFile)) {
        std::cerr << "Failed to load image" << std::endl;
        return 1;
    }
    
    std::cout << "Evaluating: " << expression << std::endl;
    
    auto* image = manager.getCurrentImage();
    TaggedValue result = image->evaluate(expression);
    
    printResult(result);
    return 0;
}

int evaluateExpression(const std::string& expression) {
    std::cout << "Evaluating in fresh image: " << expression << std::endl;
    
    auto& manager = ImageManager::getInstance();
    manager.createFreshImage();
    
    auto* image = manager.getCurrentImage();
    TaggedValue result = image->evaluate(expression);
    
    printResult(result);
    return 0;
}

int main(int argc, char* argv[]) {
    if (argc < 2) {
        printUsage(argv[0]);
        return 1;
    }
    
    std::string command = argv[1];
    
    try {
        if (command == "create") {
            if (argc != 3) {
                std::cerr << "Usage: " << argv[0] << " create <image_file>" << std::endl;
                return 1;
            }
            return createImage(argv[2]);
            
        } else if (command == "load") {
            if (argc != 4) {
                std::cerr << "Usage: " << argv[0] << " load <source_dir> <image_file>" << std::endl;
                return 1;
            }
            return loadSourceDirectory(argv[2], argv[3]);
            
        } else if (command == "loadfiles") {
            if (argc < 4) {
                std::cerr << "Usage: " << argv[0] << " loadfiles <file1> <file2> ... <image_file>" << std::endl;
                return 1;
            }
            
            std::vector<std::string> sourceFiles;
            for (int i = 2; i < argc - 1; ++i) {
                sourceFiles.push_back(argv[i]);
            }
            std::string imageFile = argv[argc - 1];
            
            return loadSourceFiles(sourceFiles, imageFile);
            
        } else if (command == "info") {
            if (argc != 3) {
                std::cerr << "Usage: " << argv[0] << " info <image_file>" << std::endl;
                return 1;
            }
            return showImageInfo(argv[2]);
            
        } else if (command == "run") {
            if (argc != 4) {
                std::cerr << "Usage: " << argv[0] << " run <image_file> <expression>" << std::endl;
                return 1;
            }
            return runExpression(argv[2], argv[3]);
            
        } else if (command == "eval") {
            if (argc != 3) {
                std::cerr << "Usage: " << argv[0] << " eval <expression>" << std::endl;
                return 1;
            }
            return evaluateExpression(argv[2]);
            
        } else {
            std::cerr << "Unknown command: " << command << std::endl;
            printUsage(argv[0]);
            return 1;
        }
        
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        return 1;
    }
}