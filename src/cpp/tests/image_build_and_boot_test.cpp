#include "smalltalk_image.h"
#include "interpreter.h"
#include "memory_manager.h"
#include "smalltalk_class.h"
#include "symbol.h"

#include <cassert>
#include <filesystem>
#include <iostream>
#include <string>

// Simple test helper macro
#define EXPECT_TRUE(cond)                                                                 \
    do {                                                                                  \
        if (!(cond)) {                                                                     \
            std::cerr << "EXPECT_TRUE failed: " #cond " at " << __FILE__ << ":" << __LINE__ \
                      << std::endl;                                                       \
            std::abort();                                                                  \
        }                                                                                 \
    } while (0)

#define EXPECT_EQ(a, b)                                                                    \
    do {                                                                                  \
        auto _va = (a);                                                                    \
        auto _vb = (b);                                                                    \
        if (!(_va == _vb)) {                                                               \
            std::cerr << "EXPECT_EQ failed: " #a " == " #b " got (" << _va << ", " << _vb  \
                      << ") at " << __FILE__ << ":" << __LINE__ << std::endl;            \
            std::abort();                                                                  \
        }                                                                                 \
    } while (0)

int main() {
    using namespace smalltalk;

    // 1) Build an image from kernel .st sources
    auto& manager = ImageManager::getInstance();
    manager.createFreshImage();

    // Load .st files from the kernel directory
    EXPECT_TRUE(manager.loadSourceDirectory("st/kernel"));

    // Save image to a temporary location
    const auto tmp = std::filesystem::temp_directory_path();
    const std::string imagePath = (tmp / "smalltalk_core_test.image").string();
    EXPECT_TRUE(manager.saveImageToFile(imagePath));

    // 2) Load the image into a fresh ImageManager
    ImageManager freshManager;
    EXPECT_TRUE(freshManager.loadImageFromFile(imagePath));
    SmalltalkImage* image = freshManager.getCurrentImage();
    EXPECT_TRUE(image != nullptr);

    // Sanity: core classes exist after load
    EXPECT_TRUE(ClassRegistry::getInstance().hasClass("Object"));
    EXPECT_TRUE(ClassRegistry::getInstance().hasClass("Integer"));

    // 3) Exercise SystemLoader>>start: through the minimal primitive
    TaggedValue started = image->evaluate("SystemLoader new start: 'kernel'");
    EXPECT_TRUE(started.isBoolean());
    EXPECT_TRUE(started.asBoolean());

    // 4) Evaluate a simple expression using the loaded image
    //    This exercises interpreter + primitives after snapshot load
    TaggedValue result = image->evaluate("1 + 2");
    EXPECT_TRUE(result.isInteger());
    EXPECT_EQ(result.asInteger(), 3);

    std::cout << "image_build_and_boot_test: PASS" << std::endl;
    return 0;
}
