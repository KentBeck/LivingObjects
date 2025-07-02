#!/bin/bash

# Simple expression runner for Smalltalk VM
# Allows quick testing of expressions without full test suite

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}=================================================="
echo -e "         SMALLTALK VM EXPRESSION RUNNER"
echo -e "==================================================${NC}"

# Build the VM if needed
if [[ ! -f "build/smalltalk-vm" ]] || [[ "src/main.cpp" -nt "build/smalltalk-vm" ]]; then
    echo -e "${YELLOW}Building Smalltalk VM...${NC}"
    make all > /dev/null 2>&1
    echo -e "${GREEN}✓ Build complete${NC}"
fi

# Function to run expression
run_expr() {
    local expr="$1"
    echo -e "\n${BLUE}Expression:${NC} ${YELLOW}$expr${NC}"

    if output=$(./build/smalltalk-vm "$expr" 2>&1); then
        result=$(echo "$output" | grep 'Result:' | sed 's/Result: //')
        echo -e "${GREEN}✓ Result:${NC} $result"

        if [[ "${SHOW_DETAILS:-}" == "1" ]]; then
            echo -e "${CYAN}Details:${NC}"
            echo "$output"
        fi
    else
        echo -e "${RED}✗ Error:${NC}"
        echo "$output"
    fi
}

# Function to demo blocks (shows parsing/compilation)
demo_block() {
    local expr="$1"
    echo -e "\n${BLUE}Block Demo:${NC} ${YELLOW}$expr${NC}"
    echo -e "${CYAN}(Shows parsing and compilation)${NC}"

    if output=$(./build/smalltalk-vm "$expr" 2>&1); then
        echo "$output" | head -6
    else
        echo "$output" | head -6
    fi
}

# If arguments provided, run them
if [[ $# -gt 0 ]]; then
    for expr in "$@"; do
        run_expr "$expr"
    done
    exit 0
fi

# Otherwise run demo expressions
echo -e "\n${CYAN}=== ARITHMETIC EXPRESSIONS ===${NC}"
run_expr "3 + 4"
run_expr "10 - 3"
run_expr "6 * 7"
run_expr "20 / 4"
run_expr "(2 + 3) * (4 + 1)"

echo -e "\n${CYAN}=== COMPARISON EXPRESSIONS ===${NC}"
run_expr "5 > 3"
run_expr "2 < 8"
run_expr "4 = 4"
run_expr "7 ~= 9"
run_expr "5 <= 5"
run_expr "8 >= 6"

echo -e "\n${CYAN}=== BOOLEAN AND NIL ===${NC}"
run_expr "true"
run_expr "false"
run_expr "nil"

echo -e "\n${CYAN}=== COMPLEX EXPRESSIONS ===${NC}"
run_expr "((10 + 5) * 2) / 6"
run_expr "(3 + 4) > (2 * 3)"
run_expr "1 + 2 * 3 + 4"

echo -e "\n${CYAN}=== BLOCK EXPRESSIONS ===${NC}"
echo -e "${YELLOW}Note: Blocks parse and compile successfully but need full VM context to execute${NC}"
demo_block "[42]"
demo_block "[3 + 4]"
demo_block "[true]"
demo_block "[3 + 4. 5 * 6]"
demo_block "[1 + 2. 3 * 4. 5 - 1]"

echo -e "\n${CYAN}=== USAGE ===${NC}"
echo -e "Run specific expressions:"
echo -e "  ${YELLOW}./run_expressions.sh \"3 + 4\" \"5 > 2\"${NC}"
echo -e ""
echo -e "Show compilation details:"
echo -e "  ${YELLOW}SHOW_DETAILS=1 ./run_expressions.sh \"3 + 4\"${NC}"
echo -e ""
echo -e "Interactive mode:"
echo -e "  ${YELLOW}./build/smalltalk-vm \"your expression here\"${NC}"

echo -e "\n${GREEN}Expression runner complete!${NC}"
