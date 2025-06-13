# Makefile for the C++ Smalltalk VM

# Directories
BUILD_DIR = build
SRC_DIR = src
INCLUDE_DIR = include
TEST_DIR = tests

# Compiler and flags
CXX = g++
CXXFLAGS = -std=c++17 -Wall -Wextra -I$(INCLUDE_DIR)
LDFLAGS = 

# Debug build flags
DEBUG_FLAGS = -g -O0 -DDEBUG

# Release build flags
RELEASE_FLAGS = -O3 -DNDEBUG

# Default build type
BUILD_TYPE ?= Debug

ifeq ($(BUILD_TYPE), Debug)
    CXXFLAGS += $(DEBUG_FLAGS)
else
    CXXFLAGS += $(RELEASE_FLAGS)
endif

# Source files
SOURCES = $(wildcard $(SRC_DIR)/*.cpp)
OBJECTS = $(SOURCES:$(SRC_DIR)/%.cpp=$(BUILD_DIR)/%.o)

# Test files
TEST_SOURCES = $(wildcard $(TEST_DIR)/*.cpp)
TEST_OBJECTS = $(TEST_SOURCES:$(TEST_DIR)/%.cpp=$(BUILD_DIR)/tests/%.o)

# Main executable
EXECUTABLE = $(BUILD_DIR)/smalltalk-vm

# Test executable
TEST_EXECUTABLE = $(BUILD_DIR)/run-tests

# Default target
all: create_dirs $(EXECUTABLE)

# Create build directories
create_dirs:
	@mkdir -p $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)/tests

# Build the main executable
$(EXECUTABLE): $(OBJECTS) $(BUILD_DIR)/main.o
	$(CXX) $(CXXFLAGS) $(LDFLAGS) $^ -o $@

# Compile source files
$(BUILD_DIR)/%.o: $(SRC_DIR)/%.cpp
	$(CXX) $(CXXFLAGS) -c $< -o $@

# Compile main.cpp
$(BUILD_DIR)/main.o: $(SRC_DIR)/main.cpp
	$(CXX) $(CXXFLAGS) -c $< -o $@

# Build and run tests
test: create_dirs $(TEST_EXECUTABLE)
	./$(TEST_EXECUTABLE)

# Build the test executable
$(TEST_EXECUTABLE): $(TEST_OBJECTS) $(filter-out $(BUILD_DIR)/main.o, $(OBJECTS))
	$(CXX) $(CXXFLAGS) $(LDFLAGS) $^ -o $@ -lgtest -lgtest_main -pthread

# Compile test files
$(BUILD_DIR)/tests/%.o: $(TEST_DIR)/%.cpp
	@mkdir -p $(BUILD_DIR)/tests
	$(CXX) $(CXXFLAGS) -I$(INCLUDE_DIR) -c $< -o $@

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Clean and rebuild
rebuild: clean all

# Install GTest for tests (Linux)
install_gtest:
	apt-get update && apt-get install -y libgtest-dev
	cd /usr/src/gtest && cmake . && make && cp lib/*.a /usr/lib

# Install GTest for tests (macOS)
install_gtest_mac:
	brew install googletest

.PHONY: all create_dirs test clean rebuild install_gtest install_gtest_mac