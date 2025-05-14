# Plan for Updating and Removing Deprecated Code

This document outlines a plan for completely migrating from the old core/classes structure to the new consolidated pile package.

## Files Updated So Far
- All VM test files have been updated to use pile package:
  - vm/factorial_test.go
  - vm/vm_test.go
  - vm/block_test.go (converted to vm_test package)
  - vm/block_bytecode_handlers_test.go (converted to vm_test package)
  - vm/bytecode_handlers_test.go
  - vm/block_scenarios_test.go (converted to vm_test package)
  - vm/block_variable_access_test.go (converted to vm_test package)
  - vm/bytecode_test.go
  - vm/array_primitives_test.go
  - vm/basicclass_test.go
  - vm/boolean_primitives_test.go
  - vm/byte_array_primitives_test.go
  - vm/byte_array_test.go
  - vm/consolidated_benchmark_test.go
  - vm/context_test.go
  - vm/method_block_literals_test.go
  - vm/nil_class_test.go
  - vm/primitives_test.go
  - vm/send_message_stack_test.go
  - vm/send_message_test.go
  - vm/string_primitives_test.go
  - vm/string_test.go
  - vm/vm_class_test.go
  - vm/vm_method_lookup_test.go
- Demo files:
  - factorial_demo.go (moved to demo package and updated to use pile)
  - demo/factorial.go (updated to use pile)
- VM implementation files:
  - vm/vm.go (fully updated to use pile)
  - vm/block_factory.go (updated to use pile)
  - vm/byte_array.go (updated to use pile)
  - vm/block_executor.go (updated to use pile)
  - vm/bytecode_handlers.go (updated to use pile)
  - vm/context.go (updated to use pile)
  - vm/executor.go (updated to use pile)
  - vm/block_bytecode_handlers.go (updated to use pile)
  - vm/class_accessors.go (updated to use pile)
  - vm/class_factory.go (updated to use pile)
  - vm/class_registry.go (updated to use pile)
  - vm/dictionary_factory.go (updated to use pile)
  - vm/method_factory.go (updated to use pile)
  - vm/symbol_factory.go (updated to use pile)
- Compiler files:
  - compiler/bytecode_compiler.go (updated to use pile)
  - compiler/method_builder.go (updated to use pile)
  - compiler/method_factory.go (updated to use pile)
  - compiler/vm_access.go (updated to use pile)
  - compiler/bytecode_compiler_test.go (updated to use pile)
  - compiler/method_builder_test.go (updated to use pile)
- Parser files:
  - parser/parser.go (updated to use pile)
  - parser/combined_test.go (updated to use pile)
  - parser/expression_parser_test.go (updated to use pile)
  - parser/parser_test.go (updated to use pile)
  - parser/simplified_test.go (updated to use pile)
- Test utility files:
  - tests/expression_test.go (updated to use pile)
  - tests/expression_tester.go (updated to use pile)
- AST files:
  - ast/ast.go (updated to use pile)

## Implementation Strategy (Completed)
- Transition one component at a time
- Run tests after each component is updated
- Maintain backward compatibility until all code is migrated

## Progress
- ✅ Successfully updated all packages to use pile instead of core/classes:
  - VM package: All implementation and test files use pile
  - Compiler package: All implementation and test files use pile
  - Parser package: All implementation and test files use pile
  - AST package: Already using pile
- ✅ All tests are now passing
- ✅ The factorial demo has been moved to demo package and updated to use pile
- ✅ Main executable builds and runs correctly
- ✅ core and classes packages have been backed up and removed

## Final Cleanup (Completed)
- ✅ Ensure all imports from core and classes are replaced with pile
- ✅ Delete core and classes packages completely
- ✅ Ensure all tests pass
- ✅ Ensure the main executable builds and runs correctly

## Project Status
✅ **Migration Complete**: All code now uses the pile package instead of the deprecated core/classes packages. The codebase is now cleaner and more maintainable, with circular dependencies removed.