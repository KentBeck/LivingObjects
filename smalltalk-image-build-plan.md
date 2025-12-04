# Smalltalk Image Build System - Detailed Implementation Plan

## Overview

Create a build system that reads Smalltalk source files, compiles them to bytecode, builds an executable image, and runs tests within that image. This directly advances the self-hosting goal while leveraging the existing C++ VM as the execution engine.

## Current State Analysis

### What We Have ✅

- **C++ VM**: 93.75% expression test success rate (60/64 tests passing)
- **Parser & Compiler**: Working Smalltalk parser and bytecode compiler in C++
- **Core Collections**: OrderedCollection, Dictionary, ByteArray fully implemented
- **SUnit Framework**: Complete testing framework ready to use
- **Expression Tests**: Comprehensive test suite covering all language features
- **Basic Class System**: C++ classes for Object, Class, Method, etc.

### What We Need 🎯

- **Source File Loader**: Read and parse `.st` files
- **Class Installation**: Create classes from Smalltalk source
- **Method Compilation**: Compile Smalltalk methods to bytecode
- **Image Building**: Serialize complete object graph
- **Image Execution**: Load and run code from images
- **Test Integration**: Run SUnit tests in loaded images

---

## Phase 1: Source File Loading & Parsing (Week 1)

### 1.1 Smalltalk Source File Loader

**Goal**: Read and parse Smalltalk source files into AST structures

**Components to Build**:

#### SourceFileLoader Class

```cpp
class SourceFileLoader {
public:
    struct ClassDefinition {
        std::string className;
        std::string superclassName;
        std::vector<std::string> instanceVariables;
        std::vector<std::string> classVariables;
        std::string packageName;
        std::vector<MethodDefinition> methods;
    };

    struct MethodDefinition {
        std::string selector;
        std::string sourceCode;
        std::vector<std::string> arguments;
        std::vector<std::string> temporaries;
        bool isClassMethod;
    };

    std::vector<ClassDefinition> loadSourceFile(const std::string& filename);
    std::vector<ClassDefinition> loadSourceDirectory(const std::string& directory);
};
```

#### Smalltalk Source Parser Enhancement

```cpp
class SmalltalkSourceParser {
public:
    ClassDefinition parseClassDefinition(const std::string& source);
    MethodDefinition parseMethodDefinition(const std::string& source);
    std::vector<std::string> parseInstanceVariables(const std::string& declaration);

private:
    void parseClassHeader(const std::string& header, ClassDefinition& classDef);
    void parseMethodsSection(const std::string& section, ClassDefinition& classDef);
};
```

**Implementation Tasks**:

1. **File Reading**: Read `.st` files and handle file I/O errors
2. **Class Definition Parsing**: Parse `Object subclass: #MyClass` syntax
3. **Instance Variable Parsing**: Parse `instanceVariableNames: 'var1 var2'`
4. **Method Section Parsing**: Parse `!MyClass methodsFor: 'category'!` sections
5. **Method Body Parsing**: Extract method selector, arguments, and body
6. **Error Handling**: Provide clear error messages with line numbers

**Acceptance Criteria**:

- Can load all existing `.st` files in `src/` directory
- Correctly parses class definitions with inheritance
- Extracts all method definitions with proper selectors
- Handles parsing errors gracefully with useful error messages
- Processes multiple files and resolves dependencies

### 1.2 Source File Discovery & Dependency Resolution

**Goal**: Automatically discover source files and resolve class dependencies

**Components to Build**:

#### SourceFileManager

```cpp
class SourceFileManager {
public:
    struct LoadOrder {
        std::vector<std::string> orderedFiles;
        std::map<std::string, std::vector<std::string>> dependencies;
    };

    LoadOrder calculateLoadOrder(const std::vector<std::string>& sourceFiles);
    std::vector<std::string> findSourceFiles(const std::string& directory);

private:
    void buildDependencyGraph(const std::vector<ClassDefinition>& classes);
    std::vector<std::string> topologicalSort();
};
```

**Implementation Tasks**:

