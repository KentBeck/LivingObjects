#include "../include/primitives.h"
#include "../include/interpreter.h"
#include "../include/method_compiler.h"
#include "../include/object.h"
#include "../include/smalltalk_class.h"
#include "../include/smalltalk_string.h"

// Forward declarations for primitive functions
namespace smalltalk {
namespace IntegerPrimitives {
TaggedValue add(TaggedValue receiver, const std::vector<TaggedValue> &args,
                Interpreter &interpreter);
TaggedValue subtract(TaggedValue receiver, const std::vector<TaggedValue> &args,
                     Interpreter &interpreter);
TaggedValue multiply(TaggedValue receiver, const std::vector<TaggedValue> &args,
                     Interpreter &interpreter);
TaggedValue divide(TaggedValue receiver, const std::vector<TaggedValue> &args,
                   Interpreter &interpreter);
TaggedValue lessThan(TaggedValue receiver, const std::vector<TaggedValue> &args,
                     Interpreter &interpreter);
TaggedValue greaterThan(TaggedValue receiver,
                        const std::vector<TaggedValue> &args,
                        Interpreter &interpreter);
TaggedValue equal(TaggedValue receiver, const std::vector<TaggedValue> &args,
                  Interpreter &interpreter);
TaggedValue notEqual(TaggedValue receiver, const std::vector<TaggedValue> &args,
                     Interpreter &interpreter);
TaggedValue lessThanOrEqual(TaggedValue receiver,
                            const std::vector<TaggedValue> &args,
                            Interpreter &interpreter);
TaggedValue greaterThanOrEqual(TaggedValue receiver,
                               const std::vector<TaggedValue> &args,
                               Interpreter &interpreter);
} // namespace IntegerPrimitives
namespace BlockPrimitives {
TaggedValue value(TaggedValue receiver, const std::vector<TaggedValue> &args,
                  Interpreter &interpreter);
}
} // namespace smalltalk

namespace smalltalk {

// PrimitiveRegistry implementation
PrimitiveRegistry &PrimitiveRegistry::getInstance() {
  static PrimitiveRegistry instance;
  return instance;
}

void PrimitiveRegistry::registerPrimitive(int primitiveNumber,
                                          PrimitiveFunction function) {
  primitives_[primitiveNumber] = function;
}

TaggedValue
PrimitiveRegistry::callPrimitive(int primitiveNumber, TaggedValue receiver,
                                 const std::vector<TaggedValue> &args,
                                 Interpreter &interpreter) {
  auto it = primitives_.find(primitiveNumber);
  if (it == primitives_.end()) {
    throw PrimitiveFailure("Primitive " + std::to_string(primitiveNumber) +
                           " not found");
  }

  return it->second(receiver, args, interpreter);
}

bool PrimitiveRegistry::hasPrimitive(int primitiveNumber) const {
  return primitives_.find(primitiveNumber) != primitives_.end();
}

std::vector<int> PrimitiveRegistry::getAllPrimitiveNumbers() const {
  std::vector<int> numbers;
  for (const auto &pair : primitives_) {
    numbers.push_back(pair.first);
  }
  return numbers;
}

void PrimitiveRegistry::clear() { primitives_.clear(); }

void PrimitiveRegistry::initializeCorePrimitives() {
  // Forward declarations for primitive registration functions
  extern void registerObjectPrimitives();
  extern void registerArrayPrimitives();
  extern void registerStringPrimitives();
  extern void registerExceptionPrimitives();
  void registerIntegerPrimitives();
  void registerBlockPrimitives();
  void registerDictionaryPrimitives();
  void registerSystemLoaderPrimitives();
  void registerCompilerPrimitives();

  // Register all core primitive groups
  registerObjectPrimitives();
  registerArrayPrimitives();
  registerStringPrimitives();
  registerExceptionPrimitives();

  // Register integer and block primitives directly here
  registerIntegerPrimitives();
  registerBlockPrimitives();
  registerDictionaryPrimitives();
  registerSystemLoaderPrimitives();
  registerCompilerPrimitives();
}

// Primitives namespace implementation
namespace Primitives {

void initialize() {
  PrimitiveRegistry::getInstance().initializeCorePrimitives();
}

TaggedValue callPrimitive(int primitiveNumber, TaggedValue receiver,
                          const std::vector<TaggedValue> &args,
                          Interpreter &interpreter) {
  return PrimitiveRegistry::getInstance().callPrimitive(
      primitiveNumber, receiver, args, interpreter);
}

void checkArgumentCount(const std::vector<TaggedValue> &args, size_t expected,
                        const std::string &primitiveName) {
  if (args.size() != expected) {
    throw PrimitiveFailure("Primitive " + primitiveName + " expects " +
                           std::to_string(expected) + " arguments, got " +
                           std::to_string(args.size()));
  }
}

void checkReceiverType(TaggedValue receiver, ObjectType expectedType,
                       const std::string &primitiveName) {
  if (!receiver.isPointer()) {
    throw PrimitiveFailure("Primitive " + primitiveName +
                           " expects object receiver");
  }

  Object *obj = receiver.asObject();
  if (obj->header.getType() != expectedType) {
    throw PrimitiveFailure("Primitive " + primitiveName +
                           " expects receiver of type " +
                           std::to_string(static_cast<int>(expectedType)));
  }
}

Class *ensureReceiverIsClass(TaggedValue receiver,
                             const std::string &primitiveName) {
  if (!receiver.isPointer()) {
    throw PrimitiveFailure("Primitive " + primitiveName +
                           " expects class receiver");
  }

  Object *obj = receiver.asObject();
  if (obj->header.getType() != ObjectType::CLASS) {
    throw PrimitiveFailure("Primitive " + primitiveName +
                           " expects class receiver");
  }

  return static_cast<Class *>(obj);
}

Object *ensureReceiverIsObject(TaggedValue receiver,
                               const std::string &primitiveName) {
  if (!receiver.isPointer()) {
    throw PrimitiveFailure("Primitive " + primitiveName +
                           " expects object receiver");
  }

  return receiver.asObject();
}

} // namespace Primitives

// Register integer primitives
void registerIntegerPrimitives() {
  auto &registry = PrimitiveRegistry::getInstance();

  // Register Integer arithmetic primitives using standard Smalltalk primitive
  // numbers
  registry.registerPrimitive(PrimitiveNumbers::SMALL_INT_ADD,
                             IntegerPrimitives::add);
  registry.registerPrimitive(PrimitiveNumbers::SMALL_INT_SUB,
                             IntegerPrimitives::subtract);
  registry.registerPrimitive(PrimitiveNumbers::SMALL_INT_MUL,
                             IntegerPrimitives::multiply);
  registry.registerPrimitive(PrimitiveNumbers::SMALL_INT_DIV,
                             IntegerPrimitives::divide);

  // Register Integer comparison primitives
  registry.registerPrimitive(PrimitiveNumbers::SMALL_INT_LT,
                             IntegerPrimitives::lessThan);
  registry.registerPrimitive(PrimitiveNumbers::SMALL_INT_GT,
                             IntegerPrimitives::greaterThan);
  registry.registerPrimitive(PrimitiveNumbers::SMALL_INT_EQ,
                             IntegerPrimitives::equal);
  registry.registerPrimitive(PrimitiveNumbers::SMALL_INT_NE,
                             IntegerPrimitives::notEqual);
  registry.registerPrimitive(PrimitiveNumbers::SMALL_INT_LE,
                             IntegerPrimitives::lessThanOrEqual);
  registry.registerPrimitive(PrimitiveNumbers::SMALL_INT_GE,
                             IntegerPrimitives::greaterThanOrEqual);
}

// Register block primitives
void registerBlockPrimitives() {
  auto &registry = PrimitiveRegistry::getInstance();
  registry.registerPrimitive(PrimitiveNumbers::BLOCK_VALUE,
                             BlockPrimitives::value);
  registry.registerPrimitive(PrimitiveNumbers::BLOCK_VALUE_ARG,
                             BlockPrimitives::value);
}

} // namespace smalltalk

