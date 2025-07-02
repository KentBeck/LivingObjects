#!/bin/bash

# Smalltalk Image System Demo
# This script demonstrates the complete workflow of:
# 1. Creating source files
# 2. Loading them into an image
# 3. Saving the binary image
# 4. Loading and using the image

echo "==============================================="
echo "Smalltalk Image System Demo"
echo "==============================================="
echo

# Build the image tool if needed
if [ ! -f "build/image-tool" ]; then
    echo "Building image tool..."
    make clean && make
    echo
fi

# Clean up any existing demo files
echo "Cleaning up existing demo files..."
rm -f demo*.image demo*.st
echo

# Create demo source files
echo "Creating demo source files..."

cat > demo_arithmetic.st << 'EOF'
"Demo arithmetic expressions"

"Basic arithmetic that works"
42
123
0

"Boolean literals"
true
false
nil

"String literals"
'Hello from Smalltalk!'
'This is a demo'
'Source code loaded from file'
EOF

cat > demo_values.st << 'EOF'
"Demo value expressions"

"More integer literals"
999
-5
7

"More strings"
'Smalltalk'
'image'
'system'

"More booleans"
true
false
nil
EOF

echo "Created demo source files:"
echo "  - demo_arithmetic.st"
echo "  - demo_values.st"
echo

# Show source file contents
echo "Source file contents:"
echo "=================="
echo "demo_arithmetic.st:"
cat demo_arithmetic.st
echo
echo "demo_values.st:"
cat demo_values.st
echo

# Test 1: Create fresh image
echo "==============================================="
echo "Test 1: Creating a fresh image"
echo "==============================================="
./build/image-tool create demo_fresh.image
echo

# Test 2: Show fresh image info
echo "==============================================="
echo "Test 2: Fresh image information"
echo "==============================================="
./build/image-tool info demo_fresh.image
echo

# Test 3: Load source files into image
echo "==============================================="
echo "Test 3: Loading source files into image"
echo "==============================================="
./build/image-tool loadfiles demo_arithmetic.st demo_values.st demo_loaded.image
echo

# Test 4: Show loaded image info
echo "==============================================="
echo "Test 4: Loaded image information"
echo "==============================================="
./build/image-tool info demo_loaded.image
echo

# Test 5: Evaluate expressions in fresh image
echo "==============================================="
echo "Test 5: Evaluating expressions in fresh image"
echo "==============================================="
echo "Evaluating: 42"
./build/image-tool eval "42"
echo

echo "Evaluating: 'hello world'"
./build/image-tool eval "'hello world'"
echo

echo "Evaluating: true"
./build/image-tool eval "true"
echo

echo "Evaluating: nil"
./build/image-tool eval "nil"
echo

# Test 6: Evaluate expressions in loaded image
echo "==============================================="
echo "Test 6: Evaluating expressions in loaded image"
echo "==============================================="
echo "Evaluating: 999"
./build/image-tool run demo_loaded.image "999"
echo

echo "Evaluating: 'Smalltalk'"
./build/image-tool run demo_loaded.image "'Smalltalk'"
echo

echo "Evaluating: false"
./build/image-tool run demo_loaded.image "false"
echo

# Test 7: Binary image properties
echo "==============================================="
echo "Test 7: Binary image file properties"
echo "==============================================="
echo "Image file sizes:"
ls -la *.image
echo

echo "Image file types:"
file *.image
echo

# Test 8: Image validation
echo "==============================================="
echo "Test 8: Image validation"
echo "==============================================="
echo "Testing if demo_loaded.image is a valid image:"
if ./build/image-tool info demo_loaded.image > /dev/null 2>&1; then
    echo "✓ demo_loaded.image is valid"
else
    echo "✗ demo_loaded.image is invalid"
fi
echo

echo "Testing if demo_arithmetic.st is a valid image:"
if ./build/image-tool info demo_arithmetic.st > /dev/null 2>&1; then
    echo "✓ demo_arithmetic.st is valid (unexpected!)"
else
    echo "✓ demo_arithmetic.st is not a valid image (expected)"
fi
echo

# Summary
echo "==============================================="
echo "Demo Summary"
echo "==============================================="
echo "✓ Created fresh Smalltalk image"
echo "✓ Loaded source files into image"
echo "✓ Saved binary image to disk"
echo "✓ Loaded binary image from disk"
echo "✓ Evaluated expressions in both fresh and loaded images"
echo "✓ Validated image file format"
echo
echo "Key capabilities demonstrated:"
echo "  - Source file loading (.st files)"
echo "  - Binary image persistence"
echo "  - Expression evaluation"
echo "  - Image introspection"
echo "  - Cross-session state preservation"
echo
echo "Next steps for full Smalltalk implementation:"
echo "  - Method compilation and storage"
echo "  - Class definition parsing"
echo "  - Object instance serialization"
echo "  - Garbage collection integration"
echo "  - Method dispatch restoration"
echo
echo "Demo complete!"