1. **File Discovery**: Recursively find all `.st` files
2. **Dependency Analysis**: Build graph of class dependencies (superclass relationships)
3. **Topological Sort**: Order classes so superclasses are loaded before subclasses
4. **Circular Dependency Detection**: Detect and report circular dependencies
5. **Missing Class Detection**: Report references to undefined classes

**Acceptance Criteria**:

- Discovers all `.st` files in source directory
- Correctly orders classes by dependency (Object first, subclasses after)
- Detects circular dependencies and reports errors
- Handles forward references appropriately

---

## Phase 2: Class Installation & Method Compilation (Week 2)

### 2.1 Dynamic Class Creation

**Goal**: Create C++ Class objects from Smalltalk source definitions

**Components to Build**:

#### ClassInstaller

```cpp
class ClassInstaller {
public:
    Class* installClass(const ClassDefinition& classDef, SmalltalkImage& image);
    void installMethod(Class* targetClass, const MethodDefinition& methodDef);

private:
    Class* createClassObject(const std::string& name, Class* superclass);
    void setInstanceVariables(Class* clazz, const std::vector<std::string>& ivars);
    Symbol* internSymbol(const std::string& symbolName);
};
```

#### Enhanced Method Compiler

```cpp
class SmalltalkMethodCompiler {
public:
    std::shared_ptr<CompiledMethod> compileMethod(
        const MethodDefinition& methodDef,
        Class* targetClass,
        SmalltalkImage& image
    );

private:
    void compileMethodBody(const std::string& source, CompiledMethod& method);
    void resolveVariableReferences(CompiledMethod& method, Class* targetClass);
    void generateBytecode(const ParseNode& ast, CompiledMethod& method);
};
```

**Implementation Tasks**:

1. **Class Object Creation**: Create Class instances with proper metadata
2. **Superclass Linking**: Establish inheritance relationships
3. **Instance Variable Setup**: Configure instance variable layouts
4. **Method Dictionary Creation**: Initialize empty method dictionaries
5. **Symbol Interning**: Ensure symbols are properly interned and unique
6. **Method Compilation**: Compile method source to bytecode
7. **Method Installation**: Add compiled methods to class method dictionaries

**Acceptance Criteria**:

- Creates proper Class objects with correct inheritance hierarchy
- Instance variables are correctly configured and accessible
- Methods compile to valid bytecode
- Method lookup works correctly through inheritance chain
- Symbol table is properly maintained

### 2.2 Primitive Method Integration

**Goal**: Integrate existing C++ primitive methods with Smalltalk classes

**Components to Build**:

#### PrimitiveMethodInstaller

```cpp
class PrimitiveMethodInstaller {
public:
    void installCorePrimitives(SmalltalkImage& image);
    void installIntegerPrimitives(Class* integerClass);
    void installStringPrimitives(Class* stringClass);
    void installArrayPrimitives(Class* arrayClass);
    void installBooleanPrimitives(Class* trueClass, Class* falseClass);

private:
    void addPrimitiveMethod(Class* clazz, const std::string& selector, int primitiveNumber);
};
```

**Implementation Tasks**:

1. **Core Primitive Installation**: Install Object, Class primitives
2. **Integer Primitive Installation**: Install arithmetic and comparison methods
3. **String Primitive Installation**: Install concatenation, size, access methods
4. **Array Primitive Installation**: Install at:, at:put:, size methods
5. **Boolean Primitive Installation**: Install ifTrue:, ifFalse:, ifTrue:ifFalse: methods
6. **Block Primitive Installation**: Install value, value: methods

**Acceptance Criteria**:

- All existing primitive methods are available in Smalltalk classes
- Method dispatch correctly calls primitive implementations
- Primitive failures fall back to Smalltalk implementations when available
- Boolean conditionals work correctly with blocks

---

## Phase 3: Image Building & Serialization (Week 3)

### 3.1 Object Graph Serialization

**Goal**: Serialize complete Smalltalk object graph to persistent storage

**Components to Build**:

#### ImageBuilder

