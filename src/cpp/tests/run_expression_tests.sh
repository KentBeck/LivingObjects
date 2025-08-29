#!/bin/bash

# Script to compile and run the expression tests

set -e

# Compile the expression test runner
g++ -std=c++17 -Wall -Wextra -Isrc/cpp/include \
  src/cpp/tests/all_expressions_test.cpp \
  src/cpp/src/interpreter.cpp \
  src/cpp/src/memory_manager.cpp \
  src/cpp/src/object.cpp \
  src/cpp/src/primitives.cpp \
  src/cpp/src/simple_compiler.cpp \
  src/cpp/src/simple_object.cpp \
  src/cpp/src/simple_parser.cpp \
  src/cpp/src/smalltalk_class.cpp \
  src/cpp/src/smalltalk_image.cpp \
  src/cpp/src/smalltalk_image_interpreter.cpp \
  src/cpp/src/smalltalk_string.cpp \
  src/cpp/src/smalltalk_vm.cpp \
  src/cpp/src/method_compiler.cpp \
  src/cpp/src/smalltalk_exception.cpp \
  src/cpp/src/symbol.cpp \
  src/cpp/src/tagged_value.cpp \
  src/cpp/src/primitives/array.cpp \
  src/cpp/src/primitives/block.cpp \
  src/cpp/src/primitives/integer.cpp \
  src/cpp/src/primitives/object.cpp \
  src/cpp/src/primitives/string.cpp \
  src/cpp/src/primitives/exception.cpp \
  -o src/cpp/tests/expression_test_runner

# Run the tests
src/cpp/tests/expression_test_runner src/cpp/tests/expression_tests.txt

echo "Expression tests completed!"
