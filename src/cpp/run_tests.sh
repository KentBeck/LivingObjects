#!/bin/bash

# Simple script to compile and run tests

set -e

# Compile the tests
g++ -std=c++17 -Wall -Wextra -I./include tests/simple_test_main.cpp src/memory_manager.cpp -o tests/run_tests

# Run the tests
./tests/run_tests

echo "Tests completed successfully!"