```cpp
class ImageBuilder {
public:
    struct ImageHeader {
        uint32_t version;
        uint32_t objectCount;
        uint32_t symbolTableOffset;
        uint32_t classTableOffset;
        uint32_t globalDictionaryOffset;
    };

    void buildImage(const std::string& filename, SmalltalkImage& image);
    void serializeObjectGraph(std::ostream& stream, SmalltalkImage& image);

private:
    void writeHeader(std::ostream& stream, const ImageHeader& header);
    void writeObjectTable(std::ostream& stream, SmalltalkImage& image);
    void writeSymbolTable(std::ostream& stream, SmalltalkImage& image);
    void writeClassTable(std::ostream& stream, SmalltalkImage& image);
};
```

#### Object Graph Traversal

```cpp
class ObjectGraphTraverser {
public:
    std::vector<Object*> collectAllObjects(SmalltalkImage& image);
    void resolveCircularReferences(std::vector<Object*>& objects);

private:
    void traverseObject(Object* obj, std::set<Object*>& visited);
    void assignObjectIDs(std::vector<Object*>& objects);
};
```

**Implementation Tasks**:

1. **Object Graph Traversal**: Find all reachable objects from roots
2. **Circular Reference Handling**: Detect and properly serialize circular references
3. **Object ID Assignment**: Assign unique IDs to all objects for serialization
4. **Binary Format Design**: Design efficient binary format for objects
5. **Symbol Table Serialization**: Serialize symbol table with proper interning
6. **Class Table Serialization**: Serialize class hierarchy and method dictionaries
7. **Global Dictionary Serialization**: Serialize global variable bindings

**Acceptance Criteria**:

- Can serialize complete object graph without data loss
- Handles circular references correctly
- Produces compact, efficient binary format
- Symbol table maintains proper interning relationships
- Class hierarchy is preserved correctly

### 3.2 Image Loading & Deserialization

**Goal**: Load serialized images and reconstruct object graph

**Components to Build**:

#### ImageLoader

```cpp
class ImageLoader {
public:
    std::unique_ptr<SmalltalkImage> loadImage(const std::string& filename);
    void deserializeObjectGraph(std::istream& stream, SmalltalkImage& image);

private:
    ImageHeader readHeader(std::istream& stream);
    void readObjectTable(std::istream& stream, SmalltalkImage& image);
    void readSymbolTable(std::istream& stream, SmalltalkImage& image);
    void readClassTable(std::istream& stream, SmalltalkImage& image);
    void resolveObjectReferences(SmalltalkImage& image);
};
```

**Implementation Tasks**:

1. **Binary Format Reading**: Read binary image format correctly
2. **Object Reconstruction**: Recreate objects with proper types and data
3. **Reference Resolution**: Resolve object references and circular references
4. **Symbol Table Reconstruction**: Rebuild symbol table with proper interning
5. **Class Hierarchy Reconstruction**: Rebuild class hierarchy and method dictionaries
6. **Global Dictionary Reconstruction**: Restore global variable bindings
7. **Primitive Method Reconnection**: Reconnect primitive methods to implementations

**Acceptance Criteria**:

- Can load images created by ImageBuilder
- All objects are correctly reconstructed with proper types
- Object references and circular references work correctly
- Symbol table maintains proper interning
- Class hierarchy and method lookup work correctly
- Primitive methods are properly connected

---

## Phase 4: Build System Integration (Week 4)

### 4.1 Build Tool Implementation

**Goal**: Create command-line tool for building Smalltalk images

**Components to Build**:

#### SmalltalkImageBuilder Tool

```cpp
class SmalltalkImageBuilder {
public:
    struct BuildOptions {
        std::string sourceDirectory;
        std::string outputImage;
        std::vector<std::string> includePaths;
        bool verbose;
        bool runTests;
    };

    int main(int argc, char** argv);
    void buildImage(const BuildOptions& options);

private:
    void parseCommandLine(int argc, char** argv, BuildOptions& options);
    void loadSources(const BuildOptions& options);
    void compileClasses();
    void runTests();
    void saveImage(const std::string& filename);
};
```

**Implementation Tasks**:

1. **Command Line Parsing**: Parse build options and source directories
2. **Source Loading**: Load all Smalltalk source files in dependency order
3. **Class Installation**: Install all classes and methods
4. **Primitive Integration**: Connect primitive methods
5. **Image Building**: Build and save final image
6. **Error Reporting**: Provide clear error messages for build failures
7. **Verbose Output**: Optional detailed build progress reporting

