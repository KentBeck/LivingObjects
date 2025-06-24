#!/bin/bash

# Script to compile and run the expression tests

set -e

# Compile the expression test runner
g++ -std=c++17 -Wall -Wextra -I../include tests/expression_test_runner.cpp src/memory_manager.cpp -o tests/expression_test_runner

# Run the tests
./tests/expression_test_runner tests/expression_tests.txt

echo "Expression tests completed!"