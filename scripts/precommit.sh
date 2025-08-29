#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel)"
CPP_DIR="$ROOT_DIR/src/cpp"
cd "$ROOT_DIR"

echo "[pre-commit] Running pre-commit checks..."

# 1) Formatting check (clang-format): do not modify index during commit
if command -v clang-format >/dev/null 2>&1; then
  echo "[pre-commit] Checking C/C++ formatting (clang-format)"
  # Check only tracked files under src/cpp to avoid large argv/xargs limits
  NEEDS_FORMAT=0
  while IFS= read -r f; do
    if [[ -f "$f" ]]; then
      if ! clang-format -n --Werror "$f" 2>/dev/null; then
        echo "  needs format: $f"
        NEEDS_FORMAT=1
      fi
    fi
  done < <(git ls-files 'src/cpp/**/*.h' 'src/cpp/**/*.hpp' 'src/cpp/**/*.hh' 'src/cpp/**/*.c' 'src/cpp/**/*.cc' 'src/cpp/**/*.cpp')
  if [[ $NEEDS_FORMAT -ne 0 ]]; then
    echo "[pre-commit] Formatting issues found. Please run clang-format on the listed files and re-stage."
    exit 1
  fi
else
  echo "[pre-commit] clang-format not found; skipping formatting check"
fi

# 2) Build with warnings captured; fail on any warnings
echo "[pre-commit] Building project (capturing warnings)"
BUILD_LOG="$(mktemp)"
set +e
(
  cd "$CPP_DIR" && make clean >/dev/null 2>&1 || true
  cd "$CPP_DIR" && make all
) >"$BUILD_LOG" 2>&1
BUILD_STATUS=$?
set -e
if [[ $BUILD_STATUS -ne 0 ]]; then
  echo "[pre-commit] Build failed; tail of log:"
  tail -200 "$BUILD_LOG" || true
  exit 1
fi

if grep -Ei "warning:" "$BUILD_LOG" >/dev/null 2>&1; then
  echo "[pre-commit] Compiler warnings detected; please fix them before committing."
  echo "--------- warnings (truncated) ---------"
  grep -E "warning:" -n "$BUILD_LOG" | head -50 || true
  echo "---------------------------------------"
  exit 1
fi

# 3) Run expression tests (fast gate)
echo "[pre-commit] Running expression test suite"
EXP_LOG="$(mktemp)"
if ! "$ROOT_DIR"/src/cpp/tests/run_expression_tests.sh >"$EXP_LOG" 2>&1; then
  echo "[pre-commit] Expression tests failed. See output below:"
  tail -200 "$EXP_LOG" || true
  exit 1
fi

# 4) Run full test suite (builds test binaries and runs them)
echo "[pre-commit] Running full test suite"
SUITE_LOG="$(mktemp)"
if ! "$ROOT_DIR"/src/cpp/run_all_tests.sh >"$SUITE_LOG" 2>&1; then
    echo "[pre-commit] Full suite failed; tail of log:" 
    tail -200 "$SUITE_LOG" || true
    exit 1
fi
if grep -Ei "warning:" "$SUITE_LOG" >/dev/null 2>&1; then
    echo "[pre-commit] Warnings detected during full suite build; please address before committing."
    grep -E "warning:" -n "$SUITE_LOG" | head -50 || true
    exit 1
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
