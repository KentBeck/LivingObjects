Engineering workflow policy

Before committing any changes to this repository:

- Run the full build and test suites, and ensure they pass.
- Address all compiler warnings (treat warnings as errors for gatekeeping).
- Format the codebase (C/C++ and scripts) so no formatting diffs remain.
- Remove or avoid committing temporary, generated, or stray files.

Notes

- The pre-commit hook installed in this repo automates these checks. It will:
  - Format C/C++ sources via clang-format if available.
  - Build the project and fail on any compiler warnings.
  - Run the expression tests and fail if they do not pass.
  - Block commits if disallowed temporary files are detected.
- You can bypass the full suite temporarily by setting SKIP_FULL_SUITE=1, but this should be used sparingly.