**Acceptance Criteria**:

- Command-line tool builds images from source directories
- Handles build errors gracefully with clear messages
- Supports verbose output for debugging
- Produces working images that can be loaded and executed

### 4.2 Makefile Integration

**Goal**: Integrate image building into existing build system

**Makefile Targets**:

```makefile
# Build Smalltalk image from sources
smalltalk-image: build/smalltalk-image-builder src/*.st
	./build/smalltalk-image-builder --source src --output build/smalltalk.image --verbose

# Run tests in Smalltalk image
test-smalltalk: build/smalltalk.image
	./build/smalltalk-vm --load-image build/smalltalk.image --run "RunExpressionTests run"

# Clean Smalltalk build artifacts
clean-smalltalk:
	rm -f build/smalltalk.image build/smalltalk-image-builder

# Full build including Smalltalk
all: smalltalk-vm smalltalk-image-builder smalltalk-image

# Test everything
test: test-cpp test-smalltalk
```

**Implementation Tasks**:

1. **Image Builder Compilation**: Add image builder to build system
2. **Source Dependencies**: Track dependencies on `.st` files
3. **Incremental Building**: Rebuild image only when sources change
4. **Test Integration**: Run Smalltalk tests as part of build process
5. **Clean Targets**: Properly clean Smalltalk build artifacts

**Acceptance Criteria**:

- `make smalltalk-image` builds image from sources
- `make test-smalltalk` runs SUnit tests in image
- `make all` builds complete system including image
- `make test` runs both C++ and Smalltalk tests
- Incremental builds work correctly

---

## Phase 5: Test Execution & Reporting (Week 5)

### 5.1 Image Test Runner

**Goal**: Execute SUnit tests within loaded Smalltalk images

**Components to Build**:

#### ImageTestRunner

```cpp
class ImageTestRunner {
public:
    struct TestResults {
        int totalTests;
        int passedTests;
        int failedTests;
        int errorTests;
        std::vector<std::string> failures;
        std::vector<std::string> errors;
    };

    TestResults runTests(SmalltalkImage& image, const std::string& testExpression);
    void reportResults(const TestResults& results);

private:
    TaggedValue executeExpression(SmalltalkImage& image, const std::string& expression);
    TestResults parseTestResults(TaggedValue result);
};
```

**Implementation Tasks**:

1. **Test Expression Execution**: Execute test runner expressions in image
2. **Result Extraction**: Extract test results from Smalltalk objects
3. **Result Formatting**: Format results for console output
4. **Error Handling**: Handle test execution errors gracefully
5. **Exit Code Management**: Return appropriate exit codes for CI/CD

**Acceptance Criteria**:

- Can execute SUnit test expressions in loaded images
- Extracts and reports test results correctly
- Handles test failures and errors appropriately
- Returns proper exit codes for automated testing

### 5.2 Continuous Integration Integration

**Goal**: Integrate Smalltalk tests into CI/CD pipeline

**GitHub Actions Workflow**:

```yaml
name: Smalltalk Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build C++ VM
        run: make -C src/cpp all
      - name: Build Smalltalk Image
        run: make -C src/cpp smalltalk-image
      - name: Run C++ Tests
        run: make -C src/cpp test
      - name: Run Smalltalk Tests
        run: make -C src/cpp test-smalltalk
```

**Implementation Tasks**:

1. **CI Configuration**: Add Smalltalk testing to CI pipeline
2. **Test Reporting**: Generate test reports in CI-friendly format
3. **Failure Handling**: Ensure CI fails on test failures
4. **Artifact Management**: Save images and test reports as artifacts

**Acceptance Criteria**:

- CI runs both C++ and Smalltalk tests
- Test failures cause CI to fail
- Test results are clearly reported
- Build artifacts are properly managed

---

## Phase 6: Bootstrap Preparation (Week 6)

### 6.1 Self-Compilation Infrastructure

**Goal**: Prepare infrastructure for Smalltalk compiler to compile itself

**Components to Build**:

#### Bootstrap Compiler Integration

