# Go Interpreter Restructuring Plan

This document outlines a comprehensive plan for restructuring the Go implementation of the Smalltalk interpreter to improve modularity and maintainability.

## Proposed Directory Structure

```
src/
└── interpreter/
    ├── main.go                  # Main entry point
    ├── go.mod                   # Module definition
    ├── Makefile                 # Build instructions
    ├── README.md                # Documentation
    ├── core/                    # Core functionality
    │   ├── object.go            # Base Object type and interface
    │   ├── object_test.go       # Tests for base Object
    │   ├── memory.go            # Memory management
    │   ├── memory_test.go       # Memory tests
    │   ├── immediate.go         # Immediate value handling
    │   └── immediate_test.go    # Immediate value tests
    ├── classes/                 # Smalltalk class implementations
    │   ├── array.go             # Array implementation
    │   ├── array_test.go        # Array tests
    │   ├── dictionary.go        # Dictionary implementation
    │   ├── dictionary_test.go   # Dictionary tests
    │   ├── string.go            # String implementation
    │   ├── string_test.go       # String tests
    │   ├── symbol.go            # Symbol implementation
    │   ├── symbol_test.go       # Symbol tests
    │   ├── block.go             # Block implementation
    │   ├── block_test.go        # Block tests
    │   ├── class.go             # Class implementation
    │   ├── class_test.go        # Class tests
    │   ├── method.go            # Method implementation
    │   └── method_test.go       # Method tests
    ├── vm/                      # Virtual Machine
    │   ├── context.go           # Execution context
    │   ├── context_test.go      # Context tests
    │   ├── bytecode.go          # Bytecode definitions
    │   ├── bytecode_test.go     # Bytecode tests
    │   ├── bytecode_handlers.go # Bytecode execution handlers
    │   ├── bytecode_handlers_test.go # Handler tests
    │   ├── vm.go                # VM implementation
    │   ├── vm_test.go           # VM tests
    │   ├── primitives.go        # Primitive methods
    │   └── primitives_test.go   # Primitive tests
    ├── compiler/                # Compiler subsystem
    │   ├── method_builder.go    # Method building
    │   ├── method_builder_test.go # Method builder tests
    │   ├── parser.go            # Smalltalk parser
    │   └── parser_test.go       # Parser tests
    └── utils/                   # Utilities
        ├── conversions.go       # Type conversions
        ├── conversions_test.go  # Conversion tests
        ├── image.go             # Image loading/saving
        └── image_test.go        # Image tests
```

## Implementation Plan

### Phase 1: Setup New Directory Structure

First, create the new directory structure without moving any files yet:

```
src/interpreter/
├── core/
├── classes/
├── vm/
├── compiler/
└── utils/
```

### Phase 2: Refactor and Migrate Core Components

#### Step 1: Create Core Package

1. Create `core/object.go`:

   - Move the base `Object` struct, `ObjectType` enum, and `ObjectInterface` from the current `object.go`
   - Update package declaration to `package core`
   - Keep core object functionality (getters/setters)

2. Create `core/memory.go`:

   - Move memory management code from the current `memory.go`
   - Update package declaration and imports

3. Create `core/immediate.go`:
   - Move immediate value handling from the current `immediate.go`
   - Update package declaration and imports

#### Step 2: Create Classes Package

1. Create `classes/array.go`:

   - Move `Array` struct and related functions from `object.go`
   - Update package declaration to `package classes`
   - Add imports for the core package

2. Create `classes/dictionary.go`:

   - Move `Dictionary` struct and related functions from `object.go`
   - Update package declaration and imports

3. Create `classes/string.go`:

   - Move `String` struct and related functions from `object.go`
   - Update package declaration and imports

4. Create `classes/symbol.go`:

   - Move `Symbol` struct and related functions from `object.go`
   - Update package declaration and imports

5. Create `classes/block.go`:

   - Move `Block` struct and related functions from `object.go`
   - Update package declaration and imports

6. Create `classes/class.go`:

   - Move `Class` struct and related functions from `object.go`
   - Update package declaration and imports

7. Create `classes/method.go`:
   - Move `Method` struct and related functions from `object.go`
   - Update package declaration and imports

#### Step 3: Create VM Package

1. Create `vm/context.go`:

   - Move `Context` struct and related functions from `context.go`
   - Update package declaration to `package vm`
   - Add imports for the core and classes packages

