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

# Detect host OS/arch to avoid mismatched architectures on macOS
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# On Apple Silicon, force arm64 for both compile and link
ifeq ($(UNAME_S),Darwin)
  ifeq ($(UNAME_M),arm64)
    CXXFLAGS += -arch arm64
    LDFLAGS += -arch arm64
  endif
endif

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

# Source files (excluding main files and image_tool)
SOURCES = $(filter-out $(SRC_DIR)/main.cpp $(SRC_DIR)/smalltalk_repl.cpp $(SRC_DIR)/image_tool.cpp, $(shell find $(SRC_DIR) -name "*.cpp"))
OBJECTS = $(SOURCES:$(SRC_DIR)/%.cpp=$(BUILD_DIR)/%.o)

# Test files
TEST_SOURCES = $(wildcard $(TEST_DIR)/*.cpp)
TEST_OBJECTS = $(TEST_SOURCES:$(TEST_DIR)/%.cpp=$(BUILD_DIR)/tests/%.o)

# Main executable
EXECUTABLE = $(BUILD_DIR)/smalltalk-vm

# REPL executable
REPL_EXECUTABLE = $(BUILD_DIR)/smalltalk-repl

# Image tool executable
IMAGE_TOOL = $(BUILD_DIR)/image-tool

# Test executable
TEST_EXECUTABLE = $(BUILD_DIR)/run-tests

# Library
LIBRARY = $(BUILD_DIR)/lib/libsmalltalk_vm.a

# Default target
all: create_dirs $(LIBRARY) $(EXECUTABLE) $(REPL_EXECUTABLE) $(IMAGE_TOOL)

# Create build directories
create_dirs:
	@mkdir -p $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)/tests
	@mkdir -p $(BUILD_DIR)/lib

# Build the main library
$(LIBRARY): $(OBJECTS)
	@mkdir -p $(dir $@)
	ar rcs $@ $^

# Build the main executable
$(EXECUTABLE): $(OBJECTS) $(BUILD_DIR)/main.o
	$(CXX) $(CXXFLAGS) $(LDFLAGS) $^ -o $@

# Build the REPL executable
$(REPL_EXECUTABLE): $(BUILD_DIR)/smalltalk_repl.o $(OBJECTS)
	$(CXX) $(CXXFLAGS) $(LDFLAGS) $^ -o $@

# Build the image tool executable
$(IMAGE_TOOL): $(BUILD_DIR)/image_tool.o $(OBJECTS)
	$(CXX) $(CXXFLAGS) $(LDFLAGS) $^ -o $@

# Compile source files
$(BUILD_DIR)/%.o: $(SRC_DIR)/%.cpp
	@mkdir -p $(dir $@)
	$(CXX) $(CXXFLAGS) -c $< -o $@

# Compile main.cpp
$(BUILD_DIR)/main.o: $(SRC_DIR)/main.cpp
	$(CXX) $(CXXFLAGS) -c $< -o $@

# Compile image_tool.cpp
$(BUILD_DIR)/image_tool.o: $(SRC_DIR)/image_tool.cpp
	$(CXX) $(CXXFLAGS) -c $< -o $@

# Build and run tests
test: create_dirs $(TEST_EXECUTABLE)
	./$(TEST_EXECUTABLE)

# Build the test executable
$(TEST_EXECUTABLE): $(TEST_DIR)/all_expressions_test.cpp $(LIBRARY)
	$(CXX) $(CXXFLAGS) $(LDFLAGS) $< -o $@ -L$(BUILD_DIR)/lib -lsmalltalk_vm

# Build individual test executables
$(BUILD_DIR)/minimal_parse_test: $(TEST_DIR)/minimal_parse_test.cpp $(LIBRARY)
	$(CXX) $(CXXFLAGS) $(LDFLAGS) $< -o $@ -L$(BUILD_DIR)/lib -lsmalltalk_vm

$(BUILD_DIR)/simple_block_test: $(TEST_DIR)/simple_block_test.cpp $(LIBRARY)
	$(CXX) $(CXXFLAGS) $(LDFLAGS) $< -o $@ -L$(BUILD_DIR)/lib -lsmalltalk_vm

$(BUILD_DIR)/minimal_block_test: $(TEST_DIR)/minimal_block_test.cpp $(LIBRARY)
	$(CXX) $(CXXFLAGS) $(LDFLAGS) $< -o $@ -L$(BUILD_DIR)/lib -lsmalltalk_vm

$(BUILD_DIR)/debug_block_layout: $(TEST_DIR)/debug_block_layout.cpp $(LIBRARY)
	$(CXX) $(CXXFLAGS) $(LDFLAGS) $< -o $@ -L$(BUILD_DIR)/lib -lsmalltalk_vm

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

# Run the REPL
repl: $(REPL_EXECUTABLE)
	./$(REPL_EXECUTABLE)

# Run all tests
test-all: all
	./run_all_tests.sh

# Run tests with verbose output
test-verbose: all
	VERBOSE=1 ./run_all_tests.sh

# Run expression demos
expressions: all
	./run_expressions.sh

# Run expression with details
expressions-verbose: all
	SHOW_DETAILS=1 ./run_expressions.sh

# Build individual test targets
test-parse: $(BUILD_DIR)/minimal_parse_test
	./$(BUILD_DIR)/minimal_parse_test

test-block-compile: $(BUILD_DIR)/simple_block_test
	./$(BUILD_DIR)/simple_block_test

test-block-runtime: $(BUILD_DIR)/minimal_block_test
	./$(BUILD_DIR)/minimal_block_test

test-block-layout: $(BUILD_DIR)/debug_block_layout
	./$(BUILD_DIR)/debug_block_layout

# Quick test - just run our main expression test
test-quick: $(TEST_EXECUTABLE)
	./$(TEST_EXECUTABLE)

# Test Object new specifically
test-object-new: all
	echo "Object new" | ./$(REPL_EXECUTABLE)

.PHONY: all create_dirs test clean rebuild install_gtest install_gtest_mac repl test-all test-verbose expressions expressions-verbose test-parse test-block-compile test-block-runtime test-block-layout test-quick test-object-new
