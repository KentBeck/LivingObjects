# Plan to Resolve Circular Module Dependency

After analyzing the code, I've identified the circular dependency:

1. The VM package needs to import the compiler package to use `NewMethodBuilder` for primitive method initialization
2. The compiler package imports the VM package to use bytecode constants (PUSH_LITERAL, PUSH_INSTANCE_VARIABLE, etc.)

## Solution Options

### Option 1: Create a Common Package for Shared Constants

1. Create a new package `src/interpreter/bytecode` that contains all bytecode constants
2. Move the bytecode constants from `vm/bytecode.go` to this new package
3. Update both VM and compiler packages to import this common package
4. This breaks the circular dependency by having both packages depend on a third package

### Option 2: Move Method Builder to VM Package

1. Move the MethodBuilder implementation from the compiler package to the VM package
2. This eliminates the need for VM to import compiler
3. However, this might not be ideal for code organization

### Option 3: Create a Factory in VM Package

1. Create a factory function in the VM package that creates methods with primitives
2. This function would not use the compiler's MethodBuilder
3. This keeps the code organization clean but requires duplicating some functionality

### Option 4: Combine Packages Temporarily

If the above solutions are too complex for immediate implementation:

1. Move the MethodBuilder code from compiler package into the VM package temporarily
2. Uncomment the primitive method initializations
3. Plan a proper refactoring for the future

## Recommended Approach

I recommend Option 1 (Create a Common Package) as the cleanest solution:

1. Create `src/interpreter/bytecode/constants.go` with all bytecode constants
2. Update imports in both VM and compiler packages
3. This is a clean solution that maintains good code organization

However, if you prefer a quicker solution, Option 4 (Combine Packages Temporarily) would allow us to make progress now and refactor later.