2. Create `vm/bytecode.go`:

   - Move bytecode constants and related functions from `bytecode.go`
   - Update package declaration and imports

3. Create `vm/bytecode_handlers.go`:

   - Move bytecode handler functions from `bytecode_handlers.go`
   - Update package declaration and imports

4. Create `vm/vm.go`:

   - Move VM implementation from existing files
   - Update package declaration and imports

5. Create `vm/primitives.go`:
   - Move primitive method implementations from `primitives.go`
   - Update package declaration and imports

#### Step 4: Create Compiler Package

1. Create `compiler/method_builder.go`:

   - Move method building code from `method_builder.go`
   - Update package declaration to `package compiler`
   - Add imports for the core and classes packages

2. Create `compiler/parser.go` (if applicable):
   - Add Smalltalk parser implementation
   - Update package declaration and imports

#### Step 5: Create Utils Package

1. Create `utils/conversions.go`:

   - Move type conversion functions from `struct_conversions.go`
   - Update package declaration to `package utils`
   - Add imports for the core and classes packages

2. Create `utils/image.go`:
   - Move image loading/saving code from `image.go`
   - Update package declaration and imports

### Phase 3: Update Tests

For each implementation file, move the corresponding test file to the same package:

1. Move and update tests for core components:

   - `core/object_test.go`
   - `core/memory_test.go`
   - `core/immediate_test.go`

2. Move and update tests for class implementations:

   - `classes/array_test.go`
   - `classes/dictionary_test.go`
   - `classes/string_test.go`
   - `classes/symbol_test.go`
   - `classes/block_test.go`
   - `classes/class_test.go`
   - `classes/method_test.go`

3. Move and update tests for VM components:

   - `vm/context_test.go`
   - `vm/bytecode_test.go`
   - `vm/bytecode_handlers_test.go`
   - `vm/vm_test.go`
   - `vm/primitives_test.go`

4. Move and update tests for compiler components:

   - `compiler/method_builder_test.go`
   - `compiler/parser_test.go`

5. Move and update tests for utility functions:
   - `utils/conversions_test.go`
   - `utils/image_test.go`

### Phase 4: Update Main Package

1. Update `main.go`:

   - Update imports to reference the new package structure
   - Ensure the main function correctly initializes and uses components from the new packages

2. Update `go.mod`:
   - Ensure module dependencies are correctly specified

### Phase 5: Code Refactoring Guidelines

When moving code to the new structure, follow these guidelines:

1. **Package Declarations**: Update package declarations at the top of each file to match the new package name.

2. **Imports**: Update import statements to reference the new package locations.

3. **Exported vs. Unexported**: Review which functions and types should be exported (capitalized) based on their usage across packages.

4. **Interface Implementations**: Ensure that interface implementations are properly maintained across package boundaries.

5. **Circular Dependencies**: Avoid circular dependencies between packages. If necessary, refactor code to break cycles.

6. **Error Handling**: Maintain consistent error handling patterns across packages.

7. **Documentation**: Update documentation comments to reflect the new structure and package relationships.

### Phase 6: Testing Strategy

To ensure the refactoring doesn't break existing functionality:

1. **Incremental Testing**: After moving each component, run its tests to verify functionality.

2. **Integration Testing**: After completing each phase, run all tests to ensure components work together.

3. **Benchmark Testing**: Run benchmarks before and after refactoring to ensure performance is maintained.

4. **Manual Testing**: Manually test key scenarios to verify end-to-end functionality.

### Phase 7: Documentation Updates

1. Update `README.md` to describe the new structure and organization.

2. Add package-level documentation comments to explain the purpose and responsibilities of each package.

3. Create a migration guide for contributors to understand the new structure.

## Example: Refactoring the Object System

Let's look at a concrete example of how to refactor the object system:

### Current `object.go`:

```go
package main

import (
    "fmt"
    "unsafe"
)

// ObjectType represents the type of a Smalltalk object
type ObjectType int

const (
    OBJ_INTEGER ObjectType = iota
    OBJ_BOOLEAN
    OBJ_NIL
    OBJ_STRING
    OBJ_ARRAY
    OBJ_DICTIONARY
    OBJ_BLOCK
    OBJ_INSTANCE
    OBJ_CLASS
    OBJ_METHOD
    OBJ_SYMBOL
)

// Object represents a Smalltalk object
type Object struct {
    type1         ObjectType
    class         *Class
    moved         bool      // Used for garbage collection
    forwardingPtr *Object   // Used for garbage collection
    instanceVars  []*Object // Instance variables stored by index
}

// ... methods for Object ...

// String represents a Smalltalk string object
type String struct {
    Object
    Value string
}

// ... methods for String ...

// Array represents a Smalltalk array object
type Array struct {
    Object
    Elements []*Object
}

// ... methods for Array ...
```