// ---- SystemLoader primitives ----
namespace smalltalk {
namespace SystemLoaderPrimitives {
// Primitive 5000: SystemLoader>>start:
// Minimal bootstrap hook; currently returns true to signal success.
TaggedValue start(TaggedValue receiver, const std::vector<TaggedValue> &args,
                  Interpreter &interpreter) {
  (void)receiver;
  (void)interpreter;
  Primitives::checkArgumentCount(args, 1, "SystemLoader>>start:");
  return TaggedValue::trueValue();
}
} // namespace SystemLoaderPrimitives

void registerSystemLoaderPrimitives() {
  PrimitiveRegistry &registry = PrimitiveRegistry::getInstance();
  registry.registerPrimitive(PrimitiveNumbers::SYSTEM_LOADER_START,
                             SystemLoaderPrimitives::start);
}
} // namespace smalltalk

// ---- Compiler bridge primitives ----
namespace smalltalk {
namespace CompilerPrimitives {

// Primitive 5100: Compiler>>compile:in:
// Uses the C++ MethodCompiler to compile and install a method from source into
// a class
TaggedValue compileIn(TaggedValue receiver,
                      const std::vector<TaggedValue> &args,
                      Interpreter &interpreter) {
  (void)receiver;
  (void)interpreter;
  Primitives::checkArgumentCount(args, 2, "Compiler>>compile:in:");

  // arg0: String source
  if (!StringUtils::isString(args[0])) {
    throw PrimitiveFailure(
        "First argument must be a String containing method source");
  }
  String *str = StringUtils::asString(args[0]);
  const std::string &source = str->getContent();

  // arg1: Class
  if (!args[1].isPointer() ||
      args[1].asObject()->header.getType() != ObjectType::CLASS) {
    throw PrimitiveFailure("Second argument must be a Class");
  }
  Class *clazz = static_cast<Class *>(args[1].asObject());

  // Compile and install with write-through to Smalltalk dictionary
  MethodCompiler::addSmalltalkMethod(clazz, source, interpreter.getMemoryManager());

  return TaggedValue::trueValue();
}
} // namespace CompilerPrimitives

void registerCompilerPrimitives() {
  PrimitiveRegistry &registry = PrimitiveRegistry::getInstance();
  registry.registerPrimitive(PrimitiveNumbers::COMPILER_COMPILE_IN,
                             CompilerPrimitives::compileIn);
}
} // namespace smalltalk
