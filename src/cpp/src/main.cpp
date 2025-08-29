#include "bytecode.h"
#include "compiled_method.h"
#include "interpreter.h"
#include "memory_manager.h"
#include "method_compiler.h"
#include "primitives.h"
#include "primitives/block.h"
#include "simple_compiler.h"
#include "simple_parser.h"
#include "smalltalk_class.h"
#include "smalltalk_image.h"
#include "smalltalk_string.h"

#include <iomanip>
#include <iostream>
#include <string>

using namespace smalltalk;

struct Options {
  bool showParseTree = false;
  bool showBytecode = false;
  bool showMethod = false;
  bool runExpression = true;
  bool runLoaderOnImage = true;
  std::string imagePath;
  std::string expression;
};

void printUsage() {
  std::cout << "Usage:" << '\n';
  std::cout << "  smalltalk-vm [options] [<expression>]" << '\n';
  std::cout << '\n';
  std::cout << "Options:" << '\n';
  std::cout << "  --parse-tree     Show the parsed AST" << '\n';
  std::cout << "  --bytecode       Show detailed bytecode analysis" << '\n';
  std::cout << "  --method         Show compiled method details" << '\n';
  std::cout << "  --no-run         Don't execute the expression" << '\n';
  std::cout << "  --image <file>   Load a Smalltalk image file and bootstrap"
            << '\n';
  std::cout
      << "  --no-loader      When used with --image, skip SystemLoader start:"
      << '\n';
  std::cout << "  --help, -h       Show this help message" << '\n';
  std::cout << '\n';
  std::cout << "Examples:" << '\n';
  std::cout << "  smalltalk-vm \"42\"" << '\n';
  std::cout << "  smalltalk-vm --parse-tree \"3 + 4\"" << '\n';
  std::cout << "  smalltalk-vm --bytecode --method \"(10 - 2) * 3\"" << '\n';
  std::cout << "  smalltalk-vm --parse-tree --no-run \"ensure: aBlock | result "
               "| result := self value\""
            << '\n';
  std::cout << "  smalltalk-vm --image build/core.image" << '\n';
}

Options parseArguments(int argc, char **argv) {
  Options opts;

  for (int i = 1; i < argc; i++) {
    std::string arg = argv[i];

    if (arg == "--parse-tree") {
      opts.showParseTree = true;
    } else if (arg == "--bytecode") {
      opts.showBytecode = true;
    } else if (arg == "--method") {
      opts.showMethod = true;
    } else if (arg == "--no-run") {
      opts.runExpression = false;
    } else if (arg == "--no-loader") {
      opts.runLoaderOnImage = false;
    } else if (arg == "--image") {
      if (i + 1 >= argc) {
        std::cerr << "--image requires a file argument" << '\n';
        printUsage();
        exit(1);
      }
      opts.imagePath = argv[++i];
    } else if (arg == "--help" || arg == "-h") {
      printUsage();
      exit(0);
    } else if (arg.substr(0, 2) == "--") {
      std::cerr << "Unknown option: " << arg << '\n';
      printUsage();
      exit(1);
    } else {
      if (!opts.expression.empty()) {
        std::cerr
            << "Multiple expressions provided. Only one expression allowed."
            << '\n';
        printUsage();
        exit(1);
      }
      opts.expression = arg;
    }
  }
  // Allow running with only --image
  if (opts.expression.empty() && opts.imagePath.empty()) {
    std::cerr << "No expression or --image provided." << '\n';
    printUsage();
    exit(1);
  }

  return opts;
}

