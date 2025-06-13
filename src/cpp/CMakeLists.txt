cmake_minimum_required(VERSION 3.10)
project(SmalltalkVM VERSION 0.1.0 LANGUAGES CXX)

# Set C++ standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# Set default build type to Debug
if(NOT CMAKE_BUILD_TYPE)
  set(CMAKE_BUILD_TYPE Debug)
endif()

# Output directories
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)

# Include directories
include_directories(
  ${PROJECT_SOURCE_DIR}/include
)

# Add source files
file(GLOB_RECURSE SOURCES "src/*.cpp")

# Add the main library
add_library(smalltalk_vm STATIC ${SOURCES})

# Add executable
add_executable(smalltalk_vm_exe src/main.cpp)
target_link_libraries(smalltalk_vm_exe smalltalk_vm)
set_target_properties(smalltalk_vm_exe PROPERTIES OUTPUT_NAME "smalltalk-vm")

# Tests
enable_testing()
add_subdirectory(tests)