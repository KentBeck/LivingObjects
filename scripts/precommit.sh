#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel)"
cd "$ROOT_DIR"

echo "[pre-commit] Running pre-commit checks..."

# 1) Optional formatting step (clang-format) for C/C++ files
if command -v clang-format >/dev/null 2>&1; then
  echo "[pre-commit] Formatting C/C++ sources with clang-format"
  find src/cpp -type f \( -name "*.h" -o -name "*.hpp" -o -name "*.hh" -o -name "*.c" -o -name "*.cc" -o -name "*.cpp" \) -print0 | xargs -0 clang-format -i || true
  git add -A || true
else
  echo "[pre-commit] clang-format not found; skipping formatting"
fi

# 2) Build with warnings captured; fail on any warnings
echo "[pre-commit] Building project (capturing warnings)"
BUILD_LOG="$(mktemp)"
(
  cd src/cpp
  make clean >/dev/null 2>&1 || true
  make all
) 2> >(tee -a "$BUILD_LOG" >&2) 1> >(tee -a "$BUILD_LOG")

if grep -Ei "warning:" "$BUILD_LOG" >/dev/null 2>&1; then
  echo "[pre-commit] Compiler warnings detected; please fix them before committing."
  echo "--------- warnings (truncated) ---------"
  grep -E "warning:" -n "$BUILD_LOG" | head -50 || true
  echo "---------------------------------------"
  exit 1
fi

# 3) Run expression tests (fast gate)
echo "[pre-commit] Running expression test suite"
./src/cpp/tests/run_expression_tests.sh >/dev/null

# 4) Run full test suite (builds test binaries and runs them)
if [[ "${SKIP_FULL_SUITE:-}" != "1" ]]; then
  echo "[pre-commit] Running full test suite"
  SUITE_LOG="$(mktemp)"
  set +e
  ./src/cpp/run_all_tests.sh 2> >(tee -a "$SUITE_LOG" >&2) 1> >(tee -a "$SUITE_LOG")
  STATUS=$?
  set -e
  if [[ $STATUS -ne 0 ]]; then
      echo "[pre-commit] Full suite failed; see logs above."
      exit 1
  fi
  if grep -Ei "warning:" "$SUITE_LOG" >/dev/null 2>&1; then
      echo "[pre-commit] Warnings detected during full suite build; please address before committing."
      grep -E "warning:" -n "$SUITE_LOG" | head -50 || true
      exit 1
  fi
fi

# 5) Detect disallowed temporary files at repo root (common mistakes)
DISALLOWED=("compile:in:" "start:")
FOUND=()
for f in "${DISALLOWED[@]}"; do
  if [[ -e "$f" ]]; then
    FOUND+=("$f")
  fi
done
if [[ ${#FOUND[@]} -gt 0 ]]; then
  echo "[pre-commit] Disallowed temporary files present: ${FOUND[*]}"
  echo "[pre-commit] Please remove them before committing."
  exit 1
fi

echo "[pre-commit] All checks passed. Proceeding with commit."