// Helper function to decode and print bytecode instructions
void printBytecodeAnalysis(const std::vector<uint8_t> &bytecodes) {
  std::cout << "\n=== Bytecode Analysis ===" << std::endl;

  // Print raw bytecode
  std::cout << "Raw bytecode (" << bytecodes.size() << " bytes): ";
  for (size_t i = 0; i < bytecodes.size(); i++) {
    std::cout << std::hex << std::setw(2) << std::setfill('0')
              << static_cast<int>(bytecodes[i]);
    if (i < bytecodes.size() - 1)
      std::cout << " ";
  }
  std::cout << std::dec << std::endl;

  // Decode instructions
  std::cout << "\nDecoded instructions:" << std::endl;
  for (size_t i = 0; i < bytecodes.size();) {
    uint8_t opcode = bytecodes[i];
    std::cout << "  " << std::setw(3) << i << ": ";

    switch (static_cast<Bytecode>(opcode)) {
    case Bytecode::PUSH_LITERAL:
      std::cout << "PUSH_LITERAL ";
      if (i + 4 < bytecodes.size()) {
        uint32_t index = bytecodes[i + 1] | (bytecodes[i + 2] << 8) |
                         (bytecodes[i + 3] << 16) | (bytecodes[i + 4] << 24);
        std::cout << index;
        i += 5;
      } else {
        std::cout << "(incomplete)";
        i++;
      }
      break;

    case Bytecode::PUSH_SELF:
      std::cout << "PUSH_SELF";
      i++;
      break;

    case Bytecode::PUSH_TEMPORARY_VARIABLE:
      std::cout << "PUSH_TEMPORARY_VARIABLE ";
      if (i + 4 < bytecodes.size()) {
        uint32_t index = bytecodes[i + 1] | (bytecodes[i + 2] << 8) |
                         (bytecodes[i + 3] << 16) | (bytecodes[i + 4] << 24);
        std::cout << index;
        i += 5;
      } else {
        std::cout << "(incomplete)";
        i++;
      }
      break;

    case Bytecode::STORE_TEMPORARY_VARIABLE:
      std::cout << "STORE_TEMPORARY_VARIABLE ";
      if (i + 4 < bytecodes.size()) {
        uint32_t index = bytecodes[i + 1] | (bytecodes[i + 2] << 8) |
                         (bytecodes[i + 3] << 16) | (bytecodes[i + 4] << 24);
        std::cout << index;
        i += 5;
      } else {
        std::cout << "(incomplete)";
        i++;
      }
      break;

    case Bytecode::SEND_MESSAGE:
      std::cout << "SEND_MESSAGE ";
      if (i + 8 < bytecodes.size()) {
        uint32_t selector = bytecodes[i + 1] | (bytecodes[i + 2] << 8) |
                            (bytecodes[i + 3] << 16) | (bytecodes[i + 4] << 24);
        uint32_t argCount = bytecodes[i + 5] | (bytecodes[i + 6] << 8) |
                            (bytecodes[i + 7] << 16) | (bytecodes[i + 8] << 24);
        std::cout << "selector=" << selector << " args=" << argCount;
        i += 9;
      } else {
        std::cout << "(incomplete)";
        i++;
      }
      break;

    case Bytecode::RETURN_STACK_TOP:
      std::cout << "RETURN_STACK_TOP";
      i++;
      break;

    case Bytecode::POP:
      std::cout << "POP";
      i++;
      break;

    case Bytecode::CREATE_BLOCK:
      std::cout << "CREATE_BLOCK ";
      if (i + 8 < bytecodes.size()) {
        uint32_t methodIndex = bytecodes[i + 1] | (bytecodes[i + 2] << 8) |
                               (bytecodes[i + 3] << 16) |
                               (bytecodes[i + 4] << 24);
        uint32_t paramCount = bytecodes[i + 5] | (bytecodes[i + 6] << 8) |
                              (bytecodes[i + 7] << 16) |
                              (bytecodes[i + 8] << 24);
        std::cout << "method=" << methodIndex << " params=" << paramCount;
        i += 9;
      } else {
        std::cout << "(incomplete)";
        i++;
      }
      break;

    case Bytecode::DUPLICATE:
      std::cout << "DUPLICATE";
      i++;
      break;

    default:
      std::cout << "UNKNOWN(" << static_cast<int>(opcode) << ")";
      i++;
      break;
    }
    std::cout << std::endl;
  }
}

int main(int argc, char **argv) {
  if (argc < 2) {
    printUsage();
    return 1;
  }

  Options opts = parseArguments(argc, argv);

  try {
    // Step 1: Initialize class system and primitives
    ClassUtils::initializeCoreClasses();
    auto &primitiveRegistry = PrimitiveRegistry::getInstance();
    primitiveRegistry.initializeCorePrimitives();
    // Add integer primitives (legacy path maintains consistency)
    Class *integerClass = ClassUtils::getIntegerClass();
    IntegerClassSetup::addPrimitiveMethods(integerClass);
    primitiveRegistry.registerPrimitive(PrimitiveNumbers::BLOCK_VALUE,
                                        BlockPrimitives::value);

    // If an image is provided, load it and optionally call the loader
    std::unique_ptr<SmalltalkImage> ownedImage;
    SmalltalkImage *imagePtr = nullptr;
    if (!opts.imagePath.empty()) {
      auto &manager = ImageManager::getInstance();
      if (!manager.loadImageFromFile(opts.imagePath)) {
        std::cerr << "Failed to load image: " << opts.imagePath << "\n";
        return 1;
      }
      imagePtr = manager.getCurrentImage();
      if (!imagePtr) {
        std::cerr << "Image failed to initialize after load" << '\n';
        return 1;
      }

      if (opts.runLoaderOnImage) {
        // Call minimal bootstrap entry
        TaggedValue started =
            imagePtr->evaluate("SystemLoader new start: 'cli'");
        (void)started; // We don't require a true result to proceed
      }
    }

    // If expression provided, evaluate it (within image if loaded)
    if (!opts.expression.empty()) {
      if (imagePtr) {
        TaggedValue result = imagePtr->evaluate(opts.expression);
        std::cout << "\n=== Result ===" << std::endl;
        if (StringUtils::isString(result)) {
          String *str = StringUtils::asString(result);
          std::cout << str->toString() << std::endl;
        } else {
          std::cout << result << std::endl;
        }
      } else {
        // Legacy compile-and-execute path without an image
        SimpleParser parser(opts.expression);
        auto methodAST = parser.parseMethod();
        if (opts.showParseTree) {
          std::cout << "\n=== Parse Tree ===" << std::endl;
          std::cout << methodAST->toString() << std::endl;
        }
        SimpleCompiler compiler;
        auto compiledMethod = compiler.compile(*methodAST);
        if (opts.showMethod) {
          std::cout << "\n=== Compiled Method ===" << std::endl;
          std::cout << "Primitive number: " << compiledMethod->primitiveNumber
                    << std::endl;
        }
        if (opts.showBytecode) {
          printBytecodeAnalysis(compiledMethod->getBytecodes());
        }
        if (opts.runExpression) {
          MemoryManager memoryManager;
          SmalltalkImage image;
          Interpreter interpreter(memoryManager, image);
          primitiveRegistry.registerPrimitive(PrimitiveNumbers::BLOCK_VALUE,
                                              BlockPrimitives::value);
          TaggedValue result =
              interpreter.executeCompiledMethod(*compiledMethod);
          std::cout << "\n=== Result ===" << std::endl;
          if (StringUtils::isString(result)) {
            String *str = StringUtils::asString(result);
            std::cout << str->toString() << std::endl;
          } else {
            std::cout << result << std::endl;
          }
        }
      }
    }

  } catch (const std::exception &e) {
    std::cerr << "Error: " << e.what() << '\n';
    return 1;
  }

  return 0;
}
