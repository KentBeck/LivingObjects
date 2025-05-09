# Smalltalk Interpreter Primitive Initialization Plan

## Task Overview

The task is to:

1. Uncomment all primitive method initializations in `src/interpreter/vm/vm.go` for classes like ObjectClass, IntegerClass, FloatClass, and BlockClass
2. Delete any primitive initializations that are not in the VM file
3. Check for circular module dependencies

## Current State Analysis

From reviewing the code, I found:

1. In `src/interpreter/vm/vm.go`:

   - There are commented-out primitive method initializations in several class initialization methods:
     - `NewObjectClass()` - Line ~60-64: commented-out basicClass method
     - `NewIntegerClass()` - Lines ~73-92: commented-out methods for +, -, \*, =, <, >
     - `NewFloatClass()` - Lines ~101-123: commented-out methods for +, -, \*, /, =, <, >
     - `NewBlockClass()` - Lines ~169-180: commented-out methods for new, value, value:

2. In `src/interpreter/classes/string.go`:

   - There's a String class implementation but no primitive initializations outside of VM

3. In `src/interpreter/vm/primitives.go`:

   - Contains the implementation of primitive methods but not their initialization

4. In `src/interpreter/classes/array.go`:
   - Contains Array class implementation but no primitive initializations

## Implementation Plan

### Step 1: Uncomment Primitive Method Initializations in VM.go

Uncomment all the primitive method initializations in the following methods:

- `NewObjectClass()`
- `NewIntegerClass()`
- `NewFloatClass()`
- `NewBlockClass()`

This will require:

1. Uncommenting the code blocks
2. Ensuring the method builder is properly imported and used
3. Checking for any syntax errors or missing dependencies

### Step 2: Check for Primitive Initializations Outside VM.go

In `src/interpreter/vm/vm.go`, the `NewStringClass()` method already has a non-commented primitive initialization for the string concatenation method (`,`). We need to check if there are any other primitive initializations in other files that should be deleted.

### Step 3: Check for Circular Module Dependencies

After making the changes, we need to check for any circular module dependencies that might be introduced. This could happen if:

- The VM package depends on classes package
- The classes package depends on VM package
- The compiler package depends on both

## Specific Code Changes

### In NewObjectClass()

```go
func (vm *VM) NewObjectClass() *classes.Class {
    result := classes.NewClass("Object", nil) // patch this up later. then even later when we have real images all this initialization can go away

    // Add basicClass method to Object class
    NewMethodBuilder(result).
        Selector("basicClass").
        Primitive(5). // basicClass primitive
        Go()

    return result
}
```

### In NewIntegerClass()

```go
func (vm *VM) NewIntegerClass() *classes.Class {
    result := classes.NewClass("Integer", vm.ObjectClass)

    // Add primitive methods to the Integer class
    builder := NewMethodBuilder(result)

    // + method (addition)
    builder.Selector("+").Primitive(1).Go()

    // - method (subtraction)
    builder.Selector("-").Primitive(4).Go()

    // * method (multiplication)
    builder.Selector("*").Primitive(2).Go()

    // = method (equality)
    builder.Selector("=").Primitive(3).Go()

    // < method (less than)
    builder.Selector("<").Primitive(6).Go()

    // > method (greater than)
    builder.Selector(">").Primitive(7).Go()

    return result
}
```

### In NewFloatClass()

```go
func (vm *VM) NewFloatClass() *classes.Class {
    result := classes.NewClass("Float", vm.ObjectClass) // patch this up later. then even later when we have real images all this initialization can go away

    // Add primitive methods to the Float class
    builder := NewMethodBuilder(result)

    // + method (addition)
    builder.Selector("+").Primitive(10).Go()

    // - method (subtraction)
    builder.Selector("-").Primitive(11).Go()

    // * method (multiplication)
    builder.Selector("*").Primitive(12).Go()

    // / method (division)
    builder.Selector("/").Primitive(13).Go()

    // = method (equality)
    builder.Selector("=").Primitive(14).Go()

    // < method (less than)
    builder.Selector("<").Primitive(15).Go()

    // > method (greater than)
    builder.Selector(">").Primitive(16).Go()

    return result
}
```

### In NewBlockClass()

```go
func (vm *VM) NewBlockClass() *classes.Class {
    result := classes.NewClass("Block", vm.ObjectClass)

    // Add primitive methods to the Block class
    builder := NewMethodBuilder(result)

    // new method (creates a new block instance)
    // fixme sketchy
    builder.Selector("new").Primitive(20).Go()

    // value method (executes the block with no arguments)
    builder.Selector("value").Primitive(21).Go()

    // value: method (executes the block with one argument)
    builder.Selector("value:").Primitive(22).Go()

    return result
}
```

## Potential Issues and Solutions

1. **Import Conflicts**: We may need to ensure that the `NewMethodBuilder` function is properly imported from the compiler package.

2. **Circular Dependencies**: If uncommenting these methods creates circular dependencies, we may need to refactor the code to break these dependencies.

3. **Consistency**: We need to ensure that all primitive methods are consistently initialized in the VM file and not elsewhere.

## Next Steps

After implementing these changes, we should:

1. Run tests to ensure the interpreter still works correctly
2. Check for any circular dependencies or import errors
3. Verify that all primitive methods are properly initialized