### New `core/object.go`:

```go
package core

import (
    "fmt"
)

// ObjectType represents the type of a Smalltalk object
type ObjectType int

const (
    OBJ_INTEGER ObjectType = iota
    OBJ_BOOLEAN
    OBJ_NIL
    OBJ_STRING
    OBJ_ARRAY
    OBJ_DICTIONARY
    OBJ_BLOCK
    OBJ_INSTANCE
    OBJ_CLASS
    OBJ_METHOD
    OBJ_SYMBOL
)

// Object represents a Smalltalk object
type Object struct {
    Type1         ObjectType
    Class         *Class
    Moved         bool      // Used for garbage collection
    ForwardingPtr *Object   // Used for garbage collection
    InstanceVars  []*Object // Instance variables stored by index
}

// ObjectInterface defines the interface for all Smalltalk objects
type ObjectInterface interface {
    Type() ObjectType
    SetType(t ObjectType)
    Class() *Object
    SetClass(class *Object)
    Moved() bool
    SetMoved(moved bool)
    ForwardingPtr() *Object
    SetForwardingPtr(ptr *Object)
    InstanceVars() []*Object
    GetInstanceVarByIndex(index int) *Object
    SetInstanceVarByIndex(index int, value *Object)
    IsTrue() bool
    String() string
}

// ... methods for Object ...
```

### New `classes/string.go`:

```go
package classes

import (
    "fmt"
    "unsafe"

    "path/to/interpreter/core"
)

// String represents a Smalltalk string object
type String struct {
    core.Object
    Value string
}

// NewString creates a new string object
func NewString(value string) *String {
    str := &String{
        Value: value,
    }
    str.Type1 = core.OBJ_STRING
    return str
}

// StringToObject converts a String to an Object
func StringToObject(s *String) *core.Object {
    return (*core.Object)(unsafe.Pointer(s))
}

// ObjectToString converts an Object to a String
func ObjectToString(o core.ObjectInterface) *String {
    return (*String)(unsafe.Pointer(o.(*core.Object)))
}

// ... other String-related functions ...
```

## Timeline and Milestones

To implement this restructuring in a manageable way, I recommend the following timeline:

1. **Week 1**: Set up directory structure and refactor core components
2. **Week 2**: Refactor class implementations and VM components
3. **Week 3**: Refactor compiler and utility components
4. **Week 4**: Update tests and documentation, perform final integration testing

## Benefits of This Structure

This structure offers several advantages:

1. **Clear separation of concerns**: Each package has a specific responsibility
2. **Improved maintainability**: Related code is grouped together
3. **Better discoverability**: Easier to find specific implementations
4. **Reduced coupling**: Dependencies between components are more explicit
5. **Easier testing**: Tests are kept close to the code they test
6. **Scalability**: New components can be added without disrupting existing code
7. **Consistency with Go standards**: Follows Go community best practices

## Potential Challenges and Solutions

### Challenge 1: Circular Dependencies

When refactoring, you might encounter circular dependencies between packages.

**Solution**:

- Identify shared interfaces or types that cause circular dependencies
- Move these shared elements to a common package (e.g., `core`)
- Use interfaces instead of concrete types when possible

### Challenge 2: Maintaining Backward Compatibility

The refactoring might break existing code that imports the current package structure.

**Solution**:

- Implement the changes incrementally
- Create adapter functions in the main package that forward to the new structure
- Deprecate old functions gradually rather than removing them immediately

### Challenge 3: Test Coverage During Refactoring

Ensuring tests continue to work during the refactoring process.

**Solution**:

- Run tests after each component is moved
- Temporarily duplicate critical tests if necessary
- Use integration tests to verify the system works as a whole

## Conclusion

This restructuring plan provides a clear roadmap for transforming the Go interpreter codebase into a more modular and maintainable structure. By following this phased approach and adhering to the refactoring guidelines, you can achieve a better organized codebase while minimizing the risk of introducing bugs.
