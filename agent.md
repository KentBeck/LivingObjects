Engineering workflow policy

Before committing any changes to this repository:

- Run the full build and test suites, and ensure they pass.
- Address all compiler warnings (treat warnings as errors for gatekeeping).
- Format the codebase (C/C++ and scripts) so no formatting diffs remain.
- Remove or avoid committing temporary, generated, or stray files.

Notes

- The pre-commit hook installed in this repo automates these checks. It will:
  - Check C/C++ formatting via clang-format (no auto-fix during commit).
  - Build the project and fail on any compiler warnings.
  - Run the expression tests and fail if they do not pass.
  - Run the full test suite and fail if any test fails or emits warnings.
  - Block commits if disallowed temporary files are detected.

Design principles

- Represent everything possible as a Smalltalk object. This includes runtime constructs like MethodContext and BlockContext, not just user-visible classes and data structures.
- Persist and restore execution state. The image should be able to serialize active contexts (stack frames), sender links, instruction pointers, and stack slots, and resume execution after loading.
- Route global references through the Smalltalk dictionary. C++ helpers should use the same global mechanism as Smalltalk code for consistency.