```cpp
class BootstrapCompiler {
public:
    void installSmalltalkCompiler(SmalltalkImage& image);
    bool canSelfCompile(SmalltalkImage& image);
    void testSelfCompilation(SmalltalkImage& image);

private:
    void loadCompilerClasses(SmalltalkImage& image);
    void testCompilerFunctionality(SmalltalkImage& image);
};
```

**Implementation Tasks**:

1. **Compiler Class Loading**: Load SmalltalkCompiler and related classes
2. **Compiler Testing**: Test that compiler can compile simple expressions
3. **Self-Compilation Test**: Test compiler compiling its own methods
4. **Bootstrap Verification**: Verify all components needed for self-hosting

**Acceptance Criteria**:

- SmalltalkCompiler class loads and functions correctly
- Can compile simple expressions and methods
- Can compile its own methods (self-compilation)
- All bootstrap components are present and working

### 6.2 Image Validation & Testing

**Goal**: Comprehensive validation of built images

**Components to Build**:

#### ImageValidator

```cpp
class ImageValidator {
public:
    struct ValidationResults {
        bool isValid;
        std::vector<std::string> errors;
        std::vector<std::string> warnings;
    };

    ValidationResults validateImage(SmalltalkImage& image);

private:
    void validateClassHierarchy(SmalltalkImage& image);
    void validateMethodDictionaries(SmalltalkImage& image);
    void validateSymbolTable(SmalltalkImage& image);
    void validatePrimitiveMethods(SmalltalkImage& image);
};
```

**Implementation Tasks**:

1. **Class Hierarchy Validation**: Verify inheritance relationships
2. **Method Dictionary Validation**: Verify all methods are properly installed
3. **Symbol Table Validation**: Verify symbol interning is correct
4. **Primitive Method Validation**: Verify primitive methods are connected
5. **Test Suite Validation**: Verify SUnit framework is complete and functional

**Acceptance Criteria**:

- Can validate image integrity comprehensively
- Detects and reports structural problems
- Validates that all expected classes and methods are present
- Confirms SUnit framework is ready for use

---

## Success Metrics & Validation

### Technical Metrics

- **Image Build Success**: Can build complete image from all `.st` files
- **Test Execution**: Can run full SUnit test suite in image
- **Expression Test Success**: All 64 expression tests pass in image
- **Self-Compilation Ready**: Infrastructure ready for bootstrap compiler

### Functional Metrics

- **Build Integration**: Image building integrated into make system
- **CI Integration**: Smalltalk tests run in continuous integration
- **Developer Experience**: Easy to add new Smalltalk classes and tests
- **Error Reporting**: Clear error messages for build and test failures

## Risk Mitigation

### Technical Risks

- **Serialization Complexity**: Start with simple format, optimize later
- **Circular References**: Use proven algorithms for graph traversal
- **Memory Management**: Leverage existing C++ memory management
- **Performance**: Focus on correctness first, optimize later

### Integration Risks

- **Build System Complexity**: Keep build steps simple and well-documented
- **CI/CD Integration**: Test locally before integrating with CI
- **Backward Compatibility**: Maintain existing C++ test functionality

## Timeline Summary

- **Week 1**: Source file loading and parsing
- **Week 2**: Class installation and method compilation
- **Week 3**: Image building and serialization
- **Week 4**: Build system integration
- **Week 5**: Test execution and reporting
- **Week 6**: Bootstrap preparation and validation

## Deliverables

### Phase Deliverables

- **Week 1**: Working source file loader and parser
- **Week 2**: Dynamic class creation and method compilation
- **Week 3**: Image serialization and loading
- **Week 4**: Complete build system integration
- **Week 5**: Automated test execution in images
- **Week 6**: Bootstrap-ready system with validation

### Final Deliverable

A complete Smalltalk image build system that:

1. **Reads Smalltalk source files** and resolves dependencies
2. **Compiles classes and methods** to bytecode
3. **Builds executable images** with complete object graphs
4. **Runs comprehensive test suites** within images
5. **Integrates with build system** and CI/CD pipeline
6. **Prepares for self-hosting** with bootstrap compiler infrastructure

This system will enable rapid Smalltalk development, comprehensive testing, and provide the foundation for the final self-hosting milestone.
