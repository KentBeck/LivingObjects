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
  - vm/bytecode_test.go (no changes needed as it doesn't use core/classes)
- Demo files:
  - factorial_demo.go (moved to demo package and updated to use pile)
  - demo/factorial.go (updated to use pile)

## Next Steps

1. **Update Tests**
   - Update test files to use the pile package instead of core/classes
   - Update any test helper code to use the pile package

2. **Update VM Package**
   - Update VM struct to use pile types instead of core
   - Update all VM methods to use pile
   - Update VM factory and helper methods to use pile
   - Update context handling to use pile

3. **Update Compiler Package**
   - Update bytecode compiler to use pile
   - Update method builder to use pile
   - Update VM access components to use pile

4. **Update Parser Package**
   - Update parser to use pile instead of core/classes

5. **Remove Deprecated Files**
   - Once all files are updated and tests pass, remove the deprecated.go files
   - Remove empty core and classes directories

## Implementation Strategy
- Transition one component at a time
- Run tests after each component is updated
- Maintain backward compatibility until all code is migrated

## Progress
- Successfully updated all VM test files to use the pile package
- All updated tests are passing
- The factorial demo has been moved to demo package and updated to use pile
- Main executable now builds and runs correctly

## Next Files to Update
1. All VM test files have been updated successfully. Next, update implementation files:
   - vm/vm.go (currently partially updated)
   - vm/bytecode_handlers.go (currently partially updated)
   - vm/block_bytecode_handlers.go
   - vm/block_executor.go
   - vm/class_accessors.go

## Final Cleanup
- Ensure all imports from core and classes are replaced with pile
- Delete core and classes packages completely
- Update build and documentation