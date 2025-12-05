# Smalltalk VM Specifications

**Machine-Readable Specifications for AI Agents and VM Implementers**

This directory contains comprehensive, machine-readable specifications for implementing a Smalltalk virtual machine. The specifications are designed to be:

- **Unambiguous**: Formal definitions with no room for interpretation
- **Executable**: Test suites that verify correctness
- **Parseable**: JSON schemas that agents can process
- **Complete**: Every aspect of the VM is specified

## Directory Structure

```
specs/
├── primitives/          # Primitive method specifications
│   ├── schema/          # JSON schemas for each primitive
│   ├── tests/           # Executable test suites
│   └── README.md        # Primitive documentation
├── bytecode/            # Bytecode instruction specifications
│   ├── schema/          # JSON schemas for each instruction
│   ├── tests/           # Executable test suites
│   └── README.md        # Bytecode documentation
├── image-format/        # Image file format specifications
│   ├── schema/          # JSON schema for image format
│   ├── tests/           # Image I/O test suites
│   └── README.md        # Image format documentation
└── README.md            # This file
```

## For AI Agents

### Quick Start

1. **Parse JSON Schemas**: Load primitive/bytecode definitions from `schema/` directories
2. **Generate Implementation**: Use schemas to generate VM code
3. **Run Tests**: Execute test suites to verify correctness
4. **Iterate**: Fix failures until all tests pass

### Schema Format

Each primitive/bytecode has a JSON schema with:

```json
{
  "primitive_number": 1,
  "selector": "+",
  "receiver": { "type": "SmallInteger" },
  "arguments": [{ "type": "SmallInteger" }],
  "returns": { "type": "SmallInteger" },
  "stack_effect": {
    "before": ["...", "receiver", "arg"],
    "after": ["...", "result"],
    "pops": 2,
    "pushes": 1
  },
  "operation": {
    "pseudocode": "result = receiver + arg",
    "steps": ["Pop arg", "Pop receiver", "Compute sum", "Push result"]
  },
  "preconditions": ["receiver is SmallInteger", "arg is SmallInteger"],
  "postconditions": ["result == receiver + arg"],
  "failures": [{ "condition": "overflow", "error": "OverflowError" }],
  "test_cases": [{ "receiver": 3, "arguments": [4], "expected": 7 }]
}
```

### Test Suite Format

Executable Python tests define expected behavior:

```python
def test_basic_addition(vm: VMInterface):
    """3 + 4 should return 7"""
    vm.push(3)  # receiver
    vm.push(4)  # argument
    result = vm.call_primitive(1)

    assert result == 7
    assert vm.stack == [7]
```

## For Human Implementers

### Implementation Workflow

1. **Read Human Documentation**: Start with markdown docs in project root
2. **Consult JSON Schemas**: Get precise specifications for each component
3. **Implement Components**: Write VM code following specifications
4. **Adapt Test Interface**: Implement `VMInterface` for your VM
5. **Run Test Suites**: Execute tests to verify correctness
6. **Debug Failures**: Fix issues until all tests pass

### Test Execution

```bash
# Install dependencies
pip install pytest

# Run all tests
pytest specs/primitives/tests/ -v
pytest specs/bytecode/tests/ -v

# Run specific test file
pytest specs/primitives/tests/test_integer_primitives.py -v

# Run specific test
pytest specs/primitives/tests/test_integer_primitives.py::TestPrimitive1_Add::test_basic_addition -v
```

## Specification Coverage

### Primitives (40+ methods)

| Range      | Category              | Count | Status       |
| ---------- | --------------------- | ----- | ------------ |
| 1-11       | Integer arithmetic    | 11    | ✅ Specified |
| 60-62      | Array operations      | 3     | ✅ Specified |
| 63-67      | String operations     | 5     | ✅ Specified |
| 70-75, 111 | Object operations     | 6     | ✅ Specified |
| 154-159    | Boolean conditionals  | 6     | ✅ Specified |
| 201-202    | Block execution       | 2     | ✅ Specified |
| 700-703    | Dictionary operations | 4     | ✅ Specified |
| 1000-1001  | Exception handling    | 2     | ✅ Specified |
| 5000+      | System operations     | 2+    | ✅ Specified |

### Bytecode Instructions (15 opcodes)

| Opcode | Mnemonic                 | Status       |
| ------ | ------------------------ | ------------ |
| 0      | PUSH_LITERAL             | ✅ Specified |
| 1      | PUSH_INSTANCE_VARIABLE   | ✅ Specified |
| 2      | PUSH_TEMPORARY_VARIABLE  | ✅ Specified |
| 3      | PUSH_SELF                | ✅ Specified |
| 4      | STORE_INSTANCE_VARIABLE  | ✅ Specified |
| 5      | STORE_TEMPORARY_VARIABLE | ✅ Specified |
| 6      | SEND_MESSAGE             | ✅ Specified |
| 7      | RETURN_STACK_TOP         | ✅ Specified |
| 8      | JUMP                     | ✅ Specified |
| 9      | JUMP_IF_TRUE             | ✅ Specified |
| 10     | JUMP_IF_FALSE            | ✅ Specified |
| 11     | POP                      | ✅ Specified |
| 12     | DUPLICATE                | ✅ Specified |
| 13     | CREATE_BLOCK             | ✅ Specified |
| 14     | EXECUTE_BLOCK            | ✅ Specified |

## Compliance Testing

A VM implementation is compliant if it passes all test suites:

- ✅ All primitive tests pass
- ✅ All bytecode tests pass
- ✅ All image format tests pass
- ✅ All integration tests pass

## Schema Validation

Schemas can be validated using JSON Schema validators:

```python
import json
import jsonschema

# Load schema
with open('specs/primitives/schema/primitive-001-add.json') as f:
    schema = json.load(f)

# Validate against meta-schema
jsonschema.Draft7Validator.check_schema(schema)
```

## Generating Code from Schemas

Example: Generate primitive dispatch code

```python
import json
import glob

def generate_primitive_dispatch():
    code = "def call_primitive(self, num, receiver, args):\n"

    for schema_file in glob.glob('specs/primitives/schema/*.json'):
        with open(schema_file) as f:
            spec = json.load(f)

        num = spec['primitive_number']
        selector = spec['selector']

        code += f"    if num == {num}:  # {selector}\n"
        code += f"        return self.primitive_{num}(receiver, args)\n"

    return code
```

## Contributing

When adding new primitives or bytecodes:

1. Create JSON schema in `schema/` directory
2. Add test cases to test suite
3. Update this README with new entries
4. Ensure all tests pass

## References

- **VM Specification**: `../smalltalk-vm-specification.md`
- **Primitives Specification**: `../smalltalk-primitives-specification.md`
- **Image Build Plan**: `../smalltalk-image-build-plan.md`

## License

These specifications are part of the Smalltalk VM project and follow the same license.

---

**For questions or issues, please refer to the main project documentation.**
