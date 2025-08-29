#include "smalltalk_exception.h"
#include "context.h"
#include "interpreter.h"
#include <sstream>
#include <stdexcept>

namespace smalltalk {

SmalltalkException::SmalltalkException(const std::string &message,
                                       const std::string &exceptionClass)
    : Object(ObjectType::OBJECT, 0), message_(message),
      exceptionClass_(exceptionClass) {
  // Object header is initialized by base constructor
}

std::string SmalltalkException::toString() const {
  return exceptionClass_ + ": " + message_;
}

void SmalltalkException::captureStackTrace(Interpreter &interpreter) {
  stackTrace_.clear();

  // Get current context chain
  MethodContext *context = interpreter.getCurrentContext();
  int frameNumber = 0;

  while (context && frameNumber < 10) { // Limit to 10 frames
    std::stringstream ss;
    ss << "Frame " << frameNumber << ": ";

    // Try to get method information
    if (context->header.hash != 0) {
      ss << "method hash=" << context->header.hash;
    } else {
      ss << "unknown method";
    }

    // Add instruction pointer info
    ss << " ip=" << context->instructionPointer;

    stackTrace_.push_back(ss.str());

    // Move to sender context
    if (context->sender.isPointer()) {
      context = static_cast<MethodContext *>(context->sender.asObject());
    } else {
      break;
    }
    frameNumber++;
  }
}

void SmalltalkException::signal() {
  // Convert to C++ exception for now
  throw std::runtime_error(toString());
}

// ExceptionHandler implementation

std::unique_ptr<SmalltalkException>
ExceptionHandler::fromStdException(const std::exception &e) {
  std::string message = e.what();

  // Map common C++ exceptions to Smalltalk exceptions
  if (message.find("Division by zero") != std::string::npos ||
      message.find("ZeroDivisionError") != std::string::npos) {
    return std::make_unique<ZeroDivisionError>(message);
  } else if (message.find("Undefined variable") != std::string::npos ||
             message.find("NameError") != std::string::npos) {
    // Extract variable name
    size_t pos = message.find(": ");
    if (pos != std::string::npos) {
      std::string varName = message.substr(pos + 2);
      return std::make_unique<NameError>(varName);
    }
    return std::make_unique<NameError>("unknown");
  } else if (message.find("Index") != std::string::npos ||
             message.find("bounds") != std::string::npos) {
    return std::make_unique<IndexError>(message);
  } else if (message.find("does not understand") != std::string::npos) {
    return std::make_unique<MessageNotUnderstood>("Object", "unknownMethod");
  } else if (message.find("Invalid argument") != std::string::npos) {
    return std::make_unique<ArgumentError>(message);
  } else {
    return std::make_unique<RuntimeError>(message);
  }
}

TaggedValue ExceptionHandler::handleException(SmalltalkException &exception,
                                              Interpreter &interpreter) {
  // Capture stack trace
  exception.captureStackTrace(interpreter);

  // For now, just convert back to C++ exception with exception class name
  // In a full implementation, this would handle Smalltalk exception semantics
  // like searching for exception handlers, unwinding contexts, etc.

  std::string errorMessage = exception.getExceptionClass();

  // Return the exception class name as the error for testing
  // This allows the test to see the proper exception type
  throw std::runtime_error(errorMessage);
}

void ExceptionHandler::throwException(
    std::unique_ptr<SmalltalkException> exception) {
  // Convert to C++ exception with proper type information
  std::string errorMessage = exception->getExceptionClass();
  throw std::runtime_error(errorMessage);
}

bool ExceptionHandler::shouldCatch(const SmalltalkException &exception,
                                   const std::string &handlerClass) {
  // Simple class hierarchy check
  std::string exceptionClass = exception.getExceptionClass();

  if (exceptionClass == handlerClass) {
    return true;
  }

  // Check inheritance hierarchy
  if (handlerClass == "Exception") {
    return true; // Exception is the root of all exceptions
  }

  return false;
}

} // namespace smalltalk